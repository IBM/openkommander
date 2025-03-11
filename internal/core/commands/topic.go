package commands

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/IBM/sarama"
)

func CreateTopic(topicName string, numPartitions, replicationFactor int) (*Success[string], *Failure) {
	client, validateFailure := GetClient()
	if validateFailure != nil {
		return nil, validateFailure
	}

	adminClient, validateFailure := GetAdminClient()
	if validateFailure != nil {
		return nil, validateFailure
	}

	// Function parameter validation
	if numPartitions < 1 || replicationFactor < 1 {
		return nil, NewFailure("Partitions and replication factor must be greater than 1", http.StatusBadRequest)
	}

	if replicationFactor > len(client.Brokers()) {
		return nil, NewFailure("Replication factor cannot be greater than the number of brokers", http.StatusBadRequest)
	}

	topicDetail := &sarama.TopicDetail{
		NumPartitions:     int32(numPartitions),
		ReplicationFactor: int16(replicationFactor),
	}

	err := adminClient.CreateTopic(topicName, topicDetail, false)
	if err != nil {
		if strings.Contains(err.Error(), "Topic with this name already exists") {
			return nil, NewFailure(fmt.Sprintf("Topic '%s' already exists", topicName), http.StatusInternalServerError)
		} else {
			return nil, NewFailure(fmt.Sprintf("Error creating topic: %v", topicName), http.StatusInternalServerError)
		}
	}

	// Could return a Success[Topic] object if we want to abstract more
	return NewSuccess(fmt.Sprintf("Successfully created topic '%s' with %d partitions and replication factor %d",
		topicName, numPartitions, replicationFactor)), nil
}

func DeleteTopic(topicName string) (*Success[string], *Failure) {
	client, validateFailure := GetAdminClient()
	if validateFailure != nil {
		return nil, validateFailure
	}

	if topicName == "" {
		return nil, NewFailure("Topic name cannot be empty", http.StatusInternalServerError)
	}

	err := client.DeleteTopic(topicName)
	if err != nil {
		return nil, NewFailure(fmt.Sprintf("Error deleting topic: %v", err), http.StatusInternalServerError)
	}

	return NewSuccess(fmt.Sprintf("Successfully deleted topic '%s", topicName)), nil
}

func ListTopics() (*Success[map[string]sarama.TopicDetail], *Failure) {
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

	return NewSuccess(topics), nil
}
