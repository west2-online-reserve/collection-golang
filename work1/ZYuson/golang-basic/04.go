package main

import "fmt"

func isPrime(x int) bool {
	if x == 1 {
		return false
	}
	if x <= 3 {
		return true
	}
	if x%2 == 0 || x%3 == 0 {
		return false
	}
	for i := 5; i*i <= x; i += 6 {
		if x%i == 0 || x%(i+2) == 0 {
			return false
		}
	}
	return true
}

func main() {
	var n int
	fmt.Scan(&n)
	if isPrime(n) {
		fmt.Println("YES")
	} else {
		fmt.Println("NO")
	}
}