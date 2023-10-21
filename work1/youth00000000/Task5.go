package main

import "fmt"

func main() {
	var a = make([]int, 0, 50) //创建容量为50，长度为零的空切片
	//利用append+for循环按序添加50个自然数
	for i := 1; i < 51; i++ {
		a = append(a, i)
	}

	//在切片a中删除3的倍数
	var k = 0
	for _, v := range a { //遍历a切片
		//筛出不是3的倍数的数并替换原下标为k的数字
		if v%3 != 0 { //筛出不是3的倍数的数字
			a[k] = v //替换原有的数字
			k++
		}
	}
	//该代码实际上是找出1-50中的所有不是3的倍数的数将其替换到数组序号靠前的数字
	a = a[:k] //通过将切片a再切片 删除后16个未被替换的原数字

	//利用append添加114514
	a = append(a, 114514)
	fmt.Println(a)
}
