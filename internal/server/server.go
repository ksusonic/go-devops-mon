package server

import (
	"log"
	"net/http"

	"github.com/ksusonic/go-devops-mon/internal/metrics"
	"github.com/ksusonic/go-devops-mon/internal/server/controller"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	Router *chi.Mux
}

func NewServer(storage metrics.MetricStorage) *Server {
	router := chi.NewRouter()

	metricController := controller.NewController(storage)
	metricController.Register(router)

	s := &Server{
		Router: router,
	}
	return s
}

func (s Server) Start() error {
	for _, route := range s.Router.Routes() {
		log.Println("Registered route:", route.Pattern)
	}
	log.Println("Server started")
	return http.ListenAndServe("127.0.0.1:8080", s.Router)
}
