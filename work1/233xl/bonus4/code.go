package main

import (
	"fmt"
)

func generate(ch chan int) {
	for i := 2; ; i++ {
		ch <- i		
	}
}	//将生成的数字写入ch,这里生成2后不会继续生成,在prime收到后才会继续									*/

func filter(in chan int, out chan int, prime int) {
	for {
		num := <-in
		if num%prime != 0 {
			out <- num		//检测质数写入out
		}
	}
}

func main() {
	ch := make(chan int)    //创建通道ch
	go generate(ch)			//go并发生成数字
	for i := 0; i < 6; i++ {
		prime := <-ch       //第一次循环将2传给prime,ch继续生成变为3
		fmt.Printf("prime:%d\n", prime)
		out := make(chan int)
		go filter(ch, out, prime)
		ch = out            //out中的数写入ch,衔接下一轮循环中的prime
	}
}

//使用了2个并发子协程goroutine,生成质数速度快于单协程