// 运行下面代码，在你认为重要的地方写好注释，同时回答下面这些问题
// 这个代码实现了什么功能？
// 这个代码实现了获取前六个质数的功能（埃氏筛法）
// 这个代码利用了golang的什么特性？
// 这个代码运用了golong的协程和管道的阻塞即生产者-消费者模式，通过channel阻塞自动实现goroutine间的同步
// 这个代码相较于普通写法，是否有性能上的提升？（性能提升：求解速度更快了）
// 并非加快了，因为巨大的协程开销和频繁的通道操作o(n^2)次操作，使得比平常写法性能更差
// 只有当每个filter阶段有昂贵的I/O操作时才会更省时
package main

import (
	"fmt"
)

func generate(ch chan int) {
	for i := 2; ; i++ {
		ch <- i // 创建无缓冲管道的同步阻塞，发送会阻塞直到有人从ch接收
	}
}

func filter(in chan int, out chan int, prime int) {
	for {
		// 第一个filter持续从generate的ch取数，后续的filter持续从上一个filter的out读值
		num := <-in         // 阻塞点1：从上游通道读取，如果上游没有数据则阻塞
		if num%prime != 0 { //如果num%prime == 0，该数字被过滤掉，不会发送到out，循环继续
			out <- num // 阻塞点2：向下游通道发送，如果下游没有接收者则阻塞
		}
	}
}

func main() {
	ch := make(chan int)
	go generate(ch) // 启动goroutine生成数字序列(2,3,4,5,6...)
	for i := 0; i < 6; i++ {
		prime := <-ch // 第1次循环：从generate读取2 后续循环：从上一个filter的输出通道读取质数
		//主程序的 <-ch 也起到驱动整个管道流动的作用
		fmt.Printf("prime:%d\n", prime) // 打印该质数
		out := make(chan int)           // 为新的filter创建输出通道
		go filter(ch, out, prime)       // 启动过滤协程，过滤掉能被prime整除的数
		ch = out                        // 更新ch指向新filter的输出通道，现在主程序将从这个新通道读取下一个质数
	}
} // main结束时，所有子goroutine会被回收
