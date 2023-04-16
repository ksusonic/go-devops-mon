package main

import (
	"math/rand"
	"os"
)

func main() {
	if rand.Intn(100) > 50 {
		os.Exit(0) // want "os.Exit call prohibited in main package, main function"
	} else {
		os.Exit(100) // want "os.Exit call prohibited in main package, main function"
	}
}
