package main

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"time"
	_ "github.com/go-sql-driver/mysql"
	"github.com/PuerkitoBio/goquery"
)

const (
	dbDriver   = "mysql"
	dbUser     = "root"
	dbPassword = "071802"
	dbName     = "fzu"
)

func dateSave(title, date, writer, content, click string) {

	//与mysql数据库建立连接
	db, err := sql.Open(dbDriver, fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s", dbUser, dbPassword, dbName))
	if err != nil {
		panic(err.Error())
	}

	//检验是否与数据库连接成功
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	//在fzu数据库中创建fnews表储存爬取数据
	const createTableQuery = `CREATE TABLE IF NOT EXISTS fnews (
	id INT AUTO_INCREMENT PRIMARY KEY,
	title VARCHAR(255),
    click VARCHAR(255),
    date DATE,
    writer VARCHAR(255),
	content LONGTEXT
)`

	//执行创建表语句
	_, err = db.Exec(createTableQuery)
	if err != nil {
		panic(err.Error())
	}

	//定义插入语句
	insertQuery := "INSERT INTO fnews (title, click, date, writer, content) VALUES (?, ?, ?, ?, ?)"

	//执行插入操作
	_, err = db.Exec(insertQuery, title, click, date, writer, content)
	if err != nil {
		panic(err.Error())
	}

	return
}

func Paqu(url1, writer string) {

	var click string

	//进入到每一页中每一个通知的界面
	url := "https://info22.fzu.edu.cn/" + url1

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1)")
	
	//响应请求
	resp, err1 := client.Do(req)
	if err1 != nil {
		panic(err1.Error())
	}

	defer resp.Body.Close()

	//解析resq.Body字段
	doc, _ := goquery.NewDocumentFromReader(resp.Body)
	//获取点击数
	scriptContent := doc.Find(".conthsj script").Text()
	reg := regexp.MustCompile(`_showDynClicks\("wbnews", (\d+), (\d+)\)`) //获取owner与clickid
	matches := reg.FindStringSubmatch(scriptContent)
	//owner clickid 与url拼接
	if len(matches) > 2 {
		owner := matches[1]
		clickid := matches[2]
		url3 := "https://info22.fzu.edu.cn/system/resource/code/news/click/dynclicks.jsp?clickid=" + clickid + "&owner=" + owner + "&clicktype=wbnews"
		resp, _ := http.Get(url3)

		buf := make([]byte, 4096)
		for {
			n, err2 := resp.Body.Read(buf)
			if n == 0 {
				break
			}
			if err2 != nil && err2 != io.EOF {
				return
			}
			click += string(buf[0:n])
		}
	}

	//获取想要爬取的内容
	title := doc.Find(".conth1").Text()
	date := doc.Find(".conthsj").Contents().First().Text()
	re := regexp.MustCompile(`\d{4}-\d{2}-\d{2}`)
	date = re.FindString(date)

	if !(date >= "2020-01-01" && date <= "2021-09-01") {
		return
	}

	content := doc.Find(".v_news_content").Text()

	//储存数据
	dateSave(title, date, writer, content, click)

}

func fzu(url string) {
	
	//建立客户端实体
	client := &http.Client{}
	
	//建立请求
	req, _ := http.NewRequest("GET", url, nil)
	
	//设置请求头
	req.Header.Set("User-Agent", "Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1)")
	
	resp, err1 := client.Do(req)
	if err1 != nil {
		panic(err1.Error())
	}
	defer resp.Body.Close()

	//解析resp.Body字段的内容并生成一个对象
	doc, _ := goquery.NewDocumentFromReader(resp.Body)

	//查询内容,writer和子链接
	doc.Find("li.clearfloat").Each(func(i int, ele *goquery.Selection) {
		link, _ := ele.Find("a").Eq(1).Attr("href")
		writer := ele.Find("a.lm_a").Text()
		Paqu(link, writer)
	})

}

func splider(i int) {

	//网站链接
	url := "https://info22.fzu.edu.cn/lm_list.jsp?totalpage=961&PAGENUM=" + strconv.Itoa(i) + "&wbtreeid=1460"
	fzu(url)
}

func working(start, end int) {
	for i := start; i <= end; i++ {
		//并发执行爬取操作
		go splider(i)
	}

}

func main() {
	start1 := time.Now()
	var start, end int
	start = 1
	end = 30
	working(start, end)
	time := time.Since(start1)
	fmt.Println(time)
}
