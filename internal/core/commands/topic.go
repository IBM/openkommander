package commands

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/IBM/sarama"
)

func CreateTopic(topicName string, numPartitions, replicationFactor int) (successMessage string, f *Failure) {
	client, validateFailure := GetClient()
	if validateFailure != nil {
		return "", validateFailure
	}

	adminClient, validateFailure := GetAdminClient()
	if validateFailure != nil {
		return "", validateFailure
	}

	if numPartitions < 1 || replicationFactor < 1 {
		return "", NewFailure("Partitions and replication factor must be greater than 1", http.StatusBadRequest)
	}

	if replicationFactor > len(client.Brokers()) {
		return "", NewFailure("Replication factor cannot be greater than the number of brokers", http.StatusBadRequest)
	}

	topicDetail := &sarama.TopicDetail{
		NumPartitions:     int32(numPartitions),
		ReplicationFactor: int16(replicationFactor),
	}

	err := adminClient.CreateTopic(topicName, topicDetail, false)
	if err != nil {
		if strings.Contains(err.Error(), "Topic with this name already exists") {
			return "", NewFailure(fmt.Sprintf("Topic '%s' already exists", topicName), http.StatusInternalServerError)
		} else {
			return "", NewFailure(fmt.Sprintf("Error creating topic: %v", topicName), http.StatusInternalServerError)
		}
	}

	return fmt.Sprintf("Successfully created topic '%s' with %d partitions and replication factor %d",
		topicName, numPartitions, replicationFactor), nil
}

// When successful, returns a success message
func DeleteTopic(topicName string) (successMessage string, f *Failure) {
	client, validateFailure := GetAdminClient()
	if validateFailure != nil {
		return "", validateFailure
	}

	if topicName == "" {
		return "", NewFailure("Topic name cannot be empty", http.StatusInternalServerError)
	}

	err := client.DeleteTopic(topicName)
	if err != nil {
		return "", NewFailure(fmt.Sprintf("Error deleting topic: %v", err), http.StatusInternalServerError)
	}

	return fmt.Sprintf("Successfully deleted topic '%s'", topicName), nil
}

// When successful, returns a map of topic names to their details
func ListTopics() (topicMap map[string]sarama.TopicDetail, f *Failure) {
	client, validateFailure := GetAdminClient()
	if validateFailure != nil {
		return nil, validateFailure
	}

	topics, err := client.ListTopics()
	if err != nil {
		return nil, NewFailure(fmt.Sprintf("Error listing topics: %v", err), http.StatusInternalServerError)
	}

	if len(topics) == 0 {
		return nil, NewFailure("No topics found", http.StatusInternalServerError)
	}

	return topics, nil
}
