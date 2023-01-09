package server

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	Router       *chi.Mux
	serveAddress string
}

func NewServer(config *Config, router *chi.Mux) *Server {
	return &Server{
		Router:       router,
		serveAddress: config.ServeAddress,
	}
}

func (s Server) Start() error {
	for _, route := range s.Router.Routes() {
		for _, r := range route.SubRoutes.Routes() {
			log.Println("Registered path:", r.Pattern)

		}
	}
	log.Println("Server started")
	return http.ListenAndServe(s.serveAddress, s.Router)
}
