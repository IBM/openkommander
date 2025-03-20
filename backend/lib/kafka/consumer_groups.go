package kafka

import (
	"fmt"
	"github.com/IBM/sarama"
	"openkommander/models"
)

func (c *Client) GetConsumerGroups() ([]models.ConsumerGroupInfo, error) {
    groups, err := c.admin.ListConsumerGroups()
    if err != nil {
        return nil, fmt.Errorf("failed to list consumer groups: %w", err)
    }

    groupIDs := make([]string, 0, len(groups))
    for groupID := range groups {
        groupIDs = append(groupIDs, groupID)
    }

    descriptions, err := c.admin.DescribeConsumerGroups(groupIDs)
    if err != nil {
        return nil, fmt.Errorf("failed to describe consumer groups: %w", err)
    }

    client, err := sarama.NewClient(c.brokers, c.config)
    if err != nil {
        return nil, fmt.Errorf("failed to create client: %w", err)
    }
    defer client.Close()

    result := make([]models.ConsumerGroupInfo, 0, len(descriptions))

    for _, desc := range descriptions {
        groupInfo := models.ConsumerGroupInfo{
            GroupID:     desc.GroupId,
            Members:     len(desc.Members),
            Topics:      0,
            Lag:         0,
            Coordinator: -1, 
            State:       string(desc.State),
            TopicLags:   []models.TopicLagInfo{},
        }

        coordinatorID, err := client.Coordinator(desc.GroupId)
        if err == nil && coordinatorID != nil {
            groupInfo.Coordinator = coordinatorID.ID()
        }

        topicsMap := make(map[string]struct{})
        for _, member := range desc.Members {
            metadata, err := member.GetMemberMetadata()
            if err != nil {
                continue
            }

            for _, topic := range metadata.Topics {
                topicsMap[topic] = struct{}{}
            }
        }

        groupInfo.Topics = len(topicsMap)
        topics := make([]string, 0, len(topicsMap))
        for topic := range topicsMap {
            topics = append(topics, topic)
        }

        for _, topic := range topics {
            partitions, err := client.Partitions(topic)
            if err != nil {
                continue
            }

            for _, partition := range partitions {
                latestOffset, err := client.GetOffset(topic, partition, sarama.OffsetNewest)
                if err != nil {
                    continue
                }

                consumerOffset, err := c.admin.ListConsumerGroupOffsets(desc.GroupId, map[string][]int32{
                    topic: {partition},
                })
                if err != nil {
                    continue
                }

                if block, ok := consumerOffset.Blocks[topic]; ok {
                    if offsetFetchResponse, ok := block[partition]; ok {
                        if offsetFetchResponse.Offset != -1 { // -1 means no offset committed
                            lag := latestOffset - offsetFetchResponse.Offset
                            groupInfo.Lag += lag
                            
                            groupInfo.TopicLags = append(groupInfo.TopicLags, models.TopicLagInfo{
                                Topic:     topic,
                                Partition: partition,
                                Lag:       lag,
                            })
                        }
                    }
                }
            }
        }

        result = append(result, groupInfo)
    }

    return result, nil
}

func (c *Client) GetConsumerGroup(groupID string) (*models.ConsumerGroupInfo, error) {
	groups, err := c.GetConsumerGroups()
	if err != nil {
		return nil, err
	}

	for _, group := range groups {
		if group.GroupID == groupID {
			return &group, nil
		}
	}

	return nil, fmt.Errorf("consumer group not found: %s", groupID)
}
