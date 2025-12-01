package main

import (
	"errors"
	"fmt"
	"math"
)

func arc017a() {
	var num int
	_, err := fmt.Scanf("%d", &num)
	if err != nil {
		fmt.Println(err)
	}
	prime, err := isPrime(num)
	if err != nil {
		fmt.Println(err)
		return
	}
	if prime {
		fmt.Println("YES")
	} else {
		fmt.Println("NO")
	}
}
func isPrime(n int) (bool, error) {
	if n < 2 {
		return false, errors.New("num is invalid")
	}
	sqrt := math.Sqrt(float64(n))
	for i := 2; i <= int(sqrt); i++ {
		if n%i == 0 {
			return false, nil
		}
	}
	return true, nil
}
