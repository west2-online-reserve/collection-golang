package main

import (
	"fmt"
)

func main() {
	var x, y, num int
	var flag [3000 + 1]bool
	fmt.Scanln(&x, &y)
	for i := x; i <= y; i++ {
		if i%4 == 0 && i%100 != 0 || i%400 == 0 {
			flag[i] = true
			num++
		}
	}
	fmt.Println(num)
	for i := x; i <= y; i++ {
		if flag[i] {
			fmt.Printf("%d ", i)
		}
	}
}
