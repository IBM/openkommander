package commands

import (
	"fmt"

	"github.com/IBM/openkommander/pkg/session"
	"github.com/IBM/sarama"
)

// When successful, returns a map of topic names to their details
func GetBrokerClient() (client sarama.Client, f *Failure) {
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
	// brokers := client.Brokers()

	return client, nil
}
