package main

// import "fmt"
// import "unsafe"
import (
	"fmt"
	"unsafe"
)

func main() {
	var num = 18
	//默认是int类型
	fmt.Printf("num的整数类型是:%T\n", num)
	//占用字节数
	fmt.Println(unsafe.Sizeof(num))

	var num2 byte = 19
	fmt.Printf("num2的整数类型是:%T\n", num2)
	fmt.Println(unsafe.Sizeof(num2))

	var num3 float32 = 3.14
	fmt.Println(num3)

	//浮点数的科学计数法
	var num4 float64 = 314e-2
	fmt.Println(num4)
	var num5 float64 = -3.14e2
	fmt.Println(num5)

	//浮点数的精度问题
	var num7 float32 = 256.0000000916
	fmt.Println(num7)
	var num8 float64 = 256.0000000916
	fmt.Println(num8)

	//浮点数默认是float64类型
	var num9 = 3.14
	fmt.Printf("num9的类型是:%T\n", num9)

	var c1 byte = 'a'
	fmt.Println(c1)
	fmt.Printf("c1的字符是:%c\n", c1)

	//使用\n表示换行
	fmt.Println("aaa\nbbb")
	//使用\t表示制表符(每八个代表一个制表符)
	fmt.Println("aaa\tbbb")
	fmt.Println("aaaaa\tbbb")
	//使用\\表示反斜杠
	fmt.Println("aaa\\bbb")
	//使用\r表示回车(回到行首)
	fmt.Println("aaa\rbbb")
	//使用\b表示退格
	fmt.Println("aaa\bbbb")

	var flag bool = true
	fmt.Println("flag=", flag)
	var flag2 bool
	fmt.Println("flag2=", flag2)

	//如果字符串没有特殊字符，字符串用双引号括起来
	//如果字符串有特殊字符，字符串用反引号括起来，反引号表示字符串原样输出
	var str2 string = "hello\nworld"
	fmt.Println("str2=", str2)
	str2 = "hello,go"
	fmt.Println("str2=", str2)
	//str2[0]='H' //字符串其中的字符的值不可变
	var str3 string = `hello\nworld`
	fmt.Println("str3=", str3)

	var str1 string = "abc" + "def" + //字符串的连接加号要放在最后
		"ghi"
	fmt.Println("str1=", str1)

}
