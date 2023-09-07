package main

import (
	"fmt"
)

func isPrime(n int) bool {
	for i := 2; i*i <= n; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func main() {
	var in int
	fmt.Scan(&in)
	if isPrime(in) {
		fmt.Print("YES")
	} else {
		fmt.Print("NO")
	}
}
