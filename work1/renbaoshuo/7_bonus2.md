## Go 语言中的切片和数组的区别有哪些？

数组长度固定，在内存中连续存储。赋值或传递时会复制整个数组。

切片的长度是动态的。可以使用切片操作（如 append, []）来修改切片的容量或得到新的切片。复制或传递时只复制切片的引用，操作的是同一个底层数组。

## Go 中创建切片有几种方式？创建 map 呢？

### 切片

```go
slice1 := []type{}
slice2 := make([]type, len, cap) // len 和 cap 可选，len 长度，cap 初始容量（类似 C++ 中的 std::vector<>::reserve()），数据初始为 0
slice3 := arr[:] // 从数组中获取（arr 也可以是 slice）
```

### Map

```go
map1 := map[type]type{}
map2 := make(map[type]type, len) // len 可选
```
