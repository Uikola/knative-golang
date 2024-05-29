package main

import (
	"encoding/json"
	"fmt"
	"github.com/Uikola/knative-golang/pkg/zlog"
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
	MTU     int
	RxPkt   int
	RxBytes int
	TxPkt   int
	TxBytes int
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

	interfaces := strings.Split(output, "\n\n")
	for _, inface := range interfaces {
		if len(inface) == 0 {
			continue
		}

		interfaceInfo := Interface{}

		lines := strings.Split(inface, "\n")
		firstElems := strings.Split(lines[0], " ")

		interfaceName := firstElems[0][:len(firstElems[0])-1]

		for _, line := range lines {
			elems := strings.Split(line, " ")
			for i, el := range elems {
				switch el {
				case "mtu":
					mtu, err := strconv.Atoi(elems[i+1])
					if err != nil {
						return NetworkInfo{}, err
					}
					interfaceInfo.MTU = mtu
				case "RX":
					if elems[i+1] == "packets" {
						rxPkt, err := strconv.Atoi(elems[i+2])
						if err != nil {
							return NetworkInfo{}, err
						}
						rxBts, err := strconv.Atoi(elems[i+5])
						if err != nil {
							return NetworkInfo{}, err
						}
						interfaceInfo.RxPkt = rxPkt
						interfaceInfo.RxBytes = rxBts
					}
				case "TX":
					if elems[i+1] == "packets" {
						rxPkt, err := strconv.Atoi(elems[i+2])
						if err != nil {
							return NetworkInfo{}, err
						}
						rxBts, err := strconv.Atoi(elems[i+5])
						if err != nil {
							return NetworkInfo{}, err
						}
						interfaceInfo.TxPkt = rxPkt
						interfaceInfo.TxBytes = rxBts
					}
				}
			}
		}
		networkInfo.Interfaces[interfaceName] = interfaceInfo
	}
	return networkInfo, nil
}
