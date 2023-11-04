package main

import (
	helpers "bilibilicrawl/helper"
	"fmt"
)

func main() {
	db := helpers.InitDB()

	fmt.Println("Input the startPage and endPage")
	var endPage, startPage int
	_, err1 := fmt.Scanf("%d %d", &startPage, &endPage)
	if err1 != nil {
		panic(err1)
	}
	for i := startPage; i <= endPage; i++ {
		tempData := helpers.DataCrawl(i)
		for _, v := range tempData.Data.MainReply {
			helpers.DataMigration(v, db)
		}
	}

	err2 := helpers.CloseDB(db)
	if err2 != nil {
		panic(err2)
	}
	fmt.Println("Close database successfully")
}
