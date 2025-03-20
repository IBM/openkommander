package brokers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"openkommander/lib/kafka"
	"openkommander/lib/utils"
)

func RegisterRoutes(router *gin.RouterGroup, client *kafka.Client) {
	router.GET("/brokers", func(c *gin.Context) {
		brokers, err := client.GetBrokerInfo()
		if err != nil {
			utils.ServerError(c, "Failed to list brokers", err)
			return
		}

		c.JSON(http.StatusOK, brokers)
	})
}
