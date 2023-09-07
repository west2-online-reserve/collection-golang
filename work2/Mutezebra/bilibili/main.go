package main

import (
	"bilibili/conf"
	"bilibili/dataStruct"
	"bilibili/db/dao"
	"bilibili/db/model"
	"bilibili/types"
	"encoding/json"
	"github.com/bytedance/gopkg/lang/fastrand"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
)

var Clients [100]*http.Client
var wg sync.WaitGroup
var Daos [100]*dao.InfoDao

func main() {
	initAll()                    // 初始化
	ch := make(chan *types.Info) // 在协程之间传递消息
	end := 1
	for i := 1; i <= end; i++ {
		wg.Add(1)
		go getTheReplyInfo(i, ch)
	}
	for i := 1; i <= end*2; i++ {
		go getSubCommentInfo(ch)
	}
	wg.Wait()
}

// getTheReplyInfo 对热门评论进行操作
func getTheReplyInfo(index int, ch chan *types.Info) {
	log.Printf("第%d号协程已开启\n", index)
	url := getTheHotReplyJsonUrl(index)
	body, err := dataBody(url)
	if err != nil {
		log.Println(err)
	}
	for _, v := range body.Data.Replies {
		ch <- &types.Info{
			Mid:     v.Mid,
			Message: v.Content.Message,
			Rpid:    v.Rpid,
			Count:   v.Count,
		}
	}
	log.Println("Done!")
	wg.Done()
}

// getSubCommentInfo 对子评论进行操作
func getSubCommentInfo(ch chan *types.Info) {
	for {
		info := <-ch
		log.Println("Have received!")
		infoSave(info.Message, info.Mid)
		var pages int
		if info.Count%10 != 0 {
			pages = info.Count/10 + 1
		} else {
			pages = info.Count / 10
		}
		for page := 0; page <= pages; page++ {
			url := getSubCommentJsonUrl(page, info.Rpid)
			body, err := dataBody(url)
			if err != nil {
				log.Println(err)
			}
			for _, v := range body.Data.Replies {
				infoSave(v.Content.Message, v.Mid)
			}
		}
	}
}

// dataBody 以json形式返回url中的数据
func dataBody(url string) (*dataStruct.Body, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	setRequestHeader(req)
	resp, err := aClient().Do(req)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if err != nil {
		return nil, err
	}
	b := &dataStruct.Body{}
	err = json.Unmarshal(body, &b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// infoSave 将数据保存到数据库
func infoSave(message string, mid int64) {
	err := aDao().Add(&model.Info{Message: message, Mid: mid})
	if err != nil {
		log.Println(err)
	}
}

// 返回一个client对象
func aClient() *http.Client {
	return Clients[fastrand.Intn(100)]
}

// 返回一个db对象
func aDao() *dao.InfoDao {
	return Daos[fastrand.Intn(100)]
}

// getTheHotReplyJsonUrl 热评url
func getTheHotReplyJsonUrl(number int) string {
	return "https://api.bilibili.com/x/v2/reply/main?csrf=f5b57399c62e18808c56eebf2cdfec02&mode=3&oid=420981979&pagination_str=%7B%22offset%22:%22%7B%5C%22type%5C%22:1,%5C%22direction%5C%22:1,%5C%22data%5C%22:%7B%5C%22pn%5C%22:" + strconv.Itoa(number) + "%7D%7D%22%7D&plat=1&type=1"
}

// getSubCommentJsonUrl 子评论url
func getSubCommentJsonUrl(page int, rpid int64) string {
	return "https://api.bilibili.com/x/v2/reply/reply?oid=420981979&type=1&root=" + strconv.FormatInt(rpid, 10) + "&ps=10&pn=" + strconv.Itoa(page) + "&web_location=333.788"
}

// setRequestHeader 设置请求头
func setRequestHeader(req *http.Request) {
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36 Edg/116.0.1938.69")
	req.Header.Set("Origin", "https://www.bilibili.com")
	req.Header.Set("Referer", "https://www.bilibili.com/video/BV12341117rG/?p=2&spm_id_from=333.880.my_history.page.click&vd_source=6d8ea21e6f2f2c3344c170907eb4ca6c")
}

// initAll 初始化
func initAll() {
	initClient()
	conf.InitConfig()
	dao.MysqlInit()
	initDao()
}

// initDao 对全局变量进行初始化
func initDao() {
	for i := 0; i < 100; i++ {
		Daos[i] = &dao.InfoDao{DB: dao.Db()}
	}
}

// initClient 对全局变量进行初始化
func initClient() {
	for i := 0; i < len(Clients); i++ {
		Clients[i] = &http.Client{}
	}
}
