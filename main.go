package main

import (
	"encoding/json"
	"github.com/Uikola/knative-golang/pkg/zlog"
	"io"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/rs/zerolog/log"
)

const (
	DEBUGLEVEL = 0
)

func main() {
	log.Logger = zlog.Default(true, "dev", DEBUGLEVEL)

	r := chi.NewRouter()
	r.Post("/cloud-event", CloudEventHandler)

	log.Info().Msg("starting server...")
	if err := http.ListenAndServe(":8000", r); err != nil {
		log.Error().Err(err).Msg(err.Error())
		os.Exit(1)
	}
}

func CloudEventHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error().Err(err).Msg("failed to read request body")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"reason": "bad request"})
		return
	}

	event := cloudevents.NewEvent()

	err = json.Unmarshal(body, &event)
	if err != nil {
		log.Error().Err(err).Msg("failed to unmarshal json into CloudEvent")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"reason": "invalid CloudEvent format"})
		return
	}

	log.Info().Msg(event.String())
}
