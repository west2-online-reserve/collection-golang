package main

import "fmt"

func isPrime(x int) bool {
	for i := 2; i <= x; i++ {
		if x%i == 0 {
			return false
		} else {
			return true
		}
	}
	return false
}

func main() {
	var a int
	fmt.Scanf("%d", &a)
	if isPrime(a) == false {
		fmt.Printf("No")
	} else {
		fmt.Printf("Yes")
	}
}
