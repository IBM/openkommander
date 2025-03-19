package broker

import (
	"fmt"
	"github.com/IBM/openkommander/pkg/session"
	"github.com/IBM/sarama"
	"os"
    "github.com/jedib0t/go-pretty/v6/table"
)

type Broker struct{
	Info string
}

func GetInfo()  {
	currentSession := session.GetCurrentSession()
	if !currentSession.IsAuthenticated() {
		fmt.Println("Error: no session found.")
		return
	}

	client, err := currentSession.GetClient()
	if err != nil {
		fmt.Printf("Error connecting to cluster: %v\n", err)
		return
	}

	brokers := client.Brokers()
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