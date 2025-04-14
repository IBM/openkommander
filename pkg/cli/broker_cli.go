package cli

import (
	"fmt"

	"github.com/IBM/openkommander/internal/core/commands"
	// "github.com/spf13/cobra"
)

type BrokerCommandList struct{}

func (BrokerCommandList) GetParentCommand() *OkParentCmd {
	return &OkParentCmd{
		Use:   "broker <command>",
		Short: "Broker management commands",
	}
}

func (m BrokerCommandList) GetCommands() []*OkCmd {
	return []*OkCmd{

		{ // List broker info
			Use:   "info",
			Short: "List all broker info",
			Run:   getBrokerInfo,
		},
	}
}

func (BrokerCommandList) GetSubcommands() []CommandList {
	return nil
}

// List Broker info
func getBrokerInfo(cmd cobraCmd, args cobraArgs) {
	client, failure := commands.GetBrokerClient()
	brokers := client.Brokers()

	if failure != nil {
		fmt.Println(failure.Err)
		return
	}

	brokerHeaders := []string{"ID", "Address", "Rack", "Connected", "ResponseSize"}
	brokerRows := [][]interface{}{}
	for _, broker := range brokers {
		connected, _ := broker.Connected()

		brokerRows = append(brokerRows, []interface{}{
			broker.ID(),
			broker.Addr(),
			broker.Rack(),
			connected,
			broker.ResponseSize(),
		})
	}
	RenderTable("Broker Information:", brokerHeaders, brokerRows)
}
