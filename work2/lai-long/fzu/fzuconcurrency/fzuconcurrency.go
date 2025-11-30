// 1m28.46624258s
// 24m57.047984887s
// 加速比：16.92
package main

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Inform struct {
	Title   string `json:"title"`
	Passage string `json:"passage"`
}

func main() {
	ch := make(chan bool, 10000)
	var times []string
	var urls []string
	start := time.Now()
	//创建文件，保存内容
	file, err := os.Create("./fzuconcurrency.txt")
	if err != nil {
		fmt.Println("creat err", err)
		return
	}
	defer file.Close()
	var flag = 1
	for count := 420; count >= 300; count-- {
		//爬取时间和对应时间的网址
		urls, times = SpiderUrlAndTime(strconv.Itoa(count))
		//go Spiders(urls, times, file, flag,ch)
		for i, _ := range urls {
			if isTrueTime(times[i]) {
				newUrl := "https://info22.fzu.edu.cn/" + urls[i]
				go Spiders(newUrl, times[i], file, &flag, ch)
			}
		}
	}
	for count := 420; count >= 300; count-- {
		urls, _ = SpiderUrlAndTime(strconv.Itoa(count))
		for i, _ := range urls {
			if isTrueTime(times[i]) {
				<-ch
			}
		}
	}
	elapsed := time.Since(start)
	fmt.Println("elapsed", elapsed)
}

func Spiders(urls, times string, file *os.File, flag *int, ch chan bool) {
	//正文所在地址
	//爬取正文
	a := SpiderInclude(urls)
	//写入文件
	file.WriteString("文件" + strconv.Itoa(*flag) + "\n")
	file.WriteString("时间" + times + "\n")
	file.WriteString("标题" + a.Title + "\n")
	file.WriteString("正文" + a.Passage)
	*flag++
	if ch != nil {
		ch <- true
	}
}
func SpiderUrlAndTime(strPage string) (UrlTemp, TimeTemp []string) {
	client := http.Client{}
	Url := "https://info22.fzu.edu.cn/lm_list.jsp?totalpage=1098&PAGENUM=%s&urltype=tree.TreeTempUrl&wbtreeid=1460"
	finalUrl := fmt.Sprintf(Url, strPage)
	req, err := http.NewRequest("GET", finalUrl, nil)
	if err != nil {
		fmt.Println("req err", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("resp err", err)
		return
	}
	defer resp.Body.Close()
	docDetails, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println("docDetails err", err)
	}
	//寻找正文的url和时间
	docDetails.Find("body > div.sy-content > div > div.right.fr > div.list.fl > ul > li").
		Each(func(i int, s *goquery.Selection) {
			Url := s.Find("p > a:nth-child(2)")
			Url2, ok := Url.Attr("href")
			Time := s.Find("p > span")
			TimeRe, _ := regexp.Compile(`\d{4}-\d{2}-\d{2}`)
			if ok {
				UrlTemp = append(UrlTemp, Url2)
				TimeTemp = append(TimeTemp, string(TimeRe.Find([]byte(Time.Text()))))
			}
		})
	return UrlTemp, TimeTemp
}

// 判断时间是否符合2020年1月1号 - 2021年9月1号
func isTrueTime(strTime string) bool {
	t, err := time.Parse("2006-01-02", strTime)
	if err != nil {
		fmt.Println("timer err", err)
	}
	startTime := time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local)
	endTime := time.Date(2021, 9, 1, 23, 59, 59, 59, time.Local)
	return t.After(startTime) && t.Before(endTime)
}

// 爬取正文
func SpiderInclude(url string) Inform {
	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("req2 err", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("resp2 err", err)
	}
	defer resp.Body.Close()
	docDetails, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println("docDetails2 err", err)
	}
	var Data Inform
	Data.Title = docDetails.Find("body > div.wa1200w > div.conth > form > div.conth1").Text()
	Data.Passage = docDetails.Find("#vsb_content > div").Text()
	return Data
}
