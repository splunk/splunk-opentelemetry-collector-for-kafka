package common

import (
	"context"
	"fmt"
	"os/exec"
	"testing"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/stretchr/testify/require"
)

func AddKafkaTopic(t *testing.T, topicName string, numberOfPartitions int, replicationFactor int) {
	t.Logf("Adding Kafka topic: %s\n", topicName)

	// Check if the topic already exists
	if checkKafkaTopicExists(t, topicName) {
		// Log a warning message and return
		t.Logf("WARN: Kafka topic '%s' already exists, skipping creation.\n", topicName)
		t.Logf("WARN: This may lead to test failures if the topic is expected to be created fresh for each test run.\n")
		return
	}

	// Admin client configuration
	adminClient, err := kafka.NewAdminClient(&kafka.ConfigMap{
		"bootstrap.servers": GetConfigVariable("KAFKA_BROKER_ADDRESS"),
	})
	require.NoError(t, err, "Failed to create Kafka admin client")
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

	require.NoError(t, err, "Failed to create topic")

	// Handle the results of the topic creation
	for _, result := range results {
		if result.Error.Code() == kafka.ErrNoError {
			t.Logf("Topic '%s' created successfully.\n", result.Topic)
		} else {
			t.Logf("Failed to create topic '%s': %v\n", result.Topic, result.Error)
		}
	}
	require.Eventually(t, func() bool { return checkKafkaTopicExists(t, topicName) }, 30*time.Second, 5*time.Second, "Expected at least one result from topic creation")
}

func getKafkaTopicsList(t *testing.T) []string {
	// Create a new admin client
	adminClient, err := kafka.NewAdminClient(&kafka.ConfigMap{
		"bootstrap.servers": GetConfigVariable("KAFKA_BROKER_ADDRESS"),
	})
	require.NoError(t, err, "Failed to create Kafka admin client")

	defer adminClient.Close()

	// List topics
	metadata, err := adminClient.GetMetadata(nil, true, 5000)
	require.NoError(t, err, "Failed to list Kafka topics")

	// Extract topic names
	var topics []string
	for topicName := range metadata.Topics {
		topics = append(topics, topicName)
		//fmt.Println(topicName)
	}
	return topics
}

func checkKafkaTopicExists(t *testing.T, topicName string) bool {
	t.Logf("Checking if Kafka topic exists: %s\n", topicName)
	topics := getKafkaTopicsList(t)
	for _, topic := range topics {
		if topic == topicName {
			t.Logf("Kafka topic '%s' exists.\n", topicName)
			return true
		}
	}
	t.Logf("Kafka topic '%s' does not exist.\n", topicName)
	return false
}

func SendMessageToKafkaTopic(t *testing.T, topicName string, message string, headers ...kafka.Header) {
	t.Logf("Adding message to Kafka topic: %s\n", topicName)
	// Create a new producer
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": GetConfigVariable("KAFKA_BROKER_ADDRESS"),
	})
	require.NoError(t, err, "Failed to create Kafka producer")
	defer producer.Close()

	// Produce a message to the topic
	deliveryChan := make(chan kafka.Event)
	err = producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topicName, Partition: kafka.PartitionAny},
		Value:          []byte(message),
		Headers:        headers,
	}, deliveryChan)

	require.NoError(t, err, "Failed to produce message")

	// Wait for the delivery report
	e := <-deliveryChan
	m := e.(*kafka.Message)
	if m.TopicPartition.Error != nil {
		t.Logf("Failed to deliver message: %v\n", m.TopicPartition.Error)
	} else {
		t.Logf("Message delivered to %s [%d] at offset %d\n", *m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
	}
	close(deliveryChan)
}

func StartKafkaPerfScript(topicName string, numMsg int, recordSize int) error {
	cmd := exec.Command(
		"docker", "exec", "-i", "cp-kafka-container", "/bin/kafka-producer-perf-test",
		"--topic", topicName,
		"--num-records", fmt.Sprintf("%d", numMsg),
		"--record-size", fmt.Sprintf("%d", recordSize),
		"--throughput", "-1",
		"--producer-props", "bootstrap.servers=localhost:9092",
	)

	err := cmd.Run()
	return err
}
