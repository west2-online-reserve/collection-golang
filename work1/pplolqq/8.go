package main
import "fmt"
func add_two_nums(target int,nums []int) []int{
	mp:=map[int]int{}
	a:=0
	for i:= range nums{
		a=target-nums[i]
		_,er:=mp[a]
		if er{
			return []int{mp[a],i}
		}else{mp[nums[i]]=i}
	}
	return []int{}


}
func main(){
	fmt.Println(add_two_nums(9,[]int{2,7,11,15}))//测试案例
	fmt.Println(add_two_nums(6,[]int{3,2,4}))//测试案例
}