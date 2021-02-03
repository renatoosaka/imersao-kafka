package main

import (
	"fmt"
	"sync"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
)

var wg sync.WaitGroup

func main() {
	configMap := &ckafka.ConfigMap{
		"bootstrap.servers": "kafka:9092",
	}

	producer, err := ckafka.NewProducer(configMap)

	if err != nil {
		panic(err)
	}

	topic := "teste"
	msg := "Hello Fullcycle"
	deliveryChan := make(chan ckafka.Event)

	message := &ckafka.Message{
		TopicPartition: ckafka.TopicPartition{Topic: &topic, Partition: ckafka.PartitionAny},
		Value:          []byte(msg),
	}

	err = producer.Produce(message, deliveryChan)
	if err != nil {
		panic(err)
	}

	wg.Add(1)

	go func() {
		for e := range deliveryChan {
			switch ev := e.(type) {
			case *ckafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Println("Delivery failed:", ev.TopicPartition)
				} else {
					fmt.Println("Delivered message:", ev.TopicPartition)
				}
			}
		}
		wg.Done()
	}()

	wg.Wait()
}
