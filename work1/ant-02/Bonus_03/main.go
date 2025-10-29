package main

import "fmt"

func findIdx(arr []int, target int) (int, int) {
	mp := make(map[int]int, 0)    
	for i, x := range arr {
		if idx, ok := mp[target-x]; ok {
			return idx, i
		}
		mp[x] = i
	}
	return -1, -1
}

func main() {
	var l, r int
	l, r = findIdx([]int{2, 7, 11, 15}, 9)
	fmt.Println(l, r)
	l, r = findIdx([]int{3, 2, 4}, 6)
	fmt.Println(l, r)
}
