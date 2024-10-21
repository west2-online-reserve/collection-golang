package main

import "fmt"

func main() {
	var beginYear int32
	var endYear int32
	fmt.Scan(&beginYear, &endYear)
	var sum = 0
	var yearArray = make([]int32, 0)
	for i := beginYear; i <= endYear; i++ {
		if (i%4 == 0 && i%100 != 0) || (i%400 == 0) {
			sum += 1
			yearArray = append(yearArray, i)
		}
	}
	fmt.Println(sum)
	for _, i := range yearArray {
		fmt.Print(i)
		fmt.Print(" ")
	}
}
