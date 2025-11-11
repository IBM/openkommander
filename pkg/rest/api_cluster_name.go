package rest

import (
	"fmt"
	"net/http"

	"github.com/IBM/openkommander/pkg/logger"
	"github.com/IBM/openkommander/pkg/session"
)

func HandleGetClusterByName(w http.ResponseWriter, r *http.Request, clusterName string) {
	err := session.InitAPI()
	if err != nil {
		logger.Error("Failed to initialize session for cluster by name", "error", err)
		SendError(w, "Failed to initialize session", err)
		return
	}

	cluster := session.GetClusterByName(clusterName)
	if cluster == nil {
		logger.Warn("Cluster not found", "cluster_name", clusterName)
		SendError(w, "Cluster not found", fmt.Errorf("cluster '%s' not found", clusterName))
		return
	}
	logger.Info("Successfully retrieved cluster details", "cluster_name", clusterName)
	SendJSON(w, http.StatusOK, Response{Status: "ok", Data: cluster})
}

func HandleDeleteClusterByName(w http.ResponseWriter, r *http.Request, clusterName string) {
	err := session.InitAPI()
	if err != nil {
		logger.Error("Failed to initialize session for cluster deletion", "error", err)
		SendError(w, "Failed to initialize session", err)
		return
	}

	success := session.Logout(clusterName)
	if !success {
		logger.Warn("Failed to logout from cluster", "cluster_name", clusterName)
		SendError(w, "Failed to logout from cluster", fmt.Errorf("failed to logout from cluster '%s'", clusterName))
		return
	}
	logger.Info("Successfully logged out from cluster", "cluster_name", clusterName)
	SendJSON(w, http.StatusOK, Response{Status: "ok", Message: fmt.Sprintf("Logged out from cluster '%s' successfully", clusterName)})
}

func HandleSelectCluster(w http.ResponseWriter, r *http.Request, clusterName string) {
	err := session.InitAPI()
	if err != nil {
		logger.Error("Failed to initialize session for select cluster", "error", err)
		SendError(w, "Failed to initialize session", err)
		return
	}

	session.SelectCluster(clusterName)
	logger.Info("Successfully set active cluster", "cluster_name", clusterName)
	SendJSON(w, http.StatusOK, Response{Status: "ok", Message: fmt.Sprintf("Active cluster set to '%s'", clusterName)})
}
