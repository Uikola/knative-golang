package cloud_event

import (
	"encoding/json"
	"github.com/Uikola/knative-golang/internal/entity"
	"github.com/Uikola/knative-golang/internal/parser"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/rs/zerolog/log"
	"net/http"
)

func ParseEvent(w http.ResponseWriter, r *http.Request) {
	event := cloudevents.NewEvent()
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		log.Error().Err(err).Msg("failed to decode cloud event")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"reason": "invalid cloud event formant"})
		return
	}

	var data entity.CloudEventData
	if err := event.DataAs(&data); err != nil {
		log.Error().Err(err).Msg("failed to unmarshal cloud event data")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"reason": "invalid cloud event format"})
		return
	}

	if data.StdinCmdIsConfig() {
		networkInfo, err := parser.ParseIfConfigOutput(data.Stdout)
		if err != nil {
			log.Error().Err(err).Msg("failed to parse ifconfig output")
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"reason": "bad request"})
			return
		}

		responseEvent, err := parser.ConvertNetworkInfoToCloudEvent(networkInfo)
		if err != nil {
			log.Error().Err(err).Msg("failed to convert network info to cloud event")
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"reason": "internal error"})
			return
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(responseEvent)
	} else if data.StdinCmdKubectlGetNs() {
		namespaceInfo, err := parser.ParseKubectlGetNsOutput(data.Stdout)
		if err != nil {
			log.Error().Err(err).Msg("failed to parse kubectlgetns output")
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"reason": "bad request"})
			return
		}

		responseEvent, err := parser.ConvertNamespaceInfoToCloudEvent(namespaceInfo)
		if err != nil {
			log.Error().Err(err).Msg("failed to convert namespace info to cloud event")
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"reason": "internal error"})
			return
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(responseEvent)
	}
}
