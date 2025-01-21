package main

import (
	"fmt"
	"math"
)

func main() {
	var n int
	fmt.Scan(&n)
	if isPrime(n) {
		fmt.Println("YES")
	} else {
		fmt.Println("NO")
	}
}

// 判断一个数是否为质数
func isPrime(n int) bool {
	if n <= 1 {
		return false
	}
	// 只需检查到sqrt(n)即可
	for i := 2; i <= int(math.Sqrt(float64(n))); i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}
