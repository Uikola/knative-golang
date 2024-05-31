package entity

import "strconv"

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
