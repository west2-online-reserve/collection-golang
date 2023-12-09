package main

import (
	"fmt"
)

func generate(ch chan int) {
	for i := 2; ; i++ {
		ch <- i //向ch传输从2开始的所有整数
	}
}

// 接受in通道中的数据，判断是否能整除素数，并将筛出的素数发送到通道out
func filter(in chan int, out chan int, prime int) {
	for {
		num := <-in
		if num%prime != 0 {
			out <- num
		}
	}
}

func main() {
	ch := make(chan int)     //创建ch通道
	go generate(ch)          //并发产生未筛选数据
	for i := 0; i < 6; i++ { //指定筛选素数个数
		prime := <-ch
		fmt.Printf("prime:%d\n", prime)
		out := make(chan int)
		go filter(ch, out, prime) //并发进行素数筛选
		ch = out                  //将out通道数据传输给ch
	}
}
