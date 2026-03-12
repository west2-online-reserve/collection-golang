package main

import "fmt"

func pass5737() {
	leftBoundary, rightBoundary := 1582, 3000
	var yearL, yearR int
	_, err := fmt.Scanf("%d %d", &yearL, &yearR)
	if err != nil {
		fmt.Println(err)
		return
	}
	if yearL < leftBoundary || yearR > rightBoundary || yearL > yearR {
		fmt.Println("bad input")
		return
	}
	var count int
	var arr []int = make([]int, 0)
	for yearL <= yearR {
		if yearL%4 == 0 {
			if yearL%100 != 0 || yearL%400 == 0 {
				arr = append(arr, yearL)
				count++
			}
		}
		yearL++
	}
	fmt.Println(count)
	for _, value := range arr {
		fmt.Printf("%d ", value)
	}
}
