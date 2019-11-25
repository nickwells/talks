// example4
package main

import "fmt"

// Created: Sat Nov 23 16:27:39 2019

func main() {
	_ = NewOrderedValOrPanic(1, 2)
}

// START-STD OMIT
// OrderedVal records v1 and v2
type OrderedVal struct {
	v1 int
	v2 int
}

// NewOrderedVal returns an OrderedVal and nil if the arguments are valid or
// nil and an error if they are not
func NewOrderedVal(v1, v2 int) (*OrderedVal, error) {
	if v1 > v2 {
		return nil,
			fmt.Errorf("Bad parameters: v1 (%d) must be <= v2 (%d)",
				v1, v2)
	}
	return &OrderedVal{v1: v1, v2: v2}, nil
}

//END-STD OMIT

//START-PANIC OMIT
// NewOrderedValOrPanic will create and return an OrderedVal. It will panic
// if the parameters are invalid
func NewOrderedValOrPanic(dim1, dim2 int) *OrderedVal {
	s, err := NewOrderedVal(dim1, dim2)
	if err != nil {
		panic(err)
	}
	return s
}

//END-PANIC OMIT
