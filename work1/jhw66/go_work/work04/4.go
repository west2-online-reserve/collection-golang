package main

import "fmt"

func isPrime(n int) bool {
	if n <= 1 {
		return false
	}
	if n == 2 {
		return true
	}
	if n%2 == 0 {
		return false
	}
	for j := 3; j*j <= n; j += 2 {
		if n%j == 0 {
			return false
		}
	}
	return true
}

func main() {
	var N int
	fmt.Scan(&N)
	if isPrime(N) {
		fmt.Println("Yes")
	} else {
		fmt.Println("No")
	}
}
