package main

import "fmt"

func main() {
	// 初始化苹果高度数组
	height := make([]int, 10)
	// 读取10个苹果的高度
	for i := 0; i < 10; i++ {
		fmt.Scan(&height[i])
	}

	// 读取陶陶能够达到的最大高度
	var H int
	fmt.Scan(&H)

	// 陶陶站在板凳上能够达到的最大高度
	H += 30

	// 初始化能够摘到的苹果数量
	var s int
	// 遍历苹果高度数组，计算能够摘到的苹果数量
	for _, h := range height {
		if H >= h {
			s++
		}
	}

	// 输出能够摘到的苹果数量
	fmt.Println(s)
}
