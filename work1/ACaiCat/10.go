package main

import (
	"fmt"
	"sync"
)

var (
	current = 1
	m       = 6
	n       = 1114514
)

func main() {
	lock := &sync.Mutex{}
	c := sync.NewCond(lock)
	wg := &sync.WaitGroup{}

	for i := 0; i < m; i++ {
		wg.Add(1)
		go printNumber(i, c, wg)
	}

	wg.Wait()
}
func printNumber(id int, c *sync.Cond, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		c.L.Lock()
		if current%m != id && current <= n {
			c.Wait()
			c.L.Unlock()
			continue
		}

		if current > n {
			c.Broadcast()
			c.L.Unlock()
			return
		}

		fmt.Printf("worker-%d: %d\n", id, current)
		current++
		c.Broadcast()
		c.L.Unlock()
	}
}
