package main

import "fmt"

/*
给定一个整数数组 nums 和一个整数目标值 target，请你在该数组中找出 和为目标值 target 的那两个 整数，并返回它们的数组下标。

你可以假设每种输入只会对应一个答案。但是，数组中同一个元素在答案里不能重复出现。

你可以按任意顺序返回答案。

示例 1：

输入：nums = [2,7,11,15], target = 9 输出：[0,1] 解释：因为 nums[0] + nums[1] == 9 ，返回 [0, 1]

示例2

输入：nums = [3,2,4], target = 6 输出：[1,2]

是否有复杂度O(n)的算法？
*/

func solve(nums []int, target int) (int, int) {
	mp := make(map[int]int)
	for i, num := range nums {
		need := target - num

		j, ok := mp[need]
		if ok {
			return i, j
		}
		mp[num] = i
	}
	return 0, 0
}

func main() {
	nums := []int{1, 2, 3, 4, 5}
	target := 5
	fmt.Println(solve(nums, target))
}
