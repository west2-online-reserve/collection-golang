package main

import (
	"fmt"
	"strconv"
	"strings"
)

func main() {
	var minYear, maxYear int
	_, _ = fmt.Scan(&minYear, &maxYear)

	var leapYears []string
	for year := minYear; year <= maxYear; year++ {
		if (year%4 == 0 && year%100 != 0) || year%400 == 0 {
			leapYears = append(leapYears, strconv.Itoa(year))
		}
	}
	fmt.Println(len(leapYears))
	fmt.Println(strings.Join(leapYears, " "))

}
