package main

import "fmt"

func main() {
	var a, b, c, d, e, f, g, h, i, j, x, n int
	fmt.Scanf("%d%d%d%d%d%d%d%d%d%d\n", &a, &b, &c, &d, &e, &f, &g, &h, &i, &j)
	fmt.Scanf("%d", &x)
	if a <= x+30 {
		n++
	}
	if b <= x+30 {
		n++
	}
	if c <= x+30 {
		n++
	}
	if d <= x+30 {
		n++
	}
	if e <= x+30 {
		n++
	}
	if f <= x+30 {
		n++
	}
	if g <= x+30 {
		n++
	}
	if h <= x+30 {
		n++
	}
	if i <= x+30 {
		n++
	}
	if j <= x+30 {
		n++
	}
	fmt.Println(n)
}
