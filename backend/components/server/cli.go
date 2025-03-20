package server

import (
	"fmt"

	"github.com/spf13/cobra"
	"openkommander/lib/server"
	"openkommander/lib/config"
	"openkommander/lib/factory"
	"openkommander/lib/utils"
)

func StartServer(configPath string) error {
	clientFactory := factory.NewClientFactory()
	
	cfg, err := clientFactory.LoadConfig(configPath)
	if err != nil {
		return err
	}

	serverConfig := server.FromServerConfig(cfg)
	srv, err := server.NewServer(serverConfig)
	if err != nil {
		return fmt.Errorf("failed to create server: %w", err)
	}

	if err := srv.Start(); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	srv.WaitForShutdown()
	return nil
}

func NewCommand(clientFactory *factory.ClientFactory) *cobra.Command {
	serverCmd := &cobra.Command{
		Use:   "server",
		Short: "Start the OpenKommander API server",
		Run: func(cmd *cobra.Command, args []string) {
			configPath, _ := cmd.Flags().GetString("config")
			err := StartServer(configPath)
			
			if err != nil {
				if config.IsRequiredConfigError(err) {
					fmt.Println("Error:", config.GetConfigError())
				} else {
					utils.HandleCLIError(err, "Server error")
				}
			}
		},
	}

	return serverCmd
}
