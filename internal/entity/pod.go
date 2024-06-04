package entity

type Pod struct {
	Name     string `json:"name"`
	Ready    string `json:"ready"`
	Status   string `json:"status"`
	Restarts int    `json:"restarts"`
	Age      string `json:"age"`
}
