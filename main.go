package main

import (
	"os"

	"github.com/IBM/openkommander/pkg/cli"
	"github.com/IBM/openkommander/pkg/logger"
)

func main() {
	// Initialize logger with pretty formatting for development
	config := &logger.Config{
		Level:     logger.LevelDebug,
		Format:    "pretty",
		AddColors: true,
	}
	logger.Init(config)

	var rootCmd = cli.Init()

	if err := rootCmd.Execute(); err != nil {
		logger.Error("Command execution failed", "error")
		os.Exit(1)
	}
}
