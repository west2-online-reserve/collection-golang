package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// 获取指定链接的body内容
func fetch(url string) string {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36")
	req.Header.Add("Cookie", "JSESSIONID=03686CCAD2429099B96798BEF5699B37")

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

// 写入文件
func writeFile(filename string, content string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	// 将内容写入文件
	_, err = file.WriteString(content)

	if err != nil {
		return err
	}
	return nil
}

// 获取该页面的所有正文链接,打印链接里的所有需要的内容
func parse(html string) {
	//替换掉换行
	html = strings.Replace(html, "\n", "-3.14", -1)
	//右侧边栏正则
	re_sidebar := regexp.MustCompile(`<div class="list fl">(.*?)<div class="clear"></div>`)
	sidebar := re_sidebar.FindString(html)

	//链接正则
	re_link := regexp.MustCompile(`href="content.(.*?)"`)

	//找到所有链接
	links := re_link.FindAllString(sidebar, -1)

	base_url := "https://info22.fzu.edu.cn/"

	for _, v := range links {
		s := v[6 : len(v)-1]
		url := base_url + s
		go parseNext(url)
		//启动另外一个线程去解析正文链接里的内容
	}

}

// 获取正文并
func parseNext(url string) {
	body := fetch(url)
	time := FindTime(body)
	isTrue, date := parseTime(time)
	fmt.Println(date)
	if !isTrue {
		fmt.Printf("时间不合法:%s\n", date)
		return
	}
	str := "\t\t\t" + FindTitle(body) + "\n时间:" + time + "\n" + Findtxt(body)
	writeFile(FindTitle(body)+".txt", str)

}

// 获取标题
func FindTitle(body string) string {

	//替换换行
	body = strings.Replace(body, "\n", "", -1)
	//页面内容<div class="clear></div>
	re_content := regexp.MustCompile(`<div class="conth">\s*(.*?)\s*<div class="clear"></div>`)
	//找到页面内容
	content := re_content.FindString(body)

	// 使用正则表达式匹配制表符、换行符和空白符
	re := regexp.MustCompile(`[\t\n\s]+`)
	// 将匹配到的制表符、换行符和空白符替换为空字符串
	result := re.ReplaceAllString(content, "")
	//fmt.Println(result)
	//标题
	re_title := regexp.MustCompile(`<divclass="conth1">(.*?)</div>`)
	//找到标题
	match := re_title.FindStringSubmatch(result)
	//fmt.Println(match)

	title := match[1]

	return title

}

// 获取时间
func FindTime(body string) (date string) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		log.Fatal(err)
	}

	date = doc.Find("div.conthsj").Text()

	// 使用正则表达式匹配制表符、换行符和空白符
	re := regexp.MustCompile(`[\t\n\s]+`)
	// 将匹配到的制表符、换行符和空白符替换为空字符串
	result := re.ReplaceAllString(date, "")

	// 使用正则表达式匹配日期
	re = regexp.MustCompile(`日期：([^&]+)信息`)
	match := re.FindStringSubmatch(result)

	if len(match) > 1 {
		match[1] = match[1][:10]
		date = match[1]

	}

	return
}

// 获取作者
// func FindAuth(body string) (auth string) {
// 	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	//Find(".conth").Find(".v_news_content").
// 	s := doc.Find(".conth").Find(".v_news_content").Find("p[style*='text-align: right;'],p[style*='text-align:right;']").First()
// 	s2 := s.Text()

// 	// 输出：党委宣传部

// 	return s2
// }

// 获取正文内容
func Findtxt(body string) (txt string) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	doc1 := doc.Find(`.v_news_content`)
	if err != nil {
		log.Fatal(err)
	}
	txt = doc1.Text()
	return txt
}

// 处理时间函数
func parseTime(t string) (bool, string) {
	layout := "2006-01-02"
	startDateStr := "2020-01-01"
	endDateStr := "2021-09-02"

	startDate, _ := time.Parse(layout, startDateStr)
	endDate, _ := time.Parse(layout, endDateStr)

	inputDate, err := time.Parse(layout, t)
	if err != nil {
		return false, t
	}

	if inputDate.After(startDate) && inputDate.Before(endDate) {
		return true, t
	}

	return false, t
}
func main() {
	url := "https://info22.fzu.edu.cn/lm_list.jsp?totalpage=962&PAGENUM=1&urltype=tree.TreeTempUrl&wbtreeid=1460"
	c := 162
	for {

		url = "https://info22.fzu.edu.cn/lm_list.jsp?totalpage=962&PAGENUM=" + strconv.Itoa(c) + "&urltype=tree.TreeTempUrl&wbtreeid=1460"
		s := fetch(url)
		go parse(s)
		c++
		fmt.Println(c)

	}

}
