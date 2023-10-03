package main
import "fmt"
func main(){
	var x,y int
	count := 0
	var year[]int
	fmt.Scanf("%d %d" ,&x, &y)
	for i := x ;i <= y; i++{
		if (i % 4 == 0 && i % 100 != 0) || (i % 100 ==0 && i % 400 == 0){
			year = append(year,i)
			count++
		}
	}
	fmt.Println(count)
	for i := 0 ; i < len(year);i++{
		fmt.Printf("%d ",year[i])
	}
}