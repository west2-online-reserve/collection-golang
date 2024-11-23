package main

import (
	"fmt"
)

// 生成从 2 开始的所有自然数，并通过通道发送
func generate(ch chan int) {
	for i := 2; ; i++ {
		ch <- i // 将生成的自然数发送到通道
	}
}

// 从输入通道读取数字，筛选出不能被 prime 整除的数字，放入输出通道
func filter(in chan int, out chan int, prime int) {
	for {
		num := <-in // 从输入通道读取一个数
		// 如果 num 不能被 prime 整除，说明它不是 prime 的倍数，放入输出通道
		if num%prime != 0 {
			out <- num
		}
	}
}

func main() {
	ch := make(chan int) // 创建一个用于生成从 2 开始的自然数的通道
	go generate(ch)      // 启动生成自然数的协程
	for i := 0; i < 6; i++ {
		prime := <-ch // 从通道中获取下一个素数
		fmt.Printf("prime:%d\n", prime)
		out := make(chan int)     // 创建一个新的通道用于筛选后数字的传递
		go filter(ch, out, prime) // 启动一个筛选协程，将 prime 的倍数过滤掉
		ch = out                  // 更新通道，将下一次筛选的内容设为 out
	}
}

/*
该程序运用了协程和通道保证数据传递同步和线程安全
懒加载：数据处理是懒加载的，只有在需要时才会执行。
*/

/*
性能上讨论：
	普通素数筛选法：
		func simpleResolution() {
			n := 204729
			for num := 2; num <= n; num++ {
				isPrime := true
				for i := 2; i*i <= num; i++ {
					if num%i == 0 {
						isPrime = false
						break
					}
				}
				if isPrime {
					//fmt.Printf("prime: %d\n", num)
				}
			}
		}

		func main() {
			simpleResolution()
		}
		这种暴力算法适合小规模数据的筛选，当数据量极大的时候，效率极大下降

	多协程、多通道算法：
		通过查询资料，多协程、多通道似乎更适合IO操作，这种计算密集型任务并行化似乎不能利用这种多核的优势，反而会带来额外的协程和通道的开销
		因为创建协程和通道通信都有一定的系统开销。大量创建和销毁通道，频繁的通信操作，都会为系统增加负担，导致性能下降
		并且通道是阻塞的，需要等待通信完成，这也可能成为性能瓶颈，尤其是当协程之间存在大量的通信时

因而似乎不能在这种算法上很大的提高性能，对于大数据计算任务，可能分段筛选、减少通道通信和控制协程数量会好一些
*/
