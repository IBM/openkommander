package rest

import (
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/IBM/openkommander/pkg/logger"
	"github.com/IBM/openkommander/pkg/session"
	"github.com/IBM/sarama"
)

// MessagesLastMinute holds the count of produced and consumed messages in the last minute and last second
type MessagesLastMinute struct {
	Topic          string `json:"topic"`
	ProducedCount  int64  `json:"produced_count"`
	ConsumedCount  int64  `json:"consumed_count"`
	ProducedPerSec int64  `json:"produced_per_sec"`
	ConsumedPerSec int64  `json:"consumed_per_sec"`
}

// OffsetHistoryEntry stores offset and timestamp
type OffsetHistoryEntry struct {
	Offset    int64
	Timestamp time.Time
}

// In-memory rolling window for offsets (not production safe)
var producedHistory = make(map[string][]OffsetHistoryEntry)
var consumedHistory = make(map[string][]OffsetHistoryEntry)

// HandleMessagesPerMinute handles requests for message metrics
func HandleMessagesPerMinute(w http.ResponseWriter, r *http.Request) {
	currentSession := session.GetCurrentSession()
	admin, err := currentSession.GetAdminClient()
	if err != nil {
		logger.Error("Failed to create admin client for messages per minute", "cluster", session.GetActiveClusterName(), "error", err)
		SendError(w, "Failed to create admin client", err)
		return
	}
	if admin == nil {
		logger.Error("Admin client creation failed for messages per minute", "cluster", session.GetActiveClusterName())
		SendError(w, "Failed to create admin client", fmt.Errorf("admin client creation failed"))
		return
	}

	kafkaClient, err := currentSession.GetClient()
	if err != nil {
		logger.Error("Failed to get Kafka client for messages per minute", "cluster", session.GetActiveClusterName(), "error", err)
		SendError(w, "Failed to get Kafka client", err)
		return
	}
	if kafkaClient == nil {
		logger.Error("Kafka client is nil for messages per minute", "cluster", session.GetActiveClusterName())
		SendError(w, "Kafka client not available", fmt.Errorf("kafka client is nil"))
		return
	}

	topics, err := admin.ListTopics()

	if err != nil {
		logger.Error("Failed to list topics for messages per minute", "cluster", session.GetActiveClusterName(), "error", err)
		SendError(w, "Failed to list topics", err)
		return
	}

	logger.Debug("Processing message metrics", "cluster", session.GetActiveClusterName(), "topic_count", len(topics))

	// Calculate produced and consumed message counts in the last minute for each topic
	counts := []MessagesLastMinute{}
	now := time.Now()
	totalAllProduced := int64(0)
	totalAllConsumed := int64(0)
	totalAllProducedSec := int64(0)
	totalAllConsumedSec := int64(0)

	for name := range topics { // Produced: sum latest offsets across all partitions
		partitions, err := kafkaClient.Partitions(name)
		if err != nil {
			logger.Warn("Failed to get partitions for topic", "cluster", session.GetActiveClusterName(), "topic", name, "error", err)
		}
		var totalProduced int64 = 0
		for _, partition := range partitions {
			offset, err := kafkaClient.GetOffset(name, partition, sarama.OffsetNewest)
			if err == nil {
				totalProduced += offset
			}
		}
		// Consumed: sum committed offsets for all consumer groups
		var totalConsumed int64 = 0
		groups, err := admin.ListConsumerGroups()
		if err != nil {
			logger.Warn("Failed to list consumer groups", "cluster", session.GetActiveClusterName(), "error", err)
		}

		for group := range groups {
			offsets, err := admin.ListConsumerGroupOffsets(group, map[string][]int32{name: partitions})
			if err == nil && offsets.Blocks != nil {
				for _, partition := range partitions {
					block := offsets.GetBlock(name, partition)
					if block != nil && block.Offset > 0 {
						totalConsumed += block.Offset
					}
				}
			}
		} // Store offset history for rolling window
		producedHistory[name] = append(producedHistory[name], OffsetHistoryEntry{Offset: totalProduced, Timestamp: now})
		consumedHistory[name] = append(consumedHistory[name], OffsetHistoryEntry{Offset: totalConsumed, Timestamp: now})

		// Remove entries older than 1 minute using slices.DeleteFunc
		pruneMinute := now.Add(-1 * time.Minute)
		producedHistory[name] = slices.DeleteFunc(producedHistory[name], func(entry OffsetHistoryEntry) bool {
			return entry.Timestamp.Before(pruneMinute)
		})
		consumedHistory[name] = slices.DeleteFunc(consumedHistory[name], func(entry OffsetHistoryEntry) bool {
			return entry.Timestamp.Before(pruneMinute)
		})

		// Create filtered copies for second-based calculations
		pruneSecond := now.Add(-1 * time.Second)
		producedHistorySec := slices.DeleteFunc(slices.Clone(producedHistory[name]), func(entry OffsetHistoryEntry) bool {
			return entry.Timestamp.Before(pruneSecond)
		})
		consumedHistorySec := slices.DeleteFunc(slices.Clone(consumedHistory[name]), func(entry OffsetHistoryEntry) bool {
			return entry.Timestamp.Before(pruneSecond)
		})

		// Calculate produced/consumed in last minute
		producedCount := int64(0)
		consumedCount := int64(0)
		if len(producedHistory[name]) > 1 {
			producedCount = producedHistory[name][len(producedHistory[name])-1].Offset - producedHistory[name][0].Offset
		}
		if len(consumedHistory[name]) > 1 {
			consumedCount = consumedHistory[name][len(consumedHistory[name])-1].Offset - consumedHistory[name][0].Offset
		}

		// Calculate produced/consumed in last second
		producedPerSec := int64(0)
		consumedPerSec := int64(0)
		if len(producedHistorySec) > 1 {
			producedPerSec = producedHistorySec[len(producedHistorySec)-1].Offset - producedHistorySec[0].Offset
		}
		if len(consumedHistorySec) > 1 {
			consumedPerSec = consumedHistorySec[len(consumedHistorySec)-1].Offset - consumedHistorySec[0].Offset
		}

		totalAllProduced += producedCount
		totalAllConsumed += consumedCount
		totalAllProducedSec += producedPerSec
		totalAllConsumedSec += consumedPerSec

		counts = append(counts, MessagesLastMinute{
			Topic:          name,
			ProducedCount:  producedCount,
			ConsumedCount:  consumedCount,
			ProducedPerSec: producedPerSec,
			ConsumedPerSec: consumedPerSec,
		})
	}

	counts = append(counts, MessagesLastMinute{
		Topic:          "total",
		ProducedCount:  totalAllProduced,
		ConsumedCount:  totalAllConsumed,
		ProducedPerSec: totalAllProducedSec,
		ConsumedPerSec: totalAllConsumedSec,
	})
	logger.Info("Successfully calculated message metrics", "cluster", session.GetActiveClusterName(), "total_topics", len(counts)-1, "total_produced", totalAllProduced, "total_consumed", totalAllConsumed)
	SendJSON(w, http.StatusOK, Response{Status: "ok", Data: counts})
}
