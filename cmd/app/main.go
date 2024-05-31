package main

import (
	"encoding/json"
	"fmt"
	"github.com/Uikola/knative-golang/pkg/zlog"
	"github.com/sirikothe/gotextfsm"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/rs/zerolog/log"
)

const (
	DEBUGLEVEL = 0
)

type Interface struct {
	MTU     int `json:"mtu"`
	RxPkt   int `json:"rx_pkt"`
	RxBytes int `json:"rx_bytes"`
	TxPkt   int `json:"tx_pkt"`
	TxBytes int `json:"tx_bytes"`
}

func NewInterface(mtuStr, rxPktStr, rxBytesStr, txPktStr, txBytesStr string) (Interface, error) {
	mtu, err := strconv.Atoi(mtuStr)
	if err != nil {
		return Interface{}, err
	}
	rxPkt, err := strconv.Atoi(rxPktStr)
	if err != nil {
		return Interface{}, err
	}
	rxBytes, err := strconv.Atoi(rxBytesStr)
	if err != nil {
		return Interface{}, err
	}
	txPkt, err := strconv.Atoi(txPktStr)
	if err != nil {
		return Interface{}, err
	}
	txBytes, err := strconv.Atoi(txBytesStr)
	if err != nil {
		return Interface{}, err
	}

	return Interface{
		MTU:     mtu,
		RxPkt:   rxPkt,
		RxBytes: rxBytes,
		TxPkt:   txPkt,
		TxBytes: txBytes,
	}, nil
}

type NetworkInfo struct {
	Interfaces map[string]Interface `json:"interfaces"`
}

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

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(networkInfo)
	} else {
		log.Info().Msg(event.String())
	}
}

func parseIfConfigOutput(output string) (NetworkInfo, error) {
	networkInfo := NetworkInfo{
		Interfaces: make(map[string]Interface),
	}

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
		return NetworkInfo{}, err
	}
	parser := gotextfsm.ParserOutput{}
	err = parser.ParseTextString(output, fsm, true)
	if err != nil {
		fmt.Printf("Error while parsing input '%s'\n", err.Error())
	}

	for _, record := range parser.Dict {
		inface, err := NewInterface(
			record["MTU"].(string),
			record["RX_PACKETS"].(string),
			record["RX_BYTES"].(string),
			record["TX_PACKETS"].(string),
			record["TX_BYTES"].(string),
		)
		if err != nil {
			return NetworkInfo{}, err
		}

		name := record["INTERFACE"].(string)
		networkInfo.Interfaces[name] = inface
	}

	return networkInfo, nil
}
