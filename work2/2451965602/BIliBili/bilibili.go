package main

import (
	"fmt"
	"github.com/antchfx/jsonquery"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

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

func biliacount(date string) (num int) {

	doc, _ := jsonquery.Parse(strings.NewReader(date))
	num1 := jsonquery.FindOne(doc, "data/page/count")
	num = int(num1.Value().(float64))
	return
}

func biliroot(str string) (allroot []float64) {

	doc, _ := jsonquery.Parse(strings.NewReader(str))
	root1 := jsonquery.Find(doc, "/data/replies/*/rpid")
	for _, va := range root1 {
		allroot = append(allroot, va.Value().(float64))
	}
	return
}

func MarkNum(url string) (markNum int) {
	date1, _ := httpget(url)
	markNum = biliacount(date1)
	return
}

func pagerange(url string) (end int) {
	num := MarkNum(url)
	if num%20 == 0 {
		end = num / 20
		return
	} else {
		end = num/20 + 1
		return
	}

}

func saveFile(text []string, index int) {
	for _, str := range text {
		path := "E:\\golang\\src\\GoCode\\Practice\\west2work2\\01\\BIliBili\\result\\" + "第" + strconv.Itoa(index) + "楼.txt"
		f, _ := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModeAppend|os.ModePerm)
		f.WriteString(str + "\n")
		f.Close()
	}
}

func fliter(url string, index int) {
	var conment []string
	date, _ := httpget(url)
	doc, _ := jsonquery.Parse(strings.NewReader(date))
	conment1 := jsonquery.Find(doc, "/data/replies/*/content/message")
	for _, va := range conment1 {
		conment = append(conment, va.Value().(string))
	}
	saveFile(conment, index)
}

func splider(root []float64) {
	for _, i := range root {
		url := "https://api.bilibili.com/x/v2/reply/reply?oid=420981979&type=1&root=" + strconv.Itoa(int(i)) + "&ps=20&pn=1"
		start := 1
		end := pagerange(url)
		for k := start; k <= end; k++ {
			url2 := "https://api.bilibili.com/x/v2/reply/reply?oid=420981979&type=1&root=" + strconv.Itoa(int(i)) + "&ps=20&pn=" + strconv.Itoa(k)
			fliter(url2, int(i))
		}

	}
}
