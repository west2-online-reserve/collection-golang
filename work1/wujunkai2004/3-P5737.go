package main

import "fmt"


func is_leap_year(year int) bool {
	return ( year % 4==0 && year % 100 != 0 || year % 400 == 0 );
}

func main()  {
	var start int = 0;
	var end   int = 0;

	var ans   int = 0;

	fmt.Scan(&start, &end);

	for year:=start; year<=end; year++ {
		if( is_leap_year(year) ){
			ans++;
		}
	}

	fmt.Printf("%d\n", ans);

	for year:=start; year<=end; year++ {
		if( is_leap_year(year) ){
			fmt.Printf("%d ", year);
		}
	}
}