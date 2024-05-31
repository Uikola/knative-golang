package cloud_event

import (
	"encoding/json"
	"fmt"
	"github.com/Uikola/knative-golang/internal/entity"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/rs/zerolog/log"
	"github.com/sirikothe/gotextfsm"
	"io"
	"net/http"
	"strings"
)

func ParseEvent(w http.ResponseWriter, r *http.Request) {
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

	var data map[string]string
	err = json.Unmarshal(event.Data(), &data)
	if err != nil {
		log.Error().Err(err).Msg("failed to unmarshal json into CloudEvent")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"reason": "invalid CloudEvent format"})
		return
	}

	if strings.Contains(data["stdin"], "ifconfig") {
		networkInfo, err := parseIfConfigOutput(data["stdout"])
		if err != nil {
			log.Error().Err(err).Msg("failed to parse ifconfig output")
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"reason": "bad request"})
			return
		}

		clEvent := cloudevents.NewEvent()
		clEvent.SetSource("example/uri")
		clEvent.SetType("example.type")
		err = clEvent.SetData(cloudevents.ApplicationJSON, networkInfo)
		if err != nil {
			log.Error().Err(err).Msg("failed to set cloud event data")
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"reason": "internal error"})
			return
		}

		clEventJSON, err := json.Marshal(event)
		if err != nil {
			log.Error().Err(err).Msg("failed to marshal cloud event")
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"reason": "internal error"})
			return
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(clEventJSON)
	} else {
		log.Info().Msg(event.String())
	}
}

func parseIfConfigOutput(output string) (map[string]entity.Interface, error) {
	networkInfo := make(map[string]entity.Interface)

	template := `Value INTERFACE (\S+)
Value MTU (\d+)
Value RX_PACKETS (\d+)
Value RX_BYTES (\d+)
Value TX_PACKETS (\d+)
Value TX_BYTES (\d+)

Start
  ^${INTERFACE}: flags=\S+  mtu ${MTU}
  ^\s+RX packets ${RX_PACKETS}  bytes ${RX_BYTES} \(\d+\.?\d* \w+\)
  ^\s+TX packets ${TX_PACKETS}  bytes ${TX_BYTES} \(\d+\.?\d* \w+\) -> Record
`
	fsm := gotextfsm.TextFSM{}
	err := fsm.ParseString(template)
	if err != nil {
		fmt.Printf("Error while parsing template '%s'\n", err.Error())
		return nil, err
	}
	parser := gotextfsm.ParserOutput{}
	err = parser.ParseTextString(output, fsm, true)
	if err != nil {
		fmt.Printf("Error while parsing input '%s'\n", err.Error())
	}

	for _, record := range parser.Dict {
		inface, err := entity.NewInterface(
			record["MTU"].(string),
			record["RX_PACKETS"].(string),
			record["RX_BYTES"].(string),
			record["TX_PACKETS"].(string),
			record["TX_BYTES"].(string),
		)
		if err != nil {
			return nil, err
		}

		name := record["INTERFACE"].(string)
		networkInfo[name] = inface
	}

	return networkInfo, nil
}
