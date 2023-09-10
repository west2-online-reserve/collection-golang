package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/antchfx/jsonquery"
)

// 爬虫检测阈值：270个goroutine左右
// aaaaa
func main() {
	ch := make(chan int) //同步
	num := 200           //end 目前测出来按照最热排序最大能到498
	start := 1           //start
	for i := start; i <= num; i++ {
		go work(i, ch)
	}
	for i := start; i <= num; i++ {
		<-ch
	}
}

// index 设置到第几页为止
// 工作函数
func work(index int, ch chan int) {
	url := "https://api.bilibili.com/x/v2/reply/main?csrf=f5b57399c62e18808c56eebf2cdfec02&mode=3&oid=420981979&pagination_str=%7B%22offset%22:%22%7B%5C%22type%5C%22:1,%5C%22direction%5C%22:1,%5C%22data%5C%22:%7B%5C%22pn%5C%22:" + strconv.Itoa(index) + "%7D%7D%22%7D&plat=1&type=1"
	data := get(url)
	SaveData(data, index, ch)
	id := getRootId(string(data))
	getReply(id)
	ch <- index
}
func getReply(list []float64) {
	for _, v := range list { //list 楼主rpid，设置url，float64
		url := "https://api.bilibili.com/x/v2/reply/reply?csrf=6f81254eb1291177725c64919659378d&oid=420981979&pn=" + strconv.Itoa(1) + "&ps=10&root=" + strconv.Itoa(int(v)) + "&type=1"
		data := get(url)
		doc, err := jsonquery.Parse(strings.NewReader(string(data)))
		if err != nil {
		}
		total := jsonquery.FindOne(doc, "/data/page/count")
		total1 := int(total.Value().(float64))
		page := total1/10 + func() int { //处理余数
			if total1%10 != 0 {
				return 1
			} else {
				return 0
			}
		}()
		ch1 := make(chan int)
		for i := 1; i <= page; i++ {
			go spiderPage(i, v, ch1)
		}

		for i := 1; i <= page; i++ {
			<-ch1
		}
		fmt.Printf("第 %d 楼的回复保存完成\n", int(v))
	}
}

// 爬取楼猪的回复并保存
func spiderPage(i int, v float64, ch chan int) {
	url := "https://api.bilibili.com/x/v2/reply/reply?csrf=6f81254eb1291177725c64919659378d&oid=420981979&pn=" + strconv.Itoa(i) + "&ps=10&root=" + strconv.Itoa(int(v)) + "&type=1"
	data := get(url)
	path := "./第 " + strconv.Itoa(int(v)) + "楼的回复.txt"
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		println(err)
	}
	defer f.Close()
	n, err2 := f.Write(data)
	if err2 != nil || n == 0 {
		println(err2)
	}
	ch <- i
}

// 获取楼主ID
func getRootId(data string) (list []float64) {
	doc, err := jsonquery.Parse(strings.NewReader(data))
	if err != nil {
		return nil
	}
	rootId_list := jsonquery.Find(doc, "/data/replies/*/rpid")
	list = make([]float64, 0)
	for _, va := range rootId_list {
		list = append(list, va.Value().(float64))
	}
	return list
}

// 判断是否已经到底
func isEnd(data []byte) bool {
	ret := regexp.MustCompile(`"is_end":(.+),"mode":3,"mode_text"`)
	result := ret.FindAllStringSubmatch(string(data), -1)
	for _, v1 := range result {
		if v1[1] == "false" {
			return false
		}
	}
	return true
}

// SaveData 保存某一页的数据
func SaveData(data []byte, index int, ch chan int) bool {
	if isEnd(data) {
		println("The end")
		ch <- index
		return false
	}
	path := "./第 " + strconv.Itoa(index) + " 页的.txt"
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		println(err)
		return false
	}
	defer f.Close()
	n, err2 := f.Write(data)
	if err2 != nil || n == 0 {
		println(err2)
		return false
	}

	fmt.Printf("第 %d 页保存完成\n", index)
	return true
}

// Get 获取数据
func get(url string) (result []byte) {

	var client http.Client
	req := SetRequest(url)
	resp, err := client.Do(req)
	if err != nil {
		println(err)
		return
	}
	result = make([]byte, 0)
	temp := make([]byte, 409600)
	for {
		n, err := resp.Body.Read(temp)
		if n == 0 {
			break
		}
		if err != nil && err != io.EOF {
			// 处理其他错误
			continue
		}
		result = append(result, temp[:n]...)
	}
	return
}

// SetRequest 设置请求头
func SetRequest(url string) (req *http.Request) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		println(err)
		return
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36 Edg/115.0.1901.203")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"Windows\"")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
	req.Header.Set("Origin", "https://www.bilibili.com")
	req.Header.Set("Referer", "https://www.bilibili.com/")
	return
}
