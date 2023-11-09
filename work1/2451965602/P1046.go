package main  
import "fmt"  
func main() {  
    var applehigh [10]int
	var manhigh int
	count := 0
	//输入苹果高度
	for i := 0; i < 10; i++{
		fmt.Scan(&applehigh[i])
	}
	//输入陶陶高度
	fmt.Scan(&manhigh)
	//比较大小
	for i := 0; i < 10; i++{
		if (manhigh + 30 >= applehigh[i]) {
			count++
		}
	}
	fmt.Println(count)

}