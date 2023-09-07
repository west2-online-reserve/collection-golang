package main

import (
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var count int
var wg sync.WaitGroup
var mute = &sync.Mutex{}

func main() {
	times := 10000
	url := getTheHotReplyJsonUrl(2)
	req, err := http.NewRequest("GET", url, nil)
	reqq := *req
	if err != nil {
		log.Fatalln(err)
	}
	setRequestHeader(req)

	wg.Add(30)
	start := time.Now().UnixMicro()
	for i := 0; i < 30; i++ {
		go do(times, reqq, start)
	}
	wg.Wait()
	end := time.Now().UnixMicro()
	result := float64(end-start) / 1e6
	log.Printf("it cost %f s,that`s %f times a second", result, float64(count)/result)
}

func do(times int, req http.Request, start int64) {
	var client = &http.Client{}
	for i := 0; i < times; i++ {
		resp, _ := client.Do(&req)
		mute.Lock()
		count++
		mute.Unlock()
		if resp.StatusCode != 200 {
			end := time.Now().UnixMicro()
			result := float64(end-start) / 1000000
			log.Fatalf("Fatal! It cost %f s,that`s %f times a second", result, float64(count)/result)
		}
	}
	wg.Done()
}

func setRequestHeader(req *http.Request) {
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36")
	req.Header.Set("Origin", "https://www.bilibili.com")
	req.Header.Set("Referer", "https://www.bilibili.com/video/BV12341117rG/?p=2&spm_id_from=333.880.my_history.page.click&vd_source=6d8ea21e6f2f2c3344c170907eb4ca6c")
}

func getTheHotReplyJsonUrl(number int) string {
	return "https://api.bilibili.com/x/v2/reply/main?csrf=f5b57399c62e18808c56eebf2cdfec02&mode=3&oid=420981979&pagination_str=%7B%22offset%22:%22%7B%5C%22type%5C%22:1,%5C%22direction%5C%22:1,%5C%22data%5C%22:%7B%5C%22pn%5C%22:" + strconv.Itoa(number) + "%7D%7D%22%7D&plat=1&type=1"

}
