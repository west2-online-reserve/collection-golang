// 最坏情况下 o(n²)

package main

import (
	"fmt"
)

func main() {
	nums := []int{7, 3, 11, 2}
	target := 9

	for x := 0; x < len(nums); x++ {
		for y := x; y < len(nums); y++ {
			sum := nums[x] + nums[y]
			if sum == target {
				fmt.Printf("[%d, %d]\n", x, y)
				return
			}
		}
	}
	fmt.Println("No Any Answer")

}
