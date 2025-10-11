package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

const (
	DefaultMaxPage    = 1090
	DefaultMaxRoutine = 100
)

/* 并行爬取fzu通信文件系统 */
func crawlFzuAsync(maxPageCount int, maxRoutineCount int) []article {
	if maxRoutineCount == 0 {
		maxRoutineCount = DefaultMaxRoutine
	}
	if maxPageCount == 0 {
		maxPageCount = DefaultMaxPage
	}

	var wg sync.WaitGroup
	articleCh := make(chan []article, maxPageCount)   // 管道用于接收每一页的数据并汇总
	routineCh := make(chan struct{}, maxRoutineCount) // 管道用于限制最大并行量

	for page := 1; page <= maxPageCount; page++ {
		routineCh <- struct{}{}
		wg.Add(1)
		go func(page int) {
			defer func() { <-routineCh }()
			defer wg.Done()
			defer fmt.Printf("第 %d 页数据爬取完成！\n", page)
			articleCh <- crawlPage(page, maxPageCount)
		}(page)
	}

	// 阻塞数据管道关闭，直到全部爬取完。否则range管道无法执行
	go func() {
		wg.Wait()
		close(articleCh)
	}()

	// 阻塞主routine，等待最后几个协程完成爬取
	for i := 0; i < cap(routineCh); i++ {
		routineCh <- struct{}{}
	}

	// 汇总所有文章数据
	var all []article
	for articles := range articleCh {
		for _, v := range articles {
			all = append(all, v)
		}
	}
	return all
}

/* 串行爬取fzu通信文件系统 */
func crawlFzu(maxPageCount int) []article {
	if maxPageCount == 0 {
		maxPageCount = DefaultMaxPage
	}

	var articles []article

	for page := 1; page <= maxPageCount; page++ {
		pageArticles := crawlPage(page, maxPageCount)
		for _, a := range pageArticles {
			articles = append(articles, a)
		}
		fmt.Printf("第 %d 页数据爬取完成！\n", page)
	}

	return articles
}

/* 爬取第page页 */
func crawlPage(page int, totalPage int) []article {
	// 发送请求
	pageStr := strconv.Itoa(page)
	totalPageStr := strconv.Itoa(totalPage)
	req, err := http.NewRequest("GET", "https://info22.fzu.edu.cn/lm_list.jsp?totalpage="+totalPageStr+"&PAGENUM="+pageStr+"&wbtreeid=1460", nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/140.0.0.0 Safari/537.36")

	client := http.Client{}
	var response *http.Response
	for i := 0; i < 3; i++ {
		response, err = client.Do(req)
		if err == nil {
			break
		}
	}
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	// 解析网页
	docDetail, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var articles []article

	docDetail.Find("body > div.sy-content > div > div.right.fr > div.list.fl > ul > li").
		Each(func(i int, s *goquery.Selection) {
			author := s.Find("p > a.lm_a").Text()
			a2 := s.Find("p > a:nth-child(2)")
			href, _ := a2.Attr("href") // 正文链接，给的是相对地址，需要进一步处理
			title, _ := a2.Attr("title")
			date := s.Find("p > span").Text()

			ymd := strings.Split(date, "-")
			year, _ := strconv.Atoi(ymd[0])
			month, _ := strconv.Atoi(ymd[1])
			day, _ := strconv.Atoi(ymd[2])
			// 只爬取 2020-01-01 ~ 2021-09-01 的文章
			if year == 2020 || (year == 2021 && month < 9) || (year == 2021 && month == 9 && day == 1) {
				norm := normalizeHref("https://info22.fzu.edu.cn", href)

				// 继续请求正文
				req, err := http.NewRequest("GET", norm, nil)
				if err != nil {
					log.Fatal(err)
				}
				req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/140.0.0.0 Safari/537.36")
				response, err := client.Do(req)
				if err != nil {
					log.Fatal(err)
				}
				defer response.Body.Close()

				docDetail, err = goquery.NewDocumentFromReader(response.Body)
				if err != nil {
					log.Fatal(err)
				}

				// body > div.wa1200w > div.conth > form > div.conthsj > span
				// body > div.wa1200w > div.conth > form > div.conthsj > script
				clicksStr := docDetail.Find("body > div.wa1200w > div.conth > form > div.conthsj > span").Text()
				clicks, _ := strconv.Atoi(clicksStr)

				// #vsb_content > div > p:nth-child(1)
				// #vsb_content > div > p:nth-child(2)
				// #vsb_content > div > p:nth-child(21)
				// 正文body分成了若干段，全部join起来
				var seg []string
				docDetail.Find("#vsb_content > div > p").
					Each(func(i int, s *goquery.Selection) {
						seg = append(seg, s.Find("span").Text())
					})
				body := strings.Join(seg, "\n")

				data := article{date, author, title, body, clicks}
				articles = append(articles, data)
			}
		})
	return articles
}

/* 标准化url地址 */
func normalizeHref(base string, href string) string {
	href = strings.TrimSpace(href)
	b, _ := url.Parse(base)
	u, _ := url.Parse(href)
	abs := b.ResolveReference(u) // 拼接根链接和相对地址
	return abs.String()
}
