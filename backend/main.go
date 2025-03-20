package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	
	"openkommander/components/brokers"
	"openkommander/components/clusters"
	"openkommander/components/consumers"
	"openkommander/components/connect"
	"openkommander/components/messages"
	"openkommander/components/server"
	"openkommander/components/topics"
	"openkommander/lib/config"
	"openkommander/lib/factory"
)

func main() {
	clientFactory := factory.NewClientFactory()

	rootCmd := &cobra.Command{
		Use:   "ok",
		Short: "OpenKommander",
		Long:  "A CLI tool for managing Apache Kafka clusters, topics, consumers, and producers.",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if cmd.Use == "connect" {
				return
			}
			
			if config.DefaultConfigPath != "" {
				//cmd.Printf("Using configuration %s \n", config.DefaultConfigPath)
			}
		},
	}

	rootCmd.PersistentFlags().String("config", "", "Path to config file (defaults to $HOME/.config/openkommander.json)")
	rootCmd.PersistentFlags().String("cluster", "", "Use named cluster from config")
	
	rootCmd.AddCommand(
		server.NewCommand(clientFactory),
		topics.NewCommand(clientFactory),
		consumers.NewCommand(clientFactory),
		brokers.NewCommand(clientFactory),
		messages.NewCommand(clientFactory),
		clusters.NewCommand(clientFactory),
		connect.NewCommand(),
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}