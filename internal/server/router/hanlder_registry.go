package router

import (
	"log"
	"net/http"

	"github.com/ksusonic/go-devops-mon/internal/metrics"

	"github.com/go-chi/chi/v5"
)

type serverHandler struct {
	method  string
	path    string
	handler func(w http.ResponseWriter, r *http.Request, c context)
}

type context struct {
	storage *metrics.MetricStorage
}

var availableHandlers []serverHandler

func registerHandler(method string, path string, handler func(w http.ResponseWriter, r *http.Request, c context)) {
	availableHandlers = append(availableHandlers, serverHandler{
		method:  method,
		path:    path,
		handler: handler,
	})
}

func propagateData(handler func(w http.ResponseWriter, r *http.Request, c context), inputContext context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, inputContext)
	}
}

func NewRouter(storage *metrics.MetricStorage) *chi.Mux {
	r := chi.NewRouter()
	c := context{storage: storage}
	for _, h := range availableHandlers {
		log.Printf("Registered %s handler %s\n", h.method, h.path)
		r.Method(h.method, h.path, propagateData(h.handler, c))
	}
	return r
}
