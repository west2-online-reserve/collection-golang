package main

import (
	"fmt"
)

// 不停地向通道中传入数据
func generate(ch chan int) {
	for i := 2; ; i++ {
		ch <- i
	}
}

// 
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
	go generate(ch) // 开启协程向ch传入
	for i := 0; i < 6; i++ { // 循环6次
		prime := <-ch // 主线程不停接受数据
		fmt.Printf("prime:%d\n", prime)
		out := make(chan int)
		go filter(ch, out, prime)
		ch = out
	}
}
