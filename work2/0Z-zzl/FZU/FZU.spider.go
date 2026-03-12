package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

type FZU_test struct {
	Title  string `gorm:"column:title"`
	Time   string `gorm:"column:time"`
	Writer string `gorm:"column:writer"`
	Test   string `gorm:"type:text"`
}

// 错误方法
//
//	func istimeright(times string) bool {
//		timeT, err := time.Parse(" 2006-01-02 ", times)
//		if err != nil {
//			fmt.Println("time字符串转换失败")
//		}
//		start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
//		end := time.Date(2021, 9, 1, 0, 0, 0, 0, time.UTC)
//
//		if timeT.Before(end) && timeT.After(start) {
//			return true
//		} else {
//			return false
//		}
//	}
func istimeright(times string) bool {
	// 先清洗时间字符串：只保留数字和-
	cleanTime := regexp.MustCompile(`[^\d-]`).ReplaceAllString(times, "")
	// 处理空值
	if cleanTime == "" {
		fmt.Println("清洗后时间为空：", times)
		return false
	}

	// 定义可能的时间格式列表
	formats := []string{
		"2006-01-02",
		"20060102",
		"2006-1-2",
	}

	var timeT time.Time
	var err error
	for _, f := range formats {
		timeT, err = time.Parse(f, cleanTime)
		if err == nil {
			break
		}
	}
	if err != nil {
		fmt.Printf("时间解析失败（原始：%s，清洗后：%s）：%v\n", times, cleanTime, err)
		return false
	}

	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2021, 9, 1, 0, 0, 0, 0, time.UTC)

	return timeT.Before(end) && timeT.After(start)
}

// 连接数据库
func Init() {
	var dsn = "root:123456@tcp(127.0.0.1:3306)/mysql?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("数据库连接失败!", err)
	}
	err = db.AutoMigrate(&FZU_test{})
	if err != nil {
		fmt.Println("数据库表迁移失败!", err)
	}
	fmt.Println("数据库连接成功")
}

// 导入数据
func spidersql(title, text, time, writer string) {
	if db == nil {
		fmt.Println("⚠️ 数据库未连接，跳过保存")
		return
	}
	title = strings.TrimSpace(title)
	time = strings.TrimSpace(time)
	writer = strings.TrimSpace(writer)
	time = strings.Join(strings.Fields(time), " ")

	FZU_tests := FZU_test{
		Title:  title,
		Time:   time,
		Writer: writer,
		Test:   text,
	}
	db.Create(&FZU_tests)
	fmt.Println(FZU_tests)
}

// `日期：\s*(.*?)\s*信息来源：\s*(.*?)\s*`
func sumSpite(sum string) (string, string) {
	match := regexp.MustCompile(`日期：\s*(.*?)\s*信息来源：\s*(.*?)\s*`).FindStringSubmatch(sum)
	times := match[1]
	writer := match[2]
	return times, writer
}
func SpiderInside(InsideUrl string) {
	//1.发送请求
	var client http.Client
	req, err := http.NewRequest("GET", InsideUrl, nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	//2.解析网页
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println("err:", err)
	}
	//3.获取节点信息
	//body > div.wa1200w > div.conth > form
	//body > div.wa1200w > div.conth > form > div.conth1
	doc.Find("body > div.wa1200w > div.conth > form").
		Each(func(i int, s *goquery.Selection) {
			title := s.Find("div.conth1").Text()
			text := s.Find("#vsb_content").Text()
			sum := s.Find(" div.conthsj").Text()
			times, writer := sumSpite(sum)
			ok := istimeright(times)
			if ok {
				spidersql(title, text, times, writer)
			}
		})
}
func SpiderList(page int) {
	//爬取发布时间，作者，标题以及正文。
	// 1.发送请求
	var client http.Client
	url := fmt.Sprintf("https://info22.fzu.edu.cn/lm_list.jsp?totalpage=1100&PAGENUM=%d&wbtreeid=1460", page)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("err:", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("发送请求失败:", err)
		return
	}
	defer resp.Body.Close()
	// 2.解析网页
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println("err:", err)
	}
	// 3.获取链接
	//body > div.sy-content > div > div.right.fr > div.list.fl > ul > li:nth-child(1)
	doc.Find("body > div.sy-content > div > div.right.fr > div.list.fl > ul > li").
		Each(func(i int, s *goquery.Selection) {
			//用于获取正文
			s.Find(" p > a:nth-child(2)").
				Each(func(i int, s *goquery.Selection) {
					href, _ := s.Attr("href")
					InsideUrl := "https://info22.fzu.edu.cn/" + href
					//fmt.Println(InsideUrl)
					SpiderInside(InsideUrl)
				})
		})

}
func main() {
	start := time.Now()
	Init()
	var i int
	for i = 430; i > 300; i-- {
		fmt.Printf("正在爬取第%d页的内容\n", i)
		SpiderList(i)
	}
	defer func() {
		end := time.Since(start)
		fmt.Printf("程序执行时间：%s\n", end)
	}()
}
