package main

import "fmt"

func main() {
	var temp int
	var nums []int
	for i := 0; i < 10; i++ {
		fmt.Scan(&temp)
		nums = append(nums, temp)
	}
	count := 0
	var max int
	fmt.Scan(&max)
	for i := 0; i < 10; i++ {
		if nums[i] <= max+30 {
			count++
		}
	}
	fmt.Println(count)
}
