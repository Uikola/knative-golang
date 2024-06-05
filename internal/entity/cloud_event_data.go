package entity

import "strings"

const (
	IFCONFIGCMD     = "ifconfig"
	KUBECTLGETNSCMD = "kubectl get ns"
	KUBECTLGETPODS  = "kubectl get pods"
)

type CloudEventData struct {
	Stdin  string `json:"stdin"`
	Stdout string `json:"stdout"`
}

func (c CloudEventData) StdinCmdIsConfig() bool {
	return strings.Contains(c.Stdin, IFCONFIGCMD)
}

func (c CloudEventData) StdinCmdKubectlGetNs() bool {
	return strings.Contains(c.Stdin, KUBECTLGETNSCMD)
}

func (c CloudEventData) StdinCmdKubectlGetPods() bool {
	return strings.Contains(c.Stdin, KUBECTLGETPODS)
}
