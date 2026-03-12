//加速比：理论上为10

package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	Insert "github.com/jhw66/mycrawler/insertdatabase"
)

func NewRequest(method string, url string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) "+
		"AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	return req, nil
}

func fetchClick(clickid string) (string, error) {
	//提取文章点击数
	//解析URL，处理相对路径
	href := "https://info22.fzu.edu.cn/system/resource/code/news/click/dynclicks.jsp?clickid=" + clickid + "&owner=1768654345&clicktype=wbnews"

	req, err := NewRequest("GET", href)
	if err != nil {
		return "", fmt.Errorf("创建详情页请求失败: %v", err)
	}
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求详情页失败: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("详情页状态码异常: %d", resp.StatusCode)
	}
	// 选择紧跟在 <script> 标签后的 <span> 兄弟元素
	click, _ := io.ReadAll(resp.Body)

	return strings.TrimSpace(string(click)), nil
}
func fetchContext(href string) (string, error) {
	//提取文章正文
	req, err := NewRequest("GET", href)
	if err != nil {
		return " ", fmt.Errorf("创建详情页请求失败: %v", err)
	}
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求详情页失败: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("详情页状态码异常: %d", resp.StatusCode)
	}
	doc, _ := goquery.NewDocumentFromReader(resp.Body)
	content := strings.TrimSpace(doc.Find(".v_news_content").Text())
	return content, nil
}
func findcurrURL(baseURL string) string {
	for {
		fmt.Print("正在查询：")
		req, err := NewRequest("GET", baseURL)
		if err != nil {
			log.Fatal("req0创建失败")
		}
		client := http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal("resp0创建失败")
		}
		defer resp.Body.Close()
		doc, _ := goquery.NewDocumentFromReader(resp.Body)
		fmt.Println(doc.Find("body > div.sy-content > div > div.right.fr > div.list.fl > ul > li:nth-child(1) > p > span").Text())
		if doc.Find("body > div.sy-content > div > div.right.fr > div.list.fl > ul > li:nth-child(1) > p > span").Text() <= "2025-12-08" {
			return baseURL
		}
		nexturl, exists := doc.Find(".p_next a").Attr("href")
		if !exists {
			log.Fatal("页面不存在")
		}
		nextURL, _ := url.Parse(nexturl)
		baseURL = resp.Request.URL.ResolveReference(nextURL).String()
	}

}

func main() {
	baseURL := "https://info22.fzu.edu.cn/lm_list.jsp?wbtreeid=1460"

	currentPage := findcurrURL(baseURL)

	db, err := Insert.OpenDB()
	if err != nil {
		log.Fatal(err)
	}

	pageCount := 0
	totalItems := 0

	for {
		pageCount++
		fmt.Printf("正在抓取第 %d 页: %s\n", pageCount, currentPage)

		req, err := NewRequest("GET", currentPage)
		if err != nil {
			log.Printf("创建请求失败: %v", err)
			break
		}

		client := http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("请求失败: %v", err)
			break
		}

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		resp.Body.Close()
		if err != nil {
			log.Printf("解析页面失败: %v", err)
			break
		}

		// 检查页面是否有内容
		if doc.Find("body > div.sy-content > div > div.right.fr > div.list.fl > ul > li:nth-child(1) > p > span").Text() < "2025-11-20" {
			fmt.Println(doc.Find("body > div.sy-content > div > div.right.fr > div.list.fl > ul > li:nth-child(1) > p > span").Text())
			fmt.Println("爬取截止时间到")
			// file.WriteString("OVER\n")
			break
		}

		// 爬取当前页面的通知列表
		// 在主循环中：收集当前页面所有通知的信息

		var items []Insert.Item
		doc.Find(".clearfloat").Each(func(i int, s *goquery.Selection) {
			author := strings.TrimSpace(s.Find(".lm_a").Text())
			title := strings.TrimSpace(s.Find("a").Eq(1).AttrOr("title", ""))
			href := strings.TrimSpace(s.Find("a").Eq(1).AttrOr("href", ""))
			atime := strings.TrimSpace(s.Find(".fr").Text())
			if href != "" {
				items = append(items, Insert.Item{
					Author: author,
					Title:  title,
					Href:   href,
					Atime:  atime,
				})
			}
		})

		// 并发抓取每个通知的点击数和正文
		var wg sync.WaitGroup
		wg.Add(len(items))
		for i := range items {
			go func(i int) {
				defer wg.Done()
				item := &items[i]
				// 提取 clickid 并请求点击数
				re := regexp.MustCompile(`wbnewsid=(\d+)`)
				matches := re.FindStringSubmatch(item.Href)
				if len(matches) > 1 {
					click, err := fetchClick(matches[1])
					if err != nil {
						errMsg := fmt.Sprintf("   详情爬取失败: %v\n", err)
						fmt.Print(errMsg)
						// file.WriteString(errMsg)
					}
					item.Click = click
				}
				// 请求正文内容
				parsedHref, _ := url.Parse(item.Href)
				currenthref := resp.Request.URL.ResolveReference(parsedHref).String()
				content, err := fetchContext(currenthref)
				if err != nil {
					errMsg := fmt.Sprintf("   详情爬取失败: %v\n", err)
					fmt.Print(errMsg)
					// file.WriteString(errMsg)
				}
				item.Content = content
			}(i)
		}
		wg.Wait()
		Insert.Insertintotable(db, items)

		// 寻找下一页链接

		nextPageHref, exists := doc.Find(".p_next a").Attr("href")
		if !exists {
			fmt.Println("已到达最后一页")
			break
		}

		// 将相对路径转换为绝对URL
		nextPageURL, err := url.Parse(nextPageHref)
		if err != nil {
			log.Printf("解析下一页URL出错: %v", err)
			break
		}

		currentPage = resp.Request.URL.ResolveReference(nextPageURL).String()

	}

	// 写入统计信息
	summary := "\n" + strings.Repeat("=", 50) + "\n"
	summary += "爬取完成!\n"
	summary += fmt.Sprintf("总页数: %d\n", pageCount)
	summary += fmt.Sprintf("总条目数: %d\n", totalItems)
	summary += fmt.Sprintf("完成时间: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Print(summary)
}
