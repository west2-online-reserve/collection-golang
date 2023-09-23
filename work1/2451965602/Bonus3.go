package main
import "fmt"
func main(){
	var nums []int{3,2,4}
	var target int
	fmt.Scan(&target)
	for x := 0; x < len(nums);x++{
		for y := (x + 1); y < len(nums); y++{
			if nums[x] + nums[y] == target{
				fmt.Printf("[%d %d]",x,y)
			}
		}
	}
}