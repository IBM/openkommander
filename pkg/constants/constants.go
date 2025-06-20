package constants

import (
	"github.com/IBM/sarama"
)

var (
	KafkaVersion                           = "4.0.0"
	SaramaKafkaVersion sarama.KafkaVersion = sarama.V4_0_0_0
	KafkaBroker                            = "kafka:9093"
)

func init() {

}
