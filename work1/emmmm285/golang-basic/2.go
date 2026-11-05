package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	a := make([]int, 10)
	for i := range a {
		fmt.Fscan(in, &a[i])
	}
	var h int
	fmt.Fscan(in, &h)
	b := 0
	for i := 0; i < 10; i++ {
		if h+30 >= a[i] {
			b++
		}
	}
	fmt.Fprintln(out, b)
}
