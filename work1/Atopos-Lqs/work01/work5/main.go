package main

import "fmt"

func main() {
    // 创建一个切片并初始化为数字1到50
    slice := make([]int, 0)
    for i := 1; i <= 50; i++ {
        slice = append(slice, i)
    }

    // 创建一个临时切片，用于存放不是3的倍数的元素
    tempSlice := make([]int, 0)

    // 遍历原切片，将不是3的倍数的元素添加到临时切片中
    for _, num := range slice {
        if num%3 != 0 {
            tempSlice = append(tempSlice, num)
        }
    }

    // 清空原切片
    slice = nil

    // 将临时切片的内容复制回原切片
    slice = append(slice, tempSlice...)

    // 在切片末尾添加数114514
    slice = append(slice, 114514)

    // 输出切片
    fmt.Println(slice)
}