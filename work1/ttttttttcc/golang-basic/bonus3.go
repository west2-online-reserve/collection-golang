package main

import (
	"fmt"
)

func generate(ch chan int) {
	for i := 2; ; i++ {
		ch <- i// 把自然数i发给通道ch
	}

}

func filter(in chan int, out chan int, prime int) {
	for {
		num := <-in
		if num%prime != 0 {
			out <- num// 如果num不能被prime整除，就把num发给输出通道out
		}
	}
}

func main() {
    ch := make(chan int) 
    go generate(ch) // 启动generate函数生成自然数序列

    for i := 0; i < 6; i++ {
        prime := <-ch // 从通道ch接收一个素数
        fmt.Printf("prime:%d\n", prime) // 打印素数

        out := make(chan int) 
        go filter(ch, out, prime) // 启动一个协程运行filter函数筛选数字
        ch = out 
    }
}