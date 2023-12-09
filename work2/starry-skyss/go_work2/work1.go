package main

import (
	"database/sql"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	_ "github.com/antchfx/htmlquery"
	_ "github.com/go-sql-driver/mysql"
	"time"

	//"golang.org/x/text/number"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

const (
	USERNAME = "root"
	PASSWORD = "root"
	HOST     = "127.0.0.1"
	PORT     = "3306"
	DBNAME   = "gowork1"
)

var DB *sql.DB

type Data struct {
	Title   string `json:"title"`
	Writer  string `json:"writer"`
	Time    string `json:"time"`
	Num     string `json:"num"`
	Content string `json:"content"`
}

// 除多余的空白
func removeExtraWhitespace(input string) string {
	// 用空格替换多个连续的空白字符
	input = strings.ReplaceAll(input, "\n", " ")
	input = strings.ReplaceAll(input, "\t", " ")
	input = strings.ReplaceAll(input, "\r", " ")
	input = strings.ReplaceAll(input, "  ", " ")
	input = strings.ReplaceAll(input, " ", " ")
	input = strings.ReplaceAll(input, "  ", " ") // 连续的空格替换为单个空格

	// 去除字符串开头和结尾的空格
	input = strings.TrimSpace(input)

	return input
}

// 发送请求
func HttpGet(url string) (result string, err error) {
	resp, err1 := http.Get(url)
	if err1 != nil {
		err = err1
		return
	}
	defer resp.Body.Close()

	buf := make([]byte, 4*1024)
	for {
		n, _ := resp.Body.Read(buf)
		if n == 0 {
			break
		}
		result += (string(buf[:n]))
	}
	return
}

// 分开作者、时间
func Spite(tmp string) (writer, time string) {
	timeRe, _ := regexp.Compile(`：(.*)    `)
	time = string(timeRe.Find([]byte(tmp)))
	time = strings.Replace(time, "： ", "", -1)
	time = strings.Replace(time, "  ", "", -1)
	writerRe, _ := regexp.Compile(`信息来源： (.*)`)
	writer = string(writerRe.Find([]byte(tmp)))
	writer = strings.Replace(writer, "信息来源： ", "", -1)
	return
}

// 爬取子页面内容
func SpiderOne(url string) (title, writer, content, time, num string, err error) {
	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("req err:", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("请求失败", err)
	}
	defer resp.Body.Close()

	docDetail, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println("解析失败", err)
	}
	docDetail.Find("body > div.wa1200w").
		Each(func(i int, s *goquery.Selection) {
			title = s.Find(" div.conth > form > div.conth1").Text()
			tmp := s.Find("div.conth > form > div.conthsj").Text()
			content = s.Find("#vsb_content").Text()
			content = removeExtraWhitespace(content)
			writer, time = Spite(tmp)
			numurl_o := s.Find("div.conth > form > div.conthsj > script").Text()
			numurltmp := regexp.MustCompile(`, (\d+), `)
			if numurltmp == nil {
				fmt.Println("numurl1 regexp.MustCompile err")
				return
			}
			numurls := numurltmp.FindAllStringSubmatch(numurl_o, -1)
			var numurl1, numurl2 string
			for _, data := range numurls {
				numurl1 = data[1]
			}

			numurltmp = regexp.MustCompile(`, (\d+\))`)
			if numurltmp == nil {
				fmt.Println("numurl2 regexp.MustCompile err")
				return
			}
			numurls = numurltmp.FindAllStringSubmatch(numurl_o, -1)
			for _, data := range numurls {
				numurl2 = data[1]
			}
			numurl2 = strings.Replace(numurl2, ")", "", -1)
			numurl2 = strings.Replace(numurl2, ", ", "", -1)
			numurl := "https://info22.fzu.edu.cn/system/resource/code/news/click/dynclicks.jsp?clickid=" + numurl2 + "&owner=" + numurl1 + "&clicktype=wbnews"
			num, err = HttpGet(numurl) // 替换为目标网站的URL
			if err != nil {
				fmt.Println("HTTP请求错误:", err)
				return
			}

		})
	return
	//总页面
	//body > div.wa1200w
	//取标题
	//body > div.wa1200w > div.conth > form > div.conth1
	//取作者、时间、点击数
	//body > div.wa1200w > div.conth > form > div.conthsj
	//内容
	//#vsb_content
	//点击数
	//body > div.wa1200w > div.conth > form > div.conthsj > script
	//https://info22.fzu.edu.cn/system/resource/code/news/click/dynclicks.jsp?clickid=18131&owner=1768654345&clicktype=wbnews
}

// 爬取主页面内容
func SpiderPage(i int, page chan int) {

	//https://info22.fzu.edu.cn/lm_list.jsp?totalpage=956&PAGENUM=154&urltype=tree.TreeTempUrl&wbtreeid=1460
	url := "https://info22.fzu.edu.cn/lm_list.jsp?totalpage=956&PAGENUM=" + strconv.Itoa(i) + "&urltype=tree.TreeTempUrl&wbtreeid=1460"
	fmt.Printf("正在爬取第%d页的网址:%s \n", i, url)
	//爬取主页面内容
	result, err := HttpGet(url)
	if err != nil {
		fmt.Println("HTTPget err:", err)
		return
	}

	//取子页面url
	//<a href="content.jsp?urltype=news.NewsContentUrl&wbtreeid=1303&wbnewsid=9799" target=_blank title=
	re := regexp.MustCompile(`</a>    
         <a href="(?s:(.*?))"`)
	if re == nil {
		fmt.Println("regexp.MustCompile err")
		return
	}
	urls := re.FindAllStringSubmatch(result, -1)
	for _, datas := range urls {
		var data Data
		datas[1] = "https://info22.fzu.edu.cn/" + datas[1]
		title, writer, content, time, num, err := SpiderOne(datas[1])
		if err != nil {
			fmt.Println("SpiderOne err=", err)
			continue
		}
		data.Title = title
		data.Writer = writer
		data.Time = time
		data.Num = num
		data.Content = content
		if InsertData(data) {
		} else {
			fmt.Println("insert fail")
			return
		}
	}
	fmt.Println("insert success")
	page <- i
}

func InitDB() {
	//root:root@tcp(127.0.0.1:3306)/
	path := strings.Join([]string{USERNAME, ":", PASSWORD, "@tcp(", HOST, ":", PORT, ")/", DBNAME, "?charset=utf8"}, "")
	DB, _ = sql.Open("mysql", path)
	DB.SetConnMaxLifetime(10)
	DB.SetMaxIdleConns(5)
	if err := DB.Ping(); err != nil {
		fmt.Println("open database fail", err)
		return
	}
	fmt.Println("connect success")
}

func InsertData(data Data) bool {
	tx, err := DB.Begin()
	if err != nil {
		fmt.Println("begin err", err)
		return false
	}
	stmt, err := tx.Prepare("INSERT INTO gowork (`title`,`writer`,`time`,`num`,`content`) VALUES (?,?,?,?,?)")
	if err != nil {
		fmt.Println("prepare err", err)
		return false
	}
	_, err = stmt.Exec(data.Title, data.Writer, data.Time, data.Num, data.Content)
	if err != nil {
		fmt.Println("exec fail", err)
		return false
	}
	_ = tx.Commit()
	return true
}

func DoWork(start, end int) {
	fmt.Printf("准备爬取第%d到第%d页的网址\n", start, end)

	page := make(chan int)
	st := time.Now()
	for i := start; i <= end; i++ {
		go SpiderPage(i, page)
	}
	for i := start; i <= end; i++ {
		fmt.Printf("第%d页爬取完毕\n", <-page)
	}
	elapsed := time.Since(st)
	fmt.Printf("NormalStart Time %s \n", elapsed)
}

func main() {
	const START = 155
	const END = 205

	InitDB()
	DoWork(START, END)
}
