package main

import "fmt"

func isprime(x int) bool {
	if (x == 1) || (x == 2) {
		return true
	} else {
		for i := 2; i*i < x; i++ {
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
	res := isprime(x)
	if res {
		fmt.Printf("YES")
	} else {
		fmt.Printf("NO")
	}
}
