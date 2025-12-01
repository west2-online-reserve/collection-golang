//4.AtCoder ARC017A：https://www.luogu.com.cn/problem/AT_arc017_1

package main

import (
	"fmt"
	"math"
)

func isPrime(x int) bool {
	if x < 2 {
		return false
	}
	if x == 2 {
		return true
	}
	if x%2 == 0 {
		return false
	}
	for i := 3; i < int(math.Sqrt(float64(x))); i++ {
		if x%i == 0 {
			return false
		}
	}
	return true
}

func main() {
	var input int
	fmt.Scan(&input)
	if isPrime(input) {
		fmt.Println("YES")
	} else {
		fmt.Println("NO")
	}
}
