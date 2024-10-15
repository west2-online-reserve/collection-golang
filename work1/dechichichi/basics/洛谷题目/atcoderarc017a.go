package main

import "fmt"

func isprime(x int) bool {
	if x < 2 {
		return false
	}
	for i := 2; i <= x/i; i++ {
		if x%i == 0 {
			return false
		}
	}
	return true
}

func main() {
	var x int
	fmt.Scan(&x)
	if isprime(x) {
		fmt.Print("YES")
	} else {
		fmt.Print("NO")
	}
}
