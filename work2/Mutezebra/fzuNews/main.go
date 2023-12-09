package main

import (
	"fmt"
	"fzuNews/conf"
	"fzuNews/db/dao"
	"fzuNews/db/model"
	"github.com/bytedance/gopkg/lang/fastrand"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Data struct {
	URL  string
	Date string
	Id   int
}

var clients [200]*http.Client
var daos [200]*dao.InfoDao

func isTrue(stringDate string) bool {
	its := strings.Split(stringDate, "-")
	str := strings.Join(its, "")
	date, _ := strconv.Atoi(str)
	if date >= 20200102 && date <= 20210901 {
		return true
	}
	return false
}
func aClient() *http.Client {
	return clients[fastrand.Intn(200)]
}
func aDao() *dao.InfoDao {
	return daos[fastrand.Intn(200)]
}

func getTheUrlData() []*Data {
	startPage := 150
	endPage := 269
	resultData := make([]*Data, 0)

	for i := startPage; i <= endPage; i++ {
		url := "https://info22.fzu.edu.cn/lm_list.jsp?totalpage=952&PAGENUM=" + strconv.Itoa(i) + "&wbtreeid=1460"

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			panic(err)
		}
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36 Edg/113.0.1774.42")

		resp, err := aClient().Do(req)
		if err != nil {
			log.Println(err)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
		}
		err = resp.Body.Close()
		if err != nil {
			log.Println(err)
		}

		text := string(body)

		pattern := "(?s)</a>.*?<a.href=.(?P<url>.*?)wbnewsid=(?P<id>.*?)..target=_blank.title=.*?<span class=\"fr\">(?P<date>.*?)</span>"

		rePattern := regexp.MustCompile(pattern)
		results := rePattern.FindAllStringSubmatch(text, -1)
		for _, result := range results {
			if isTrue(result[rePattern.SubexpIndex("date")]) {
				aData := &Data{}
				aData.URL = "https://info22.fzu.edu.cn/" + result[rePattern.SubexpIndex("url")] + "wbnewsid=" + result[rePattern.SubexpIndex("id")]
				aData.Date = result[rePattern.SubexpIndex("date")]
				aData.Id, _ = strconv.Atoi(result[rePattern.SubexpIndex("id")])
				//fmt.Println(aData.Id)
				resultData = append(resultData, aData)
			}
		}
	}
	//fmt.Println(resultData[0].URL)
	return resultData
}
func allocation(data []*Data) {
	clientSum := len(data)/5 + 1
	wg.Add(clientSum)
	var i int
	for i = 0; i < clientSum-1; i++ {
		fmt.Printf("%d号协程启动\n", i+1)
		go spiderPage(i*5, data, 5)
	}
	go spiderPage(i*5, data, len(data)%5)
	wg.Wait()
}
func spiderPage(start int, data []*Data, end int) {
	for i := start; i < start+end; i++ {
		url := data[i].URL
		id := data[i].Id
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			continue
		}
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36 Edg/113.0.1774.42")

		resp, _ := aClient().Do(req)
		body, _ := io.ReadAll(resp.Body)
		text := string(body)
		resp.Body.Close()
		info := todo(text)
		if info != nil {
			url = "https://info22.fzu.edu.cn/system/resource/code/news/click/dynclicks.jsp?clickid=" + strconv.Itoa(id) + "&owner=1768654345&clicktype=wbnews"
			req, err = http.NewRequest("GET", url, nil)
			if err != nil {
				log.Println(err)
			}
			req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36 Edg/113.0.1774.42")

			resp, _ = aClient().Do(req)
			body, _ = io.ReadAll(resp.Body)
			text = string(body)
			info.Views = text
			write(info)
		}
	}
	wg.Done()
}
func todo(text string) *model.Info {
	patternTitle := "(?s)<title>(?P<title>(.*?))</title><ME"
	patternDate := "(?s)<div class=\"conthsj\" >日期： (?P<date>.*?)  &nbsp"
	patternAuthor := "(?s)信息来源..(?P<author>.*?)\n.."

	reTitle := regexp.MustCompile(patternTitle)
	reDate := regexp.MustCompile(patternDate)
	reAuthor := regexp.MustCompile(patternAuthor)

	result1 := reTitle.FindStringSubmatch(text)
	result2 := reDate.FindStringSubmatch(text)
	result3 := reAuthor.FindStringSubmatch(text)

	if len(result1) > 1 && len(result2) > 1 && len(result3) > 1 {
		return &model.Info{
			Title:  result1[reTitle.SubexpIndex("title")],
			Date:   result2[reDate.SubexpIndex("date")],
			Author: result3[reAuthor.SubexpIndex("author")],
		}
	}
	return nil
}
func write(info *model.Info) {
	adao := aDao()
	err := adao.Add(*info)
	if err != nil {
		log.Println(err)
	}
}

var wg sync.WaitGroup
var mute sync.Mutex

func main() {
	start := time.Now().Unix()
	initAll()
	resultData := getTheUrlData()
	allocation(resultData)
	end := time.Now().Unix()
	fmt.Println(end - start)
}

func initClient() {
	for i := 0; i < 200; i++ {
		clients[i] = &http.Client{}
	}
}
func initDao() {
	for i := 0; i < 200; i++ {
		daos[i] = &dao.InfoDao{DB: dao.Db()}
	}
}
func initLog() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetOutput(os.Stdout)
}
func initAll() {
	initClient()
	initLog()
	conf.InitConfig()
	dao.MysqlInit()
	initDao()
}
