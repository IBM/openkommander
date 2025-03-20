package kafka

import (
	"io"
	"os"
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/IBM/sarama"
)

type consumerGroupHandler struct {
	handler func(msg *sarama.ConsumerMessage) error
}

func (h *consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h *consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		if err := h.handler(msg); err != nil {
			return err
		}
		session.MarkMessage(msg, "")
	}
	return nil
}

func (c *Client) ProduceMessage(topic, key string, value interface{}) error {
	producer, err := sarama.NewSyncProducer(c.brokers, c.config)
	if err != nil {
		return fmt.Errorf("failed to create producer: %w", err)
	}
	defer producer.Close()

	var msgValue []byte
	switch v := value.(type) {
	case string:
		msgValue = []byte(v)
	case []byte:
		msgValue = v
	default:
		msgValue, err = json.Marshal(value)
		if err != nil {
			return fmt.Errorf("failed to marshal message value: %w", err)
		}
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(msgValue),
	}

	if key != "" {
		msg.Key = sarama.StringEncoder(key)
	}

	_, _, err = producer.SendMessage(msg)
	return err
}

func (c *Client) ProduceMessageFromReader(topic, key string, reader io.Reader, isJSON bool) error {
	body, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	var msgValue interface{}
	if isJSON {
		if err := json.Unmarshal(body, &msgValue); err != nil {
			return err
		}
	} else {
		msgValue = string(body)
	}

	return c.ProduceMessage(topic, key, msgValue)
}

func (c *Client)  ProduceMessageFromFile(topic, key, filename string, isJSON bool) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return c.ProduceMessageFromReader(topic, key, file, isJSON)
}


func (c *Client) ConsumeMessages(ctx context.Context, topic string, group string, handler func(msg *sarama.ConsumerMessage) error) error {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	consumer, err := sarama.NewConsumerGroup(c.brokers, group, config)
	if err != nil {
		return fmt.Errorf("failed to create consumer group: %w", err)
	}
	defer consumer.Close()

	msgHandler := &consumerGroupHandler{
		handler: handler,
	}

	for {
		if err := consumer.Consume(ctx, []string{topic}, msgHandler); err != nil {
			return err
		}

		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
}

func (c *Client) ConsumeMessagesWithOptions(ctx context.Context, topic string, group string, initialOffset int64, handler func(msg *sarama.ConsumerMessage) error) error {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = initialOffset
	
	if group == "" {
		return c.consumeWithoutGroup(ctx, topic, config, handler)
	}
	
	// otherwise use defined consumer group 
	consumer, err := sarama.NewConsumerGroup(c.brokers, group, config)
	if err != nil {
		return fmt.Errorf("failed to create consumer group: %w", err)
	}
	defer consumer.Close()

	msgHandler := &consumerGroupHandler{
		handler: handler,
	}

	for {
		if err := consumer.Consume(ctx, []string{topic}, msgHandler); err != nil {
			return err
		}

		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
}

func (c *Client) consumeWithoutGroup(ctx context.Context, topic string, config *sarama.Config, handler func(msg *sarama.ConsumerMessage) error) error {
	consumer, err := sarama.NewConsumer(c.brokers, config)
	if err != nil {
		return fmt.Errorf("failed to create consumer: %w", err)
	}
	defer consumer.Close()

	partitions, err := consumer.Partitions(topic)
	if err != nil {
		return fmt.Errorf("failed to get partitions: %w", err)
	}

	var wg sync.WaitGroup
	
	errors := make(chan error, len(partitions))

	for _, partition := range partitions {
		wg.Add(1)
		
		go func(partition int32) {
			defer wg.Done()
			
			partitionConsumer, err := consumer.ConsumePartition(topic, partition, config.Consumer.Offsets.Initial)
			if err != nil {
				errors <- fmt.Errorf("failed to create partition consumer: %w", err)
				return
			}
			defer partitionConsumer.Close()
			
			for {
				select {
				case msg := <-partitionConsumer.Messages():
					if err := handler(msg); err != nil {
						errors <- err
						return
					}
				case err := <-partitionConsumer.Errors():
					errors <- err
					return
				case <-ctx.Done():
					return
				}
			}
		}(partition)
	}

	select {
	case err := <-errors:
		return err
	case <-ctx.Done():
		wg.Wait()
		return ctx.Err()
	}
}
