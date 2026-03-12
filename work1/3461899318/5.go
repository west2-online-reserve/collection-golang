package main

import (
	"fmt"
)

func main() {
	slice1 := make([]int, 0, 50)
	for i := 1; i <= 50; i++ {
		slice1 = append(slice1, i)
	}
	//理论上第一个切片十分多余，但更加符合题目的模拟流程
	slice2 := make([]int, 0)
	for i := 1; i <= 50; i++ {
		if i%3 != 0 {
			slice2 = append(slice2, i)
		}
	}
	slice2 = append(slice2, 114514)
	fmt.Println(slice2)
}
