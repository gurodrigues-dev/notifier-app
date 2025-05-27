package queue

import (
	"context"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaImpl struct {
	brokers []string
}

func NewKafkaImpl(brokers []string) *KafkaImpl {
	return &KafkaImpl{
		brokers: brokers,
	}
}

func (k *KafkaImpl) Produce(topic, message string) error {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  k.brokers,
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	})

	defer writer.Close()

	return writer.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(time.Now().Format(time.RFC3339)),
		Value: []byte(message),
	})
}

func (k *KafkaImpl) Consumer(topic, group string, handler func(message string)) error {

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        k.brokers,
		GroupID:        group,
		Topic:          topic,
		MinBytes:       10e3,
		MaxBytes:       10e6,
		CommitInterval: time.Second,
	})

	go func() {
		log.Println("Kafka consumer started...")
		defer reader.Close()

		for {
			m, err := reader.ReadMessage(context.Background())
			if err != nil {
				log.Printf("Kafka read error: %v", err)
				continue
			}
			log.Println("Message received")
			handler(string(m.Value))
		}
	}()

	return nil
}
