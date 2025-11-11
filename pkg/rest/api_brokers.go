package rest

import (
	"fmt"
	"net/http"

	"github.com/IBM/openkommander/pkg/logger"
	"github.com/IBM/openkommander/pkg/session"
)

func HandleCreateBroker(w http.ResponseWriter, r *http.Request) {
	_ = r // Not used in this stub implementation
	response := Response{
		Status:  "ok",
		Message: "Broker creation is not implemented yet",
	}

	SendJSON(w, http.StatusNotImplemented, response)
}

func HandleGetBrokers(w http.ResponseWriter, r *http.Request) {
	currentSession := session.GetCurrentSession()
	client, err := currentSession.GetClient()
	if err != nil {
		logger.Error("Failed to create Kafka client for brokers operation", "cluster", session.GetActiveClusterName(), "error", err)
		SendError(w, "Failed to create Kafka client", err)
		return
	}
	if client == nil {
		logger.Error("Client creation failed for brokers operation", "cluster", session.GetActiveClusterName())
		SendError(w, "Failed to create Kafka client", fmt.Errorf("client creation failed"))
		return
	}

	brokers := client.Brokers()
	brokerList := make([]map[string]interface{}, 0)

	for _, brokerInfo := range brokers {
		connected, err := brokerInfo.Connected()
		if err != nil {
			connected = false
		}

		tlsState, _ := brokerInfo.TLSConnectionState()

		brokerData := map[string]interface{}{
			"id":        brokerInfo.ID(),
			"addr":      brokerInfo.Addr(),
			"connected": connected,
			"rack":      brokerInfo.Rack(),
			"state":     tlsState,
		}
		brokerList = append(brokerList, brokerData)
	}

	logger.Info("Successfully retrieved brokers", "cluster", session.GetActiveClusterName(), "broker_count", len(brokerList))
	SendJSON(w, http.StatusOK, Response{Status: "ok", Data: brokerList})
}
