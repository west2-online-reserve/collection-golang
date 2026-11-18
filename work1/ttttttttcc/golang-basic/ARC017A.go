package main

import (
	"fmt"
	"math"
)

// 判断一个数是否为质数
func isPrime(x int) bool {
	if x <= 1 {
		return false
	}
	if x <= 3 {
		return true
	}
	if x%2 == 0 || x%3 == 0 {
		return false
	}
	sqrtX := int(math.Sqrt(float64(x)))
	for i := 5; i <= sqrtX; i += 6 {
		if x%i == 0 || x%(i+2) == 0 {
			return false
		}
	}
	return true
}

func main() {
	var n int
	fmt.Scan(&n) // 读取输入的整数

	if isPrime(n) {
		fmt.Println("YES") // 如果是质数，输出 YES
	} else {
		fmt.Println("NO")  // 如果不是质数，输出 NO
	}
}