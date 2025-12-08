package main

import (
	"fmt"
	"sync"
)

func indexPrint(n, m int) {
	vis := make([]bool, n+7)
	vis[0] = true

	var mu sync.Mutex
	wg := sync.WaitGroup{}

	wg.Add(m)
	for id := 1; id <= m; id++ {
		go func(tid int) {
			defer wg.Done()

			num := tid
			for num <= n {
				mu.Lock()
				
				if vis[num - 1] {
					fmt.Printf("%d ", num)
					vis[num] = true
					num += m
				}

				mu.Unlock()
			}
		}(id)
	}

	wg.Wait()
}

func main() {
	indexPrint(10, 3)
}