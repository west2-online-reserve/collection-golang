package main

import "fmt"

func isprime(a int) bool {
	if a == 1 {
		return false
	}
	for i := 2; i*i <= a; i++ {
		if a%i == 0 {
			return false
		}
	}
	return true
}

func main() {
	var a int
	fmt.Scanf("%d", &a)
	if isprime(a) {
		fmt.Println("YES")
	} else {
		fmt.Println("NO")
	}

}
