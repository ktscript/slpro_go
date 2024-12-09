package kafka

import (
	"log"
	"os"
	"github.com/segmentio/kafka-go"
)

var kafkaBroker = os.Getenv("KAFKA_CONSUMER_URL")
var kafkaTopic = "output_stream"

func StartKafkaProducer(resultChannel chan string) (*kafka.Writer, error) {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(kafkaBroker), 
		Topic:    kafkaTopic,              
		Balancer: &kafka.LeastBytes{},         
	}

	go func() {
		for result := range resultChannel {
			msg := kafka.Message{
				Value: []byte(result),
			}
			err := writer.WriteMessages(nil, msg)
			if err != nil {
				log.Printf("Failed to write message: %v", err)
			} else {
				log.Printf("Successfully sent message: %s", result)
			}
		}
	}()

	return writer, nil
}
