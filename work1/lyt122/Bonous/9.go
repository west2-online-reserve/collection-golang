package Bonous

// 1.这个代码实现了用筛法筛选素数的功能
// 2.利用goroutine和channel的特性
// 3.有性能上的提升,函数运行速度更快了
import (
	"fmt"
)

// 遍历从2开始的所有整数
func generate(ch chan int) {
	for i := 2; ; i++ {
		ch <- i
	}
}

// 筛出不是合数的数字
func filter(in chan int, out chan int, prime int) {
	for {
		num := <-in
		//是素数的话就把数字传到out里
		if num%prime != 0 {
			out <- num
		}
	}
}

func main() {
	ch := make(chan int)
	//把所有整数传到ch里
	go generate(ch)
	//循环获取素数
	for i := 0; i < 6; i++ {
		//把ch里的数据给prime
		prime := <-ch
		fmt.Printf("prime:%d\n", prime)
		out := make(chan int)
		//把素数传到out里
		go filter(ch, out, prime)
		//把素数传到ch里
		ch = out
	}
}
