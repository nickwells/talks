// example4
package main

import "errors"

// Created: Sat Nov 23 16:27:39 2019

func main() {
	_ = NewObjOrPanic(true)
}

// START-STD OMIT
type Obj struct{}

func NewObj(v bool) (*Obj, error) {
	if !v {
		return nil, errors.New("False!")
	}
	return &Obj{}, nil
}

func NewObjOrPanic(v bool) *Obj {
	if s, err := NewObj(v); err != nil {
		panic(err)
	} else {
		return s
	}
}

//END-STD OMIT
