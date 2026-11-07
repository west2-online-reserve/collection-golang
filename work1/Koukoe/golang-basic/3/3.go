package main

import "fmt"

func main() {
	var start, end int
	var leapYears []int

	fmt.Scan(&start, &end)

	for i := start; i <= end; i++ {
		if (i%4 == 0 && i%100 != 0) || (i%400 == 0) {
			leapYears = append(leapYears, i)
		}
	}

	fmt.Println(len(leapYears))

	for i, leapYear := range leapYears {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(leapYear)
	}
	fmt.Println()
}
