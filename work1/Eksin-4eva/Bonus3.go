package main

import "fmt"

func twoSum(nums []int, target int) []int {
	m := make(map[int]int)
	for i, num := range nums {
		complete := target - num
		if index, ok := m[complete]; ok {
			return []int{index, i}
		}
		m[num] = i
	}
	return nil
}

func main() {
	var n int
	fmt.Print("请输入数组长度: ")
	fmt.Scan(&n)

	nums := make([]int, n)

	fmt.Printf("请依次输入 %d 个数字 (用空格隔开): \n", n)
	for i := 0; i < n; i++ {
		fmt.Scan(&nums[i])
	}

	var target int
	fmt.Print("请输入Target: ")
	fmt.Scan(&target)

	result := twoSum(nums, target)
	if result != nil {
		fmt.Printf("%v", result)
	} else {
		fmt.Println("No result")
	}
}

// O(N)
