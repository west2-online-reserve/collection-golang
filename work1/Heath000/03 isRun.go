package main

import (
	"fmt"
)

func isRun(a int) bool {
	var flag bool = false
	if a%4 == 0 {
		flag = true
	}
	if a%100 == 0 {
		flag = false
		if a%400 == 0 {
			flag = true
		}
	}
	return flag
}
func main() {
	var (
		x     int
		y     int
		count int
	)
	fmt.Scanln(&x, &y)
	var a [10001]int
	count = 0
	for i := x; i <= y; i++ {
		if isRun(i) {
			a[count] = i
			count++
		}
	}
	fmt.Printf("%d\n", count)
	for i := 0; i < count; i++ {
		fmt.Printf("%d ", a[i])
	}

}
