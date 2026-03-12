package main

import (
	"fmt"
	"sync"
)

func queuePrint(n, m int) {
	queue := make([]int, 0, n)
	for i := 1; i <= n; i++ {
		queue = append(queue, i)
	}

	var mu sync.Mutex
	current := 1
	wg := sync.WaitGroup{}

	wg.Add(m)
	for i := 0; i < m; i++ {
		go func() {
			defer wg.Done()

			for {
				mu.Lock()

				if len(queue) == 0 {
					mu.Unlock()
					return
				}

				if queue[0] == current {
					fmt.Printf("%d ", queue[0])
					queue = queue[1:]
					current++
				}

				mu.Unlock()
			}
		}()
	}

	wg.Wait()
}

func main() {
	queuePrint(10, 3)
}
