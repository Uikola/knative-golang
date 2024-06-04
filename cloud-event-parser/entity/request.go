package entity

import "strings"

const (
	IFCONFIGCMD = "ifconfig"
)

type CloudEventDataRequest struct {
	Stdin  string `json:"stdin"`
	Stdout string `json:"stdout"`
}

func (r CloudEventDataRequest) StdinCmdIsConfig() bool {
	return strings.Contains(r.Stdin, IFCONFIGCMD)
}
