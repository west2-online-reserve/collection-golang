package main

import "fmt"

func mary(n int) bool {
	for i := 2; i <= (n - 1); i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func main() {
	var n int
	fmt.Scan(&n)
	if mary(n) {
		fmt.Println("Yes")
	} else {
		fmt.Println("No")
	}

}
