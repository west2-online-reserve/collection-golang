# Bonus

## Go语言中切片和数组的区别有哪些？
1. Go语言中的切片可以动态扩充内存，而数组只能在声明时指定大小
2. 切片可以从数组上截取
3. 数组可以是多维的，切片只能是一维的
4. 数组在定义时就会分配内存空间，切片在运行时才会分配


```go
// 创建数组的方式  
    var array1 [5]int
    array2 := [3]int{1,2,3}
	
//创建切片的方式
    slice1 := array[0:2]
    var slice2 []int
    slice3 := make([]int,1)
    slice4 := []int{1,2,3}
    slice5 := append(slice4,4)
```

## Question four
1. 筛选素数
2. channel和go routine
3. 同时筛选多个素数，效率更高
