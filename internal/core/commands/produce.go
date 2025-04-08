package commands

import (
	"fmt"
	"log"
	"net/http"

	"github.com/IBM/openkommander/pkg/session"
	"github.com/IBM/sarama"
)

func ProduceMessage(topicName, key, msg string, partition, acks int) (successMessage string, f *Failure) {
	// validate the session is open
	_, validateFailure := GetClient()
	if validateFailure != nil {
		return "", validateFailure
	}

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.RequiredAcks(acks)
	config.Producer.Return.Successes = true

	if partition >= 0 {
		config.Producer.Partitioner = sarama.NewManualPartitioner
	} else {
		config.Producer.Partitioner = sarama.NewRandomPartitioner
	}

	message := &sarama.ProducerMessage{Topic: topicName, Partition: int32(partition)}

	if key != "" {
		message.Key = sarama.StringEncoder(key)
	}

	message.Value = sarama.StringEncoder(msg)

	producer, err := sarama.NewSyncProducer(session.GetCurrentSession().GetBrokers(), config)
	if err != nil {
		return "", NewFailure("Failed to open Kafka producer", http.StatusBadRequest)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			log.Println("Failed to close Kafka producer cleanly:", err)
		}
	}()

	part, offset, err := producer.SendMessage(message)
	if err != nil {
		return "", NewFailure(fmt.Sprintf("Failed to produce message: %s", err), http.StatusBadRequest)
	}

	return fmt.Sprintf("successfully written to partition %d with offset %d", part, offset), nil
}
