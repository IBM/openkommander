package factory

import (
	"fmt"
	"github.com/spf13/cobra"
	
	"openkommander/lib/config"
	"openkommander/lib/kafka"
	"openkommander/models"
)

type ClientFactory struct {
	defaultConfigPath string
}

func NewClientFactory() *ClientFactory {
	return &ClientFactory{
		defaultConfigPath: config.DefaultConfigPath,
	}
}

func (f *ClientFactory) LoadConfig(configPath string) (*models.ServerConfig, error) {
	return config.LoadServerConfig(configPath)
}

func (f *ClientFactory) CreateClientFromConfig(cfg *models.ServerConfig, clusterName string) (*kafka.Client, error) {
	return kafka.GetClientFromConfig(cfg, clusterName)
}

func (f *ClientFactory) CreateClientFromFlags(cmd *cobra.Command) (*kafka.Client, error) {
	configPath, _ := cmd.Flags().GetString("config")
	clusterName, _ := cmd.Flags().GetString("cluster")
	
	cfg, err := config.LoadServerConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	
	return f.CreateClientFromConfig(cfg, clusterName)
}

func (f *ClientFactory) CreateMultiClusterClients(configPath string) (map[string]*kafka.Client, error) {
	cfg, err := config.LoadServerConfig(configPath) 
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	
	clients := make(map[string]*kafka.Client)
	
	for _, cluster := range cfg.Clusters {
		client, err := kafka.CreateClientFromKafkaConfig(cluster.KafkaConfig)
		if err != nil {
			fmt.Printf("Warning: Failed to create client for cluster '%s': %v\n", cluster.Name, err)
			continue
		}
		
		clients[cluster.Name] = client
	}
	
	return clients, nil
}

func (f *ClientFactory) DefaultClient(configPath string) (*kafka.Client, error) {
	cfg, err := config.LoadServerConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	
	return kafka.CreateClientFromKafkaConfig(cfg.KafkaConfig)
}