package main

import (
	"fmt"
	"strconv"
	"strings"
)

func main() {
	str := "golang你好" //一个汉字三个字节
	fmt.Println((len(str)))

	for i, value := range str {
		fmt.Println(i, string(value))
	}
	r := []rune(str)
	for i := 0; i < len(r); i++ {
		fmt.Printf("%c\n", r[i])
	}

	num1, _ := strconv.Atoi("666")
	fmt.Println(num1)
	str1 := strconv.Itoa(88)
	fmt.Println(str1)

	//统计有多少个字串
	count := strings.Count("golangandhavaga", "ga")
	fmt.Println(count)

	//不区分大小写的字符串比较
	flag := strings.EqualFold("hello", "HELLO")
	fmt.Println(flag)

	//区分大小写的字符串比较
	fmt.Println("hello" == "HELLO")

	//返回子串第一次出现的索引值，没有返回-1
	index := strings.Index("golangandjava", "an")
	fmt.Println(index)

	//字符串的替换(n=-1表示全部替换，n=2表示替换两个)
	str2 := strings.Replace("golangjavagogo", "go", "golang", -1)
	str3 := strings.Replace("golangjavagogo", "go", "golang", 2)
	fmt.Println(str2)
	fmt.Println(str3)

	//按照指定的某个字符，为分割标记，将一个字符串拆分成字符串数组
	arr := strings.Split("go-java-c++", "-")
	fmt.Println(arr)

	fmt.Println(strings.ToLower("Go"))
	fmt.Println(strings.ToUpper("go"))

	//去除字符串左右两边空格
	fmt.Println(strings.TrimSpace("   hello go   "))

	//将字符串左/右/左和右两边指定的字符去掉
	fmt.Println(strings.Trim("~golang~", "~"))
	fmt.Println(strings.TrimLeft("~golang~", "~"))
	fmt.Println(strings.TrimRight("~golang~", "~"))

	//判断字符串是否以指定的字符串开头/结尾
	fmt.Println(strings.HasPrefix("http://java.sum.com", "http"))
	fmt.Println(strings.HasSuffix("demo.png", ".jpg"))
}
