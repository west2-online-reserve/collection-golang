package main

import "fmt"

func main() {
	var start, end int
	var arr []int
	_, _ = fmt.Scan(&start, &end)
	for i := start; i <= end; i++ {
		if (i%4 == 0 && i%100 != 0) || (i%400 == 0) {
			arr = append(arr, i)
		}
	}
	fmt.Println(len(arr))
	for i := 0; i < len(arr); i++ {
		fmt.Printf("%d ", arr[i])
	}
}
