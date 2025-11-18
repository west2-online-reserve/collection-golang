package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	"west2/crawler"
	"west2/db"
	"west2/model"
	"west2/util"

	"github.com/joho/godotenv"
	"golang.org/x/net/html"
)

func getFzuInfoByGoroutine() {
	defer fmt.Println("run by goroutine over")

	pool := crawler.NewPool(50, 500)
	defer pool.Stop()

	for i := 295; i <= 413; i++ {
		url := os.Getenv("BASE_URL") + "lm_list.jsp?totalpage=1092&PAGENUM=" + strconv.Itoa(i) + "&wbtreeid=1460"

		task := &crawler.Task{
			URL:     url,
			PageNum: i,
			Type:    crawler.TaskTypePage,
		}

		if !pool.Submit(task) {
			log.Printf("Failed to submit task for page %d", i)
		}

	}
	pool.Stop()
}

func getFzuInfo() {
	for i := 295; i <= 413; i++ {
		url := os.Getenv("BASE_URL") + "lm_list.jsp?totalpage=1092&PAGENUM=" + strconv.Itoa(i) + "&wbtreeid=1460"
		n := crawler.GetHtmlNode(url)

		niDb := db.InitDb()

		crawler.ParseNodeAndDeal(n, &model.Node{
			Type:      html.ElementNode,
			Data:      "li",
			ClassName: "clearfloat",
		}, func(n *html.Node) {
			res := &model.NotiInfo{}
			var href string
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if c.Type == html.ElementNode && c.Data == "p" {
					for nc := c.FirstChild; nc != nil; nc = nc.NextSibling {
						if nc.Type != html.ElementNode {
							continue
						}
						if nc.Data == "a" {
							if util.GetHtmlNodeValByKey(nc, "class") == "lm_a" {
								res.Author = strings.TrimSpace(nc.FirstChild.Data)
							} else {
								res.Title = util.GetHtmlNodeValByKey(nc, "title")
								href = util.GetHtmlNodeValByKey(nc, "href")
							}
						} else {
							res.Time = strings.TrimSpace(nc.FirstChild.Data)
						}
					}
				}
			}

			if res.Time < "2020-01-01" || res.Time > "2021-09-01" {
				return
			}

			crawler.ParseNodeAndDeal(crawler.GetHtmlNode(os.Getenv("BASE_URL")+href), &model.Node{
				Type:      html.ElementNode,
				Data:      "div",
				ClassName: "conthsj",
			}, func(n *html.Node) {
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					if c.Data != "script" {
						continue
					}
					if clicktype, owner, clickid, ok := util.ParseShowDynClicks(strings.TrimSpace(c.FirstChild.Data)); ok {
						cnt, err := crawler.GetFZUClickCount(clicktype, owner, clickid)
						if err != nil {
							log.Fatalf("获取点击数失败: %v", err)
							break
						}
						res.Count = cnt
					}
					break
				}
			})

			niDb.AddNi(res)
		})
	}
	fmt.Println("run over")
}

func getRunTime(f func()) time.Duration {
	start := time.Now()
	f()
	return time.Since(start)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	elapsed1 := getRunTime(getFzuInfoByGoroutine)
	fmt.Printf("程序执行时间: %v\n", elapsed1)
	elapsed2 := getRunTime(getFzuInfo)
	fmt.Printf("程序执行时间: %v\n", elapsed2)

	fmt.Printf("加速比：%.2f", float64(elapsed1)/float64(elapsed2))
}
