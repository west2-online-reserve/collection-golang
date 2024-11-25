package main

import (
	"fmt"
)

func main() {
	var a, b, s int

	fmt.Scan(&a, &b)

	for i := a; i <= b; i++ {
		if i%400 == 0 || (i%4 == 0 && i%100 != 0) {
			s++
		}
	}

	fmt.Println(s)

	for i := a; i <= b; i++ {
		if i%400 == 0 || (i%4 == 0 && i%100 != 0) {
			fmt.Print(i, " ")
		}
	}

	fmt.Println()
}
