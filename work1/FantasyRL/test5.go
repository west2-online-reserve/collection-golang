package main

import "fmt"

func main() {
	var homo, Acceed []int
	for i := 0; i < 50; i++ {
		homo = append(homo, i+1)
	}
	for i := 0; i < 50; i++ {
		if homo[i]%3 != 0 {
			Acceed = append(Acceed, homo[i])
		}
	}
	Acceed = append(Acceed, 114514)
	fmt.Println(Acceed)
}
