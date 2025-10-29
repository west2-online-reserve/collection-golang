package main

import (
	"fmt"
)

func main() {
	var x, y int
	fmt.Scan(&x, &y)
	arr := make([]int, 0)
	for x <= y {
		if (x%4 == 0 && x%100 != 0) || x%400 == 0 {
			arr = append(arr, x)
		}
		x++
	}
	fmt.Println(len(arr))
	for _, x := range arr {
		fmt.Print(x, " ")
	}
}
