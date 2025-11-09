package main

import "fmt"

func sum(a, b int) int {
	return a + b
}
func main() {
	var a int
	var b int
	fmt.Scan(&a, &b)
	fmt.Println(sum(a, b))
}
