package main

import (
	"fmt"
)

func generate(ch chan int) {  //创建generate函数输出i到ch通道
	for i := 2; ; i++ {
		ch <- i
	}
}

func filter(in chan int, out chan int, prime int) {  //filter函数筛选质数
	for {
		num := <-in
		if num%prime != 0 {
			out <- num
		}
	}
}

func main() {
	ch := make(chan int)
	go generate(ch)  //新创建一个goroutine运行generate函数
	for i := 0; i < 6; i++ {
		prime := <-ch 
		fmt.Printf("prime:%d\n", prime)
		out := make(chan int)
		go filter(ch, out, prime)  //新运行一个goroutine运行filter函数
		ch = out
	}
}


/*该代码实现了输出前6个质数
利用了go的高并发的特点
提升了运行速度
*/
