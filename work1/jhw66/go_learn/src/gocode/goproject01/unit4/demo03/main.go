package main

import "fmt"

func main() {
	var str = "hello"
label1:
	for i := 0; i < 4; i++ {
		//label2:
		for index, value := range str {
			// if index == 3 { //与下面的效果相同
			// 	break label2
			// }
			if index == 3 {
				continue label1
			}
			fmt.Printf("index=%v value=%c\n", index, value)
		}
	}

	fmt.Println("hello golang1")
	fmt.Println("hello golang2")
	if true {
		goto label3
	}
	fmt.Println("hello golang3")
	fmt.Println("hello golang4")
	fmt.Println("hello golang5")
	fmt.Println("hello golang6")
label3:
	fmt.Println("hello golang7")
	fmt.Println("hello golang8")
}
