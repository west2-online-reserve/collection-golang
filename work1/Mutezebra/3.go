package main

import "fmt"

func main() {
	var begin, end, total int
	_, _ = fmt.Scanf("%d %d", &begin, &end)
	for i := begin; i <= end; i++ {
		if isOrNot(i) {
			total++
		}
	}
	fmt.Println(total)
	for i := begin; i <= end; i++ {
		if isOrNot(i) {
			fmt.Printf("%d ", i)
		}
	}
}

func isOrNot(year int) bool {
	if year%400 == 0 || (year%4 == 0 && year%100 != 0) {
		return true
	}
	return false
}
