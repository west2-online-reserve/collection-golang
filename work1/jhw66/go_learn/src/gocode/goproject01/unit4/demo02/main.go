package main

import "fmt"

func main() {
	var sum int

	//不能用var i int在for上定义i
	for i := 1; i <= 3; i++ {
		sum += i
	}

	fmt.Println("1+2+3=", sum)

	i := 1
	for i <= 3 {
		sum += i
		i++
	}
	fmt.Println("1+2+3=", sum)

	i = 0
	for {
		fmt.Println("死循环")
		i++
		if i == 3 {
			break
		}
	}

	var str = "hello golang"
	for j := 0; j < len(str); j++ {
		fmt.Println(string(str[j])) //byte类型
		//or fmt.Printf("%c\n", str[j])
	}

	//for range遍历字符串
	for index, value := range str {
		if value == ' ' {
			fmt.Printf("index=%d value=空格\n", index)
		} else {
			fmt.Printf("index=%d value=%c\n", index, value) //rune类型
		}
	}
}
