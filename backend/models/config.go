package models

type KafkaConfig struct {
	Brokers       []string `json:"brokers"`
	SASLEnabled   bool     `json:"sasl_enabled"`
	SASLUsername  string   `json:"sasl_username,omitempty"`
	SASLPassword  string   `json:"sasl_password,omitempty"`
	SASLMechanism string   `json:"sasl_mechanism,omitempty"`
	TLSEnabled    bool     `json:"tls_enabled"`
}

type ClusterConfig struct {
	KafkaConfig
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type ServerConfig struct {
	KafkaConfig
	Port       string          `json:"port"`
	LogLevel   string          `json:"log_level"`
	MetricsURL string          `json:"metrics_url,omitempty"`
	Clusters   []ClusterConfig `json:"clusters,omitempty"`
}

func DefaultKafkaConfig() KafkaConfig {
	return KafkaConfig{
		Brokers:       []string{"localhost:9092"},
		SASLEnabled:   false,
		SASLMechanism: "PLAIN",
		TLSEnabled:    false,
	}
}

func DefaultServerConfig() ServerConfig {
	return ServerConfig{
		KafkaConfig: DefaultKafkaConfig(),
		Port:        "8080",
		LogLevel:    "info",
		Clusters:    []ClusterConfig{},
	}
}

