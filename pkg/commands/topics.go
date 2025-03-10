package commands

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/IBM/openkommander/pkg/session"
	"github.com/IBM/sarama"
	"github.com/jedib0t/go-pretty/v6/table"
)

func createTopicCommand() {
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

	fmt.Print("Enter topic name: ")
	var topicName string
	_, err = fmt.Scanln(&topicName)
	if err != nil {
		fmt.Println("Error reading input:", err)
	}

	if topicName == "" {
		fmt.Println("Topic name cannot be empty")
		return
	}

	fmt.Print("Enter number of partitions (default 1): ")
	var partitionsStr string
	_, err = fmt.Scanln(&partitionsStr)
	if err != nil {
		fmt.Println("Error reading input:", err)
	}
	partitions := 1
	if partitionsStr != "" {
		if p, err := strconv.Atoi(partitionsStr); err == nil && p > 0 {
			partitions = p
		}
	}

	fmt.Print("Enter replication factor (default 1): ")
	var replicationStr string
	_, err = fmt.Scanln(&replicationStr)
	if err != nil {
		fmt.Println("Error reading input:", err)
	}
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
		topicName, partitions, replicationFactor)
}

func deleteTopicCommand() {
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

	fmt.Print("Enter topic name to delete: ")
	var topicName string
	_, err = fmt.Scanln(&topicName)
	if err != nil {
		fmt.Println("Error reading input:", err)
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

func listTopicsCommand() {
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

func describeTopicCommand() {
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

	fmt.Print("Enter topic name: ")
	var topicName string
	_, err = fmt.Scanln(&topicName)
	if err != nil {
		fmt.Printf("Error reading input: %v\n", err)
		return
	}

	if topicName == "" {
		fmt.Println("Topic name cannot be empty")
		return
	}

	// Fetch topic metadata
	topicMetadata, err := client.DescribeTopics([]string{topicName})
	if err != nil || len(topicMetadata) == 0 {
		fmt.Printf("Error describing topic '%s': %v\n", topicName, err)
		return
	}

	topic := topicMetadata[0]

	fmt.Println("\nTopic Metadata:")
	fmt.Printf("  Topic Name: %s\n", topic.Name)
	fmt.Printf("  Replication Factor: %d\n", len(topic.Partitions[0].Replicas))
	fmt.Printf("  Version: %d\n", topic.Version)
	fmt.Printf("  UUID: %s\n", topic.Uuid)
	fmt.Printf("  Is Internal: %t\n", topic.IsInternal)
	fmt.Printf("  Authorized Operations: %d\n", topic.TopicAuthorizedOperations)

	// Table for partition details
	fmt.Println("\nTopic Partitions:")
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Partition ID", "Leader", "Replicas", "In-Sync Replicas (ISR)"})

	for _, partition := range topic.Partitions {
		t.AppendRow(table.Row{
			partition.ID,
			partition.Leader,
			fmt.Sprintf("%v", partition.Replicas),
			fmt.Sprintf("%v", partition.Isr),
		})
	}

	t.SetStyle(table.StyleLight)
	t.Render()

	// Fetch and display topic configurations
	configs, err := client.DescribeConfig(sarama.ConfigResource{Type: sarama.TopicResource, Name: topicName})
	if err != nil {
		fmt.Printf("Error describing configs for topic: %v\n", err)
		return
	}

	fmt.Println("\nTopic Configurations:")
	configTable := table.NewWriter()
	configTable.SetOutputMirror(os.Stdout)
	configTable.AppendHeader(table.Row{"Config Name", "Value"})

	for _, config := range configs {
		configTable.AppendRow(table.Row{config.Name, config.Value})
	}

	configTable.SetStyle(table.StyleLight)
	configTable.Render()

}
