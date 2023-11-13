package main

import (
	"database/sql"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"regexp"
)

const (
	dbDriver   = "mysql"
	dbUser     = "casaos"
	dbPassword = "casaos"
	dbName     = "casaos"
)

func dateSave(title, date, writer, content, click string) {
	db, err := sql.Open(dbDriver, fmt.Sprintf("%s:%s@tcp(192.168.1.201:3306)/%s", dbUser, dbPassword, dbName))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	const createTableQuery = `CREATE TABLE IF NOT EXISTS FZUnews (
	id INT AUTO_INCREMENT PRIMARY KEY,
	title VARCHAR(255),
    click VARCHAR(255),
    date DATE,
    writer VARCHAR(255),
	content LONGTEXT
)`

	_, err = db.Exec(createTableQuery)
	if err != nil {
		panic(err.Error())
	}

	insertQuery := "INSERT INTO FZUnews (title, click, date, writer, content) VALUES (?, ?, ?, ?, ?)"

	_, err = db.Exec(insertQuery, title, click, date, writer, content)
	if err != nil {
		panic(err.Error())
	}

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}
	return
}

func fliter(url1, writer string) {
	var click string
	url := "https://info22.fzu.edu.cn/" + url1
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1)")
	resp, err1 := client.Do(req)
	if err1 != nil {
		fmt.Println("httpget出错")
	}
	defer resp.Body.Close()

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
	title := doc.Find(".conth1").Text()
	date := doc.Find(".conthsj").Contents().First().Text()
	re := regexp.MustCompile(`\d{4}-\d{2}-\d{2}`)
	date = re.FindString(date)

	if !(date >= "2020-01-01" && date <= "2021-09-01") {
		return
	}
	content := doc.Find(".v_news_content").Text()
	dateSave(title, date, writer, content, click)

}

func FZU(url string) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1)")
	resp, err1 := client.Do(req)
	if err1 != nil {
		fmt.Println("httpget出错")
	}
	defer resp.Body.Close()

	doc, _ := goquery.NewDocumentFromReader(resp.Body)
	doc.Find("li.clearfloat").Each(func(i int, ele *goquery.Selection) {
		link, _ := ele.Find("a").Eq(1).Attr("href")
		writer := ele.Find("a.lm_a").Text()
		fliter(link, writer)
	})

}
