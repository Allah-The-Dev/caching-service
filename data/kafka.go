package data

import (
	"encoding/json"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

func (emp *Employee) PublishToKafka() {

	CLogger.Println("publishing emp data to kafka topic")

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
	if err != nil {
		panic(err)
	}

	defer p.Close()

	// Delivery report handler for produced messages
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					CLogger.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					CLogger.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	// Produce messages to topic (asynchronously)
	topic := "employee"
	empBytesArr, err := json.Marshal(emp)
	if err != nil {
		return err
	}
	p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          empBytesArr,
	}, nil)

	// Wait for message deliveries before shutting down
	p.Flush(15 * 1000)
}
