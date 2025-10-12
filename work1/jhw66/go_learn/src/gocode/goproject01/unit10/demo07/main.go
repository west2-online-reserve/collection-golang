package main

import (
	"fmt"
	"time"
)

func main() {
	// 示例1：向已满管道写入的阻塞
	ch1 := make(chan int, 2)
	ch1 <- 1
	ch1 <- 2
	go func() {
		time.Sleep(2 * time.Second)
		<-ch1 // 2秒后取出一个元素，腾出空间
	}()
	fmt.Println("准备写入第三个元素...")
	ch1 <- 3 // 这里会阻塞2秒，直到goroutine取出数据
	fmt.Println("第三个元素写入成功")

	// 示例2：从空管道读取的阻塞
	ch2 := make(chan int, 2)
	go func() {
		time.Sleep(2 * time.Second)
		ch2 <- 100 // 2秒后发送数据
	}()
	fmt.Println("准备读取数据...")
	data := <-ch2 // 这里会阻塞2秒，直到goroutine发送数据
	fmt.Printf("读取到数据: %d\n", data)

	// 示例3：无缓冲管道的同步阻塞
	// 	同步阻塞的含义：
	// 发送操作阻塞：直到有接收方准备好接收数据
	// ch <- data // 阻塞，直到有人执行 <-ch
	// 接收操作阻塞：直到有发送方准备好发送数据
	// data := <-ch // 阻塞，直到有人执行 ch <- data
	// 同步机制：发送和接收必须同时发生，像接力棒传递
	ch3 := make(chan int) // 无缓冲
	go func() {
		fmt.Println("goroutine 准备发送数据")
		ch3 <- 42
		fmt.Println("goroutine 发送完成")
	}()
	time.Sleep(1 * time.Second)
	fmt.Println("主程准备接收数据")
	result := <-ch3
	fmt.Printf("主程收到: %d\n", result)
	time.Sleep(1 * time.Second)
}
