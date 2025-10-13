数组大小不可变，切片大小可变。
数组赋值时完全拷贝，切片赋值时浅拷贝
a := make([]int)
a := []int{}

a := map[int]int
a := make(map[int]int)