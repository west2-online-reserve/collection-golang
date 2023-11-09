package main

import "fmt"

func getLeapYear(startYear, endYear int) []int {
	var leapYear []int
	for year := startYear; year <= endYear; year++ {
		if (year%4 == 0 && year%100 != 0) || year%400 == 0 {
			leapYear = append(leapYear, year)
		}
	}
	return leapYear
}
func main() {
	var (
		startYear int
		endYear   int
	)
	fmt.Scan(&startYear)
	fmt.Scan(&endYear)
	getLeapYear(startYear, endYear)
	leapYears := getLeapYear(startYear, endYear)
	fmt.Println(len(leapYears))
	for _, leapYear := range leapYears {
		fmt.Printf("%d ", leapYear)
	}
}
