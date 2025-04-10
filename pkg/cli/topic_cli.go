package cli

import (
	"fmt"
	"os"

	"github.com/IBM/openkommander/internal/core/commands"
	"github.com/jedib0t/go-pretty/v6/table"
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
			RequiredFlags: []string{"partitions", "replication-factor"},
			Args:          cobra.ExactArgs(1),
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

	fmt.Println("\nTopics:")
	fmt.Println("--------")
	for name, detail := range topics {
		fmt.Printf("Name: %s\n", name)
		fmt.Printf("  Partitions: %d\n", detail.NumPartitions)
		fmt.Printf("  Replication Factor: %d\n", detail.ReplicationFactor)
		fmt.Println()
	}
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

	// Print Topic Metadata
	fmt.Println("\nTopic Metadata:")
	fmt.Printf("  Topic Name: %s\n", metadata.Name)
	fmt.Printf("  Replication Factor: %d\n", len(metadata.Partitions[0].Replicas))
	fmt.Printf("  Version: %d\n", metadata.Version)
	fmt.Printf("  UUID: %s\n", metadata.Uuid)
	fmt.Printf("  Is Internal: %t\n", metadata.IsInternal)
	fmt.Printf("  Authorized Operations: %d\n", metadata.TopicAuthorizedOperations)

	// Table for partition details
	fmt.Println("\nTopic Partitions:")
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Partition ID", "Leader", "Replicas", "In-Sync Replicas (ISR)"})

	for _, partition := range metadata.Partitions {
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
	configs, failure := commands.DescribeTopicConfig(topicName)
	if failure != nil {
		fmt.Printf("Error describing configs for topic: %v\n", failure.Err)
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
