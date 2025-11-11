package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/IBM/openkommander/pkg/logger"
)

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type TopicRequest struct {
	Name              string `json:"name"`
	Partitions        int32  `json:"partitions"`
	ReplicationFactor int16  `json:"replication_factor"`
}

func SendJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		logger.Error("Failed to encode JSON response", "error", err)
	}
}

func SendError(w http.ResponseWriter, message string, err error) {
	logger.Error(message, "error", err)
	SendJSON(w, http.StatusInternalServerError, Response{
		Status:  "error",
		Message: fmt.Sprintf("%s: %v", message, err),
	})
}
