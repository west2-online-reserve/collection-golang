package main

import "fmt"

func main() {
	var arr = [10]int{0}
	var hand_get_height int
	ans := 0
	for i := 0; i < 10; i++ {
		var temp int
		fmt.Scan(&temp)
		arr[i] = temp
	}
	chair := 30
	fmt.Scan(&hand_get_height)
	for i := 0; i < len(arr); i++ {
		if arr[i] <= chair+hand_get_height {
			ans++
		}
	}
	fmt.Println(ans)

}
