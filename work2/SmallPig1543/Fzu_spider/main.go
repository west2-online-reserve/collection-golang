package main

//运行速度取决于校园网，
import (
	"fmt"
	"github.com/gocolly/colly"
	"strconv"
	"strings"
	"sync"
)

var nextLink string
var year int
var month int
var day int
var wg sync.WaitGroup
var date, author, text, title string
var nums int

// 设置请求头并请求
func setHeader() *colly.Collector {
	c := colly.NewCollector(colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36 Edg/117.0.2045.36"))
	return c
}

// 查询页面，调用查看时间
func queryWeb(url string) {
	c := setHeader()
	c.OnHTML("ul .clearfloat", func(e *colly.HTMLElement) {
		e.ForEachWithBreak(".fr", func(i int, element *colly.HTMLElement) bool {
			date = element.Text
			d := strings.Split(element.Text, "-")
			year, _ = strconv.Atoi(d[0])
			month, _ = strconv.Atoi(d[1])
			_, _ = strconv.Atoi(d[2])
			if year == 2020 || (year == 2021 && (month < 9 || (month == 9 && day == 1))) {
				e.ForEach("a:nth-child(2)", func(i int, elem *colly.HTMLElement) {
					href := elem.Attr("href")
					wg.Add(1)
					go queryText("https://info22.fzu.edu.cn/" + href)
				})
				return true
			}
			return false
		})
	})
	_ = c.Visit(url)
	if year > 2019 {
		clickNext(url)
		queryWeb(nextLink)
	}
}

// 点击下一页
func clickNext(url string) {
	c := setHeader()
	c.OnHTML("body > div.sy-content > div > div.right.fr > div.list.fl > div > span.p_pages > span.p_next.p_fun > a", func(e *colly.HTMLElement) {
		href := e.Attr("href")
		nextLink = "https://info22.fzu.edu.cn/lm_list.jsp?" + href
		fmt.Println(nextLink)
	})
	//https://info22.fzu.edu.cn/lm_list.jsp?totalpage=958&PAGENUM=3&wbtreeid=1460
	_ = c.Visit(url)
}

// 爬取访问人数
func queryNums(url string) {
	c := setHeader()
	c.OnHTML("body", func(e *colly.HTMLElement) {
		nums, _ = strconv.Atoi(e.Text)
	})
	_ = c.Visit(url)
}

// 查看内容
func queryText(url string) {
	c := setHeader()

	//爬取访问人数
	c.OnHTML("body > div.wa1200w > div.conth > form > div.conthsj > script", func(e *colly.HTMLElement) {
		split := strings.Split(e.Text, ",")
		num1, num2 := strings.TrimSpace(strings.Trim(split[2], ")")), strings.TrimSpace(split[1])
		link := "https://info22.fzu.edu.cn/system/resource/code/news/click/dynclicks.jsp?clickid=" + num1 + "&owner=" + num2 + "&clicktype=wbnews"
		queryNums(link)
	})

	c.OnHTML("body > div.wa1200w > div.dqlm.fl > div > div > a:nth-child(4)", func(e *colly.HTMLElement) {
		author = e.Text
	})
	c.OnHTML("div.conth1", func(e *colly.HTMLElement) {
		title = e.Text
	})
	c.OnHTML("#vsb_content > div", func(e *colly.HTMLElement) {
		text = strings.ReplaceAll(e.Text, " ", " ")
	})
	_ = c.Visit(url)
	data := Data{
		Author: author,
		Title:  title,
		Num:    nums,
		Time:   date,
		Text:   text,
	}
	insert(&data)
	defer wg.Done()
}

func main() {
	//数据库的一些准备工作
	Init()
	CreatTable()
	rootUrl := "https://info22.fzu.edu.cn/lm_list.jsp?wbtreeid=1460"
	queryWeb(rootUrl)
	wg.Wait()
}
