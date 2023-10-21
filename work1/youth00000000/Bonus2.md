# Q1 : 切片和数组的区别

1. 切片是指针类型，数组是值类型。
2. 数组的赋值形式为值传递，切片的赋值形式为引用传递。
3. 数组的长度是固定的，而切片长度可以任意调整。
4. 切片比数组多一个容量属性。

# Q2 : 创建切片的方法

1. 使用内置的`make`函数创建切片：

	``` go
	slice := make([]type, length, capacity)
	```

2. 使用字面量创建切片

	```go
	slice := []type{value1, value2, ...}
	```

3. 通过数组创建切片：

	```go
	slice := &array[start:end]
	```

	

# Q3 : 创建map的方法

1. 使用`make`函数创建`map`：

	```go
	Map := make(map[keyType]valueType)
	```

2. 使用字面量创建`map`：

	```go
	myMap := map[keyType]valueType{
	    key1: value1,
	    key2: value2,
	    // ...
	}
	```

	