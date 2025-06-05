package cluster

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/spf13/viper"
)

type Cluster struct {
	Brokers []string
	Config  *sarama.Config
}

func NewCluster(brokers []string) *Cluster {
	config := sarama.NewConfig()
	defaultVersion := viper.GetString("kafka.version")
	config.Version, _ = sarama.ParseKafkaVersion(defaultVersion)
	return &Cluster{
		Brokers: brokers,
		Config:  config,
	}
}

func (c *Cluster) Connect(ctx context.Context) (sarama.Client, error) {
	client, err := sarama.NewClient(c.Brokers, c.Config)
	if err != nil {
		return nil, fmt.Errorf("error creating sarama client (brokers: %v): %w", c.Brokers, err)
	}
	return client, nil
}

func (c *Cluster) ConnectAdmin(ctx context.Context) (sarama.ClusterAdmin, error) {
	admin, err := sarama.NewClusterAdmin(c.Brokers, c.Config)
	if err != nil {
		return nil, fmt.Errorf("error creating sarama cluster admin (brokers: %v): %w", c.Brokers, err)
	}
	return admin, nil
}
