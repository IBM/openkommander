package constants

import (
	"github.com/IBM/sarama"
)

var (
	OpenKommanderConfigFilename                     = ".openkommander_config"
	KafkaVersion                                    = "3.9.0"
	SaramaKafkaVersion          sarama.KafkaVersion = sarama.V3_9_0_0
	KafkaBroker                                     = "kafka:9093"
)

func init() {

}
