package main

import (
	"fmt"
)

//1. 这个代码实现了什么功能？
//输出前6个素数

//2.这个代码利用了golang的什么特性？
//Goroutine（轻量级线程）
//Channel（通道）

//3.这个代码相较于普通写法，是否有性能上的提升？（性能提升：求解速度更快了）
/*如果只是输出6个素数的话其实和普通写法差不多，但是大量筛选素数的话，
利用多线程同时生成和处理素数，确有性能提升。*/

// 功能：从2开始依次生成自然数，并传入管道中
// 入口参数：一个整型的通道
// 返回值：无
func generate(ch chan int) {
	for i := 2; ; i++ {
		ch <- i
	}
}

// 功能：筛选素数，并传入输出通道
// 入口参数：整型的输入和输出通道，和待检测的自然数
// 返回值：无
// 主要利用埃氏筛的算法
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
	go generate(ch) //启动一个生成自然数的子协程
	for i := 0; i < 6; i++ {
		prime := <-ch
		fmt.Printf("prime:%d\n", prime) //2不做判断直接输出，作为第一个素数进行后续的筛选
		out := make(chan int)
		go filter(ch, out, prime) //每一轮循环都会创建一个过滤任意某一素数倍数的函数的子协程
		ch = out                  //将上一轮过滤完的数字给到通道，并且能保证通道的第一个数字一定是素数
	}
}
