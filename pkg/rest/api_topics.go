package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/IBM/openkommander/pkg/logger"
	"github.com/IBM/openkommander/pkg/session"
	"github.com/IBM/sarama"
)

func ListTopics(w http.ResponseWriter, r *http.Request) {
	currentSession := session.GetCurrentSession()
	admin, err := currentSession.GetAdminClient()

	if err != nil {
		logger.Error("Failed to create admin client for listing topics", "cluster", session.GetActiveClusterName(), "error", err)
		SendError(w, "Failed to create admin client", err)
		return
	}

	topics, err := admin.ListTopics()

	if err != nil {
		logger.Error("Failed to list topics from Kafka", "cluster", session.GetActiveClusterName(), "error", err)
		SendError(w, "Failed to list topics", err)
		return
	}

	logger.Debug("Successfully retrieved topics", "cluster", session.GetActiveClusterName(), "topic_count", len(topics))

	topicList := make([]map[string]interface{}, 0, len(topics))
	for name, details := range topics {
		replicas := int(details.NumPartitions) * int(details.ReplicationFactor)
		inSyncReplicas := replicas
		topicList = append(topicList, map[string]interface{}{
			"name":               name,
			"partitions":         details.NumPartitions,
			"replication_factor": details.ReplicationFactor,
			"replicas":           replicas,
			"in_sync_replicas":   inSyncReplicas,
		})
	}

	SendJSON(w, http.StatusOK, Response{Status: "ok", Data: topicList})
}

func CreateTopic(w http.ResponseWriter, r *http.Request) {
	currentSession := session.GetCurrentSession()
	admin, err := currentSession.GetAdminClient()

	if err != nil {
		logger.Error("Failed to create admin client for topic creation", "cluster", session.GetActiveClusterName(), "error", err)
		SendError(w, "Failed to create admin client", err)
		return
	}

	var req TopicRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Invalid request body for topic creation", "cluster", session.GetActiveClusterName(), "error", err)
		SendError(w, "Invalid request body", err)
		return
	}

	logger.Info("Topic creation request details",
		"cluster", session.GetActiveClusterName(),
		"topic_name", req.Name,
		"partitions", req.Partitions,
		"replication_factor", req.ReplicationFactor)

	err = admin.CreateTopic(req.Name, &sarama.TopicDetail{
		NumPartitions:     req.Partitions,
		ReplicationFactor: req.ReplicationFactor,
	}, false)

	if err != nil {
		logger.Error("Failed to create topic in Kafka", "cluster", session.GetActiveClusterName(), "topic_name", req.Name, "error", err)
		SendError(w, "Failed to create topic", err)
		return
	}

	logger.Info("Topic created successfully", "cluster", session.GetActiveClusterName(), "topic_name", req.Name, "partitions", req.Partitions, "replication_factor", req.ReplicationFactor)
	SendJSON(w, http.StatusCreated, Response{Status: "ok", Message: fmt.Sprintf("Topic '%s' created successfully", req.Name)})
}

func DeleteTopic(w http.ResponseWriter, r *http.Request) {
	currentSession := session.GetCurrentSession()
	admin, err := currentSession.GetAdminClient()

	if err != nil {
		logger.Error("Failed to create admin client for topic deletion", "cluster", session.GetActiveClusterName(), "error", err)
		SendError(w, "Failed to create admin client", err)
		return
	}

	topicName := r.PathValue("name")

	if topicName == "" {
		logger.Warn("Topic name is required for deletion", "cluster", session.GetActiveClusterName())
		SendError(w, "Topic name is required", nil)
		return
	}

	logger.Info("Topic deletion request details", "cluster", session.GetActiveClusterName(), "topic_name", topicName)

	err = admin.DeleteTopic(topicName)
	if err != nil {
		logger.Error("Failed to delete topic from Kafka", "cluster", session.GetActiveClusterName(), "topic_name", topicName, "error", err)
		SendError(w, "Failed to delete topic", err)
		return
	}

	logger.Info("Topic deleted successfully", "cluster", session.GetActiveClusterName(), "topic_name", topicName)
	SendJSON(w, http.StatusOK, Response{Status: "ok", Message: fmt.Sprintf("Topic '%s' deleted successfully", topicName)})
}
