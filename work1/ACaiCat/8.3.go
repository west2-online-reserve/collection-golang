// AI优化的，把初始化hashMap和循环合并了

package main

import (
	"fmt"
)

func main() {
	nums := []int{3, 7, 11, 2}
	target := 9

	hashMap := map[int]int{}

	for numIndex, num := range nums {
		num2 := target - num

		if num2Index, ok := hashMap[num2]; ok {
			fmt.Printf("[%d, %d]\n", numIndex, num2Index)
			return
		}

		hashMap[num] = numIndex

	}

	fmt.Println("No Any Answer")

}
