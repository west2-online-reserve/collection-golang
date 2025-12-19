/*
爬取福大通知、文件系统
爬取福州大学通知、文件系统（https://info22.fzu.edu.cn/lm_list.jsp?wbtreeid=1460)

包含发布时间，作者，标题以及正文。
可自动翻页（爬虫可以自动对后续页面进行爬取，而不需要我们指定第几页）
范围：2020年1月1号 - 2021年9月1号（不要爬太多了）。

Bonus
使用并发爬取，同时给出加速比（加速比：相较于普通爬取，快了多少倍）
搜集每个通知的访问人数
将爬取的数据存入数据库，原生SQL或ORM映射都可以
*/

package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 数据库连接配置
const (
	dbUsername = "root"        // MySQL用户名
	dbPassword = "forEG0630"   // MySQL密码
	dbHost     = "127.0.0.1"   // MySQL地址
	dbPort     = "3306"        // MySQL端口
	dbName     = "fzu_notices" // 目标数据库名（自动创建）
)

// ORM结构体（与MySQL表一一映射）
// gorm标签说明：
// - column: 数据库字段名
// - type: 字段类型
// - not null: 非空约束
// - uniqueIndex: 唯一索引（避免重复爬取同一URL）
// - comment: 字段注释
type Notice struct {
	gorm.Model           // 内置字段：ID（主键自增）、CreatedAt（插入时间）、UpdatedAt（更新时间）、DeletedAt（软删除）
	Title      string    `gorm:"column:title;type:varchar(255);not null;comment:'通知标题'"`
	Writer     string    `gorm:"column:writer;type:varchar(100);comment:'发布作者'"`
	NoticeTime time.Time `gorm:"column:notice_time;type:datetime;not null;comment:'通知发布时间'"`
	URL        string    `gorm:"column:url;type:varchar(255);not null;uniqueIndex;comment:'详情页链接（唯一）'"`
	Content    string    `gorm:"column:content;type:text;comment:'通知正文'"`
}

// 原有配置参数
const (
	baseURL      = "https://info22.fzu.edu.cn/"
	listPath     = "/lm_list.jsp?wbtreeid=1460&totalpage=1099"
	startDateStr = "2020-01-01"
	endDateStr   = "2021-09-01"
	userAgent    = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36 Edg/142.0.0.0"
	cookie       = "_ga=GA1.3.134312462.1727272898; _ga_59QEB25NR7=GS1.3.1741261869.114.1.1741262227.0.0.0; _gscu_1331749010=61830325hxpll553; JSESSIONID=04DB7CB5D790CF19079A07BAB3571575"
	concurrency  = 5                      // 最大并发数
	detailDelay  = 500 * time.Millisecond // 单个详情页延迟
)

var (
	StartDate time.Time                                     // 解析后的时间起点
	EndDate   time.Time                                     // 解析后的时间终点
	dateRegex = regexp.MustCompile(`\d{4}-\d{1,2}-\d{1,2}`) // 日期匹配正则
)

// init：初始化时间范围
func init() {
	var err error
	StartDate, err = time.Parse("2006-01-02", startDateStr)
	if err != nil {
		log.Fatalf("解析开始时间失败：%v", err)
	}
	EndDate, err = time.Parse("2006-01-02 15:04:05", endDateStr+" 23:59:59")
	if err != nil {
		log.Fatalf("解析结束时间失败：%v", err)
	}
	log.Printf("时间范围初始化完成：%s ~ %s", StartDate.Format("2006-01-02"), EndDate.Format("2006-01-02 15:04:05"))
}

// 初始化数据库连接（ORM核心）
func initDB() *gorm.DB {
	// 1. 拼接MySQL连接字符串（gorm要求格式）
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		dbUsername, dbPassword, dbHost, dbPort, dbName)

	// 2. 连接MySQL（开启日志，便于调试）
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 打印SQL日志（生产环境可改为Silent关闭）
	})
	if err != nil {
		// 若数据库不存在，自动创建数据库
		if err := createDBIfNotExist(); err != nil {
			log.Fatalf("数据库连接失败且创建数据库失败：%v", err)
		}
		// 重新连接创建后的数据库
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
		if err != nil {
			log.Fatalf("数据库连接失败：%v", err)
		}
	}

	// 3. 设置数据库连接池（避免连接泄露，提升性能）
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("获取连接池失败：%v", err)
	}
	sqlDB.SetMaxIdleConns(10)               // 最大空闲连接数
	sqlDB.SetMaxOpenConns(100)              // 最大打开连接数
	sqlDB.SetConnMaxLifetime(1 * time.Hour) // 连接最大生命周期（1小时）

	// 4. 自动迁移：根据Notice结构体创建/更新表（无需手动写SQL）
	if err := db.AutoMigrate(&Notice{}); err != nil {
		log.Fatalf("创建数据表失败：%v", err)
	}

	log.Println("MySQL数据库初始化成功！")
	return db
}

// 自动创建数据库（如果不存在）
func createDBIfNotExist() error {
	// 连接MySQL服务器（不指定具体数据库）
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=true&loc=Local",
		dbUsername, dbPassword, dbHost, dbPort)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	// 执行创建数据库SQL
	sql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s DEFAULT CHARACTER SET utf8mb4", dbName)
	return db.Exec(sql).Error
}

// spiderDetail返回ORM结构体（收集完整数据）
func spiderDetail(detailURL string) (Notice, error) {
	// 初始化ORM结构体（存储爬取数据）
	notice := Notice{URL: detailURL}

	client := http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", detailURL, nil)
	if err != nil {
		return notice, fmt.Errorf("构造请求失败：%v", err)
	}

	// 设置请求头
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Cookie", cookie)

	resp, err := client.Do(req)
	if err != nil {
		return notice, fmt.Errorf("发送请求失败：%v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return notice, fmt.Errorf("状态码错误：%d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return notice, fmt.Errorf("解析HTML失败：%v", err)
	}

	// 提取字段并赋值给ORM结构体
	notice.Title = strings.TrimSpace(doc.Find("body > div.wa1200w > div.conth > form > div.conth1").Text())
	notice.Writer = strings.TrimSpace(doc.Find("body > div.wa1200w > div.dqlm.fl > h3").Text())
	timeRaw := doc.Find("body > div.wa1200w > div.conth > form > div.conthsj").Text()
	timeStr := dateRegex.FindString(timeRaw)

	// 日期提取失败
	if timeStr == "" {
		return notice, fmt.Errorf("未提取到发布时间")
	}

	// 解析发布时间（ORM自动映射time.Time到MySQL的DATETIME）
	notice.NoticeTime, err = time.Parse("2006-01-02", timeStr)
	if err != nil {
		return notice, fmt.Errorf("解析时间失败：%v", err)
	}

	// 筛选目标时间范围（不在范围则返回错误，不存入数据库）
	inRange := notice.NoticeTime.After(StartDate.Add(-time.Second)) && notice.NoticeTime.Before(EndDate.Add(time.Second))
	if !inRange {
		return notice, fmt.Errorf("时间不在范围：%s", notice.NoticeTime.Format("2006-01-02"))
	}

	// 提取并清理正文
	contentRaw := doc.Find("#vsb_content > div").Text()
	notice.Content = strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(contentRaw, "\n\n", "\n"), "  ", ""))

	// // 输出爬取结果
	// fmt.Println("\n" + strings.Repeat("=", 80))
	// fmt.Printf("标题：%s\n", notice.Title)
	// fmt.Printf("作者：%s\n", notice.Writer)
	// fmt.Printf("时间：%s\n", notice.NoticeTime.Format("2006-01-02"))
	// fmt.Printf("链接：%s\n", notice.URL)
	// fmt.Printf("正文：%s\n", notice.Content)
	// fmt.Println(strings.Repeat("=", 80))

	return notice, nil
}

// spiderListPage接收DB连接，批量存入数据库
func spiderListPage(page int, db *gorm.DB) bool {
	listURL := fmt.Sprintf("%s%s&PAGENUM=%d", baseURL, listPath, page)
	fmt.Printf("\n开始爬取列表页：第%d页（URL：%s）\n", page, listURL)

	// 列表页请求、解析逻辑（不变）
	client := http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", listURL, nil)
	if err != nil {
		fmt.Printf("构造列表页请求失败（%s）：%v\n", listURL, err)
		return false
	}
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Cookie", cookie)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("发送列表页请求失败（%s）：%v\n", listURL, err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		fmt.Printf("列表页第%d页不存在（404），停止翻页\n", page)
		return false
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("列表页第%d页请求失败：状态码 %d\n", page, resp.StatusCode)
		return false
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Printf("解析列表页第%d页失败：%v\n", page, err)
		return false
	}

	// 通知条目选择器
	noticeSelector := "body > div.sy-content > div > div.right.fr > div.list.fl > ul > li"
	crawledURLs := make(map[string]bool)
	hasValidNotice := false
	earliestTime := EndDate

	// 并发安全收集爬取数据
	var validNotices []Notice // 存储当前页所有有效通知（待存入数据库）
	var dataMutex sync.Mutex  // 保护validNotices的并发写入（避免数据竞争）

	// 并发控制核心代码
	sem := make(chan struct{}, concurrency)
	var wg sync.WaitGroup
	var mutex sync.Mutex

	// 遍历通知条目
	doc.Find(noticeSelector).Each(func(index int, s *goquery.Selection) {
		detailA := s.Find("a[href*='news.NewsContentUrl']")
		href, exists := detailA.Attr("href")
		if !exists {
			fmt.Printf("列表页第%d页第%d个条目无详情页a标签，跳过\n", page, index+1)
			return
		}

		// 拼接详情页URL
		var detailURL string
		if strings.HasPrefix(href, "http://") || strings.HasPrefix(href, "https://") {
			detailURL = href
		} else {
			detailURL = baseURL + href
		}

		if crawledURLs[detailURL] {
			fmt.Printf("详情链接（%s）已爬取，跳过\n", detailURL)
			return
		}
		crawledURLs[detailURL] = true

		// 列表页提前筛选时间
		listTimeStr := strings.TrimSpace(s.Find("span.fr").Text())
		if listTimeStr != "" {
			listTimeMatch := dateRegex.FindString(listTimeStr)
			if listTimeMatch != "" {
				listTime, err := time.Parse("2006-01-02", listTimeMatch)
				if err == nil && listTime.Before(StartDate) {
					fmt.Printf("列表页第%d页第%d个条目时间（%s）早于%s，跳过详情页\n", page, index+1, listTimeMatch, StartDate.Format("2006-01-02"))
					mutex.Lock()
					if listTime.Before(earliestTime) {
						earliestTime = listTime
					}
					mutex.Unlock()
					return
				}
			}
		}

		// 并发爬取详情页（收集数据到validNotices）
		wg.Add(1)
		tempURL := detailURL
		go func() {
			defer wg.Done()
			sem <- struct{}{}
			defer func() {
				<-sem
				time.Sleep(detailDelay)
			}()

			// 调用修改后的spiderDetail，获取ORM结构体
			notice, err := spiderDetail(tempURL)
			if err != nil {
				fmt.Printf("爬取详情页失败（%s）：%v\n", tempURL, err)
				return
			}

			// 加锁写入有效通知切片（并发安全）
			dataMutex.Lock()
			validNotices = append(validNotices, notice)
			dataMutex.Unlock()

			// 更新共享变量
			mutex.Lock()
			defer mutex.Unlock()
			hasValidNotice = true
			if notice.NoticeTime.Before(earliestTime) {
				earliestTime = notice.NoticeTime
			}
		}()
	})

	// 等待当前页所有详情页爬取完成
	wg.Wait()
	close(sem)

	// ORM批量存入数据库
	if len(validNotices) > 0 {
		// CreateInBatches：批量插入（每10条一批，提升效率）
		err := db.CreateInBatches(validNotices, 10).Error
		if err != nil {
			log.Printf("第%d页数据批量插入失败：%v", page, err)
		} else {
			log.Printf("第%d页成功插入%d条通知数据到MySQL", page, len(validNotices))
		}
	} else {
		log.Printf("第%d页无有效通知数据，无需插入数据库", page)
	}

	// 翻页判断逻辑
	needContinue := hasValidNotice || earliestTime.After(StartDate.Add(-time.Second))
	if !needContinue {
		fmt.Printf("列表页第%d页最早时间（%s）早于%s，后续无符合条件通知，停止翻页\n", page, earliestTime.Format("2006-01-02"), StartDate.Format("2006-01-02"))
	}

	return needContinue
}

// main函数初始化DB，传入爬虫逻辑
func main() {
	log.Println("开始爬取...")
	start := time.Now()
	defer func() {
		elapsed := time.Since(start)
		log.Printf("爬取任务结束，总耗时：%v", elapsed)
	}()
	log.Println("开始爬取福州大学信息公开网通知（2020-01-01 ~ 2021-09-01）")
	defer log.Println("爬取任务结束")

	// 初始化数据库连接
	db := initDB()
	// 获取SQL连接对象，用于程序结束时关闭
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("获取SQL连接失败：%v", err)
	}
	defer sqlDB.Close() // 程序结束自动关闭数据库连接

	page := 1
	for {
		// 传入db连接到spiderListPage，批量存入数据库
		needContinue := spiderListPage(page, db)
		if !needContinue {
			break
		}
		time.Sleep(2 * time.Second) // 翻页延迟防反爬
		page++
	}
}

/*
	Bonus1:当从第300页开始爬取，爬取日期为 startDateStr = "2021-08-30" endDateStr   = "2021-09-01"时
  	并发爬取：总耗时：21.9250086s
  	普通爬取：总耗时：2m17.331705s
	加速比：6.2637013059141
*/
