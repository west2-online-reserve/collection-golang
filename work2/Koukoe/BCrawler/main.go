package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

const (
	// ReplyURL is 评论区URL, oid是av号, type是1表示视频
	ReplyURL = "https://api.bilibili.com/x/v2/reply?oid=420981979&type=1"
	// ReReplyURL is 评论的回复URL, root后面接主评论rpid
	ReReplyURL = "https://api.bilibili.com/x/v2/reply/reply?oid=420981979&type=1&root="
	// RequestInterval 请求频率间隔
	RequestInterval = 50 * time.Millisecond
)

// APIresponse 返回根对象
type APIresponse struct {
	Code    int       `json:"code"`    //返回值 0为成功
	Message string    `json:"message"` //错误信息
	Data    ReplyData `json:"data"`    //数据本体
}

// ReplyData 数据本体
type ReplyData struct {
	Page    PageInfo `json:"page"`    //页信息
	Replies []Reply  `json:"replies"` //评论列表
}

// PageInfo 页信息
type PageInfo struct {
	Num   int `json:"num"`   //当前页码
	Size  int `json:"size"`  //每页数量
	Count int `json:"count"` //总评论数
}

// Reply 评论结构体
type Reply struct {
	Rpid    int64        `json:"rpid"`    //评论ID
	Root    int64        `json:"root"`    //根评论ID，0表示为根评论
	Rcount  int          `json:"rcount"`  //回复数
	Ctime   int64        `json:"ctime"`   //评论时间戳
	Content ReplyContent `json:"content"` //评论内容
	Like    int          `json:"like"`    //点赞数
	Member  MemberInfo   `json:"member"`  //评论用户信息
	Replies []Reply      `json:"replies"` //楼中楼回复列表
}

// MemberInfo 用户信息结构体
type MemberInfo struct {
	Mid   string `json:"mid"`   //用户ID
	Uname string `json:"uname"` //用户名
	Sex   string `json:"sex"`   //性别
}

// ReplyContent 评论内容结构体
type ReplyContent struct {
	Message string `json:"message"` //评论文本内容
}

func printReplies(replies []Reply, page PageInfo, file *os.File) {
	var indent string
	fmt.Printf("当前页码: %d, 每页数量: %d, 总评论数: %d\n", page.Num, page.Size, page.Count)
	for n, reply := range replies {
		time.Sleep(RequestInterval)
		if reply.Root != 0 {
			indent = "  " // 如果不是根评论，增加缩进
			fmt.Println("  ↑这里是在说子评论↓")
		}
		fmt.Printf("正在处理第 %d 条评论...\n", n+1)
		t := time.Unix(reply.Ctime, 0)
		fmt.Fprintf(file, "%s用户名: %s  性别: %s\n%s评论: %s\n%s时间: %s, 点赞数: %d, 回复数: %d\n", indent, reply.Member.Uname, reply.Member.Sex, indent, reply.Content.Message, indent, t.Format("2006-01-02 15:04:05"), reply.Like, reply.Rcount)
		fmt.Fprintln(file, "------")
		// 如果有子评论
		if reply.Rcount > 0 {
			rrURL := fmt.Sprintf("%s%d", ReReplyURL, reply.Rpid) // 构造楼中楼请求URL
			crawler(rrURL, file)
		}
	}
}

func crawler(url string, file *os.File) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	for i := 1; ; i++ {
		time.Sleep(RequestInterval)
		fmt.Printf("正在爬取第 %d 页评论...\n", i)
		pageURL := fmt.Sprintf("%s&pn=%d", url, i)

		req, err := http.NewRequest("GET", pageURL, nil)
		if err != nil {
			fmt.Println("Error creating the request:", err)
			return
		}
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36 Edg/142.0.0.0")
		//这里填写Cookie里的SESSDATA
		req.Header.Set("Cookie", "SESSDATA=")
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error making the request:", err)
			return
		}
		defer resp.Body.Close()

		var apiResp APIresponse
		err = json.NewDecoder(resp.Body).Decode(&apiResp)
		if err != nil {
			fmt.Println("Error decoding JSON response:", err)
			return
		}

		if apiResp.Code != 0 {
			fmt.Println("API returned error code:", apiResp.Code, "message:", apiResp.Message)
			return
		}
		//子评论的pn似乎是要大于501才会返回非0的code
		if (apiResp.Data.Page.Num-1)*apiResp.Data.Page.Size > apiResp.Data.Page.Count {
			fmt.Println("  子评论爬取完毕")
			return
		}

		printReplies(apiResp.Data.Replies, apiResp.Data.Page, file)
	}
}

func main() {
	start := time.Now()

	file, err := os.Create("BiliReplies.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	crawler(ReplyURL, file)

	elapsed := time.Since(start)
	fmt.Printf("程序运行时间: %s\n", elapsed)
}
