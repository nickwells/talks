// ex1
package main

import (
	"flag"
	"os"
)

// Created: Fri Sep  6 15:00:21 2019

// START OMIT
func main() {
	var ip1 = flag.Int("f", 1, "help")
	var ip2 = flag.Int("f", 2, "help")

	// END OMIT
	if ip1 == ip2 {
		os.Exit(1)
	}
}
