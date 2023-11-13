package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/antchfx/jsonquery"
)

// 获取指定链接的body部分内容
func fetch(url1 string) string {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url1, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36")
	cookieValue := "buvid3=A36897B1-7DCD-F8D3-B3F1-4C74AC532E4389438infoc; b_nut=1699160589; CURRENT_FNVAL=4048; _uuid=DF10875710-DFD1-E654-A1A7-AC6DB3129C2D89545infoc; buvid4=165A30B8-48AD-B007-73D4-1045105A2C5B90669-023110513-%2FKPDqUY3StGHPPV6m9xpeg%3D%3D; rpdid=|(JlklRl)~lR0J'uYmmuY|lmu; enable_web_push=DISABLE; header_theme_version=CLOSE; DedeUserID=316938455; DedeUserID__ckMd5=917dbc1e05a26731; home_feed_column=4; browser_resolution=1279-707; bili_ticket=eyJhbGciOiJIUzI1NiIsImtpZCI6InMwMyIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTk4MDA0MjYsImlhdCI6MTY5OTU0MTE2NiwicGx0IjotMX0.WP9JksQWrZrdJS4fwCGxcQKIoXOhjyWhWqYV0ci-gp4; bili_ticket_expires=1699800366; SESSDATA=794fc10d%2C1715093238%2C331a9%2Ab2CjAQQ8WuVcSwWU93J8zyxiSkz5mkCvfUo7PEVjTaEDUoZ_J9CAn-g2zMGC8xpzqCmUoSVkZqSkZ2T0dYSC04Z3pNclZZNnkwck9tNFJ3VENiUnJvZGJxYTM3aWF6WVdqSV9Fang5dE9yZ0l0azVnc0ZkX3JvRTVjb3lfSV9KNzJ0WV8zcmFRMEFBIIEC; bili_jct=e7c3e65bb012cd7a14fd5ec1688d8947; sid=8002h806; PVID=1; bp_video_offset_316938455=862745134624669718; fingerprint=27b0fd1df8f7f4395047fe0215fcbe5d;"
	encodedCookieValue := url.QueryEscape(cookieValue)
	req.Header.Add("Cookie", encodedCookieValue)
	req.Header.Set("Sec-Ch-Ua-Platform", "\"Windows\"")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
	req.Header.Set("Origin", "https://www.bilibili.com")
	req.Header.Set("Referer", "https://www.bilibili.com/")

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Http get err:", err)
		return ""
	}

	if resp.StatusCode != 200 {
		fmt.Println("Http status err:", err)
		return ""
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Read error", err)
		return ""
	}

	return string(body)
}

// 获取指定页码的评论数据保存并返回
func getPage(page int) (txt string) {
	//sort为1代表按点赞数评论
	baseurl := "https://api.bilibili.com/x/v2/reply?&type=1&oid=420981979&pn=" + strconv.Itoa(page) + "&sort=1"
	txt = fetch(baseurl)
	fileName := "第" + strconv.Itoa(page) + "页.txt"
	checkAndWriteToFile(fileName, txt)
	return

}

// 传入子评论的id获取子评论并保存
func getSecondreply(Root string, pn float64) {
	baseURL := "https://api.bilibili.com/x/v2/reply/reply"
	params := url.Values{}
	params.Set("type", "1")
	//oid为视频编号
	params.Set("oid", "450678162")
	params.Set("root", Root)
	params.Set("ps", "20")
	params.Set("pn", "1")

	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.URL.RawQuery = params.Encode()
	req.Header.Set("Cookie", "SESSDATA=xxx")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}
	strPn := strconv.FormatFloat(pn, 'f', -1, 64)
	fileName := "第" + strPn + "页.txt"
	checkAndWriteToFile(fileName, string(body))

}

// 获取子评论的id以及当前评论的页数
func getSecondId(data string) (list []string, pn float64) {
	doc, err := jsonquery.Parse(strings.NewReader(data))
	if err != nil {
		return nil, 0.0
	}
	rootId_list := jsonquery.Find(doc, "/data/replies/*/rpid")
	page := jsonquery.FindOne(doc, "/data/page/num")
	pn = page.Value().(float64)
	list = make([]string, 0)
	for _, va := range rootId_list {
		id := va.Value().(float64)
		strid := strconv.FormatFloat(id, 'f', -1, 64)
		list = append(list, strid)
	}
	return list, pn
}

func checkAndWriteToFile(filename string, content string) error {
	// 检查文件是否存在
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		// 文件不存在，创建文件
		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer file.Close()
	} else if err != nil {
		// 其他错误，返回错误信息
		return err
	}

	// 打开文件以追加形式写入内容
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// 写入内容
	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}

// getPage,不断递增数据来获取该页的传入的json数据格式
func main() {

	for i := 1; i <= 200; i++ {
		data := getPage(i)
		ids, pn := getSecondId(data)
		for _, id := range ids {
			getSecondreply(id, pn)
		}

	}

}
