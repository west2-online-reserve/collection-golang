package main

import "fmt"

func main() {
	var count int = 30
	if count > 20 {
		fmt.Println("count大于20")
	} else { //else必须和if的}在同一行
		fmt.Println("count小于等于20")
	}

	//在if语句中声明变量
	if amount := 100; amount > 50 {
		fmt.Println("amount大于50")
	}

	var score int = 65
	//var score2 int64=5   //类型不匹配
	var score3 int = 6
	//多分支选择
	switch score / 10 {
	case 10, 9:
		fmt.Println("优秀")
	case 8:
		fmt.Println("良好")
	case 7:
		fmt.Println("中等")
	case score3:
		fmt.Println("及格")
		fallthrough //穿透一层
	// case score2:
	// 	fmt.Println("差")
	default:
		fmt.Println("差！")
	}

	//switch后不带表达式，类似if...else...结构
	var a int = 5
	switch {
	case a > 0:
		fmt.Println("a是正数")
	case a < 0:
		fmt.Println("a是负数")
	default:
		fmt.Println("a是0")
	}

	//在switch语句中声明变量
	switch b := 2; b {
	case 1:
		fmt.Println("b=1")
	case 2:
		fmt.Println("b=2")
	case 3:
		fmt.Println("b=3")
	default:
		fmt.Println("b不在1-3之间")
	}
}
