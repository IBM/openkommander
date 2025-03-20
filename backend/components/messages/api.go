package messages

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"openkommander/lib/kafka"
	"openkommander/lib/utils"
)

func RegisterRoutes(router *gin.RouterGroup, client *kafka.Client) {
	router.POST("/messages/:topic", func(c *gin.Context) {
		topic := c.Param("topic")
		key := c.Query("key")

		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			utils.BadRequest(c, "Failed to read request body", err)
			return
		}

		var msgValue interface{}
		contentType := c.GetHeader("Content-Type")

		if strings.Contains(contentType, "application/json") {
			if err := json.Unmarshal(body, &msgValue); err != nil {
				utils.BadRequest(c, "Invalid JSON format", err)
				return
			}
		} else {
			msgValue = string(body)
		}

		if err := client.ProduceMessage(topic, key, msgValue); err != nil {
			utils.ServerError(c, "Failed to produce message", err)
			return
		}

		c.JSON(http.StatusAccepted, gin.H{"status": "message sent"})
	})
}
