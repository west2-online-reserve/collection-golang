package main

import "fmt"

func main() {
	var a, b int
	fmt.Scanf("%d %d", &a, &b) // 确保输入的格式正确
	sum := a + b               // 使用 := 进行变量声明
	fmt.Println(sum)
}
