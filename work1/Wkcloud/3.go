package main

import "fmt"

func main() {
	var x, y int
	years := make([]int, 0)
	fmt.Scan(&x, &y)
	for i := x; i <= y; i++ {
		if (i%4 == 0 && i%100 != 0) || (i%400 == 0) {
			years = append(years, i)
		}
	}
	fmt.Println(len(years))
	for _, year := range years {
		fmt.Printf("%d", year)
		fmt.Print(" ")
	}
}
