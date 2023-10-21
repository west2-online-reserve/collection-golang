package main

import "fmt"

func main() {
	var x, y, z int
	i := 0
	_ = z
	fmt.Scanf("%d %d", &x, &y)
	for z := x; z <= y; z++ {
		if z%4 == 0 && z%100 != 0 {
			i++
		}
	}
	fmt.Printf("%d\n", i)

	for z := x; z <= y; z++ {
		if z%4 == 0 && z%100 != 0 {
			fmt.Printf("%d ", z)
		}
	}
}
