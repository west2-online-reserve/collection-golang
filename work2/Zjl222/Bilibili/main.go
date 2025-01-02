package main

import (
	"fmt"
	"github.com/antchfx/jsonquery" //用于json查询的解析库
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

//httpget 函数发送 HTTP GET 请求并返回响应内容
func httpget(url string) (result string, err error) {

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36 Edg/115.0.1901.203")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"Windows\"")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
	req.Header.Set("Origin", "https://www.bilibili.com")
	req.Header.Set("Referer", "https://www.bilibili.com/")
	resp, err1 := client.Do(req)
	if err1 != nil {
		fmt.Println("httpget出错")
		err = err1
		return
	}
	defer resp.Body.Close()

	buf := make([]byte, 4096)
	for {
		n, err2 := resp.Body.Read(buf)
		if n == 0 {
			break
		}
		if err2 != nil && err2 != io.EOF {
			err = err2
			return
		}
		result += string(buf[0:n])
	}

	return
}

//获取评论总数
func commentsum(date string) (num int) {

	doc, _ := jsonquery.Parse(strings.NewReader(date))
	num1 := jsonquery.FindOne(doc, "data/page/count")
	num = int(num1.Value().(float64))
	return
}

//从 JSON 响应中提取所有评论的根 ID
func rootId(str string) (allroot []float64) {

	doc, _ := jsonquery.Parse(strings.NewReader(str))
	root1 := jsonquery.Find(doc, "/data/replies/*/rpid")
	for _, va := range root1 {
		allroot = append(allroot, va.Value().(float64))
	}
	return
}


//获取该url下的评论总数
func commentsum1(url string) (sum int) {
	date1, _ := httpget(url)
	sum = commentsum(date1)
	return
}

//获取评论总页数
func pagesum(url string) (end int) {
	num := commentsum1(url)
	if num%20 == 0 {
		end = num / 20
		return
	} else {
		end = num/20 + 1
		return
	}

}

//保存评论内容
func savecomment(text []string, index int) {
	for _, str := range text {
		path := "D:\\GO语言--gopath\\src\\github.com\\" + "第" + strconv.Itoa(index) + "楼.txt"
		f, _ := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModeAppend|os.ModePerm)
		f.WriteString(str + "\n")
		f.Close()
	}
}

//获取指定url下的评论内容
func fliter(url string, index int) {
	var conment []string
	date, _ := httpget(url)
	doc, _ := jsonquery.Parse(strings.NewReader(date))
	conment1 := jsonquery.Find(doc, "/data/replies/*/content/message")
	for _, va := range conment1 {
		conment = append(conment, va.Value().(string))
	}
	savecomment(conment, index)
}

//获取子评论
func commentson(root []float64) {
	for _, i := range root {
		url := "https://api.bilibili.com/x/v2/reply/reply?oid=420981979&type=1&root=" + strconv.Itoa(int(i)) + "&ps=20&pn=1"
		start := 1
		end := pagesum(url)
		for k := start; k <= end; k++ {
			url2 := "https://api.bilibili.com/x/v2/reply/reply?oid=420981979&type=1&root=" + strconv.Itoa(int(i)) + "&ps=20&pn=" + strconv.Itoa(k)
			fliter(url2, int(i))
		}

	}
}

func main() {
	var start, end int
	start = 1
	end = pagesum("https://api.bilibili.com/x/v2/reply?type=1&oid=420981979&sort=1")
	for i := start; i <= end; i++ {
		url2 := "https://api.bilibili.com/x/v2/reply?type=1&oid=420981979&sort=1&pn=" + strconv.Itoa(i)
		date2, _ := httpget(url2)
		root := rootId(date2)
		commentson(root)
	}

}