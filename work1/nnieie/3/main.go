package main

import "fmt"

func main() {
	var since, until int
	fmt.Scanf("%d %d", &since, &until)
	var leapYears []int
	for year := since; year <= until; year++ {
		if isLeap(year) {
			leapYears = append(leapYears, year)
		}
	}
	fmt.Printf("%d\n", len(leapYears))
	for _, y := range leapYears {
		fmt.Printf("%d ", y)
	}
}

func isLeap(y int) bool {
	return y%400 == 0 || (y%4 == 0 && y%100 != 0)
}