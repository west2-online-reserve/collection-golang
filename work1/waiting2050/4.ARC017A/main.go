package main

import "fmt"

func isPrime(a int) bool {

	if a <= 1 || a % 2 == 0 {
		return false
	}

	for i := 3; i <= a / i; i += 2 {
		if a % i == 0{
			return false
		}
	}

	return true
}

func main(){

	var a int
	fmt.Scan(&a)

	if isPrime(a) {
		fmt.Println("YES")
	} else {
		fmt.Println("NO")
	}
}