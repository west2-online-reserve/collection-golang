package main

import "fmt"

func main() {
	var startYear, endYear int
	var list []int
	fmt.Scan(&startYear, &endYear)
	for i := startYear; i <= endYear; i++ {
		if i%400 == 0 || (i%4 == 0 && i%100 != 0) {
			list = append(list, i)
		}
	}
	fmt.Println(len(list))
	for i := 0; i < len(list); i++ {
		fmt.Print(list[i], " ")
	}
}
