package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func main() {
	client := http.Client{}
	URL := "https://api.bilibili.com/x/v2/reply/wbi/main?oid=420981979&type=1&mode=3&pagination_str=%7B%22offset%22:%22%22%7D&plat=1&seek_rpid=&web_location=1315875&w_rid=3b7e1afeb49f4098bc477a53f1928526&wts=1763985592"
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		fmt.Println("requesterr", err)
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
	req.Header.Set("Origin", "https://www.bilibili.com")
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Referer", "https://www.bilibili.com/video/BV12341117rG/?vd_source=c1f176934af8d46f3c0b213af546afe5")
	req.Header.Set("User-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36 Edg/142.0.0.0")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("responseerr", err)
	}

	defer resp.Body.Close()

	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("ioerr", err)
	}

	var rslList BiliData
	_ = json.Unmarshal(bodyText, &rslList)
	for _, Preplies := range rslList.Data.Replies {
		fmt.Println("Primary reply:", Preplies.Content.Message)
		for _, Sreplies := range Preplies.Replies {
			fmt.Println("Secondary reply:", Sreplies.Content.Message)
		}
	}
}
