package main

import "fmt"

func main() {
	var length int
	var target int
	fmt.Print("输入数组长度: ")
	fmt.Scan(&length)
	nums := make([]int, length)
	fmt.Println("输入数组元素: ")
	for i := range nums {
		fmt.Scan(&nums[i])
	}
	fmt.Print("target = ")
	fmt.Scan(&target)

	result := findThem(nums, target)

	fmt.Print("[")
	for i, index := range result {
		if i > 0 {
			fmt.Print(",")
		}
		fmt.Print(index)
	}
	fmt.Println("]")
}

func findThem(nums []int, target int) []int {
	m := make(map[int]int)
	for i, j := range nums {
		if k, ok := m[target-j]; ok {
			return []int{k, i}
		}
		m[j] = i
	}
	return nil
}
