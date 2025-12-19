package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	_ "github.com/go-sql-driver/mysql"
)

const (
	USERNAME = "root"
	PASSWORD = "Aa133944"
	HOST     = "192.168.0.3"
	PORT     = "3306"
	DBNAME   = "fzunews"
)

func main() {
	client := http.Client{} //创建客户端
	InitDB()
	defer DB.Close()

	Normal(client)      //普通爬虫耗时 32m44.189886s
	Concurrency(client) //并发爬虫耗时 2m32.7746359s
}

func Concurrency(client http.Client) {
	ch := make(chan bool)
	start := time.Now()
	baseURL := "https://info22.fzu.edu.cn/lm_list.jsp?totalpage=1099&PAGENUM="
	endURL := "&wbtreeid=1460"
	for i := 302; i <= 419; i++ {
		go func(i int) {
			fmt.Println("正在爬取第", i, "页信息")
			URL := baseURL + strconv.Itoa(i) + endURL
			Spider(client, URL, ch)
		}(i)
	}
	for i := 302; i <= 419; i++ {
		<-ch
	}

	t := time.Since(start)
	fmt.Println("并发爬虫耗时", t)
}

func Normal(client http.Client) {
	start := time.Now()
	baseURL := "https://info22.fzu.edu.cn/lm_list.jsp?totalpage=1099&PAGENUM="
	endURL := "&wbtreeid=1460"
	for i := 302; i <= 419; i++ {
		fmt.Println("正在爬取第", i, "页信息")
		URL := baseURL + strconv.Itoa(i) + endURL
		Spider(client, URL, nil)

	}
	t := time.Since(start)
	fmt.Println("普通爬虫耗时", t)
}

func Spider(client http.Client, URL string, ch chan bool) {
	req, err := http.NewRequest("GET", URL, nil) //创建请求
	if err != nil {
		fmt.Println("request error", err)
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-")
	req.Header.Set("Accept-language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
	req.Header.Set("Cache-control", "max-age=0")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Referer", "https://info22-443.webvpn.fzu.edu.cn/lm_list.jsp?totalpage=1099&PAGENUM=303&urltype=tree.TreeTempUrl&wbtreeid=1460")
	req.Header.Set("User-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36 Edg/142.0.0.0")
	req.Header.Set("Upgrade-insecure-requests", "1")

	resp, err := client.Do(req) //发送请求
	if err != nil {
		fmt.Println("Spider fail", err)
	}

	defer resp.Body.Close()

	docDetail, err := goquery.NewDocumentFromReader(resp.Body) //解析网页
	if err != nil {
		fmt.Println("解析失败", err)
	}

	docDetail.Find("body > div.sy-content > div > div.right.fr > div.list.fl > ul > li").
		Each(func(i int, s *goquery.Selection) {
			date := s.Find(" p > span").Text()
			title := s.Find(" p > a:nth-child(2)").Text()
			author := s.Find(" p > a.lm_a").Text()
			text := s.Find("p > a:nth-child(2)")
			textURL, ok := text.Attr("href")
			var fzu_news fzunews
			if ok {
				text, click := Find_MainText(client, textURL)
				fzu_news.Date = date
				fzu_news.Title = title
				fzu_news.Author = author
				fzu_news.Text = text
				fzu_news.Click = click
			}
			if !InsertDate(fzu_news) {
				fmt.Println("insert fail")
				return
			}
		})
	fmt.Println("insert seccuss")
	if ch != nil {
		ch <- true
	}
}

func Find_MainText(client http.Client, textURl string) (maintext, click string) {
	mainTextURl := "https://info22.fzu.edu.cn/" + textURl
	req1, err := http.NewRequest("GET", mainTextURl, nil)
	if err != nil {
		fmt.Println("request error", err)
	}

	req1.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req1.Header.Set("Accept-language", "zh-CN")
	req1.Header.Set("Cache-control", "max-age=0")
	req1.Header.Set("Connection", "keep-alive")
	req1.Header.Set("User-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36 Edg/142.0.0.0")
	req1.Header.Set("Upgrade-insecure-requests", "1")

	resp1, err := client.Do(req1)
	if err != nil {
		fmt.Println("Spider fail", err)
	}

	defer resp1.Body.Close()

	docDetail1, err := goquery.NewDocumentFromReader(resp1.Body)
	if err != nil {
		fmt.Println("解析失败", err)
	}

	maintext = docDetail1.Find("#vsb_content > div").Text()
	//补充了点击量的获取，函数名懒得改了
	cl := docDetail1.Find("body > div.wa1200w > div.conth > form > div.conthsj > script").Text()
	parts := strings.Split(cl, ",")
	owner := strings.TrimSpace(parts[1])
	clickid := strings.TrimSpace(parts[2])
	clickid = strings.TrimSuffix(clickid, ")")
	clickURL := "https://info22.fzu.edu.cn/system/resource/code/news/click/dynclicks.jsp?clickid=" + clickid + "&owner=" + owner + "&clicktype=wbnews"
	req2, err := http.NewRequest("GET", clickURL, nil)
	if err != nil {
		fmt.Println("request error", err)
	}
	resp2, err := client.Do(req2)
	if err != nil {
		fmt.Println("Spider fail", err)
	}

	defer resp2.Body.Close()
	docDetail2, err := goquery.NewDocumentFromReader(resp2.Body)
	if err != nil {
		fmt.Println("解析失败", err)
	}
	click = docDetail2.Text()

	return
}

type fzunews struct {
	Date   string `json:"date"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Text   string `jsom:"text"`
	Click  string `json:"click"`
}

var DB *sql.DB

func InitDB() {
	path := strings.Join([]string{USERNAME, ":", PASSWORD, "@tcp(", HOST, ":", PORT, ")/", DBNAME, "?charset=utf8"}, "")
	var err error
	DB, err = sql.Open("mysql", path)
	if err != nil {
		fmt.Println("connect fail", err)
	}
	DB.SetConnMaxLifetime(time.Minute * 30)
	DB.SetMaxIdleConns(50)
	DB.SetMaxOpenConns(120)
	if err := DB.Ping(); err != nil {
		fmt.Println("opon database fail", err)
		return
	}
	fmt.Println("connect success")
}

func InsertDate(fzu_news fzunews) bool {
	if DB == nil {
		fmt.Println("DB is uninitialized")
		return false
	}
	tx, err := DB.Begin()
	if err != nil {
		fmt.Println("begin err", err)
		return false
	}
	stmt, err := tx.Prepare("INSERT INTO info_data(`Date`,`Title`,`Author`,`Text`,`Click`)VALUES(?,?,?,?,?)")
	if err != nil {
		fmt.Println("prepare fail", err)
		return false
	}
	defer stmt.Close()

	_, err = stmt.Exec(fzu_news.Date, fzu_news.Title, fzu_news.Author, fzu_news.Text, fzu_news.Click)
	if err != nil {
		fmt.Println("exec fail", err)
		tx.Rollback()
		return false
	}
	_ = tx.Commit()
	return true
}
