package connect

import (
	"fmt"

	"github.com/spf13/cobra"
	"openkommander/lib/config"
)

func NewCommand() *cobra.Command {
	connectCmd := &cobra.Command{
		Use:   "connect",
		Short: "Connect to your Kafka cluster by initializing default configuration file",
		Long:  `Creates a default configuration file at $HOME/.config/openkommander.json`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Creating default configuration file at %s\n", config.DefaultConfigPath)
			
			if err := config.WriteDefaultConfig(); err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
			
			fmt.Println("Default configuration file created successfully.")
			fmt.Println("You should modify this file with your Kafka connection details.")
		},
	}

	return connectCmd
}