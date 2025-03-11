package functions

import (
	"fmt"

	"github.com/IBM/openkommander/pkg/session"
)

func ListTopics() {
	currentSession := session.GetCurrentSession()
	if !currentSession.IsAuthenticated() {
		fmt.Println("Error: no session found.")
		return
	}

	client, err := currentSession.GetAdminClient()
	if err != nil {
		fmt.Printf("Error connecting to cluster: %v\n", err)
		return
	}

	topics, err := client.ListTopics()
	if err != nil {
		fmt.Printf("Error listing topics: %v\n", err)
		return
	}

	if len(topics) == 0 {
		fmt.Println("No topics found")
		return
	}

	fmt.Println("\nTopics:")
	fmt.Println("--------")
	for name, detail := range topics {
		fmt.Printf("Name: %s\n", name)
		fmt.Printf("  Partitions: %d\n", detail.NumPartitions)
		fmt.Printf("  Replication Factor: %d\n", detail.ReplicationFactor)
		fmt.Println()
	}
}
