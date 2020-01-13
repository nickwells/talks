// example3
package main

import (
	"errors"
	"fmt"
	"runtime"
)

// Created: Fri Sep  6 17:06:55 2019

func main() {
	// f1()
	// f2()
	// f3()
	// f4()
	f5()
	safeCall(f99, "f99")
	safeCall(f98, "f98")
	safeCall(f97, "f97")
}

// safeCall automatically recovers from any panic
func safeCall(f func() (int, error), name string) {
	defer func() {
		if p := recover(); p != nil {
			fmt.Println("recovered from panic in ", name, ": panic val:", p)
		}
	}()
	x, err := f()
	fmt.Println(name, "- no panic, returned: ", x, err)
}

// f1
func f1() {
	// START-recover OMIT
	defer func() {
		if p := recover(); p != nil {
			fmt.Println("Panic - recovered:", p)
		}
	}()
	// END-recover OMIT

	// START-panic OMIT
	panic("Whoops!")
	// END-panic OMIT
}

// f2
func f2() {
	// START-recover-runtime OMIT
	defer func() {
		if p := recover(); p != nil {
			if _, ok := p.(runtime.Error); ok {
				fmt.Println("runtime error - recovered:", p)
			}
		}
	}()
	// END-recover-runtime OMIT
	zero := 0
	x := 1 / zero
	println(x)
}

// f3
func f3() {
	// STARTf3 OMIT
	defer func() {
		if p := recover(); p != nil {
			fmt.Println("Panic - recovered:", p) // this is never run
		}
	}()
	panic(nil)
	// ENDf3 OMIT
}

// f4
func f4() {
	defer func() {
		if p := recover(); p != nil {
			fmt.Println("Panic - recovered:", p)
		}
	}()

	// START-panic-err OMIT
	panic(errors.New("Whoops!"))
	// END-panic-err OMIT
}

// START-type-assertion OMIT
type myType string

// f5 re-panics with an interface conversion error if the panic value is not
// of the right type
func f5() {
	defer func() {
		if p := recover(); p != nil {
			myVal := p.(myType) // panics if p is not a 'myType'
			fmt.Println(myVal)
		}
	}()
	panic("Whoops!")
	panic(myType("Whoops!"))
}

// END-type-assertion OMIT

// f99
// START-panic-behaviour OMIT
func f99() (int, error) {
	defer func() {
		fmt.Println("deferred function 1")
	}()
	defer func() {
		fmt.Println("deferred function 2")
	}()

	fmt.Println("About to panic")
	panic(errors.New("Whoops!"))
	fmt.Println("Just panicked") // nolint: unreachable

	defer func() {
		fmt.Println("deferred function 3")
	}()

	return 99, nil
}

// END-panic-behaviour OMIT

type Error string

// Error converts the Error back into a string
func (e Error) Error() string {
	return string(e)
}

func f97() (x int, err error) {
	defer func() {
		if p := recover(); p != nil {
			err = p.(Error)
		}
	}()
	panic("Whoops")
}

func f98() (x int, err error) {
	defer func() {
		if p := recover(); p != nil {
			err = p.(Error)
		}
	}()
	panic(Error("Whoops"))
}
