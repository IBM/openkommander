package commands

import (
	"fmt"

	"github.com/IBM/openkommander/pkg/session"
)

func brokersListCommand() {
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

	brokers, controllerID, err := client.DescribeCluster()
	if err != nil {
		fmt.Printf("Error getting cluster metadata: %v\n", err)
		return
	}

	if len(brokers) == 0 {
		fmt.Println("No brokers found")
		return
	}

	fmt.Println("\nBrokers:")
	fmt.Println("--------")
	for _, broker := range brokers {
		fmt.Printf("ID: %d", broker.ID())
		if broker.ID() == controllerID {
			fmt.Printf(" (controller)")
		}
		fmt.Println()
		fmt.Printf("  Address: %s\n", broker.Addr())

		// Check if the broker is connected
		connected, connectErr := broker.Connected()
		if connectErr != nil {
			fmt.Printf("  Connected: false (error checking connection: %v)\n", connectErr)
		} else {
			fmt.Printf("  Connected: %t\n", connected)
		}

		if rack := broker.Rack(); rack != "" {
			fmt.Printf("  Rack: %s\n", rack)
		}

		fmt.Println()
	}
}
