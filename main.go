package main

import (
	"context"
	"fmt"
	"github.com/IBM/openkommander/pkg/api"
	"github.com/IBM/openkommander/pkg/commands"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "ok",
		Short: "OpenKommander - A CLI tool for Apache Kafka management",
		Long: `OpenKommander is a command line utility for Apache Kafka compatible brokers.
Complete documentation is available at https://github.com/IBM/openkommander`,
		Aliases: []string{"openkommander", "kommander", "okm"},
	}

	commands.RegisterCommands(rootCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	router := api.NewRouter()
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Error during server shutdown: %v", err)
		}
	}()

	log.Println("REST API server running on port 8080...")
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}
}
