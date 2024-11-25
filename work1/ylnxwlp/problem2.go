package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	applesStr, _ := reader.ReadString('\n')
	heightStr, _ := reader.ReadString('\n')
	applesStr = strings.TrimSpace(applesStr)
	heightStr = strings.TrimSpace(heightStr)
	applesStrs := strings.Split(applesStr, " ")
	apples := make([]int, len(applesStrs))
	for i, s := range applesStrs {
		apple, err := strconv.Atoi(s)
		if err != nil {
			fmt.Println("转换错误:", err)
			return
		}
		apples[i] = apple
	}
	height, _ := strconv.Atoi(heightStr)
	var sum = 0
	for i := range apples {
		if apples[i] <= height+30 {
			sum += 1
		}
	}
	fmt.Println(sum)
}
