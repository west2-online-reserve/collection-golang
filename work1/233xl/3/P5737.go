package main

import "fmt"

func is4x(year int) bool {
	if (year % 400) == 0 {
		return true
	}
	if (year%100) == 0 && (year%400) != 0 {
		return false
	}
	if (year % 4) == 0 {
		return true
	}
	return false
}

func main() {
	var startYear, endYear int
	fmt.Scan(&startYear, &endYear)
	yearRange := endYear - startYear
	var res []int

	for i := 0; i <= yearRange; i++ {
		currentYear := startYear + i
		if is4x(currentYear) {
			res = append(res, currentYear)
		}
	}

	fmt.Println(len(res))
	fmt.Print(res)
}
