/*
爬取Bilibili视频评论
爬取 https://www.bilibili.com/video/BV12341117rG 的全部评论

全部评论，包含子评论
Bonus
给出Bilibili爬虫检测阈值（请求频率高于这个阈值将会被ban。也可以是你被封时的请求频率）
给出爬取的流程图，使用mermaid或者excalidraw
给出接口返回的json中每个参数所代表的意义
*/

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

// BV2AidResp: BV号转Aid接口返回结构
type BV2AidResp struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Aid int64 `json:"aid"` // 视频AID（评论接口需要）
	} `json:"data"`
}

// MainReplyResp: 旧版主评论接口返回结构
type MainReplyResp struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
	Ttl     int64  `json:"ttl"`
	Data    struct {
		Cursor struct {
			IsBegin         bool  `json:"is_begin"`
			Prev            int64 `json:"prev"`
			Next            int64 `json:"next"`
			IsEnd           bool  `json:"is_end"` // 是否是最后一页
			PaginationReply struct {
				NextOffset string `json:"next_offset"` // 下一页偏移量
			} `json:"pagination_reply"`
			SessionID string `json:"session_id"`
			Mode      int64  `json:"mode"`
			AllCount  int64  `json:"all_count"` // 评论总数
		} `json:"cursor"`
		Replies []struct { // 一级评论列表
			Rpid    int64 `json:"rpid"`  // 评论唯一ID
			Oid     int64 `json:"oid"`   // 视频AID
			Type    int64 `json:"type"`  // 评论类型（1=视频）
			Mid     int64 `json:"mid"`   // 评论用户ID
			Count   int64 `json:"count"` // 子评论数
			Content struct {
				Message string `json:"message"` // 评论内容
			} `json:"content"`
			Like int64 `json:"like"` // 点赞数
		} `json:"replies"`
	} `json:"data"`
}

// SubReplyResp: 旧版子评论接口返回结构
type SubReplyResp struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Replies []struct {
			Rpid    int64 `json:"rpid"`
			Content struct {
				Message string `json:"message"`
			} `json:"content"`
		} `json:"replies"`
		Page struct {
			Num   int64 `json:"num"`   // 当前页码
			Size  int64 `json:"size"`  // 每页条数
			Count int64 `json:"count"` // 总子评论数
		} `json:"page"`
	} `json:"data"`
}

// BV号转AID（评论接口需要AID作为oid参数）
func bv2Aid(bv string, cookie string) (int64, error) {
	apiUrl := fmt.Sprintf("https://api.bilibili.com/x/web-interface/view?bvid=%s", bv)
	req, _ := http.NewRequest("GET", apiUrl, nil)
	req.Header.Set("Cookie", cookie)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var bvResp BV2AidResp
	if err := json.Unmarshal(body, &bvResp); err != nil {
		return 0, err
	}
	if bvResp.Code != 0 {
		return 0, fmt.Errorf("BV转AID失败: %s (code=%d)", bvResp.Message, bvResp.Code)
	}
	return bvResp.Data.Aid, nil
}

// 爬取单页主评论
func crawlMainReply(aid int64, offset string, cookie string) (*MainReplyResp, error) {
	// 旧版主评论接口URL（核心：mode=3表示全部评论，mode=2是热门）
	params := map[string]string{
		"oid":            strconv.FormatInt(aid, 10),
		"type":           "1", // 固定值：视频评论
		"mode":           "3", // 3=全部评论
		"pagination_str": fmt.Sprintf(`{"offset":"%s"}`, offset),
		"plat":           "1",
		"seek_rpid":      "",
		"web_location":   "1315875",
	}

	// 构造URL
	apiUrl := "https://api.bilibili.com/x/v2/reply/main?"
	var queryStr string
	for k, v := range params {
		queryStr += k + "=" + url.QueryEscape(v) + "&"
	}
	queryStr = strings.TrimSuffix(queryStr, "&")
	apiUrl += queryStr

	// 发送请求
	req, _ := http.NewRequest("GET", apiUrl, nil)
	// 核心请求头（必须包含）
	req.Header.Set("Cookie", cookie)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Referer", fmt.Sprintf("https://www.bilibili.com/video/%s/", bv))
	req.Header.Set("authority", "api.bilibili.com")
	req.Header.Set("sec-ch-ua", `"Not_A Brand";v="8", "Chromium";v="120", "Microsoft Edge";v="120"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"Windows"`)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取完整返回（关键：打印完整JSON用于调试）
	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("主评论接口完整返回：\n%s\n", string(body))

	var mainResp MainReplyResp
	if err := json.Unmarshal(body, &mainResp); err != nil {
		return nil, fmt.Errorf("解析失败: %v", err)
	}
	if mainResp.Code != 0 {
		return nil, fmt.Errorf("接口返回错误: %s (code=%d)", mainResp.Message, mainResp.Code)
	}
	return &mainResp, nil
}

// 爬取某条主评论的所有子评论
func crawlSubReply(aid, rootRpid int64, cookie string) ([]string, error) {
	var allSubReplies []string
	pageNum := 1
	pageSize := 20

	for {
		params := map[string]string{
			"oid":          strconv.FormatInt(aid, 10),
			"type":         "1",
			"root":         strconv.FormatInt(rootRpid, 10),
			"ps":           strconv.FormatInt(int64(pageSize), 10),
			"pn":           strconv.FormatInt(int64(pageNum), 10),
			"web_location": "1315875",
		}

		apiUrl := "https://api.bilibili.com/x/v2/reply/reply?"
		var queryStr string
		for k, v := range params {
			queryStr += k + "=" + url.QueryEscape(v) + "&"
		}
		queryStr = strings.TrimSuffix(queryStr, "&")
		apiUrl += queryStr

		req, _ := http.NewRequest("GET", apiUrl, nil)
		req.Header.Set("Cookie", cookie)
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("子评论请求失败: %v", err)
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		var subResp SubReplyResp
		if err := json.Unmarshal(body, &subResp); err != nil {
			return nil, fmt.Errorf("子评论解析失败: %v", err)
		}
		if subResp.Code != 0 {
			return nil, fmt.Errorf("子评论接口错误: %s (code=%d)", subResp.Message, subResp.Code)
		}

		// 收集子评论
		for _, sub := range subResp.Data.Replies {
			allSubReplies = append(allSubReplies, sub.Content.Message)
		}

		// 分页判断
		totalPage := (subResp.Data.Page.Count + int64(pageSize) - 1) / int64(pageSize)
		if pageNum >= int(totalPage) {
			break
		}
		pageNum++
		time.Sleep(1 * time.Second)
	}
	return allSubReplies, nil
}

// -------------------------- 主函数 --------------------------
var (
	bv     = "BV12341117rG"
	cookie = `buvid4=BAB197E7-FAFE-8308-9981-2E906056B5AF35025-024090608-orrT4cKxJW9LzNOjOKIXgNs3LEXPZEPoTxXYKI9TD5j/5CBxdM16LQPbuIGNovuu; CURRENT_FNVAL=4048; CURRENT_QUALITY=0; _uuid=B88F1D22-F8F7-35610-10573-C10E21810BDBE791479infoc; DedeUserID=2001406901; DedeUserID__ckMd5=d169b0925b7be907; theme-tip-show=SHOWED; buvid3=141D7FBE-F01C-BAD9-CFB3-0BC4D3ED5E1843163infoc; b_nut=1761208643; buvid_fp=6bfb273c26ca5676173bc32e2fbd0e37; rpdid=|(kJYJlRk|u)0J'u~Yuk|)YRR; bp_t_offset_2001406901=1137696240679518208; theme-avatar-tip-show=SHOWED; b_lsid=5362FD58_19ACE95D27C; bili_ticket=eyJhbGciOiJIUzI1NiIsImtpZCI6InMwMyIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjQ2NjE3MTQsImlhdCI6MTc2NDQwMjQ1NCwicGx0IjotMX0.XLpCd9EZSfnA2A1cLEcAI1ZOZJtrgJAM4vvW_Nm3Fxk; bili_ticket_expires=1764661654; SESSDATA=d5950205%2C1779954532%2C5bc95%2Ab1CjBYKxkLE69X3LNs2ft6_p6HOhjJuxUxSZK6qYAl62cqdhPIOHT2O8PX5jQ4CE5cpMcSVnhIZGNOS2VDN0JUdFJsNFk1UjdiUUMwbDNkbW9uMFpEejNOb2VwS0hXTi1fYUtHUUFETGFrYWphN1pvQlVqR3l3dG84eFA1dXNtaXFtWnVhN2syUzNnIIEC; bili_jct=093222597bd07eb85f7e5ea47c2bb289; sid=mfmkls4b`
)

func main() {
	// 1. BV转AID
	aid, err := bv2Aid(bv, cookie)
	if err != nil {
		fmt.Printf("BV转AID失败: %v\n", err)
		return
	}
	fmt.Printf("成功 目标视频BV=%s → AID=%d\n", bv, aid)

	// 2. 初始化分页参数
	var allMainReplies []struct {
		Rpid    int64
		Content string
		Subs    []string
	}
	offset := "" // 初始偏移量
	pageNum := 1

	fmt.Println("开始爬取评论（旧版接口+修正结构）...")
	for {
		fmt.Printf("\n===== 爬取第%d页主评论 =====\n", pageNum)
		// 爬取单页主评论
		mainResp, err := crawlMainReply(aid, offset, cookie)
		if err != nil {
			fmt.Printf("错误 爬取主评论失败: %v\n", err)
			break
		}

		// 打印评论总数（确认有数据）
		fmt.Printf("统计 该视频总评论数: %d\n", mainResp.Data.Cursor.AllCount)
		fmt.Printf("当前页 一级评论数: %d\n", len(mainResp.Data.Replies))

		// 无更多评论则退出
		if mainResp.Data.Cursor.IsEnd || len(mainResp.Data.Replies) == 0 {
			if len(mainResp.Data.Replies) == 0 && pageNum == 1 {
				fmt.Println("警告 未获取到评论，可能原因：BV号无效/视频无评论/Cookie权限不足/mode参数错误")
			} else {
				fmt.Println("成功 已爬完所有主评论")
			}
			break
		}

		// 处理当前页主评论
		for _, mainReply := range mainResp.Data.Replies {
			fmt.Printf("\n一级评论 [%d]: %s\n", mainReply.Rpid, mainReply.Content.Message)

			// 爬取子评论
			subReplies, err := crawlSubReply(aid, mainReply.Rpid, cookie)
			if err != nil {
				fmt.Printf("警告 爬取子评论失败(rpid=%d): %v\n", mainReply.Rpid, err)
				subReplies = []string{}
			}

			// 收集结果
			allMainReplies = append(allMainReplies, struct {
				Rpid    int64
				Content string
				Subs    []string
			}{
				Rpid:    mainReply.Rpid,
				Content: mainReply.Content.Message,
				Subs:    subReplies,
			})

			// 打印子评论
			for i, sub := range subReplies {
				fmt.Printf("  二级评论 [%d]: %s\n", i+1, sub)
			}
		}

		// 更新偏移量（关键：从cursor.pagination_reply.next_offset获取）
		offset = mainResp.Data.Cursor.PaginationReply.NextOffset
		fmt.Printf("\n提示 下一页偏移量: %s\n", offset)

		// 频率控制
		time.Sleep(2 * time.Second)
		pageNum++
	}

	// 3. 保存结果
	outputPath := "bilibili_comments_old_api.txt"
	outputFile, err := os.Create(outputPath)
	if err != nil {
		fmt.Printf("错误 创建文件失败: %v\n", err)
		return
	}
	defer outputFile.Close()

	for _, item := range allMainReplies {
		outputFile.WriteString(fmt.Sprintf("一级评论 (rpid=%d): %s\n", item.Rpid, item.Content))
		if len(item.Subs) > 0 {
			outputFile.WriteString("  子评论列表:\n")
			for i, sub := range item.Subs {
				outputFile.WriteString(fmt.Sprintf("    %d. %s\n", i+1, sub))
			}
		}
		outputFile.WriteString("-------------------------\n")
	}

	// 4. 统计结果
	fmt.Printf("\n完成 爬取完成！\n")
	fmt.Printf("汇总 统计：共爬取 %d 条一级评论\n", len(allMainReplies))
	fmt.Printf("保存 结果已保存到：%s\n", outputPath)
}

