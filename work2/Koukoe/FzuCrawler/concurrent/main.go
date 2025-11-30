package main

import (
	"FzuCrawler/common"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// 我写的时候页数是303到420，但是再多发6条通知的话页数要改了
// 是的过了一天现在2020的第一条被挤到421页了
const (
	startPage = 303
	endPage   = 421
	baseURL   = "https://info22.fzu.edu.cn/lm_list.jsp?wbtreeid=1460"
	// 并发爬取详情页的 worker 数量，太高的话首先是点击数爬取失败，然后我不敢试了
	// 正文也会失败，也许是我一开始没看到
	numDetailWorkers = 100
)

type task struct {
	URL    string
	Title  string
	Date   string
	Source string
}

type result struct {
	Title   string
	Date    string
	Source  string
	Content string
	Clicks  string
}

func listWorker(pageNum int, tasks chan<- task, wg *sync.WaitGroup) {
	defer wg.Done()

	client := &http.Client{}
	pageURL := fmt.Sprintf("%s&PAGENUM=%d", baseURL, pageNum)
	req, err := http.NewRequest("GET", pageURL, nil)
	if err != nil {
		fmt.Printf("创建列表页请求失败 (Page %d): %v\n", pageNum, err)
		return
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36 Edg/142.0.0.0")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("请求列表页失败 (Page %d): %v\n", pageNum, err)
		return
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Printf("解析列表页失败 (Page %d): %v\n", pageNum, err)
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
		source := strings.TrimSpace(s.Find("a.lm_a").Text())

		if common.IsDateInRange(date) {
			tasks <- task{
				URL:    detailURL,
				Title:  title,
				Date:   date,
				Source: source,
			}
			fmt.Printf("  [列表Worker] Page %d - 添加任务: %s (日期 %s)\n", pageNum, title, date)
		}
	})
}

func detailWorker(tasks <-chan task, results chan<- result, wg *sync.WaitGroup) {
	defer wg.Done()
	for task := range tasks { // 从 tasks channel 持续读取任务，直到 channel 关闭
		content, err := common.FetchContent(task.URL)
		if err != nil {
			fmt.Printf("获取正文失败 (URL: %s): %v\n", task.URL, err)
			content = "获取正文失败"
		}
		clicks, err := common.FetchClicks(task.URL)
		if err != nil {
			fmt.Printf("获取点击量失败 (URL: %s): %v\n", task.URL, err)
			clicks = "获取点击量失败"
		}

		results <- result{
			Title:   task.Title,
			Date:    task.Date,
			Source:  task.Source,
			Content: content,
			Clicks:  clicks,
		}
	}
}

func crawler(file *os.File) {
	var listWg sync.WaitGroup
	var detailWg sync.WaitGroup

	// 我没算错的话是有2347条通知
	tasks := make(chan task, 2350)
	results := make(chan result, 2350)

	fmt.Printf("启动 %d 个列表页 Worker...\n", endPage-startPage+1)
	for i := startPage; i <= endPage; i++ {
		listWg.Add(1)
		go listWorker(i, tasks, &listWg)
	}

	go func() {
		listWg.Wait()
		close(tasks)
		fmt.Println("所有列表页爬取完成，关闭任务通道")
	}()

	fmt.Printf("启动 %d 个详情页 Worker...\n", numDetailWorkers)
	for w := 0; w < numDetailWorkers; w++ {
		detailWg.Add(1)
		go detailWorker(tasks, results, &detailWg)
	}

	go func() {
		detailWg.Wait()
		close(results)
		fmt.Println("所有详情页处理完成，关闭结果通道")
	}()

	for result := range results {
		fmt.Printf("  [详情Worker] 处理结果: %s (日期 %s)\n", result.Title, result.Date)
		output := fmt.Sprintf("标题: %s\n日期: %s\n来源: %s\n点击量: %s\n正文: %s\n------\n",
			result.Title, result.Date, result.Source, result.Clicks, result.Content)
		_, err := file.WriteString(output)
		if err != nil {
			fmt.Printf("写入文件失败: %v\n", err)
		}
	}
}

func main() {
	start := time.Now()

	file, err := os.Create("FzuNotices.txt") //为什么通知文件系统的URL里写的是news
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	crawler(file)
	elapsed := time.Since(start)
	fmt.Printf("\n用时: %v\n", elapsed)
}
