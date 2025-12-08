package main

import "fmt"

func main() {
	var len int
	var target int
	fmt.Print("输入数组长度: ")
	fmt.Scanln(&len)
	nums := make([]int, len)
	fmt.Println("每行输入一个数组元素: ")
	for i := 0; i < len; i++ {
		fmt.Scanln(&nums[i])
	}
	fmt.Print("target= ")
	fmt.Scanln(&target)
	fmt.Println(findThem(nums, target))
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
