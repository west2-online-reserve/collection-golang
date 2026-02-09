package main

import "fmt"

func main() {
	var x, y int
	_, _ = fmt.Scan(&x, &y)

	if x > y {
		x, y = y, x
	}

	var leapYears []int

	for i := x; i <= y; i++ {
		if (i%4 == 0 && i%100 != 0) || (i%400 == 0) {
			leapYears = append(leapYears, i)
		}
	}

	fmt.Println(len(leapYears))
	for i, year := range leapYears {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(year)
	}
}
