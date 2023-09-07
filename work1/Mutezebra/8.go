package main

func generate(ch chan int) { //用于向ch中发送所有的自然数
	for i := 2; ; i++ {
		ch <- i
	}
}
func filter(in chan int, out chan int, prime int) { //过滤掉所有不能整除prime的素数到out通道
	for {
		num := <-in         //接受ch中的值用于筛选
		if num%prime != 0 { //判断是否为"素"数
			out <- num
		}
	}
}
func main() {
	ch := make(chan int)
	go generate(ch)          //开启generate像第一个ch中发送自然数
	for i := 0; i < 6; i++ { //开始一个循环输出6个素数
		prime := <-ch //从ch中读取素数，用于输出且用于在filter中筛选素数
		out := make(chan int)
		go filter(ch, out, prime)
		ch = out // 将ch通道更新为out通道
	}
}
