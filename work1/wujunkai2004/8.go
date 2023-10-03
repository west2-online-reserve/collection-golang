package main

import "fmt"

func search(nums []int, sum int) []int {
	num_map := make( map[int]int );

	for index, value := range nums {
		want := sum - value;
		gets, ok := num_map[ want ];

		if(ok){
			return []int {gets, index};
		}

		num_map[ value ] = index;

	}

	return []int{};
}

func main(){
	fmt.Print( search( []int{2, 7, 11, 15}, 9 ) )
	fmt.Print( search( []int{2, 7, 11, 15}, 20 ) )
	fmt.Print( search( []int{3,2,4}, 6 ))
}