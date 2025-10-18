package main

import "fmt"

func main() {
	var arr = make([]int, 50)
	for i := 0; i < 50; i++ {
		arr[i] = i + 1
	}
	var remaining int = 0
	for i := 0; i < len(arr); i++ {
		if arr[i]%3 != 0 {
			arr[remaining] = arr[i]
			remaining++
		}
	}
	arr = append(arr, 114514)
	fmt.Println(arr)
}
