package main

import (
	"fmt"
)

func pl(nums []int, tagret int) []int {
	m := make(map[int]int)
	for i, v := range nums {
		d := tagret - v
		if idx, ok := m[d]; ok {
			return []int{idx, i}
		}
		m[v] = i
	}
	return nil
}
func main() {
	nums := []int{}
	var tagret int
	fmt.Printf("输入目标值：")
	fmt.Scan(&tagret)
	fmt.Printf("输入数组，输入-1停止输入\n")
	fmt.Printf("nums= ")
	var n int
	for {
		_, err := fmt.Scan(&n)
		if err != nil {
			break
		}
		if n == -1 {
			break
		}
		nums = append(nums, n)
	}
	y := []int{}
	y = append(y, pl(nums, tagret)...)
	fmt.Printf("%s%d%s%d%s", "[", y[0], ",", y[1], "]")
}
