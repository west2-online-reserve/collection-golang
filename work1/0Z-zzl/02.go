package main

import "fmt"

func main() {
	var arr [10]int
	//输入苹果的高度
	for i := 0; i < 10; i++ {
		fmt.Scan(&arr[i])
		if arr[i] < 100 || arr[i] > 200 {
			fmt.Printf("第%d个数无法输入", i+1)
			arr[i] = 0
		}
	}
	//输入伸直达到的最大高度
	var a int
	fmt.Scan(&a)
	if a < 100 || a > 120 {
		fmt.Println("输入失败")
		a = 0
	}
	//判断是否能摘到苹果
	var n int = 0
	for i := 0; i < 10; i++ {
		if (a + 30) >= arr[i] {
			n = n + 1
		}
	}
	fmt.Println(n)
}
