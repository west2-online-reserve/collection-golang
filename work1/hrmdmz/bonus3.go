package main

import "fmt"

func main() {
	m := map[int]int{}
	var n int
	var a int
	fmt.Scan(&n)
	for i := 0; i < n; i++ {
		fmt.Scan(&a)
		m[a] = i
	}
	var t int
	fmt.Scan(&t)
	for b, num := range m {
		c, ok := m[t-b]
		if c != num && ok {
			fmt.Printf("%d %d", num, c)
			break

		}
	}

}
