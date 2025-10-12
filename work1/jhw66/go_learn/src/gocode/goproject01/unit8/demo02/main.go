package main

import "fmt"

func main() {
	b := map[int]string{
		999: "aaa",
		888: "sss",
	}
	//增加/修改操作
	b[777] = "kkk"
	//删除操作
	delete(b, 999)
	fmt.Println(b)
	//清空操作
	//1，遍历key逐个删除
	//2，重新make是原来那个成为垃圾被回收
	b = make(map[int]string, 10) //这时候就不是:=了
	fmt.Println(b)
	//查找操作
	value, flag := b[9]
	fmt.Println(value, flag)
	b[777] = "kkk"
	b[787] = "kkk"
	b[797] = "kkk"

	//获取长度
	fmt.Println(len(b))

	//遍历
	for k, v := range b {
		fmt.Print(k, v, "\t")
	}
	fmt.Println()

	a := make(map[string]map[int]string)
	a["班级一"] = make(map[int]string, 3)
	a["班级一"][2009] = "wjh"
	a["班级一"][2008] = "jhw"
	a["班级二"] = make(map[int]string, 3)
	a["班级二"][2009] = "jjj"

	for k1, v1 := range a {
		fmt.Println(k1)
		for k2, v2 := range v1 {
			fmt.Println(k2, v2)
		}
		fmt.Println()
	}
}
