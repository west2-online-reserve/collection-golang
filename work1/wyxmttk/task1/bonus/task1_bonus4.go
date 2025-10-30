package main

import (
	"fmt"
)

//原理:如果一个整数 n ≥ 2 不被 比它小的所有素数 整除，那么它一定是素数
//实现的功能:从2开始，打印指定个素数
//特性:闭包以及无容量管道的操作会堵塞的特性，前者实现协程的运行不会因为外部变量的引用改变而异常，后者我不知道有什么用
//性能提升:个人认为没有较明显的提升，因为本质上还是串行执行

// 永久持有ch的地址值，不会因为ch变量被赋别的值而影响协程工作(闭包特性)
func generate(ch chan int) {
	for i := 2; ; i++ {
		ch <- i
	}
}

// in 为传输 prime之后且不被prime之前所有素数整除的数 的管道
// out为 传输不被prime及之前所有素数整除的数 的管道，因此out管道的第一个元素一定是素数
func filter(in chan int, out chan int, prime int) {
	for {
		num := <-in
		//隐含num>prime的条件
		if num%prime != 0 {
			out <- num
		}
	}
}

func test() {
	//创建一个没有容量的管道，放入与拿取操作都会堵塞至另一端有相应操作接应它
	ch := make(chan int)
	//管道的值为地址值
	go generate(ch)
	//循环几次就是获取几个素数
	for i := 0; i < 6; i++ {
		prime := <-ch
		fmt.Printf("prime:%d\n", prime)
		out := make(chan int)
		go filter(ch, out, prime)
		ch = out
	}
}
