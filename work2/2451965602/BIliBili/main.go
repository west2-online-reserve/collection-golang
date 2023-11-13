package main

import (
	"strconv"
)

func main() {
	var start, end int
	start = 1
	end = pagerange("https://api.bilibili.com/x/v2/reply?type=1&oid=420981979&sort=1")
	for i := start; i <= end; i++ {
		url2 := "https://api.bilibili.com/x/v2/reply?type=1&oid=420981979&sort=1&pn=" + strconv.Itoa(i)
		date2, _ := httpget(url2)
		root := biliroot(date2)
		splider(root)
	}

}
