package main

func main(){

	var nums []int
	nums =[]int{2,7,11,15}
	target:=9
	for i:=0;i<len(nums);i++{
		for j:=i+1;j<len(nums);j++{
			if(nums[i]+nums[j]==target){
				fmt.Printf("%d %d\n",i,j)
			}
		}
	}
}

//有复杂度O(n)的算法--利用哈希表(map)