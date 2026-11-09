package main

import "fmt"

func main() {
	var a, b int
	fmt.Scan(&a, &b)
	var res []int
	for i := a; i <= b; i++ {
		sum := 0
		for j := 1; j < i; j++ {
			if i%j == 0 {
				sum += j
			}
		}
		if sum == i { // 完数定义：所有真因子之和等于它本身
			res = append(res, i)
		}
	}
	for _, v := range res {
		fmt.Printf("%d ", v)
	}
	if len(res) == 0 {
		fmt.Println()
	}
}
