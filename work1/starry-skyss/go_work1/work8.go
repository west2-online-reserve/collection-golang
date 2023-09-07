package main

import "fmt"

// 以下为用数组进行运算
func main() {
	nums := [4]int{2, 7, 11, 15}
	var target int
	fmt.Scanln(&target)
	for k, v := range nums {
		for tmp := k + 1; tmp < 4; tmp++ {
			if (v + nums[tmp]) == target {
				fmt.Println(k, tmp)
			}

		}
	}

}

//以下为复杂度为O（n）的算法，只想到用map
/*
func main() {
	var num = map[int]int{
		2:  1,
		7:  2,
		11: 3,
		15: 4,
	}
	var target int
	fmt.Scanf("%d", &target)
	for k, v := range num {
		val, ok := num[target-k]
		if ok {
			fmt.Println(v, val)
			delete(num, k)
		}
	}
}
*/
