1、数组长度固定，切片长度可扩张
2、在函数内对数组进行修改，在使用函数后数组仍不变（不使用指针的话）；在函数内修改切片则使用函数后切片会发生变化（切片是引用类型）
########################

1、s:=make([]Type,size,cap)
2、a:=[4]int{1,2,3,4}
   s:=a[:3]
3、s:=[]type{}
4、var s []type
########################

1、var mapname map[keyType]valueType
2、mp:=make(map[keyType]valueType)
