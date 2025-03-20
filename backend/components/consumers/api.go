package consumers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"openkommander/lib/kafka"
	"openkommander/lib/utils"
)

func RegisterRoutes(router *gin.RouterGroup, client *kafka.Client, consumerTracker *kafka.ConsumerTracker) {
	consumers := router.Group("/consumers")
	{
		
		consumers.DELETE("/:id", func(c *gin.Context) {
			id := c.Param("id")
			if removed := consumerTracker.RemoveConsumer(id); !removed {
				utils.NotFound(c, "Consumer", id)
				return
			}

			c.JSON(http.StatusOK, gin.H{"status": "stopped"})
		})

		consumers.GET("", func(c *gin.Context) {
			groups, err := client.GetConsumerGroups()
			if err != nil {
				utils.ServerError(c, "Failed to list consumer groups", err)
				return
			}

			c.JSON(http.StatusOK, groups)
		})

		consumers.GET("/:group", func(c *gin.Context) {
			groupID := c.Param("group")
			group, err := client.GetConsumerGroup(groupID)
			if err != nil {
				utils.NotFound(c, "Consumer group", groupID)
				return
			}

			c.JSON(http.StatusOK, group)
		})
	}
}
