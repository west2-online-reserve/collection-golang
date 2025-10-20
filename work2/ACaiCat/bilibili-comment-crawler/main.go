package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	url2 "net/url"
	"strconv"
	"time"

	"golang.org/x/time/rate"
)

const (
	sessionData = "..."
	oid         = "420981979"
	baseurl     = "https://api.bilibili.com/x/v2/reply"
	maxRate     = 0.01
	maxBurst    = 1
)

func main() {
	InitDB()
	ctx, cancel := context.WithCancel(context.Background())
	limiter := rate.NewLimiter(rate.Limit(maxRate), maxBurst)
	c := make(chan int)

	go func() {
		for page := 1; ; page++ {
			select {
			case <-ctx.Done():
				return
			default:
				c <- page
			}
		}
	}()

	for range maxBurst {
		go startCrawler(ctx, cancel, c, limiter)
	}

	<-ctx.Done()

}

func startCrawler(ctx context.Context, cancel context.CancelFunc, c <-chan int, limiter *rate.Limiter) {
	for {
		select {
		case <-ctx.Done():
			return
		case page := <-c:
			end := fetchComment(limiter, page)
			if end {
				cancel()
				return
			}
		}
	}
}

func fetchComment(limiter *rate.Limiter, page int) bool {
	for {
		_ = limiter.Wait(context.Background())
		client := &http.Client{Timeout: 5 * time.Second}
		parm := url2.Values{}
		parm.Set("oid", oid)
		parm.Set("type", "1")
		parm.Set("sort", "1")
		parm.Set("pn", strconv.Itoa(page))
		tagetURL, _ := url2.ParseRequestURI(baseurl)
		tagetURL.RawQuery = parm.Encode()

		req, _ := http.NewRequest("GET", tagetURL.String(), nil)
		req.AddCookie(&http.Cookie{
			Name:  "SESSDATA",
			Value: sessionData,
		})

		response, err := client.Do(req)

		fmt.Println("GET", page, response.Status)

		if err != nil {
			fmt.Println("Failed to request Bili API:", err.Error())
		}

		if response.StatusCode != http.StatusOK {
			continue
		}

		defer response.Body.Close()

		var result APIResponse

		err = json.NewDecoder(response.Body).Decode(&result)
		if err != nil {
			fmt.Println("Failed to decode JSON:", err.Error())
			return false
		}
		processComments(result.Data.Replies)

		return result.Code == http.StatusBadRequest
	}

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
