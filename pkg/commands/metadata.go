package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/IBM/openkommander/pkg/session"
	"github.com/IBM/sarama"
)

func metadataCommand() {
	currentSession := session.GetCurrentSession()
	if !currentSession.IsAuthenticated() {
		fmt.Println("Error: no session found.")
		return
	}

	client := currentSession.GetClient()
	if client == nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_, err := currentSession.Connect(ctx)
		if err != nil {
			fmt.Printf("Error connecting to cluster: %v\n", err)
			return
		}
		client = currentSession.GetClient()
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
