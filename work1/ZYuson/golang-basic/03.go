package main

import "fmt"

func main() {
	var x, y int
	var leapyear []int
	fmt.Scan(&x, &y)
	for year := x; year <= y; year++ {
		if isleap(year) {
			leapyear = append(leapyear, year)
		}
	}
	fmt.Println(len(leapyear))
	for i, year := range leapyear {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(year)
	}
}
func isleap(y int) bool {
	return (y%100 != 0 && y%4 == 0) || y%400 == 0
}
