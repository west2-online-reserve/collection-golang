package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	_ "github.com/mattn/go-sqlite3"
)

const (
	baseURL = "https://info22.fzu.edu.cn/"
	listURL = "lm_list.jsp?urltype=tree.TreeTempUrl&wbtreeid=1460"

	fmtWithBaseURL  = baseURL + "%s"
	fmtWithListURL  = baseURL + listURL + "&totalpage=%d&PAGENUM=%d"
	fmtWithClickURL = baseURL + "system/resource/code/news/click/dynclicks.jsp?clicktype=%s&owner=%s&clickid=%s"

	// 使用 diff 的方式进行计算，以免页面总数变化对程序的影响（前提：历史文件不删除）
	startPageDiffWithLastPage = (1028 - 230)
	endPageDiffWithLastPage   = (1028 - 250)

	// 选择器
	selectorTotalPagesLink = `.p_last a`
	selectorList           = `body > div.sy-content > div > div.right.fr > div.list.fl > ul > li`
	selectorListLink       = `a[href^="content.jsp"]`
	selectorTitle          = `body > div.wa1200w > div.conth > form > div.conth1`
	selectorContent        = `#vsb_content`

	// 正则
	patternTotalPages = `(?s)totalpage=(?P<totalPages>\d+)`
	patternDate       = `(?s)<div class="conthsj" >日期： (?P<date>.*?)  &nbsp;`
	patternAuthor     = `(?s)信息来源..(?P<author>.*?)\n..`
	patternClicks     = `_showDynClicks\("(?P<clickType>\w+)", (?P<owner>\d+), (?P<clickId>\d+)\)`

	// 数据库
	databasePath = "res/fzu.db"
)

var (
	startTime, _ = time.Parse("2006-01-02 15:04:05", "2020-01-01 00:00:00")
	endTime, _   = time.Parse("2006-01-02 15:04:05", "2021-09-30 23:59:59")
)

func httpGet(url string) (string, error) {
	client := &http.Client{}

	resp, err := client.Get(url)
	if err != nil {
		fmt.Printf("Error occurred while getting %s: %s", url, err)
		return "", err
	}

	defer resp.Body.Close()

	buf := make([]byte, 4096)
	res := ""

	for {
		n, _ := resp.Body.Read(buf)
		if n == 0 {
			break
		}
		res += string(buf[:n])
	}

	return res, nil
}

func parseTotalPages(html string) (int, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return 0, err
	}

	href, exists := doc.Find(selectorTotalPagesLink).Attr("href")
	if !exists {
		return 0, fmt.Errorf("no last page link found")
	}
	regexp := regexp.MustCompile(patternTotalPages)
	matches := regexp.FindStringSubmatch(href)
	if len(matches) < 2 {
		return 0, fmt.Errorf("no total page number found")
	}

	totalPages, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, err
	}

	return totalPages, nil
}

func getTotalPages() (int, error) {
	url := fmt.Sprintf(fmtWithBaseURL, listURL)

	html, err := httpGet(url)
	if err != nil {
		return 0, err
	}

	return parseTotalPages(html)
}

func parseListPage(html string) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	res := []string{}

	doc.Find(selectorList).Each(func(i int, s *goquery.Selection) {
		href, exists := s.Find(selectorListLink).Attr("href")
		if !exists {
			return
		}

		pubDate := s.Find("span").Text()
		timePubDate, err := time.Parse("2006-01-02", pubDate)
		if err != nil {
			return
		}

		if timePubDate.Before(startTime) || timePubDate.After(endTime) {
			return
		}

		res = append(res, href)
	})

	return res, nil
}

func getListPage(page int, totalPages int) ([]string, error) {
	url := fmt.Sprintf(fmtWithListURL, totalPages, page)

	html, err := httpGet(url)
	if err != nil {
		return nil, err
	}

	return parseListPage(html)
}

func getClicks(clickType, owner, clickId string) (int, error) {
	url := fmt.Sprintf(fmtWithClickURL, clickType, owner, clickId)

	data, err := httpGet(url)
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(data)
}

type ArticleResult struct {
	Title    string
	Content  string
	Author   string
	PostDate string // YYYY-MM-DD
	Clicks   int
}

func parseArticlePage(html string) (ArticleResult, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return ArticleResult{}, err
	}

	title := doc.Find(selectorTitle).Text()
	content, err := doc.Find(selectorContent).Html()
	if err != nil {
		return ArticleResult{}, err
	}

	regexpDate := regexp.MustCompile(patternDate)
	regexpAuthor := regexp.MustCompile(patternAuthor)
	regexpClicks := regexp.MustCompile(patternClicks)

	matchesDate := regexpDate.FindStringSubmatch(html)
	matchesAuthor := regexpAuthor.FindStringSubmatch(html)
	matchesClicks := regexpClicks.FindStringSubmatch(html)

	if len(matchesDate) < 2 || len(matchesAuthor) < 2 {
		return ArticleResult{}, fmt.Errorf("no date or author found")
	}

	postDate := matchesDate[1]
	author := matchesAuthor[1]
	clickType := matchesClicks[1]
	owner := matchesClicks[2]
	clickId := matchesClicks[3]

	clicks, err := getClicks(clickType, owner, clickId)
	if err != nil {
		return ArticleResult{}, err
	}

	return ArticleResult{
		Title:    title,
		Author:   author,
		Content:  content,
		PostDate: postDate,
		Clicks:   clicks,
	}, nil
}

func getArticlePage(urlSuffix string) (ArticleResult, error) {
	url := fmt.Sprintf(fmtWithBaseURL, urlSuffix)

	html, err := httpGet(url)
	if err != nil {
		return ArticleResult{}, err
	}

	return parseArticlePage(html)
}

func initDb() *sql.DB {
	os.Remove(databasePath)

	db, err := sql.Open("sqlite3", databasePath)
	if err != nil {
		log.Fatalln("Error occurred while opening database:", err)
	}

	sqlStmt := `
	create table articles (
		id integer not null primary key autoincrement,
		title text,
		content text,
		author text,
		post_date text,
		clicks integer
	);
	`

	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatalf("%q: %s\n", err, sqlStmt)
	}

	return db
}

func insertDb(db *sql.DB, article ArticleResult) {
	stmt, err := db.Prepare("INSERT INTO articles(title, content, author, post_date, clicks) VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatalln("Error occurred while preparing statement:", err)
	}

	defer stmt.Close()

	_, err = stmt.Exec(article.Title, article.Content, article.Author, article.PostDate, article.Clicks)
	if err != nil {
		log.Fatalln("Error occurred while executing statement:", err)
	}
}

func main() {
	db := initDb()
	defer db.Close()

	totalPages, err := getTotalPages()
	if err != nil {
		fmt.Println("Error occurred while getting total pages:", err)
		return
	}

	fmt.Println("Total pages:", totalPages)

	urlList := []string{}
	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, 20) // 限制并发

	for page := (totalPages - startPageDiffWithLastPage); page <= (totalPages - endPageDiffWithLastPage); page++ {
		wg.Add(1)
		go func(page int) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			list, err := getListPage(page, totalPages)
			if err != nil {
				fmt.Println("Error occurred while getting list page:", err)
				return
			}
			mu.Lock()
			urlList = append(urlList, list...)
			mu.Unlock()
		}(page)
	}

	wg.Wait()

	fmt.Println("Total articles:", len(urlList))

	articleChan := make(chan ArticleResult, len(urlList))

	for _, url := range urlList {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			article, err := getArticlePage(url)
			if err != nil {
				fmt.Println("Error occurred while getting article page:", err)
				return
			}
			articleChan <- article
		}(url)
	}

	go func() {
		wg.Wait()
		close(articleChan)
	}()

	for article := range articleChan {
		insertDb(db, article)
	}

	fmt.Println("All articles saved to database:", databasePath)
}
