package rest

import (
	"net/http"
	"time"

	"github.com/IBM/openkommander/pkg/logger"
	"github.com/IBM/openkommander/pkg/session"
)

func HandleStatus(w http.ResponseWriter, r *http.Request, startTime time.Time) {
	clusters := session.GetClusterConnections()
	kafkaStatus := "disconnected"
	for _, cluster := range clusters {
		if cluster.IsAuthenticated {
			kafkaStatus = "connected"
			break
		}
	}

	uptime := time.Since(startTime).Seconds()
	logger.Info("Status check completed", "kafka_status", kafkaStatus, "clusters_count", len(clusters), "uptime_seconds", uptime)

	response := Response{
		Status:  "ok",
		Message: "OpenKommander REST API is running",
		Data: map[string]interface{}{
			"status":         "running",
			"kafka_status":   kafkaStatus,
			"clusters_count": len(clusters),
			"uptime_seconds": uptime,
		},
	}
	SendJSON(w, http.StatusOK, response)
}
