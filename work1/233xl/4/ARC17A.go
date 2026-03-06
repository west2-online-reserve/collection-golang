package main

import "fmt"

func isPrime(num int) bool {
	if num == 1 || num % 2 == 0 {
		return false
	}
	for i:= 3; i*i <= num; i += 2 {
		if num % i == 0 {
			return false
		}
	}
	return true
}

func main() {
	var n int
	_, _ = fmt.Scan(&n)

	if isPrime(n) {
		fmt.Println("Yes")
	} else {fmt.Println("No")}
}