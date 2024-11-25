```go
/*该段代码实现了打印前六个质数*/
package main

import (
	"fmt"
)
//生成从2开始的整数并传入通道ch
func generate(ch chan int) {
	for i := 2; ; i++ {
		ch <- i
	}
}
//将in中的整数过滤掉prime的整数再传给out
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
	go generate(ch)
	for i := 0; i < 6; i++ {//仅循环输出6个质数
		prime := <-ch//从ch中接收一个数 
		fmt.Printf("prime:%d\n", prime)
		out := make(chan int)
		go filter(ch, out, prime)//将in通道中的数过滤掉2的倍数后传送给out；
		ch = out//再将in通道中的数更新为已过滤的out通道里的数
	}
}
```

