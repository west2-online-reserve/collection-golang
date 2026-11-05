package main

import (
	"bufio"
	"fmt"
	"os"
)

func ln(x int) bool {
	if x%4 == 0 && x%100 != 0 {
		return true
	} else if x%100 == 0 && x%400 == 0 {
		return true
	} else {
		return false
	}
}
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var x, y int
	fmt.Fscan(in, &x, &y)
	z := 0
	s := []int{}
	for i := x; i <= y; i++ {
		if ln(i) == true {
			z++
			s = append(s, i)
		}
	}
	fmt.Fprintln(out, z)
	for i := 0; i < len(s); i++ {
		fmt.Fprintf(out, "%d ", s[i])
	}
}
