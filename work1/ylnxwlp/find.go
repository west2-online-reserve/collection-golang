package main

import (
	"fmt"
)

func twoSum(nums []int, target int) []int {
	m := make(map[int]int)

	for i, num := range nums {
		targetNumber := target - num
		if index, found := m[targetNumber]; found {
			return []int{index, i}
		}
		m[num] = i
	}
	return []int{}
}

func main() {
	nums := []int{2, 7, 11, 15}
	target := 9
	fmt.Println(twoSum(nums, target))
}

//方法的时间复杂度是 O(n)，其中 n 是数组 nums 的长度
