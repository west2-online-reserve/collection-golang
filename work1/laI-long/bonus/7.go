package main

import (
	"fmt"
)

func answer(nums []int, target int) []int {
	numMap := make(map[int]int)
	for i, num := range nums {
		complement := target - num
		if j, ok := numMap[complement]; ok {
			return []int{j, i}
		}
		numMap[num] = i
	}
	return nil
}
func main() {
	var length int
	fmt.Println("输入数组长度")
	fmt.Scan(&length)
	nums := make([]int, length)
	fmt.Println("请逐个输入数组中的元素")
	for i := 0; i < length; i++ {
		fmt.Scan(&nums[i])
	}
	fmt.Println("请输入target")
	var target int
	fmt.Scan(&target)
	fmt.Println(answer(nums, target))
}
