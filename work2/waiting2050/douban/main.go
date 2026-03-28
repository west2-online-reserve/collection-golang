package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	_ "github.com/go-sql-driver/mysql"
)

const (
	USERNAME = "root"
	PASSWORD = "123456"
	HOST     = "localhost"
	PORT     = "3306"
	DBNAME   = "douban_movie"
)

var DB *sql.DB

type MovieData struct {
	Title    string `json:"title"`
	Director string `json:"director"`
	Picture  string `json:"picture"`
	Actor    string `json:"actor"`
	Year     string `json:"year"`
	Score    string `json:"score"`
	Quote    string `json:"quote"`
}

func main() {
	InitDB()

	for i := 0; i < 10; i++ {
		fmt.Printf("正在爬取第 %d 页\n", i+1)
		Spider(strconv.Itoa(i * 25))
	}
}

func Spider(page string) {
	// 1．发送请求
	client := http.Client{}
	req, err := http.NewRequest("GET", "https://movie.douban.com/top250?start="+page, nil)
	if err != nil {
		fmt.Println("req err", err)
	}

	//加请求头伪造浏览器访问
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36 Edg/141.0.0.0")
	req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("accept-language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
	req.Header.Set("referer", "https://cn.bing.com/")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("请求失败", err)
	}

	defer resp.Body.Close()
	//2.解析网页
	docDetail, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println("解析失败", err)
	}

	//3.获取节点信息
	// #content > div > div.article > ol > li:nth-child(1) > div > div.info > div.hd > a > span:nth-child(1)
	// #content > div > div.article > ol > li:nth-child(1)
	// #content > div > div.article > ol > li:nth-child(2)
	// #content > div > div.article > ol > li:nth-child(1) > div > div.pic > a > img
	// #content > div > div.article > ol > li:nth-child(1) > div > div.info > div.bd > p:nth-child(1)
	// #content > div > div.article > ol > li:nth-child(1) > div > div.info > div.bd > div > span.rating_num
	// #content > div > div.article > ol > li:nth-child(1) > div > div.info > div.bd > p.quote > span
	docDetail.Find("#content > div > div.article > ol > li").
		Each(func(i int, s *goquery.Selection) {
			var data MovieData
			title := s.Find("div > div.info > div.hd > a > span:nth-child(1)").Text()
			img := s.Find("div > div.pic > a > img")
			imgTmp, ok := img.Attr("src")
			info := s.Find("div > div.info > div.bd > p:nth-child(1)").Text()
			score := s.Find("div > div.info > div.bd > div > span.rating_num").Text()
			quote := s.Find("div > div.info > div.bd > p.quote > span").Text()

			if ok {
				director, actor, year := InfoSpite(info)
				data.Title = title
				data.Director = director
				data.Picture = imgTmp
				data.Actor = actor
				data.Year = year
				data.Score = score
				data.Quote = quote

				// 4.保存信息
				if InsertData(data) {
					// fmt.Println("")
				} else {
					fmt.Println("插入失败")
					return
				}
				// fmt.Println("data", data)
			}
		})
	fmt.Println("插入成功")
	return
}

func InitDB() {
	path := strings.Join([]string{USERNAME, ":", PASSWORD, "@tcp(", HOST, ":", PORT, ")/", DBNAME, "?charset=utf8"}, "")
	DB, _ = sql.Open("mysql", path)
	DB.SetConnMaxLifetime(10 * time.Minute)
	DB.SetMaxIdleConns(5)
	if err := DB.Ping(); err != nil {
		fmt.Println("open database fail")
		return
	}
	fmt.Println("connect success")
}

func InfoSpite(info string) (director, actor, year string) {

	directorRe, _ := regexp.Compile(`导演:(.*)主演:`)
	director = string(directorRe.Find([]byte(info)))

	actorRe, _ := regexp.Compile(`主演:(.*)`)
	actor = string(actorRe.Find([]byte(info)))

	yearRe, _ := regexp.Compile(`(\d+)`)
	year = string(yearRe.Find([]byte(info)))

	return
}

func InsertData(movieData MovieData) bool {
	// 开始事务
	tx, err := DB.Begin()
	if err != nil {
		fmt.Println("begin transaction error:", err)
		return false
	}

	// 准备 SQL 语句（使用预处理防止 SQL 注入）
	stmt, err := tx.Prepare(`
        INSERT INTO movie_data (
            Title, Director, Picture, Actor, Year, Score, Quote
        ) VALUES (?, ?, ?, ?, ?, ?, ?)
    `)
	if err != nil {
		fmt.Println("prepare statement error:", err)
		tx.Rollback() // 预处理失败，回滚事务
		return false
	}
	defer stmt.Close() // 延迟关闭预处理语句

	// 执行插入
	_, err = stmt.Exec(
		movieData.Title,
		movieData.Director,
		movieData.Picture,
		movieData.Actor,
		movieData.Year,
		movieData.Score,
		movieData.Quote,
	)
	if err != nil {
		fmt.Println("execute statement error:", err)
		tx.Rollback() // 执行失败，回滚事务
		return false
	}

	// 提交事务
	err = tx.Commit()
	if err != nil {
		fmt.Println("commit transaction error:", err)
		tx.Rollback() // 提交失败，尝试回滚
		return false
	}

	return true
}
