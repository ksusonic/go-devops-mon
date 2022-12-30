package server

import (
	"log"
	"net/http"

	"github.com/ksusonic/go-devops-mon/internal/metrics"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	Storage metrics.MetricStorage
	Router  chi.Router
}

func (s Server) registerHandlers() {
	s.Router.Get("/", s.GetAllMetrics)
	s.Router.Post("/update/{type}/{name}/{value}", s.UpdateMetric)
	s.Router.Get("/value/{type}/{name}", s.GetMetric)
}

func NewServer(storage metrics.MetricStorage) *Server {
	s := &Server{
		Storage: storage,
		Router:  chi.NewRouter(),
	}
	s.registerHandlers()
	return s
}

func (s Server) Start() error {
	log.Println("Server started")
	return http.ListenAndServe("127.0.0.1:8080", s.Router)
}
