package main

import (
	"fmt"
)

func generate(ch chan int) {
	for i := 2; ; i++ {
		ch <- i //把2, 3, 4, 5, ...写入ch
	}
}

func filter(in chan int, out chan int, prime int) {
	for {
		num := <-in
		if num%prime != 0 {
			out <- num //如果是质数， 写入out
		}
	}
}

func main() {
	ch := make(chan int)
	go generate(ch)
	for i := 0; i < 6; i++ {
		prime := <-ch // 从ch取出一个质数(后续filter后保证取出的是质数)
		fmt.Printf("prime:%d\n", prime)
		out := make(chan int)
		go filter(ch, out, prime)
		ch = out //我认为是类似把out的指针赋值给ch
	}
}

//代码实现了取质数的功能，体现了go的高并发性
