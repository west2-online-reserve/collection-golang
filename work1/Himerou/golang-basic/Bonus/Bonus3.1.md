# 实现功能

实现了质数筛选

# 特性

使用了golang的并发编程的特性

# 提升

这个写法在性能上比普通写法更慢，计算到第n个质数的时候已经开启了n个goroutine ，channel通信很慢 这里的数字基本全是通过channel通信并且第n个质数至少通过n层通信，取模判断上一个质数要通过比他小的所有的质数的筛选器，比√n要大，更非时间

```go
package main

import (
	"fmt"
)

// generate 从 2 	开始递增的向管道发送数字
func generate(ch chan int) {
	for i := 2; ; i++ {
		ch <- i
	}
}

// filter 接收从in发送的数字，丢弃prime的整数倍数字，结果发送到out
func filter(in chan int, out chan int, prime int) {
	for {
		num := <-in
		if num%prime != 0 { //剔除所有prime的整数倍的数字
			out <- num
		}
	}
}

func main() {
	ch := make(chan int)
	go generate(ch)
	for i := 0; i < 6; i++ {
		prime := <-ch //当前管道的第一个数一定是质数
		fmt.Printf("prime:%d\n", prime)
		out := make(chan int)
		go filter(ch, out, prime) //接上上一级的过滤器，只有通过了上一级过滤器的数才会进入这一级（上一级的out变成这一级的in)
		ch = out                  //ch指向新的末端，下一轮从这里取质数，并进入下一级过滤
	}
}

```

