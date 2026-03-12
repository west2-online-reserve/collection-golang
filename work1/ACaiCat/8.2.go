// 偷看了提示说用hashMap, 这个确实是o(n)

package main

import (
	"fmt"
)

func main() {
	nums := []int{3, 7, 11, 2}
	target := 9

	hashMap := map[int]int{}
	for index, num := range nums {
		hashMap[num] = index
	}

	for numIndex, num := range nums {
		num2 := target - num

		if num2Index, ok := hashMap[num2]; ok && num2Index != numIndex {
			fmt.Printf("[%d, %d]\n", numIndex, num2Index)
			return
		}

	}

	fmt.Println("No Any Answer")

}
