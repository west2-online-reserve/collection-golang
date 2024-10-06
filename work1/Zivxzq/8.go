package main

import (
	"fmt"
)

func generate(ch chan int) {
	for i := 2; ; i++ { //无限循环
		ch <- i
	}
}

func filter(in chan int, out chan int, prime int) {
	for { //while（1）
		num := <-in
		if num%prime != 0 {
			out <- num //返回下一个素数
		}
	}
}

func main() {
	ch := make(chan int) //创建通道
	go generate(ch)      //启动，并行goroutine，就可以一直给ch输入++的整数
	for i := 0; i < 8; i++ {
		prime := <-ch
		fmt.Printf("prime:%d\n", prime)
		out := make(chan int) //创建新管道
		go filter(ch, out, prime)
		ch = out //返回素数
	}
}

/*
1.这个代码实验了素数的筛选算法，即 Eratosthenes 筛法。
2.运用了go的goroutine并行和channel通道的安全数据传输，两者结合，使数据传输变得简单安全。
3.查资料知道在面对更多素数时可以用埃拉托斯特尼筛法，即用bool数组去标记是否是素数，感觉跟上一题要求O(n)的代码实现类似。
*/
