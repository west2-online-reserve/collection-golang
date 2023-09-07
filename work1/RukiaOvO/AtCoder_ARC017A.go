package main

import "fmt"

func isPrime(a int) bool {
	for i := 2; 2*i <= a; i++ {
		if a%i == 0 {
			return false
		}
	}

	return true
}

func main() {
	var num int
	_, _ = fmt.Scan(&num)

	if isPrime(num) == true {
		fmt.Println("YES")
	} else {
		fmt.Println("NO")
	}
}
