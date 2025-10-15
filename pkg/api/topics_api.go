package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/IBM/openkommander/internal/core/commands"
)

type CreateTopicRequest struct {
	Name              string `json:"name"`
	Partitions        int    `json:"partitions"`
	ReplicationFactor int    `json:"replication_factor"`
}

func CreateTopicHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateTopicRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON input", http.StatusBadRequest)
		return
	}

	successMessage, failure := commands.CreateTopic(req.Name, req.Partitions, req.ReplicationFactor)
	if failure != nil {
		http.Error(w, failure.Err.Error(), failure.HttpCode)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]string{"message": successMessage}); err != nil {
		fmt.Println("Error encoding response:", err)
	}
}

func ListTopicsHandler(w http.ResponseWriter, r *http.Request) {
	broker := r.PathValue("broker")

	fmt.Println("API - Listing topics for broker:", broker)

	topics, failure := commands.ListTopics()
	if failure != nil {
		http.Error(w, failure.Err.Error(), failure.HttpCode)
		return
	}

	err := json.NewEncoder(w).Encode(topics)
	if err != nil {
		fmt.Println("Error encoding response:", err)
	}
}

func DeleteTopicHandler(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")

	successMessage, failure := commands.DeleteTopic(name)
	if failure != nil {
		http.Error(w, failure.Err.Error(), failure.HttpCode)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"message": successMessage}); err != nil {
		fmt.Println("Error encoding response:", err)
	}
}
