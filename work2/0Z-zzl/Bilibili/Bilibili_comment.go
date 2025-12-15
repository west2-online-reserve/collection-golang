package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

// DedeUserID=449728349
// SESSDATA=ffe1880d%2C1781137300%2C31d91%2Ac1CjAr-dkihLyhqQGtvDf-rucVJOvNTPyx-o8el50wFl_AiRUKwuifPiSIzYTAHshQXsYSVjRvVTdnUWJMNHEwUGM0bWY2d1ZvS1lvU003TWdmbkJHQm1VWHdaY2V5MXJyZEVKSUswWG9YWVVGY0FWa3lFblRmaWJOakFvYTUyZ3lQRnBnR3dZdEZRIIEC;
// bili_jct=44cda4a655b36998e2267002dca671ff
const (
	BV_ID          = "BV12341117rG"
	COOKIE         = "SESSDATA=ffe1880d%2C1781137300%2C31d91%2Ac1CjAr-dkihLyhqQGtvDf-rucVJOvNTPyx-o8el50wFl_AiRUKwuifPiSIzYTAHshQXsYSVjRvVTdnUWJMNHEwUGM0bWY2d1ZvS1lvU003TWdmbkJHQm1VWHdaY2V5MXJyZEVKSUswWG9YWVVGY0FWa3lFblRmaWJOakFvYTUyZ3lQRnBnR3dZdEZRIIEC; bili_jct=44cda4a655b36998e2267002dca671ff; DedeUserID=449728349;"
	USER_AGENT     = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
	REQUEST_DELAY  = 5000 * time.Millisecond // 延长至2秒，降低风控概率
	MAX_PAGE_LIMIT = 5000                    // 最大页码限制，避免超限
)

// 获取评论总数的响应结构体
type ReplyCountResponse struct {
	Code int `json:"code"`
	Data struct {
		Count int `json:"count"` // 评论总数
	} `json:"data"`
}

// BV转AID的响应结构体
type ViewResponse struct {
	Code int `json:"code"`
	Data struct {
		Aid int `json:"aid"`
	} `json:"data"`
}

// 评论接口响应结构体
type ReplyResponse struct {
	Code int `json:"code"`
	Data struct {
		Cursor struct {
			Pn    int  `json:"pn"`     // 当前页码
			IsEnd bool `json:"is_end"` // 是否最后一页
		} `json:"cursor"`
		Replies []Comment `json:"replies"` // 评论列表
	} `json:"data"`
}

// 评论结构体
type Comment struct {
	Rpid   int64 `json:"rpid"`
	Member struct {
		Uname string `json:"uname"`
	} `json:"member"`
	Content struct {
		Message string `json:"message"` // 评论内容
	} `json:"content"`
	Reply int   `json:"reply"` // 子评论数
	Ctime int64 `json:"ctime"` // 发布时间
}

func getCommentCount(aid int) (int, error) {
	url := fmt.Sprintf("https://api.bilibili.com/x/v2/reply/count?oid=%d&type=1", aid)
	time.Sleep(REQUEST_DELAY)
	req, _ := http.NewRequest("GET", url, nil)
	setReqHeader(req)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("获取评论总数失败：%v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var countResp ReplyCountResponse
	if err := json.Unmarshal(body, &countResp); err != nil {
		return 0, fmt.Errorf("解析评论总数失败：%v", err)
	}
	if countResp.Code != 0 {
		return 0, fmt.Errorf("评论总数接口返回错误：code=%d", countResp.Code)
	}
	return countResp.Data.Count, nil
}

// BV转AID
func getAidByBV(bv string) (int, error) {
	url := fmt.Sprintf("https://api.bilibili.com/x/web-interface/view?bvid=%s", bv)
	req, _ := http.NewRequest("GET", url, nil)
	setReqHeader(req)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("BV转AID请求失败：%v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var viewResp ViewResponse
	if err := json.Unmarshal(body, &viewResp); err != nil {
		return 0, fmt.Errorf("解析AID失败：%v", err)
	}
	if viewResp.Code != 0 {
		return 0, fmt.Errorf("BV转AID接口返回错误：code=%d", viewResp.Code)
	}
	return viewResp.Data.Aid, nil
}

// 设置请求头
func setReqHeader(req *http.Request) {
	req.Header.Set("User-Agent", USER_AGENT)
	req.Header.Set("Cookie", COOKIE)
	req.Header.Set("Referer", fmt.Sprintf("https://www.bilibili.com/video/%s/", BV_ID))
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Origin", "https://www.bilibili.com")
}

// 请求主评论列表
func getMainComments(aid int, pn int) (*ReplyResponse, error) {
	url := fmt.Sprintf("https://api.bilibili.com/x/v2/reply/main?oid=%d&type=1&mode=3&pn=%d&ps=20", aid, pn)
	return requestComment(url)
}

// 请求子评论列表
func getSubComments(aid int, rootRpid int64, pn int) (*ReplyResponse, error) {
	url := fmt.Sprintf("https://api.bilibili.com/x/v2/reply/reply?oid=%d&type=1&root=%s&pn=%d&ps=20&mode=1",
		aid, strconv.FormatInt(rootRpid, 10), pn)
	time.Sleep(REQUEST_DELAY)
	return requestComment(url)
}

// 通用评论请求方法
func requestComment(url string) (*ReplyResponse, error) {
	time.Sleep(REQUEST_DELAY)
	req, _ := http.NewRequest("GET", url, nil)
	setReqHeader(req)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败：%v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var replyResp ReplyResponse
	if err := json.Unmarshal(body, &replyResp); err != nil {
		return nil, fmt.Errorf("解析评论失败：%v，响应体：%s", err, string(body))
	}
	if replyResp.Code != 0 {
		return &replyResp, fmt.Errorf("API返回错误：code=%d，响应体：%s", replyResp.Code, string(body))
	}
	return &replyResp, nil
}

func main() {
	// 1. BV转AID
	aid, err := getAidByBV(BV_ID)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("成功获取视频AID：%d\n", aid)

	// 2. 获取评论总数
	totalCount, err := getCommentCount(aid)
	if err != nil {
		fmt.Printf("获取评论总数失败（不影响爬取）：%v\n", err)
	} else {
		fmt.Printf("视频总评论数：%d条，预计总页数：%d页\n", totalCount, (totalCount+19)/20) // 每页20条，向上取整
	}

	// 3. 爬取所有主评论+子评论
	var allComments []Comment
	pn := 1 // 主评论页码
	for {
		// 触发最大页码限制，停止爬取
		if pn > MAX_PAGE_LIMIT {
			fmt.Printf("已达到最大页码限制（%d页），停止爬取\n", MAX_PAGE_LIMIT)
			break
		}

		fmt.Printf("\n===== 爬取主评论第%d页 =====\n", pn)
		mainResp, err := getMainComments(aid, pn)
		if err != nil {
			fmt.Println("主评论爬取失败：", err)
			pn++
			continue
		}

		// 遍历主评论
		for _, mainCom := range mainResp.Data.Replies {
			allComments = append(allComments, mainCom)
			fmt.Printf("主评论[%s] %s：%s\n", time.Unix(mainCom.Ctime, 0).Format("2006-01-02 15:04:05"), mainCom.Member.Uname, mainCom.Content.Message)
			// 爬取子评论
			//if mainCom.Reply > 0 {
			subPn := 1
			for {
				fmt.Printf("  ----- 爬取子评论第%d页 -----\n", subPn)
				subResp, err := getSubComments(aid, mainCom.Rpid, subPn)
				if err != nil {
					fmt.Println("子评论爬取失败：", err)
					break
				}

				// 遍历子评论
				for _, subCom := range subResp.Data.Replies {
					allComments = append(allComments, subCom)
					fmt.Printf("  子评论[%s] %s：%s\n", time.Unix(subCom.Ctime, 0).Format("2006-01-02 15:04:05"), subCom.Member.Uname, subCom.Content.Message)
				}
				if len(subResp.Data.Replies) == 0 {
					fmt.Println("  子评论无数据，停止爬取")
					break
				}
				subPn++
			}
			//}
		}

		// 主评论是否到底
		if mainResp.Data.Cursor.IsEnd {
			break
		}
		pn++
	}

	// 输出结果
	fmt.Printf("\n===== 爬取完成 =====\n")
	fmt.Printf("总计爬取评论数：%d条\n", len(allComments))
}
