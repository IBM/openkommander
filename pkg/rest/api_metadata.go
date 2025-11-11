package rest

import (
	"fmt"
	"net/http"

	"github.com/IBM/openkommander/internal/core/commands"
	"github.com/IBM/openkommander/pkg/logger"
	"github.com/IBM/openkommander/pkg/session"
)

func HandleClusterMetadata(w http.ResponseWriter, r *http.Request) {
	currentSession := session.GetCurrentSession()
	client, err := currentSession.GetClient()

	if err != nil {
		logger.Error("Failed to create Kafka client for cluster metadata operation", "cluster", session.GetActiveClusterName(), "error", err)
		SendError(w, "Failed to create Kafka client", err)
		return
	}
	if client == nil {
		logger.Error("Client creation failed for cluster metadata operation", "cluster", session.GetActiveClusterName())
		SendError(w, "Failed to create Kafka client", fmt.Errorf("client creation failed"))
		return
	}

	metadata, failure := commands.GetClusterMetadata()
	if failure != nil {
		logger.Error("Failed to get cluster metadata", "cluster", session.GetActiveClusterName(), "error", failure.Err)
		SendError(w, "Failed to get cluster metadata", failure.Err)
		return
	}

	logger.Info("Successfully retrieved cluster metadata", "cluster", session.GetActiveClusterName())
	SendJSON(w, http.StatusOK, Response{Status: "ok", Data: metadata})
}
