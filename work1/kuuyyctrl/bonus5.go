package main

import (
	"fmt"
	"sync"
)

// 可以使用chan接受信息的同时并在一方输出时候对其上锁，以保证其有序性。
func main() {
	var lock sync.Mutex
	var num1 = [3]int{1, 2, 3}
	var num2 = [3]int{4, 5, 6}
	ch := make(chan int)
	wg := sync.WaitGroup{}
	wg.Add(6)
	go func() {
		for {
			fmt.Println(<-ch)
			wg.Done()
		}
	}()
	go func() {
		lock.Lock()
		for _, i := range num1 {
			ch <- i
		}
		lock.Unlock()
	}()
	go func() {
		lock.Lock()
		for _, i := range num2 {
			ch <- i
		}
		lock.Unlock()
	}()
	wg.Wait()
}
