
/*
这份代码利用了 Go 语言的并发特性，多线程运行代码，使得质数筛选的过程更加高效。
*/

package main

import (
	"fmt"
)

func generate(ch chan int) {
	for i := 2; ; i++ {
		ch <- i
	}
}

func filter(in chan int, out chan int, prime int) {
	for {
		num := <-in

		if num%prime != 0 {
			out <- num
		}
	}
}

func main() {
	// Channel 在 Goroutine 之间通信
	ch := make(chan int)

	// Goroutine 并发
	// https://go.dev/tour/concurrency/1
	go generate(ch)

	for i := 0; i < 6; i++ {
		prime := <-ch
		fmt.Printf("prime: %d\n", prime)
		out := make(chan int)
		go filter(ch, out, prime)
		ch = out
	}
}
