package main

import (
	"fmt"
	"strconv"
)

func main() {
	var num1 int = 10
	fmt.Println(num1)

	var num2 float32 = float32(num1) //类型转换必须显示转换
	fmt.Println(num2)

	var num3 int64 = 88888
	var num4 int8 = int8(num3) //大类型转换为小类型，可能会溢出
	fmt.Println(num4)

	var num5 int64 = 12
	var num6 int8 = int8(num5) + 127 //编译会通过，但是会溢出
	//var num7 int8 = int8(num5) + 128 //编译不会通过
	fmt.Println(num6)

	//其他类型转换为字符串类型
	var n1 int = 19
	var n2 float32 = 4.78
	var n3 bool = false
	var n4 byte = 'a'
	//使用fmt.Sprintf转换
	var s1 string = fmt.Sprintf("%d", n1)
	fmt.Printf("s1=%v\n", s1)
	var s2 string = fmt.Sprintf("%f", n2)
	fmt.Printf("s2=%q\n", s2)
	var s3 string = fmt.Sprintf("%t", n3)
	fmt.Printf("s3=%v\n", s3)
	var s4 string = fmt.Sprintf("%c", n4)
	fmt.Printf("s4=%q\n", s4)
	//使用strconv包转换
	var str1 string = strconv.FormatInt(int64(n1), 10)
	//10表示十进制
	var str2 string = strconv.FormatFloat(float64(n2), 'f', -1, 32)
	//-1表示默认的精度(需要特定精度要自己设置)，32表示float32
	var str3 string = strconv.FormatBool(n3)
	//返回字符串"true"或"false"
	var str4 string = strconv.FormatUint(uint64(n4), 10)
	//10表示十进制
	fmt.Printf("str1=%q\n", str1)
	fmt.Printf("str2=%q\n", str2)
	fmt.Printf("str3=%q\n", str3)
	fmt.Printf("str4=%q\n", str4)

	//字符串类型转换为其他类型
	var str5 string = "true"
	var str6 string = "12345"
	var str7 string = "3.14"
	var str8 string = "a"
	//转换函数都返回两个值，第一个是转换后的值，第二个是错误信息
	//如果不需要错误信息，可以使用_忽略
	//如果转换失败，第一个值是该类型的零值，第二个值是错误信息
	b, err := strconv.ParseBool(str5)
	if err != nil {
		fmt.Println("转换错误")
	}
	i, _ := strconv.ParseInt(str6, 10, 64)
	f, _ := strconv.ParseFloat(str7, 64)
	c, err := strconv.ParseUint(str8, 10, 8)
	if err != nil {
		fmt.Println("转换错误")
	}
	fmt.Printf("b=%v i=%v f=%v c=%v\n", b, i, f, c)
}
