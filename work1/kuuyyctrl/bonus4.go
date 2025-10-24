package main

import (
	"fmt"
	"time"
)

// 1.这个代码实现了素数的枚举
// 2.利用了GO的并发特性
// 3.百度上叫埃拉托斯特尼筛法。此程序将并发与这个筛法结合，通过类似于递归的方法枚举N个素数，时间复杂度为n的2,但通过并发所以会更快一些
// 与常规做法根号N累加来枚举时间复杂度大概是N的3/2相比更慢，缺点在于没有给通道缓存空间，这样预存一些数据下次枚举出来就能更快，具体给多少空间应该与CPU性能有关。
// 不过用我的电脑还是跑不过常规做法。
func generate(ch chan int) {
	for i := 2; ; i++ {
		ch <- i
	}
}

func filter(in chan int, out chan int, prime int) {
	for {
		num := <-in
		if num%prime != 0 {
			out <- num
		}
	}
}

//func isprime(digit int) bool {
//	for i := 2; i*i <= digit; i++ {
//		if digit%i == 0 {
//			return false
//		}
//	}
//	return true
//}

func main() {
	start := time.Now()
	var prime int
	ch := make(chan int, 10)
	go generate(ch)
	for i := 0; i < 10000; i++ {
		prime = <-ch
		fmt.Printf("prime:%d\n", prime)
		out := make(chan int, 10)
		go filter(ch, out, prime)
		ch = out
	}
	//cot := 10000
	//for i := <-ch; cot != 0; i++ {
	//	if isprime(i) {
	//		prime = i
	//		fmt.Printf("prime:%d\n", prime)
	//		cot--
	//	}
	//}
	//第一次循环生成除2的倍数以外的数的通道
	//第二次循环生成除3及2的倍数以外的数的通道
	//第三次筛5
	duration := time.Since(start)
	fmt.Println(duration)
}
