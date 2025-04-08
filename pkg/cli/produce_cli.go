package cli

import (
	"fmt"

	"github.com/IBM/openkommander/internal/core/commands"
)

type ProduceCommandList struct{}

func (ProduceCommandList) GetParentCommand() *OkParentCmd {
	return &OkParentCmd{
		Use:   "produce",
		Short: "Produce command",
		Run:   produceMessage,
		Flags: []OkFlag{
			NewOkFlag(OkFlagString, "topic", "t", "topic to write to"),
			NewOkFlag(OkFlagInt, "partition", "p", "[optional] partition to write message to", -1),
			NewOkFlag(OkFlagInt, "acks", "a", "[optional] acks flag, default -1 (full ISR).", -1),
			NewOkFlag(OkFlagString, "msg", "m", "message payload"),
			NewOkFlag(OkFlagString, "key", "k", "[optional] message key"),
		},
		RequiredFlags: []string{"topic", "msg"},
	}
}

func (m ProduceCommandList) GetCommands() []*OkCmd {
	return nil
}

func (ProduceCommandList) GetSubcommands() []CommandList {
	return nil
}

func produceMessage(cmd cobraCmd, args cobraArgs) {
	acks, _ := cmd.Flags().GetInt("acks")
	partition, _ := cmd.Flags().GetInt("partition")
	topic, _ := cmd.Flags().GetString("topic")
	msg, _ := cmd.Flags().GetString("msg")
	key, _ := cmd.Flags().GetString("key")

	successMessage, failure := commands.ProduceMessage(topic, key, msg, partition, acks)
	if failure != nil {
		fmt.Println(failure.Err)
		return
	}

	fmt.Println(successMessage)
}
