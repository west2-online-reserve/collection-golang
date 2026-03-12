package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"time"
)

var (
	bvid string
	aid  int64
	mix  string
)

func init() {
	bvid = "BV12341117rG"
	aid = getAidByBvid(bvid) // BV->aid
	log.Printf("aid = %d\n", aid)

	imgKey, subKey := getWbiKeys() // 取 WBI key

	log.Printf("imgKey = %s subKey = %s\n", imgKey, subKey)

	mix = mixinKey(imgKey, subKey) // 生成 mixinKey
	log.Printf("mix = %s\n", mix)
}
func getCommentURL(pg string) string {
	q := map[string]string{
		"oid":            fmt.Sprintf("%d", aid),
		"type":           "1",
		"mode":           "3",
		"pagination_str": pg,
		"plat":           "1",
	}

	wts, wRid := signWbi(q, mix)

	// 拼 URL
	u := "https://api.bilibili.com/x/v2/reply/wbi/main?" +
		"oid=" + q["oid"] +
		"&type=" + q["type"] +
		"&mode=" + q["mode"] +
		"&pagination_str=" + q["pagination_str"] +
		"&plat=1" +
		"&wts=" + wts +
		"&w_rid=" + wRid
	return u
}

// 生成子评论接口 URL：/x/v2/reply/wbi/reply
func getReplyURL(root int64, pn, ps int, mix string, aid int64) string {
	q := map[string]string{
		"oid":  fmt.Sprintf("%d", aid),
		"type": "1",
		"root": fmt.Sprintf("%d", root),
		"pn":   fmt.Sprintf("%d", pn),
		"ps":   fmt.Sprintf("%d", ps),
		"plat": "1",
	}
	wts, wRid := signWbi(q, mix) // 来自 webkey.go
	return "https://api.bilibili.com/x/v2/reply/reply?" +
		"oid=" + q["oid"] + "&type=1&root=" + q["root"] +
		"&pn=" + q["pn"] + "&ps=" + q["ps"] + "&plat=1" +
		"&wts=" + wts + "&w_rid=" + wRid
}

// CommentResponse 表示 Bilibili 评论 API 的完整响应
type CommentResponse struct {
	Code int `json:"code"`
	Data struct {
		Replies []Comment `json:"replies"`
		Cursor  struct {
			PaginationReply struct {
				NextOffset string `json:"next_offset"`
			} `json:"pagination_reply"`
		} `json:"cursor"`
	} `json:"data"`
}

// Comment 表示单条评论
type Comment struct {
	Rpid    int64   `json:"rpid"`
	Ctime   int64   `json:"ctime"`
	Like    int     `json:"like"`
	Member  Member  `json:"member"`
	Content Content `json:"content"`
}

// Member 表示评论者信息
type Member struct {
	Uname string `json:"uname"`
}

// Content 表示评论内容
type Content struct {
	Message string `json:"message"`
}

// 抓取指定主评(root=rpid)下所有子评
func fetchSubReplies(root int64, mix string, aid int64) ([]Comment, error) {
	var all []Comment
	for pn := 1; ; pn++ {
		u := getReplyURL(root, pn, 20, mix, aid) // 每页 20 条
		b := fetch(u)
		var r struct {
			Code int `json:"code"`
			Data struct {
				Replies []Comment `json:"replies"`
				Page    struct {
					Count int `json:"count"`
					Num   int `json:"num"`
					Size  int `json:"size"`
				} `json:"page"`
			} `json:"data"`
		}

		if err := json.Unmarshal(b, &r); err != nil {
			return nil, err
		}
		all = append(all, r.Data.Replies...)
		total := r.Data.Page.Count
		if len(r.Data.Replies) == 0 || pn*r.Data.Page.Size >= total {
			break
		}

		time.Sleep(200 * time.Millisecond)
	}
	return all, nil
}
func main() {
	offset := "" // 首页 offset 为空

	db := openDB("comments.db")

	for page := 1; ; page++ {
		//  query（不含 w_rid/wts）
		pg := url.QueryEscape(fmt.Sprintf(`{"offset":"%s"}`, offset))

		b := fetch(getCommentURL(pg))

		var r CommentResponse

		if err := json.Unmarshal(b, &r); err != nil {
			panic(err)
		}
		go func() {
			for _, c := range r.Data.Replies {
				SaveComment(db, c)
				subs, err := fetchSubReplies(c.Rpid, mix, aid)
				log.Printf("主评 %d 的子评数量 = %d", c.Rpid, len(subs))
				if err == nil {
					for _, sc := range subs {
						SaveComment(db, sc)
					}
				}
			}
		}()

		// 下一页 offset
		if r.Data.Cursor.PaginationReply.NextOffset == "" {
			break
		}
		offset = r.Data.Cursor.PaginationReply.NextOffset
		time.Sleep(300 * time.Millisecond) // 时间限制，防止封禁
	}
}
