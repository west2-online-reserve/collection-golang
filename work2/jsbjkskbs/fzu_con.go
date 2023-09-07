package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	_ "github.com/go-sql-driver/mysql"
)

// GoLimit(协程限制)
// ----------------------------------------------------
type GoLimit struct {
	cnt chan int
}

func NewGoLimit(maxRoutine int) *GoLimit {
	return &GoLimit{make(chan int, maxRoutine)}
}

func (goLimit *GoLimit) Add() {
	goLimit.cnt <- 1
}

func (goLimit *GoLimit) Done() {
	<-goLimit.cnt
}

//----------------------------------------------------

const (
	// MySQL信息
	username     = "root"           //username 		like 	`root`
	password     = "357920"         //password 		like 	`123456`
	hostname     = "127.0.0.1:3306" //hostname 		like 	`127.0.0.1:3306`
	databasename = "test"           //databasename 	like 	`databasename`

	// FZU网站信息
	FzuInfoUrl        = "https://info22.fzu.edu.cn/"
	FzuNewsList       = "lm_list.jsp"
	FzuNewsListTreeId = "1460"

	// 通过分析 `function _showDynClicks(clicktype, owner, clickid)`后
	// 发现有段可以直接获得点击量的url `var url = '/system/resource/code/news/click/dynclicks.jsp?clickid='+clickid+'&owner='+owner+'&clicktype='+clicktype;`
	rGetClicksPath = "system/resource/code/news/click/dynclicks.jsp"

	// 跳过页数，可加快数据定位(后面想改成用大步长翻页,但懒起来了)
	SkipPage = 150
)

// 协程循环体统一结束标志
var DataCatchOver = false

// 爬虫日期限制
var TimeEnd, _ = time.Parse("2006-01-02", "2020-01-01")
var TimeBeg, _ = time.Parse("2006-01-02", "2021-09-01")

// 正则表达式初始化
var filterAlnum = regexp.MustCompile(`[0-9a-zA-z]+`)
var rDateRegex = regexp.MustCompile(`\d+\-\d+\-\d+`)
var rAuthorRegex = regexp.MustCompile(`\S+`)
var rClicksRegex = regexp.MustCompile(`"\w+".*\d+.*\d+`)

// 需求数据字段
type NewsInfo struct {
	title  string
	pDate  string
	author string
	text   string
	clicks string
}

// 返回DSN
func DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, databasename)
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
	stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS `fzu_info`(`index` INT UNSIGNED AUTO_INCREMENT,`title` VARCHAR(128) NOT NULL,`author` VARCHAR(16) NOT NULL,`date` DATE,`text` MEDIUMTEXT,`clicks` INT UNSIGNED,PRIMARY KEY ( `index` ))ENGINE=InnoDB DEFAULT CHARSET=utf8;")
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
func InsertToDatabase(db *sql.DB, newsInfo NewsInfo) {
	const insertCmd = "insert into fzu_info(title,author,date,text,clicks)values (?,?,?,?,?)"
	t, _ := time.Parse("2006-01-02", newsInfo.pDate)
	if t.After(TimeBeg) {
		log.Printf("[InsertToDatabase]Info:Insert Rejected(Due to Date:%s)\n", newsInfo.pDate)
		return
	} else if t.Before(TimeEnd) {
		log.Printf("[InsertToDatabase]Info:Insert Rejected(Due to Date:%s)\n", newsInfo.pDate)
		DataCatchOver = true
		return
	}
	result, err := db.Exec(insertCmd, newsInfo.title, newsInfo.author, newsInfo.pDate, newsInfo.text, newsInfo.clicks)
	if err != nil {
		log.Printf("[MySQL]Insert failed,%s\n", err)
		return
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("[MySQL]Insert failed,%s\n", err)
		return
	}
	log.Printf("[MySQL]Info:The [%v]field inserted successfully\n", id)
}

// 构造特定页链接
func FzuInfoPageLink(page int) string {
	return FzuInfoUrl + FzuNewsList + "?PAGENUM=" + strconv.Itoa(page) + "&wbtreeid=" + FzuNewsListTreeId
}

// 构造点击量链接
func FzuInfoClicksLink(clicktype, owner, clickid string) string {
	return FzuInfoUrl + rGetClicksPath + "?clickid=" + clickid + "&owner=" + owner + "&clicktype=" + clicktype
}

// 获取跳转链接
func FzuInfoChildUrl(parentResp *http.Response, childUrl chan string) {
	defer parentResp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(parentResp.Body)
	if err != nil {
		log.Printf("[FzuInfoChildUrl]Error:%s\n", err)
	}
	doc.Find("li.clearfloat").Each(func(index int, sel *goquery.Selection) {
		link, exist := sel.Find("a:nth-of-type(2)").Attr("href")
		if exist {
			childUrl <- link
		}
	})
}

// 获取点击量
func GetClicks(scriptOnHtml string) string {
	paramString := rClicksRegex.FindString(scriptOnHtml)
	params := strings.Split(paramString, ",")
	for i := 0; i < len(params); i++ {
		params[i] = filterAlnum.FindString(params[i])
	}
	resp, err := http.Get(FzuInfoClicksLink(params[0], params[1], params[2]))
	if err != nil {
		log.Printf("[GetClicks]Error:%s\n", err)
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Printf("[GetClicks]Error:%s\n", err)
	}
	return doc.Find("body").Text()
}

// 访问跳转链接并传入字段
func GetChildUrlInfo(childUrl string, newsInfoList chan NewsInfo) {
	resp, err := http.Get(FzuInfoUrl + childUrl)
	if err != nil {
		log.Printf("[GetChildUrlInfo]Error:%s\n", err)
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Printf("[GetChildUrlInfo]Error:%s\n", err)
	}
	var newsInfo NewsInfo
	newsInfo.author = rAuthorRegex.FindString(doc.Find("h3.fl").Text())
	newsInfo.title = doc.Find("div.conth1").Text()
	newsInfo.pDate = rDateRegex.FindString(doc.Find("div.conthsj").Text())
	newsInfo.clicks = GetClicks(doc.Find("div.conthsj").Find("script").Text())
	newsInfo.text = doc.Find("div.v_news_content").Find("Span").Text()
	if !DataCatchOver {
		newsInfoList <- newsInfo
	}
}

// 快速定位指定日期区间
func QuickLocatePage() (int, int) {
	pageBeg, pageEnd := SkipPage, SkipPage
	for findBeg := false; !findBeg; {
		resp, err := http.Get(FzuInfoPageLink(pageBeg))
		if err != nil {
			log.Print("[QuickLocatePage]Error:", err)
		}
		defer resp.Body.Close()
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			log.Print("[QuickLocatePage]Error:", err)
		}
		doc.Find("span.fr").Each(func(index int, sel *goquery.Selection) {
			t, _ := time.Parse("2006-01-02", sel.Text())
			if t.Before(TimeBeg) || t.Equal(TimeBeg) {
				pageEnd = pageBeg
				findBeg = true
			}
		})
		if !findBeg {
			pageBeg++
		}
	}
	for findEnd := false; !findEnd; {
		resp, err := http.Get(FzuInfoPageLink(pageEnd))
		if err != nil {
			log.Print("[QuickLocatePage]Error:", err)
		}
		defer resp.Body.Close()
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			log.Print("[QuickLocatePage]Error:", err)
		}
		doc.Find("span.fr").Each(func(index int, sel *goquery.Selection) {
			t, _ := time.Parse("2006-01-02", sel.Text())
			if t.Before(TimeEnd) || t.Equal(TimeEnd) {
				findEnd = true
			}
		})
		if !findEnd {
			pageEnd++
		}
	}
	return pageBeg, pageEnd
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

	log.Print("[Main]Info:正在寻找指定时间区间页...\n")
	pageBeg, pageEnd := QuickLocatePage()
	log.Print("[Main]Info:PageBeg-", pageBeg, " pageEnd-", pageEnd, "\n")

	//构造通道
	childUrlList := make(chan string, 32)
	newsInfoList := make(chan NewsInfo, 32)

	//独立协程池
	childUrlCatchRoutine := NewGoLimit(16)
	insertDataRoutine := NewGoLimit(32)

	//协程1:翻页并导入子链接
	go func() {
		for i := pageBeg; i <= pageEnd; i++ {
			resp, err := http.Get(FzuInfoPageLink(i))
			if err != nil {
				log.Printf("[Main]Error:访问失败(%s)\n", err)
			}
			FzuInfoChildUrl(resp, childUrlList)
		}
	}()

	//协程2:获取子链接(各通知文件)详细信息
	go func() {
		for !DataCatchOver {
			childUrlCatchRoutine.Add()
			go func() {
				defer childUrlCatchRoutine.Done()
				if childUrl, ok := <-childUrlList; ok {
					GetChildUrlInfo(childUrl, newsInfoList)
				}
			}()
		}
	}()

	//协程3:插入数据库(做日期检验)
	go func() {
		for !DataCatchOver {
			insertDataRoutine.Add()
			go func() {
				defer insertDataRoutine.Done()
				if newsInfo, ok := <-newsInfoList; ok {
					t, _ := time.Parse("2006-01-02", newsInfo.pDate)
					if t.Before(TimeEnd) {
						close(newsInfoList)
					}
					InsertToDatabase(db, newsInfo)
				}
			}()
		}
	}()

	//主协程:等待数据抓取结束(每5s检测一次)
	for !DataCatchOver {
		time.Sleep(5 * time.Second)
	}

	log.Printf("[Main]Info:Task Finished.")
	defer db.Close()
}
