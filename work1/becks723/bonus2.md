## Golang数组 vs 切片

Q：Go语言中的切片和数组的区别有哪些？答案越详细越好。Go中创建切片有几种方式？创建map 呢？

A：

- 长度区别：数组的长度在声明时就写死了，不可变；切片声明时不用写死长度，长度可由`make`或`append`等方法改变。

  ```go
  var arr [8]int  // 声明数组
  var sli  []int  // 声明切片
  
  sli = append(sli, 8)  // 切片的长度后续可变
  ```

- 传参区别：数组在传参时需要写死长度；切片不用，更灵活。

  ```go
  func funcArray(arr [8]int) {}
  func funcSlice(sli []int) {}
  ```

  另外，数组作为参数是值传递，方法内部的修改不会影响原有数组；切片则是<small>*引用传递*</small>（严格来说仍是值传递，传递切片头元素），外面的切片值将同步被修改。



创建切片的方式：

```go
// 1. var声明。声明一个没有长度的切片slice1，然后用make开辟3个长度
var slice1 []int
slice1 = make([]int, 3)

// 2. 声明一个初始长度为4的s2
slice2 := []int { 4, 3, 2, 1 }

// 3. 变体1
var slice3 []int = make([]int, 3)

// 4. 变体2
slice4 := make([]int, 3)

// 5. make传入第三个参数，指定初始容量（capacity）
slice5 := make([]int, 3, 5)

// 但是不能这样写 ×
slice_err := []int  // 错误格式
```

在用`{}`<small>（初始化器）</small>时要注意：**换行初始化时即使最后一个元素后也要加逗号**

```go
s1 := []int { 1, 2, 3, 4 }  // 单行最后一个元素不用加逗号
s2 := []string {
    "one",
    "two",
    "three",
    "four",  // 多行要加！
}
```



创建map的方式：

```go
// 声明1
var map1 map[int]string  // 以int为键，string为值
map1 = make(map[int]string, 10)
map1[1] = "C"
map1[2] = "C++"
map1[3] = "Java"
map1[4] = "Python"

// 声明2
map2 := make(map[int]string)

// 声明3（注意每个键值对用逗号隔开）
map3 := map[string]string {
	"one": "C",
	"two": "C++",
	"three": "Java",
	"four": "Python",
}
```

