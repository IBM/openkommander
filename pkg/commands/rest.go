package commands

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/IBM/openkommander/pkg/rest"
	"github.com/spf13/cobra"
)

var restCmd = &cobra.Command{
	Use:   "rest",
	Short: "Start the REST API server",
	Long: `Starts the OpenKommander REST API server for Kafka management.
Example: ok rest --port 8080`,
	Run: handleRestServer,
}

func init() {
	restCmd.Flags().StringP("port", "p", "8080", "Port to run the REST server on")
	restCmd.Flags().StringSliceP("brokers", "b", []string{"localhost:9092"}, "Kafka broker list")
}

func handleRestServer(cmd *cobra.Command, args []string) {
	port, _ := cmd.Flags().GetString("port")
	brokers, _ := cmd.Flags().GetStringSlice("brokers")
	
	// Create a new server instance
	srv, err := rest.NewServer(port, brokers)
	if err != nil {
		fmt.Printf("Failed to create server: %v\n", err)
		os.Exit(1)
	}

	// Create a channel to listen for interrupt signals
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Start the server in a goroutine
	go func() {
		if err := srv.Start(); err != nil {
			fmt.Printf("Server error: %v\n", err)
			os.Exit(1)
		}
	}()

	fmt.Printf("REST server started on port %s\n", port)
	fmt.Println("\nAvailable endpoints:")
	fmt.Println("  GET    /api/v1/status        - Check server status")
	fmt.Println("  GET    /api/v1/topics        - List all topics")
	fmt.Println("  POST   /api/v1/topics        - Create a new topic")
	fmt.Println("  DELETE /api/v1/topics?name=X - Delete a topic")
	fmt.Println("\nPress Ctrl+C to stop the server")

	// Wait for interrupt signal
	<-done
	fmt.Println("\nShutting down server...")

	// Create a deadline for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := srv.Stop(ctx); err != nil {
		fmt.Printf("Server forced to shutdown: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Server stopped gracefully")
}