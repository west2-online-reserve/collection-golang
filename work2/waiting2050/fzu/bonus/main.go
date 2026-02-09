package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	_ "github.com/go-sql-driver/mysql"
)

type Notice struct {
	Title     string
	Author    string
	Date      string
	Content   string
	URL       string
	ClickNums int
}

type DBConfig struct {
	Username string
	Password string
	Host     string
	Port     string
	DBName   string
}

var (
	client http.Client
	db     *sql.DB
)

func main() {
	dbConfig := DBConfig{
		Username: "root",
		Password: "123456",
		Host:     "localhost",
		Port:     "3306",
		DBName:   "fzu_notices",
	}

	initDB(dbConfig)
	defer db.Close()

	start := time.Now()

	urlList := getLinks(290, 420)
	fmt.Printf("一共拿到%d个链接\n", len(urlList))

	noticeList := getDetails(urlList)
	fmt.Printf("爬到了%d条通知\n", len(noticeList))

	saveToDB(noticeList)
	fmt.Printf("成功保存 %d 条到数据库\n", len(noticeList))

	// 测试
	for i, n := range noticeList {
		if i >= 5 {
			break
		}
		fmt.Printf("\n%d. %s\n", i+1, n.Title)
		fmt.Printf("   时间: %s\n", n.Date)
		fmt.Printf("   来源: %s\n", n.Author)
		fmt.Printf("   点击量: %d\n", n.ClickNums)
		fmt.Printf("   内容: %.100s...\n", n.Content)
	}

	fmt.Printf("\n完事了，花了: %s\n", time.Since(start))
}

func initDB(config DBConfig) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True",
		config.Username, config.Password, config.Host, config.Port, config.DBName)

	db, _ = sql.Open("mysql", dsn)
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(time.Hour)

	createTable()
}

func createTable() {
	sql := `
	CREATE TABLE IF NOT EXISTS notices (
		id INT AUTO_INCREMENT PRIMARY KEY,
		title VARCHAR(500) NOT NULL,
		author VARCHAR(100),
		date DATE,
		content TEXT,
		url VARCHAR(500) UNIQUE,
		click_nums INT DEFAULT 0,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		INDEX idx_date (date),
		INDEX idx_author (author)
	)`

	db.Exec(sql)
	fmt.Println("数据表准备就绪")
}

func saveToDB(notices []Notice) {
	sql := `
	INSERT INTO notices (title, author, date, content, url, click_nums) 
	VALUES (?, ?, ?, ?, ?, ?)
	ON DUPLICATE KEY UPDATE 
		title=VALUES(title), 
		author=VALUES(author), 
		date=VALUES(date), 
		content=VALUES(content),
		click_nums=VALUES(click_nums)`

	tx, _ := db.Begin()

	successCount := 0
	for _, n := range notices {
		_, err := tx.Exec(sql, n.Title, n.Author, n.Date, n.Content, n.URL, n.ClickNums)
		if err == nil {
			successCount++
		}
	}

	tx.Commit()
	fmt.Printf("实际成功保存 %d 条记录\n", successCount)
}

func getDetails(urls []string) []Notice {
	var notices []Notice
	var lock sync.Mutex
	var wg sync.WaitGroup

	for _, v := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			notice, err := getDetail(url)
			if err == nil && notice.Date != "" {
				if isTrueDate(notice.Date) {
					lock.Lock()
					notices = append(notices, notice)
					lock.Unlock()
				} else {
					//fmt.Printf("超出时间范围: %s - %s\n", notice.Date, notice.Title)
				}
			}
		}(v)
		time.Sleep(50 * time.Millisecond)
	}

	wg.Wait()
	return notices
}

func getDetail(url string) (Notice, error) {
	var nt Notice
	nt.URL = url

	req, _ := http.NewRequest("GET", url, nil)
	setHeaders(req)

	resp, err := client.Do(req)
	if err != nil {
		return nt, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nt, err
	}

	nt.Title = strings.TrimSpace(doc.Find("div.conth1").Text())
	infoHtml, _ := doc.Find("div.conthsj").Html()
	if infoHtml != "" {
		dateMatch := regexp.MustCompile(`日期：\s*(\d{4}-\d{2}-\d{2})`).FindStringSubmatch(infoHtml)
		if len(dateMatch) > 1 {
			nt.Date = dateMatch[1]
		}
		authorMatch := regexp.MustCompile(`信息来源：\s*([^&<]+)`).FindStringSubmatch(infoHtml)
		if len(authorMatch) > 1 {
			nt.Author = authorMatch[1]
		}
	}
	nt.Content = strings.TrimSpace(doc.Find(".v_news_content").Text())

	clickNums := getClickNums(doc)
	nt.ClickNums = clickNums

	return nt, nil
}

func getClickNums(doc *goquery.Document) int {
	scriptText := doc.Find("div.conthsj script").Text()
	reg := regexp.MustCompile(`[0-9]+`)
	matches := reg.FindAllString(scriptText, -1)
	
	if len(matches) >= 2 {
		owner := matches[0]
		clickid := matches[1]
		
		clickURL := fmt.Sprintf("https://info22-443.webvpn.fzu.edu.cn/system/resource/code/news/click/dynclicks.jsp?clickid=%s&owner=%s&clicktype=wbnews", clickid, owner)
		
		req, _ := http.NewRequest("GET", clickURL, nil)
		setHeaders(req)
		
		resp, err := client.Do(req)
		if err != nil {
			return 0
		}
		defer resp.Body.Close()
		
		clickDoc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			return 0
		}
		
		clickText := strings.TrimSpace(clickDoc.Find("body").Text())
		clickNums, _ := strconv.Atoi(clickText)
		return clickNums
	}
	
	return 0
}

func isTrueDate(dateStr string) bool {
	if dateStr == "" {
		return false
	}
	return dateStr >= "2020-01-01" && dateStr <= "2021-09-01"
}

func getLinks(st, ed int) []string {
	var allUrls []string
	var lock sync.Mutex
	var wg sync.WaitGroup

	for page := st; page <= ed; page++ {
		wg.Add(1)
		go func(p int) {
			defer wg.Done()
			curUrls := getPagelink(p)
			lock.Lock()
			allUrls = append(allUrls, curUrls...)
			lock.Unlock()
		}(page)
		time.Sleep(100 * time.Millisecond)
	}

	wg.Wait()
	return Unique(allUrls)
}

func getPagelink(page int) []string {
	url := fmt.Sprintf("https://info22-443.webvpn.fzu.edu.cn/lm_list.jsp?totalpage=1093&PAGENUM=%d&urltype=tree.TreeTempUrl&wbtreeid=1460", page)

	req, _ := http.NewRequest("GET", url, nil)
	setHeaders(req)

	resp, err := client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil
	}

	var urls []string
	doc.Find("li a:nth-child(2)").Each(func(i int, s *goquery.Selection) {
		href, ok := s.Attr("href")
		if ok && strings.Contains(href, "content.jsp") {
			fullUrl := fixUrl(href)
			urls = append(urls, fullUrl)
		}
	})

	fmt.Printf("第%d页找到了%d个链接\n", page, len(urls))
	return urls
}

func setHeaders(req *http.Request) {
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Cookie", "_gscu_1331749010=27499092kmwl5c23; Ecp_ClientId=c241216221301948425; _webvpn_key=eyJhbGciOiJIUzI1NiJ9.eyJ1c2VyIjoiMTAyNDAwMzE0IiwiZ3JvdXBzIjpbNl0sImlhdCI6MTc2MTIzMTAzOSwiZXhwIjoxNzYxMzE3NDM5fQ.xLxoSW-1jSbLqbt-e2FWR0uJDSAqFsyar447QJXcFpY; webvpn_username=102400314%7C1761231039%7C1003dcdabd9154f3ef34b608f6c54268f4dec3b2; JSESSIONID=059D3CD2C3636C0E2EF76EB2E35A5442")
}

func fixUrl(href string) string {
	if strings.HasPrefix(href, "http") {
		return href
	}
	return "https://info22-443.webvpn.fzu.edu.cn/" + href
}

func Unique(urls []string) []string {
	vis := make(map[string]bool)
	var res []string
	for _, v := range urls {
		if !vis[v] {
			vis[v] = true
			res = append(res, v)
		}
	}
	return res
}