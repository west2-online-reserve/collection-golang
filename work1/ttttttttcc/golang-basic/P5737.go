package main

import (
	"fmt"
)


func isLeapYear(year int) bool {
	if year%4 == 0 && year%100 != 0 || year%400 == 0 {
		return true
	}
	return false
}

func main() {
	var x, y int

	
	fmt.Scan(&x, &y)

	
	if x > y {
		x, y = y, x
	}

	count := 0
	var years []int

	
	for year := x; year <= y; year++ {
		if isLeapYear(year) {
			count++
			years = append(years, year)
		}
	}

	
	fmt.Println(count)

	
	for _, year := range years {
		fmt.Print(year, " ")
	}
	fmt.Println() 
}