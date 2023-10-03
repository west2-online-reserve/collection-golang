package main

import "fmt"

func isPrime(x int) bool {
	var i int
	var ret = true
	for i = 2; i < x; i++ {
		if x%i == 0 {
			ret = false
		}
	}
	return ret
}
func main() {
	var i int
	fmt.Scanf("%d", &i)
	if isPrime(i) {
		fmt.Println("yes")
	} else {
		fmt.Println("no")
	}

}
