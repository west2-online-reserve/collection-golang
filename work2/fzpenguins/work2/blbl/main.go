package main

import (
	"blbl/dataRow"
	"encoding/json"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"io"
	"net/http"
	"strconv"
	"sync"
)

var (
	db *gorm.DB
)

type Info struct {
	Mid     int64
	Message string
	Rpid    int64
	Count   int
}

type MyModel struct {
	Id   uint
	Text string
}

func dataBody(url string) (b *dataRow.Body, err error) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("err = ", err)
		return nil, err
	}
	defer resp.Body.Close()
	resp.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36 Edg/116.0.1938.69")
	resp.Header.Set("Origin", "https://www.bilibili.com")
	resp.Header.Set("Referer", "https://www.bilibili.com/video/BV12341117rG/?p=2&spm_id_from=333.880.my_history.page.click&vd_source=6d8ea21e6f2f2c3344c170907eb4ca6c")

	buf := make([]byte, 1024*4)
	buf, err = io.ReadAll(resp.Body) //
	if err != nil {
		fmt.Println("err = ", err)
		return
	}
	fmt.Println(string(buf))
	b = &dataRow.Body{}
	err = json.Unmarshal(buf, &b)
	if err != nil {
		fmt.Println("err = ", err)
		return nil, err
	}

	return b, nil
}

var wg sync.WaitGroup

func main() {
	ch := make(chan *Info)
	end := 1
	for i := 1; i <= end; i++ {
		wg.Add(1)
		go getMainReply(i, ch)
	}
	for i := 1; i <= end*2; i++ {
		go getSonReply(ch)
	}
	wg.Wait()

}

func getMainReply(index int, ch chan *Info) {
	url := "https://api.bilibili.com/x/v2/reply/main?csrf=f5b57399c62e18808c56eebf2cdfec02&mode=3&oid=420981979&pagination_str=%7B%22offset%22:%22%7B%5C%22type%5C%22:1,%5C%22direction%5C%22:1,%5C%22data%5C%22:%7B%5C%22pn%5C%22:" + strconv.Itoa(index) + "%7D%7D%22%7D&plat=1&type=1" //"https://api.bilibili.com/x/v2/reply/main?csrf=f5b57399c62e18808c56eebf2cdfec02&mode=3&oid=420981979&pagination_str=%7B%22offset%22:%22%7B%5C%22type%5C%22:1,%5C%22direction%5C%22:1,%5C%22data%5C%22:%7B%5C%22pn%5C%22:" + strconv.Itoa(index) + "%7D%7D%22%7D&plat=1&type=1"
	//url := "https://api.bilibili.com/x/v2/reply/main?csrf=f5b57399c62e18808c56eebf2cdfec02&mode=3&oid=420981979&pagination_str=%7B%22offset%22:%22%7B%5C%22type%5C%22:1,%5C%22direction%5C%22:1,%5C%22data%5C%22:%7B%5C%22pn%5C%22:1%7D%7D%22%7D&plat=1&type=1"
	buf, err := dataBody(url)
	if err != nil {
		fmt.Println("err = ", err)
		return
	}

	for _, v := range buf.Data.Replies {
		ch <- &Info{
			Mid:     v.Mid,
			Message: v.Content.Message,
			Rpid:    v.Rpid,
			Count:   v.Count,
		}
		saveDB(v.Content.Message)
	}
	wg.Done()
}

func getSonReply(ch chan *Info) {
	for {

		info := <-ch
		//infoSave
		var index int
		if info.Count%10 != 0 {
			index = info.Count/10 + 1
		} else {
			index = info.Count / 10
		}
		for i := 0; i <= index; i++ {
			url := "https://api.bilibili.com/x/v2/reply/reply?oid=420981979&type=1&root=" + strconv.FormatInt(info.Rpid, 10) + "&ps=10&pn=" + strconv.Itoa(i) + "&web_location=333.788"
			buf, err := dataBody(url)
			if err != nil {
				fmt.Println("err = ", err)
				return
			}
			//fmt.Println(buf.Data.Replies)

			for _, v := range buf.Data.Replies {
				//保存信息infoSave
				saveDB(v.Content.Message)
				fmt.Println(v.Content.Message)
			}
		}
	}
}

func saveDB(m string) {
	//user和pass，以及dbname根据实际情况填写，
	dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("err = ", err)
		return
	}
	db.AutoMigrate(&MyModel{})
	model := MyModel{Text: m}
	db.Create(&model)
}
