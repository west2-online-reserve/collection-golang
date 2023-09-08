package main

import "fmt"

func isPrime(x int) bool {
	if x == 1 {
		return false
	} else if x == 2 {
		return true
	} else {
		for i := 2; i*i <= x; i++ {
			if x%i == 0 {
				return false
			}
		}
		return true
	}
}
func main() {
	var x int
	fmt.Scanf("%d", &x)
	res := isPrime(x)
	if res {
		fmt.Println("YES")
	} else {
		fmt.Println("NO")
	}
}

