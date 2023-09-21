package main

import "fmt"

func main() {
    // 创建切片并初始化数字1-50
    slice := make([]int, 0)
    for i := 1; i <= 50; i++ {
        slice = append(slice, i)
    }

    // 删除切片中为3的倍数的数字
    for i := len(slice) - 1; i >= 0; i-- {
        if slice[i]%3 == 0 {
            slice = append(slice[:i], slice[i+1:]...)
        }
    }

    // 在末尾添加数114514
    slice = append(slice, 114514)

    // 输出切片
    fmt.Println(slice)
}