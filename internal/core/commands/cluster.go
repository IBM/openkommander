package commands

import (
	"net/http"

	"github.com/IBM/sarama"
)

// ClusterInfo represents information about a Kafka cluster/broker
type ClusterInfo struct {
	ID        int32  `json:"id"`
	Address   string `json:"address"`
	Status    string `json:"status"`
	Rack      string `json:"rack"`
	Connected bool   `json:"connected"`
}

// ListClusters returns information about all available clusters/brokers
// When successful, returns a slice of ClusterInfo structs
func ListClusters() (clusters []ClusterInfo, f *Failure) {
	client, validateFailure := GetClient()
	if validateFailure != nil {
		return nil, validateFailure
	}

	brokers := client.Brokers()
	if len(brokers) == 0 {
		return nil, NewFailure("No clusters found", http.StatusNotFound)
	}

	clusters = make([]ClusterInfo, 0, len(brokers))

	for _, broker := range brokers {
		status := "Disconnected"
		connected := false

		// Check if broker is connected
		isConnected, err := broker.Connected()
		if err == nil && isConnected {
			status = "Connected"
			connected = true
		} else if err := broker.Open(client.Config()); err == nil || err == sarama.ErrAlreadyConnected {
			status = "Connected"
			connected = true
		}

		rack := broker.Rack()
		if rack == "" {
			rack = "N/A"
		}

		clusterInfo := ClusterInfo{
			ID:        broker.ID(),
			Address:   broker.Addr(),
			Status:    status,
			Rack:      rack,
			Connected: connected,
		}

		clusters = append(clusters, clusterInfo)
	}

	return clusters, nil
}

// GetClusterMetadata returns metadata about the current cluster
// When successful, returns cluster metadata information
func GetClusterMetadata() (metadata map[string]interface{}, f *Failure) {
	client, validateFailure := GetClient()
	if validateFailure != nil {
		return nil, validateFailure
	}

	brokers := client.Brokers()

	metadata = make(map[string]interface{})
	metadata["broker_count"] = len(brokers)
	metadata["cluster_id"] = client.Config().ClientID

	// Get broker details
	brokerDetails := make([]map[string]interface{}, 0, len(brokers))
	for _, broker := range brokers {
		connected, _ := broker.Connected()
		brokerInfo := map[string]interface{}{
			"id":        broker.ID(),
			"address":   broker.Addr(),
			"connected": connected,
			"rack":      broker.Rack(),
		}
		brokerDetails = append(brokerDetails, brokerInfo)
	}
	metadata["brokers"] = brokerDetails

	return metadata, nil
}
