package kafka

import (
	"context"
	"log"
	"os"
	"github.com/segmentio/kafka-go"
)

var kafkaBroker = os.Getenv("KAFKA_CONSUMER_URL")
var kafkaTopic = "input_stream"

func StartKafkaConsumer(resultChannel chan string) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{kafkaBroker}, 
		Topic:   kafkaTopic,            
		GroupID: "consumer-group-id",   
	})

	defer reader.Close()

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Fatalf("Error reading message from Kafka: %s", err)
		}

		log.Printf("Received message: %s", string(msg.Value))

		resultChannel <- string(msg.Value)
	}
}
