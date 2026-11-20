package main

import (
	"fmt"
)

func generate(ch chan int) {	// 生成自然数并塞入channel
	for i := 2; ; i++ {
		ch <- i
	}
}

func filter(in chan int, out chan int, prime int) {	// 筛选素数，用in通道传入数据，筛选后放入out通道
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
	for i := 0; i < 6; i++ {	// 控制筛选素数的个数
		prime := <-ch 
		fmt.Printf("prime:%d\n", prime)
		out := make(chan int)
		go filter(ch, out, prime)
		ch = out	// 将out通道赋值给ch通道，继续筛选
	}
}

/*
1.从2开始无限生成自然数并筛选素数
2.使用goroutine和channel实现
3.性能无提升
*/