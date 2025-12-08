package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	_ "github.com/go-sql-driver/mysql"
)

// 单线时间： 24m13.5399886s 多线时间： 1m40.7177208s 加速比： 14，玩游戏的时候运行的，实际并发速度应该会更快。

type fzu struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Date   string `json:"date"`
	Click  string `json:"click"`
	Text   string `json:"text"`
}

var DB *sql.DB

func main() {
	//InitDB()
	var wg sync.WaitGroup
	var duration1, duration2 time.Duration
	pagn_l := 296
	pagn_r := 413
	wg.Add(2)
	go func() {
		start1 := time.Now()
		nor_spider(pagn_l, pagn_r)
		wg.Done()
		duration1 = time.Since(start1)
	}()
	go func() {
		start2 := time.Now()
		mul_spider(pagn_l, pagn_r)
		wg.Done()
		duration2 = time.Since(start2)
	}()
	wg.Wait()
	fmt.Println("单线时间：", duration1, "多线时间：", duration2, "加速比：", duration1/duration2)
}

func nor_spider(l int, r int) {
	for i := l; i <= r; i++ {
		spider(strconv.Itoa(i))
	}
}
func mul_spider(l int, r int) {
	var wg sync.WaitGroup
	wg.Add(r - l + 1)
	for i := l; i <= r; i++ {
		go func() {
			spider(strconv.Itoa(i))
			wg.Done()
		}()
	}
	wg.Wait()
}

func spider(pagn string) {
	ul := "https://info22.fzu.edu.cn/lm_list.jsp?totalpage=1091&PAGENUM=" + pagn + "&wbtreeid=1460"
	req, _ := http.NewRequest("GET", ul, nil)
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36 Edg/141.0.0.0")
	req.Header.Set("connection", "keep-alive")
	resp, _ := http.DefaultClient.Do(req)
	defer resp.Body.Close()
	doc, _ := goquery.NewDocumentFromResponse(resp)
	var f fzu
	doc.Find("body > div.sy-content > div > div.right.fr > div.list.fl > ul >li").Each(func(i int, s *goquery.Selection) {
		author := s.Find("p > a.lm_a").Text()
		f.Author = author
		title := s.Find("p > a:nth-child(2)").Text()
		f.Title = title
		date := s.Find("p > span").Text()
		f.Date = date
		sec_ul, _ := s.Find("p > a:nth-child(2)").Attr("href")
		fmt.Println("作者", author)
		fmt.Println("标题", title)
		fmt.Println("日期", date)
		sec_ul = "https://info22.fzu.edu.cn/" + sec_ul
		sec_spider(sec_ul, &f)
		//if insertDB(f) {
		//	fmt.Println("插入成功")
		//}
	})
}

func sec_spider(sec_ul string, f *fzu) {
	req, _ := http.NewRequest("GET", sec_ul, nil)
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36 Edg/141.0.0.0")
	req.Header.Set("connection", "keep-alive")
	resp, _ := http.DefaultClient.Do(req)
	defer resp.Body.Close()
	//https://info22.fzu.edu.cn/system/resource/code/news/click/dynclicks.jsp?clickid=6584&owner=1768654345&clicktype=wbnews
	//body > div.wa1200w > div.conth > form > div.conthsj > script
	doc, _ := goquery.NewDocumentFromResponse(resp)
	click_ul := doc.Find("body > div.wa1200w > div.conth > form > div.conthsj > script").Text()
	refex := regexp.MustCompile(`\d+`)
	owner := refex.FindString(click_ul)
	refex = regexp.MustCompile(`\d+$`)
	clickid := refex.FindString(sec_ul)
	click_ul = "https://info22.fzu.edu.cn/system/resource/code/news/click/dynclicks.jsp?clickid=" + clickid + "&owner=" + owner + "&clicktype=wbnews"
	click_time(click_ul, f)
	text := doc.Find("#vsb_content > div > p").Text()
	(*f).Text = text
	fmt.Println("正文:")
	fmt.Println(text)
}

func click_time(ul string, f *fzu) {
	req, _ := http.NewRequest("GET", ul, nil)
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36 Edg/141.0.0.0")
	req.Header.Set("connection", "keep-alive")
	resp, _ := http.DefaultClient.Do(req)
	defer resp.Body.Close()
	doc, _ := goquery.NewDocumentFromResponse(resp)
	click := doc.Find("body").Text()
	(*f).Click = click
	fmt.Println("访问人数", click)
}

func InitDB() {
	DB, _ = sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/fzu")
	DB.SetConnMaxLifetime(10)
	err := DB.Ping()
	if err != nil {
		fmt.Println("connect database fail")
		return
	}
	fmt.Println("connect database success")
}

func insertDB(f fzu) bool {
	ts, err := DB.Begin()
	if err != nil {
		fmt.Println("err1:", err)
		return false
	}
	stmt, err := ts.Prepare("INSERT INTO `table`(`Title`,`Author`,`Date`,`Click`,`Text`) VALUES (?, ?, ?,?,?)")
	if err != nil {
		fmt.Println("err2:", err)
	}
	_, err = stmt.Exec(f.Title, f.Author, f.Date, f.Click, f.Text)
	if err != nil {
		fmt.Println("err3", err)
		return false
	}
	err = ts.Commit()
	return true
}
