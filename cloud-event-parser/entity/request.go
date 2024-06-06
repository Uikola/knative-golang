package entity

import "strings"

const (
	IFCONFIGCMD = "ifconfig"
)

type CloudEventRequest struct {
	Data struct {
		Stdin  string `json:"stdin"`
		Stdout string `json:"stdout"`
	} `json:"data"`
}

func (r CloudEventRequest) StdinCmdIsConfig() bool {
	return strings.Contains(r.Data.Stdin, IFCONFIGCMD)
}
