package main

import (
	"fmt"
	"strconv"
	"strings"
)

func main() {
	var a, b, cnt int
	var list []string
	fmt.Scanf("%d %d", &a, &b)
	for i := a; i <= b; i++ {
		if (i%4 == 0 && i%100 != 0) || i%400 == 0 {
			cnt++
			list = append(list, strconv.Itoa(i))
		}
	}

	fmt.Printf("%d\n%s", cnt, strings.Join(list, " "))
}
