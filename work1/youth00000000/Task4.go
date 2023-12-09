package main

import "fmt"

func isPrime(x int) bool {
	for i := 2; i < x; i++ {
		if x%i == 0 {
			return false
		}
	}
	return true
}

func main() {
	var num int
	fmt.Scanln(&num)
	ret := isPrime(num)
	if ret == true {
		fmt.Println("YES")
	} else {
		fmt.Println("NO")
	}

}
