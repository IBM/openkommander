package constants

import (
	"os"
	"path/filepath"

	"github.com/IBM/sarama"
)

var (
	OpenKommanderFolder         string
	OpenKommanderConfigFilename string
	KafkaVersion                                    = "3.9.0"
	SaramaKafkaVersion          sarama.KafkaVersion = sarama.V4_1_0_0
	KafkaBroker                                     = "localhost:9092"
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
