package main

import (
	"errors"
	"fmt"
)

// START OMIT
func main() {
	f()
}

func f() {
	defer fmt.Println("deferred func in f()")
	_, _ = panicker()
}

func panicker() (int, error) {
	defer fmt.Println("deferred func first")
	defer fmt.Println("deferred func last")
	fmt.Println("About to panic")

	panic(errors.New("Whoops!"))

	fmt.Println("Just panicked") // nolint: unreachable
	defer fmt.Println("deferred func 3")
	return 99, nil
}

// END OMIT
