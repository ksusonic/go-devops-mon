package server

import (
	"net/http"

	"github.com/ksusonic/go-devops-mon/internal/storage"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	MemStorage storage.Storage
	Router     chi.Router
}

func (s Server) registerHandlers() {
	s.Router.Get("/", s.GetAllMetrics)
	s.Router.Post("/update/{type}/{name}/{value}", s.UpdateMetric)
	s.Router.Get("/value/{type}/{name}", s.GetMetric)
}

func NewServer() Server {
	s := Server{
		MemStorage: &storage.MemStorage{
			GaugeStorage:   storage.GaugeStorage{},
			CounterStorage: storage.CounterStorage{},
		},
		Router: chi.NewRouter(),
	}
	s.registerHandlers()
	return s
}

func (s Server) Start() error {
	return http.ListenAndServe("127.0.0.1:8080", s.Router)
}
