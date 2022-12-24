package server

import (
	"github.com/ksusonic/go-devops-mon/internal/storage"
	"net/http"
)

type Server struct {
	MemStorage storage.Storage
}

func (s Server) registerHandlers() {
	http.HandleFunc(UpdateHandlerName, s.UpdateMetric)
}

func NewServer() Server {
	s := Server{MemStorage: &storage.MemStorage{
		GaugeStorage:   storage.GaugeStorage{},
		CounterStorage: storage.CounterStorage{},
	}}
	s.registerHandlers()
	return s
}

func (s Server) Start() error {
	return http.ListenAndServe("127.0.0.1:8080", nil)
}
