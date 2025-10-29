package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand/v2"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Config struct {
	Cookie string `json:"cookie"`
}

var Cookie = func() string {
	data, _ := os.ReadFile("config.json")
	var cfg Config
	_ = json.Unmarshal(data, &cfg)
	return cfg.Cookie
}()

type URL struct {
	BaseURL string
}

func (a URL) Page(page int) string {
	return fmt.Sprintf(a.BaseURL, page)
}

type Article struct {
	date, title, url string
	author, text     string
	visitNumber      int
}

func (a Article) save(page int) {
	dir := fmt.Sprintf("articles/page%d", page)
	// 确保目录存在（若不存在则创建）
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Println("创建目录失败：", err)
		return
	}
	safeTitle := regexp.MustCompile(`[<>:"/\\|?*]+`).ReplaceAllString(a.title, "_")
	if safeTitle == "" {
		safeTitle = "untitled"
	}

	path := fmt.Sprintf("%s/%s.md", dir, safeTitle)
	f, err := os.Create(path)
	if err != nil {
		fmt.Println("保存失败：", err)
	}
	defer f.Close()

	fmt.Fprintf(f, "发布时间：%s\n作者：%s\n标题：%s\n访问人数：%d\n正文：%s", a.date, a.author, a.title, a.visitNumber, a.text)
}
func getBody(url string) []byte {
	time.Sleep(100 * time.Millisecond)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("Cookie", Cookie)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("请求失败:", err)
		return []byte{}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取正文失败：", err)
		return []byte{}
	}
	return body
}
func getClickCount(wbnewsid, owner, typ string) int {
	api := fmt.Sprintf(
		"https://info22-443.webvpn.fzu.edu.cn/system/resource/code/news/click/clicktimes.jsp?wbnewsid=%s&owner=%s&type=%s&randomid=%f",
		wbnewsid, owner, typ, rand.Float64(),
	)

	b := fetch(api)
	if len(b) == 0 {
		return 0
	}

	var c struct {
		WBShowTimes int `json:"wbshowtimes"`
	}
	if err := json.Unmarshal(b, &c); err != nil {
		fmt.Println("解析点击数失败:", err)
		return 0
	}
	return c.WBShowTimes
}

func parseFzuCatalog(body []byte) []Article {
	articles := make([]Article, 0)
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		panic(err)
	}
	base, _ := url.Parse("https://info22-443.webvpn.fzu.edu.cn/")
	doc.Find("div.list.fl ul li").Each(func(_ int, li *goquery.Selection) {
		as := li.Find("a")
		if as.Length() < 2 {
			return
		}
		title, _ := as.Eq(1).Attr("title")
		href, _ := as.Eq(1).Attr("href")
		u, err := url.Parse(href)
		if err != nil {
			return
		}
		fullURL := base.ResolveReference(u).String()

		date := li.Find(".fr").First().Text()
		articles = append(articles, Article{
			date:  date,
			title: title,
			url:   fullURL,
		})

	})
	if len(articles) == 0 {
		fmt.Println("articles 空无一物！！！")
	}
	return articles
}
func (a *Article) parseFzuArticle(body []byte) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		panic(err)
	}
	//form := doc.Find(`form[name="_newscontent_frommname"]`)
	// 提取作者信息
	conthsj := doc.Find("div.conthsj").Text()
	conthsj = strings.ReplaceAll(conthsj, "\u00A0", " ") // 去掉 &nbsp; 造成的奇怪空格

	//fmt.Println("提取到的信息栏 ", conthsj)
	// 用正则提取 “信息来源：” 后面的内容
	re := regexp.MustCompile(`信息来源[:：]\s*([^\s]+)`)
	match := re.FindStringSubmatch(conthsj)
	if len(match) > 1 {
		a.author = match[1]
	}
	a.text = doc.Find("div.v_news_content").Text()

	re = regexp.MustCompile(`_showDynClicks\("(\w+)",\s*(\d+),\s*(\d+)\)`)
	match = re.FindStringSubmatch(string(body))
	if len(match) == 4 {
		typ := string(match[1])
		owner := string(match[2])
		wbnewsid := string(match[3])

		clicks := getClickCount(wbnewsid, owner, typ)
		//fmt.Printf("点击数：%d\n", clicks)
		a.visitNumber = clicks
	} else {
		fmt.Printf("点击数解析失败\n")
	}

}
func saveArticle(articles []Article) error {
	f, err := os.Create("articles.csv")
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()

	// 写表头
	if err := w.Write([]string{"date", "title", "url"}); err != nil {
		return fmt.Errorf("写入表头失败: %w", err)
	}
	// 写内容
	for _, a := range articles {
		if err := w.Write([]string{a.date, a.title, a.url}); err != nil {
			return fmt.Errorf("写入行失败: %w", err)
		}
	}
	fmt.Println("已保存到articles.csv")

	return nil
}

func searchPage(l, r int, t string) int {
	x, _ := time.Parse("2006-01-02", t)
	for l < r {
		mid := (l + r + 1) >> 1
		time.Sleep(100 * time.Millisecond)
		fmt.Printf("cherck page = %d\n", mid)
		articles := parseFzuCatalog(getBody(articleBaseurl.Page(mid)))
		y, _ := time.Parse("2006-01-02", articles[0].date)
		if x.Before(y) {
			l = mid
		} else {
			r = mid - 1
		}
	}
	return l
}

var articleBaseurl = URL{
	BaseURL: "https://info22-443.webvpn.fzu.edu.cn/lm_list.jsp?totalpage=1091&PAGENUM=%d&wbtreeid=1460",
}
var (
	sem         = make(chan struct{}, 5)                 // 并发上限 5
	rateLimiter = time.NewTicker(100 * time.Millisecond) // 每 100ms 一个请求
)

func fetch(url string) []byte {
	<-rateLimiter.C // 等令牌（控制速率）

	sem <- struct{}{}        // 占一个并发槽
	defer func() { <-sem }() // 请求结束释放
	return getBody(url)
}

type ArticleModel struct {
	ID          uint      `gorm:"primaryKey"`           // 主键
	Date        string    `gorm:"type:text"`            // 文章发布时间（字符串原样存）
	Title       string    `gorm:"size:255;index"`       // 标题，加索引便于查询
	URL         string    `gorm:"size:512;uniqueIndex"` // 唯一索引，避免重复插入同一篇
	Author      string    `gorm:"size:128"`             // 作者/信息来源
	Text        string    `gorm:"type:text"`            // 正文（text 类型）
	VisitNumber int       `gorm:"index"`                // 点击量（加索引便于范围/排序查询）
	CreatedAt   time.Time // 自动维护：创建时间
	UpdatedAt   time.Time // 自动维护：更新时间
}

func toModel(a Article) ArticleModel {
	return ArticleModel{
		Date:        a.date,        // 映射发布时间
		Title:       a.title,       // 映射标题
		URL:         a.url,         // 映射文章链接（作为唯一约束）
		Author:      a.author,      // 映射作者/来源
		Text:        a.text,        // 映射正文
		VisitNumber: a.visitNumber, // 映射点击量
	}
}
func openDB(dbPath string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{ // 打开 SQLite 数据库文件
		PrepareStmt: true, // 预编译 SQL 语句（适合大量重复插入）
	})
	if err != nil {
		return nil, err // 连接失败则返回错误
	}

	// 自动迁移：根据 ArticleModel 结构创建/更新表结构
	if err := db.AutoMigrate(&ArticleModel{}); err != nil {
		return nil, err // 迁移失败则返回错误
	}
	return db, nil // 返回已就绪的 *gorm.DB
}
func saveArticleORM(db *gorm.DB, a Article) error {
	m := toModel(a) // 结构体转换：Article -> ArticleModel

	// 使用 OnConflict 实现 UPSERT：
	// 冲突目标是 URL（uniqueIndex），发生冲突时更新这些列
	return db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "url"}}, // 唯一冲突列
		DoUpdates: clause.AssignmentColumns([]string{
			"title", "author", "text", "visit_number", "date", "updated_at",
		}), // 冲突时要更新的列
	}).Create(&m).Error // Create 会在无冲突时插入，有冲突时按上面规则更新
}
func main() {

	// body := getBody(articleBaseurl.Page(1))
	// os.WriteFile("fzu.html", body, 0644)

	// articles := parseFzuCatalog(body)
	// if err := saveArticle(articles); err != nil {
	// 	fmt.Println("保存失败：", err)
	// }

	l, r := 1, 1091

	//l, r = searchPage(l, r, "2021-09-01"), searchPage(l, r, "2019-12-31")

	l, r = 294, 411 // 算过的不用再算一遍，以加快测试速度

	fmt.Printf("pageL=%d pageR=%d \n", l, r)

	// jsBody := getBody("https://info22-443.webvpn.fzu.edu.cn/system/resource/js/ajax.js")
	// os.WriteFile("ajax.js", jsBody, 0644)
	db, err := openDB("articles.db") // 打开/创建本地 SQLite 文件数据库
	if err != nil {                  // 连接或迁移失败
		log.Fatalf("open db failed: %v", err) // 直接中止并打印错误
	}

	var wg sync.WaitGroup
	defer rateLimiter.Stop()

	for page := l; page <= r; page++ {
		time.Sleep(100 * time.Millisecond)

		articles := parseFzuCatalog(fetch(articleBaseurl.Page(page)))
		for i, a := range articles {
			wg.Add(1)
			go func(i int, a Article) {
				defer wg.Done()
				// os.WriteFile("article1.html", fetch(a.url), 0644)
				// return
				articles[i].parseFzuArticle(fetch(a.url))
				log.Printf("已爬取 %d 页 第 %d 篇 %s 的文章 \n", page, i, articles[i].author)
				articles[i].save(page)
				if err := saveArticleORM(db, articles[i]); err != nil { // 保存/更新到数据库
					log.Printf("save failed: %v", err) // 若失败打印错误
				}
			}(i, a)
		}

	}

	wg.Wait()
}
