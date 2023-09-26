package main

import "fmt"

func add(a, b int) int {
	return a + b
}
func main() {
	var (
		a int
		b int
	)
	_, err := fmt.Scan(&a)
	if err != nil {
		return
	}
	_, err = fmt.Scan(&b)
	if err != nil {
		return
	}
	fmt.Println(add(a, b))
}
