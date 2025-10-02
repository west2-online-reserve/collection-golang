package main // 第一行必须是package声明

import "fmt" //编译前记得保存

// 第一题：a+b
var (
	//用了全局，虽然没什么必要
	a   int64
	b   int64
	sum int64
)

func main() {
	fmt.Scanln(&a, &b)
	sum = a + b
	fmt.Println(sum)
}
