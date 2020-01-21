// exampleRecoverInMain
package main

import "fmt"

// Created: Sat Jan 18 11:12:33 2020

// START OMIT
func main() {
	defer func() {
		if p := recover(); p != nil {
			fmt.Println("Panic!")
		}
	}()

	for i := 0; i < 5; i++ {
		panicker()
	}
}

// panicker prints a message and panics
func panicker() {
	fmt.Println("Panicking")
	panic("Whoops!")
}

// END OMIT
