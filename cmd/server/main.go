package main

import (
	"log"

	"github.com/ksusonic/go-devops-mon/internal/server"
)

func main() {
	s := server.NewServer()
	log.Println("Server started")
	log.Fatal(s.Start())
}
