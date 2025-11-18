package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Inform struct {
	Title   string `json:"title"`
	Passage string `json:"passage"`
	Time    string `json:"time"`
	Author  string `json:"author"`
	Url     string `json:"url"`
}

func main() {
	var urls []string
	var times []string
	for count := 2; count < 15; count++ {
		urls, times = SpiderUrlAndTime(strconv.Itoa(count))
		for i, v := range times {
			if isTrueTime(times[i]) {
				fmt.Println(i, v)
				fmt.Println(urls[i])
			}
		}
	}
}

// 读取时间与正文所在网址
func SpiderUrlAndTime(strPage string) (UrlTemp, TimeTemp []string) {
	client := http.Client{}
	Url := "https://jwch.fzu.edu.cn/gsgg/%s.htm"
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
	//body > div.page > div.w-main > div.wapper > div > div > div:nth-child(3) > div.box-gl.clearfix > ul > li:nth-child(20) > a
	//body > div.page > div.w-main > div.wapper > div > div > div:nth-child(3) > div.box-gl.clearfix > ul > li:nth-child(1)
	//body > div.page > div.w-main > div.wapper > div > div > div:nth-child(3) > div.box-gl.clearfix > ul > li:nth-child(10) > a
	docDetails.Find("body > div.page > div.w-main > div.wapper > div > div > div:nth-child(3) > div.box-gl.clearfix > ul >li").
		Each(func(i int, s *goquery.Selection) {
			Url := s.Find("a")
			Url2, ok := Url.Attr("href")
			Time := s.Find("span")
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
func SpiderInclude(url string) {
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
	title := docDetails.Find("").Text()
}
