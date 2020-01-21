// example6
package main

import (
	"fmt"
	"os"
)

// Created: Fri Jan 17 18:09:20 2020

// START OMIT
func main() {
	err := exPanicAndRecover()
	if err != nil {
		fmt.Println("Something went wrong:", err)
		os.Exit(1)
	}
	fmt.Println("All's well.")
}

// exPanicAndRecover should return an error if anything goes wrong
func exPanicAndRecover() error {
	defer func() {
		if p := recover(); p != nil {
			fmt.Println("Panic!")
		}
	}()
	panic("Whoops!")
}

// END OMIT
