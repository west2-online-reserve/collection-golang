声明数组时需要声明数组长度
var array [20]int

切片是动态的引用类型


创建切片方式：
var s []int
make([]int, len, cap)
切取数组 arr[1:3]
var s []int{1, 2}



创建MAP：

var m map[int]string

mp:=make(map[keyType]valueType)