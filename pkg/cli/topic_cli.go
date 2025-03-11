package cli

import (
	"fmt"

	"github.com/IBM/openkommander/internal/core"
)

type TopicCommandList struct{}

func (TopicCommandList) GetParentCommand() *OkParentCmd {
	return &OkParentCmd{
		Use:   "topics",
		Short: "Topic management commands",
	}
}

func (m TopicCommandList) GetCommands() []*OkCmd {
	return []*OkCmd{
		{ // Create
			Use:           "create",
			Short:         "Create a new topic",
			Run:           createTopic,
			Flags:         m.getCreateFlags(),
			RequiredFlags: []string{"name", "partitions", "replication-factor"},
		},
		{ // Delete
			Use:   "delete",
			Short: "Delete a topic",
			Run:   deleteTopic,
			Flags: []OkFlag{
				{
					Name:      "name",
					ShortName: "n",
					ValueType: "string",
					Usage:     "Specify the name of the topic to delete",
				},
			},
			RequiredFlags: []string{"name"},
		},
		{ // List
			Use:   "list",
			Short: "List all topics",
			Run:   listTopics,
		},
	}
}

func (TopicCommandList) GetSubcommands() []CommandList {
	return nil
}

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
			ShortName: "rf",
			ValueType: "int",
			Usage:     "Specify the replication factor of the new topic",
		},
	}
}

func createTopic(cmd cobraCmd, args cobraArgs) {
	name, _ := cmd.Flags().GetString("name")
	numPartitions, _ := cmd.Flags().GetInt("partitions")
	replicationFactor, _ := cmd.Flags().GetInt("replication-factor")

	fmt.Print(core.CreateTopic(name, numPartitions, replicationFactor))
}

func deleteTopic(cmd cobraCmd, args cobraArgs) {
	name, _ := cmd.Flags().GetString("name")

	fmt.Print(core.DeleteTopic(name))
}

func listTopics(cmd cobraCmd, args cobraArgs) {
	core.ListTopics()
}
