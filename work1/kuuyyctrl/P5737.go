package main

import "fmt"

func main() {
	var arr [10]int
	for i := 0; i < 10; i++ {
		fmt.Scan(&arr[i])
	}
	var hei int
	cot := 0
	fmt.Scan(&hei)
	hei += 30
	for i := 0; i < 10; i++ {
		if hei >= arr[i] {
			cot += 1
		}
	}
	fmt.Println(cot)
}
