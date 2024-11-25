package main

import (
	"fmt"
)

// 从2到无穷向ch通道发送整数
func generate(ch chan int) {
	for i := 2; ; i++ {
		ch <- i
	}
}

// 从in通道读取数 把不能整除prime的数发送给out通道
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
	//启动线程1
	go generate(ch)
	//六轮
	for i := 0; i < 6; i++ {
		prime := <-ch
		fmt.Printf("prime:%d\n", prime)
		out := make(chan int)
		//启动线程2
		go filter(ch, out, prime)
		//把out作为下一轮的in
		ch = out
	}
}

//golang的多并发
//不需要等待一个进程结束再进行另外一个
