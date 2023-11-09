package main

import "fmt"

func isPrime(x int, a int) bool {
	if x%a == 0 {
		return true
	} else {

		return false
	}

}
func main() {
	var n int
	var bool1 bool
	fmt.Scan(&n)
	for i := 2; i < n; i++ {
		temptBool := isPrime(n, i)
		if temptBool == true {
			bool1 = true
		}
	}
	if bool1 == true {
		fmt.Println("NO")
	} else {
		fmt.Println("YES")

	}

}
