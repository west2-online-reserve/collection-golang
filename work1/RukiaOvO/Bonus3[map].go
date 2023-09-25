package main

import (
	"fmt"
)

func main() {
	var numlist []int
	var temp int

	for {
		_, err := fmt.Scanf("%d", &temp)
		if err != nil {
			break
		}

		numlist = append(numlist, temp)
	}
	fmt.Scan(&temp)

	list := make(map[int]int, len(numlist))
	for i := 0; i < len(numlist); i++ {
		list[numlist[i]] = i
	}
	for i := 0; i < len(numlist); i++ {
		val, ok := list[temp-numlist[i]]
		if ok && (temp != 2*numlist[i]) {
			fmt.Printf("%d %d", i, val)
			return
		}
	}
	fmt.Println("None")
}

// O(n)
