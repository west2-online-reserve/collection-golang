### 运行下面代码，在你认为重要的地方写好注释，同时回答下面这些问题
   - 这个代码实现了什么功能？  
   **实现了计算素数的功能**
   - 这个代码利用了golang的什么特性？  
   **使用go关键字启动并发执行的goroutines充分利用多核处理器;使用通道实现了 goroutines 之间的通信，实现数据传递**
   - 这个代码相较于普通写法，是否有性能上的提升（性能提升：求解速度更快了）   
    **开启协程加快了计算速度。**

```go
package main

import (
	"fmt"
)

func generate(ch chan int) {
	for i := 2; ; i++ {
		ch <- i 
	}
}//从2开始不断向通道中发送连续整数

func filter(in chan int, out chan int, prime int) {
	for {
		num := <-in
		if num%prime != 0 {
			out <- num
		}   
	}
}//fileter函数从in通道接受数据，filt出不是素数的数，传给out

func main() {
	ch := make(chan int)
	go generate(ch)  //启动协程
	for i := 0; i < 6; i++ {
		prime := <-ch 
		fmt.Printf("prime:%d\n", prime)
		out := make(chan int)
		go filter(ch, out, prime)//启动协程，用之前筛选的质数进行过滤
		ch = out //将新生成的质数赋值给ch
	}
}
```
