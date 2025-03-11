package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/IBM/openkommander/internal/core/commands"
	"github.com/IBM/openkommander/pkg/session"
	"github.com/gorilla/mux"
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

	success, failure := commands.CreateTopic(req.Name, req.Partitions, req.ReplicationFactor)
	if failure != nil {
		http.Error(w, failure.Err.Error(), failure.HttpCode)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]string{"message": success.Body}); err != nil {
		fmt.Println("Error encoding response:", err)
	}
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

	err = json.NewEncoder(w).Encode(topics)
	if err != nil {
		fmt.Println("Error encoding response:", err)
	}
}

func DeleteTopicHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	success, failure := commands.DeleteTopic(name)
	if failure != nil {
		http.Error(w, failure.Err.Error(), failure.HttpCode)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"message": success.Body}); err != nil {
		fmt.Println("Error encoding response:", err)
	}
}
