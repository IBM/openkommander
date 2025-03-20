package clusters

import (
	"fmt"

	"openkommander/lib/kafka"
	"openkommander/models"
)

func ListClusters(clients map[string]*kafka.Client) []models.ClusterInfo {
	clusters := make([]models.ClusterInfo, 0, len(clients))
	
	for name, client := range clients {
		status := "connected"
		
		_, err := client.GetBrokerInfo()
		if err != nil {
			status = "error"
		}
		
		clusterInfo := models.ClusterInfo{
			Name:    name,
			Brokers: client.GetBrokers(),
			Status:  status,
		}
		
		clusters = append(clusters, clusterInfo)
	}
	
	return clusters
}

func GetCluster(clients map[string]*kafka.Client, name string) (map[string]interface{}, error) {
	client, exists := clients[name]
	if !exists {
		return nil, fmt.Errorf("cluster not found: %s", name)
	}
	
	brokers, err := client.GetBrokerInfo()
	if err != nil {
		return nil, fmt.Errorf("error fetching broker info: %w", err)
	}
	
	topics, err := client.GetTopicInfo()
	topicCount := 0
	if err == nil {
		topicCount = len(topics)
	}
	
	return map[string]interface{}{
		"name":          name,
		"brokers":       client.GetBrokers(),
		"status":        "connected",
		"broker_count":  len(brokers),
		"topic_count":   topicCount,
		"broker_details": brokers,
	}, nil
}