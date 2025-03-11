package core

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/IBM/openkommander/pkg/session"
	"github.com/IBM/sarama"
)

func CreateTopic(topicName string, numPartitions, replicationFactor int) *Status[string] {
	currentSession := session.GetCurrentSession()
	if !currentSession.IsAuthenticated() {
		return NewStatus[string]("No active session found", false, http.StatusUnauthorized)
	}

	client, err := currentSession.GetAdminClient()
	if err != nil {
		return NewStatus[string](fmt.Sprintf("Error connecting to cluster: %v\n", err), false, http.StatusInternalServerError)
	}

	// Function parameter validation
	if numPartitions < 1 || replicationFactor < 1 {
		return NewStatus[string]("Partitions and replication factor must be greater than 1", false, http.StatusBadRequest)
	}

	if replicationFactor > numPartitions {
		return NewStatus[string]("Replication factor cannot be larger than partitions", false, http.StatusBadRequest)
	}

	topicDetail := &sarama.TopicDetail{
		NumPartitions:     int32(numPartitions),
		ReplicationFactor: int16(replicationFactor),
	}

	err = client.CreateTopic(topicName, topicDetail, false)
	if err != nil {
		if strings.Contains(err.Error(), "Topic with this name already exists") {
			return NewStatus[string](fmt.Sprintf("Topic '%s' already exists\n", topicName), false, http.StatusInternalServerError)
		} else {
			return NewStatus[string](fmt.Sprintf("Error creating topic: %v\n", topicName), false, http.StatusInternalServerError)
		}
	}

	return NewStatus[string](fmt.Sprintf("Successfully created topic '%s' with %d partitions and replication factor %d\n",
		topicName, numPartitions, replicationFactor), true, 200)
}

func DeleteTopic(topicName string) *Status {
	currentSession := session.GetCurrentSession()
	if !currentSession.IsAuthenticated() {
		return status.Failure("No active session found", http.StatusUnauthorized)
	}

	client, err := currentSession.GetAdminClient()
	if err != nil {
		status.Failure(fmt.Sprintf("Error connecting to cluster: %v\n", err), http.StatusInternalServerError)
	}

	if topicName == "" {
		return status.Failure("Topic name cannot be empty", http.StatusInternalServerError)
	}

	err = client.DeleteTopic(topicName)
	if err != nil {
		return status.Failure(fmt.Sprintf("Error deleting topic: %v\n", err), http.StatusInternalServerError)
	}

	return status.Success(fmt.Sprintf("Successfully deleted topic '%s'\n", topicName))
}

func ListTopic() *Status {

}

currentSession := session.GetCurrentSession()
	if !currentSession.IsAuthenticated() {
		http.Error(w, "No active session", http.StatusUnauthorized)
		return
	}

	client, err := currentSession.GetAdminClient()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error connecting to Kafka: %v", err), http.StatusInternalServerError)
		return
	}

	topics, err := client.ListTopics()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error listing topics: %v", err), http.StatusInternalServerError)
		return
	}
