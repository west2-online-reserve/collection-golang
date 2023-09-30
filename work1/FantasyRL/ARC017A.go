package main

import (
	"fmt"
)

func main() {
	var a int
	fmt.Scanf("%d", &a)
	b := isPrime(a)
	if b == true {
		fmt.Println("YES")
	} else {
		fmt.Println("NO")
	}
}
func isPrime(x int) bool {
	for i := 2; i < x; i++ {
		if x%i == 0 {
			return false
		}
	}
	return true
}
