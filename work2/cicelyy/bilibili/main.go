package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/antchfx/jsonquery"
)

// httpget 发送 HTTP GET 请求并返回响应体内容
func httpget(url string) (result string, err error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36 Edg/115.0.1901.203")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"Windows\"")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
	req.Header.Set("Origin", "https://www.bilibili.com")
	req.Header.Set("Referer", "https://www.bilibili.com/")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("httpget failed with status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// biliacount 解析 JSON 数据并返回评论数量
func biliacount(date string) (num int, err error) {
	if date == "" {
		return 0, fmt.Errorf("date string is empty")
	}
	doc, err := jsonquery.Parse(strings.NewReader(date))
	if err != nil {
		return 0, err
	}
	num1 := jsonquery.FindOne(doc, "data/page/count")
	if num1 == nil {
		return 0, fmt.Errorf("cannot find count in date string")
	}
	return int(num1.Value().(float64)), nil
}

// biliroot 解析 JSON 数据并返回所有根评论的 ID
func biliroot(str string) (allroot []float64, err error) {
	if str == "" {
		return nil, fmt.Errorf("date string is empty")
	}
	if !strings.HasPrefix(str, "{") {
		return nil, fmt.Errorf("response is not JSON")
	}
	doc, err := jsonquery.Parse(strings.NewReader(str))
	if err != nil {
		return nil, err
	}
	root1 := jsonquery.Find(doc, "/data/replies/*/rpid")
	if root1 == nil {
		return nil, fmt.Errorf("cannot find root in date string")
	}
	for _, va := range root1 {
		allroot = append(allroot, va.Value().(float64))
	}
	return allroot, nil
}

// MarkNum 获取评论数量
func MarkNum(url string) (markNum int, err error) {
	date1, err := httpget(url)
	if err != nil {
		return 0, err
	}
	markNum, err = biliacount(date1)
	if err != nil {
		return 0, err
	}
	return markNum, nil
}

// pagerange 计算需要爬取的页数范围
func pagerange(url string) (end int, err error) {
	num, err := MarkNum(url)
	if err != nil {
		return 0, err
	}
	if num%20 == 0 {
		end = num / 20
	} else {
		end = num/20 + 1
	}
	return end, nil
}

// saveFile 将评论保存到文件
func saveFile(text []string, index int) {
	for _, str := range text {
		path := "E:\\golang\\src\\GoCode\\Practice\\west2work2\\01\\BIliBili\\result\\" + "第" + strconv.Itoa(index) + "楼.txt"
		f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModeAppend|os.ModePerm)
		if err != nil {
			fmt.Println("Error opening file:", err)
			continue
		}
		defer f.Close()
		_, err = f.WriteString(str + "\n")
		if err != nil {
			fmt.Println("Error writing to file:", err)
		}
	}
}

// fliter 获取评论并保存
func fliter(url string, index int) {
	var conment []string
	date, err := httpget(url)
	if err != nil {
		fmt.Println("Error fetching page:", err)
		return
	}
	if !strings.HasPrefix(date, "{") {
		fmt.Println("Response is not JSON:", date)
		return
	}
	doc, err := jsonquery.Parse(strings.NewReader(date))
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}
	conment1 := jsonquery.Find(doc, "/data/replies/*/content/message")
	if conment1 == nil {
		fmt.Println("No comments found")
		return
	}
	for _, va := range conment1 {
		conment = append(conment, va.Value().(string))
	}
	saveFile(conment, index)
}

// splider 递归获取所有子评论
func splider(root []float64) {
	for _, i := range root {
		url := "https://api.bilibili.com/x/v2/reply/reply?oid=420981979&type=1&root=" + strconv.Itoa(int(i)) + "&ps=20&pn=1"
		start := 1
		end, err := pagerange(url)
		if err != nil {
			fmt.Println("Error getting page range:", err)
			continue
		}
		for k := start; k <= end; k++ {
			url2 := "https://api.bilibili.com/x/v2/reply/reply?oid=420981979&type=1&root=" + strconv.Itoa(int(i)) + "&ps=20&pn=" + strconv.Itoa(k)
			fliter(url2, int(i))
		}
	}
}

func main() {
	// B站视频页面URL
	url := "https://www.bilibili.com/video/BV12341117rG/?vd_source=8ab12435fa47b222be4f0de0b7d79be0"

	// 获取页面内容
	date, err := httpget(url)
	if err != nil {
		fmt.Println("Error fetching page:", err)
		return
	}

	// 提取评论根ID
	root, err := biliroot(date)
	if err != nil {
		fmt.Println("Error extracting root IDs:", err)
		return
	}

	// 递归获取所有子评论
	splider(root)
}
