package kafka

import (
	"fmt"
	"github.com/IBM/sarama"
	"openkommander/lib/utils"
	"openkommander/models"
)

func (c *Client) GetBrokerInfo() ([]models.BrokerInfo, error) {
    client, err := sarama.NewClient(c.brokers, c.config)
    if err != nil {
        return nil, fmt.Errorf("failed to create client: %w", err)
    }
    defer client.Close()

    brokers := client.Brokers()
    if len(brokers) == 0 {
        return nil, fmt.Errorf("no brokers found")
    }

    brokerInfoList := make([]models.BrokerInfo, 0, len(brokers))

    topics, err := client.Topics()
    if err != nil {
        return nil, fmt.Errorf("failed to get topics: %w", err)
    }

    for _, broker := range brokers {
        brokerInfo := models.BrokerInfo{
            ID:               broker.ID(),
            Host:             "",
            Port:             0,
            PartitionsLeader: 0,
            Partitions:       0,
            InSyncPartitions: 0,
        }

        if addr := broker.Addr(); addr != "" {
            host, port, err := utils.SplitHostPort(addr)
            if err == nil {
                brokerInfo.Host = host
                brokerInfo.Port = int32(port)
            } else {
                fmt.Printf("Warning: Failed to parse broker address %s: %v\n", addr, err)
            }
        }

        for _, topic := range topics {
            partitions, err := client.Partitions(topic)
            if err != nil {
                continue
            }

            for _, partition := range partitions {
                replicas, err := client.Replicas(topic, partition)
                if err != nil {
                    continue
                }

                isBrokerReplica := false
                for _, replicaID := range replicas {
                    if replicaID == broker.ID() {
                        isBrokerReplica = true
                        brokerInfo.Partitions++
                        break
                    }
                }

                if !isBrokerReplica {
                    continue
                }

                leader, err := client.Leader(topic, partition)
                if err == nil && leader.ID() == broker.ID() {
                    brokerInfo.PartitionsLeader++
                }

                isr, err := client.InSyncReplicas(topic, partition)
                if err != nil {
                    continue
                }

                for _, isrID := range isr {
                    if isrID == broker.ID() {
                        brokerInfo.InSyncPartitions++
                        break
                    }
                }
            }
        }

        brokerInfoList = append(brokerInfoList, brokerInfo)
    }

    return brokerInfoList, nil
}