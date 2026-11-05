package main

import "fmt"

func isprime(a int) bool {
	for i := 2; i*i <= a; i++ {
		if a%i == 0 {
			return false
		}
	}
	return true
}

func main() {
	var n int
	fmt.Scan(&n)

	if isprime(n) {
		fmt.Print("YES")
	} else {
		fmt.Print("NO")
	}
}
