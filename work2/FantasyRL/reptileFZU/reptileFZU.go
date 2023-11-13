package main

import (
	"database/sql"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

// 日期正则:2021\-0[1-8]\-[0-9][0-9]?|2020\-[0-9][0-9]?\-[0-9][0-9]?|2021\-09\-01
// 使用javascript:alert(document.cookie)查看cookie
// cookie:_gscu_1331749010=90730337klq08s14; JSESSIONID=4145000C75579B69F0D41E8521365227
// <a href="content.jsp?urltype=news.NewsContentUrl&amp;wbtreeid=1303&amp;wbnewsid=34833" target="_blank" title="关于福建省社会科学研究基地福州大学物流研究中心2023年度重大项目申报的通知" style="">关于福建省社会科学研究基地福州大学物流研究中心2023年度重大项目申报的通知</a>
// <span class="fr">2023-09-29</span>
var Db *sql.DB

func init() {
	var err error
	Db, err = sql.Open("mysql", "root:114514@tcp(127.0.0.1:3306)/fzuA")
	if err != nil {
		panic(err)
	}
}

type announce struct {
	Id      int
	Title   string
	Time    string
	Content string
	Origin  string
	Click   int
}

// CREATE TABLE `FZUannounce` (
// `id` bigint unsigned AUTO_INCREMENT,
// `title` varchar(1000) DEFAULT NULL,
// `time` varchar(1000) DEFAULT NULL,
// `origin` varchar(1000) DEFAULT NULL,
// `click` int DEFAULT NULL,
// `content` longtext,
// PRIMARY KEY (`id`),
// UNIQUE KEY `id` (`id`)
// ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
// 存至mysql
func SaveToSQL(ann announce) {
	sql := "insert into FZUannounce (title, time, origin, click, content) values (?, ?, ?, ?, ?)"
	res, err := Db.Exec(sql, ann.Title, ann.Time, ann.Origin, ann.Click, ann.Content)
	if err != nil {
		fmt.Println("save err:", err)
	}
	annId, _ := res.LastInsertId()
	ann.Id = int(annId)
}

// 使用此函数来获取处理后的response
func httpGetFZU(url0 string) (resp string, err error) {
	request, err0 := http.NewRequest(http.MethodGet, url0, nil)
	if err0 != nil {
		err = err0
		return
	}
	//伪装浏览器请求
	request.Header.Add("user-agent", "Mozilla/5.0")
	//校外时使用
	////使用本人cookie实名上网
	//request.Header.Add("Cookie", "_gscu_1331749010=90730337klq08s14; JSESSIONID=4145000C75579B69F0D41E8521365227; _webvpn_key=eyJhbGciOiJIUzI1NiJ9.eyJ1c2VyIjoiMTAyMzAxNTE3IiwiZ3JvdXBzIjpbNl0sImlhdCI6MTY5NjM4ODU5NiwiZXhwIjoxNjk2NDc0OTk2fQ.k245_4oBc3em4U8mCihcPmsxuHKBSdggCZc665FpiWw; webvpn_username=102301517%7C1696388596%7Cc249bf173e78bcae5de3fba7e026301b098c8242")
	////修改referer(不知道有没有用)
	//request.Header.Add("Referer", "https://info22-443.webvpn.fzu.edu.cn/")
	r, err1 := http.DefaultClient.Do(request)
	if err1 != nil {
		err = err1
		return
	}
	defer func() { _ = r.Body.Close() }()
	buf := make([]byte, 4096)
	for {
		content, err2 := r.Body.Read(buf)
		if err2 != nil && err2 != io.EOF {
			err = err2
			return
		}
		if content == 0 {
			break
		}
		resp += string(buf[:content])
	}
	//fmt.Println(url) //test
	return
}

// 使用此函数获取response给goquery用
func goqueryFZU(url0 string) (r *http.Response, err error) {
	request, err0 := http.NewRequest(http.MethodGet, url0, nil)
	if err0 != nil {
		err = err0
		return
	}
	request.Header.Add("user-agent", "Mozilla/5.0")
	request.Header.Add("Cookie", "_gscu_1331749010=90730337klq08s14; JSESSIONID=4145000C75579B69F0D41E8521365227; _webvpn_key=eyJhbGciOiJIUzI1NiJ9.eyJ1c2VyIjoiMTAyMzAxNTE3IiwiZ3JvdXBzIjpbNl0sImlhdCI6MTY5NjM4ODU5NiwiZXhwIjoxNjk2NDc0OTk2fQ.k245_4oBc3em4U8mCihcPmsxuHKBSdggCZc665FpiWw; webvpn_username=102301517%7C1696388596%7Cc249bf173e78bcae5de3fba7e026301b098c8242")
	request.Header.Add("Referer", "https://info22-443.webvpn.fzu.edu.cn/")
	r, err1 := http.DefaultClient.Do(request)
	if err1 != nil {
		err = err1
		return
	}
	return
}

// 使用此函数来进行有效信息筛选
func regexpFZU(urlsp string) (ann announce) {
	//校外使用
	//url := "https://info22-443.webvpn.fzu.edu.cn/" + urlsp
	url := "https://info22.fzu.edu.cn/" + urlsp
	resp, err := httpGetFZU(url)
	if err != nil {
		fmt.Println("pageUrlGet err:", err)
	}

	//标题通过正则实现
	//<div class="conth1" >转发2021年度国家自然科学基金委员会与韩国国家研究基金会合作研究项目指南有关文件的通知</div>
	retTitle := regexp.MustCompile(`<div class="conth1" >(.*)</div>`)
	title := retTitle.FindAllStringSubmatch(resp, -1)

	//时间通过正则实现
	//<div class="conthsj">日期： 2023-09-06  &nbsp; &nbsp; &nbsp;信息来源：
	retTime := regexp.MustCompile(`日期： (.*)  `)
	time := retTime.FindAllStringSubmatch(resp, -1)

	//来源通过正则实现
	//<a href="list.jsp?urltype=tree.TreeTempUrl&wbtreeid=1324">科技处</a>
	//>>
	//	正文
	retOrigin := regexp.MustCompile(`>(.*)</a>
                >>
                正文
`)
	origin := retOrigin.FindAllStringSubmatch(resp, -1)

	//点击量通过正则提取clickid,owner后访问这个是什么东西实现
	//&nbsp; &nbsp; &nbsp;点击数:<script>_showDynClicks("wbnews", 1768654345, 21397)</script></div>
	retResp := regexp.MustCompile(`点击数:<script>_showDynClicks\("wbnews",(.*)\)</script></div>`)
	respClick := retResp.FindAllStringSubmatch(resp, -1)
	var respC string
	for _, v := range respClick {
		respC += v[0]
	}
	retOwner := regexp.MustCompile(`, (.*), [0-9]+\)</script></div>`)
	retClickid := regexp.MustCompile(`[0-9], (.*?)\)</script></div>`)
	owner := retOwner.FindAllStringSubmatch(respC, -1)
	clicked := retClickid.FindAllStringSubmatch(respC, -1)
	//https://info22-443.webvpn.fzu.edu.cn/system/resource/code/news/click/dynclicks.jsp?clickid=13210&owner=1768654345&clicktype=wbnews
	urlClick := "https://info22.fzu.edu.cn/system/resource/code/news/click/dynclicks.jsp?clickid=" + clicked[0][1] + "&owner=" + owner[0][1] + "&clicktype=wbnews"
	click, err := httpGetFZU(urlClick)
	if err != nil {
		fmt.Println("clickGet err:", err)
	}
	Click, _ := strconv.Atoi(click)
	//正文部分使用goquery包解析
	respQuery, err := goqueryFZU(url)
	if err != nil {
		fmt.Println("goquery err:", err)
	}
	defer func() { _ = respQuery.Body.Close() }()
	doc, err := goquery.NewDocumentFromReader(respQuery.Body)
	if err != nil {
		fmt.Println("text err:", err)
	}
	var content string
	doc.Find(".v_news_content").Each(func(i int, s *goquery.Selection) {
		content = s.Find("span").Text()
	})
	/*
		test:
		fmt.Println(resp)
		fmt.Println(title)
		fmt.Println(time)
		fmt.Println(origin)
		fmt.Println(retResp)
		fmt.Println(owner)
		fmt.Println(clicked)
		fmt.Println(click)
		fmt.Println(text)
	*/
	//PrintTest(title[0][1], time[0][1], origin[0][1], click, content)
	ann = announce{
		Title:   title[0][1],
		Time:    time[0][1],
		Origin:  origin[0][1],
		Click:   Click,
		Content: content,
	}
	return
}

// 通过通知文件页获取符合条件的url
func urlGet(url0 string) (urlAll []string, err error) {
	url, err := httpGetFZU(url0)
	if err != nil {
		return
	}
	ret := regexp.MustCompile(`<a href="(.*?)" target=_blank title=(?s:(.*?))(2021\-0[1-8]\-[0-9][0-9]?|2020\-[0-9][0-9]?\-[0-9][0-9]?|2021\-09\-01)`)
	alls := ret.FindAllStringSubmatch(url, -1)
	for i := 0; i < len(alls); i++ {
		urlAll = append(urlAll, alls[i][1])
	}
	//fmt.Println(urlAll) test
	return
}

func main() {
	start := time.Now()
	for i := 150; i < 280; i++ {
		//for i := 160; i < 161; i++ { //test
		//校外使用
		//url0 := "https://info22-443.webvpn.fzu.edu.cn/lm_list.jsp?PAGENUM=" + strconv.Itoa(i) + "&urltype=tree.TreeTempUrl&wbtreeid=1460"
		url0 := "https://info22.fzu.edu.cn/lm_list.jsp?PAGENUM=" + strconv.Itoa(i) + "&wbtreeid=1460"
		urlAll, err := urlGet(url0)
		if err != nil {
			fmt.Println("url0Get err:", err)
		}
		//fmt.Println(urlAll) //test
		for _, v := range urlAll {
			ann := regexpFZU(v)
			SaveToSQL(ann)
		}
	}
	end := time.Since(start)
	fmt.Println(end)

}

// 这是一个确认能打印出有效信息的东西
//func PrintTest(title string, time string, origin string, click string, content string) {
//	fmt.Printf("标题:%s\n时间:%s 来源:%s 点击量:%s\n%s", title, time, origin, click, content)
//
//}
