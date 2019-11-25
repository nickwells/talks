// panic
package main

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/nickwells/param.mod/v3/param"
	"github.com/nickwells/param.mod/v3/param/paramset"
	"github.com/nickwells/param.mod/v3/param/psetter"
)

// Created: Mon Aug 26 10:50:05 2019

const (
	sf_none    = "none"
	sf_unnamed = "unnamedReturnVals"
	sf_named   = "namedReturnVals"
)

var funcMap = map[string]func() (int, string, error){
	sf_unnamed: unnamedReturnVals,
	sf_named:   namedReturnVals,
}

const (
	rt_none   = "none"
	rt_basic  = "basic"
	rt_seterr = "seterr"
)

const indent = "   "

var callPanic = true
var panicInF2Defer bool
var recoverType = "none"
var panicStartFunc = sf_unnamed
var panicValType = "string"

var panicVal interface{}

type panicStruct struct {
	i int
	s string
}

// (ps panicStruct)String returns a string representing the panicStruct
func (ps panicStruct) String() string {
	return fmt.Sprintf("panicStruct{i: %d, s: %q}", ps.i, ps.s)
}

func main() {
	ps := paramset.NewOrDie(addParams,
		param.SetProgramDescription("panic in various ways"))
	ps.Parse()

	// this is purely to have another user-created goroutine so we can
	// explore the different behaviour seen with alterative settings of the
	// GOTRACEBACK environment variable
	go func() {
		time.Sleep(10 * time.Second)
	}()

	f, ok := funcMap[panicStartFunc]
	if ok {
		i, s, err := f()
		fmt.Printf("%s returned: i: %d   s: %q   err: %v\n",
			panicStartFunc, i, s, err)
	}
}

// addParams will add parameters to the passed ParamSet
func addParams(ps *param.PSet) error {
	ps.Add("starting-func",
		psetter.Enum{
			AVM: param.AVM{
				AllowedVals: param.AValMap{
					sf_none: "do nothing - no function will be called",
					sf_unnamed: "panic and (optionally) recover" +
						" starting from a function with unnamed return values",
					sf_named: "panic and (optionally) recover" +
						" starting from a function with named return values",
				},
			},
			Value: &panicStartFunc,
		},
		"which type of function do you want to start with. This allows you"+
			" to see the differences in behaviour when functions"+
			" have named return values. Note that the functions take no"+
			" parameters and return an int, a string and an error",
		param.AltName("f"))

	ps.Add("recover-type",
		psetter.Enum{
			AVM: param.AVM{
				AllowedVals: param.AValMap{
					rt_none: "do not recover",
					rt_basic: "recover, simply reporting that" +
						" you have recovered",
					rt_seterr: "recover and, if a panic was" +
						" detected, set an error value." +
						" If you choose this then the starting" +
						" function will be set to " + sf_named,
				},
			},
			Value: &recoverType,
		},
		"how do you want to recover",
		param.AltName("r"))

	ps.Add("panic-val",
		psetter.Enum{
			AVM: param.AVM{
				AllowedVals: param.AValMap{
					"nil":    "set the panic value to nil",
					"int":    "set the panic value to an int (42)",
					"string": "set the panic value to a string ('a string')",
					"err":    "set the panic value to an error ('an error')",
					"struct": "set the panic value to a struct" +
						" (an int set to 42 and a string set to 'a string')",
				},
			},
			Value: &panicValType,
		},
		"what do you want the panic value set to",
		param.AltName("v"))

	ps.Add("no-panic",
		psetter.Bool{
			Value:  &callPanic,
			Invert: true,
		},
		"do everything but don't panic (to see the effect of panic)")

	ps.Add("panic-in-defer", psetter.Bool{Value: &panicInF2Defer},
		"if we panic then also panic in the defered function on the way back up")

	ps.AddFinalCheck(func() error {
		if recoverType == rt_seterr {
			if panicStartFunc != sf_none &&
				panicStartFunc != sf_named {
				return fmt.Errorf(
					"the starting function is set to %s"+
						" but a recovery type of %s only works with"+
						" a starting function of %s",
					panicStartFunc, recoverType, sf_named)
			}
			panicStartFunc = sf_named
		}
		return nil
	})
	ps.AddFinalCheck(func() error {
		switch panicValType {
		case "string":
			panicVal = "a string"
		case "int":
			panicVal = 42
		case "struct":
			panicVal = panicStruct{i: 42, s: "a string"}
		case "nil":
			panicVal = nil
		case "err":
			panicVal = errors.New("an error")
		default:
			return fmt.Errorf("unknown panic value type: %q",
				panicValType)
		}
		return nil
	})

	return nil
}

// recoverFunc will recover and print the panic val
func recoverFunc(name string) {
	fmt.Println(indent+name+"(defered)", "- about to attempt to recover")
	if pv := recover(); pv != nil {
		fmt.Println(indent+name+"(defered)", "- recovered")
		fmt.Println("\n*** panic value:", pv)
		fmt.Println()
	} else {
		fmt.Println(indent+name+"(defered)", "- recover returned nil")
	}
}

// unnamedReturnVals will optionally recover from any passed panics. If there
// is no panic it will return 99 and a string: "desired value". Otherwise it
// returns nil values
func unnamedReturnVals() (int, string, error) {
	depth := 0
	name := "unnamedReturnVals"
	fmt.Println(name, "- entered")

	defer func() { fmt.Println(indent+name+"(defered)", "- first defered func") }()
	if recoverType == rt_basic {
		defer recoverFunc(name)
	}
	defer func() { fmt.Println(indent+name+"(defered)", "- last defered func") }()

	f1(depth + 1)

	fmt.Println(name, "- left")
	return 9, "no panic", nil
}

// namedReturnVals will optionally recover from any panic. If there is no
// panic it will return 99 and a string: "desired value". Otherwise it
// returns 42 and "pre-panic value" (the values set before panic)
func namedReturnVals() (i int, s string, err error) {
	depth := 0
	name := "namedReturnVals"
	fmt.Println(name, "- entered")

	defer func() { fmt.Println(indent+name+"(defered)", "- first defered func") }()
	switch recoverType {
	case rt_basic:
		defer recoverFunc(name)
	case rt_seterr:
		defer func() {
			fmt.Println(name, "- about to attempt to recover")
			if pv := recover(); pv != nil {
				fmt.Println(name, "- recovered")
				fmt.Println("\n*** panic value:", pv)
				fmt.Println()
				fmt.Println(name, "- setting error")
				err = errors.New("panic recovered")
			} else {
				fmt.Println(name, "- recover returned nil")
			}
		}()
	}
	defer func() { fmt.Println(indent+name+"(defered)", "- last defered func") }()

	f1(depth + 1)

	fmt.Println(name, "- left")
	return 9, "no panic", nil
}

// f1 just calls f2 - it's here to demonstrate the call stack you see
// with inlined functions
func f1(depth int) {
	f2(depth + 1)
}

// f2 calls f3 but has defered functions and some logging to help you see
// what's happening
func f2(depth int) {
	name := strings.Repeat(indent, depth) + "f2"
	fmt.Println(name, "- entered")
	defer func() { fmt.Println(indent+name+"(defered)", "- first defered func") }()
	defer func() {
		if callPanic && panicInF2Defer {
			fmt.Println(indent+name+"(defered)", "- panicking again")
			panic("2nd panic")
		}
	}()
	defer func() { fmt.Println(indent+name+"(defered)", "- last defered func") }()
	f3(depth + 1)
	fmt.Println(name, "- left")
}

// f3 just calls panicker - it's here to demonstrate the call stack you see
// with inlined functions
func f3(depth int) {
	panicker(depth + 1)
}

// panicker will call panic passing the panicVal. It logs what's happening to
// help you understand the behaviour
func panicker(depth int) {
	name := strings.Repeat(indent, depth) + "panicker"
	fmt.Println(name, "- entered")
	defer func() {
		fmt.Println(name, "- defered func")
	}()
	if callPanic {
		fmt.Println(name, "- about to panic")
		panic(panicVal)
	}
	fmt.Println(name, "- no panic")
	fmt.Println(name, "- left")

}
