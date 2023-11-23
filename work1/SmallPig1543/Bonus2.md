
### Question1:
    -1.数组是静态的，切片是动态的
    -2.数组是聚合类型，切片是引用类型
    -3.切片的底层仍是数组

### Question2:
    -slice：
        -有5种，分别为：
            -1.name := make([]T,length,capacity)
            -2.name := make([]T,length)
            -3.var name []T{}
            -4.var name []T{val1,va2,...}
            -5.name := arr[n:m]
    -map:
        -有2种，分别为：
            -1.name := make(map[T1]T2){}
            -2.var name map[T1]T2
