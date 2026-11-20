package main

import "fmt"

func main() {
	var nums []int
	for i := 1; i <= 50; i++ {
		if i%3 != 0 {
			nums = append(nums, i)
		}
	}
	nums = append(nums, 114514)
	fmt.Println(nums)
}
