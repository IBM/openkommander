package constants

import (
	"github.com/IBM/sarama"
	"os"
	"path/filepath"
)

var (
	OpenKommanderFolder         string
	OpenKommanderConfigFilename                     string
	KafkaVersion                                    = "3.9.0"
	SaramaKafkaVersion          sarama.KafkaVersion = sarama.V3_9_0_0
	KafkaBroker                                     = "kafka:9093"
)

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Fallback to current directory if home directory can't be determined
		homeDir = "."
	}
	
	OpenKommanderFolder = filepath.Join(homeDir, ".ok")
	OpenKommanderConfigFilename = filepath.Join(homeDir, ".ok", ".ok_config")
}
