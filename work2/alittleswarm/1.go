package main

import (
	"fmt"
	"github.com/antchfx/htmlquery"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// 并行：28439345600
// 串行：32270195700
// 加速比：1:11左右
func main() {
	startTime := time.Now()
	ch := make(chan int) //同步
	start, end := setStart()

	for i := start; i <= end; i++ {
		go work(i, ch)

	}
	for i := start; i <= end; i++ {
		<-ch
	}

	endTime := time.Since(startTime)
	println(endTime)
}

// index 设置到第几页为止
// 工作函数
func work(index int, ch chan int) {
	url := "https://info22.fzu.edu.cn/lm_list.jsp?PAGENUM=" + strconv.Itoa(index) + "&wbtreeid=1460"
	data := get(url)
	fmt.Printf("正在爬取第%d页的20条通知\n", index)
	l0, l1, l2, l3 := getUrl(data)
	for i := 0; i < len(l1); i++ {
		if !isWant(l2, i) {
			continue
		} else {
			url2 := "https://info22.fzu.edu.cn/" + l1[i]
			pageData := get(url2)

			savePage(pageData, l0, l2, l3, i)
		}
	}
	ch <- index
}

// 设置起始和结束页
func setStart() (start int, end int) {
	totalPage := getTotalPage()
	start = 170 + totalPage - 951 //以951为基准
	end = 270 + totalPage - 951
	return
}

// 判断日期
func isWant(date []string, index int) bool {
	if !(date[index] >= "2020-01-01" && date[index] <= "2021-09-01") {
		return false
	}
	return true
}

// 获取总共的页数
func getTotalPage() (totalPage int) {
	url := "https://info22.fzu.edu.cn/lm_list.jsp?PAGENUM=1&wbtreeid=1460"
	data := get(url)
	ret := regexp.MustCompile(`totalpage=(.+?)&PAGENUM`)
	result := ret.FindAllStringSubmatch(string(data), 1)
	totalPage, _ = strconv.Atoi(result[0][1])
	return
}

// 保存正文
func savePage(pageData []byte, author []string, date []string, title []string, index int) bool {

	doc, err := htmlquery.Parse(strings.NewReader(string(pageData)))
	if err != nil {
		fmt.Println(err)
		return false
	}
	divContent := htmlquery.Find(doc, "//div[contains(@class, 'v_news_content')]")
	spanContent := htmlquery.Find(divContent[0], "//p")
	var b strings.Builder
	content := ""
	for _, v := range spanContent {
		b.WriteString(htmlquery.InnerText(v))
		b.WriteString("\n")
	}
	divClick := htmlquery.Find(doc, "//div[contains(@class, 'conthsj')]") // 尝试解析访问人数
	tempScript := htmlquery.FindOne(divClick[0], "//script")
	tempString := htmlquery.InnerText(tempScript)
	re := regexp.MustCompile(`_showDynClicks\("[^"]+", (\d+), (\d+)\)`)
	match := re.FindStringSubmatch(tempString) //解析获取click的url
	clickUrl := "https://info22.fzu.edu.cn/system/resource/code/news/click/dynclicks.jsp?clickid=" + match[2] + "&owner=" + match[1] + "&clicktype=wbnews"
	res := get(clickUrl)
	content = b.String()
	tmp := info{}
	tmp.Date = date[index]
	tmp.Content = content
	tmp.Author = author[index]
	tmp.Title = title[index]
	tmp.Clicks, _ = strconv.Atoi(string(res))
	DB.Create(&tmp)
	if err != nil {
		return false
	}
	return true
}

// list0 作者 list1 URL list2 日期 list3 标题
func getUrl(data []byte) (list0 []string, list1 []string, list2 []string, list3 []string) {
	doc, err := htmlquery.Parse(strings.NewReader(string(data)))
	if err != nil {
		fmt.Println(err)
		return
	}
	li := htmlquery.Find(doc, "//li[contains(@class, 'clearfloat')]")
	for i, _ := range li {
		listP := htmlquery.Find(li[i], ".//p")
		for _, v := range listP {

			a := htmlquery.Find(v, "//a")
			list1 = append(list1, htmlquery.SelectAttr(a[1], "href"))
			list0 = append(list0, htmlquery.InnerText(a[0]))
			list3 = append(list3, htmlquery.InnerText(a[1]))
			list2 = append(list2, htmlquery.InnerText(htmlquery.FindOne(v, "//span")))
		}
	}

	return
}

// Get 获取数据
func get(url string) (result []byte) {

	var client http.Client
	req := SetRequest(url)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	result = make([]byte, 0)
	temp := make([]byte, 409600)
	for {
		n, err := resp.Body.Read(temp)
		if n == 0 {
			break
		}
		if err != nil && err != io.EOF {
			// 处理其他错误
			continue
		}
		result = append(result, temp[:n]...)
	}
	return
}

// SetRequest 设置请求头
func SetRequest(url string) (req *http.Request) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {

		return
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36 Edg/115.0.1901.203")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"Windows\"")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
	req.Header.Set("Origin", "https://www.bilibili.com")
	req.Header.Set("Referer", "https://www.bilibili.com/")
	return
}

type info struct {
	Date    string
	Title   string
	Content string
	Author  string
	Clicks  int
}

var DB *gorm.DB

// 链接数据库
func init() {
	username := "root"
	password := "root"
	host := "127.0.0.1"
	port := 3306
	dbname := "gorm"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, dbname)
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: dsn,
	}), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		fmt.Printf("open db failed,err:", err)
		return
	}
	DB = db
	err = DB.AutoMigrate(&info{})
	if err != nil {
		fmt.Println(err)
		return
	}
}
