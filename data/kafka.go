package data

import (
	"encoding/json"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

var KafkaHost string

//PublishToKafka ...
func (emp *Employee) PublishToKafka() {

	CLogger.Println("publishing emp data to kafka topic")

	p, err := kafka.NewProducer(
		&kafka.ConfigMap{
			"bootstrap.servers": KafkaHost,
		})
	if err != nil {
		CLogger.Println(err)
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
		CLogger.Println(err)
	}
	p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          empBytesArr,
	}, nil)

	// Wait for message deliveries before shutting down
	p.Flush(15 * 1000)
}

//KafkaConsumer ...
func KafkaConsumer() {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": KafkaHost,
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})
	defer c.Close()

	if err != nil {
		CLogger.Println(err)
	}

	c.SubscribeTopics([]string{"employee", "^aRegex.*[Ee]mployee.*"}, nil)

	for {
		msg, err := c.ReadMessage(-1)
		if err == nil && msg != nil {
			CLogger.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
		} else {
			// The client will automatically try to recover from all errors.
			CLogger.Printf("Consumer error: %v (%v)\n", err, msg)
		}

		emp := &Employee{}
		if err = json.Unmarshal(msg.Value, emp); err != nil {
			CLogger.Println(err)
		} else {
			emp.UpdateEmployeeCache()
		}
	}
}
