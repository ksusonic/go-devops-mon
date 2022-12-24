package main

import (
	"github.com/ksusonic/go-devops-mon/internal/server"
	"log"
)

func main() {
	s := server.NewServer()
	log.Fatal(s.Start())
}
