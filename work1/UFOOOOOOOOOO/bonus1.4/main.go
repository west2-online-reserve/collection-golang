package main

import (
	"fmt"
)

func generate(ch chan int) {
	for i := 2; ; i++ {//没有i递增的条件，i无限递增
		ch <- i//向通道中输入>=2的整数
	}
}

func filter(in chan int, out chan int, prime int) {
	for {
		num := <-in//
		if num%prime != 0 {
			out <- num//
		}
	}
}

func main() {
	ch := make(chan int)//创建通道
	go generate(ch)//开启generate协程
	
	for i := 0; i < 6; i++ {//循环六次
		prime := <-ch//接收一个>=2的整数，防止generate协程死锁
		fmt.Printf("prime:%d\n", prime)
		out := make(chan int)//每次循环生成一个新的out通道
		go filter(ch, out, prime)
		ch = out//将ch替换为out，下一次循环中out作为输入通道，新的out通道作为输出通道，然后新的out通道代替旧的out通道
	}
	
}