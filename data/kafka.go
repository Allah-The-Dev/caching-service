package data

import (
	"caching-service/config"
	"encoding/json"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

var (
	kafkaConsumer kafka.Consumer
)

//PublishToKafka ...
func (emp *Employee) PublishToKafka() {

	config.EmpAPILogger.Println("publishing emp data to kafka topic")

	p, err := kafka.NewProducer(
		&kafka.ConfigMap{
			"bootstrap.servers": config.KafkaHost,
		})
	if err != nil {
		config.EmpAPILogger.Println(err)
	}

	defer p.Close()

	// Delivery report handler for produced messages
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					config.EmpAPILogger.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					config.EmpAPILogger.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	// Produce messages to topic (asynchronously)
	topic := "employee"
	empBytesArr, err := json.Marshal(emp)
	if err != nil {
		config.EmpAPILogger.Println(err)
	}
	p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          empBytesArr,
	}, nil)

	// Wait for message deliveries before shutting down
	p.Flush(15 * 1000)
}

//StartKafkaConsumer ...
func StartKafkaConsumer() {
	kafkaConsumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": config.KafkaHost,
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})
	defer kafkaConsumer.Close()

	if err != nil {
		config.EmpAPILogger.Println(err)
	}

	kafkaConsumer.SubscribeTopics([]string{"employee", "^aRegex.*[Ee]mployee.*"}, nil)

	for {
		msg, err := kafkaConsumer.ReadMessage(-1)
		if err == nil && msg != nil {
			config.EmpAPILogger.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))

			emp := &Employee{}
			if err = json.Unmarshal(msg.Value, emp); err != nil {
				config.EmpAPILogger.Println(err)
			} else {
				emp.UpdateEmployeeCache()
			}
		} else {
			// The client will automatically try to recover from all errors.
			config.EmpAPILogger.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}
}
