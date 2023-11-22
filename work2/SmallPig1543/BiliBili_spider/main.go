package main

import (
	"encoding/json"
	"github.com/gocolly/colly"
	"os"
	"strconv"
	"sync"
)

var dist *os.File

var wg sync.WaitGroup

// 爬取子评论的并发限制数
var channel = make(chan int, 5)

// 爬取主评论的并发限制数
var MainChannel = make(chan int, 30)

type JSONData struct {
	Data struct {
		Page struct {
			Size  int `json:"size"`  //每一页最多有多少评论或者回复
			Count int `json:"count"` //有多少评论或者回复
		} `json:"page"`
		Replies []struct {
			Rpid    int64 `json:"rpid"`
			Content struct {
				Message string `json:"message"`
			} `json:"content,omitempty"`
		} `json:"replies"`
	} `json:"data"`
}

func initRequest(c *colly.Collector) {
	c.OnRequest(func(req *colly.Request) {
		req.Headers.Set("Origin", "https://www.bilibili.com")
		req.Headers.Set("Referer", "https://www.bilibili.com/video/BV12341117rG/?vd_source=3b1d5937bf9e893fdc2d1a55a1b3dc4d")
		req.Headers.Set("User-Agent", "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Mobile Safari/537.36 Edg/117.0.2045.43")
		//req.Headers.Set("Cookie", "buvid3=10BD871C-59B2-6810-4D7C-F135A46E268B83471infoc; i-wanna-go-back=-1; _uuid=83CBA3610-8DE9-18DC-23EF-BC14BF3107EDD83000infoc; home_feed_column=5; buvid4=199E897E-DFDC-4594-578D-A2B50BC8BD9984375-023061719-X83v1qigvaWhlP6YG5eNsw%3D%3D; b_nut=1687000086; rpdid=|(u)luk)YYk)0J'uY)~|||uJY; b_ut=5; header_theme_version=CLOSE; buvid_fp_plain=undefined; hit-new-style-dyn=1; hit-dyn-v2=1; nostalgia_conf=-1; LIVE_BUVID=AUTO1716901134952324; DedeUserID=252709338; DedeUserID__ckMd5=f35fbf6efdab8d68; FEED_LIVE_VERSION=V8; CURRENT_BLACKGAP=0; CURRENT_QUALITY=80; CURRENT_FNVAL=4048; PVID=1; bili_ticket=eyJhbGciOiJIUzI1NiIsImtpZCI6InMwMyIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTYxNTQ4MjUsImlhdCI6MTY5NTg5NTU2NSwicGx0IjotMX0.EBY4TW7j6lMloRAz7sYwLjEqNfVHmaBbipDBgj9Y1Kk; bili_ticket_expires=1696154765; browser_resolution=1659-992; fingerprint=b7497332bda3d40c022ad14b9df899d6; bp_video_offset_252709338=846741442247262261; b_lsid=BCFBA176_18AE12669AE; buvid_fp=b7497332bda3d40c022ad14b9df899d6; SESSDATA=e9f91e33%2C1711546511%2C96ea5%2A91CjCAyvFV4aVr9VE8q-iKhCGUtjCSsQ89mdqvpCEhTD8g4WSdTbOC9ZIdf7Gw8fqxELISVkI5ODhFeGI1bHdvUl9XT3hGanJhRTJ1WUNBV1l5b0JETWJNZ2RMSE1LOEY4eFJYOW5BcWNPR2tLQnFDQndoWWJmNzB4cEsxT0xtVVZWNlBCamNtUUFBIIEC; bili_jct=ddaa1862b5ab6404ce5c4f38bd6b401c; sid=p2wvrg3n")
	})

}

// 获得评论或者回复的页数
func getPages(url string) int {
	c := colly.NewCollector()
	object := JSONData{}
	initRequest(c)
	c.OnResponse(func(res *colly.Response) {
		err := json.Unmarshal(res.Body, &object)
		if err != nil {
			panic(err)
		}
	})
	_ = c.Visit(url)
	if object.Data.Page.Count == 0 {
		return 0
	}
	pages := int(float64(object.Data.Page.Count/object.Data.Page.Size)) + 1
	return pages
}

// 查询主评论的信息
func queryMainText(url string) {
	c := colly.NewCollector()
	data := JSONData{}
	initRequest(c)
	c.OnResponse(func(res *colly.Response) {
		err := json.Unmarshal(res.Body, &data)
		if err != nil {
			panic(err)
		}
	})
	_ = c.Visit(url)
	for _, reply := range data.Data.Replies {
		//fmt.Println(reply.Content.Message)
		_, _ = dist.WriteString(reply.Content.Message + "\n")
		wg.Add(1)
		channel <- 1
		go querySubText(int(reply.Rpid))
	}
	<-MainChannel
	wg.Done()
}

// 查询回复信息
func querySubText(rpid int) {
	pages := getPages(getSubURL(rpid))
	for page := 1; page <= pages; page++ {
		url := "https://api.bilibili.com/x/v2/reply/reply?csrf=6f81254eb1291177725c64919659378d&oid=420981979&pn=" + strconv.Itoa(page) + "&ps=10&root=" + strconv.Itoa(rpid) + "&type=1"
		c := colly.NewCollector()
		data := JSONData{}
		initRequest(c)
		c.OnResponse(func(res *colly.Response) {
			err := json.Unmarshal(res.Body, &data)
			if err != nil {
				panic(err)
			}
		})
		_ = c.Visit(url)
		for _, reply := range data.Data.Replies {
			//fmt.Println(reply.Content.Message)
			_, _ = dist.WriteString(reply.Content.Message + "\n")
		}
	}
	<-channel
	wg.Done()
}

// 获得主评论的url
func getURL(pages int) string {
	return "https://api.bilibili.com/x/v2/reply?&jsonp=jsonp&pn=" + strconv.Itoa(pages) + "&type=1&oid=420981979&sort=2"
	//https://api.bilibili.com/x/v2/reply?&jsonp=jsonp&pn=[评论的页数]&type=1&oid=[AV号]&sort=2
}

// 获取子评论首页的url
func getSubURL(rpid int) string {
	return "https://api.bilibili.com/x/v2/reply/reply?csrf=6f81254eb1291177725c64919659378d&oid=420981979&pn=1&ps=10&root=" + strconv.Itoa(rpid) + "&type=1"
}

func main() {
	fileName := "./mytest.txt"
	dist, _ = os.Create(fileName)
	pages := getPages("https://api.bilibili.com/x/v2/reply?&jsonp=jsonp&pn=1&type=1&oid=420981979&sort=2")
	for page := 1; page <= pages; page++ {
		url := getURL(page)
		wg.Add(1)
		MainChannel <- 1
		go queryMainText(url)
	}
	wg.Wait()
}
