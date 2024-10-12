package main

import (
	"fmt"
)

// 生成器函数
func generate(ch chan int) {
	for i := 2; ; i++ {
		ch <- i
	}
}

// 筛选器函数
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
	//执行generate函数，将生成的整数发送到ch
	go generate(ch)
	for i := 0; i < 6; i++ {
		//每从ch接收一个素数就打印出来
		prime := <-ch
		fmt.Printf("prime:%d\n", prime)
		out := make(chan int)
		//启动一个新的协程，执行filter函数，将筛选后的整数发送到out
		go filter(ch, out, prime)
		//将ch更新为out，以便下一个筛选器可以使用
		ch = out
	}
}

//功能
//生成一个素数流，使用多个筛选器来过滤这些数
//每个筛选器只允许不被当前素数整除的数通过

//这个代码利用了：协程、通道、并发、管道和迭代、匿名函数

//有，可以进行并行和流水线处理
