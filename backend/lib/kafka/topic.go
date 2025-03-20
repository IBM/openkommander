package kafka

import (
	"fmt"
	"strings"

	"github.com/IBM/sarama"
	"openkommander/models"
)

func (c *Client) ListTopics() (map[string]sarama.TopicDetail, error) {
	return c.admin.ListTopics()
}

func (c *Client) CreateTopic(name string, partitions int32, replicationFactor int16) error {
	topicDetail := &sarama.TopicDetail{
		NumPartitions:     partitions,
		ReplicationFactor: replicationFactor,
	}
	return c.admin.CreateTopic(name, topicDetail, false)
}

func (c *Client) DeleteTopic(name string) error {
	return c.admin.DeleteTopic(name)
}

func (c *Client) DescribeTopic(name string) (*sarama.TopicDetail, error) {
	topics, err := c.admin.ListTopics()
	if err != nil {
		return nil, err
	}

	if detail, exists := topics[name]; exists {
		return &detail, nil
	}

	return nil, fmt.Errorf("topic %s not found", name)
}


func (c *Client) GetTopicPartitions(topic string) ([]int32, error) {
	client, err := sarama.NewClient(c.brokers, c.config)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}
	defer client.Close()

	return client.Partitions(topic)
}

func (c *Client) GetTopicInfo() ([]models.TopicInfo, error) {
    client, err := sarama.NewClient(c.brokers, c.config)
    if err != nil {
        return nil, fmt.Errorf("failed to create client: %w", err)
    }
    defer client.Close()

    topics, err := c.admin.ListTopics()
    if err != nil {
        return nil, fmt.Errorf("failed to list topics: %w", err)
    }

    result := make([]models.TopicInfo, 0, len(topics))

    for name, detail := range topics {
        info := models.TopicInfo{
            Name:             name,
            Partitions:       detail.NumPartitions,
            ReplicationFactor: detail.ReplicationFactor,
            Internal:         strings.HasPrefix(name, "_"),
        }

        configEntries, err := c.admin.DescribeConfig(sarama.ConfigResource{
            Type: sarama.TopicResource,
            Name: name,
        })
        if err == nil {
            for _, entry := range configEntries {
                if entry.Name == "cleanup.policy" {
                    info.CleanupPolicy = entry.Value
                    break
                }
            }
        }

        partitions, err := client.Partitions(name)
        if err != nil {
            fmt.Printf("Warning: Failed to get partitions for topic %s: %v\n", name, err)
            continue
        }

        var totalReplicas, totalInSyncReplicas int

        for _, partition := range partitions {
            replicas, err := client.Replicas(name, partition)
            if err == nil {
                totalReplicas += len(replicas)
            }

            inSyncReplicas, err := client.InSyncReplicas(name, partition)
            if err == nil {
                totalInSyncReplicas += len(inSyncReplicas)
            }
        }

        info.Replicas = totalReplicas
        info.InSyncReplicas = totalInSyncReplicas

        result = append(result, info)
    }

    return result, nil
}

func (c *Client) GetTopic(name string) (*models.TopicDetail, error) {
	detail, err := c.DescribeTopic(name)
	if err != nil {
		return nil, err
	}
	
	partitions, err := c.GetTopicPartitions(name)
	if err != nil {
		return nil, fmt.Errorf("failed to get topic partitions: %w", err)
	}
	
	return &models.TopicDetail{
		Name:              name,
		Partitions:        detail.NumPartitions,
		ReplicationFactor: detail.ReplicationFactor,
		PartitionIDs:      partitions,
	}, nil
}
