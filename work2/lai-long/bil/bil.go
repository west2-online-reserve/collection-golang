package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

// 接收主评论的结构体
type Result struct {
	Code int64    `json:"code"` //状态
	Data struct { //数据
		Cursor struct {
			AllCount        int64  `json:"all_count"` //评论总数
			IsBegin         bool   `json:"is_begin"`  //是否为第一页
			IsEnd           bool   `json:"is_end"`    //是否为最后一页
			Mode            int64  `json:"mode"`
			ModeText        string `json:"mode_text"`
			Name            string `json:"name"`
			Next            int64  `json:"next"`
			PaginationReply struct {
				NextOffset string `json:"next_offset"` //下一页
			} `json:"pagination_reply"`
			Prev        int64   `json:"prev"`
			SessionID   string  `json:"session_id"`
			SupportMode []int64 `json:"support_mode"`
		} `json:"cursor"`
		Replies []struct {
			Content struct {
				Device  string        `json:"device"`
				JumpURL struct{}      `json:"jump_url"`
				MaxLine int64         `json:"max_line"`
				Members []interface{} `json:"members"`
				Message string        `json:"message"`
				Plat    int64         `json:"plat"`
			} `json:"content"`
			Count  int64 `json:"count"`
			Rpid   int64 `json:"rpid"` //id
			Folder struct {
				HasFolded bool   `json:"has_folded"`
				IsFolded  bool   `json:"is_folded"`
				Rule      string `json:"rule"`
			} `json:"folder"`
			Like    int64 `json:"like"`
			Replies []struct {
				Action  int64 `json:"action"`
				Assist  int64 `json:"assist"`
				Attr    int64 `json:"attr"`
				Content struct {
					Device  string   `json:"device"`
					JumpURL struct{} `json:"jump_url"`
					MaxLine int64    `json:"max_line"`
					Message string   `json:"message"`
					Plat    int64    `json:"plat"`
				} `json:"content"`
				Rcount  int64       `json:"rcount"`
				Replies interface{} `json:"replies"`
			} `json:"replies"`
			Type int64 `json:"type"`
		} `json:"replies"`
	} `json:"data"`
	Message string `json:"message"`
}

// 接收二级评论的结构体
type Results struct {
	Code int64 `json:"code"`
	Data struct {
		Cursor struct {
			AllCount        int64  `json:"all_count"`
			IsBegin         bool   `json:"is_begin"`
			IsEnd           bool   `json:"is_end"`
			Mode            int64  `json:"mode"`
			ModeText        string `json:"mode_text"`
			Name            string `json:"name"`
			Next            int64  `json:"next"`
			PaginationReply struct {
				NextOffset string `json:"next_offset"`
			} `json:"pagination_reply"`
			Prev        int64   `json:"prev"`
			SessionID   string  `json:"session_id"`
			SupportMode []int64 `json:"support_mode"`
		} `json:"cursor"`
		Replies []struct {
			Content struct {
				Device  string        `json:"device"`
				JumpURL struct{}      `json:"jump_url"`
				MaxLine int64         `json:"max_line"`
				Members []interface{} `json:"members"`
				Message string        `json:"message"`
				Plat    int64         `json:"plat"`
			} `json:"content"`
			Count  int64 `json:"count"`
			Folder struct {
				HasFolded bool   `json:"has_folded"`
				IsFolded  bool   `json:"is_folded"`
				Rule      string `json:"rule"`
			} `json:"folder"`
			Like int64 `json:"like"`
			Type int64 `json:"type"`
		} `json:"replies"`
	} `json:"data"`
	Message string `json:"message"`
}

var (
	all_count    int64
	total_pages  int
	current_page int
)

func main() {
	fmt.Println("begin")
	current_page = 1
	nextOffset := CreatUrlAndSpider("")
	for nextOffset != "" {
		current_page++
		time.Sleep(5 * time.Second)
		nextOffset = CreatUrlAndSpider(nextOffset)
	}
}

// 进行爬取
func CreatUrlAndSpider(offset string) string {
	//生成评论所在的URL
	wts := strconv.FormatInt(time.Now().Unix(), 10)
	var paginationStr string
	//获取pn
	if offset == "" {
		paginationStr = "%7B%22offset%22%3A%22%22%7D"
	} else {
		paginationStr = "%7B%22offset%22%3A%22" + offset + "%22%7D"
	}
	//获取wrid
	v := "mode=3&oid=420981979&pagination_str=" + paginationStr + "&plat=1&seek_rpid=&type=1&web_location=1315875&wts=" + wts
	a := "ea1db124af3c7062474693fa704f4ff8"
	data := []byte(v + a)
	has := md5.Sum(data)
	w_rid := fmt.Sprintf("%x", has)
	var Url string
	//拼接url
	if offset == "" {
		Url = "https://api.bilibili.com/x/v2/reply/wbi/main?oid=420981979&type=1&mode=3&pagination_str=%7B%22offset%22:%22%22%7D&plat=1&seek_rpid=&web_location=1315875&w_rid=" + w_rid + "&wts=" + wts
	} else {
		Url = "https://api.bilibili.com/x/v2/reply/wbi/main?oid=420981979&type=1&mode=3&pagination_str=%7B%22offset%22:%22" + offset + "%22%7D&plat=1&web_location=1315875&w_rid=" + w_rid + "&wts=" + wts
	}
	fmt.Printf("Page %d---------------\n", current_page)
	return GetMessage(Url)
}
func GetMessage(str string) string {
	client := &http.Client{}
	//发送请求
	req, err := http.NewRequest("GET", str, nil)
	if err != nil {
		fmt.Println("req1 err:", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36")
	//获取返回值
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("resp1 err:", err)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("io err:", err)
	}
	var result Result
	err = json.Unmarshal(bodyText, &result)
	if err != nil {
		fmt.Println("json err:", err)
	}
	if all_count == 0 {
		all_count = result.Data.Cursor.AllCount
		total_pages = int(all_count) / 20
		fmt.Printf("All comments: %d, pages: %d\n", all_count, total_pages)
	}
	fmt.Println("一级：", len(result.Data.Replies))
	for _, result := range result.Data.Replies {
		rpid := result.Rpid
		count := result.Count
		fmt.Printf("主评论: %s\n", result.Content.Message)
		if count > 0 {
			fmt.Println("二级：", count)
			GetSonComment(rpid, int(count))
			fmt.Printf("----------------------")
		}
		fmt.Println()
	}
	var nextOffset string
	if result.Data.Cursor.PaginationReply.NextOffset != "" && !result.Data.Cursor.IsEnd {
		nextOffset = result.Data.Cursor.PaginationReply.NextOffset
	} else {
		nextOffset = ""
		fmt.Println("ended")
	}
	return nextOffset
}

// 通过主评论id获取子评论内容
func GetSonComment(rootRpid int64, totalCount int) {
	pages := (totalCount + 9) / 10
	for page := 1; page <= pages; page++ {
		url := fmt.Sprintf("https://api.bilibili.com/x/v2/reply/reply?oid=420981979&type=1&root=%d&ps=10&pn=%d&web_location=333.788", rootRpid, page)
		time.Sleep(5 * time.Second)
		client := &http.Client{}
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println("req sonComment err", err)
		}
		req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36")
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("req sonComment err", err)
		}
		bodyText, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			fmt.Println("resp sonComment err", err)
		}
		var resultList Results
		err = json.Unmarshal(bodyText, &resultList)
		if err != nil {
			fmt.Println("JSON sonComments err:", err)
		}
		for i, reply := range resultList.Data.Replies {
			fmt.Printf("二级 %d: %s\n", (page-1)*10+i+1, reply.Content.Message)
		}
	}
}
