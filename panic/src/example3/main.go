// example3
package main

import (
	"fmt"
	"runtime"
)

// Created: Fri Sep  6 17:06:55 2019

func main() {
	f1()
	f2()
}

// f1
func f1() {
	// STARTf1 OMIT
	defer func() {
		if p := recover(); p != nil {
			fmt.Println("Panic - recovered:", p)
		}
	}()
	// ENDf1 OMIT
	panic("whoops")
}

// f2
func f2() {
	// STARTf2 OMIT
	defer func() {
		if p := recover(); p != nil {
			if _, ok := p.(runtime.Error); ok {
				fmt.Println("runtime error - recovered:", p)
			}
		}
	}()
	// ENDf2 OMIT
	zero := 0
	x := 1 / zero
	println(x)
}
