package main

import "fmt"

func main() {
	nums := []int{2, 7, 11, 15}
	target := 9
	m := map[int]int{}
	for i, x := range nums {
		m[x] = i
	}
	res := [][]int{}
	for i, x := range nums {
		if j, ok := m[target-x]; ok {
			if i < j {
				res = append(res, []int{i, j})
			}
		}
	}
	fmt.Println(res)
}
