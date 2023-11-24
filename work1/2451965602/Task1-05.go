package main
import "fmt"
func main(){
	var arr []int
	for i := 1; i <=50; i++{
		arr = append(arr,i)
	}
	for i := 0; i < len(arr);i++{
		if arr[i] % 3 == 0{
			arr = append(arr[:i],arr[i+1:]...)
			i--
		}
	}
	arr = append(arr,114514)
	fmt.Println(arr[:])
}