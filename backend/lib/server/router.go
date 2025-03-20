package server

import (
	"github.com/gin-gonic/gin"

	"openkommander/lib/kafka"
	
	"openkommander/components/brokers"
	"openkommander/components/clusters"
	"openkommander/components/consumers"
	"openkommander/components/health"
	"openkommander/components/messages"
	"openkommander/components/topics"
)

func SetupRoutes(router *gin.Engine, client *kafka.Client, consumerTracker *kafka.ConsumerTracker) {
	RegisterStaticRoutes(router)
	
	v1 := router.Group("/api/v1")
	
	health.RegisterRoutes(v1)
	brokers.RegisterRoutes(v1, client)
	topics.RegisterRoutes(v1, client)
	messages.RegisterRoutes(v1, client)
	consumers.RegisterRoutes(v1, client, consumerTracker)
}

func SetupRoutesWithClusters(
	router *gin.Engine, 
	client *kafka.Client, 
	clusterClients map[string]*kafka.Client, 
	consumerTracker *kafka.ConsumerTracker,
) {
	RegisterStaticRoutes(router)
	
	v1 := router.Group("/api/v1")
	
	health.RegisterRoutes(v1)
	brokers.RegisterRoutes(v1, client)
	topics.RegisterRoutes(v1, client)
	messages.RegisterRoutes(v1, client)
	consumers.RegisterRoutes(v1, client, consumerTracker)
	
	if len(clusterClients) > 0 {
		clusters.RegisterRoutes(v1, clusterClients)
		
		for name, clusterClient := range clusterClients {
			clusterGroup := v1.Group("/clusters/" + name)
			
			brokers.RegisterRoutes(clusterGroup, clusterClient)
			topics.RegisterRoutes(clusterGroup, clusterClient)
			messages.RegisterRoutes(clusterGroup, clusterClient)
			consumers.RegisterRoutes(clusterGroup, clusterClient, consumerTracker)
		}
	}
}
