package main
import (
		"fmt"
)
func generate(ch chan int) {
		for i := 2; ; i++ {
			ch <- i//遍历大于1的自然数
		}
}
func filter(in chan int, out chan int, prime int) {
		for {
			num := <-in
			if num%prime != 0 {		
				out <- num   //不断晒去此时ch中所有prime的倍数
			}
	}
}
func main() {
		ch := make(chan int)
		go generate(ch)
		for i := 0; i < 6; i++ {
			prime := <-ch      //接受下一个质数
			fmt.Printf("prime:%d\n", prime)	
		out := make(chan int)
		go filter(ch, out, prime) //建立新的一个prime过滤器直到main结束
		ch = out //将晒后的结果返回ch
	}
}
//此方法可依次得到i-1个质数
//这个方法利用go的多线程处理和管道
//处理上优势在于不必对每个数都重复的进行质数检查，相当于每得到一个数就可以晒掉一批数，节省了大量的处理时间