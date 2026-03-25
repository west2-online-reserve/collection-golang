package main

import (
	"fmt"
	"io"
)

func twoSum(array []int, target int) []int {
	lookupTable := make(map[int]int)

	for index, num := range array {
		complement := target - num
		if val, ok := lookupTable[complement]; ok {   // Comma:ok := hashmap[ele]
			return []int{index, val}
		} else {
			lookupTable[num] = index
		}
	}

	return []int{}
}

func main() {
	var simpleNum int
	var numList []int

	for {
		_, err := fmt.Scan(&simpleNum)
		if err == io.EOF {
			fmt.Println("List Ready")     // 顺序matter
			break
		}
		if err != nil {
			fmt.Println("Invalid input")
			break
		}                                 //这段读取输入太变态了,给到拉完了   

		numList = append(numList, simpleNum)
	}

	var targetNum int
	fmt.Scan(&targetNum)

	fmt.Println(twoSum(numList, targetNum))
}

// 哈希表