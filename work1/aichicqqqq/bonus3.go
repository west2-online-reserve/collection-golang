package main

import "fmt"

func main() {
	nums := [4]int{2, 7, 11, 15}
	var target = 9
	nums1 := &nums
	var i, j int
	var slice []int
	for i = 1; i < 4; i++ {
		for j = 0; j < i; j++ {
			if nums1[i]+nums1[j] == target {
				slice = append(slice, j)
				slice = append(slice, i)
			}

		}
	}
	fmt.Println(slice)

}
