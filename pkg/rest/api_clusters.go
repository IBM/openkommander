package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/IBM/openkommander/pkg/constants"
	"github.com/IBM/openkommander/pkg/logger"
	"github.com/IBM/openkommander/pkg/session"
)

func HandleListClusters(w http.ResponseWriter, r *http.Request) {
	clusters := session.GetClusterConnections()
	activeCluster := session.GetActiveClusterName()

	if len(clusters) == 0 {
		// Return an empty slice of maps when there are no clusters
		SendJSON(w, http.StatusOK, Response{Status: "ok", Data: []map[string]interface{}{}})
		return
	}

	rows := make([]map[string]interface{}, 0, len(clusters))

	for i, cluster := range clusters {
		status := "Disconnected"
		if cluster.IsAuthenticated {
			status = "Connected"
		}

		active := "No"
		if cluster.Name == activeCluster {
			active = "Yes"
		}

		rows = append(rows, map[string]interface{}{
			"id":      i + 1,
			"name":    cluster.Name,
			"brokers": cluster.Brokers,
			"status":  status,
			"active":  active,
		})
	}

	logger.Info("Successfully retrieved clusters", "cluster_count", len(clusters))
	SendJSON(w, http.StatusOK, Response{Status: "ok", Data: rows})
}

func HandleLoginCluster(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name    string `json:"name"`
		Broker  string `json:"broker"`
		Version string `json:"version"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		logger.Error("Invalid request body for cluster login", "error", err)
		SendError(w, "Invalid request body", err)
		return
	}

	err = session.InitAPI()
	if err != nil {
		logger.Error("Failed to initialize session for cluster login", "error", err)
		SendError(w, "Failed to initialize session", err)
		return
	}

	brokers := []string{req.Broker}

	success, message := session.LoginWithParams(brokers, constants.KafkaVersion, req.Name)
	if !success {
		logger.Error("Failed to login to cluster", "cluster_name", req.Name, "message", message)
		SendError(w, "Failed to login to cluster: "+message, nil)
		return
	}

	logger.Info("Successfully logged in to cluster", "cluster_name", req.Name)
	SendJSON(w, http.StatusOK, Response{Status: "ok", Message: fmt.Sprintf("Logged in to cluster '%s' successfully", req.Name)})
}
