package main

import (
	"encoding/json"
	"fmt"
	"github.com/Uikola/knative-golang/pkg/zlog"
	"io"
	"net/http"
	"os"
	"regexp"
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
	MTU     int
	RxPkt   int
	RxBytes int
	TxPkt   int
	TxBytes int
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
	Interfaces map[string]Interface
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

		for name, values := range networkInfo.Interfaces {
			fmt.Printf("%s:\n", name)
			fmt.Printf("	mtu: %d\n", values.MTU)
			fmt.Printf("	rx_pkt: %d\n", values.RxPkt)
			fmt.Printf("	rx_bytes: %d\n", values.RxBytes)
			fmt.Printf("	tx_pkt: %d\n", values.TxPkt)
			fmt.Printf("	tx_bytes: %d\n", values.TxBytes)
		}

	} else {
		log.Info().Msg(event.String())
	}
}

func parseIfConfigOutput(output string) (NetworkInfo, error) {
	networkInfo := NetworkInfo{
		Interfaces: make(map[string]Interface),
	}

	regex := regexp.MustCompile("([A-Za-z0-9.]+): .*?mtu ([0-9]+)")
	nameAndMTU := regex.FindAllStringSubmatch(output, -1)

	regex = regexp.MustCompile("RX packets ([0-9]+)  bytes ([0-9]+)")
	rx := regex.FindAllStringSubmatch(output, -1)

	regex = regexp.MustCompile("TX packets ([0-9]+)  bytes ([0-9]+)")
	tx := regex.FindAllStringSubmatch(output, -1)

	for i := 0; i < len(nameAndMTU); i++ {
		inface, err := NewInterface(nameAndMTU[i][2], rx[i][1], rx[i][2], tx[i][1], tx[i][2])
		if err != nil {
			return NetworkInfo{}, err
		}
		networkInfo.Interfaces[nameAndMTU[i][1]] = inface
	}
	return networkInfo, nil
}
