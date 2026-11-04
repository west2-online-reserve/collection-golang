package main

import (
	"fmt"
)

// 不停地向通道中传入数据
func generate(ch chan int) {
	for i := 2; ; i++ {
		ch <- i
	}
}

// 对所有不满足num%prime(2, 3, ...) == 0的num过滤
func filter(in chan int, out chan int, prime int) {
	// 无限循环过滤
	for {
		num := <-in
		if num%prime != 0 {
			out <- num // 若不可整除, 送给下一个过滤器
		}
	}
}

func main() {
	ch := make(chan int)
	go generate(ch) // 开启协程, 当ch中无数据时传递
	for i := 0; i < 6; i++ { // 循环6次
		prime := <-ch // 主线程接受数据
		fmt.Printf("prime:%d\n", prime)
		out := make(chan int) // 每次循环都新建一个和之前无关的out
		go filter(ch, out, prime) // 后台启动协程, 记录prime(2, 3, 5, ...)和当前的ch,out的指针,并不断判断ch值是否为质数(不可被prime整除)
		ch = out // 将out指针赋值给ch
	}
}

// 举个例子, 当ch为6:
// filter(ch1, out1, 2) // 6 % 2 == 0, 直接丢弃, 进入下一轮循环
// filter(out1, out2, 3) // 由于num := <-out1堵塞, 等待
// filter(out2, out3, 5) // 由于num := <-out2堵塞, 等待
// 后台generate执行, ch为7
// prime := 7, generate执行, ch为8

// ch从2开始, 循环6次, 每次记录ch传出的数据作为prime值(2, 3, 5, ...), 当新的ch可以整除之前的prime值(不是质数)时, 过滤
// 主线程(generate)
//     -fliter(2) 第一次循环加入
//     -fliter(3) 第二次循环
//     -fliter(5) 第三次循环
//     -fliter(7) ...
//     -fliter(11)
//     -fliter(13)
//     -...

// 实现了查找质数
// 使用了Go的协程和通道
// 相较于试根法(从2到√n+1循环), 这么写只记录质数值(prime), 每次for循环判断当前值%prime(2, 3, 5, ...)是否为0
// 减少了判断次数, 查找速度更快, 可以实现按需查找(要几个给几个)
