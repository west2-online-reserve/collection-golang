package main

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/bytedance/sonic"
	_ "github.com/go-sql-driver/mysql"
)

var Db *sql.DB

type Reply struct {
	Ctime   int `json:"ctime"`
	Content struct {
		Message string `json:"message"`
	} `json:"content"`
	Member struct {
		Uname string `json:"uname"`
	} `json:"member"`
}

type BilibiliResp struct {
	Code int `json:"code"`
	Data struct {
		Replies []Reply `json:"replies"`
	} `json:"data"`
}

//CREATE DATABASE IF NOT EXISTS BiliCrawl;
//USE BiliCrawl;
//CREATE TABLE `bilibili` (
//`id` bigint unsigned AUTO_INCREMENT,
//`uname` varchar(1000) DEFAULT NULL,
//`time` varchar(100) DEFAULT NULL,
//`content` longtext,
//PRIMARY KEY (`id`)
//) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

// 初始化数据库
func initDB() {
	var err error
	dsn := "root:YOURPWD@tcp(127.0.0.1:3306)/BiliCrawl?charset=utf8mb4&parseTime=True"
	Db, err = sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	Db.SetMaxOpenConns(20)
	Db.SetMaxIdleConns(10)

	if err := Db.Ping(); err != nil {
		fmt.Println("Database connection failed.")
		panic(err)
	}
	fmt.Println("Database Connected")
}

func timeProcess(timestamp int) string {
	t := time.Unix(int64(timestamp), 0)
	return t.Format("2006-01-02 15:04:05")
}

// 导入数据库
func InsertDB(replies []Reply) {
	if len(replies) == 0 {
		return
	}

	sqlStr := "INSERT INTO bilibili (uname, time, content) VALUES (?, ?, ?)"

	for _, reply := range replies {
		fmtTime := timeProcess(reply.Ctime)
		_, err := Db.Exec(sqlStr, reply.Member.Uname, fmtTime, reply.Content.Message)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		} else {
			fmt.Printf("Saved: %s\n", reply.Content.Message)
		}
	}
}

func main() {
	Limit := make(chan struct{}, 20)
	initDB()
	defer Db.Close()
	startPage, endPage := 1, 20
	var wg sync.WaitGroup
	for i := startPage; i <= endPage; i++ {
		wg.Add(1)
		go func(page int) {
			defer wg.Done()
			Limit <- struct{}{}

			defer func() {
				<-Limit
			}()

			url := fmt.Sprintf("https://api.bilibili.com/x/v2/reply?type=1&oid=420981979&sort=2&pn=%d", page)

			req, _ := http.NewRequest(http.MethodGet, url, nil)
			req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36 Edg/117.0.2045.41")
			req.Header.Add("authority", "api.bilibili.com")
			req.Header.Add("Referer", "https://www.bilibili.com/video/BV12341117rG/?")

			var client http.Client
			resp, err := client.Do(req)
			if err != nil {
				fmt.Printf("Error:%v", err)
			}
			defer resp.Body.Close()
			var result BilibiliResp
			body, _ := io.ReadAll(resp.Body)
			err = sonic.Unmarshal(body, &result)
			if err != nil {
				fmt.Printf("Error:%v", err)
			}
			InsertDB(result.Data.Replies)

		}(i)
	}
	wg.Wait()
	fmt.Println("done")
}
