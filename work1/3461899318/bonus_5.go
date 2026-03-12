package main

import (
	"fmt"
	"sync"
)

var m, n int

// 协程循环处理，直到收到终止信号
func printNum(ch chan int, wait *sync.WaitGroup) {
	defer wait.Done() // 确保协程退出时一定会减少计数
	for {
		i := <-ch // 循环接收数据，直到收到终止信号
		if i > n {
			ch <- i // 将终止信号放回通道，让其他协程也能收到
			return  // 退出协程
		}
		fmt.Printf("协程打印: %d\n", i)
		ch <- i + 1 // 发送下一个数字
	}
}

func main() {
	fmt.Printf("请输入协程数量和打印的数字数量(默认从1开始打印)：")
	fmt.Scan(&m, &n)

	// 创建带缓冲的通道（缓冲大小1），避免发送时立即阻塞
	ch := make(chan int, 1)
	var wait sync.WaitGroup

	// 先启动所有协程
	wait.Add(m)
	for i := 0; i < m; i++ {
		go printNum(ch, &wait)
	}

	// 启动打印流程（此时已有协程等待接收，不会阻塞）
	ch <- 1

	// 等待所有协程退出
	wait.Wait()
	close(ch) // 关闭通道，释放资源
}
