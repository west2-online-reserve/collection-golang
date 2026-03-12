package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

// 主评论和子评论函数都有cookie，如果爬不出来或者不完整可以更新自己对应的cookie,

var all_count int64

func main() {
	Spider()
}

type Resp struct {
	Code int64 `json:"code"`
	Data struct {
		Cursor struct {
			AllCount        int64 `json:"all_count"`
			PaginationReply struct {
				NextOffset string `json:"next_offset"`
			} `json:"pagination_reply"`
		}
		Replies []struct {
			Count   int   `json:"count"`
			Rpid    int64 `json:"rpid"`
			Content struct {
				Device  string        `json:"device"`
				JumpURL struct{}      `json:"jump_url"`
				MaxLine int64         `json:"max_line"`
				Members []interface{} `json:"members"`
				Message string        `json:"message"`
				Plat    int64         `json:"plat"`
			} `json:"content"`
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
				Replies interface{} `json:"replies"`
			} `json:"replies"`
			Type int64 `json:"type"`
		} `json:"replies"`
	} `json:"data"`
	Message string `json:"message"`
}

func getbranch(str string) {
	req, err := http.NewRequest("GET", str, nil)
	if err != nil {
		fmt.Println("err ", err)
	}

	req.Header.Set("cookie", "enable_web_push=DISABLE; buvid4=804115CA-FBBF-D554-CAD5-69E89F0A8B0681546-024062002-GH5AR28nChgirqANeReP5TlAN3V8N73k9kxQ236OQF3xcmQwoOfZhH2IYLw2aU0E; DedeUserID=106280196; DedeUserID__ckMd5=aef786c7f7ab1d24; buvid_fp_plain=undefined; LIVE_BUVID=AUTO5817301262924707; enable_feed_channel=ENABLE; hit-dyn-v2=1; fingerprint=6019c772ba55b49829bec8af7d21a415; buvid_fp=6019c772ba55b49829bec8af7d21a415; header_theme_version=OPEN; theme-tip-show=SHOWED; theme-avatar-tip-show=SHOWED; theme-switch-show=SHOWED; theme_style=light; buvid3=53250049-B76D-C90D-248B-046795A41C5B21562infoc; b_nut=1750389521; _uuid=6CA142E6-CD67-D8F2-767D-8AEAD102CB7F723649infoc; rpdid=0zbfVHc1lp|1dYiSBAGq|4q|3w1UGjxp; PVID=2; CURRENT_QUALITY=80; home_feed_column=5; browser_resolution=2000-1026; bili_ticket=eyJhbGciOiJIUzI1NiIsImtpZCI6InMwMyIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjEzMTM4OTMsImlhdCI6MTc2MTA1NDYzMywicGx0IjotMX0._B1zfJaIAXicKX8hVlHUFGwA-5PxHCWmdVPcpsrgca8; bili_ticket_expires=1761313833; SESSDATA=c6be17ae%2C1776606693%2C17b6b%2Aa2CjD-WaIRFMHi8SMN-UQ1xSiI8uL0xdswP2obvivJjdHvdsMMz_4cKAVdYFrs57KsYqUSVjQyZGtPYnUzUmNoWGVtQnBmZ1Znb1pnbzB5RFp2T3kzTTh3UjZ1VlZEWFBRa3pUQU5VREZIVElVWFIxeXRndWNXczB4ME9CY0RtY3NjOV9IQXVDcUN3IIEC; bili_jct=489cfefaff4b7eaf363bfd706f7afc23; sid=5smd70hl; CURRENT_FNVAL=4048; b_lsid=2CEABC10B_19A0BFE75F4; bp_t_offset_106280196=1126574891017961472")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36 Edg/141.0.0.0")
	req.Header.Set("referer", "https://www.bilibili.com/video/BV1ACxKzBEJp/?spm_id_from=333.1007.tianma.1-1-1.click&vd_source=3de3373552f474732a2060121a97b2b4")
	req.Header.Set("sec-fetch-site", "same-site")
	req.Header.Set("sec-fetch-mode", "cors")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var resultList Resp
	_ = json.Unmarshal(bodyText, &resultList)
	for _, result := range resultList.Data.Replies {
		fmt.Println("二级评论：", result.Content.Message)
	}
}

func mainmessage(str string) string {
	uurl := str
	req, err := http.NewRequest("GET", uurl, nil)
	if err != nil {
		fmt.Println("err ", err)
	}
	req.Header.Set("cookie", "enable_web_push=DISABLE; buvid4=804115CA-FBBF-D554-CAD5-69E89F0A8B0681546-024062002-GH5AR28nChgirqANeReP5TlAN3V8N73k9kxQ236OQF3xcmQwoOfZhH2IYLw2aU0E; DedeUserID=106280196; DedeUserID__ckMd5=aef786c7f7ab1d24; buvid_fp_plain=undefined; LIVE_BUVID=AUTO5817301262924707; enable_feed_channel=ENABLE; hit-dyn-v2=1; fingerprint=6019c772ba55b49829bec8af7d21a415; buvid_fp=6019c772ba55b49829bec8af7d21a415; header_theme_version=OPEN; theme-tip-show=SHOWED; theme-avatar-tip-show=SHOWED; theme-switch-show=SHOWED; theme_style=light; buvid3=53250049-B76D-C90D-248B-046795A41C5B21562infoc; b_nut=1750389521; _uuid=6CA142E6-CD67-D8F2-767D-8AEAD102CB7F723649infoc; rpdid=0zbfVHc1lp|1dYiSBAGq|4q|3w1UGjxp; PVID=2; CURRENT_QUALITY=80; home_feed_column=5; browser_resolution=2000-1026; bili_ticket=eyJhbGciOiJIUzI1NiIsImtpZCI6InMwMyIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjEzMTM4OTMsImlhdCI6MTc2MTA1NDYzMywicGx0IjotMX0._B1zfJaIAXicKX8hVlHUFGwA-5PxHCWmdVPcpsrgca8; bili_ticket_expires=1761313833; SESSDATA=c6be17ae%2C1776606693%2C17b6b%2Aa2CjD-WaIRFMHi8SMN-UQ1xSiI8uL0xdswP2obvivJjdHvdsMMz_4cKAVdYFrs57KsYqUSVjQyZGtPYnUzUmNoWGVtQnBmZ1Znb1pnbzB5RFp2T3kzTTh3UjZ1VlZEWFBRa3pUQU5VREZIVElVWFIxeXRndWNXczB4ME9CY0RtY3NjOV9IQXVDcUN3IIEC; bili_jct=489cfefaff4b7eaf363bfd706f7afc23; sid=5smd70hl; CURRENT_FNVAL=4048; b_lsid=2CEABC10B_19A0BFE75F4; bp_t_offset_106280196=1126574891017961472")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36 Edg/141.0.0.0")
	req.Header.Set("referer", "https://www.bilibili.com/video/BV12341117rG/?vd_source=3de3373552f474732a2060121a97b2b4")
	req.Header.Set("sec-fetch-site", "same-site")
	req.Header.Set("sec-fetch-mode", "cors")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var resultList Resp
	_ = json.Unmarshal(bodyText, &resultList)
	for _, result := range resultList.Data.Replies {
		rpid := result.Rpid
		count := result.Count
		fmt.Println("评论：", result.Content.Message)
		for i := 1; i*10 <= count+9; i++ {
			branch := "https://api.bilibili.com/x/v2/reply/reply?oid=420981979&type=1&root=" + strconv.FormatInt(rpid, 10) + "&ps=10&pn=" + strconv.Itoa(i) + "&web_location=333.788"
			getbranch(branch)
			time.Sleep(1 * time.Second)
		}
		if count != 0 {
			fmt.Printf("共%d条二级评论\n\n", count)
		}
	}
	var nextpagn string
	all_count = resultList.Data.Cursor.AllCount
	if resultList.Data.Cursor.PaginationReply.NextOffset != "" {
		nextpagn = resultList.Data.Cursor.PaginationReply.NextOffset
	} else {
		nextpagn = ""
	}
	return nextpagn
}

func Spider() {
	v := "mode=3&oid=420981979&pagination_str=%7B%22offset%22%3A%22%22%7D&plat=1&seek_rpid=&type=1&web_location=1315875&wts=" + strconv.FormatInt(time.Now().Unix(), 10)
	a := "ea1db124af3c7062474693fa704f4ff8"
	data := []byte(v + a)
	has := md5.Sum(data)
	w_rid := fmt.Sprintf("%x", has)
	uurl := "https://api.bilibili.com/x/v2/reply/wbi/main?oid=420981979&type=1&mode=3&pagination_str=%7B%22offset%22:%22%22%7D&plat=1&seek_rpid=&web_location=1315875&w_rid=" + w_rid + "&wts=" + strconv.FormatInt(time.Now().Unix(), 10)
	//uurl := "https://api.bilibili.com/x/v2/reply/wbi/main?oid=420981979&type=1&mode=3&pagination_str=%7B%22offset%22:%22%22%7D&plat=1&seek_rpid=&web_location=1315875&w_rid=0a324c573d83f5241d17d3e6ba093153&wts=1761139606"
	nextpagn := mainmessage(uurl)
	for nextpagn != "" {
		time.Sleep(1 * time.Second)
		v := "mode=3&oid=420981979&pagination_str=%7B%22offset%22%3A%22" + nextpagn + "%22%7D&plat=1&type=1&web_location=1315875&" + "wts=" + strconv.FormatInt(time.Now().Unix(), 10)
		a := "ea1db124af3c7062474693fa704f4ff8"
		data := []byte(v + a)
		has := md5.Sum(data)
		w_rid := fmt.Sprintf("%x", has)
		uurl = "https://api.bilibili.com/x/v2/reply/wbi/main?oid=420981979&type=1&mode=3&pagination_str=%7B%22offset%22:%22" + nextpagn + "%22%7D&plat=1&web_location=1315875&w_rid=" + w_rid + "&wts=" + strconv.FormatInt(time.Now().Unix(), 10)
		nextpagn = mainmessage(uurl)
	}
	fmt.Println("共", all_count, "条评论")
}
