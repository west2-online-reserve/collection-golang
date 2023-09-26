package main

import (
	"fmt"
	"net/http"

	// "regexp"
	// "strconv"
	// "strings"
	// "time"
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	// "gorm.io/driver/mysql"
	// "gorm.io/gorm"
)

func main() {
	client := http.Client{}
	url := "https://www.bilibili.com/video/BV12341117rG/?vd_source=31f8c1ecdf6759e9f2c51d8466d5f534"
	cookie := &http.Cookie{Name: "bili_jct", Value: "272f9c1816f54e6ce527a7b11c5d6333", HttpOnly: false, Path: "/"}
	req, _ := http.NewRequest("GET", url, nil)
	req.AddCookie(cookie)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36")
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	html, _ := html.Parse(res.Body)
	listnode, _ := htmlquery.QueryAll(html, "//span[@class='reply-content']")
	for _, li := range listnode {
		text := htmlquery.InnerText(li)
		fmt.Println(text)
	}
}
