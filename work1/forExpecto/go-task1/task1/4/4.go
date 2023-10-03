package main

import "fmt"

func isPrime(num int) bool {
	if num == 1 {
		return false
	}
	if num == 2 || num == 0 {
		return true
	}
	for i := 2; i < num; i++ {
		if num%i == 0 {
			return false
		}
	}
	return true
}
func main() {
	var num int
	fmt.Scan(&num)
	if isPrime(num) {
		fmt.Println("YES")
	} else {
		fmt.Println("NO")
	}
}
