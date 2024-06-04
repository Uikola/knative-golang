package entity

import (
	"strconv"
)

type Interface struct {
	Name    string `json:"name"`
	MTU     int    `json:"mtu"`
	RxPkt   int    `json:"rx_pkt"`
	RxBytes int    `json:"rx_bts"`
	TxPkt   int    `json:"tx_pkt"`
	TxBytes int    `json:"tx_bts"`
}

func NewInterface(name, mtuStr, rxPktStr, rxBytesStr, txPktStr, txBytesStr string) Interface {
	mtu, _ := strconv.Atoi(mtuStr)
	rxPkt, _ := strconv.Atoi(rxPktStr)
	rxBytes, _ := strconv.Atoi(rxBytesStr)
	txPkt, _ := strconv.Atoi(txPktStr)
	txBytes, _ := strconv.Atoi(txBytesStr)

	return Interface{
		Name:    name,
		MTU:     mtu,
		RxPkt:   rxPkt,
		RxBytes: rxBytes,
		TxPkt:   txPkt,
		TxBytes: txBytes,
	}
}
