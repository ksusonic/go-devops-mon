package main

import (
	"github.com/ksusonic/go-devops-mon/internal/server"
	"log"
)

func main() {
	s := server.NewServer()
	log.Println("Server started")
	log.Fatal(s.Start())
}
