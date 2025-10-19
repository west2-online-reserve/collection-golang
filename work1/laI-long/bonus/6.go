package main

import (
	"fmt"
	"os"
)

func main() {
	file, err := os.OpenFile("./work1/lai-long/bonus/ninenine.txt", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	for i := 1; i < 10; i++ {
		for j := 1; j <= i; j++ {
			str := fmt.Sprintf("%d * %d = %d\t", j, i, i*j)
			_, err := file.WriteString(str)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		_, err = file.WriteString("\n")
	}
	defer file.Close()
}
