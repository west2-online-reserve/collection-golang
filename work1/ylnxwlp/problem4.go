package main

import "fmt"

func main() {
	var x int
	fmt.Scan(&x)
	check(x)

}

func check(x int) {
	if x == 1 {
		fmt.Println("YES")
		return
	}
	for i := 2; i < x; i++ {
		if x%i == 0 {
			fmt.Println("NO")
			return
		}
	}
	fmt.Println("YES")
}
