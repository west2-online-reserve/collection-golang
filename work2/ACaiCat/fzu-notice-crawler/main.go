package main

import "C"
import (
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	startPage = 295
	endPage   = 372
	baseUrl   = "https://info22.fzu.edu.cn/"
)

func main() {
	var notices []*Notice
	client := &http.Client{
		Timeout: 120 * time.Second,
	}
	c := make(chan *Notice)

	var wg sync.WaitGroup
	for currentPage := startPage; currentPage <= endPage; currentPage++ {
		wg.Add(1)
		go fetchMenuPage(client, &wg, c, currentPage)
	}

	go func() {
		for notice := range c {
			notices = append(notices, notice)
		}
	}()

	wg.Wait()
	close(c)

	for _, notice := range notices {
		wg.Add(1)
		go fetchPage(client, &wg, notice)
	}

	wg.Wait()

	slices.SortStableFunc(notices, func(a, b *Notice) int {
		return a.ReleaseTime.Compare(b.ReleaseTime)
	})

	InitDB()
	db := GetDb()
	db.Create(notices)
	fmt.Printf("Crawler get %d records", len(notices))
}

func fetchPage(client *http.Client, wg *sync.WaitGroup, notice *Notice) {
	fmt.Println("Get " + notice.URL)
	response, err := client.Get(notice.URL)
	if err != nil {
		fmt.Println("Error on get page: " + err.Error())
		wg.Done()
		return
	}

	doc, err := goquery.NewDocumentFromReader(response.Body)
	notice.Body = doc.Find("#vsb_content").Text()
	wg.Done()
}

func fetchMenuPage(client *http.Client, wg *sync.WaitGroup, c chan<- *Notice, pageNum int) {
	url := baseUrl + "lm_list.jsp?wbtreeid=1460&PAGENUM=" + strconv.Itoa(pageNum)
	fmt.Println("Get " + url)
	response, err := client.Get(url)
	if err != nil {
		fmt.Println("Error on get page: " + err.Error())
		wg.Done()
		return
	}

	doc, err := goquery.NewDocumentFromReader(response.Body)
	// body > div.sy-content > div > div.right.fr > div.list.fl > ul > li:nth-child(1)
	doc.Find("body > div.sy-content > div > div.right.fr > div.list.fl > ul > li").Each(
		func(i int, selection *goquery.Selection) {
			// p > a.lm_a
			author := selection.Find("p > a.lm_a").Text()
			// p > a:nth-child(2)
			title := selection.Find("p > a:nth-child(2)").Text()
			// p > span
			date, _ := time.Parse("2006-01-02", selection.Find("p > span").Text())
			// p > a:nth-child(2)
			href, _ := selection.Find("p > a:nth-child(2)").Attr("href")
			url := baseUrl + href
			c <- &Notice{
				ReleaseTime: date,
				Author:      author,
				Title:       title,
				URL:         url,
			}
		})
	wg.Done()
}
