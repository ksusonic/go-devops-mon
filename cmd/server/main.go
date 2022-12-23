package main

import (
	"github.com/ksusonic/go-devops-mon/internal/handlers"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc(handlers.UpdateHandlerName, handlers.UpdateMetric)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
