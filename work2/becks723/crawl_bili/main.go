package main

import (
	"fmt"
	"log"
)

func main() {
	// 通过BV号获取目标视频oid（即av号）
	oid := bv2av("BV12341117rG")

	// 爬取视频评论（由oid指定哪个视频）
	comments := crawlComments(oid)

	for _, c := range comments {
		fmt.Println(c)
	}

	// 写入数据库
	db, err := initDB()
	if err != nil {
		log.Fatal(err)
	}
	if err = insertComments(db, comments); err != nil {
		log.Fatal(err)
	}
	fmt.Println("ok!")
}
