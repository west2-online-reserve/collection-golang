package main

import "fmt"

func main() {
	var x, y int
	fmt.Scan(&x, &y)
	count := 0
	var leapYears []int
	for year := x; year <= y; year++ {
		if isLeapYear(year) {
			count++
			leapYears = append(leapYears, year)
		}
	}
	fmt.Println(count)
	for _, year := range leapYears {
		fmt.Print(year, " ")
	}
}

// 判断是否为闰年
func isLeapYear(year int) bool {
	return (year%4 == 0 && year%100 != 0) || year%400 == 0
}
