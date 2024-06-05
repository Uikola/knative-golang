package entity

type Service struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	ClusterIP  string `json:"cluster_ip"`
	ExternalIP string `json:"external_ip"`
	Ports      string `json:"ports"`
	Age        string `json:"age"`
}
