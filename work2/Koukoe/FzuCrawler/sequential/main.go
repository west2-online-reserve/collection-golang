package main

import (
	"FzuCrawler/common"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// 2025年11月27的时候页数是303到420，但是再多发6条通知的话页数要改了
// 是的过了一天现在2020的第一条被挤到421页了
const (
	startPage = 303
	endPage   = 421
	baseURL   = "https://info22.fzu.edu.cn/lm_list.jsp?wbtreeid=1460"
)

func crawler(file *os.File) {
	client := &http.Client{}
	for j := startPage; j <= endPage; j++ {
		fmt.Printf("正在爬取第 %d 页...\n", j)
		pageURL := fmt.Sprintf("%s&PAGENUM=%d", baseURL, j)
		req, err := http.NewRequest("GET", pageURL, nil)
		if err != nil {
			fmt.Println("Error creating the request:", err)
			return
		}
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36 Edg/142.0.0.0")

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error making the request:", err)
			return
		}
		defer resp.Body.Close()

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			fmt.Println("Error parsing the response body:", err)
			return
		}

		doc.Find("li.clearfloat").Each(func(i int, s *goquery.Selection) {
			titleLink := s.Find("a").Eq(1)
			title := strings.TrimSpace(titleLink.AttrOr("title", titleLink.Text()))
			href, exists := titleLink.Attr("href")
			if !exists {
				return
			}

			// 构造完整详情页链接
			var detailURL string
			if strings.HasPrefix(href, "http") {
				detailURL = href
			} else {
				if strings.HasPrefix(href, "/") {
					detailURL = "https://info22.fzu.edu.cn" + href
				} else {
					detailURL = "https://info22.fzu.edu.cn/" + href
				}
			}

			date := strings.TrimSpace(s.Find("span.fr").Text())
			fmt.Printf("找到通知 %d:\n", i+1)

			if common.IsDateInRange(date) {
				source := strings.TrimSpace(s.Find("a.lm_a").Text())
				content, err := common.FetchContent(detailURL)
				if err != nil {
					fmt.Println("获取正文失败: ", err)
					content = "获取正文失败"
				}
				clicks, err := common.FetchClicks(detailURL)
				if err != nil {
					fmt.Println("获取点击量失败: ", err)
					clicks = "获取点击量失败"
				}
				fmt.Printf("  标题: %s\n", title)
				fmt.Printf("  日期: %s\n", date)
				fmt.Printf("  来源: %s\n", source)
				fmt.Printf("  点击量: %s\n", clicks)
				fmt.Printf("  正文: %s\n", content)
				output := fmt.Sprintf("标题: %s\n日期: %s\n来源: %s\n点击量: %s\n正文: %s\n------\n",
					title, date, source, clicks, content)
				_, err = file.WriteString(output)
				if err != nil {
					fmt.Printf("写入文件失败: %v\n", err)
				}

			} else {
				fmt.Printf("日期不符合要求 (%s)，跳过\n", date)
			}
			fmt.Println("------")
		})
	}
}

func main() {
	start := time.Now()

	file, err := os.Create("FzuNotices.txt") //  为什么通知文件系统的URL里写的是news
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	crawler(file)
	elapsed := time.Since(start)
	fmt.Printf("\n任务用时: %v\n", elapsed) //  爬了一遍用时:43m34.9631875s，好慢。。。
}
