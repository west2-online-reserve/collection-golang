package main

import "fmt"

func isPrime(x int) bool {
	var n int = 0
	for i := 1; i <= x; i++ {
		if x%i == 0 {
			n = n + 1
		}
	}
	if n == 2 {
		return true
	} else {
		return false
	}
}
func main() {
	var N int
	fmt.Scanln(&N)
	if isPrime(N) == true {
		fmt.Println("YES")
	} else {
		fmt.Println("NO")
	}
}
