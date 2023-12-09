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
	go generate(ch)
	//当数据传入ch时，会进入filter判断是否为素数，是的话就传入out里，再由out传入ch，之后打印出来,打印六次数据后，循环结束，主程序也结束了。
	//一开始当ch里没有数据时，会把第一个传入的数据2打印出来
	//这个循环进行时，其他两个函数也在同时进行，out不断给ch传入数据,使ch的第一个数据始终为素数。
	for i := 0; i < 6; i++ {
		prime := <-ch
		fmt.Printf("prime:%d\n", prime)
		out := make(chan int)
		go filter(ch, out, prime)
		ch = out
	}
}
