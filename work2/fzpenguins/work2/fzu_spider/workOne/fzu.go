package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type MyModel struct {
	ID   uint
	Text string //用于存储字符串的字段
}

func Work() {
	start := time.Now()
	page := make(chan int)
	for i := 155; i <= 272; i++ {
		go spiderPage(i, page)
	}

	for i := 155; i <= 272; i++ {
		fmt.Printf("第%d个页面爬取完成\n", <-page)
	}
	elapsed := time.Since(start)
	fmt.Println("程序执行时间为", elapsed)
	//测试普通爬取的话，把上面的go去掉即可
}

func spiderPage(index int, page chan<- int) {
	url := "https://info22.fzu.edu.cn/lm_list.jsp?totalpage=955&PAGENUM=" + strconv.Itoa(index) + "&wbtreeid=1460"

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("err = ", err)
		return
	}
	defer resp.Body.Close()
	//这里可以设置这个
	//resp.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36 Edg/113.0.1774.42")
	filename := "第" + strconv.Itoa(index) + "页" + ".html"
	f, err := os.Create(filename)
	if err != nil {
		fmt.Println("err = ", err)
		return
	}

	urlCollection := getChildLinks(resp)
	//接下来需要同理爬取每一个子链接的文本

	for i := 0; i < len(urlCollection); i++ {

		f.WriteString(urlCollection[i])
		f.WriteString("\n\n")
		save(urlCollection[i])
	}

	f.Close()
	page <- index
}

func getChildLinkContent(url string) (content string) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("err = ", err)
		return
	}
	//这里可以设置这个
	//req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36 Edg/113.0.1774.42")
	defer resp.Body.Close()

	var buf = make([]byte, 1024*4)
	buf, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("err = ", err)
		return
	}

	content = string(buf)
	return content
}

func getChildLinks(resp *http.Response) (urlCollection []string) {
	url := "https://info22.fzu.edu.cn/"
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println("err = ", err)
		return
	}
	index := 0
	doc.Find("ul").Each(func(i int, s1 *goquery.Selection) {
		s1.Find("p").Each(func(j int, s2 *goquery.Selection) {
			s2.Find("a").Each(func(k int, s3 *goquery.Selection) {
				href, exists := s3.Attr("href")
				if exists { //漏把文本存进去了
					if k == 1 {
						href = strings.TrimSpace(href)
						pageUrl := url + href
						urlCollection = append(urlCollection, getChildLinkContent(pageUrl)) //这里不知是拼接还是不同组
						index++
						//fmt.Println(pageUrl)
					}

				}
			})
		})
	})
	return
}

func save(m string) {
	//user,pass,dbname 都是要根据实际情况确定的
	dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("err = ", err)
		return
	}
	//自动创建表
	db.AutoMigrate(&MyModel{})
	model := MyModel{Text: m}
	db.Create(&model)

}

func main() {
	Work()
}
