package cli

import (
	"fmt"
	"os"

	"github.com/IBM/openkommander/internal/core/commands"
	"github.com/IBM/sarama"
	"github.com/jedib0t/go-pretty/v6/table"
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
	fmt.Println("Cluster Brokers:")
	t := table.NewWriter()
    t.SetOutputMirror(os.Stdout)
    t.AppendHeader(table.Row{"ID", "Address", "Rack", "Connected", "ResponseSize"})

	for _, b := range brokers {
		if err := b.Open(client.Config()); err == nil || err == sarama.ErrAlreadyConnected {
			connected, connErr := b.Connected()
			if connErr != nil {
				fmt.Printf(" - %s (ID: %d) Rack: %s - error checking connection: %v\n", b.Addr(), b.ID(), b.Rack(), connErr)
			} else {
				t.AppendRows([]table.Row{
					{b.ID(), b.Addr(), b.Rack(), connected, b.ResponseSize()},
				})
				t.AppendSeparator()
			}
		} else {
			fmt.Printf(" - %s (ID: %d) - error connecting: %v\n", b.Addr(), b.ID(), err)
		}
	}
	t.Render()
}
  