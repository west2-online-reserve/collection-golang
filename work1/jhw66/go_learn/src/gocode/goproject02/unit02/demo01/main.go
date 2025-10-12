// new：分配内存，返回指针
// 语法   p := new(T)
// 作用:
// 为 类型 T 分配一块零值内存空间。
// 返回的是指向该类型的指针 (*T)。
// 不会初始化复杂结构，只是分配并清零。
// 适用对象:
// 所有类型都可以使用 new：包括基本类型（int、float、struct、数组等）。

// make：初始化内建引用类型
// 语法   v := make(T, args)
// 作用:
// 为 内建引用类型 分配并初始化内存。
// 返回的是该类型本身（不是指针）。
// 仅适用于三种类型：
// slice map chan 适用对象
// 仅这三种类型可用 make，否则会编译错误。

// new([]int)	*[]int	❌ 不可直接用，需要解引用	几乎不用
// make([]int, len, cap)	[]int	✅ 可直接用	✅ 推荐
// new(map[string]int)	*map[string]int	❌ 会 panic	几乎不用
// make(map[string]int)	map[string]int	✅ 可直接用	✅ 推荐
// new(chan int)	*chan int	❌ 死锁风险	几乎不用
// make(chan int, n)	chan int	✅ 可直接用	✅ 推荐

// 语法	runtime 调用	是否分配底层结构	返回值类型	是否可立即用
// new(T)	runtime.newobject	❌ 仅分配零值	*T	❌
// make([]T, n)	runtime.makeslice	✅ 分配底层数组	[]T	✅
// make(map[K]V)	runtime.makemap	✅ 分配桶结构	map[K]V	✅
// make(chan T)	runtime.makechan	✅ 分配 channel 缓冲区	chan T	✅

// 如果就是想new要配合make使用
package main

import "fmt"

type Manager struct {
	Name  string
	Tasks []string
	Data  map[string]int
	Msgs  chan string
}

func NewManager(name string) *Manager {
	m := new(Manager) // 分配零值结构体
	m.Name = name
	m.Tasks = make([]string, 0) // 初始化 slice
	m.Data = make(map[string]int)
	m.Msgs = make(chan string, 10)
	return m
}

func main() {
	mgr := NewManager("Alice")
	mgr.Tasks = append(mgr.Tasks, "Build API")
	mgr.Data["done"] = 1
	mgr.Msgs <- "hello"

	fmt.Println(mgr.Name)
	fmt.Println(mgr.Tasks)
	fmt.Println(mgr.Data)
	fmt.Println(<-mgr.Msgs)
}
