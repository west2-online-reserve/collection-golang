package main

import "fmt"

// 判断是否为素数
func judge_prime() bool {
	n := 1
	fmt.Scanln(&n)
	if n <= 1 || n%2 == 0 {
		return false
	}
	for i := 3; i*i <= n; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func main() {
	//else必须紧跟在if的右大括号
	if judge_prime() {
		fmt.Println("YES")
	} else {
		fmt.Println("NO")
	}
}
