package main

import (
	"fmt"
	"time"
)

func main() {
	intChan := make(chan int, 1)
	go func() {
		time.Sleep(time.Second)
		fmt.Println("存int数据")
		intChan <- 10
	}()
	stringChan := make(chan string, 2)
	go func() {
		time.Sleep(time.Second)
		fmt.Println("存string数据")
		stringChan <- "asdfghjkl"
	}()

	// fmt.Println(<-intChan)
	select {
	case v := <-intChan:
		fmt.Println(v)
	case v := <-stringChan:
		fmt.Println(v)
	default:
		fmt.Println("管道为空读取失败")
	}
}
