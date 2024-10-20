// 使用哈希表，时间复杂度为O(n)
package main

import "fmt"

func twoSum(nums []int, target int) []int {
	//创建一个哈希表，用于存储数组元素的值和对应的索引
	numMap := make(map[int]int)

	//遍历数组
	for i, num := range nums {
		//计算与当前元素相加需要的值
		complement := target - num
		// 如果哈希表中存在需要的值，则返回结果
		if j, ok := numMap[complement]; ok {
			return []int{j, i}
		}
		//将当前元素的值和索引存入哈希表
		numMap[num] = i
	}

	//如果没有找到答案，则返回空切片
	return []int{}
}

func main() {
	nums := []int{2, 7, 11, 15}
	target := 9
	result := twoSum(nums, target)
	fmt.Println("[", result[0], ",", result[1], "]")
}
