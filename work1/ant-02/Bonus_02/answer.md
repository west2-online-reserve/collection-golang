### 1 切片与数组的区别

- 数组的长度为固定，切片的长度可变。

- `len`和`cap`既可用于切片，也可用于数组，但数组的`len`和`cap`是固定的值且一样。

- 数组为Go语言单一类型，切片使用数组作为底层结构，包含容量、长度和数组指针，如下。

  ```go
  type SliceHeader struct {
  	Data uintptr
  	Len  int
  	Cap  int
  }
  ```

- 数组声明`var arr [0]int`默认为空数组，不等于`nil`,切片声明`var arr []int`为`nil`

- 数组传递通常是赋值传递，会分配新的内存空间，切片传递为引用传递，共用同一内存空间。

### 2 切片创建的几种方式

```go
var (
	a []int               // nil 切片, 和 nil 相等, 一般用来表示一个不存在的切片
	b = []int{}           // 空切片, 和 nil 不相等, 一般用来表示一个空的集合
	c = []int{1, 2, 3}    // 有 3 个元素的切片, len 和 cap 都为 3
	d = c[:2]             // 有 2 个元素的切片, len 为 2, cap 为 3
	e = c[0:2:cap(c)]     // 有 2 个元素的切片, len 为 2, cap 为 3
	f = c[:0]             // 有 0 个元素的切片, len 为 0, cap 为 3
	g = c[1:2]            // 有 1 个元素的切片, len 为 1, cap 为 2
	h = make([]int, 3)    // 有 3 个元素的切片, len 和 cap 都为 3
	i = make([]int, 2, 3) // 有 2 个元素的切片, len 为 2, cap 为 3
	j = make([]int, 0, 3) // 有 0 个元素的切片, len 为 0, cap 为 3
)
```

### 3 创建map

```go
// 仅声明
m1 := make(map[string]int)
// 声明时初始化
m2 := map[string]string{
	"Sam": "Male",
	"Alice": "Female",
}
```





