package main

import (
	"fmt"
)

func generate(ch chan int) {
	for i := 2; ; i++ {
		ch <- i //将i传送给通道ch
	}
}

func filter(in chan int, out chan int, prime int) {
	for {
		num := <-in //从in通道中获取数据
		if num%prime != 0 {
			out <- num //将结果传送给out
		}
	}
}

func main() {
	ch := make(chan int) //创建通道
	go generate(ch)      //开启一个新goroutine
	for i := 0; i < 6; i++ {
		prime := <-ch //从ch获取数据
		fmt.Printf("prime:%d\n", prime)
		out := make(chan int)     //创建out通道用于接收结果
		go filter(ch, out, prime) //开启goroutine，用于搜集质数
		ch = out                  //将out中的结果传给ch，用于下一轮的质数判断
	}
}

/*
1.该代码实现了输出质数的功能
2.利用了go天生支持高并发的特性
3.将输入数据和判断质数同时进行，提升了性能
*/
