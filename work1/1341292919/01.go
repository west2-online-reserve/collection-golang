package main

import "fmt"

func getSum(a int, b int) int64 {
	return int64(a + b)
}
func main() {
	var a, b int
	fmt.Scanf("%d%d", &a, &b)
	fmt.Println(getSum(a, b))
}
