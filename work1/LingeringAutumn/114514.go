package main

import "fmt"

func main() {
	// 创建一个切片，使其元素为数字1-50
	slice := make([]int, 0, 51)
	for i := 1; i <= 50; i++ {
		slice = append(slice, i)
	}

	// 从切片中删掉数字为3的倍数的数
	for i := 0; i < len(slice); {
		if slice[i]%3 == 0 {
			slice = append(slice[:i], slice[i+1:]...)
		} else {
			i++
		}
	}

	// 在末尾增加一个数114514
	slice = append(slice, 114514)

	// 输出切片
	fmt.Println(slice)
}
