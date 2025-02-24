package commands

import (
	"fmt"

	"github.com/IBM/openkommander/pkg/session"
	"github.com/IBM/sarama"
)

func metadataCommand() {
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
	for _, b := range brokers {
		if err := b.Open(client.Config()); err == nil || err == sarama.ErrAlreadyConnected {
			fmt.Printf(" - %s (ID: %d)\n", b.Addr(), b.ID())
		} else {
			fmt.Printf(" - %s (ID: %d) - error connecting: %v\n", b.Addr(), b.ID(), err)
		}
	}
}
