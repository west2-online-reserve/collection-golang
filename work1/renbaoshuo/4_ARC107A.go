package main

import "fmt"

func main() {
	var n int

	fmt.Scan(&n)

	if isPrime(n) {
		fmt.Println("YES")
	} else {
		fmt.Println("NO")
	}
}

func isPrime(x int) bool {
	for i := 2; i*i <= x; i++ {
		if x%i == 0 {
			return false
		}
	}

	return true
}
