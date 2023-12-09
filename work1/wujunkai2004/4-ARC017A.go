package main

import "fmt"
import "math"

func is_prime(num int) bool {
	end := int( math.Sqrt( float64( num ) ) );
	for chk:=2; chk<=end; chk++ {
		if( num % chk == 0 ){
			return false;
		}
	}
	return true;
}

func main(){
	var num int = 0;
	fmt.Scan( &num );

	if( is_prime( num ) ){
		fmt.Printf("YES\n");
	} else {
		fmt.Printf("NO\n");
	}
}