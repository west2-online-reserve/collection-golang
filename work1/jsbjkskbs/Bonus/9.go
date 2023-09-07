package main

import (
	"fmt"
)

// 生成一段从2到x的整型序列(个人感觉x的大小取决于主线程的结束时间(即main()函数的生命周期))
func generate(ch chan int) {
	for i := 2; ; i++ {
		ch <- i
	}
}

// 形似"欧拉筛"的素数判别法
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
	for i := 0; i < 6; i++ {
		prime := <-ch
		fmt.Printf("prime:%d\n", prime)
		out := make(chan int)
		go filter(ch, out, prime)
		ch = out
	}
}

//以下为个人答案，可能会有错误：
//1.这个代码实现了什么功能？
//  输出一段素数(只输出从2开始算起的6个素数)
//2.这个代码利用了golang的什么特性？
//  并发(通过go func()创建并发体，并通过chan实现并发体之间的通信)
//3.这个代码相较于普通写法，是否有性能上的提升？（性能提升：求解速度更快了）
//  采用欧拉筛的思路，比一般的素数判断法更快；
//  通过并发机制，将输出素数这一操作分配给三个协程进行操作，分别是生成序列、筛选序列、输出素数
//  原本单一的执行顺序被拆分为多项任务，使得硬件资源得到充分利用；
//  故相较于普通写法，可以有性能上的提升。
