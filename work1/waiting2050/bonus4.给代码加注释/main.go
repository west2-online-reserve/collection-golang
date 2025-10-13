// 1.实现功能：计算前6个素数。类似于埃氏筛法，每次将素数的倍数（即合数）筛掉，留下的就只有素数了
// 2.利用了golang的并发编程特性，协程实现生成和筛选并行执行，通道实现不同通道之间安全的数据传递
// 3.性能提升：感觉不如单线程的埃氏筛，当数量级大了以后多个协程没有及时终止，效率反而更低了，占用的内存也高

package main

import (
	"fmt"
)

// 生成无限递增的整数序列并存入ch通道
func generate(ch chan int) {
	for i := 2; ; i++ {
		ch <- i
	}
}

// 在in通道筛选掉某个素数的倍数，并将未被筛掉的数保存到out通道
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
	go generate(ch)
	for i := 0; i < 6; i++ { // 循环六次，找到前六个素数
		prime := <-ch // 接收第一个数，这个数必定是素数，因为它已经被比他小的所有素数筛过一遍了
		fmt.Printf("prime:%d\n", prime)
		out := make(chan int) // 接收新一轮筛选过的数
		go filter(ch, out, prime)
		ch = out
	}
}