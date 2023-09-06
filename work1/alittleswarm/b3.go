package main

func main() {
	num := []int{2, 7, 11, 15}
	println(findTarget(9, num))
}
func findTarget(target int, num []int) (ok bool, l int, r int) {
	ok = false
	for i := 0; i < len(num)-1; i++ {
		for j := i + 1; j < len(num); j++ {
			if num[i]+num[j] == target {
				l = i
				r = j
				ok = true
			}
		}
	}
	return
}

/*可以用哈希表（GPT说的
 */
