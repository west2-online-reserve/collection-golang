//1.该代码实现了筛选素数的功能,但只会输出前六个素数
//2.该代码利用了Go语言的goroutine和channel特性。
//	通过goroutine，可以在不同的线程上并行执行代码；通过channel，可以在goroutine之间安全地传递数据。
//3.求解速度有所提高

package main

import (
	"fmt"
)

// generate函数不断生成递增的整数，通过channel发送这些整数
func generate(ch chan int) {
	for i := 2; ; i++ {
		ch <- i
	}
}

// filter函数过滤掉所有能被prime整除的数，只将不是prime倍数的数发送到输出channel。  
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

	//启动generate goroutine来不断发送整数 
	go generate(ch)

	for i := 0; i < 6; i++ {
		prime := <-ch 
		fmt.Printf("prime:%d\n", prime)
		out := make(chan int)
		
		//启动filter函数过滤
		go filter(ch, out, prime)
		//更新channel，以便下一次循环从这个新的过滤后的channel中读取数 
		ch = out
	}
}
