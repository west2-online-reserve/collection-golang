package main

import (
	"fmt"
)

func generate(ch chan int) { //一直在生成不同ch值
	for i := 2; ; i++ {
		ch <- i //i发送到通道ch
	}
}

func filter(in chan int, out chan int, prime int) { //过滤器，筛选素数
	for { //无限循环
		num := <-in
		if num%prime != 0 {
			out <- num
		} //out为删去prime及其倍数的新通道
	}
}

func main() {
	ch := make(chan int) //声明通道，不带缓冲区,前一个数被接收后才会产生新一个数等待接收
	go generate(ch)      //开了个新的运行期线程运行generate函数
	for i := 0; i < 6; i++ {
		prime := <-ch
		fmt.Printf("prime:%d\n", prime)
		out := make(chan int)
		go filter(ch, out, prime)
		ch = out //用通道out覆盖通道ch
	}
}

/*1)输出前6小的素数
2）并发性
3）速度快于普通算法 */
