package main

import "fmt"

func isPrime(x int) bool {
	if x < 2 {
		return false
	} else if x == 2 {
		return true
	} else if x%2 == 0 {
		return false
	} else {
		for j := 3; j*j <= x; j += 2 {
			if x%j == 0 {
				return false
			}
		}
		return true
	}
}
func main() {
	var x int
	fmt.Scanf("%d", &x)
	if isPrime(x) == true {
		fmt.Printf("YES\n")
	} else {
		fmt.Printf("NO\n")
	}
}
