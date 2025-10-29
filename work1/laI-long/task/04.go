package main

import (
	"fmt"
)

func isPrime(x int) bool {
	var num int
	for i := 1; i <= x; i++ {
		if x%i == 0 {
			num++
		}
	}
	if num <= 2 && x != 1 {
		return true
	} else {
		return false
	}
}
func main() {
	var num1 int
	fmt.Scan(&num1)
	var test bool = isPrime(num1)
	if test {
		fmt.Println("Yes")
	} else {
		fmt.Println("No")
	}
}
