package main

import (
	"fmt"
)

// 生成整数并将其发送到ch
func generate(ch chan int) {
	for i := 2; ; i++ {
		ch <- i
	}
}

func filter(in chan int, out chan int, prime int) {
	for {
		//num从in(ch)中接收一个值
		num := <-in
		//如果num可以被prime整除，则将num的值发送到out
		//这将可以被整除的数删掉了
		if num%prime != 0 {
			out <- num
		}
	}
}

func main() {
	//定义ch
	ch := make(chan int)
	//启用generate，生成整数并将其发送到ch
	go generate(ch)
	for i := 0; i < 6; i++ {
		//prime从ch中接收一个值
		prime := <-ch
		//打印接收到的值
		fmt.Printf("prime:%d\n", prime)
		//定义out
		out := make(chan int)
		//使用filter
		go filter(ch, out, prime)
		//ch从out接收数据
		ch = out
	}
}

/*
1·这段代码通过一个无限循环从2开始生成数字
若这个数字能够被prime整除则从ch中移除
用来筛选质数
2·运用了并发，通道
3·性能有提升
*/
