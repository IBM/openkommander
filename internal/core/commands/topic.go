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

func DescribeTopic(topicName string) (*sarama.TopicMetadata, *Failure) {
	client, validateFailure := GetAdminClient()
	if validateFailure != nil {
		return nil, validateFailure
	}

	metadata, err := client.DescribeTopics([]string{topicName})
	if err != nil || len(metadata) == 0 {
		return nil, NewFailure(fmt.Sprintf("Error describing topic '%s': %v", topicName, err), http.StatusInternalServerError)
	}

	return metadata[0], nil
}

func DescribeTopicConfig(topicName string) ([]sarama.ConfigEntry, *Failure) {
	client, validateFailure := GetAdminClient()
	if validateFailure != nil {
		return nil, validateFailure
	}

	configs, err := client.DescribeConfig(sarama.ConfigResource{Type: sarama.TopicResource, Name: topicName})
	if err != nil {
		return nil, NewFailure(fmt.Sprintf("Error describing configs for topic '%s': %v", topicName, err), http.StatusInternalServerError)
	}

	return configs, nil
}

func UpdateTopic(topicName string, newPartitions int) (successMessage string, f *Failure) {
	client, validateFailure := GetAdminClient()
	if validateFailure != nil {
		return "", validateFailure
	}

	if topicName == "" {
		return "", NewFailure("Topic name cannot be empty", http.StatusInternalServerError)
	}

	topicMetadata, err := client.DescribeTopics([]string{topicName})
	if err != nil || len(topicMetadata) == 0 {
		return "", NewFailure(fmt.Sprintf("Error describing topic '%s': %v", topicName, err), http.StatusInternalServerError)
	}

	topic := topicMetadata[0]
	existingPartitions := len(topic.Partitions)
	if newPartitions <= existingPartitions {
		return "", NewFailure("New partition count must be greater than the existing partitions", http.StatusBadRequest)
	}

	err = client.CreatePartitions(topicName, int32(newPartitions), nil, false)
	if err != nil {
		return "", NewFailure(fmt.Sprintf("Error updating partitions for topic '%s': %v", topicName, err), http.StatusInternalServerError)
	}

	return fmt.Sprintf("Successfully updated topic '%s' to %d partitions.", topicName, newPartitions), nil
}
