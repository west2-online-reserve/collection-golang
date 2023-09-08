package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	// MySQL信息
	Username     = "root"           //username 		like 	`root`
	Password     = "357920"         //password 		like 	`123456`
	Hostname     = "127.0.0.1:3306" //hostname 		like 	`127.0.0.1:3306`
	Databasename = "test"           //databasename 	like 	`databasename`
)

const (
	//b站获取评论信息的网页
	//第一次发现的网页是"https://api.bilibili.com/x/v2/reply/wbi/"
	//但突然发现删除其尾部的部分不必要信息后会出现403，而不删除则不会有这个结果

	//eg.https://api.bilibili.com/x/v2/reply/wbi/main?oid=317458729&type=1&mode=3&pagination_str=%7B%22offset%22:%22%7B%5C%22type%5C%22:1,%5C%22direction%5C%22:1,%5C%22session_id%5C%22:%5C%221734794499468244%5C%22,%5C%22data%5C%22:%7B%7D%7D%22%7D&plat=1&web_location=1315875&w_rid=bf59645ab210aadd4a7d6c1373bb5440&wts=1694135644
	//->.https://api.bilibili.com/x/v2/reply/wbi/main?oid=317458729&type=1&mode=3
	//这一操作会出现403

	//神奇的是，将wbi/路径去掉后，居然可以在删掉不必要的信息下访问了
	//->.https://api.bilibili.com/x/v2/reply/main?oid=317458729&type=1&mode=3
	//这一操作会获得评论区的信息

	BilibiliReplyApi = "https://api.bilibili.com/x/v2/reply/"
	Bv               = "BV12341117rG"
	UserCookie       = ``
)

// 参考`https://mholt.github.io/json-to-go/`生成
// Json2Go
// 主评论数据结构
type RepliesMainPage struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	TTL     int    `json:"ttl"`
	Data    struct {
		Cursor struct {
			IsEnd bool `json:"is_end"`
		} `json:"cursor"`
		Replies []struct {
			Rpid   uint64 `json:"rpid"`
			Root   uint64 `json:"root"`   //主评论的root,parent值均为0
			Parent uint64 `json:"parent"` //但我想要让其与子评论数据对称，方便存储
			Like   uint   `json:"like"`
			Ctime  int64  `json:"ctime"`
			Member struct {
				UserName  string `json:"uname"`
				Mid       string `json:"mid"`
				Sex       string `json:"sex"`
				Sign      string `json:"sign"`
				LevelInfo struct {
					CurrentLevel int `json:"current_level"`
				} `json:"level_info"`
				Vip struct {
					VipType int `json:"vipType"`
				} `json:"vip"`
			} `json:"member"`
			Content struct {
				Message string `json:"message"`
			} `json:"content"`
			ChildReplies []struct{} `json:"replies"`
		} `json:"replies"`
	}
}

// 子评论数据结构
type RepliesReplyPage struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	TTL     int    `json:"ttl"`
	Data    struct {
		Page struct {
			Num   int `json:"num"`
			Size  int `json:"size"`
			Count int `json:"count"`
		} `json:"page"`
		Replies []struct {
			Rpid   uint64 `json:"rpid"`
			Root   uint64 `json:"root"`
			Parent uint64 `json:"parent"`
			Like   uint   `json:"like"`
			Ctime  int64  `json:"ctime"`
			Member struct {
				UserName  string `json:"uname"`
				Mid       string `json:"mid"`
				Sex       string `json:"sex"`
				Sign      string `json:"sign"`
				LevelInfo struct {
					CurrentLevel int `json:"current_level"`
				} `json:"level_info"`
				Vip struct {
					VipType int `json:"vipType"`
				} `json:"vip"`
			} `json:"member"`
			Content struct {
				Message string `json:"message"`
			} `json:"content"`
		} `json:"replies"`
	} `json:"data"`
}

// Reply结构，用于存入数据库
type Reply struct {
	replyId      uint64
	replyRoot    uint64
	replyParent  uint64
	replyLike    uint
	replyTime    string
	replyName    string
	replyUserId  uint64
	replySex     string
	replySign    string
	replyLevel   int
	replyVip     int
	replyContent string
}

// GoLimit(协程限制)
// ----------------------------------------------------
type GoLimit struct {
	routine chan int
}

func NewGoLimit(maxRoutine int) *GoLimit {
	return &GoLimit{make(chan int, maxRoutine)}
}

func (goLimit *GoLimit) Add() {
	goLimit.routine <- 1
}

func (goLimit *GoLimit) Done() {
	<-goLimit.routine
}

//----------------------------------------------------

// 协程循环体统一结束标志
var DataCatchOver = false
var PageReachEnd = false

// 返回DSN
func DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", Username, Password, Hostname, Databasename)
}

// 连接到sql数据库
func ConnectDataBase(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Printf("[MySQL]%s\n", err.Error())
	}
}

// 初始化数据表
func IniDataStruct(db *sql.DB) error {
	stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS `" + Bv + "`(`ReplyId` BIGINT UNSIGNED NOT NULL,`ReplyRoot` BIGINT UNSIGNED NOT NULL,`ReplyParent` BIGINT UNSIGNED NOT NULL,`ReplyLike` INT UNSIGNED NOT NULL,`ReplyTime` DATETIME,`Name` VARCHAR(32) NOT NULL,`Id`   BIGINT UNSIGNED NOT NULL,`Sex`  VARCHAR(4),`Sign` TEXT,`Level` INT,`VIP` INT,`ReplyContent` MEDIUMTEXT NOT NULL, PRIMARY KEY ( `ReplyId` ))ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;")
	if err != nil {
		log.Printf("[MySQL]%s\n", err.Error())
		return errors.New("failed to ini")
	}
	_, err = stmt.Exec()
	if err != nil {
		log.Printf("[MySQL]%s\n", err.Error())
		return errors.New("failed to ini")
	}
	log.Print("[MySQL]Info:DataStruct initialized successfully\n")
	return nil
}

// 插入数据库
func InsertToDatabase(db *sql.DB, reply Reply) {
	const insertCmd = "insert into " + Bv + "(ReplyId,ReplyRoot,ReplyParent,ReplyLike,ReplyTime,Name,Id,Sex,Sign,Level,Vip,ReplyContent)values (?,?,?,?,?,?,?,?,?,?,?,?)"
	result, err := db.Exec(insertCmd,
		reply.replyId,
		reply.replyRoot,
		reply.replyParent,
		reply.replyLike,
		reply.replyTime,
		reply.replyName,
		reply.replyUserId,
		reply.replySex,
		reply.replySign,
		reply.replyLevel,
		reply.replyVip,
		reply.replyContent)
	if err != nil {
		log.Printf("[MySQL]Insert failed,%s\n", err)
		return
	}
	_, err = result.LastInsertId()
	if err != nil {
		log.Printf("[MySQL]Insert failed,%s\n", err)
		return
	}
	log.Printf("[MySQL]Info:The [%v]Reply inserted successfully\n", reply.replyId)
}

// 获取Response
func GetResp(url string) *http.Response {
	//User-Agent----Map
	userAgentMap := []string{
		`Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.95 Safari/537.36 OPR/26.0.1656.60`,
		`Opera/8.0 (Windows NT 5.1; U; en)`,
		`Mozilla/5.0 (Windows NT 5.1; U; en; rv:1.8.1) Gecko/20061208 Firefox/2.0.0 Opera 9.50`,
		`Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1; en) Opera 9.50`,
		`Opera/9.80 (Macintosh; Intel Mac OS X 10.6.8; U; en) Presto/2.8.131 Version/11.11`,
		`Opera/9.80 (Windows NT 6.1; U; en) Presto/2.8.131 Version/11.11`,
		`Opera/9.80 (Android 2.3.4; Linux; Opera Mobi/build-1107180945; U; en-GB) Presto/2.8.149 Version/11.10`,
		`Mozilla/5.0 (Windows NT 6.1; WOW64; rv:34.0) Gecko/20100101 Firefox/34.0`,
		`Mozilla/5.0 (X11; U; Linux x86_64; zh-CN; rv:1.9.2.10) Gecko/20100922 Ubuntu/10.10 (maverick) Firefox/3.6.10`,
		`Mozilla/5.0 (Macintosh; Intel Mac OS X 10.6; rv,2.0.1) Gecko/20100101 Firefox/4.0.1`,
		`Mozilla/5.0 (Windows NT 6.1; rv,2.0.1) Gecko/20100101 Firefox/4.0.1`,
		`Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.71 Safari/537.36`,
		`Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.11 (KHTML, like Gecko) Chrome/23.0.1271.64 Safari/537.11`,
		`Mozilla/5.0 (Windows; U; Windows NT 6.1; en-US) AppleWebKit/534.16 (KHTML, like Gecko) Chrome/10.0.648.133 Safari/534.16`,
		`Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_0) AppleWebKit/535.11 (KHTML, like Gecko) Chrome/17.0.963.56 Safari/535.11`,
		`Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/30.0.1599.101 Safari/537.36`,
		`Mozilla/5.0 (Windows NT 6.1; WOW64; Trident/7.0; rv:11.0) like Gecko`,
		`Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; 360SE)`,
		`Mozilla/5.0 (Windows NT 5.1) AppleWebKit/535.11 (KHTML, like Gecko) Chrome/17.0.963.84 Safari/535.11 SE 2.X MetaSr 1.0`,
		`Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; Trident/4.0; SV1; QQDownload 732; .NET4.0C; .NET4.0E; SE 2.X MetaSr 1.0)`,
		`Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; Trident/4.0; SE 2.X MetaSr 1.0; SE 2.X MetaSr 1.0; .NET CLR 2.0.50727; SE 2.X MetaSr 1.0)`,
		`Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/38.0.2125.122 UBrowser/4.0.3214.0 Safari/537.36`,
		`Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.87 UBrowser/6.2.4094.1 Safari/537.36`,
		`Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; WOW64; Trident/5.0; SLCC2; .NET CLR 2.0.50727; .NET CLR 3.5.30729; .NET CLR 3.0.30729;`,
		`Media Center PC 6.0; .NET4.0C; .NET4.0E; QQBrowser/7.0.3698.400)`,
		`Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1; SV1; QQDownload 732; .NET4.0C; .NET4.0E)`}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Print("[GetResp]Error:", err, "\n")
	}
	request.Header.Set("Cookie", UserCookie)
	request.Header.Set("User-Agent", userAgentMap[rand.Intn(len(userAgentMap))])
	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(request)
	if err != nil {
		log.Printf("[GetResp]Error:%s\n", err)
	}
	return resp
}

// 参考知乎问答中https://www.zhihu.com/question/381784377
// mcfx的思路
// 获取oid字符串
func BvDecode(bv string) string {
	const (
		table = `fZodR9XQDSUm21yCkr6zBqiveYah8bt4xsWpHnJE7jL5VG3guMTKNPAwcF`
		xor   = 177451812
		add   = 8728348608
	)
	s := [...]int64{11, 10, 3, 8, 4, 6}
	r := int64(0)
	var tr [128]int64
	for i := int64(0); i < 58; i++ {
		tr[table[i]] = i
	}
	for i := 0; i < 6; i++ {
		r += (int64)(math.Pow(58, float64(i)) * float64(tr[bv[s[i]]]))
	}
	return strconv.FormatInt((r-add)^xor, 10)
}

// 组装主评论Url
func GetBvCommentWeb(bili_bv, bili_mode, bili_type, bili_next string) string {
	return BilibiliReplyApi + "main?oid=" + BvDecode(bili_bv) + "&type=" + bili_type + "&mode=" + bili_mode + "&next=" + bili_next
}

// 组装子评论Url
func GetChildCommentWeb(bili_bv, root, bili_type, bili_ps, bili_pn string) string {
	return BilibiliReplyApi + "reply?oid=" + BvDecode(bili_bv) + "&root=" + root + "&type=" + bili_type + "&ps=" + bili_ps + "&pn=" + bili_pn
}

// 获取主评论的json数据
func GetMainRepliesData(resp *http.Response, repliesPageList chan RepliesMainPage) {
	defer resp.Body.Close()
	var repliesMainPage RepliesMainPage
	info, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[GetMainRepliesData]Error:%s\n", err)
	}
	err = json.Unmarshal(info, &repliesMainPage)
	if err != nil {
		log.Printf("[GetMainRepliesData]Error:%s\n", err)
	}
	if repliesMainPage.Data.Cursor.IsEnd {
		PageReachEnd = true
	}
	repliesPageList <- repliesMainPage
}

// 获取子评论的json数据
func GetChildRepliesData(resp *http.Response) RepliesReplyPage {
	defer resp.Body.Close()
	var childRepliesPage RepliesReplyPage
	info, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[GetChildRepliesData]Error:%s\n", err)
	}
	err = json.Unmarshal(info, &childRepliesPage)
	if err != nil {
		log.Printf("[GetChildRepliesData]Error:%s\n", err)
	}
	return childRepliesPage
}

// 获取主评论的Reply数据结构
func GetMainReply(repliesMainPage *RepliesMainPage, replies chan Reply) {
	for i := 0; i < len(repliesMainPage.Data.Replies); i++ {
		//母评论
		userId, _ := strconv.Atoi(repliesMainPage.Data.Replies[i].Member.Mid)
		replies <- Reply{
			repliesMainPage.Data.Replies[i].Rpid,
			repliesMainPage.Data.Replies[i].Root,
			repliesMainPage.Data.Replies[i].Parent,
			repliesMainPage.Data.Replies[i].Like,
			time.Unix(repliesMainPage.Data.Replies[i].Ctime, 0).Format("2006-01-02 15:04:05"),
			repliesMainPage.Data.Replies[i].Member.UserName,
			uint64(userId),
			repliesMainPage.Data.Replies[i].Member.Sex,
			repliesMainPage.Data.Replies[i].Member.Sign,
			repliesMainPage.Data.Replies[i].Member.LevelInfo.CurrentLevel,
			repliesMainPage.Data.Replies[i].Member.Vip.VipType,
			repliesMainPage.Data.Replies[i].Content.Message}
	}
}

// 获取子评论的Reply数据结构
func GetReplyReply(repliesReplyPage *RepliesReplyPage, replies chan Reply) {
	for i := 0; i < len(repliesReplyPage.Data.Replies); i++ {
		//子评论
		userId, _ := strconv.Atoi(repliesReplyPage.Data.Replies[i].Member.Mid)
		replies <- Reply{
			repliesReplyPage.Data.Replies[i].Rpid,
			repliesReplyPage.Data.Replies[i].Root,
			repliesReplyPage.Data.Replies[i].Parent,
			repliesReplyPage.Data.Replies[i].Like,
			time.Unix(repliesReplyPage.Data.Replies[i].Ctime, 0).Format("2006-01-02 15:04:05"),
			repliesReplyPage.Data.Replies[i].Member.UserName,
			uint64(userId),
			repliesReplyPage.Data.Replies[i].Member.Sex,
			repliesReplyPage.Data.Replies[i].Member.Sign,
			repliesReplyPage.Data.Replies[i].Member.LevelInfo.CurrentLevel,
			repliesReplyPage.Data.Replies[i].Member.Vip.VipType,
			repliesReplyPage.Data.Replies[i].Content.Message}
	}
}

func main() {
	//打开数据库
	db, err := sql.Open("mysql", DSN())
	if err != nil {
		log.Printf("[MySQL]Error:%s", err)
		return
	}

	//初始化
	err = IniDataStruct(db)
	if err != nil {
		log.Printf("[MySQL]Error:%s", err)
		return
	}

	next := 1
	repliesMainPageDataList := make(chan RepliesMainPage, 32)
	repliesReplyPageRootList := make(chan uint64, 32)
	repliesInfoDataList := make(chan Reply, 64)

	pageReplyCatchRoutine := NewGoLimit(32)
	repliesMainCatchRoutine := NewGoLimit(32)
	repliesMainInsertRoutine := NewGoLimit(32)

	//协程池2:处理Main页面评论信息整体,并传输子评论链接
	go func() {
		for !DataCatchOver {
			repliesMainCatchRoutine.Add()
			go func() {
				defer repliesMainCatchRoutine.Done()
				if repliesMainPageData, ok := <-repliesMainPageDataList; ok {
					for i := 0; i < len(repliesMainPageData.Data.Replies); i++ {
						repliesReplyPageRootList <- repliesMainPageData.Data.Replies[i].Rpid
					}
					GetMainReply(&repliesMainPageData, repliesInfoDataList)
				}
			}()
		}
	}()

	//协程池3:处理Reply页面的评论信息整体
	go func() {
		for !DataCatchOver {
			pageReplyCatchRoutine.Add()
			go func() {
				defer pageReplyCatchRoutine.Done()
				if rootId, ok := <-repliesReplyPageRootList; ok {
					for i := 1; ; i++ {
						childReplyUrl := GetChildCommentWeb(Bv, strconv.FormatUint(rootId, 10), "1", "10", strconv.Itoa(i))
						resp := GetResp(childReplyUrl)
						childRepliesData := GetChildRepliesData(resp)
						GetReplyReply(&childRepliesData, repliesInfoDataList)
						if childRepliesData.Data.Page.Num*childRepliesData.Data.Page.Size >= childRepliesData.Data.Page.Count {
							break
						}
					}
				}
			}()
		}
	}()

	//协程池4:插入到mySQL中
	go func() {
		for !DataCatchOver {
			repliesMainInsertRoutine.Add()
			go func() {
				defer repliesMainInsertRoutine.Done()
				if reply, ok := <-repliesInfoDataList; ok {
					InsertToDatabase(db, reply)
				}
			}()
		}
	}()

	//协程2:控制运行状态
	go func() {
		for !DataCatchOver {
			if len(repliesInfoDataList) == 0 && len(repliesMainPageDataList) == 0 && len(repliesReplyPageRootList) == 0 && PageReachEnd {
				DataCatchOver = true
				close(repliesMainPageDataList)
				close(repliesInfoDataList)
				close(repliesReplyPageRootList)
			}
			time.Sleep(5 * time.Second)
		}
	}()

	//主协程:翻页并抓取Main页面评论信息整体
	for !PageReachEnd {
		resp := GetResp(GetBvCommentWeb(Bv, "3", "1", strconv.Itoa(next)))
		next++
		GetMainRepliesData(resp, repliesMainPageDataList)
		time.Sleep(time.Duration(rand.Intn(500)+500) * (time.Millisecond))
	}

	//主协程:等待数据抓取结束(每5s检测一次)
	for !DataCatchOver {
		time.Sleep(5 * time.Second)
	}

	log.Print("[Main]Info:Catch Over")
	db.Close()
}
