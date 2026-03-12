package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
)

type bv2avResp struct {
	Data struct {
		Aid int64
	}
}

/* bv号转av号 */
func bv2av(bvid string) (avid string) {

	/*
	  b站提供的互转接口
	  bv转av: http://api.bilibili.com/x/web-interface/view?bvid=
	  av转bv: http://api.bilibili.com/x/web-interface/view?avid=
	*/

	req, err := http.NewRequest("GET", "http://api.bilibili.com/x/web-interface/view?bvid="+bvid, nil)
	if err != nil {
		log.Fatal(err)
	}

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

	var bv2avobj bv2avResp
	err = json.Unmarshal(body, &bv2avobj)
	if err != nil {
		log.Fatal(err)
	}

	avid = strconv.FormatInt(bv2avobj.Data.Aid, 10)
	return
}
