# 结论：存在O(n)算法

> 结果如下

```go
package main

import "fmt"

func main() {
	nums := []int{}
	var x int
	var target int
    //假设通过检测到输入为0 而停止数组的输入
	for {
		if _, err := fmt.Scan(&x); err != nil || x == 0 {
			break
		}
		nums = append(nums, x)
	}
	fmt.Scan(&target)
	m_index := make(map[int]int)//使用map将已经出现的数字的下标对应起来
	for i, num := range nums {
		m_index[num] = i
		if index1, ok := m_index[target-num]; ok == true && index1 != i {
			fmt.Printf("[%d,%d]", index1, i)
		}
	}
}

```

