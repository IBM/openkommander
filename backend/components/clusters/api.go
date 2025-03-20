package clusters

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"openkommander/lib/kafka"
	"openkommander/lib/utils"
)

func RegisterRoutes(router *gin.RouterGroup, clients map[string]*kafka.Client) {
	router.GET("/clusters", func(c *gin.Context) {
		clusters := ListClusters(clients)
		c.JSON(http.StatusOK, clusters)
	})
	
	router.GET("/clusters/:name", func(c *gin.Context) {
		name := c.Param("name")
		
		clusterDetails, err := GetCluster(clients, name)
		if err != nil {
			utils.NotFound(c, "Cluster", name)
			return
		}
		
		c.JSON(http.StatusOK, clusterDetails)
	})
}
