package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand/v2"
	"net/http"
	"time"
)

const oid = 420981979 // 视频oid

type Reply struct {
	Rpid    int `json:"rpid"` // 评论ID，可能是字符串或数字
	Mid     int `json:"mid"`  // 用户ID，可能是字符串或数字
	Like    int `json:"like"` // 点赞数
	Content struct {
		Message string `json:"message"`
	} `json:"content"` // 评论内容
	Replies []Reply `json:"replies"` // 子评论
	Count   int     `json:"count"`   // 子评论数量
	Member  Member  `json:"member"`  // 用户信息
}
type Member struct {
	Mid   interface{} `json:"mid"`   // 可能是字符串或数字
	Uname string      `json:"uname"` //用户名字
}

// 主评论API响应
type MainReplyResponse struct {
	Code    int      `json:"code"`    //请求是否成功
	Message string   `json:"message"` //响应说明
	Data    struct { //实际响应数据
		Cursor struct {
			IsEnd bool `json:"is_end"` //是否是最后一页
			Next  int  `json:"next"`   //下一页的标识
		} `json:"cursor"`
		Replies []Reply `json:"replies"` //回复的数据切片
	} `json:"data"`
}

// 子评论API响应
type ChildReplyResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Cursor struct {
			IsEnd bool `json:"is_end"`
		} `json:"cursor"`
		Replies []Reply `json:"replies"`
	} `json:"data"`
}

func fetch(url string) ([]byte, error) {
	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// 设置请求头，模拟浏览器
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Referer", "https://www.bilibili.com")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Set("Origin", "https://www.bilibili.com")
	req.Header.Set("DNT", "1")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie",
		"buvid4=0FDDC2EA-F74E-F24F-44D6-0ED6B1F528B277148-022012423-L1Y8EPWwa9dhGSei6IkI7w%3D%3D; buvid_fp_plain=undefined; enable_web_push=DISABLE; CURRENT_BLACKGAP=0; DedeUserID=406757077; DedeUserID__ckMd5=1b6c9f88905884d2; hit-dyn-v2=1; rpdid=|(JYYRlYlm)m0J'u~J)kRYuJ~; enable_feed_channel=ENABLE; fingerprint=9042d82b1e463466f794093be33df918; buvid_fp=9042d82b1e463466f794093be33df918; go-back-dyn=1; header_theme_version=OPEN; theme-tip-show=SHOWED; theme-avatar-tip-show=SHOWED; theme-switch-show=SHOWED; theme_style=light; buvid3=34FB6418-BD76-43C7-66A1-E70A4EAC369E27019infoc; b_nut=1756908227; PVID=4; _uuid=1F8FC10108-41E8-EA10B-42B4-7BCD77D1498E25220infoc; CURRENT_QUALITY=112; ogv_device_support_hdr=0; CURRENT_FNVAL=4048; home_feed_column=4; browser_resolution=225-834; bp_t_offset_406757077=1142101915052539904; b_lsid=10110714A4_19AE494EC2A; bili_ticket=eyJhbGciOiJIUzI1NiIsImtpZCI6InMwMyIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjUwMzA3NTUsImlhdCI6MTc2NDc3MTQ5NSwicGx0IjotMX0._KA0lpowO1I_GP6KPC70VuzJzaIhZ5QTgPs9f1himeA; bili_ticket_expires=1765030695; SESSDATA=94f452b9%2C1780323559%2C52316%2Ac2CjDbodQoPuUm5PU1GdOx9QkQpMemVu7f9WgUHTkkWPEF9vviLTEdMo4YULrqtk8r7FUSVi13SkJSbmt5d0JTbmhzb2RuVEFvMUlSM3BEZXZ6ZURMcE9qa3FPWUdHYVMyM2FiTHJSc2FneWtHQXg3b3IzWGFBaEJKRmgxdjV4TW5SX1h6SDZma2l3IIEC; bili_jct=9d3b816375fcfce23059e8db2deb4261; sid=h0rrzi0j")
	// 随机延迟 1~3 秒
	delay := 1 + rand.Float64()*2
	time.Sleep(time.Duration(delay * float64(time.Second)))

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func fetchParentReplies(next int) ([]Reply, int, bool, error) {
	//设置URL：https://api.bilibili.com/x/v2/reply/main?jsonp=jsonp&type=1&oid=%d&next=%d&mode=3&plat=1 这是b站主评论url
	url := fmt.Sprintf(
		"https://api.bilibili.com/x/v2/reply/main?jsonp=jsonp&type=1&oid=%d&next=%d&mode=3&plat=1",
		oid, next, //接收的数据为视频oid，以及游标
	)

	//接收响应数据
	body, err := fetch(url)
	if err != nil {
		return nil, 0, false, err
	}

	//接收Json数据
	var resp MainReplyResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		log.Printf("JSON解析错误: %v", err)
		return nil, 0, false, err
	}
	if resp.Code != 0 {
		log.Printf("API返回错误: %d, %s", resp.Code, resp.Message)
		return nil, 0, false, fmt.Errorf("API错误: %s", resp.Message)
	}

	//返回Data中的数据，下一页游标，是否到末页
	return resp.Data.Replies, resp.Data.Cursor.Next, resp.Data.Cursor.IsEnd, nil
}

func getAllChildReplies(root int, count int) ([]Reply, error) {
	//接收root该评论下所有子评论数据
	var allChildReplies []Reply

	//提前预知评论数量
	if count == 0 {
		return allChildReplies, nil
	}

	log.Printf("获取评论 %d 的子评论，共 %d 条", root, count)

	pn := 1
	for {
		//设置URL：https://api.bilibili.com/x/v2/reply/reply?jsonp=jsonp&type=1&oid=%d&root=%d&pn=%d&ps=10 这是b站子评论url
		url := fmt.Sprintf(
			"https://api.bilibili.com/x/v2/reply/reply?jsonp=jsonp&type=1&oid=%d&root=%d&pn=%d&ps=10",
			oid, root, pn, //接收数据为视频oid，父评论的rpid，页码
		)

		//爬取评论
		body, err := fetch(url)
		if err != nil {
			return allChildReplies, err
		}

		//接收该页码下的reply数据
		var resp ChildReplyResponse
		if err := json.Unmarshal(body, &resp); err != nil {
			log.Printf("子评论JSON解析错误: %v", err)
			return nil, err
		}
		if resp.Code != 0 {
			log.Printf("子评论API返回错误: %d, %s", resp.Code, resp.Message)
			return nil, fmt.Errorf("API错误: %s", resp.Message)
		}
		//该页子评论
		childReplies := resp.Data.Replies
		//判断是否到末页
		isEnd := resp.Data.Cursor.IsEnd

		//与总子评论合并
		allChildReplies = append(allChildReplies, childReplies...)

		// 如果已经获取了全部子评论或者API返回已结束
		if isEnd || len(childReplies) == 0 || len(allChildReplies) >= count {
			break
		}

		pn++
		// 随机延迟 1~3 秒
		delay := 1 + rand.Float64()*2
		time.Sleep(time.Duration(delay * float64(time.Second)))
	}

	return allChildReplies, nil
}
func main() {
	fmt.Println("开始爬取 B 站评论...")

	//总评论
	var allReplies []Reply
	//游标
	next := 0

	// 获取所有父评论
	for {
		fmt.Printf("\n===== 获取父评论：next = %d =====\n", next)

		replies, nextPage, isEnd, err := fetchParentReplies(next)
		if err != nil {
			log.Printf("获取父评论失败: %v", err)
			break
		}
		fmt.Printf("获取到 %d 条父评论\n", len(replies))

		//设置截断评论函数
		truncateString := func(s string, maxLength int) string {
			if len(s) <= maxLength {
				return s
			}
			return s[:maxLength] + "..."
		}

		// 处理每条父评论
		for i, reply := range replies {
			fmt.Printf("[%d] 父评论 rpid=%d | 用户: %s | 点赞: %d | 内容: %s\n", i+1, reply.Rpid, reply.Member.Uname, reply.Like, reply.Content.Message)
			fmt.Println(reply.Count)
			// 如果有子评论，获取并处理所有子评论
			if reply.Count > 0 {
				childReplies, err := getAllChildReplies(reply.Rpid, reply.Count)
				if err != nil {
					log.Printf("获取评论 %d 的子评论失败: %v", reply.Rpid, err)
				} else {
					reply.Replies = childReplies
					fmt.Printf("\t获取到 %d 条子评论\n", len(childReplies))
					for i, reply := range childReplies {
						fmt.Printf("\t\t[%d] 子评论 rpid=%d | 用户: %s | 点赞: %d | 内容: %s\n", i+1, reply.Rpid, reply.Member.Uname, reply.Like, truncateString(reply.Content.Message, 50))
					}
				}
			}
			//合并到总评论
			allReplies = append(allReplies, reply)
		}
		//父评论到达末页
		if isEnd {
			fmt.Println("父评论已全部抓取完毕")
			break
		}

		next = nextPage
		time.Sleep(500 * time.Millisecond) // 避免请求过快
	}

	fmt.Printf("\n爬取完成，共获取 %d 条父评论\n", len(allReplies))

	// 统计子评论总数
	totalChildReplies := 0
	for _, reply := range allReplies {
		totalChildReplies += len(reply.Replies)
	}
	fmt.Printf("共获取 %d 条子评论\n", totalChildReplies)

}
