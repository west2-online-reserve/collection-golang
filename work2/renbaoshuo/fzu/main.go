package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	baseURL = "https://info22.fzu.edu.cn/"
	listURL = "lm_list.jsp?urltype=tree.TreeTempUrl&wbtreeid=1460"

	fmtWithBaseURL = baseURL + "%s"
	fmtWithListURL = baseURL + listURL + "&totalpage=%d&PAGENUM=%d"

	// 使用 diff 的方式进行计算，以免页面总数变化对程序的影响（前提：历史文件不删除）
	startPageDiffWithLastPage = (1028 - 230)
	endPageDiffWithLastPage   = (1028 - 347)

	// 选择器
	selectorTotalPagesLink = `.p_last a`
	selectorList           = `body > div.sy-content > div > div.right.fr > div.list.fl > ul > li`
	selectorListLink       = `a[href^="content.jsp"]`
	selectorTitle          = `body > div.wa1200w > div.conth > form > div.conth1`
	selectorContent        = `#v_news_content`

	// 正则
	patternTotalPages = `(?s)totalpage=(?P<totalPages>\d+)`
	patternDate       = `(?s)<div class="conthsj" >日期： (?P<date>.*?)  &nbsp;`
	patternAuthor     = `(?s)信息来源..(?P<author>.*?)\n..`
)

var (
	startTime, _ = time.Parse("2006-01-02 15:04:05", "2020-01-01 00:00:00")
	endTime, _   = time.Parse("2006-01-02 15:04:05", "2021-09-30 23:59:59")
)

func httpGet(url string) (string, error) {
	client := &http.Client{}

	resp, err := client.Get(url)
	if err != nil {
		fmt.Printf("Error occurred while getting %s: %s", url, err)

		return "", err
	}

	defer resp.Body.Close()

	buf := make([]byte, 4096)
	res := ""

	for {
		n, _ := resp.Body.Read(buf)
		if n == 0 {
			break
		}
		res += string(buf[:n])
	}

	return res, nil
}

func parseTotalPages(html string) (int, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return 0, err
	}

	href, exists := doc.Find(selectorTotalPagesLink).Attr("href")
	if !exists {
		return 0, fmt.Errorf("no last page link found")
	}
	regexp := regexp.MustCompile(patternTotalPages)
	matches := regexp.FindStringSubmatch(href)
	if len(matches) < 2 {
		return 0, fmt.Errorf("no total page number found")
	}

	totalPages, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, err
	}

	return totalPages, nil
}

func getTotalPages() (int, error) {
	url := fmt.Sprintf(fmtWithBaseURL, listURL)

	html, err := httpGet(url)
	if err != nil {
		return 0, err
	}

	return parseTotalPages(html)
}

func parseListPage(html string) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	res := []string{}

	doc.Find(selectorList).Each(func(i int, s *goquery.Selection) {
		href, exists := s.Find(selectorListLink).Attr("href")
		if !exists {
			return
		}

		pubDate := s.Find("span").Text()
		timePubDate, err := time.Parse("2006-01-02", pubDate)
		if err != nil {
			return
		}

		if timePubDate.Before(startTime) || timePubDate.After(endTime) {
			return
		}

		res = append(res, href)
	})

	return res, nil
}

func getListPage(page int, totalPages int) ([]string, error) {
	url := fmt.Sprintf(fmtWithListURL, totalPages, page)

	html, err := httpGet(url)
	if err != nil {
		return nil, err
	}

	return parseListPage(html)
}

type ArticleResult struct {
	Title    string
	Content  string
	Author   string
	PostDate string // YYYY-MM-DD
}

func parseArticlePage(html string) (ArticleResult, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return ArticleResult{}, err
	}

	title := doc.Find(selectorTitle).Text()
	content := doc.Find(selectorContent).Text()

	regexpDate := regexp.MustCompile(patternDate)
	regexpAuthor := regexp.MustCompile(patternAuthor)

	matchesDate := regexpDate.FindStringSubmatch(html)
	matchesAuthor := regexpAuthor.FindStringSubmatch(html)

	if len(matchesDate) < 2 || len(matchesAuthor) < 2 {
		return ArticleResult{}, fmt.Errorf("no date or author found")
	}

	postDate := matchesDate[1]
	author := matchesAuthor[1]

	return ArticleResult{
		Title:    title,
		Author:   author,
		Content:  content,
		PostDate: postDate,
	}, nil
}

func getArticlePage(urlSuffix string) (ArticleResult, error) {
	url := fmt.Sprintf(fmtWithBaseURL, urlSuffix)

	html, err := httpGet(url)
	if err != nil {
		return ArticleResult{}, err
	}

	return parseArticlePage(html)
}

func main() {
	totalPages, err := getTotalPages()
	if err != nil {
		fmt.Println("Error occurred while getting total pages:", err)
		return
	}

	fmt.Println("Total pages:", totalPages)

	urlList := []string{}

	for page := (totalPages - startPageDiffWithLastPage); page <= (totalPages - endPageDiffWithLastPage); page++ {
		list, err := getListPage(page, totalPages)
		if err != nil {
			fmt.Println("Error occurred while getting list page:", err)
			return
		}

		urlList = append(urlList, list...)
	}

	fmt.Println("Total articles:", len(urlList))

	articles := []ArticleResult{}

	for _, url := range urlList {
		article, err := getArticlePage(url)
		if err != nil {
			fmt.Println("Error occurred while getting article page:", err)
			return
		}

		articles = append(articles, article)
	}

	fmt.Println("Total articles parsed:", len(articles))

	// TODO: Save articles to database

	// print articles basic info
	for _, article := range articles {
		fmt.Printf("%s\t%s\t\t\t%s\n", article.PostDate, article.Author, article.Title)
	}
}
