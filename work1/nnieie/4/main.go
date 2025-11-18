package main

import "fmt"

func main() {
	var no int
	fmt.Scan(&no)
	if isPrime(no) {
		fmt.Printf("YES")
	} else {
		fmt.Printf("NO")
	}

}	
func isPrime(x int) bool {
	if x < 2 {
		return false
	}
	for i := 2; i*i <= x; i++ {
		if x%i == 0 {
			return false
		}
	}
	return true
}