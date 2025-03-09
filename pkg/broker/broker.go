package broker

import (
	"fmt"
	"github.com/IBM/openkommander/pkg/session"
	"github.com/IBM/sarama"
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
	for _, b := range brokers {
		if err := b.Open(client.Config()); err == nil || err == sarama.ErrAlreadyConnected {
			connected, connErr := b.Connected()
			if connErr != nil {
				fmt.Printf(" - %s (ID: %d) Rack: %s - error checking connection: %v\n", b.Addr(), b.ID(), b.Rack(), connErr)
			} else {
				fmt.Printf(" - %s (ID: %d) Rack: %s Connected: %t\n", b.Addr(), b.ID(), b.Rack(), connected)
			}
		} else {
			fmt.Printf(" - %s (ID: %d) - error connecting: %v\n", b.Addr(), b.ID(), err)
		}
	}
}