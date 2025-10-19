package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	url2 "net/url"
	"time"

	"golang.org/x/time/rate"
)

const (
	sessionData = "258c5504%2C1776402420%2C5b6aa%2Aa1CjB0kJAW0APVakrP1hyrCZeLf-5IeASCuq04kwZEKhAGYBWyDY4urU8hVk0OBIcYMtMSVkRRYm1TaTZHb3Z0Q2cyemtmNWpTRDVhV0JMcDhxVjJaMDk4SDBQSTVRYUVMNXl0cHhqcjA3d3VCc3BLQWZRYVBpTEd1azRiSTQzWHRDSlpLZDgteV93IIEC"
	oid         = "420981979"
	baseurl     = "https://api.bilibili.com/x/v2/reply/wbi/main"
	maxRate     = 10
	maxBurst    = 5
)

var offset string

func main() {
	InitDB()
	offset, _ = fetchComment("")
	ctx, cancel := context.WithCancel(context.Background())
	limiter := rate.NewLimiter(rate.Limit(maxRate), maxBurst)
	for range maxBurst {
		startCrawler(ctx, cancel, limiter)
	}

	<-ctx.Done()

}

func startCrawler(ctx context.Context, cancel context.CancelFunc, limiter *rate.Limiter) {
	for {
		_ = limiter.Wait(context.Background())
		select {
		case <-ctx.Done():
			return
		default:
			nextOffset, end := fetchComment(offset)
			if end {
				cancel()
				return
			}
			offset = nextOffset
		}
	}
}

func fetchComment(offset string) (string, bool) {
	client := &http.Client{Timeout: 5 * time.Second}
	parm := url2.Values{}
	parm.Set("oid", oid)
	parm.Set("type", "1")
	parm.Set("mode", "3")
	parm.Set("plat", "1")
	parm.Set("web_location", "1315875")
	parm.Set("seek_rpid", "")
	parm.Set("pagination_str", fmt.Sprintf(`{"offset":"%s"}`, offset))

	tagetURL, _ := url2.ParseRequestURI(baseurl)
	tagetURL.RawQuery = parm.Encode()
	err := Sign(tagetURL)
	if err != nil {
		fmt.Println("Failed to sign url:", tagetURL.String())
	}

	req, _ := http.NewRequest("GET", tagetURL.String(), nil)
	req.AddCookie(&http.Cookie{
		Name:  "SESSDATA",
		Value: sessionData,
	})

	response, err := client.Do(req)

	fmt.Println("GET ", offset, response.Status)

	if err != nil {
		fmt.Println("Failed to request Bili API:", err.Error())
	}

	defer response.Body.Close()

	var result APIResponse

	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		fmt.Println("Failed to decode JSON:", err.Error())
		return "", false
	}
	processComments(result.Data.Replies)

	return result.Data.Cursor.PaginationReply.NextOffset, result.Data.Cursor.IsEnd
}

func processComments(rawComments []Reply) {
	var comments []Comment

	for _, c := range rawComments {
		var replies []Comment

		for _, r := range c.Replies {
			replies = append(replies, Comment{
				AuthorName: r.Member.Uname,
				AuthorUID:  r.Member.Mid,
				Time:       time.Unix(c.Ctime, 0),
				Like:       r.Like,
				Content:    r.Content.Message,
				IsReply:    true,
			})
		}

		comments = append(comments, Comment{
			AuthorName: c.Member.Uname,
			AuthorUID:  c.Member.Mid,
			Time:       time.Unix(c.Ctime, 0),
			Like:       c.Like,
			Content:    c.Content.Message,
			IsReply:    false,
			Relies:     replies,
		})

	}
	db := GetDB()
	db.Create(&comments)

}
