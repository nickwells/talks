// multipleGoRoutines
//
// This demonstrates that deferred functions are only called in the panicking
// goroutine
package main

import (
	"fmt"
	"os"
	"time"
)

// Created: Sat Jan 11 11:17:12 2020

func main() {
	go sleeper()
	go sleeper()
	panicker()
}

// sleeper defers a function and sleeps for 10 seconds
func sleeper() {
	defer func() {
		fmt.Fprintf(os.Stderr, "sleeper: defer func called\n")
	}()

	time.Sleep(10 * time.Second)
}

// panicker defers a function and then panics
func panicker() {
	defer func() {
		fmt.Fprintf(os.Stderr, "panicker: defer func called\n")
	}()
	panic("Whoops!")

}
