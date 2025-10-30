// 解释用法以及运行结果
// 1.运行结果：此程序用于筛选前六个素数
package main

import (
	"fmt"
)

// 从2开始持续生成整数序列，并把生成结果传送到通道中
func generate(ch chan int) {
	for i := 2; ; i++ {
		ch <- i
	}
}

// 从通道中输入整数序列，对整数序列进行处理，筛出不能被整除的数，传送到输出通道中
func filter(in chan int, out chan int, prime int) {
	for {
		num := <-in
		if num%prime != 0 {
			out <- num
		}
	}
}
func main() {
	//ch是由chan组成的切片，用于存放通道中的数
	ch := make(chan int)
	//运行函数generate，生成整数序列，放到ch中
	go generate(ch)
	//循环遍历第1-6个素数
	for i := 0; i < 6; i++ {
		//prime即为ch中最后存放的素数
		prime := <-ch
		fmt.Printf("prime:%d\n", prime)
		out := make(chan int)
		//运行函数filter，筛出素数，返回到ch中
		go filter(ch, out, prime)
		ch = out
	}
}
