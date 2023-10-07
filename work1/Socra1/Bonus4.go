package main

import (
	"fmt"
)

func generate(ch chan int) {
	for i := 2; ; i++ {
		ch <- i
	}
} //让ch接收所有自然数的函数

func filter(in chan int, out chan int, prime int) { //筛选出质数的函数，其原理为任何大于一的数字都可以表示为素数的乘积形式
	for {
		num := <-in
		if num%prime != 0 {
			out <- num
		} //判断是否能被质数整除并将质数放于out通道中
	}
}

func main() {
	ch := make(chan int) //创建ch通道
	go generate(ch)      //多协程启用generate函数
	for i := 0; i < 6; i++ {
		prime := <-ch
		fmt.Printf("prime:%d\n", prime)
		out := make(chan int)
		go filter(ch, out, prime) //多协程启用filter函数
		ch = out
	} //输出6个质数的For循环
}

//1.这个代码实现了质数的查找
//2.利用了golang的并发特性
//3.利用go与channel多协程并发，提高了质数筛选速度
