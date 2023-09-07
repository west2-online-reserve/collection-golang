package main

import (
	"fmt"
)

func generate(ch chan int) { //产生原始数据样本并发送给通道ch
	for i := 2; ; i++ {
		ch <- i
	}
}

func filter(in chan int, out chan int, prime int) { //质数筛，通过接受in通道的数据，借以能否整除质数判断是否为质数，将质数发送给通道out
	for {
		num := <-in
		if num%prime != 0 {
			out <- num
		}
	}
}

func main() {
	ch := make(chan int)     //创建通道ch
	go generate(ch)          //启用协程产生原始数据
	for i := 0; i < 6; i++ { //利用循环次数来控制筛选的质数个数
		prime := <-ch //接收质数
		fmt.Printf("prime:%d\n", prime)
		out := make(chan int)
		go filter(ch, out, prime) //启用协程进行质数筛
		ch = out                  //filter中筛出的数赋值给通道ch，以进行下一阶段质数筛
	}
}

//1.代码实现了质数筛，通过控制main函数中for循环次数来调整所得质数个数
//2.代码利用了go原生支持并发的特性，利用channel实现并发处理质数筛
//3.通过并发处理提高了求解速度
