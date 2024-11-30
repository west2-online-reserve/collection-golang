/*
校外ip地址需要通过VPN或者WebVPN访问通知系统，而由VPN登录页面的"通知，文件系统"跳转到
的文件系统https://info22-fzu-edu-cn.fzu.edu.cn/（而不是https://info22.fzu.edu.cn/）
是动态网站，用goquery读取的页面会返回如<div class="conth1">Loading..</div>使用Loading...
暂时占用未加载资源的情况，所以该页面是动态加载的，使用chromedp可以等待资源加载好再爬取.
*/
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/chromedp/chromedp"
	_ "github.com/go-sql-driver/mysql"
)

type Notification struct {
	Title   string
	Author  string
	Date    string
	Content string
}

const (
	username = "root"
	password = "zth20041017"
	hostname = "localhost"
	port     = "3306"
	dbname   = "fzu_notification"
)

func dsn(dbname string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, hostname, port, dbname)
}

func connectDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn(dbname))
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func insertNotification(db *sql.DB, notification Notification) {
	query := "INSERT INTO notifications (title, author, date, content) VALUES (?, ?, ?, ?)"
	_, err := db.Exec(query, notification.Title, notification.Author, notification.Date, notification.Content)
	if err != nil {
		log.Printf("插入数据失败: %v\n", err)
	} else {
		log.Printf("成功插入通知: %s\n", notification.Title)
	}
}

func parseNotification(pageURL string, db *sql.DB) {
	var title string
	var date string
	var content string
	var author string

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	err := chromedp.Run(ctx,
		chromedp.Navigate(pageURL),
		chromedp.WaitVisible(`div.conth1`, chromedp.ByQuery),
		chromedp.OuterHTML(`div.conth1`, &title),
		chromedp.OuterHTML(`div.conthsj`, &date),
		chromedp.OuterHTML(`div.v_news_content`, &content),
		chromedp.OuterHTML(`h3.fl`, &author),
	)

	if err != nil {
		log.Fatal(err)
	}

	re := regexp.MustCompile(`<div class="conth1">(.*?)<\/div>`)
	matches := re.FindStringSubmatch(title)
	title = matches[1]

	re = regexp.MustCompile(`(\d{4}-\d{2}-\d{2})`)
	matches = re.FindStringSubmatch(date)
	date = matches[1]

	re = regexp.MustCompile(`>([^<>]+)<`)
	contentMatches := re.FindAllStringSubmatch(content, -1)
	content = ""

	for _, i := range contentMatches {
		content += i[1]
	}

	re = regexp.MustCompile(`<h3 class="fl">\s*(.*?)<\/h3>`)
	matches = re.FindStringSubmatch(author)
	author = matches[1]

	notification := Notification{
		Title:   title,
		Author:  author,
		Date:    date,
		Content: content,
	}

	insertNotification(db, notification)
}

func spiderFunc(crawlValue bool, db *sql.DB) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	//ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	//defer cancel()

	pageIndex := 227
	layout := "2006-01-02"
	var strPageIndex string
	var mainPageURL string
	var startPageA string
	var wg sync.WaitGroup
	var sem = make(chan struct{}, 20) // 限制最大并发数

	startDate, _ := time.Parse(layout, "2019-12-31")
	endDate, _ := time.Parse(layout, "2021-09-02")
	//2020-01-01---2021-09-01

outerLook:
	for {
		strPageIndex = strconv.Itoa(pageIndex)
		mainPageURL = fmt.Sprintf("https://info22-fzu-edu-cn.fzu.edu.cn/lm_list.jsp?totalpage=1023&PAGENUM=%s&urltype=tree.TreeTempUrl&wbtreeid=1460", strPageIndex)

		var dateNow string

		err := chromedp.Run(ctx,
			chromedp.Navigate(mainPageURL),
			chromedp.WaitVisible(`span.fr`),
			chromedp.OuterHTML(`html`, &dateNow),
		)

		if err != nil {
			log.Fatal(err)
		}
		//fmt.Printf(dateNow)

		re := regexp.MustCompile(`<span class="fr">(.*?)</span>`)
		matchesDate := re.FindAllStringSubmatch(dateNow, -1)

		re = regexp.MustCompile(`<a href="(.*?)" target="_blank" title=`)
		matchesDateA := re.FindAllStringSubmatch(dateNow, -1)
		//fmt.Println(matchesDate[1][1])
		//var startPageA, endPageA string

		for index, eachDate := range matchesDate {
			date, _ := time.Parse(layout, eachDate[1])
			if date.After(startDate) && date.Before(endDate) {
				startPageA = matchesDateA[index][1]
				startPageA = strings.Replace(startPageA, "amp;", "", -1)
				pageUrl := "https://info22-fzu-edu-cn.fzu.edu.cn/" + startPageA
				if crawlValue {
					wg.Add(1)
					go func(p string) {
						sem <- struct{}{}
						defer func() { <-sem }()
						parseNotification(p, db)
					}(pageUrl)
				} else {
					parseNotification(pageUrl, db)
				}

			}
			if date.Before(startDate) {
				break outerLook
			}
		}
		pageIndex += 1
	}
	wg.Wait()
}

func measureSpeed(singleThreadFunc, concurrentFunc func()) {
	start := time.Now()
	concurrentFunc()
	concurrentDuration := time.Since(start)

	start = time.Now()
	singleThreadFunc()
	singleThreadDuration := time.Since(start)

	speedup := float64(singleThreadDuration) / float64(concurrentDuration)
	fmt.Printf("加速比: %.2f\n", speedup)
}

func main() {
	db, err := connectDB()
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS notifications (
	id INT AUTO_INCREMENT PRIMARY KEY,
	title VARCHAR(255),
	author VARCHAR(100),
	date VARCHAR(50),
	content TEXT
	);`)
	if err != nil {
		log.Fatalf("创建表失败: %v", err)
	}

	measureSpeed(
		func() { spiderFunc(false, db) },
		func() { spiderFunc(true, db) },
	)

}
