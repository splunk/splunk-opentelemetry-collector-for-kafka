package tests

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/stretchr/testify/require"
)

func addKafkaTopic(t *testing.T, topicName string, numberOfPartitions int, replicationFactor int) {
	fmt.Printf("Adding Kafka topic: %s\n", topicName)

	// Check if the topic already exists
	if checkKafkaTopicExists(topicName) {
		// Log a warning message and return
		log.Printf("WARN: Kafka topic '%s' already exists, skipping creation.\n", topicName)
		log.Printf("WARN: This may lead to test failures if the topic is expected to be created fresh for each test run.\n")
		return
	}

	// Admin client configuration
	adminClient, err := kafka.NewAdminClient(&kafka.ConfigMap{
		"bootstrap.servers": GetConfigVariable("KAFKA_BROKER_ADDRESS"),
	})
	if err != nil {
		log.Fatalf("Failed to create Kafka admin client: %s\n", err)
	}
	defer adminClient.Close()

	// Create the topic
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	results, err := adminClient.CreateTopics(
		ctx,
		[]kafka.TopicSpecification{
			{
				Topic:             topicName,
				NumPartitions:     numberOfPartitions,
				ReplicationFactor: replicationFactor,
			},
		},
		kafka.SetAdminOperationTimeout(10*time.Second),
	)

	if err != nil {
		log.Fatalf("Failed to create topic: %s\n", err)
	}

	// Handle the results of the topic creation
	for _, result := range results {
		if result.Error.Code() == kafka.ErrNoError {
			fmt.Printf("Topic '%s' created successfully.\n", result.Topic)
		} else {
			fmt.Printf("Failed to create topic '%s': %v\n", result.Topic, result.Error)
		}
	}
	require.Eventually(t, func() bool { return checkKafkaTopicExists(topicName) }, 30*time.Second, 5*time.Second, "Expected at least one result from topic creation")
}

func getKafkaTopicsList() []string {
	// Create a new admin client
	adminClient, err := kafka.NewAdminClient(&kafka.ConfigMap{
		"bootstrap.servers": GetConfigVariable("KAFKA_BROKER_ADDRESS"),
	})
	if err != nil {
		log.Fatalf("Failed to create Kafka admin client: %s\n", err)

	}
	defer adminClient.Close()
	// List topics
	metadata, err := adminClient.GetMetadata(nil, true, 5000)
	if err != nil {
		log.Fatalf("Failed to list Kafka topics: %s\n", err)

	}
	// Extract topic names
	var topics []string
	for topicName := range metadata.Topics {
		topics = append(topics, topicName)
		//fmt.Println(topicName)
	}
	return topics
}

func checkKafkaTopicExists(topicName string) bool {
	fmt.Printf("Checking if Kafka topic exists: %s\n", topicName)
	topics := getKafkaTopicsList()
	for _, topic := range topics {
		if topic == topicName {
			fmt.Printf("Kafka topic '%s' exists.\n", topicName)
			return true
		}
	}
	fmt.Printf("Kafka topic '%s' does not exist.\n", topicName)
	return false
}

func sendRandomizedMessages(topicName string, numOfMsg int, msgSize int) (time.Time, time.Time) {
	fmt.Printf("Sending randomized messages to kafka topic: %s\n", topicName)
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": GetConfigVariable("KAFKA_BROKER_ADDRESS"),
	})
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %s\n", err)
	}
	defer producer.Close()
	var firstMsgTime time.Time
	var lastMsgTime time.Time
	deliveryChan := make(chan kafka.Event, numOfMsg)
	for i := 0; i < numOfMsg; i++ {
		randomBytes := make([]byte, msgSize)
		_, err := rand.Read(randomBytes)
		if err != nil {
			log.Fatalf("Couldn't generate random message!")
		}
		kafkaMsg := kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topicName, Partition: kafka.PartitionAny},
			Value:          randomBytes,
		}
		if i == 0 {
			firstMsgTime = time.Now()
		}
		if i == numOfMsg-1 {
			lastMsgTime = time.Now()
		}
		if i%int(0.1*float32(numOfMsg)) == 0 {
			log.Printf("Sending %d message \n", i)
		}
		err = producer.Produce(&kafkaMsg, deliveryChan)
		if err != nil {
			log.Fatalf("Failed to send all messages to Kafka broker!")
		}
	}
	unFlushed := producer.Flush(10000)
	if unFlushed != 0 {
		log.Fatalf("Failed to send all messages to Kafka broker!")
	}
	close(deliveryChan)
	return firstMsgTime, lastMsgTime
}

func sendMessageToKafkaTopic(topicName string, message string, headers ...kafka.Header) {
	fmt.Printf("Adding message to Kafka topic: %s\n", topicName)
	// Create a new producer
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": GetConfigVariable("KAFKA_BROKER_ADDRESS"),
	})
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %s\n", err)
	}
	defer producer.Close()

	// Produce a message to the topic
	deliveryChan := make(chan kafka.Event)
	err = producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topicName, Partition: kafka.PartitionAny},
		Value:          []byte(message),
		Headers:        headers,
	}, deliveryChan)

	if err != nil {
		log.Fatalf("Failed to produce message: %s\n", err)
	}

	// Wait for the delivery report
	e := <-deliveryChan
	m := e.(*kafka.Message)
	if m.TopicPartition.Error != nil {
		log.Printf("Failed to deliver message: %v\n", m.TopicPartition.Error)
	} else {
		fmt.Printf("Message delivered to %s [%d] at offset %d\n", *m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
	}
	close(deliveryChan)
}
