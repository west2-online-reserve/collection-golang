package main

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

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
func getContent(doc *html.Node) bool { //获取一个通知的内容
	lilist, err := htmlquery.QueryAll(doc, "//li[@class='clearfloat']")
	if err != nil {
		fmt.Println(err)
	}
	for _, li := range lilist {
		from, title, time, details, flag := getDetails(li) //flag用来判断时间是否符合要求
		if !flag {
			return false //此时时间已经不符合，返回false
		}
		storeInfo(from, title, time, details)
	}
	return true
}

// 传入数据并导入文件
func storeInfo(from string, title string, time string, details string) {
	fileName := "notice.txt"
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	fmt.Fprintln(file, from, title, time)
	fmt.Fprintln(file, details)
	fmt.Fprintln(file, "________________________________________________________________________________________________")

}

// 向通知具体内容的链接发起请求
func getDetails(node *html.Node) (string, string, string, string, bool) {
	linkNode := htmlquery.FindOne(node, "//a[2]")
	link := "https://info22.fzu.edu.cn/" + htmlquery.SelectAttr(linkNode, "href")
	html := getHtml(link)
	fromNode := htmlquery.FindOne(html, "//h3[@class='fl']")
	titleNode := htmlquery.FindOne(html, "//div[@class='conth1']")
	timeNode := htmlquery.FindOne(html, "//div[@class='conthsj']")
	detailsNode := htmlquery.FindOne(html, "//div[@class='v_news_content']")
	from := htmlquery.InnerText(fromNode)
	title := htmlquery.InnerText(titleNode)
	time := htmlquery.InnerText(timeNode)
	dateRegex := regexp.MustCompile("\\d{4}-\\d{2}-\\d{2}")
	trimmedTime := dateRegex.FindString(time) //正则过滤掉除时间以外的文字
	var flag bool
	startTime := "2020-01-01"
	endTime := "2021-09-01"
	if trimmedTime >= startTime && trimmedTime <= endTime {
		flag = true
	} else {
		flag = false
	}
	details := htmlquery.InnerText(detailsNode)
	return from, title, trimmedTime, details, flag
}
func main() {
	start := time.Now()
	startPage := 153                  //这个数字是2021-09-01时间的页码
	for page := startPage; ; page++ { //获取一页通知的内容
		url := "https://info22.fzu.edu.cn/lm_list.jsp?totalpage=954&PAGENUM=" + strconv.Itoa(page) + "&wbtreeid=1460" //观察发现不同页码下url只有PAGENUM有变化
		doc := getHtml(url)                                                                                           //获取初始页html
		sign := getContent(doc)
		if sign {
			fmt.Printf("第%d页已爬取完毕\n", page)
		} else {
			break
		}
	}
	end := time.Now()
	diff := end.Sub(start)
	fmt.Println(diff)

} //并发版以及导入数据库未解决，以及有些通知的内容会出现html内容
//9m13.3077443s
