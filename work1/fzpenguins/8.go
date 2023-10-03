package main

import "fmt"

func twoSum(nums []int, target int) []int {
	numMap := make(map[int]int)

	for i, num := range nums {
		complement := target - num
		if _, found := numMap[complement]; found {
			return []int{numMap[complement], i}
		}
		numMap[num] = i
	}

	return nil // 如果没有匹配的结果，返回nil或其他适当的值
}

func main() {
	var nums = [4]int{2, 7, 11, 15}
	target := 9
	output := twoSum(nums[:4], target)
	fmt.Println(output)
}
