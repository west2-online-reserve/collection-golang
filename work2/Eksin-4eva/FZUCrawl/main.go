package main

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	_ "github.com/go-sql-driver/mysql"
)

var Database *sql.DB

type passage struct {
	Id      int
	Title   string
	Time    string
	Author  string
	Click   int
	Content string
}

func initDB() {
	var err error
	dsn := "root:Jack12345@tcp(127.0.0.1:3306)/fzucrawl?charset=utf8mb4&parseTime=True"
	Database, err = sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println("数据库连接失败")
	}
	fmt.Println("数据库连接成功")
}

func insertPassage(item *passage) {
	const query = `
        INSERT INTO FZUCrawl 
        (title, time, origin, click, content) 
        VALUES (?, ?, ?, ?, ?)
    `
	result, err := Database.Exec(query, item.Title, item.Time, item.Author, item.Click, item.Content)
	if err != nil {
		return
	}

	if id, err := result.LastInsertId(); err == nil {
		item.Id = int(id)
	}
}

func fetchHTML(url string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func parsePage(HTML string) (*passage, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(HTML))
	if err != nil {
		return nil, err
	}

	newPsg := &passage{}
	newPsg.Title = strings.TrimSpace(doc.Find(".conth1").Text())
	newPsg.Content = strings.TrimSpace(doc.Find(".v_news_content").Text())

	// 获取包含日期和来源的那一行文字,即<div class="conthsj"></div>的DOM元素
	metaText := doc.Find(".conthsj").Text()
	// "日期： 2020-07-03 信息来源： 教务处"
	// 切成：["日期：", "2020-07-03", "信息来源：", "教务处"]
	parts := strings.Fields(metaText)

	for _, part := range parts {
		if strings.Contains(part, "日期：") {
			val := strings.ReplaceAll(part, "日期：", "")
			newPsg.Time = val
		}
		if strings.Contains(part, "信息来源：") {
			val := strings.ReplaceAll(part, "信息来源：", "")
			newPsg.Author = val
		}
	}

	return newPsg, nil
}

func getURLs(HTML string) []string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(HTML))
	if err != nil {
		fmt.Println("cannot parse html")
		return nil
	}
	var URLs []string
	startTime, _ := time.Parse("2006-01-02", "2020-01-01")
	endTime, _ := time.Parse("2006-01-02", "2021-09-01")
	items := doc.Find(".list li")
	items.Each(func(i int, s *goquery.Selection) {
		dateStr := strings.TrimSpace(s.Find("span.fr").Text())

		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			return
		}

		// 判断日期范围
		if (date.After(startTime) || date.Equal(startTime)) && (date.Before(endTime) || date.Equal(endTime)) {

			targetLink := s.Find("a[href*='content.jsp']")

			href, success := targetLink.Attr("href")
			if success {
				href = "https://info22.fzu.edu.cn/" + href
				URLs = append(URLs, href)
			}
		}
	})
	return URLs
}

func main() {
	initDB()
	defer Database.Close()
	var wg sync.WaitGroup
	Limit := make(chan struct{}, 20)
	startPage, endPage := 400, 500
	fmt.Printf("爬取P%d~P%d", startPage, endPage)
	startTime := time.Now()
	for i := startPage; i <= endPage; i++ {
		wg.Add(1)

		go func(pageIndex int) {
			defer wg.Done()
			Limit <- struct{}{}
			url := fmt.Sprintf("https://info22.fzu.edu.cn/lm_list.jsp?PAGENUM=%d&wbtreeid=1460", pageIndex)
			html, _ := fetchHTML(url)
			links := getURLs(html)
			<-Limit
			for _, link := range links {
				wg.Add(1)
				go func(lk string) {
					defer wg.Done()

					Limit <- struct{}{}
					defer func() {
						<-Limit
					}()

					pas, _ := parsePage(lk)
					insertPassage(pas)
					fmt.Printf("保存：%s", pas.Title)
				}(link)
			}
		}(i)
	}
	wg.Wait()
	fmt.Printf("完成，总耗时: %v\n", time.Since(startTime))
}

//实现了把数据保存到本地MySQL数据库
//使用了Goroutine 同时开启20个协程爬取
