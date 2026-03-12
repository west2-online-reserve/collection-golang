package main

import "fmt"

// 第二题：淘淘采苹果
var (
	//全局变量单纯想用
	arr [10]int64
	h   int64
	num int64
)

func catch() {
	for i := 0; i < 10; i++ {
		fmt.Scan(&arr[i])
	}
	fmt.Scan(&h) //注意加取地址
	for i := 0; i < 10; i++ {
		if h+30 >= arr[i] {
			num++
		}
	}
	fmt.Println(num)
}

func main() {
	catch()
}
