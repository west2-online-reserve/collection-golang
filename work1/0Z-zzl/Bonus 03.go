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
	ch := make(chan int)
	go generate(ch)
	for i := 0; i < 6; i++ {
		prime := <-ch
		fmt.Printf("prime:%d\n", prime)
		out := make(chan int)
		go filter(ch, out, prime)
		ch = out
	}
}

//1.实现了筛选素数的作用，并以并发的方式实现的
//2.Goroutine轻量型线程
//channel运用通道实现各线程之间的联系与传输
//3.没有，因为为筛选素数会在运行期间创建大量的线程和通道，导致速度变慢
