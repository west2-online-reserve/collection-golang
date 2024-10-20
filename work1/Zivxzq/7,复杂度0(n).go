package main

import (
	"fmt"
)

func main() {
	//局限于数字较小的情况，数组应该开不了太大
	var tar int
	var num [10]int
	//标记每个合数的存在与否
	var innum [10000]int
	for i := 0; i < 10; i++ {
		fmt.Scan(&num[i])
		innum[num[i]] = i //标记存在，且指向原本下标
	}
	fmt.Scan(&tar)
	for i := 0; i < 10; i++ {
		if tar-num[i] < 0 {
			continue
		}
		if innum[tar-num[i]] != 0 {
			fmt.Println(i, innum[tar-num[i]])
			break
		}
	}

}
