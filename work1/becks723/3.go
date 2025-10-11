package main

import "fmt"

func main() {
	var x, y int
	fmt.Scanf("%d %d", &x, &y)
	var ans [3000]int
	var cnt int
	for year := x; year <= y; year++ {
		if isLeapYear(year) {
			ans[cnt] = year
			cnt++
		}
	}

	fmt.Println(cnt)
	for i := 0; i < cnt; i++ {
		fmt.Printf("%d ", ans[i])
	}
}

func isLeapYear(year int) bool {
	return (year%4 == 0 && year%100 != 0) || year%400 == 0
}
