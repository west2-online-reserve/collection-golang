package main

import "fmt"

func main() {
	slice := []int{}
	fmt.Println("请输入数组，元素间用空格隔开：")
	for c := 'a'; ; {
		var data int
		fmt.Scanf("%d%c", &data, &c)
		slice = append(slice, data)
		if c == '\n' {
			break
		}
	}
	fmt.Println("请输入目标和：")
	var sum int
	fmt.Scanf("%d", &sum)
	fmt.Println(sum)

	gap := []int{}
	for _, value := range slice {
		gap = append(gap, sum-value)
	}

	for i := 0; i < len(slice); i++ {
		for j := i + 1; j < len(gap); j++ {
			if gap[j] == slice[i] {
				fmt.Printf("[%d,%d] ", i, j)
				fmt.Printf("Num[%d](%d)+Num[%d](%d)=%d\n", i, slice[j], j, slice[i], sum)
				return
			}
		}
	}
}
