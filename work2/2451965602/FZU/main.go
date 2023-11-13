package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"time"
)

var a = make(chan int)

func splider(i int) {
	url := "https://info22.fzu.edu.cn/lm_list.jsp?totalpage=961&PAGENUM=" + strconv.Itoa(i) + "&wbtreeid=1460"
	FZU(url)
	a <- i
}

func working(start, end int) {
	for i := start; i <= end; i++ {
		go splider(i)
	}

	for i := start; i <= end; i++ {
		fmt.Println(<-a)
	}

}

func main() {
	start1 := time.Now()
	var start, end int
	start = 159
	end = 279
	working(start, end)
	elapsed := time.Since(start1)
	fmt.Println("该函数执行完成耗时：", elapsed)
}
