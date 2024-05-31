package http

import (
	"github.com/Uikola/knative-golang/internal/server/http/cloud_event"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func NewServer() http.Handler {
	router := chi.NewRouter()

	addRoutes(router)

	var handler http.Handler = router

	return handler
}

func addRoutes(router *chi.Mux) {
	router.Post("/cloud-event", cloud_event.ParseEvent)
}
