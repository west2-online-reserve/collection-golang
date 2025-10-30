package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	type_        string = "1"
	web_location string = "1315875"
	ps           int    = 10
)

/* 一页的评论数据 */
type replyData struct {
	Data struct {
		Cursor struct {
			Pagination_reply struct {
				Next_offset string
			}
		}
		Replies []reply
	}
}

/* 单个评论 */
type reply struct {
	Content struct {
		Message string
	}
	Member struct {
		Uname string
		Sex   string
	}
	Ctime    int64
	Like     int
	Mid      int64
	Rcount   int    // 子评论数量
	Rpid     int64  // 子评论的root，该字段用于请求子评论url
	Rpid_str string // rpid的字符串形式
}

/* 某个主评论下所有的子评论数据 */
type subReplyData struct {
	Data struct {
		Replies []reply
	}
}

/* 爬取指定视频下的所有评论 */
func crawlComments(oid string) (comments []mainComment) {
	p := 1
	offset := ""
	for {
		fmt.Printf("正在爬取第 %d 页的评论\n", p)
		offset, mc := crawl(oid, offset)
		for _, c := range mc {
			comments = append(comments, c)
		}
		p++

		if offset == "" {
			break
		}
	}
	return
}

/* 爬取一个评论页，并返回下一页的offset。若返回为空，则说明为最后一页 */
func crawl(oid string, offset string) (nextOffset string, comments []mainComment) {
	unix := time.Now().Unix()
	nowtime := strconv.FormatInt(unix, 10)
	w_rid, pagination_str := getSign(oid, offset, nowtime)

	baseUrl := "https://api.bilibili.com/x/v2/reply/wbi/main"

	params := url.Values{}
	params.Set("oid", oid)
	params.Set("type", type_)
	params.Set("mode", "3")
	params.Set("pagination_str", pagination_str)
	params.Set("plat", "1")
	if offset == "" {
		params.Set("seek_rpid", "")
	}
	params.Set("web_location", web_location)
	params.Set("w_rid", w_rid)
	params.Set("wts", nowtime)

	fullUrl := fmt.Sprintf("%s?%s", baseUrl, params.Encode())

	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/140.0.0.0 Safari/537.36")

	client := http.Client{}
	response, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", body)

	var rd replyData
	if err := json.Unmarshal(body, &rd); err != nil {
		log.Fatal(err)
	}

	for i, r := range rd.Data.Replies {
		// 获取子评论
		var subComments []comment
		pcount := ceil(float64(r.Rcount) / float64(ps)) // 总页数
		for pn := 1; pn <= pcount; pn++ {
			fmt.Printf("正在爬取第 %d 页的子评论……\n", pn)
			for _, c := range crawlSub(oid, r.Rpid_str, pn) {
				subComments = append(subComments, c)
			}
			time.Sleep(1 * time.Second)
		}

		fmt.Printf("%d: user: %s, comment: %s\n", i, r.Member.Uname, r.Content.Message)
		mc := mainComment{
			comment: comment{
				Rpid:    r.Rpid,
				Ctime:   r.Ctime,
				Like:    r.Like,
				Message: r.Content.Message,
				User: user{
					Uid:  r.Mid,
					Name: r.Member.Uname,
					Sex:  r.Member.Sex,
				},
			},
			SubComments: subComments,
		}
		comments = append(comments, mc)
	}

	nextOffset = rd.Data.Cursor.Pagination_reply.Next_offset
	return
}

/* 爬取一页的子评论 */
func crawlSub(oid string, root string, pn int) (subComments []comment) {
	baseUrl := "https://api.bilibili.com/x/v2/reply/reply"

	params := url.Values{}
	params.Set("oid", oid)
	params.Set("type", type_)
	params.Set("root", root)
	params.Set("ps", strconv.Itoa(ps))
	params.Set("pn", strconv.Itoa(pn))
	params.Set("web_location", "333.788")

	fullUrl := fmt.Sprintf("%s?%s", baseUrl, params.Encode())

	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var srd subReplyData
	if err := json.Unmarshal(body, &srd); err != nil {
		log.Fatal(err)
	}

	for i, r := range srd.Data.Replies {
		fmt.Printf("    %d: user: %s, comment: %s\n", i, r.Member.Uname, r.Content.Message)
		c := comment{
			Rpid:    r.Rpid,
			Ctime:   r.Ctime,
			Like:    r.Like,
			Message: r.Content.Message,
			User: user{
				Uid:  r.Mid,
				Name: r.Member.Uname,
				Sex:  r.Member.Sex,
			},
		}
		subComments = append(subComments, c)
	}
	return
}

func getSign(oid string, offset string, nowtime string) (w_rid string, pagination_str string) {
	pagination_str = fmt.Sprintf(`{"offset":"%s"}`, offset)

	params := url.Values{}
	params.Set("oid", oid)
	params.Set("type", type_)
	params.Set("mode", "3")
	params.Set("pagination_str", pagination_str)
	params.Set("plat", "1")
	if offset == "" {
		params.Set("seek_rpid", "")
	}
	params.Set("web_location", web_location)
	params.Set("wts", nowtime)
	v := params.Encode()
	a := "ea1db124af3c7062474693fa704f4ff8"
	w_rid = getMD5Hash(v + a)
	return
}

func getMD5Hash(s string) string {
	hash := md5.Sum([]byte(s))
	return hex.EncodeToString(hash[:])
}

func ceil(num float64) int {
	return int(math.Ceil(num))
}
