/*
给定一个整数数组 nums 和一个整数目标值 target，
请你在该数组中找出 和为目标值 target 的那 两个 整数，并返回它们的数组下标。
*/
package main

import "fmt"

func main() {
	a := [...]int{1, 2, 3, 4, 5}
	//二次循环，在数组中找到和为目标值的两个整数
	for i := 0; i < len(a); i++ {
		for j := i + 1; j < len(a); j++ {
			if a[i]+a[j] == 6 {
				fmt.Printf("(%d,%d)", i, j) //返回下标
			}
		}
	}
}
