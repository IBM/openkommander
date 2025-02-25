package main

import (
	"fmt"
	"os"

	"github.com/IBM/openkommander/pkg/commands"
	"github.com/spf13/cobra"
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
}
