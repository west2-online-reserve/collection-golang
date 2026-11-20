package main

import (
	"fmt"
)

func generate(ch chan int) { //创建一个通道并不断往里面添加2之后的数
	for i := 2; ; i++ {
		ch <- i
	}
}

func filter(in chan int, out chan int, prime int) { //传进来一个素数 in通道 out通道 in通过素数的筛选到Out里面
	//in到out里面的数应该是无穷多
	for {
		num := <-in //num不断更新为in的输出
		if num%prime != 0 {
			out <- num //如果num的因数没有prime 那么num输入到out通道
		}
	}
}

func main() {
	ch := make(chan int)
	go generate(ch) //运行添加in线程
	for i := 0; i < 6; i++ {
		prime := <-ch //in里面的数变为素数 一开始为2 第二轮为 3
		fmt.Printf("prime:%d\n", prime)
		out := make(chan int)     //添加out通道
		go filter(ch, out, prime) //筛选出无限不以prime为因数的数
		ch = out                  //经过初步筛选的数变为下一轮的数
	}
}

//整体就是不断的筛选 从最小数2 每轮更换最小的数做筛子 同时 筛子本身也是 素数满足要求 当满足足够的素数时 for 循环结束 两个分线程也结束
//那么可以可以不用 goroutine的方式来解决呢 go routine的优势在哪里

//尝试用一般方法仿写

/*
var in []int
for i:=2;;i++{
in=append(in,i)
陷入死循环 程序结束
那么上面这段代码的优势就是可以多任务同时处理 在获取足量数据的同时不会进入死循环 非常的灵活 动态 精准
*/

/* 关于m个线程打印n个数
比如多个线程打印123456789
那么应该是 一个线程打印 1 下一个打印 1后面的数字
但不能规定 1357 2468这样的打印 虽然打印出来了 但是无序
那么多个同样的线程 打印一个递增的全局变量如何
*/
