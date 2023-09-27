package main

import "fmt"

func main() {
	var arr []int
	for i := 1; i <= 50; i++ {
		arr = append(arr, i)
	}
	for i := 0; i < len(arr); {
		if arr[i]%3 == 0 {
			arr = append(arr[:i], arr[i+1:]...)
		} else {
			i++
		}
	}
	arr = append(arr, 114514)
	for _, value := range arr {
		fmt.Printf("%d ", value)
	}
}
