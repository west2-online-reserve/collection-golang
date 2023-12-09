package main

import (
    "fmt"
    "math"
)

func isPrime(n int) bool {
    if n <= 1 {
        return false
    }
    if n == 2 {
        return true
    }
    if n%2 == 0 {
        return false
    }

    sqrtN := int(math.Sqrt(float64(n)))
    for i := 3; i <= sqrtN; i += 2 {
        if n%i == 0 {
            return false
        }
    }
    return true
}

func main() {
    var n int
    fmt.Scan(&n)

    if isPrime(n) {
        fmt.Println("YES")
    } else {
        fmt.Println("NO")
    }
}
