package main

import "fmt"

//常规做法双指针遍历复杂度n^2, 用哈希表查删改复杂度0(1),代替第二根指针这样复杂度就能达到0(1)

func main() {
	// 1 3 5 2 6
	cot := 0
	var n [100]int
	m := make(map[int]int)
	fmt.Println("当输入-1时停止输入数组")
	var x int
	fmt.Scanf("%d", &x)
	for ; x != -1; cot++ {
		n[cot] = x
		m[x] = cot
		fmt.Scanf("%d", &x)
	}
	var taget int
	fmt.Printf("输入taget:")
	fmt.Scanf("%d", &taget)
	for i := 0; i < cot; i++ {
		r, ok := m[taget-n[i]]
		if ok && r > i {
			fmt.Printf("[%d,%d]\n", i, r)
		}
	}
}
