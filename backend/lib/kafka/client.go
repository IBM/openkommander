package kafka

import (
	"fmt"

	"github.com/IBM/sarama"

	"openkommander/lib/utils"
	"openkommander/models"
)

type Client struct {
	config  *sarama.Config
	brokers []string
	admin   sarama.ClusterAdmin
}

func NewClient(brokers []string, opts ...ClientOption) (*Client, error) {
	config := sarama.NewConfig()
	config.Version = sarama.V2_8_0_0

	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true
	config.Consumer.Return.Errors = true

	client := &Client{
		config:  config,
		brokers: brokers,
	}

	for _, opt := range opts {
		opt(client)
	}

	admin, err := sarama.NewClusterAdmin(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create cluster admin: %w", err)
	}
	client.admin = admin

	return client, nil
}

func (c *Client) GetBrokers() []string {
	return c.brokers
}

type ClientOption func(*Client)

func WithSASL(username, password string, mechanism sarama.SASLMechanism) ClientOption {
	return func(c *Client) {
		c.config.Net.SASL.Enable = true
		c.config.Net.SASL.Mechanism = mechanism
		c.config.Net.SASL.User = username
		c.config.Net.SASL.Password = password
	}
}

func WithTLS() ClientOption {
	return func(c *Client) {
		c.config.Net.TLS.Enable = true
	}
}

func WithVersion(version string) ClientOption {
	return func(c *Client) {
		if v, err := sarama.ParseKafkaVersion(version); err == nil {
			c.config.Version = v
		}
	}
}

func (c *Client) Close() error {
	if c.admin != nil {
		return c.admin.Close()
	}
	return nil
}

func GetKafkaOptionsFromConfig(cfg models.KafkaConfig) []ClientOption {
	var opts []ClientOption

	if cfg.SASLEnabled {
		mechanism := utils.GetSASLMechanism(cfg.SASLMechanism)
		opts = append(opts, WithSASL(cfg.SASLUsername, cfg.SASLPassword, mechanism))
	}

	if cfg.TLSEnabled {
		opts = append(opts, WithTLS())
	}

	return opts
}

func CreateClientFromKafkaConfig(cfg models.KafkaConfig) (*Client, error) {
	opts := GetKafkaOptionsFromConfig(cfg)
	return NewClient(cfg.Brokers, opts...)
}

func GetClientFromConfig(cfg *models.ServerConfig, clusterName string) (*Client, error) {
	if clusterName != "" {
		var clusterCfg *models.ClusterConfig
		for i := range cfg.Clusters {
			if cfg.Clusters[i].Name == clusterName {
				clusterCfg = &cfg.Clusters[i]
				break
			}
		}

		if clusterCfg == nil {
			return nil, fmt.Errorf("cluster '%s' not found in config", clusterName)
		}

		return CreateClientFromKafkaConfig(clusterCfg.KafkaConfig)
	}

	return CreateClientFromKafkaConfig(cfg.KafkaConfig)
}

func (c *Client) Test() error {
	_, err := c.admin.ListTopics()
	if err != nil {
		return fmt.Errorf("failed to connect to Kafka cluster: %w", err)
	}
	return nil
}