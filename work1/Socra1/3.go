package main

import (
	"fmt"
)

func isrun(x int) bool {
	if x%400 == 0 {
		return true
	} else if x%4 == 0 && x%100 != 0 {
		return true
	}
	return false
}
func main() {
	var year1, year2 int
	var a [1000]int
	cnt := 0
	fmt.Scan(&year1)
	fmt.Scan(&year2)
	for i := year1; i <= year2; i++ {
		if isrun(i) {
			a[cnt] = i
			cnt++
		}
	}

	fmt.Printf("%v\n", cnt)
	for i := 0; i < cnt; i++ {
		fmt.Printf("%d ", a[i])

	}
}
