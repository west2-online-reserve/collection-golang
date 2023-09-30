package main

import "fmt"

func check(a int) int {
	if a%400 == 0 {
		return 1
	} else if a%100 != 0 && a%4 == 0 {
		return 1
	}
	return 0
}

func main() {
	var start, end int
	_, _ = fmt.Scanf("%d %d", &start, &end)

	count := 0
	list := make([]int, 0)

	for i := start; i <= end; i++ {
		if check(i) == 1 {
			count++
			list = append(list, i)
		}
	}

	fmt.Println(count)
	if count != 0 {
		fmt.Print(list[0])
		for i := 1; i < count; i++ {
			fmt.Printf(" %d", list[i])
		}
	}
}
