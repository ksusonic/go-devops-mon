package server

import (
	"log"
	"net/http"

	"github.com/ksusonic/go-devops-mon/internal/metrics"
	"github.com/ksusonic/go-devops-mon/internal/server/router"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	Storage metrics.MetricStorage
	Router  *chi.Mux
}

func NewServer(storage metrics.MetricStorage) *Server {
	s := &Server{
		Storage: storage,
	}
	s.Router = router.NewRouter(&s.Storage)
	return s
}

func (s Server) Start() error {
	log.Println("Server started")
	return http.ListenAndServe("127.0.0.1:8080", s.Router)
}
