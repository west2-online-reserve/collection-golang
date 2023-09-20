package main

import (
	"fmt"
)

func generate(ch chan int) {
	for i := 2; ; i++ {
		ch <- i
	}
}

func filter(in chan int, out chan int, prime int) { //判断质数
	for {
		num := <-in
		if num%prime != 0 {
			out <- num //判断到质数后将其写入out
		}
	}
}

func main() {
	ch := make(chan int)
	go generate(ch) //启动generate()协程
	for i := 0; i < 6; i++ {
		prime := <-ch
		fmt.Printf("prime:%d\n", prime)
		out := make(chan int)
		go filter(ch, out, prime) //每次for都创建了一个filter，判断之后的数是否为质数
		//那读取了num的不是只有一个filter吗，怎么会这么神奇?????????
		ch = out //等候至读取out，将out写入ch，进入下一轮for循环
	}
}

//这个代码生成了6个质数
//利用了goroutine并发
//在速度上快于普通写法
