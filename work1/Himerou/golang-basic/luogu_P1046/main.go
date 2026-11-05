package main

import "fmt"

const chairHeight = 30

func main() {
	var temp int
	nums := make([]int, 0, 10)
	for i := 0; i < 10; i++ {
		_, _ = fmt.Scan(&temp)
		nums = append(nums, temp)
	}
	count := 0
	var maxHeight int
	_, _ = fmt.Scan(&maxHeight)
	for i := 0; i < 10; i++ {
		if nums[i] <= maxHeight+chairHeight {
			count++
		}
	}
	fmt.Println(count)
}
