package main

//创建一个切片(slice) 使其元素为数字1-50，从切⽚删掉数字为3的倍数的数，并且在末尾再增加⼀个数114514，输出切⽚。 输出示例
//
//[1 2 4 5 7 8 10 11 13 14 16 17 19 20 22 23 25 26 28 29 31 32 34 35 37 38 40 41 43 44 46 47 49 50 114514]

import "fmt"

func main() {
	var arr []int
	for i := 0; i < 50; i++ {
		arr = append(arr, i+1)
	}
	for i := 0; i < len(arr); i++ {
		if arr[i]%3 == 0 {
			arr = append(arr[:i], arr[i+1:]...)
		}
	}
	arr = append(arr, 114514)
	fmt.Println(arr)
}
