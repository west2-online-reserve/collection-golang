/*4.运行下面代码，在你认为重要的地方写好注释，同时回答下面这些问题
这个代码实现了什么功能？
这个代码利用了golang的什么特性？
这个代码相较于普通写法，是否有性能上的提升？（性能提升：求解速度更快了
*/

package main

import (
	"fmt"
)

func generate(ch chan int) {
	for i := 2; ; i++ { // 无限循环，从2开始生成整数
		ch <- i
	}
}

func filter(in chan int, out chan int, prime int) {
	for {
		num := <-in
		if num%prime != 0 { // 只保留不能被当前质数整除的数字
			out <- num
		}
	}
}

func main() {
	ch := make(chan int)
	go generate(ch) // 启动goroutine生成自然数，与主程序并发执行

	// 循环6次，筛选出前6个质数
	for i := 0; i < 6; i++ {
		prime := <-ch
		fmt.Printf("prime:%d\n", prime)
		out := make(chan int) // 创建新的channel，用于传递本轮筛选后的数字
		go filter(ch, out, prime)
		ch = out // 更新channel为筛选后的输出channel，下一轮从这里接收数据
	}
}

/*
1.代码实现了质数筛选功能
2.利用了 Go 语言的并发特性
3.普通写法（如单循环标记非质数）需要预先分配内存存储数字范围，且筛选过程串行执行；而此代码通过多个协程并行处理不同质数的筛选步骤（每个质数对应一个filter协程），可利用多核 CPU 资源，且无需预先分配大内存（数据流式处理），在处理大规模数据时效率更高
*/
