//1数组不可变长，切片可动态变化长度。//
//2函数传值时，数组传临时变量，切片传指针//
//3切片只能一维//
//数组在内存中是一段连续的内存空间，切片在内存中由一个指向底层数组的指针、长度和容量构成//
//-----------------------------//
//1.从数组中截取//
var sliceName = arr[idx1:idx2]
//2.用make函数//
var sliceName []type= make([]type, len, [cap])、
//3.指定数组创建//
var sliceName []type = []type{"contents"}
//-----------------------------//
//1.make函数创建//
mapName := make(map[keyType]valueType, len)
//2.初始化创建//
mapName := map[keyType]valueType{
    "KEY":"Value"
}

