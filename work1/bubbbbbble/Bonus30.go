package main

import "fmt"

func main() {
	nums := []int{1, 4, 5, 9}
	target := 10
	hashmap := make(map[int]int)
	res := make([]int, 0)
	for index, value := range nums {
		diff := target - value
		i, exists := hashmap[diff]
		if exists {
			res = append(res, i, index)
			break
		}
		hashmap[value] = index
	}
	fmt.Println(res) //不知道为啥有两个bonus3，普通O(n方)的算法没写，O(n)的方法看了下gpt的思路

}
