package main

import (
	"fmt"
)

func isLeapYear(year int) bool {
	return (year%4 == 0 && year%100 != 0) || (year%400 == 0)
}

func main() {
	var x, y int
	fmt.Scanf("%d %d", &x, &y)

	if x > y {
		return
	}

	leapYears := []int{}
	leapYearCount := 0

	for year := x; year <= y; year++ {
		if isLeapYear(year) {
			leapYears = append(leapYears, year)
			leapYearCount++
		}
	}

	fmt.Println(leapYearCount)
	for i, year := range leapYears {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(year)
	}
	fmt.Println()
}