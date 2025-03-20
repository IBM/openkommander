package utils

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

func HandleCLIError(err error, message string) {
	if err != nil {
		log.Fatalf("%s: %v", message, err)
	}
}

func CheckCommandArgCount(cmd *cobra.Command, args []string, expected int, usage string) bool {
	if len(args) != expected {
		fmt.Fprintf(os.Stderr, "Error: expected %d argument(s), got %d\n", expected, len(args))
		fmt.Fprintf(os.Stderr, "Usage: %s\n", usage)
		return false
	}
	return true
}

func HandleAPIError(c *gin.Context, statusCode int, err error, customMessage string) {
	message := err.Error()
	if customMessage != "" {
		message = customMessage + ": " + message
	}
	
	c.JSON(statusCode, gin.H{"error": message})
}

func MustGetString(cmd *cobra.Command, flag string) string {
	value, err := cmd.Flags().GetString(flag)
	if err != nil {
		log.Fatalf("Failed to get flag '%s': %v", flag, err)
	}
	return value
}

func MustGetBool(cmd *cobra.Command, flag string) bool {
	value, err := cmd.Flags().GetBool(flag)
	if err != nil {
		log.Fatalf("Failed to get flag '%s': %v", flag, err)
	}
	return value
}

func MustGetStringSlice(cmd *cobra.Command, flag string) []string {
	value, err := cmd.Flags().GetStringSlice(flag)
	if err != nil {
		log.Fatalf("Failed to get flag '%s': %v", flag, err)
	}
	return value
}

func MustGetInt32(cmd *cobra.Command, flag string) int32 {
	value, err := cmd.Flags().GetInt32(flag)
	if err != nil {
		log.Fatalf("Failed to get flag '%s': %v", flag, err)
	}
	return value
}

func MustGetInt16(cmd *cobra.Command, flag string) int16 {
	value, err := cmd.Flags().GetInt16(flag)
	if err != nil {
		log.Fatalf("Failed to get flag '%s': %v", flag, err)
	}
	return int16(value)
}

type HTTPErrorResponse struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Error      string `json:"error,omitempty"`
}

func RespondWithError(c *gin.Context, statusCode int, message string, err error) {
	errorStr := ""
	if err != nil {
		errorStr = err.Error()
	}
	
	c.JSON(statusCode, HTTPErrorResponse{
		StatusCode: statusCode,
		Message:    message,
		Error:      errorStr,
	})
}

func NotFound(c *gin.Context, resourceType, resourceID string) {
	RespondWithError(c, http.StatusNotFound, fmt.Sprintf("%s '%s' not found", resourceType, resourceID), nil)
}

func BadRequest(c *gin.Context, message string, err error) {
	RespondWithError(c, http.StatusBadRequest, message, err)
}

func ServerError(c *gin.Context, message string, err error) {
	RespondWithError(c, http.StatusInternalServerError, message, err)
}

