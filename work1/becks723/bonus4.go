/*
bonus4.go -- 打印前n个素数
*/
package main

import (
	"fmt"
)

/* 从2开始不断往管道ch中添加值 */
func generate(ch chan int) {
	for i := 2; ; i++ {
		ch <- i
	}
}

/* 用于从管道in中过滤掉prime的倍数，并且放进管道out中 */
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
	go generate(ch) // 开启一个goroutine无限往ch中加数字2，3，4，5，6……

	for i := 0; i < 6; i++ {
		prime := <-ch
		fmt.Printf("prime:%d\n", prime) // 从ch管道中取数并打印，这个数一定是素数
		out := make(chan int)
		go filter(ch, out, prime) // 开启一个新goroutine，将ch管道中的prime倍数全部筛去
		ch = out
	}
}

/*
  代码执行流程
round1: 从管道取值，为2，打印，筛去2的倍数（4，6，8……）
round2：取值，为3，打印，筛去3的倍数（9，15…… 6，12等已经在上轮被筛）
round3：取值，为5（4被筛），打印，筛去5的倍数（25，35……）
round4：取值，为7（6被筛），打印，筛去7的倍数
round5：取值，为11（8，9，10都被筛了），打印，筛去11的倍数
round6：取值，为13（12被筛），打印，筛去13的倍数
……
*/

/*
常规判素数写法：

func isPrime(num int) bool {
	if num < 2 {
		return false
	}
	for i := 2; i <= int(math.Sqrt(float64(num))); i++ {
		if num % i == 0 {
			return false
		}
	}
	return true
}
*/

/*
Q：这个代码实现了什么功能？
A：打印前n个素数，这个方法叫“埃拉托色尼筛法”。

Q：这个代码利用了golang的什么特性？
A：协程goroutine+管道channel。

Q：这个代码相较于普通写法，是否有性能上的提升？（性能提升：求解速度更快了）
A: 这个并行的埃拉托色尼筛法可能更适合演示golang语言的特性，但这样每碰到一个素数就开goroutine+channel的做法开销比较大，因此实际开发中不会选择这样写。
*/
