package main

import "fmt"

func main() {
	var nums []int
	for i := 0; i < 50; i++ {
		nums = append(nums, i)
	}

	n := 0

	for j := 0; j < 50; j++ {
		if nums[j]%3 != 0 {
			nums[n] = nums[j]
			n++
		}
	}

	nums = nums[:n]

	nums = append(nums, 114514)
	fmt.Println(nums)
}
