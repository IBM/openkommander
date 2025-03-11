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
			Use:           "create",
			Short:         "Create a new topic",
			Run:           createTopic,
			Flags:         m.getCreateFlags(),
			RequiredFlags: m.getCreateRequiredFlags(),
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

func (TopicCommandList) getCreateRequiredFlags() []string {
	return []string{"name", "partitions", "replication-factor"}
}

func (TopicCommandList) getCreateFlags() []OkFlag {
	return []OkFlag{
		{
			Name:      "name",
			ShortName: "n",
			ValueType: "string",
			Usage:     "Specify the name of the new topic",
		},
		{
			Name:      "partitions",
			ShortName: "p",
			ValueType: "int",
			Usage:     "Specify the number of partitions of the new topic",
		},
		{
			Name:      "replication-factor",
			ShortName: "r",
			ValueType: "int",
			Usage:     "Specify the replication factor of the new topic",
		},
	}
}

func createTopic(cmd cobraCmd, args cobraArgs) {
	name, _ := cmd.Flags().GetString("name")
	numPartitions, _ := cmd.Flags().GetInt("partitions")
	replicationFactor, _ := cmd.Flags().GetInt("replication-factor")

	success, failure := commands.CreateTopic(name, numPartitions, replicationFactor)
	if failure != nil {
		fmt.Println(failure.Err)
		return
	}

	fmt.Println(success.Body)
}

// Delete topic

func deleteTopic(cmd cobraCmd, args cobraArgs) {
	name := cmd.Flags().Arg(0)

	success, failure := commands.DeleteTopic(name)
	if failure != nil {
		fmt.Println(failure.Err)
		return
	}

	fmt.Println(success.Body)
}

// List topics

func listTopics(cmd cobraCmd, args cobraArgs) {
	success, failure := commands.ListTopics()
	if failure != nil {
		fmt.Println(failure.Err)
		return
	}

	topics := success.Body

	fmt.Println("\nTopics:")
	fmt.Println("--------")
	for name, detail := range topics {
		fmt.Printf("Name: %s\n", name)
		fmt.Printf("  Partitions: %d\n", detail.NumPartitions)
		fmt.Printf("  Replication Factor: %d\n", detail.ReplicationFactor)
		fmt.Println()
	}
}
