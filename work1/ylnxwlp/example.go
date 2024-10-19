package main

import (
	"fmt"
)

// 生成从 2 开始的所有自然数，并通过通道发送
func generate(ch chan int) {
	for i := 2; ; i++ {
		ch <- i // 将生成的自然数发送到通道
	}
}

// 从输入通道读取数字，筛选出不能被 prime 整除的数字，放入输出通道
func filter(in chan int, out chan int, prime int) {
	for {
		num := <-in // 从输入通道读取一个数
		// 如果 num 不能被 prime 整除，说明它不是 prime 的倍数，放入输出通道
		if num%prime != 0 {
			out <- num
		}
	}
}

func main() {
	ch := make(chan int) // 创建一个用于生成从 2 开始的自然数的通道
	go generate(ch)      // 启动生成自然数的协程
	for i := 0; i < 6; i++ {
		prime := <-ch // 从通道中获取下一个素数
		fmt.Printf("prime:%d\n", prime)
		out := make(chan int)     // 创建一个新的通道用于筛选后数字的传递
		go filter(ch, out, prime) // 启动一个筛选协程，将 prime 的倍数过滤掉
		ch = out                  // 更新通道，将下一次筛选的内容设为 out
	}
}

//该程序运用了协程和通道保证数据传递同步和线程安全
//懒加载：数据处理是懒加载的，只有在需要时才会执行。
//性能上的提升：
//并发处理： 由于每个筛选器都运行在不同的协程中，多个筛选器可以同时处理不同的数字，这样就实现了并行筛选，处理大量数据时能够提高性能
//惰性生成： 数字的生成和筛选是按需进行的，未经过滤的数字会在进入下一个协程时被即时处理，从而避免了不必要的存储与运算
