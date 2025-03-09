package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/IBM/openkommander/pkg/session"
	"github.com/IBM/sarama"
	"github.com/gorilla/mux"
)

type TopicRequest struct {
	Name              string `json:"name"`
	Partitions        int    `json:"partitions"`
	ReplicationFactor int    `json:"replication_factor"`
}

func CreateTopicHandler(w http.ResponseWriter, r *http.Request) {
	var req TopicRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON input", http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.Partitions < 1 || req.ReplicationFactor < 1 {
		http.Error(w, "Invalid topic parameters", http.StatusBadRequest)
		return
	}

	currentSession := session.GetCurrentSession()
	if !currentSession.IsAuthenticated() {
		http.Error(w, "No active session", http.StatusUnauthorized)
		return
	}

	client, err := currentSession.GetAdminClient()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error connecting to Kafka: %v", err), http.StatusInternalServerError)
		return
	}

	topicDetail := &sarama.TopicDetail{
		NumPartitions:     int32(req.Partitions),
		ReplicationFactor: int16(req.ReplicationFactor),
	}

	err = client.CreateTopic(req.Name, topicDetail, false)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating topic: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Topic created successfully"})
}

func ListTopicsHandler(w http.ResponseWriter, r *http.Request) {
	currentSession := session.GetCurrentSession()
	if !currentSession.IsAuthenticated() {
		http.Error(w, "No active session", http.StatusUnauthorized)
		return
	}

	client, err := currentSession.GetAdminClient()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error connecting to Kafka: %v", err), http.StatusInternalServerError)
		return
	}

	topics, err := client.ListTopics()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error listing topics: %v", err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(topics)
}

func DeleteTopicHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	topicName := vars["topicName"]

	if topicName == "" {
		http.Error(w, "Topic name is required", http.StatusBadRequest)
		return
	}

	currentSession := session.GetCurrentSession()
	if !currentSession.IsAuthenticated() {
		http.Error(w, "No active session", http.StatusUnauthorized)
		return
	}

	client, err := currentSession.GetAdminClient()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error connecting to Kafka: %v", err), http.StatusInternalServerError)
		return
	}

	err = client.DeleteTopic(topicName)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error deleting topic: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Topic deleted successfully"})
}
