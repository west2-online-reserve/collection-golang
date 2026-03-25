Q1:
数组是固定长度的,且长度是类型的一部分。数组的长度必须是常量表达式，因为数组的长度需要在编译阶段确定。
切片是动态的,可以动态收缩和扩增.切片不支持==比较内容,只能和nil用==比较


Q2:
创建切片的方式:
1. 直接创建nil, 然后append() 
    var Slicing []Type
    Slicing = append(Slicing, element)
    
2. 使用make() 
    Slicing := make([]Type, len, cap)

3. 从数组或其他切片截取 
    arr := []int{1, 2, 3, 4, 5}
    Slicing := arr[1:5] // [2, 3, 4, 5]


Q3:
创建map:
1. 依旧map函数
    randomMap := make(map[KType]VType, cap)
    
2. 直接创建
    randomMap := map[KType]VType{键值对...} 
    // 直接创建一个nil map后续写入数据时会panic
    // 因此创建一个空map建议使用make()