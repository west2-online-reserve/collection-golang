package main

import (
	"fmt"
)

func main() {
	var a int
	fmt.Scan(&a)
	if isPrime(a) {
		fmt.Println("YES")
	} else {
		fmt.Println("NO")
	}
}
func isPrime(a int) bool {
	for i := 2; i*i <= a; i++ {
		if a%i == 0 {
			return false
		}
	}
	return true
}
