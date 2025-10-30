package main

import (
	"fmt"
	"math/rand"
)

func getNumber(out chan int) {
	for {
		out <- rand.Intn(10000)
	}
}

func main() {
	n := 10
	m := 3
	ch := make([]chan int, m)

	for i := 0; i < m; i++ {
		ch[i] = make(chan int)
		go getNumber(ch[i])
	}
	for i, j := 0, 0; i < n; i, j = i+1, (j+1)%m {
		x := <-ch[j]
		fmt.Printf("channel %d: %d\n", j, x)
	}
}
