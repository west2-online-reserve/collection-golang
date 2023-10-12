package main

import "fmt"

func main() {
	var x int
	var y int
	var sum int
	var slicelist []int
	fmt.Scanf("%d%d", &x, &y)
	for i := x; i <= y; i++ {
		if i%100 == 0 {
			if i%400 == 0 {
				slicelist = append(slicelist, i)
				sum += 1

			}

		} else {
			if i%4 == 0 {
				slicelist = append(slicelist, i)
				sum += 1
			}

		}

	}
	fmt.Println(sum)
	for i := 0; i <= (sum - 1); i++ {
		fmt.Printf("%d ", slicelist[i])
	}
}
