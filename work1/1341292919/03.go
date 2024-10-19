package main

import "fmt"

func isLeapyear(year int) bool {
	return (year%4 == 0 && year%100 != 0) || (year%400 == 0)
}

func main() {
	var begin, end int
	var Leapyears []int
	fmt.Scanf("%d%d", &begin, &end)
	count := 0
	for i := begin; i <= end; i++ {
		if isLeapyear(i) {
			count++
			Leapyears = append(Leapyears, i)
		}
	}
	fmt.Println(count)
	//fmt.Println(Leapyears)
	for _, i := range Leapyears {
		fmt.Printf("%d ", i)
	}
}
