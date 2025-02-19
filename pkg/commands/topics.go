package commands

import (
	"fmt"
	"strconv"
	"strings"

	"openkommander/pkg/session"

	"github.com/IBM/sarama"
)

func init() {
	Register("Topics Management", "topic-create", createTopicCommand)
	Register("Topics Management", "topic-delete", deleteTopicCommand)
	Register("Topics Management", "topic-list", listTopicsCommand)
}

func createTopicCommand() {
	currentSession := session.GetCurrentSession()
	if !currentSession.IsAuthenticated() {
		fmt.Println("Error: no session found.")
		return
	}

	client := currentSession.GetAdminClient()
	if client == nil {
		fmt.Println("Error: not connected to a cluster.")
		return
	}

	fmt.Print("Enter topic name: ")
	var topicName string
	fmt.Scanln(&topicName)

	if topicName == "" {
		fmt.Println("Topic name cannot be empty")
		return
	}

	fmt.Print("Enter number of partitions (default 1): ")
	var partitionsStr string
	fmt.Scanln(&partitionsStr)
	partitions := 1
	if partitionsStr != "" {
		if p, err := strconv.Atoi(partitionsStr); err == nil && p > 0 {
			partitions = p
		}
	}

	fmt.Print("Enter replication factor (default 1): ")
	var replicationStr string
	fmt.Scanln(&replicationStr)
	replicationFactor := 1
	if replicationStr != "" {
		if rf, err := strconv.Atoi(replicationStr); err == nil && rf > 0 {
			replicationFactor = rf
		}
	}

	topicDetail := &sarama.TopicDetail{
		NumPartitions:     int32(partitions),
		ReplicationFactor: int16(replicationFactor),
	}

	err := client.CreateTopic(topicName, topicDetail, false)
	if err != nil {
		if strings.Contains(err.Error(), "Topic with this name already exists") {
			fmt.Printf("Topic '%s' already exists\n", topicName)
		} else {
			fmt.Printf("Error creating topic: %v\n", err)
		}
		return
	}

	fmt.Printf("Successfully created topic '%s' with %d partitions and replication factor %d\n",
		topicName, partitions, replicationFactor)
}

func deleteTopicCommand() {
	currentSession := session.GetCurrentSession()
	if !currentSession.IsAuthenticated() {
		fmt.Println("Error: no session found.")
		return
	}

	client := currentSession.GetAdminClient()
	if client == nil {
		fmt.Println("Error: not connected to a cluster.")
		return
	}

	fmt.Print("Enter topic name to delete: ")
	var topicName string
	fmt.Scanln(&topicName)

	if topicName == "" {
		fmt.Println("Topic name cannot be empty")
		return
	}

	err := client.DeleteTopic(topicName)
	if err != nil {
		fmt.Printf("Error deleting topic: %v\n", err)
		return
	}

	fmt.Printf("Successfully deleted topic '%s'\n", topicName)
}

func listTopicsCommand() {
	currentSession := session.GetCurrentSession()
	if !currentSession.IsAuthenticated() {
		fmt.Println("Error: no session found.")
		return
	}

	client := currentSession.GetAdminClient()
	if client == nil {
		fmt.Println("Error: not connected to a cluster.")
		return
	}

	topics, err := client.ListTopics()
	if err != nil {
		fmt.Printf("Error listing topics: %v\n", err)
		return
	}

	if len(topics) == 0 {
		fmt.Println("No topics found")
		return
	}

	fmt.Println("\nTopics:")
	fmt.Println("--------")
	for name, detail := range topics {
		fmt.Printf("Name: %s\n", name)
		fmt.Printf("  Partitions: %d\n", detail.NumPartitions)
		fmt.Printf("  Replication Factor: %d\n", detail.ReplicationFactor)
		fmt.Println()
	}
}
