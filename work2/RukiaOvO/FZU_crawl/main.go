package main

import (
	"fmt"
	model "fzu_crawl/data"
	model2 "fzu_crawl/data/model"
	"github.com/PuerkitoBio/goquery"
	"gorm.io/gorm"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

const (
	headerSetting = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36"
	startPage     = 155
	endPage       = 278
)

func main() {
	db := model.InitDB()

	normalStart(db) //33 seconds
	//conStart(db) //1.8 seconds
	//加速比约为18

	err := model.CloseDB(db)
	if err != nil {
		panic(err)
	}
}

func normalStart(db *gorm.DB) {
	sTime := time.Now()
	emp := make(chan bool, 1)
	for i := startPage; i <= endPage; i++ {
		Crawl(db, strconv.Itoa(i), emp)
		<-emp
	}
	eTime := time.Now()

	fmt.Printf("耗时 %s\n", eTime.Sub(sTime))
}

func conStart(db *gorm.DB) {
	sTime := time.Now()
	ch := make(chan bool)
	for i := startPage; i <= endPage; i++ {
		go Crawl(db, strconv.Itoa(i), ch)
	}
	for i := startPage; i <= endPage; i++ {
		<-ch
	}
	eTime := time.Now()

	fmt.Printf("耗时 %s\n", eTime.Sub(sTime))
}

func Crawl(db *gorm.DB, p string, ch chan bool) {
	homeUrl := "https://info22.fzu.edu.cn/lm_list.jsp?totalpage=956&PAGENUM=" + p + "&wbtreeid=1460"
	homeClient := http.Client{}
	homeReq, err := http.NewRequest("GET", homeUrl, nil)
	if err != nil {
		panic(err)
	}
	homeReq.Header.Set("User-Agent", headerSetting)

	homeRes, err := homeClient.Do(homeReq)
	if err != nil {
		panic(err)
	}

	homeDetail, err := goquery.NewDocumentFromReader(homeRes.Body)
	if err != nil {
		panic(err)
	}

	for i := 1; i <= 20; i++ {
		tempSelect := fmt.Sprintf("body > div.sy-content > div > div.right.fr > div.list.fl > ul > li:nth-child(%d)", i)
		date := homeDetail.Find(tempSelect).Find("p > span").Text()
		if !dateCheck(date) {
			continue
		}

		author := homeDetail.Find(tempSelect).Find("p > a.lm_a").Text()
		author = author[3 : len(author)-3]
		title, _ := homeDetail.Find(tempSelect).Find("p > a:nth-child(2)").Attr("title")
		link, _ := homeDetail.Find(tempSelect).Find("p > a:nth-child(2)").Attr("href")
		text := "https://info22.fzu.edu.cn/" + link
		id := pickId(text)
		innerUrl := fmt.Sprintf("https://info22.fzu.edu.cn/system/resource/code/news/click/dynclicks.jsp?clickid=%d&owner=1768654345&clicktype=wbnews", id)

		innerClient := http.Client{}
		innerReq, err := http.NewRequest("GET", innerUrl, nil)
		if err != nil {
			panic(err)
		}
		innerReq.Header.Set("User-Agent", headerSetting)
		innerRes, err := innerClient.Do(innerReq)
		if err != nil {
			panic(err)
		}
		tempNums, err := io.ReadAll(innerRes.Body)
		if err != nil {
			panic(err)
		}
		nums := string(tempNums)

		newData := model2.News{}
		newData.Title = title
		newData.Author = author
		newData.Date = date
		newData.Text = text
		newData.Nums = nums

		db.Create(&newData)
		fmt.Println("Insert data, page =", p)

		innerRes.Body.Close()
	}

	homeRes.Body.Close()
	if ch != nil {
		ch <- true
	}
}

func dateCheck(d string) bool {
	numCheck := regexp.MustCompile("[0-9]+")
	numList := numCheck.FindAllString(d, -1)
	tempS := numList[0] + numList[1] + numList[2]
	tempI, err := strconv.Atoi(tempS)
	if err != nil {
		panic(err)
	}
	if 20200101 <= tempI && tempI <= 20210901 {
		return true
	}
	return false
}

func pickId(s string) int {
	idCheck := regexp.MustCompile("[0-9]+")
	sId := idCheck.FindAllString(s, -1)
	Id, err := strconv.Atoi(sId[2])
	if err != nil {
		panic(err)
	}
	return Id
}
