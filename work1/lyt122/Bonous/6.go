package Bonous

import (
	"fmt"
	"os"
)

func Bonous1() {
	file, err := os.Create("./Bonous/ninenine.txt")
	if err != nil {
		fmt.Println("无法创建文件:", err)
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)
	for i := 1; i <= 9; i++ {
		for j := 1; j <= i; j++ {
			result := i * j
			_, err := fmt.Fprintf(file, "%d * %d = %d\t", j, i, result)
			if err != nil {
				return
			}
		}
		_, err := fmt.Fprintln(file, "")
		if err != nil {
			return
		}
	}
}
