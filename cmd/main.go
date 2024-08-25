package main

import (
	"flag"
	"kafka-sarama-example/internal/consumer"
	"os"
	"sync"

	"github.com/IBM/sarama"
)

var (
	brokers = flag.String("brokers", os.Getenv("KAFKA_BROKERS"), "The Kafka brokers to connect to, as a comma separated list")
	version = flag.String("version", sarama.V3_6_0_0.String(), "Kafka cluster version")
	topic   = flag.String("topic", "default_topic", "The Kafka topic to use")
)

func main() {
	flag.Parse()

	var wg sync.WaitGroup
	wg.Add(1)

	// go func() {
	// 	defer wg.Done()
	// 		producer.Run(*brokers, *version, *topic)
	// }()

	go func() {
		defer wg.Done()
		consumer.Run(*brokers, *version, *topic)
	}()

	wg.Wait()
}