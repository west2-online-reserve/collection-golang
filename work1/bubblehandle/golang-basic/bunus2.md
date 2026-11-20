# Go语言中的切片和数组的区别
切片的容量是动态可变的 更灵活 

# 创建方式

- 数组的创建方式
```aiignore
- var arr [length]type  
- arr:=[...]type{element,element,element}
```
- 切片的创建方式  

```aiignore
var slice []type 
slice := []type{element}
slice:=make([]int,3)
slice:=make([]int,3,5)
slice:=[]int //错误写法
```

- map的创建方式

```aiignore
- var map1 map[int]string
map1:=make(map[int]string,9)
map1:=make(map[int]string]
map1:=map[int]string{
    map1[1]="go"
    ...
}

```




