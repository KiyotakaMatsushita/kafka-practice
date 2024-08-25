package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/IBM/sarama"
)

var (
	logger = log.New(os.Stdout, "[KafkaConsumer] ", log.LstdFlags)
)

func Run(brokers, version, topic string) {
	if brokers == "" {
		logger.Fatalln("At least one broker is required")
	}
	splitBrokers := strings.Split(brokers, ",")

	kafkaVersion, err := sarama.ParseKafkaVersion(version)
	if err != nil {
		log.Panicf("Error parsing Kafka version: %v", err)
	}

	config := sarama.NewConfig()
	config.Version = kafkaVersion
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()

	consumer := Consumer{
		ready: make(chan bool),
	}

	ctx, cancel := context.WithCancel(context.Background())
	client, err := sarama.NewConsumerGroup(splitBrokers, "default_group", config)
	if err != nil {
		log.Panicf("Error creating consumer group: %v", err)
	}

	go func() {
		for {
			if err := client.Consume(ctx, []string{topic}, &consumer); err != nil {
				log.Panicf("Error from consumer: %v", err)
			}
			if ctx.Err() != nil {
				return
			}
			consumer.ready = make(chan bool)
		}
	}()

	<-consumer.ready
	logger.Println("Consumer is ready")

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, os.Interrupt)
	<-sigterm

	logger.Println("Received termination signal. Initiating shutdown...")
	cancel()
	if err = client.Close(); err != nil {
		log.Panicf("Error closing client: %v", err)
	}
}

type Consumer struct {
	ready chan bool
}

func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	close(consumer.ready)
	return nil
}

func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		err := consumer.processMessage(message)
		if err != nil {
			logger.Printf("Error processing message: %v", err)
		}
		session.MarkMessage(message, "")
	}
	return nil
}

func (consumer *Consumer) processMessage(message *sarama.ConsumerMessage) error {
	logger.Printf("Received message: topic=%v, partition=%v, offset=%v, key=%s, value=%s\n",
		message.Topic, message.Partition, message.Offset, string(message.Key), string(message.Value))

	var data map[string]interface{}
	err := json.Unmarshal(message.Value, &data)
	if err != nil {
		return fmt.Errorf("error decoding JSON: %v", err)
	}

	err = saveToDatabase(data)
	if err != nil {
		return fmt.Errorf("error saving to database: %v", err)
	}

	return nil
}

func saveToDatabase(data map[string]interface{}) error {
	logger.Printf("Saving to database: %v", data)
	return nil
}