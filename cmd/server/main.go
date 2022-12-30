package main

import (
	"log"

	"github.com/ksusonic/go-devops-mon/internal/server"
	"github.com/ksusonic/go-devops-mon/internal/storage"
)

func main() {
	memStorage := storage.NewMemStorage()
	s := server.NewServer(memStorage)
	log.Fatal(s.Start())
}
