package main

import (
	"fmt"
	"sync"
)

var (
	consumerNum int = 5
	producerNum int = 1
	wg          sync.WaitGroup
)

func main() {
	wg.Add(consumerNum + producerNum)
	OrderlyPrinting([]int{1, 3, 10, 19, 20, 21, 22, 23, 24, 100, 101})
	wg.Wait()
}

// 假设切片有序
func OrderlyPrinting(slice []int) {
	var ch = make(chan int)
	var signal = make(chan int, 1)
	signal <- 1
	go func() {
		for _, v := range slice {
			<-signal
			ch <- v
		}
		close(ch)
		wg.Done()
	}()
	for i := 0; i < consumerNum; i++ {
		go func() {
			for v := range ch {
				fmt.Println(v)
				signal <- v
			}
			wg.Done()
		}()
	}
}
