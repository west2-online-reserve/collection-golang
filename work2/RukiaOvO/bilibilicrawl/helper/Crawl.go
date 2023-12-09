package helper

import (
	"bilibilicrawl/database/model"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

func DataCrawl(p int) model.BilibiliComment {
	url := "https://api.bilibili.com/x/v2/reply/main?csrf=f5b57399c62e18808c56eebf2cdfec02&mode=3&oid=420981979&pagination_str=%7B%22offset%22:%22%7B%5C%22type%5C%22:1,%5C%22direction%5C%22:1,%5C%22data%5C%22:%7B%5C%22pn%5C%22:" + strconv.Itoa(p) + "%7D%7D%22%7D&plat=1&type=1"
	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	temp := model.BilibiliComment{}
	err = json.Unmarshal(bodyText, &temp)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Crawl page%d successfully\n", p)
	return temp
}
