package entity

import "strings"

const (
	IFCONFIGCMD       = "ifconfig"
	KUBECTLGETNSCMD   = "kubectl get ns"
	KUBECTLGETPODSCMD = "kubectl get pods"
	KUBECTLGETSVCCMD  = "kubectl get svc"
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
	return strings.Contains(c.Stdin, KUBECTLGETPODSCMD)
}

func (c CloudEventData) StdinCmdKubectlGetSvc() bool {
	return strings.Contains(c.Stdin, KUBECTLGETSVCCMD)
}
