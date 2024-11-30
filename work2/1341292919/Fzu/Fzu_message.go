package main

import (
	"database/sql"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

var DB *sql.DB

type Fzumessage struct {
	Title     string
	Author    string
	Date      string
	ClickNums int
	Content   string
}

const (
	USERNAME = "casaos"
	PASSWORD = "casaos"
	HOST     = "192.168.1.201"
	PORT     = "3306"
	DBNAME   = "casaos"
)

type date struct {
	Year  int
	Month int
	Day   int
}

var a = make(chan int, 10)

func main() {
	InitDB()
	var wg sync.WaitGroup
	start1 := time.Now()
	for i := 240; i <= 340; i++ { //2020年1月1号 - 2021年9月1号
		a <- i
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			Spider(i)
		}(i)
	}
	wg.Wait()
	elapsed := time.Since(start1)
	fmt.Println("程序耗时：", elapsed)

}

func Spider(i int) {
	//定义用户
	client := http.Client{}
	req, err := http.NewRequest("GET", "https://info22.fzu.edu.cn/lm_list.jsp?totalpage=1028&PAGENUM="+strconv.Itoa(i)+"&urltype=tree.TreeTempUrl&wbtreeid=1460", nil)
	if err != nil {
		fmt.Println("creat err", err)
	}
	req.Header.Set("User-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36 Edg/131.0.0.0")

	//发送请求
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		fmt.Println("request err", err)
	}
	defer resp.Body.Close()

	//解析网页
	docDetail, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println("parse err", err)
	}

	//获取结点信息
	flag := true
	docDetail.Find("body > div.sy-content > div > div.right.fr > div.list.fl > ul > li").
		Each(func(i int, s *goquery.Selection) {
			if flag {
				var Data Fzumessage
				//爬取标题，日期，作者
				Title := s.Find("p > a:nth-child(2)").Text()
				Date := s.Find(" p > span").Text()
				Author := s.Find("p > a.lm_a").Text()
				//获取连接信息以爬取通知内容
				link := "https://info22.fzu.edu.cn/"
				suffix, _ := s.Find("p > a:nth-child(2)").Attr("href")
				link += suffix
				//爬取内容与点击量
				ClickNums, Content := GetDetail(link)
				//存在Data内
				Data.Title = Title
				Data.Date = Date
				Data.Author = Author
				Data.Content = Content
				Data.ClickNums, _ = strconv.Atoi(ClickNums)
				d := DateConvert(Date)
				//检查日期是否符合要求
				if !IsFit(d) {
					return
				}
				InsertData(Data) //将数据存入数据库
			}
		})
	fmt.Printf("爬取完第%d页\n", i)
	<-a
}

func GetDetail(link string) (ClickNums, Content string) {
	//定义用户
	client := http.Client{}
	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		fmt.Println("Enter The Link Failed", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36 Edg/131.0.0.0")

	//发送请求
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("request next link err", err)
	}
	defer resp.Body.Close()

	//解析网页
	docDetail, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println("parse next link err:", link, err)
	}

	//由点击量表达式再访问点击量
	expression := docDetail.Find("body > div.wa1200w > div.conth > form > div.conthsj > script").Text()
	ClickNums = GetClickNums(expression)

	//获取通知文本
	Content, _ = docDetail.Find("#vsb_content").Html()
	Content = HtmlToString(Content)
	return ClickNums, Content
}

func HtmlToString(html string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Fatal(err)
	}
	// 选择文章的容器元素
	articleContent := doc.Find("div.v_news_content")

	// 还原文章结构
	var result strings.Builder
	articleContent.Contents().Each(func(i int, s *goquery.Selection) {
		if goquery.NodeName(s) == "#text" {
			// 处理文本节点
			text := strings.TrimSpace(s.Text())
			if text != "" {
				result.WriteString(text + "\n")
			}
		} else if goquery.NodeName(s) == "p" {
			// 处理段落节点
			pText := strings.TrimSpace(s.Text())
			if pText != "" {
				// 检查段落是否有缩进样式
				indent := s.AttrOr("style", "")
				if strings.Contains(indent, "text-indent:") {
					// 模拟缩进
					result.WriteString("    ") // 4个空格表示缩进
				}
				result.WriteString(pText + "\n")
			}
		} else if goquery.NodeName(s) == "br" {
			// 处理换行节点
			result.WriteString("\n")
		}
	})
	return result.String()
}

// 获取点击量
func GetClickNums(Click string) string {
	//获取点击量链接参数
	reg := regexp.MustCompile("[0-9]+")
	matches := reg.FindAllString(Click, -1)
	owner := matches[0]
	clickid := matches[1]

	//拼接点击量链接
	link := "https://info22.fzu.edu.cn/system/resource/code/news/click/dynclicks.jsp?clickid=" + clickid + "&owner=" + owner + "&clicktype=wbnews"
	client := http.Client{}
	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		fmt.Println("Enter the web failed", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36 Edg/131.0.0.0")

	//发送请求
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("request next link err", err)
	}
	defer resp.Body.Close()

	//解析网页
	docDetail, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println("parse next link err", err)
	}
	nums := docDetail.Find("body").Text()
	return nums

}

func DateConvert(Date string) date {
	reg := regexp.MustCompile("[0-9]+")
	matches := reg.FindAllString(Date, -1)

	var d date
	d.Year, _ = strconv.Atoi(matches[0])
	d.Month, _ = strconv.Atoi(matches[1])
	d.Day, _ = strconv.Atoi(matches[2])
	return d
}

// 判断日期是否符合
func IsFit(Date date) bool {
	EeariestDate := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	LastDate := time.Date(2021, 9, 1, 0, 0, 0, 0, time.UTC)
	// 将输入的日期转换为time.Time类型
	myDate := time.Date(Date.Year, time.Month(Date.Month), Date.Day, 0, 0, 0, 0, time.UTC)
	return myDate.After(EeariestDate) && myDate.Before(LastDate)
}

// 初始化数据库
func InitDB() {
	path := strings.Join([]string{USERNAME, ":", PASSWORD, "@tcp(", HOST, ":", PORT, ")/", DBNAME, "?charset=utf8"}, "")
	var err error
	DB, _ = sql.Open("mysql", path)
	err = DB.Ping()
	if err != nil {
		fmt.Println("error", err)
		os.Exit(1)
	}
	DB.SetConnMaxLifetime(10)
	DB.SetMaxIdleConns(5)
	if err := DB.Ping(); err != nil {
		fmt.Println("OPEN DATABASE FAILED")
		return
	}
	fmt.Println("Open DATABASE DATA successfully")
}

// 数据写入数据库
func InsertData(data Fzumessage) bool {
	tx, err := DB.Begin()
	if err != nil {
		fmt.Println("begin error", err)
		return false
	}
	stmt, err := tx.Prepare("INSERT INTO  FzuMessage(`Title`,`Date`,`Author`,`ClickNums`,`Content`) " +
		"VALUES (?,?,?,?,?)")

	if err != nil {
		fmt.Println("prepare error", err)
		return false
	}
	_, err = stmt.Exec(data.Title, data.Date, data.Author, data.ClickNums, data.Content)
	if err != nil {
		fmt.Println("exec error", err)
		return false
	}
	_ = tx.Commit()
	return true
}
