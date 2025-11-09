package main

import (
	"fmt"
	"sync"
)

func main() {
	m := 3  //线程为0到m-1
	n := 10 //打印0到n-1
	printInOrder(m, n)
}

func printInOrder(m, n int) {
	var wg sync.WaitGroup
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
