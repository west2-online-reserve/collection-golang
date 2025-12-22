package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	fixedOID          = "420981979"                                 // 目标视频的id
	baseURL           = "https://api.bilibili.com/x/v2/reply"       // 一级评论接口地址
	replyURL          = "https://api.bilibili.com/x/v2/reply/reply" // 二级评论接口地址
	firstReqInterval  = 500 * time.Millisecond                     // 一级评论请求时间间隔
	secondReqInterval = 500 * time.Millisecond                     // 二级评论请求时间间隔

	// 数据库一条龙
	mysqlUser     = "root"
	mysqlPassword = "123456"
	mysqlHost     = "localhost"
	mysqlPort     = "3306"
	mysqlDBName   = "bilibili_comments"
)

var DB *sql.DB

// 一级评论结构体
type CommentResp struct {
	Code    int      `json:"code"`    // 接口响应状态码（0成功，1错）
	Message string   `json:"message"` // 接口响应信息（成功为空，失败返回错误信息）
	Data    struct { // 评论的实际内容
		Replies []struct { // 一级评论的相关参数
			Rpid   int64    `json:"rpid"` // 评论唯一id，用于区分评论
			Mid    int64    `json:"mid"`  // 评论用户的id，区分用户
			Member struct { // 评论用户的信息
				Uname     string   `json:"uname"` // 用户名
				LevelInfo struct { // 用户等级
					CurrentLevel int `json:"current_level"`
				} `json:"level_info"`
				Sex string `json:"sex"` // 性别
				Vip struct {
					VipStatus int `json:"vipStatus"` // 是否是会员（1是，0不是）
				}
			} `json:"member"`
			ReplyControl struct {
				Location string `json:"location"` // ip属地
			} `json:"reply_control"`

			Content struct {
				Message string `json:"message"` // 评论内容
			} `json:"content"`

			Ctime int64 `json:"ctime"` // 评论发布时间，秒级时间戳
			Like  int   `json:"like"`  // 评论的点赞数
			Count int   `json:"count"` // 下属的二级评论数
		} `json:"replies"`
		Page struct {
			Pn int `json:"pn"` // 当前页码
		} `json:"page"`
	} `json:"data"`
}

// 二级评论结构体（与一级评论结构体相同的注释会省略）
type SubReplyResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Replies []struct { 
			Rpid   int64 `json:"rpid"`
			Member struct {
				Uname     string `json:"uname"`
				LevelInfo struct {
					CurrentLevel int `json:"current_level"`
				} `json:"level_info"`
				Sex string `json:"sex"`
				Vip struct {
					VipStatus int `json:"vipStatus"`
				}
			} `json:"member"`
			ReplyControl struct {
				Location string `json:"location"`
			} `json:"reply_control"`
			Content struct {
				Message string `json:"message"`
			} `json:"content"`
			Ctime int64 `json:"ctime"`
			Like  int   `json:"like"`
		} `json:"replies"`
		Page struct { // 二级评论分页信息（比一级评论多Ps和Count字段）
			Pn    int `json:"pn,omitempty"`    // 当前页码
			Ps    int `json:"ps,omitempty"`    // 每页条数（默认20条）
			Count int `json:"count,omitempty"` // 二级评论总条数
		} `json:"data"`
	} `json:"data"`
}

func main() {
	if err := initDB(); err != nil {
		fmt.Println("MYSQL连接失败:", err)
		return
	}

	crawlAll()
}

func crawlAll() {
	h := getHeader()
	c := &http.Client{Timeout: 10 * time.Second}
	total, page := 0, 1 // 评论总数，一级评论当前页码

	for {
		params := url.Values{}
		params.Set("oid", fixedOID)
		params.Set("type", "1")
		params.Set("pn", strconv.Itoa(page))
		params.Set("sort", "1")
		reqURL := baseURL + "?" + params.Encode()

		req, _ := http.NewRequest("GET", reqURL, nil)
		req.Header = h
		resp, err := c.Do(req)
		if err != nil {
			time.Sleep(firstReqInterval)
			continue
		}
		
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		var com CommentResp
		json.Unmarshal(body, &com)
		if com.Code != 0 || len(com.Data.Replies) == 0 {
			break
		}

		pageCnt := 0
		for _, f := range com.Data.Replies {
			pageCnt++
			ip := strings.TrimPrefix(f.ReplyControl.Location, "IP属地：")
			if ip == "" {
				ip = "未知"
			}

			vip := 0;
			if f.Member.Vip.VipStatus == 1 {
				vip = 1
			}

			insertToMySQL(1, f.Rpid, 0, f.Member.Uname, f.Member.LevelInfo.CurrentLevel, f.Member.Sex, vip, ip, time.Unix(f.Ctime, 0), f.Like, f.Count, f.Content.Message)

			if f.Count > 0 {
				subCnt := crawlSub(f.Rpid, h, c)
				total += subCnt
			}

			time.Sleep(firstReqInterval)
		}

		total += pageCnt
		fmt.Printf("第%d页完成，累计%d条\n", page, total)
		page++
		time.Sleep(firstReqInterval)
	}

	defer DB.Close()
	fmt.Printf("爬取完成, 共%d条\n", total)
}

func crawlSub(rpid int64, h http.Header, c *http.Client) int {
	cnt, page := 0, 1
	for {
		params :=  url.Values{}
		params.Set("oid", fixedOID)
		params.Set("type", "1")
		params.Set("root", strconv.FormatInt(rpid, 10))
		params.Set("pn", strconv.Itoa(page))
		params.Set("sort", "20")

		reqURL := replyURL + "?" + params.Encode()
		req, _ := http.NewRequest("GET", reqURL, nil)
		req.Header = h
		resp, err := c.Do(req)
		if err != nil {
			time.Sleep(secondReqInterval) 
			continue
		}

		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		var sub SubReplyResp
		json.Unmarshal(body, &sub)
		if sub.Code != 0 || len(sub.Data.Replies) == 0 {
			break
		}

		for _, s := range sub.Data.Replies {
			cnt++ 
			ip := strings.TrimPrefix(s.ReplyControl.Location, "IP属地：")
			if ip == "" {
				ip = "未知"
			}
			vip := 0
			if s.Member.Vip.VipStatus == 1 {
				vip = 1
			}
			insertToMySQL(2, s.Rpid, rpid, s.Member.Uname, s.Member.LevelInfo.CurrentLevel, s.Member.Sex, vip, ip, time.Unix(s.Ctime, 0), s.Like, 0, s.Content.Message)
		}

		if page*sub.Data.Page.Ps >= sub.Data.Page.Count {
			break
		}
		page++
		time.Sleep(secondReqInterval)
	}

	return cnt
}

func insertToMySQL(tp int, rpid, parent int64, name string, level int, sex string, vip int, ip string, ct time.Time, like, sub int, cnt string) error {
	stmt, _ := DB.Prepare(`INSERT INTO comments (comment_type,comment_rpid,parent_rpid,username,user_level,user_sex,user_vip,ip_location,create_time,like_count,sub_reply_count,content) VALUES (?,?,?,?,?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE username=?,user_level=?,user_vip=?,like_count=?,content=?`)
	defer stmt.Close()

	_, err := stmt.Exec(
		tp, rpid, parent, name, level, sex, vip, ip, ct, like, sub, cnt,
		name, level, vip, like, cnt,
	)

	return err
}

func getHeader() http.Header {
	cookie, _ := os.ReadFile("bili_cookie.txt")

	header := http.Header{}
	header.Set("Cookie", strings.TrimSpace(string(cookie)))
	header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36")
	header.Set("Referer", fmt.Sprintf("https://www.bilibili.com/video/av%s", fixedOID))

	return header
}

func initDB() error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true",
		mysqlUser, mysqlPassword, mysqlHost, mysqlPort, mysqlDBName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetMaxIdleConns(5)

	if err := db.Ping(); err != nil {
		return err
	}
	DB= db

	return nil
}
