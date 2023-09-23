#### 1.Go语言中的切片和数组的区别有哪些？答案越详细越好。Go中创建切片有几种方式？创建map 呢？

1. GO语言中数组创建时要声明长度，切片无需声明长度
2. 数组不可以使用 `append`增加元素，切片可以使用`append`增加元素

1. GO可以使用`var name [] type`或 `var name []type = make([]type,len)` 或`name := make([]type,len)`创建切片

1. GO中创建 map 可以使用`name := make(map[keyType]valueType)` 或 
 ```
name := map([string]int){
        "zhangShang": 1
        "liShi": 0
 ```
