package commands

import (
	"fmt"

	"openkommander/pkg/session"

	"github.com/IBM/sarama"
)

func init() {
	Register("Cluster Data", "metadata", metadataCommand)
}

func metadataCommand() {
	currentSession := session.GetCurrentSession()

	if !currentSession.IsAuthenticated() {
		fmt.Println("Error: no session found.")
		return
	}

	client := currentSession.GetClient()

	if client == nil {
		fmt.Println("Error: not connected to a cluster.")
		return
	}

	brokers := client.Brokers()
	fmt.Println("Cluster Brokers:")
	for _, b := range brokers {
		if err := b.Open(client.Config()); err == nil || err == sarama.ErrAlreadyConnected {
			fmt.Printf(" - %s (ID: %d)\n", b.Addr(), b.ID())
		} else {
			fmt.Printf(" - %s (ID: %d) - error connecting: %v\n", b.Addr(), b.ID(), err)
		}
	}
}
