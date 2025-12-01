package main

import (
	"fmt"
)

// 不断向通道 ch 发送从 2 开始的连续整数
func generate(ch chan int) {
	for i := 2; ; i++ {
		ch <- i
	}
}

// 从 in 通道接收数字，如果不能被 prime 整除，则发送到 out 通道
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
	// 启动 generate，往 ch 里塞数
	go generate(ch)
	// 循环 6 次，找出 6 个素数
	for i := 0; i < 6; i++ {
		// 从通道 ch 中接收一个素数
		prime := <-ch 
		fmt.Printf("prime:%d\n", prime)
		// 创建一个新的通道
		out := make(chan int)
		// 启动 filter，过滤掉 ch 中能被 prime 整除的数，剩下的送到 out
		go filter(ch, out, prime)
		// 将 out 赋值给 ch，下一轮循环继续从 ch 中接收数
		ch = out
	}
}

// 这个代码实现了什么功能？
// 通过生成连续整数，并不断过滤掉能被已知素数整除的数，输出前 6 个素数

// 这个代码利用了golang的什么特性？
// goroutine 和 channel

// 这个代码相较于普通写法，是否有性能上的提升？
// 性能大大降低，创建和管理的 goroutine 和 channel 的开销太大