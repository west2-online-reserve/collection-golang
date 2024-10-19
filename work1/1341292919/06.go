package main

import (
	"fmt"
	"os"
)

func main() {
	filename := "ninenine.txt"
	file, err := os.Create(filename) //创建并打开文件
	if err != nil {
		fmt.Println("create file failed", err)
		return
	}
	defer file.Close()

	for i := 1; i <= 9; i++ {
		for j := i; j <= 9; j++ {
			temp := fmt.Sprintf("%d*%d=%-*d ", i, j, 2, i*j) //-左对齐，*为宽度
			_, err = file.WriteString(temp)                  //将字符串写入对应文件
			if err != nil {                                  //处理错误
				fmt.Println("Writing fail", err)
				return
			}
		}
		temp := "\n"
		_, err = file.WriteString(temp)
		if err != nil { //处理错误
			fmt.Println("Writing fail", err)
			return
		}
	}

}
