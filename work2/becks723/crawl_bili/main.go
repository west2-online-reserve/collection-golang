package main

import (
	"fmt"
)

func main() {
	// 115344467365636
	// 爬取某个视频的评论（由oid指定哪个视频）
	comments := crawlComments("420981979")

	for _, c := range comments {
		fmt.Println(c)
	}

	// 写入数据库
	/*db, err := initDB()
	if err != nil {
		log.Fatal(err)
	}
	if err = insertComments(db, comments); err != nil {
		log.Fatal(err)
	}*/
	fmt.Println("ok!")
}
