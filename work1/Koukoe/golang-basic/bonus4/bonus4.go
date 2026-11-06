package main

import (
	"fmt"
	"sync"
)

func main() {
	m := 61
	n := 67

	printInOrder(m, n)
}

func printInOrder(m, n int) {
	var wg sync.WaitGroup
	// 通道中只有一个数字在goroutine之间传递,所以同一时刻只有一个goroutine能处理和打印
	ch := make(chan int, 1)
	for i := 0; i < m; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for num := range ch {
				if num >= n {
					ch <- num
					return
				}
				// 每个goroutine只负责打印编号与数字取模结果相同的数字
				if num%m == id {
					fmt.Printf("线程 %d: %d\n", id, num)
					num++
				}
				ch <- num
			}
		}(i)
	}
	ch <- 0
	wg.Wait()
	close(ch)
}
