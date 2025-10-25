package main

import "fmt"

//创建一个切片(slice) 使其元素为数字1-50，从切⽚删掉数字为3的倍数的数，
// 并且在末尾再增加⼀个数114514，输出切⽚。
//[1 2 4 5 7 8 10 11 13 14 16 17 19 20 22 23 25 26 28 29
//31 32 34 35 37 38 40 41 43 44 46 47 49 50 114514]

func main() {
	slice := make([]int, 50)
	for i := 1; i <= 50; i++ {
		slice[i-1] = i
	}

	count := 0
	for _, value := range slice {
		if value%3 != 0 {
			slice[count] = value
			count++
		}
	}
	slice = slice[:count]
	slice = append(slice, 114514)
	fmt.Println(slice)
}
