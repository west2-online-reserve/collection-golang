package main

import (
	"fmt"
	"log"
	"sort"
	"time"
)

/* 文章数据 */
type article struct {
	Date   string // 日期
	Author string // 作者
	Title  string // 标题
	Body   string // 正文
	Clicks int    // 点击数
}

func main() {
	// 初始化mysql数据库
	db, err := initDB()
	if err != nil {
		log.Fatal(err)
	}

	start := time.Now()

	// 并发爬取（普通爬取可以调用非Async版本：crawlFzu）
	articles := crawlFzuAsync(1091, 100)

	elapsed := time.Since(start)
	fmt.Printf("Crawl finished in %s.\n", elapsed)

	// 对爬取的文章数据进行排序整理
	// 按时间倒序
	sort.SliceStable(articles, func(i, j int) bool {
		return articles[i].Date > articles[j].Date
	})

	// 写入数据库
	err = insertArticles(db, articles)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("ok!")
}
