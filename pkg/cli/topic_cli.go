package cli

import (
	"fmt"
	"sort"

	"github.com/IBM/openkommander/internal/core/commands"
	"github.com/spf13/cobra"
)

type TopicCommandList struct{}

func (TopicCommandList) GetParentCommand() *OkParentCmd {
	return &OkParentCmd{
		Use:   "topic <command>",
		Short: "Topic management commands",
	}
}

func (m TopicCommandList) GetCommands() []*OkCmd {
	return []*OkCmd{
		{ // Create topic
			Use:   "create [TOPIC NAME]",
			Short: "Create a new topic",
			Run:   createTopic,
			Flags: []OkFlag{
				NewOkFlag(OkFlagInt, "partitions", "p", "Specify the number of partitions of the new topic"),
				NewOkFlag(OkFlagInt, "replication-factor", "r", "Specify the replication factor of the new topic"),
			},
			Args: cobra.ExactArgs(1),
		},
		{ // Delete topic
			Use:   "delete [TOPIC NAME]",
			Short: "Delete a topic",
			Run:   deleteTopic,
			Args:  cobra.ExactArgs(1),
		},
		{ // List topics
			Use:   "list",
			Short: "List all topics",
			Run:   listTopics,
		},
		{ // Describe topic
			Use:   "describe [TOPIC NAME]",
			Short: "Describe a topic",
			Run:   describeTopic,
			Args:  cobra.ExactArgs(1),
		},
		{ // Update topic
			Use:   "update [TOPIC NAME]",
			Short: "Update an existing topic to create new partitions",
			Run:   updateTopic,
			Flags: []OkFlag{
				NewOkFlag(OkFlagInt, "new-partitions", "p", "Specify the new partition count for the topic"),
			},
			RequiredFlags: []string{"new-partitions"},
			Args:          cobra.ExactArgs(1),
		},
	}
}

func (TopicCommandList) GetSubcommands() []CommandList {
	return nil
}

// Create topic

func createTopic(cmd cobraCmd, args cobraArgs) {
	name := cmd.Flags().Arg(0)

	if name == "" {
		fmt.Println("Error: Topic name is required.")
		return
	}

	numPartitions, _ := cmd.Flags().GetInt("partitions")
	replicationFactor, _ := cmd.Flags().GetInt("replication-factor")

	if numPartitions <= 0 {
		fmt.Print("Enter number of partitions: ")
		if _, err := fmt.Scanln(&numPartitions); err != nil {
			fmt.Println("Error reading number of partitions:", err)
			return
		}
	}

	if numPartitions <= 0 {
		fmt.Println("Error: Invalid partition count.")
		return
	}

	if replicationFactor <= 0 {
		fmt.Print("Enter replication factor: ")
		if _, err := fmt.Scanln(&replicationFactor); err != nil {
			fmt.Println("Error reading replication factor:", err)
			return
		}
	}

	if replicationFactor <= 0 {
		fmt.Println("Error: Invalid replication factor.")
		return
	}

	successMessage, failure := commands.CreateTopic(name, numPartitions, replicationFactor)
	if failure != nil {
		fmt.Println(failure.Err)
		return
	}

	fmt.Println(successMessage)
}

// Delete topic

func deleteTopic(cmd cobraCmd, args cobraArgs) {
	name := cmd.Flags().Arg(0)

	if name == "" {
		fmt.Println("Error: Topic name is required.")
		return
	}

	successMessage, failure := commands.DeleteTopic(name)
	if failure != nil {
		fmt.Println(failure.Err)
		return
	}

	fmt.Println(successMessage)
}

// List topics

func listTopics(cmd cobraCmd, args cobraArgs) {
	topics, failure := commands.ListTopics()
	if failure != nil {
		fmt.Println(failure.Err)
		return
	}

	sortedTopicNames := make([]string, 0, len(topics))
	for name := range topics {
		sortedTopicNames = append(sortedTopicNames, name)
	}
	sort.Strings(sortedTopicNames)

	topicHeaders := []string{"Name", "Partitions", "Replication Factor"}
	topicRows := [][]interface{}{}
	for _, name := range sortedTopicNames {
		detail := topics[name]
		topicRows = append(topicRows, []interface{}{
			name,
			detail.NumPartitions,
			detail.ReplicationFactor,
		})
	}
	RenderTable("Topics:", topicHeaders, topicRows)
}

// Describe a topic
func describeTopic(cmd cobraCmd, args cobraArgs) {
	topicName := cmd.Flags().Arg(0)

	if topicName == "" {
		fmt.Println("Error: Topic name is required.")
		return
	}

	metadata, failure := commands.DescribeTopic(topicName)
	if failure != nil {
		fmt.Println(failure.Err)
		return
	}

	metadataHeaders := []string{"Property", "Value"}
	metadataRows := [][]interface{}{
		{"Topic Name", metadata.Name},
		{"Replication Factor", len(metadata.Partitions[0].Replicas)},
		{"Version", metadata.Version},
		{"UUID", metadata.Uuid},
		{"Is Internal", metadata.IsInternal},
		{"Authorized Operations", metadata.TopicAuthorizedOperations},
	}
	RenderTable("Topic Metadata:", metadataHeaders, metadataRows)

	partitionHeaders := []string{"Partition ID", "Leader", "Replicas", "In-Sync Replicas (ISR)"}
	partitionRows := [][]interface{}{}
	for _, partition := range metadata.Partitions {
		partitionRows = append(partitionRows, []interface{}{
			partition.ID,
			partition.Leader,
			fmt.Sprintf("%v", partition.Replicas),
			fmt.Sprintf("%v", partition.Isr),
		})
	}
	RenderTable("Topic Partitions:", partitionHeaders, partitionRows)

	configs, failure := commands.DescribeTopicConfig(topicName)
	if failure != nil {
		fmt.Printf("Error describing configs for topic: %v\n", failure.Err)
		return
	}

	configHeaders := []string{"Config Name", "Value"}
	configRows := [][]interface{}{}
	for _, config := range configs {
		configRows = append(configRows, []interface{}{config.Name, config.Value})
	}
	RenderTable("Topic Configurations:", configHeaders, configRows)
}

// Update topic
func updateTopic(cmd cobraCmd, args cobraArgs) {
	topicName := cmd.Flags().Arg(0)
	newPartitions, _ := cmd.Flags().GetInt("new-partitions")

	if topicName == "" {
		fmt.Println("Error: Topic name is required.")
		return
	}

	if newPartitions <= 0 {
		fmt.Println("Error: Invalid partition count.")
		return
	}

	successMessage, failure := commands.UpdateTopic(topicName, newPartitions)
	if failure != nil {
		fmt.Println(failure.Err)
		return
	}
	fmt.Println(successMessage)
}
