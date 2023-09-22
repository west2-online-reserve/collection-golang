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
		out := make(chan int) //创建out，每轮都会创建出一个新的out
		go filter(ch, out, prime)
		//每次for都创建了一个新的filter，每个filter判断成功后就会传入out的value实际就是filter的ch，如此反复直到传给最后生成的out被prime接收，进行print
		ch = out //令ch覆写为out的地址
	}
}

//这个代码能够判断质数并打印
//利用了goroutine并发
//在速度上快于普通写法
