package main

import (
	"fmt"
)

// 从2开始生成数字，并发送到通道ch中
func generate(ch chan int) {
	for i := 2; ; i++ {
		ch <- i //发送
	}
}

// 从in中接受数字，并判断是否为素数，最后发送到out中
func filter(in chan int, out chan int, prime int) {
	for {
		num := <-in         //接受
		if num%prime != 0 { //判断
			out <- num //发送
		}
	}
}

func main() {
	ch := make(chan int)
	go generate(ch) //开一个协程，开始生成
	for i := 0; i < 100; i++ {
		prime := <-ch //接受上面产生的2，避免阻塞；更新素数,接收filter中的素数
		fmt.Printf("prime:%d\n", prime)
		out := make(chan int)
		go filter(ch, out, prime)
		ch = out //复制新素数到ch
	}
}

/*
Q1:生成6个素数
Q2:多个协程同时进行
Q3：应该有吧？毕竟生成的同时就进行了判断
*/
