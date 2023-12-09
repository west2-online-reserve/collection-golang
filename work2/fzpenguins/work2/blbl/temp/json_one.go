package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type Info struct {
	Mid     int64
	Message string
	Rpid    int64
	Count   int
}

type Body struct {
	Data Data        `json:"data"`
	Info interface{} `json:"info"`
}

type Data struct {
	Replies []Replies   `json:"replies"`
	Info    interface{} `json:"info"`
}

type Replies struct {
	Content Content     `json:"content"`
	Count   int         `json:"count"`
	Rpid    int64       `json:"rpid"`
	Mid     int64       `json:"mid"`
	Info    interface{} `json:"info"`
}

type Content struct {
	Message string      `json:"message"`
	Info    interface{} `json:"info"`
}

var wg sync.WaitGroup

func main() {

	url := "https://api.bilibili.com/x/v2/reply/wbi/main?oid=420981979&type=1&mode=3&pagination_str=%7B%22offset%22:%22%22%7D&plat=1&seek_rpid=&web_location=1315875&w_rid=3fd03ab98adc51a878f6e75aeee3d824&wts=1696045913"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("err = ", err)
		return
	}
	buf := make([]byte, 1024*4)
	buf, err = io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("err=", err)
		return
	}
	b := &Body{}
	err = json.Unmarshal(buf, &b)
	if err != nil {
		fmt.Println("err. =", err)
		return
	}
	ch := make(chan *Info)
	for _, v := range b.Data.Replies {
		fmt.Println(v.Content)
	}
	for i := 1; i <= 1; i++ {
		wg.Add(1)
		go getTheReplyInfo(i, ch)
	}
	for i := 1; i <= 2; i++ {
		go getSubCommentInfo(ch)
	}
	wg.Wait()
}

func getSubCommentInfo(ch chan *Info) {
	for {
		info := <-ch

		var pages int

		if info.Count%10 != 0 {
			pages = info.Count/10 + 1
		} else {
			pages = info.Count / 10
		}
		for page := 0; page <= pages; page++ {
			url := "https://api.bilibili.com/x/v2/reply/reply?oid=420981979&type=1&root=" + strconv.FormatInt(info.Rpid, 10) + "&ps=10&pn=" + strconv.Itoa(page) + "&web_location=333.788"
			buf, err := httpGet(url)

			if err != nil {
				log.Println(err)
			}
			b := &Body{}
			err = json.Unmarshal(buf, &b)
			if err != nil {
				fmt.Println("err. =", err)
				return
			}

			fmt.Println(b.Data.Replies)
		}
	}
}

func httpGet(url string) (buf []byte, err error) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("err =", err)
		return
	}
	defer resp.Body.Close()

	buf = make([]byte, 1024*4)
	buf, err = io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("err=", err)
		return
	}

	return buf, err
}
func getTheReplyInfo(index int, ch chan *Info) {

	url := "https://api.bilibili.com/x/v2/reply/wbi/main?oid=420981979&type=1&mode=3&pagination_str=%7B%22offset%22:%22%22%7D&plat=1&seek_rpid=&web_location=1315875&w_rid=3fd03ab98adc51a878f6e75aeee3d824&wts=1696045913"

	buf, err := httpGet(url)
	if err != nil {
		log.Println(err)
	}
	b := &Body{}
	err = json.Unmarshal(buf, &b)
	if err != nil {
		fmt.Println("err. =", err)
		return
	}
	for _, v := range b.Data.Replies {
		ch <- &Info{
			Mid:     v.Mid,
			Message: v.Content.Message,
			Rpid:    v.Rpid,
			Count:   v.Count,
		}
	}
	wg.Done()
}
