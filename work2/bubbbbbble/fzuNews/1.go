package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

type Notice struct {
	From    string
	Title   string
	Date    string
	Details string
}

func getHtml(url string) (top *html.Node) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
	}
	//请求头UA设置
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36")
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	top, err = htmlquery.Parse(res.Body)
	if err != nil {
		fmt.Println(err)
	}
	return top
}
func getLink(startPage int, linkChan chan string, resChan chan bool) { //获取通知详情页的链接
	for page := startPage; ; page++ {
		url := "https://info22.fzu.edu.cn/lm_list.jsp?totalpage=955&PAGENUM=" + strconv.Itoa(page) + "&wbtreeid=1460" //观察发现不同页码下url只有PAGENUM有变化
		doc := getHtml(url)
		lilist, err := htmlquery.QueryAll(doc, "//li[@class='clearfloat']")
		if err != nil {
			fmt.Println(err)
		}
		for _, li := range lilist {
			linkNode := htmlquery.FindOne(li, "//a[2]")
			link := "https://info22.fzu.edu.cn/" + htmlquery.SelectAttr(linkNode, "href")
			linkChan <- link
		}
		if len(resChan) != 0 { //找到截至时间就关闭管道
			close(linkChan)
			break
		}
	}
}

// 向通知具体内容的链接发起请求
func getContent(linkChan chan string, resChan chan bool) {
	startTime := "2020-01-01"
	endTime := "2021-09-01"
	for {
		link := <-linkChan
		html := getHtml(link)
		fromNode := htmlquery.FindOne(html, "//h3[@class='fl']")
		titleNode := htmlquery.FindOne(html, "//div[@class='conth1']")
		dateNode := htmlquery.FindOne(html, "//div[@class='conthsj']")
		detailsNode := htmlquery.FindOne(html, "//div[@class='v_news_content']")
		from := htmlquery.InnerText(fromNode)
		from = strings.TrimSpace(from) //过滤一下前后空格
		title := htmlquery.InnerText(titleNode)
		date := htmlquery.InnerText(dateNode)
		dateRegex, _ := regexp.Compile("\\d{4}-\\d{2}-\\d{2}")
		date = dateRegex.FindString(date) //正则过滤掉除时间以外的文字
		if date < startTime || date > endTime {
			resChan <- true
			break
		}
		details := htmlquery.InnerText(detailsNode)
		notice := Notice{from, title, date, details}
		DB.Create(&notice)
	}
}
func init() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/mydata?charset=utf8mb4&parseTime=True&loc=Local" // 连接到数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
	}
	DB = db
	DB.AutoMigrate(&Notice{})
}
func main() {
	start := time.Now()
	startPage := 154                   //这个数字是2021-09-01时间的页码
	linkChan := make(chan string, 100) //这是通知详情链接的管道
	resChan := make(chan bool, 4)      //判断协程是否都工作完毕
	go getLink(startPage, linkChan, resChan)
	for i := 0; i < 4; i++ {
		go getContent(linkChan, resChan)
	}
	for {
		if len(resChan) == 4 {
			break
		}
	}
	end := time.Now()
	diff := end.Sub(start)
	fmt.Println(diff)
} //14.9156975s
