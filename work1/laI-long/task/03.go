package main

import "fmt"

func main() {
	var x, y int
	fmt.Scan(&x, &y)
	var num int
	for i := x; i <= y; i++ {
		if (i%4 == 0 && i%100 != 0) || (i%400 == 0) {
			num++
		}
	}
	fmt.Println(num)
	for i := x; i <= y; i++ {
		if (i%4 == 0 && i%100 != 0) || (i%400 == 0) {
			fmt.Printf("%d ", i)
		}
	}
}
