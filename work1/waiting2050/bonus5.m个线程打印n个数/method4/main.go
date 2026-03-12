package main

import (
	"fmt"
	"sort"
	"sync"
)

func sortPrint(n, m int) {
	nums := make([]int, 0, n)
	mu := sync.Mutex{}
	wg := sync.WaitGroup{}

	wg.Add(m)
	for id := 1; id <= m; id++ {
		go func(tid int) {
			defer wg.Done()
			for num := id; num <= n; num += m {
				mu.Lock()
				nums = append(nums, num)
				mu.Unlock()
			}
		}(id)
	}

	wg.Wait()

	sort.Ints(nums)

	for _, num := range nums {
		fmt.Printf("%d ", num)
	}
}

func main() {
	sortPrint(10, 3)
}
