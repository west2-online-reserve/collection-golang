package main

import (
	"fmt"
)

func isLeap(year int) bool {
	return (year%400 == 0) || ((year%100 != 0) && (year%4 == 0))
}

func main() {
	var min, max int
	cnt := 0
	fmt.Scanf("%v%v", &min, &max)
	leapMap := make(map[int]int, 10)
	for i := min; i <= max; i++ {
		if isLeap(i) {
			leapMap[cnt] = i
			cnt++
		}
	}
	fmt.Printf("%v\n", cnt)
	for i := 0; i < cnt; i++ {
		fmt.Printf("%v ", leapMap[i])
	}
}
