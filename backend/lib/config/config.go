package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"openkommander/models"
)

var DefaultConfigPath string

func init() {
	homeDir, err := os.UserHomeDir()
	if err == nil {
		DefaultConfigPath = filepath.Join(homeDir, ".config", "openkommander.json")
	}
}

type ConfigError struct {
	Message string
}

func (e ConfigError) Error() string {
	return e.Message
}

func RequiredConfigError() error {
	return ConfigError{Message: "config file path is required"}
}

func LoadServerConfig(path string) (*models.ServerConfig, error) {
	if path == "" {
		if DefaultConfigPath == "" {
			return nil, fmt.Errorf("failed to determine home directory for default config")
		}
		
		path = DefaultConfigPath
		
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return nil, fmt.Errorf("default config file not found at: %s", path)
		}
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	config := models.DefaultServerConfig()
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

func LoadKafkaConfig(path string) (*models.KafkaConfig, error) {
	if path == "" {
		if DefaultConfigPath == "" {
			return nil, fmt.Errorf("failed to determine home directory for default config")
		}
		
		path = DefaultConfigPath
		
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return nil, fmt.Errorf("default config file not found at: %s", path)
		}
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	config := models.DefaultKafkaConfig()
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

func GetConfigError() string {
	return fmt.Sprintf("Configuration file not found. Please create one at %s", DefaultConfigPath)
}

func IsRequiredConfigError(err error) bool {
	if configErr, ok := err.(ConfigError); ok {
		return configErr.Message == RequiredConfigError().Error()
	}
	return false
}

func EnsureConfigDirectoryExists() error {
	if DefaultConfigPath == "" {
		return fmt.Errorf("failed to determine home directory for default config")
	}
	
	configDir := filepath.Dir(DefaultConfigPath)
	
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return fmt.Errorf("failed to create config directory: %w", err)
		}
	}
	
	return nil
}

func WriteDefaultConfig() error {
	if err := EnsureConfigDirectoryExists(); err != nil {
		return err
	}
	
	defaultConfig := models.DefaultServerConfig()
	
	data, err := json.MarshalIndent(defaultConfig, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal default config: %w", err)
	}
	
	if err := os.WriteFile(DefaultConfigPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write default config file: %w", err)
	}
	
	return nil
}