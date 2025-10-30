package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func pass1046() {
	scanner := bufio.NewScanner(os.Stdin)
	var line1, line2 string
	if scanner.Scan() {
		line1 = scanner.Text()
	}
	if scanner.Scan() {
		line2 = scanner.Text()
	}
	split := strings.Split(line1, " ")
	if len(split) != 10 {
		fmt.Println("bad input")
		return
	}
	atoi, err := strconv.Atoi(line2)
	if err != nil {
		fmt.Println("bad input")
		return
	}
	var count int8
	atoi += 30
	for _, v := range split {
		i, err := strconv.Atoi(v)
		if err != nil {
			fmt.Println("bad input")
			return
		}
		if i < 100 || i > 200 {
			fmt.Println("bad input")
			return
		}
		if atoi >= i {
			count++
		}
	}
	fmt.Println(count)
}
