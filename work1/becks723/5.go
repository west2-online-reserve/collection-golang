package main

import "fmt"

func main() {
	var slice []int = make([]int, 50)
	for i := 0; i < 50; i++ {
		slice[i] = i + 1
	}

	// 快慢指针
	j := 0
	for _, v := range slice {
		if v%3 != 0 {
			slice[j] = v
			j++
		}
	}
	slice = slice[:j] // 慢指针多计了1，回退

	slice = append(slice, 114514) // 在末尾追加114514

	for _, v := range slice {
		fmt.Printf("%d ", v)
	}
}
