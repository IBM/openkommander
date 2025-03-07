package functions

import (
	"fmt"
	"strings"

	"github.com/IBM/openkommander/pkg/session"
	"github.com/IBM/sarama"
)

func CreateTopic(topicName string, numPartitions, replicationFactor int) {
	currentSession := session.GetCurrentSession()
	if !currentSession.IsAuthenticated() {
		fmt.Println("Error: no session found.")
		return
	}

	client, err := currentSession.GetAdminClient()
	if err != nil {
		fmt.Printf("Error connecting to cluster: %v\n", err)
		return
	}

	topicDetail := &sarama.TopicDetail{
		NumPartitions:     int32(numPartitions),
		ReplicationFactor: int16(replicationFactor),
	}

	err = client.CreateTopic(topicName, topicDetail, false)
	if err != nil {
		if strings.Contains(err.Error(), "Topic with this name already exists") {
			fmt.Printf("Topic '%s' already exists\n", topicName)
		} else {
			fmt.Printf("Error creating topic: %v\n", err)
		}
		return
	}

	fmt.Printf("Successfully created topic '%s' with %d partitions and replication factor %d\n",
		topicName, numPartitions, replicationFactor)
}

func DeleteTopic(topicName string) {
	currentSession := session.GetCurrentSession()
	if !currentSession.IsAuthenticated() {
		fmt.Println("Error: no session found.")
		return
	}

	client, err := currentSession.GetAdminClient()
	if err != nil {
		fmt.Printf("Error connecting to cluster: %v\n", err)
		return
	}

	if topicName == "" {
		fmt.Println("Topic name cannot be empty")
		return
	}

	err = client.DeleteTopic(topicName)
	if err != nil {
		fmt.Printf("Error deleting topic: %v\n", err)
		return
	}

	fmt.Printf("Successfully deleted topic '%s'\n", topicName)
}

func ListTopics() {
	currentSession := session.GetCurrentSession()
	if !currentSession.IsAuthenticated() {
		fmt.Println("Error: no session found.")
		return
	}

	client, err := currentSession.GetAdminClient()
	if err != nil {
		fmt.Printf("Error connecting to cluster: %v\n", err)
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
