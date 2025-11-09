package main

import "fmt"

func pass1001() {
	var a int
	var b int
	_, err := fmt.Scanf("%d %d", &a, &b)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(a + b)
}
