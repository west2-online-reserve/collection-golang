package main

import (
	"fmt"
	"math"
)

func isZhi(a int) bool {
	var flag bool = true
	x := int(math.Sqrt(float64(a))) + 1
	for i := 2; i <= x; i++ {
		if a%i == 0 {
			flag = false
			break
		}
	}
	if a == 1 {
		flag = false
	}
	if a == 2 {
		flag = true
	}
	return flag
}
func main() {
	var a int
	fmt.Scanln(&a)
	if isZhi(a) {
		fmt.Println("YES")
	} else {
		fmt.Println("NO")
	}
}
