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

type Result struct {
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

type Parameters struct {
	Oid            string `json:"oid"`
	Type           string `json:"type"`
	Mode           string `json:"mode"`
	Pagination_str string `json:"pagination_str"`
	Plat           string `json:"plat"`
	Web_location   string `json:"web_location"`
	w_rid          string `json:"w_rid"`
	wts            string `json:"wts"`
}

func main() {
	params := Parameters{
		Oid:            "420981979",
		Type:           "1",
		Mode:           "3",
		Web_location:   "1315875",
		Plat:           "1",
		Pagination_str: `{"offset":""}`,
	}
	url := "https://api.bilibili.com/x/v2/reply/wbi/main?"
	GetWrid(&params)
	Spider(&url, &params)
	fmt.Println("分割")
	fmt.Println("分割")
	fmt.Println("分割")
	time.Sleep(1 * time.Second)
	params.Pagination_str = "%7B%22offset%22:%22CAESEDE4MDYyMTIzNjY4NDA2MDAiAggB%22%7D"
	fmt.Println("分割")
	fmt.Println("分割")
	fmt.Println("分割")
	url = "https://api.bilibili.com/x/v2/reply/wbi/main?"
	GetWrid(&params)
	Spider(&url, &params)
}

func Spider(url *string, params *Parameters) {
	Creaturl(url, params)
	client := &http.Client{}
	req, err := http.NewRequest("GET", *url, nil)
	if err != nil {
		fmt.Println("req err", err)
		return
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("resp err", err)
		return
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("io err", err)
		return
	}
	var result Result
	_ = json.Unmarshal(bodyText, &result)
	for _, v := range result.Data.Replies {
		fmt.Println("一级", v.Content.Message)
		for _, vv := range v.Replies {
			fmt.Println("二级", vv.Content.Message)
		}
	}
}

func Creaturl(url *string, params *Parameters) {
	*url = *url + "oid=" + params.Oid + "&type=" + params.Type + "&mode=" + params.Mode + "&pagination_str=" +
		params.Pagination_str + "&plat=" + params.Plat + "&seek_rpid=&web_location=" + params.Web_location +
		"&w_rid=" + params.w_rid + "&wts=" + params.wts
	fmt.Println(*url)
}
func Getwts(params *Parameters) {
	temp := time.Now().Unix()
	wtsTemp := strconv.Itoa(int(temp))
	params.wts = wtsTemp
}
func GetWrid(params *Parameters) {
	Getwts(params)
	queryString := "mode=" + params.Mode + "&oid=" + params.Oid + "&pagination_str=%7B%22offset%22%3A%22%22%7D&plat=" + params.Plat +
		"&seek_rpid=&type=" + params.Type + "&web_location=1315875&wts=" + params.wts
	salt := "ea1db124af3c7062474693fa704f4ff8"
	finalString := queryString + salt
	hash := md5.Sum([]byte(finalString))
	params.w_rid = fmt.Sprintf("%x", hash)
}
