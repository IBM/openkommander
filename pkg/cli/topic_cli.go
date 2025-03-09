package cli

import "github.com/IBM/openkommander/pkg/functions"

type TopicCommands struct {
	Topic *OkParentCmd

	Create *OkCmd
	Delete *OkCmd
	List   *OkCmd
}

type OkFlagList interface {
	GetFlags() []OkFlag
}

type CreateTopicFlags struct {
	GetFlags func() []OkFlag
}

func GetTopicCommands() *TopicCommands {
	topicCommands := &TopicCommands{
		Topic: &OkParentCmd{
			Use:   "topics",
			Short: "Topic management commands",
		},
		Create: &OkCmd{
			Use:   "create",
			Short: "Create a new topic",
			Run:   createTopic,
			Flags: []OkFlag{
				{
					Name:      "name",
					ShortName: "n",
					ValueType: "string",
					Usage:     "Specify the name of the new topic",
				},
				{
					Name:      "partitions",
					ValueType: "int",
					Usage:     "Specify the number of partitions of the new topic",
				},
				{
					Name:      "replication-factor",
					ValueType: "int",
					Usage:     "Specify the replication factor of the new topic",
				},
			},
			EnforceFlagConstraints: func(cmd cobraCmd) {
				cmd.MarkFlagRequired("name")
				cmd.MarkFlagRequired("partitions")
				cmd.MarkFlagRequired("replication-factor")
			},
		},
		Delete: &OkCmd{
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
			EnforceFlagConstraints: func(cmd cobraCmd) {
				cmd.MarkFlagRequired("name")
			},
		},
		List: &OkCmd{
			Use:   "list",
			Short: "List all topics",
			Run:   listTopics,
		},
	}

	return topicCommands
}

func createTopic(cmd cobraCmd, args cobraArgs) {
	name, _ := cmd.Flags().GetString("name")
	numPartitions, _ := cmd.Flags().GetInt("partitions")
	replicationFactor, _ := cmd.Flags().GetInt("replication-factor")

	functions.CreateTopic(name, numPartitions, replicationFactor)
}

func deleteTopic(cmd cobraCmd, args cobraArgs) {
	name, _ := cmd.Flags().GetString("name")

	functions.DeleteTopic(name)
}

func listTopics(cmd cobraCmd, args cobraArgs) {
	functions.ListTopics()
}
