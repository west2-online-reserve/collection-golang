//功能：筛选素数，并输出前六个素数
//特性：使用了goroutine并发处理，Channel确保 goroutine 之间能够安全通信
//性能提升：算法的时间复杂度是：O(N * loglogN)，执行速度大大提升
package main

import (
	"fmt"
)
//generate函数生成从2开始的递增整数并将它写入通道ch
func generate(ch chan int) {
	for i := 2; ; i++ {
		ch <- i
	}
}
//filter函数筛选出所有不能被prime整除的数并写入通道out
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
	//启动generate，并行goroutine
	go generate(ch)
	for i := 0; i < 6; i++ {//只输出前6个素数
		prime := <-ch 
		fmt.Printf("prime:%d\n", prime)//输出素数
		out := make(chan int)//创建输出管道
		go filter(ch, out, prime)
		ch = out//将out管道的数据福祉到ch管道中进行下一轮筛选
	}
}