package fib

import "fmt"

func fib(n int){
	i1, i2 := 0, 1

	for i := 0; i < n; i++ {
		fmt.Println(i1)
		i1, i2 = i2, i1+i2
	}
}
