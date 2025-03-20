package models

type ClusterInfo struct {
	Name    string   `json:"name"`
	Brokers []string `json:"brokers"`
	Status  string   `json:"status"`
}


