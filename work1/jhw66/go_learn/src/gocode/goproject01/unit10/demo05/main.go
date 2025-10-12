package main

import "fmt"

//管道本质就是一个数据结构--队列
//自身线程是安全的，多协程访问时不需要加锁,不会发生资源争抢问题
//管道是引用类型，必须初始化后(make)才能写入数据
func main() {
	// 初始化容量为3的管道
	intChan := make(chan int, 3)

	// 向管道存数据（不超过容量）
	intChan <- 10
	num := 20
	intChan <- num
	intChan <- 40

	// 读取并打印
	fmt.Println(<-intChan)                                // 10
	fmt.Println(<-intChan)                                // 20
	fmt.Println("剩余:", len(intChan), "容量:", cap(intChan)) // 剩余: 1 容量: 3

	// 重新初始化容量合适的管道
	intChan = make(chan int, 100) // 增大容量以容纳所有数据

	// 写入数据
	for i := 9; i < 100; i++ {
		intChan <- i
	}

	// 关闭并遍历
	close(intChan)
	for v := range intChan {
		fmt.Printf("%v ", v)
	}
	fmt.Println()

	//声明只写管道：
	var intChan2 chan<- int = make(chan<- int, 3)
	intChan2 <- 20
	//声明只读管道：
	//var intChan3 <-chan int = make(<-chan int)

}
