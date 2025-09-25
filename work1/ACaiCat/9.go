package main

import (
	"fmt"
)

// generate 生成从2开始的整数序列，并且向管道逐一发送。
func generate(ch chan int) {
	for i := 2; ; i++ {
		ch <- i
	}
}

// filter 过滤可以被当前filter储存prime整除的数。
// 从in管道接受从generate或者其他goroutine的filter传来的数，
// 并且排除掉可以被当前prime整除的数，并传给out管道输出或者传递给下一个goroutine的filter
func filter(in chan int, out chan int, prime int) {
	for {
		num := <-in
		if num%prime != 0 { // 筛选掉可以被prime整除的数，并且传给channel
			out <- num
		}
	}
}

func main() {
	ch := make(chan int)
	go generate(ch) // 创建一个用于生成并传递自增序列的生成器

	// 主循环，找出6个质数
	for i := 0; i < 6; i++ {
		prime := <-ch // 接收质数，即当前管道传递的第一个数
		fmt.Printf("prime:%d\n", prime)
		out := make(chan int)
		go filter(ch, out, prime) // 启动一个新的filter过滤当前的质数的倍数
		ch = out                  // 扩展管道链
	}
}
