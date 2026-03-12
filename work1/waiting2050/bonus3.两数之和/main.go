func twoSum(nums []int, target int) []int {
    
    n := len(nums)
    mp := make(map[int]int, n + 1)

    for i := 0; i < n; i++ {

        nd := target - nums[i]
        
        if v, ok := mp[nd]; ok{
            return []int{i, v}
        }
        mp[nums[i]] = i;
    }

    return []int{}
}