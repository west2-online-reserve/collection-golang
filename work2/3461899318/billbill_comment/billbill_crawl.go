package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// API响应结构体
// API文档都有每个json字段的详细注释OWO
type BiliBiliResponse struct {
	Code    int    `json:"code"`    // 状态码，0表示成功，非0表示错误
	Message string `json:"message"` // 状态消息，成功时为"0"，错误时为错误描述
	TTL     int    `json:"ttl"`     // 生存时间，通常为1
	Data    struct {
		Cursor struct {
			IsBegin     bool   `json:"is_begin"`     // 是否为第一页
			Prev        int    `json:"prev"`         // 上一页页码
			Next        int    `json:"next"`         // 下一页页码
			IsEnd       bool   `json:"is_end"`       // 是否为最后一页
			Mode        int    `json:"mode"`         // 排序模式
			ModeText    string `json:"mode_text"`    // 排序模式文本描述
			AllCount    int    `json:"all_count"`    // 评论总数
			SupportMode []int  `json:"support_mode"` // 支持的排序模式列表
		} `json:"cursor"` // 分页信息
		Replies []Comment `json:"replies"` // 评论列表
	} `json:"data"` // 响应数据
}

// 评论结构体
type Comment struct {
	Rpid      int64 `json:"rpid"`      // 评论ID
	Oid       int64 `json:"oid"`       // 对象ID（视频ID）
	Type      int   `json:"type"`      // 评论区类型，1为视频
	Mid       int64 `json:"mid"`       // 发送者UID
	Root      int64 `json:"root"`      // 根评论ID（对话树最顶层）
	Parent    int64 `json:"parent"`    // 父评论ID
	Dialog    int64 `json:"dialog"`    // 对话ID
	Count     int   `json:"count"`     // 子评论数（当前返回数量）
	Rcount    int   `json:"rcount"`    // 子评论总数（实际总数）
	State     int   `json:"state"`     // 评论状态，0正常，其他为删除/屏蔽等
	Fansgrade int   `json:"fansgrade"` // 粉丝等级
	Attr      int   `json:"attr"`      // 评论属性位掩码
	Ctime     int64 `json:"ctime"`     // 评论创建时间戳（秒）
	Like      int   `json:"like"`      // 点赞数
	Action    int   `json:"action"`    // 当前用户操作状态
	Member    struct {
		Mid         string `json:"mid"`          // 用户UID（字符串形式）
		Uname       string `json:"uname"`        // 用户名
		Sex         string `json:"sex"`          // 性别，男/女/保密
		Sign        string `json:"sign"`         // 个性签名
		Avatar      string `json:"avatar"`       // 头像URL
		Rank        string `json:"rank"`         // 排名（未知用途）
		DisplayRank string `json:"display_rank"` // 显示排名
		LevelInfo   struct {
			CurrentLevel int `json:"current_level"` // 当前等级
		} `json:"level_info"` // 等级信息
	} `json:"member"` // 评论者信息
	Content struct {
		Message string `json:"message"` // 评论内容
		Plat    int    `json:"plat"`    // 发布平台
		Device  string `json:"device"`  // 设备信息
	} `json:"content"` // 评论内容信息
	Replies []Comment `json:"replies"` // 子评论列表（嵌套结构）
}

// 配置结构体
type Config struct {
	VideoID  string `json:"video_id"`
	PageSize int    `json:"page_size"`
	Delay    int    `json:"delay"`
}

// 获取视频响应的结构体
func fetchComments(videoID string, pageSize, next int) (*BiliBiliResponse, error) {
	url := fmt.Sprintf("https://api.bilibili.com/x/v2/reply?type=1&oid=%s&sort=2&ps=%d&pn=%d",
		videoID, pageSize, next)

	// 创建HTTP客户端
	client := &http.Client{
		//请求过程最大时间
		Timeout: 30 * time.Second,
	}

	// 创建请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Referer", "https://www.bilibili.com")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Cookie", "enable_web_push=DISABLE; DedeUserID=3546702799702502; DedeUserID__ckMd5=0f46fd111a3d496d; buvid_fp_plain=undefined; is-2022-channel=1; LIVE_BUVID=AUTO8017305556069754; hit-dyn-v2=1; enable_feed_channel=ENABLE; header_theme_version=OPEN; theme-tip-show=SHOWED; theme-avatar-tip-show=SHOWED; theme-switch-show=SHOWED; _uuid=314BAC9D-25D9-AE65-B4A9-C7231C6D64B437450infoc; theme_style=dark; buvid3=4972663A-5E84-E9EE-3D39-8193DDEF4B5644706infoc; b_nut=1756556044; rpdid=0zbfVG8DAW|UP7xSKPA|Mly|3w1USktD; buvid4=753677A4-54D9-7EE3-5E60-C0FC1378718B66895-024072605-oIEDEzHlzvzBSAXhnXyjpsry44fwXthi7mG+aLybzsnGp7y1FM+1Ij2DJ9IzTmc5; PVID=6; home_feed_column=5; browser_resolution=1592-820; fingerprint=53c14abdf5404077c0487d5910eddcac; buvid_fp=53c14abdf5404077c0487d5910eddcac; CURRENT_FNVAL=4048; b_lsid=66DF2F102_19A1E1FFC38; bsource=search_bing; bili_ticket=eyJhbGciOiJIUzI1NiIsImtpZCI6InMwMyIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjE3MDEyMTAsImlhdCI6MTc2MTQ0MTk1MCwicGx0IjotMX0.3ZnNbtA1dvzVJoBLhx0-CaeX7nqvv0iUTloqvT2ufVc; bili_ticket_expires=1761701150; bp_t_offset_3546702799702502=1127873929646440448; CURRENT_QUALITY=0; SESSDATA=9cf7b751%2C1776994020%2Ce9d59%2Aa1CjCR9Ztm30vuaA9p1sbwg6pTj6bUV7ca27692n-WhYp1afPUvXtGErtGR9F5sbYXxI0SVmlWcnd1M1dBeW8ydUFMMkZnVlRjVTRkSnhaeWxHeWJIOGxocjM0ZE4wTDF5MHdqa1VsX3BFN2lVeDhFQzBTV2lNRG15dFYtNXlWc3ZzWGp3VUNvM1FRIIEC; bili_jct=af2469c180c46c672bcd8a283772af98; sid=8cbicf3r")
	// 获得响应体
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应体信息
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	// 解析JSON进行反序列化
	var result BiliBiliResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, fmt.Errorf("解析JSON失败: %v", err)
	}

	// 检查API返回状态
	//0：成功
	//-400：请求错误
	//-404：无此项
	//12002：评论区已关闭
	//12009：评论主体的type不合法
	if result.Code != 0 {
		return nil, fmt.Errorf("API错误: %s (代码: %d)", result.Message, result.Code)
	}

	return &result, nil
}

// 获取所有二级评论（逻辑和获取一级评论差不多，但是需要一级评论的响应体的参数）
func getReplies(comment *Comment, videoID string, delay int) error {
	if comment.Rcount <= 0 {
		return nil
	}

	var allReplies []Comment
	page := 1
	pageSize := 20 // B站二级评论每页固定20条

	for {
		// 构建二级评论请求URL
		url := fmt.Sprintf("https://api.bilibili.com/x/v2/reply/reply?oid=%s&type=1&root=%d&ps=%d&pn=%d",
			videoID, comment.Rpid, pageSize, page)

		client := &http.Client{Timeout: 30 * time.Second}
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return fmt.Errorf("创建二级评论请求失败: %v", err)
		}

		// 设置请求头
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
		req.Header.Set("Referer", "https://www.bilibili.com")
		req.Header.Set("Accept", "application/json, text/plain, */*")
		//cookie根据自己的ip来写
		req.Header.Set("Cookie", "enable_web_push=DISABLE; DedeUserID=3546702799702502; DedeUserID__ckMd5=0f46fd111a3d496d; buvid_fp_plain=undefined; is-2022-channel=1; LIVE_BUVID=AUTO8017305556069754; hit-dyn-v2=1; enable_feed_channel=ENABLE; header_theme_version=OPEN; theme-tip-show=SHOWED; theme-avatar-tip-show=SHOWED; theme-switch-show=SHOWED; _uuid=314BAC9D-25D9-AE65-B4A9-C7231C6D64B437450infoc; theme_style=dark; buvid3=4972663A-5E84-E9EE-3D39-8193DDEF4B5644706infoc; b_nut=1756556044; rpdid=0zbfVG8DAW|UP7xSKPA|Mly|3w1USktD; buvid4=753677A4-54D9-7EE3-5E60-C0FC1378718B66895-024072605-oIEDEzHlzvzBSAXhnXyjpsry44fwXthi7mG+aLybzsnGp7y1FM+1Ij2DJ9IzTmc5; PVID=6; home_feed_column=5; browser_resolution=1592-820; fingerprint=53c14abdf5404077c0487d5910eddcac; buvid_fp=53c14abdf5404077c0487d5910eddcac; CURRENT_FNVAL=4048; b_lsid=66DF2F102_19A1E1FFC38; bsource=search_bing; bili_ticket=eyJhbGciOiJIUzI1NiIsImtpZCI6InMwMyIsInT5cCI6IkpXVCJ9.eyJleHAiOjE3NjE3MDEyMTAsImlhdCI6MTc2MTQ0MTk1MCwicGx0IjotMX0.3ZnNbtA1dvzVJoBLhx0-CaeX7nqvv0iUTloqvT2ufVc; bili_ticket_expires=1761701150; bp_t_offset_3546702799702502=1127873929646440448; CURRENT_QUALITY=0; SESSDATA=9cf7b751%2C1776994020%2Ce9d59%2Aa1CjCR9Ztm30vuaA9p1sbwg6pTj6bUV7ca27692n-WhYp1afPUvXtGErtGR9F5sbYXxI0SVmlWcnd1M1dBeW8ydUFMMkZnVlRjVTRkSnhaeWxHeWJIOGxocjM0ZE4wTDF5MHdqa1VsX3BFN2lVeDhFQzBTV2lNRG15dFYtNXlWc3ZzWGp3VUNvM1FRIIEC; bili_jct=af2469c180c46c672bcd8a283772af98; sid=8cbicf3r")

		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("请求二级评论失败: %v", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("读取二级评论响应失败: %v", err)
		}

		var replyResp BiliBiliResponse
		err = json.Unmarshal(body, &replyResp)
		if err != nil {
			return fmt.Errorf("解析二级评论JSON失败: %v", err)
		}

		if replyResp.Code != 0 {
			return fmt.Errorf("二级评论API错误: %s (代码: %d)", replyResp.Message, replyResp.Code)
		}

		// 添加当前页的二级评论
		allReplies = append(allReplies, replyResp.Data.Replies...)
		// 检查是否还有更多页
		if replyResp.Data.Cursor.IsEnd || len(replyResp.Data.Replies) == 0 {
			break
		}

		fmt.Printf("  已获取主评论rpid为 %d 的第 %d 页二级评论，当前共 %d 条\n", comment.Rpid, page, len(allReplies))

		page++

		// 延迟以避免请求过快
		time.Sleep(time.Duration(delay) * time.Millisecond)
	}

	comment.Replies = allReplies

	fmt.Printf("主评论 %d 的二级评论获取完成，共 %d 条\n", comment.Rpid, len(comment.Replies))
	return nil
}

// 打印评论统计信息
func printCommentStats(comments []Comment) {
	totalComments := len(comments)
	totalReplies := 0

	for _, comment := range comments {
		totalReplies += len(comment.Replies)
	}

	fmt.Printf("爬取完成！\n")
	fmt.Printf("主评论数量: %d\n", totalComments)
	fmt.Printf("子评论数量: %d\n", totalReplies)
	fmt.Printf("评论总数: %d\n", totalComments+totalReplies)
}

// 主函数
func main() {
	// 配置参数
	config := Config{
		VideoID:  "420981979", // 视频ID（oid），可以在视频URL中找到
		PageSize: 20,          // 每页评论数量
		Delay:    200,         // 请求延迟（毫秒）
	}
	fmt.Printf("开始爬取视频 的评论...\n")

	var allComments []Comment
	nextPage := 1
	totalFetched := 0

	for {
		fmt.Printf("正在获取第 %d 页评论...\n", nextPage)

		// 获取评论
		response, err := fetchComments(config.VideoID, config.PageSize, nextPage)
		if err != nil {
			log.Printf("获取评论失败: %v", err)
			break
		}

		// 如果没有评论，退出循环
		if len(response.Data.Replies) == 0 {
			fmt.Println("没有更多评论")
			break
		}

		// 处理当前页的评论
		for i, _ := range response.Data.Replies {
			comment := &response.Data.Replies[i]
			// 如果有子评论，获取所有子评论
			if comment.Rcount > 0 {
				fmt.Printf("获取主评论 %d 的二级评论 (共 %d 条)...\n", comment.Rpid, comment.Rcount)
				err := getReplies(comment, config.VideoID, config.Delay)
				if err != nil {
					log.Printf("获取二级评论失败: %v", err)
				}
			}
		}

		// 添加到总列表
		allComments = append(allComments, response.Data.Replies...)
		totalFetched += len(response.Data.Replies)

		fmt.Printf("已获取 %d 条主评论\n", totalFetched)

		// 检查是否还有更多页
		if response.Data.Cursor.IsEnd {
			fmt.Println("已到达最后一页")
			break
		}
		printCommentStats(response.Data.Replies)
		// 下一页
		nextPage++

		// 延迟以避免请求过快
		time.Sleep(time.Duration(config.Delay) * time.Millisecond)

	}

	// 保存评论到文件
	//filename := "bilibili_comments_" + config.VideoID + ".json"
	//err := saveCommentsToFile(allComments, filename)
	//if err != nil {
	//	log.Fatalf("保存文件失败: %v", err)
	//}
	// 打印统计信息

	//printCommentStats(allComments)
	// 打印前几条评论
	// 因为b站评论显示顺序的权重有自己的算法，这里是按照默认的时间顺序爬取，所以显示顺序与b站中可能出现出入
	//if len(allComments) > 0 {
	//	fmt.Println("\n=== 前5条评论示例 ===")
	//	for i := 0; i < len(allComments); i++ {
	//		comment := allComments[i]
	//		fmt.Printf("用户: %s\n", comment.Member.Uname)
	//		fmt.Printf("内容: %s\n", comment.Content.Message)
	//		fmt.Printf("时间: %s\n", time.Unix(comment.Ctime, 0).Format("2006-01-02 15:04:05"))
	//		//fmt.Printf("点赞: %d, 回复: %d\n", comment.Like, len(comment.Replies))
	//		fmt.Println("---")
	//	}
	//}
}
