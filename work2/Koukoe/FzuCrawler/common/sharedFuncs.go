package common

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// IsDateInRange 检查日期是否在指定范围内
func IsDateInRange(dateStr string) bool {
	parsed, err := time.Parse("2006-01-02", dateStr) // 太神秘了
	if err != nil {
		fmt.Printf("日期格式错误，跳过: %s\n", dateStr)
		return false
	}

	start, _ := time.Parse("2006-01-02", "2020-01-01")
	end, _ := time.Parse("2006-01-02", "2021-09-01")

	return !parsed.Before(start) && !parsed.After(end)
}

// FetchContent 获取详情页内容
func FetchContent(detailURL string) (string, error) {
	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	req, err := http.NewRequest("GET", detailURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36 Edg/142.0.0.0")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	// 根据真实 HTML 结构查找正文
	contentText := doc.Find("div#vsb_content div.v_news_content").Text()
	if contentText == "" {
		contentText = doc.Find("#vsb_content").Text()
	}
	if contentText == "" {
		contentText = doc.Find(".v_news_content").Text()
	}
	if contentText == "" {
		contentText = doc.Find("body").Text()
	}

	// 清理文本
	re := regexp.MustCompile(`\s+`)
	cleaned := strings.TrimSpace(re.ReplaceAllString(contentText, " "))

	return cleaned, nil
}

// FetchClicks 获取点击量
func FetchClicks(detailURL string) (string, error) {
	re := regexp.MustCompile(`wbnewsid=(\d+)`)
	matches := re.FindStringSubmatch(detailURL)
	if len(matches) < 2 {
		return "", nil
	}
	clickid := matches[1]
	clicksURL := fmt.Sprintf("https://info22.fzu.edu.cn/system/resource/code/news/click/dynclicks.jsp?clickid=%s&owner=1768654345&clicktype=wbnews", clickid)
	client := &http.Client{
		Timeout: 15 * time.Second,
	}
	req, err := http.NewRequest("GET", clicksURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36 Edg/142.0.0.0")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}
	clicks := doc.Find("body").Text()
	return clicks, nil
}
