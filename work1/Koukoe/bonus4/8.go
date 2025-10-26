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
			out <- num //除去能被当前的prime整除的数
		}
	}
}

func main() {
	ch := make(chan int)
	go generate(ch) //从2开始生成整数
	for i := 0; i < 6; i++ {
		prime := <-ch
		fmt.Printf("prime:%d\n", prime)
		out := make(chan int)
		go filter(ch, out, prime)
		ch = out //将过滤了能被prime整除的数后的通道传递给下一个循环
	}
}

//实现了从2开始打印6个质数的功能
//利用了golang的goroutine和channel特性
//相较于普通写法，没有性能上的提升
