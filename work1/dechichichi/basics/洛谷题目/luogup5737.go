package main

import (
	"fmt"
)

func findrun(A int, B int) (int, []int) {
	var tag int
	list := make([]int, 0)
	for i := A; i <= B; i++ {
		if (i%4 == 0 && i%100 != 0) || i%400 == 0 {
			tag++
			list = append(list, i)
		}
	}
	return tag, list
}

func main() {
	var A, B int
	fmt.Scan(&A, &B)
	tag, list := findrun(A, B)
	fmt.Println(tag)
	for _, val := range list {
		fmt.Print(val, " ")
	}
}
