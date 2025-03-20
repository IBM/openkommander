package server

import (
	"github.com/IBM/sarama"
	"openkommander/lib/kafka"
	"openkommander/models"
)

type Config struct {
	Port          string
	Brokers       []string
	SASLEnabled   bool
	SASLUsername  string
	SASLPassword  string
	SASLMechanism string
	TLSEnabled    bool
	LogLevel      string
	Clusters      []models.ClusterConfig
}

func (c *Config) KafkaClientOptions() []kafka.ClientOption {
	var opts []kafka.ClientOption

	if c.SASLEnabled {
		var mechanism sarama.SASLMechanism
		switch c.SASLMechanism {
		case "SCRAM-SHA-256":
			mechanism = sarama.SASLTypeSCRAMSHA256
		case "SCRAM-SHA-512":
			mechanism = sarama.SASLTypeSCRAMSHA512
		default:
			mechanism = sarama.SASLTypePlaintext
		}
		opts = append(opts, kafka.WithSASL(c.SASLUsername, c.SASLPassword, mechanism))
	}

	if c.TLSEnabled {
		opts = append(opts, kafka.WithTLS())
	}

	return opts
}

func FromServerConfig(cfg *models.ServerConfig) *Config {
	return &Config{
		Port:          cfg.Port,
		Brokers:       cfg.Brokers,
		SASLEnabled:   cfg.SASLEnabled,
		SASLUsername:  cfg.SASLUsername,
		SASLPassword:  cfg.SASLPassword,
		SASLMechanism: cfg.SASLMechanism,
		TLSEnabled:    cfg.TLSEnabled,
		LogLevel:      cfg.LogLevel,
		Clusters:      cfg.Clusters,
	}
}

