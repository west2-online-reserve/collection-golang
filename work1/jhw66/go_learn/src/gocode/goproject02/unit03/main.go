package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	fmt.Println("随机整数：", rand.Int())        // 任意非负整数
	fmt.Println("随机 0~99：", rand.Intn(100)) // [0, 100)
	fmt.Println("随机浮点数：", rand.Float64())   // [0.0, 1.0)

	// 在 Go 1.20，官方更推荐 使用局部随机数生成器（非全局状态）。
	// 也就是说，不要修改全局随机源，而是创建自己的 rand.Rand 实例。
	// 创建一个独立随机源
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)

	fmt.Println(r.Intn(100)) // 0~99
	fmt.Println(r.Float64()) // 0.0~1.0
	// ✅ 优点：
	// 不依赖全局状态；
	// 线程安全；
	// 适合在多线程（goroutine）或库函数中使用；
	// 官方推荐写法。
}
