package topics

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"openkommander/lib/kafka"
	"openkommander/models"
	"openkommander/lib/utils"
)

func RegisterRoutes(router *gin.RouterGroup, client *kafka.Client) {
	topics := router.Group("/topics")
	{
		topics.GET("", func(c *gin.Context) {
			topicInfo, err := client.GetTopicInfo()
			if err != nil {
				utils.ServerError(c, "Failed to list topics", err)
				return
			}
			c.JSON(http.StatusOK, topicInfo)
		})

		topics.POST("", func(c *gin.Context) {
			var req models.TopicCreateRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				utils.BadRequest(c, "Invalid request", err)
				return
			}

			if err := client.CreateTopic(req.Name, req.Partitions, req.ReplicationFactor); err != nil {
				utils.ServerError(c, "Failed to create topic", err)
				return
			}

			c.JSON(http.StatusCreated, gin.H{"status": "created"})
		})

		topics.GET("/:name", func(c *gin.Context) {
			name := c.Param("name")
			
			detail, err := client.GetTopic(name)
			if err != nil {
				utils.NotFound(c, "Topic", name)
				return
			}

			c.JSON(http.StatusOK, detail)
		})

		topics.DELETE("/:name", func(c *gin.Context) {
			name := c.Param("name")
			if err := client.DeleteTopic(name); err != nil {
				utils.ServerError(c, "Failed to delete topic", err)
				return
			}

			c.JSON(http.StatusOK, gin.H{"status": "deleted"})
		})
	}
}
