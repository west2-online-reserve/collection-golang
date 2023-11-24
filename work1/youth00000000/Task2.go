package main

import "fmt"

func main() {
	var a [10]int // 定义一个长度为10的数组
	var height int
	var sum int
	var err error
	//写入数组内的十个元素，即十个苹果高度
	for i := 0; i < 10; i++ {
		_, err := fmt.Scan(&a[i])
		if err != nil {
			fmt.Println(err)
		}
	}
	//写入人的高度
	_, err = fmt.Scan(&height)
	if err == nil { //即上行函数不报错的情况下
		for i := 0; i < len(a); i++ {
			if a[i] <= height { //比较人和苹果的高度
				sum++
			} else {
				m := height + 30 //加上板凳的高度
				if a[i] <= m {
					sum++
				}
			}
		}
		fmt.Println(sum)
	}

}
