package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func parse_content(html string) { //爬取内容、标题、作者、时间

	r, _ := http.Get(html)
	dom, err := goquery.NewDocumentFromReader(r.Body)
	if err != nil {
		log.Fatalln(err)
	}
	defer r.Body.Close()
	temp := Result{}

	dom.Find(".v_news_content").Each(func(i int, s *goquery.Selection) {
		s2 := s.Text()
		temp.Content = s2
	}) //爬取内容
	dom.Find(".conthsj").Each(func(i int, s *goquery.Selection) {
		s2 := string(s.Text())
		temp.Day = s2
	}) //爬取日期
	dom.Find(".fl:has(h3)").Each(func(i int, s *goquery.Selection) {

		s2 := s.Text()
		temp.Author = s2

	}) //爬取作者
	dom.Find(".conth1").Each(func(i int, s *goquery.Selection) {

		s2 := s.Text()
		temp.Name = s2

	}) //爬取标题
	temp.Click = click
	DB.Create(&temp)
	if err != nil {
		fmt.Println("write into gorm fail")
	}
}

func parse_next_url(html string) string {
	r, _ := http.Get(html)
	dom, err := goquery.NewDocumentFromReader(r.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var next_url string
	defer r.Body.Close()
	dom.Find(".p_pages").Each(func(i int, s *goquery.Selection) {
		ret, err2 := s.Html()
		if err2 != nil {
			fmt.Println("get html err", err)
			return
		}
		re_link := regexp.MustCompile(`href="(.*?)"`)
		links := re_link.FindAllString(ret, -1)
		base_url := "https://info22.fzu.edu.cn/lm_list.jsp?"
		url := base_url + links[len(links)-2][6:len(links[len(links)-2])-1] //爬取下一个页面的url
		next_url = strings.Replace(url, "&amp;", "&", -1)
	})
	return next_url
}

var click int

func judge(html string) bool { //判断爬取日期是否需要并得到clicks

	r, _ := http.Get(html)
	dom, err := goquery.NewDocumentFromReader(r.Body)
	if err != nil {
		log.Fatalln(err)
	}
	defer r.Body.Close()
	flag := false
	//2020年1月1号 - 2021年9月1号
	dom.Find(".conthsj").Each(func(i int, s *goquery.Selection) {
		s2 := string(s.Text())
		//fmt.Print(s2[10:20])
		var year, month, day int
		year, _ = strconv.Atoi(s2[10:14])
		month, _ = strconv.Atoi(s2[15:17])
		day, _ = strconv.Atoi(s2[18:20])
		l := len(s2)
		num1 := s2[l-6 : l-1]
		num1 = strings.Replace(num1, " ", "", -1) //num1有时5个有时4个，四个时会读取一个“ ”,将空格删除
		clickUrl := "https://info22.fzu.edu.cn/system/resource/code/news/click/dynclicks.jsp?clickid=" + num1 + "&owner=1768654345&clicktype=wbnews"
		c, _ := Httpget(clickUrl)
		click, err = strconv.Atoi(c)
		if year == 2020 {
			flag = true
		} else if year == 2021 && month < 9 {
			flag = true
		} else if year == 2021 && month == 9 && day == 1 {
			flag = true
		}
	})
	return flag
}
func Httpget(url string) (res string, err error) {
	req, _ := http.NewRequest("GET", url, nil)
	// 设置头部信息
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.130 Safari/537.36 OPR/66.0.3515.115")
	resp, err1 := (&http.Client{}).Do(req)
	if err1 != nil {
		fmt.Println("httpget err", err)
		return
	}
	defer resp.Body.Close()

	buf := make([]byte, 4096)
	for {
		n, _ := resp.Body.Read(buf)
		if n == 0 {
			break
		}
		res += string(buf[:n])
	}
	return
}
func parse_final_url(html string) {
	r, _ := http.Get(html)
	dom, err := goquery.NewDocumentFromReader(r.Body)
	if err != nil {
		log.Fatalln(err)
	}
	defer r.Body.Close()
	dom.Find(".clearfloat").Each(func(i int, s *goquery.Selection) {
		ret, err2 := s.Html()
		if err2 != nil {
			fmt.Println("get html err", err)
			return
		}
		re_link := regexp.MustCompile(`href="(.*?)"`)
		links := re_link.FindAllString(ret, -1)
		base_url := "https://info22.fzu.edu.cn/"
		for _, v := range links {
			s := v[6 : len(v)-1]
			url := base_url + s
			furl := strings.Replace(url, "&amp;", "&", -1)
			//fmt.Printf("furl: %v\n", furl)
			if len(furl) > 80 && judge(furl) { //得到的furl有两种，还有个是侧边栏的http，其长度为74，所需http长度为93
				parse_content(furl)
			}
		}
	})
}

func go_work() {
	start := time.Now()
	page := make(chan int)
	url_next := "https://info22.fzu.edu.cn/lm_list.jsp?totalpage=964&PAGENUM=162&wbtreeid=1460"
	for i := 0; i <= 3; i++ {
		go spiderpage(url_next, i, page)
		url_next = parse_next_url(url_next)
	}
	for i := 0; i <= 3; i++ {
		fmt.Printf("第%d页爬取完毕\n", 162+<-page)
	}
	tm := time.Since(start)
	fmt.Printf("tm: %v\n", tm)

}
func spiderpage(url string, i int, page chan int) {
	parse_final_url(url)
	page <- i
}

var DB *gorm.DB

func init_fu() {
	dsn := "root:1914@tcp(127.0.0.1:3306)/myfugorm?charset=utf8mb4&parseTime=True&loc=Local"
	d, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	DB = d
	err = DB.AutoMigrate(&Result{})
	if err != nil {
		fmt.Println(err)
		return
	}
}

type Result struct {
	Click   int
	Day     string
	Content string
	Author  string
	Name    string
}

func main() {
	init_fu()
	go_work()
}
