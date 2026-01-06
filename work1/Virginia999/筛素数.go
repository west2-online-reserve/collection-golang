package main

import "fmt"

func main() {
	var a int
	fmt.Scanf("%d", &a)
	//引入布尔型变量，防止循环遍历中多次输出结果

	result := false
	for i := 0; i < a; i++ {
		for j := 0; j < a; j++ {
			if a == i*j { //若布尔型变量保持不变，不用特别说明
			} else if a != i*j {
				//说明布尔型变量的转化
				result = true
			}
		}
	}
	//把需要打印的值放在循环外
	if result {
		fmt.Println("YES")
	} else {
		fmt.Println("NO")
	}
}

//对于这道题，请编写一个判断质数的函数isPrime(x int) bool ，并且在主函数中调用它
