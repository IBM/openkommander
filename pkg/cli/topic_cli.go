package cli

import (
	"fmt"

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
			Use:   "create",
			Short: "Create a new topic",
			Run:   createTopic,
			Flags: []OkFlag{
				NewOkFlag(OkFlagString, "name", "n", "Specify the name of the new topic"),
				NewOkFlag(OkFlagInt, "partitions", "p", "Specify the number of partitions of the new topic"),
				NewOkFlag(OkFlagInt, "replication-factor", "r", "Specify the replication factor of the new topic"),
			},
			RequiredFlags: []string{"name", "partitions", "replication-factor"},
		},
		{ // Delete topic
			Use:   "delete [NAME]",
			Short: "Delete a topic",
			Run:   deleteTopic,
			Args:  cobra.ExactArgs(1),
		},
		{ // List topics
			Use:   "list",
			Short: "List all topics",
			Run:   listTopics,
		},
	}
}

func (TopicCommandList) GetSubcommands() []CommandList {
	return nil
}

// Create topic

func createTopic(cmd cobraCmd, args cobraArgs) {
	name, _ := cmd.Flags().GetString("name")
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
