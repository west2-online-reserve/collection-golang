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

func main() {
	start := 1
	count := 1 //记录读取的主评论数目
	end := GetLastPage("https://api.bilibili.com/x/v2/reply?type=1&oid=420981979&sort=1&pn=1")
	for i := start; i <= end; i++ {
		url := "https://api.bilibili.com/x/v2/reply?type=1&oid=420981979&sort=1&pn=" + strconv.Itoa(i)
		Spider(url, &count)
		fmt.Printf("爬取完第%d页\n", i)
	}
}

func Spider(url string, count *int) {
	var comment []string
	Alljson := URLtoJson(url)
	doc, err := jsonquery.Parse(strings.NewReader(Alljson))
	if err != nil {
		fmt.Println("parse error", err)
	}
	//data/replies/content/message
	commentNodes := jsonquery.Find(doc, "/data/replies/*/content/message") //评论内容结点
	namesNodes := jsonquery.Find(doc, "/data/replies/*/member/uname")      //用户名结点
	rpidsNodes := jsonquery.Find(doc, "/data/replies/*/rpid")              //获取主评论用户的id用于爬取其下子评论

	for i, va := range commentNodes {
		comment = append(comment, "第"+strconv.Itoa(*count)+"条主评论:\n"+namesNodes[i].Value().(string)+":\n"+va.Value().(string)+"\n\n其子评论:")
		SavaDate(comment[i])
		*count++
		GetSonComments(strconv.Itoa(int(rpidsNodes[i].Value().(float64))))
		SavaDate("\n\n\n")
	}
}

// 计算最后一页的页码
func GetLastPage(url string) int {

	Alljson := URLtoJson(url)
	doc, err := jsonquery.Parse(strings.NewReader(Alljson))
	if err != nil {
		fmt.Println("parse error", err)
	}
	Nums := jsonquery.FindOne(doc, "data/page/count")
	RepliesNums := int(Nums.Value().(float64))
	if RepliesNums == 0 { //一条评论也没有页码为0
		return 0
	}
	if RepliesNums%20 == 0 {
		return RepliesNums / 20
	} else {
		return RepliesNums/20 + 1
	}
}

// 解析出该页的json格式
func URLtoJson(url string) (Alljson string) {
	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("creat request err", err)
	}
	//请求头设置，需要自行补充
	req.Header.Set("authority", "api.bilibili.com")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36 Edg/131.0.0.0")
	req.Header.Set("accept", "*/*")
	req.Header.Set("referer", "https://www.bilibili.com/video/BV12341117rG/?vd_source=78a47b2870883e95e07951ed94d31674")
	req.Header.Set("authority", "api.bilibili.com")
	req.Header.Set("accept-language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
	req.Header.Set("cookie", "xxx") //需要补充cookie
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("request err", err)
	}

	defer resp.Body.Close()

	lines := make([]byte, 4096)

	for {
		n, err2 := resp.Body.Read(lines)
		if n == 0 {
			break
		}
		if err2 != nil && err2 != io.EOF {
			err = err2
			return
		}
		Alljson += string(lines[0:n])
	}
	return
}

// 获取并写入主评论的子评论
func GetSonComments(rpid string) {
	url := "https://api.bilibili.com/x/v2/reply/reply?type=1&oid=420981979&root=" + rpid + "&pn=1"
	end := GetLastPage(url)
	if end == 0 { //一条子评论也没有
		return
	}
	for i := 1; i < end; i++ {
		url = "https://api.bilibili.com/x/v2/reply/reply?type=1&oid=420981979&root=" + rpid + "&pn=" + strconv.Itoa(i)
		Alljsons := URLtoJson(url)
		doc, err := jsonquery.Parse(strings.NewReader(Alljsons))
		if err != nil {
			fmt.Println("parse error", err)
		}
		//data/replies/content/message
		commentNodes := jsonquery.Find(doc, "/data/replies/*/content/message")
		namesNodes := jsonquery.Find(doc, "/data/replies/*/member/uname")
		for i, va := range commentNodes {
			SavaDate(namesNodes[i].Value().(string) + ": " + va.Value().(string)) //在文件中写入子评论
		}
	}

}

// 写入文件
func SavaDate(text string) {
	path := "C:\\code\\go\\spider_demo\\Task02\\Bilibili\\Comment.txt"
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModeAppend|os.ModePerm)
	if err != nil {
		fmt.Println("Open File Failed", err)
	}
	defer f.Close()
	_, err = f.WriteString(text + "\n")
	if err != nil {
		fmt.Println("Write File Failed", err)
	}
}
