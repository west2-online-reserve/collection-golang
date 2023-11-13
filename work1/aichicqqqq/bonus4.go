package main

import (
	"fmt"
)

// 开启channel1 从2开始循环
func generate(ch chan int) {
	for i := 2; ; i++ {
		ch <- i
	}
}

// 开启channel2 判断质数
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
	//找六个质数
	for i := 0; i < 6; i++ {
		prime := <-ch
		fmt.Printf("prime:%d\n", prime)
		out := make(chan int)
		go filter(ch, out, prime)
		ch = out
	}
}

/*1.该代码实现了多协程的寻找质数
2.利用了golang的并发编程支持
3.有 性能更快了*/
