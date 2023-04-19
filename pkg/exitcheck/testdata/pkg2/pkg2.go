package main

import (
	"math/rand"
	"os"
)

func main() {
	notMain()
}

func notMain() {
	if rand.Intn(100) > 50 {
		os.Exit(0)
	} else {
		os.Exit(100)
	}
}
