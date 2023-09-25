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

	for i := 0; i < len(numlist)-1; i++ {
		for j := i + 1; j < len(numlist); j++ {
			if numlist[i]+numlist[j] == temp {
				fmt.Printf("%d %d", i, j)
				return
			}
		}
	}
	fmt.Println("None")
}

// O(n^2)
