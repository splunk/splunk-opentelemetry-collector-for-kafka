package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	eventSearchQueryString = "| search "
	configFilesDir         = "testdata/configs"
	testCaseDuration       = 30 * time.Second
	testCaseTick           = 5 * time.Second
)

func Test_Functions(t *testing.T) {

	t.Run("basic scenario with single topic", testBasicScenarioWithSingleTopic)
	t.Run("scenario with multiple topics", testScenarioWithMultipleTopic)
	t.Run("scenario with custom headers", testScenarioWithCustomHeaders)
	t.Run("scenario with timestamp extraction", testScenarioTimestampExtraction)
}

func testBasicScenarioWithSingleTopic(t *testing.T) {
	fmt.Println("Running basic scenario with single topic test")
	topicName := "kafka-test-topic"
	event := "Hello, Kafka!"
	index := "kafka"
	sourcetype := "otel-basic-test"
	source := "otel"
	configFileTemplate := "basic_test.yaml.tmpl"

	replacements := map[string]any{
		"KafkaBrokerAddress": GetConfigVariable("KAFKA_BROKER_ADDRESS"),
		"KafkaTopicName":     topicName,
		"SplunkHECToken":     GetConfigVariable("HEC_TOKEN"),
		"SplunkHECEndpoint":  fmt.Sprintf("https://%s:8088/services/collector", GetConfigVariable("HOST")),
		"Source":             source,
		"Index":              index,
		"Sourcetype":         sourcetype,
	}

	configFileName := prepareConfigFile(t, configFileTemplate, replacements, configFilesDir)
	connectorHandler := startOTelKafkaConnector(t, configFileName, configFilesDir)

	addKafkaTopic(t, topicName, 1, 1)
	sendMessageToKafkaTopic(topicName, event)

	// check events in Splunk
	searchQuery := eventSearchQueryString + "index=" + index + " sourcetype=" + sourcetype + " source=" + source
	startTime := "-1m@m"
	require.Eventually(t, func() bool {
		events := checkEventsFromSplunk[any](t, searchQuery, startTime)
		if len(events) < 1 {
			return false
		}
		fmt.Println(" =========>  Events received: ", len(events))
		assert.Equal(t, 1, len(events), "Expected one event for topic %s, but got %d", topicName, len(events))
		assert.Equal(t, event, events[0].(map[string]interface{})["_raw"], "Expected event body does not match")
		return true
	}, testCaseDuration, testCaseTick, "Search query: \n\"%s\"\n returned NO events for topic %s", searchQuery, topicName)
	defer stopOTelKafkaConnector(t, connectorHandler)
}

func testScenarioWithMultipleTopic(t *testing.T) {
	fmt.Println("Running scenario with multiple topics")
	topicName1 := "kafka-test-topic-1"
	topicName2 := "kafka-test-topic-2"
	event := "Hello "
	index := "kafka"
	sourcetype := "otel-multiple-topics"
	sourceSuf := "otel"
	source1 := sourceSuf + "-1"
	source2 := sourceSuf + "-2"
	configFileTemplate := "multiple_topics_test.yaml.tmpl"

	replacements := map[string]any{
		"KafkaBrokerAddress": GetConfigVariable("KAFKA_BROKER_ADDRESS"),
		"KafkaTopicName1":    topicName1,
		"KafkaTopicName2":    topicName2,
		"SplunkHECToken":     GetConfigVariable("HEC_TOKEN"),
		"SplunkHECEndpoint":  fmt.Sprintf("https://%s:8088/services/collector", GetConfigVariable("HOST")),
		"Source1":            source1,
		"Source2":            source2,
		"Index":              index,
		"Sourcetype":         sourcetype,
	}

	configFileName := prepareConfigFile(t, configFileTemplate, replacements, configFilesDir)
	connectorHandler := startOTelKafkaConnector(t, configFileName, configFilesDir)
	defer stopOTelKafkaConnector(t, connectorHandler)

	addKafkaTopic(t, topicName1, 1, 1)
	addKafkaTopic(t, topicName2, 1, 1)
	sendMessageToKafkaTopic(topicName1, event+topicName1)
	sendMessageToKafkaTopic(topicName2, event+topicName2)

	searchQuery := eventSearchQueryString + "index=" + index + " sourcetype=" + sourcetype + " source=" + sourceSuf + "*"
	startTime := "-1m@m"
	require.Eventually(t, func() bool {
		events := checkEventsFromSplunk[any](t, searchQuery, startTime)
		if len(events) < 1 {
			return false
		}
		fmt.Println(" =========>  Events received: ", len(events))
		assert.Equal(t, 2, len(events), "Expected two events, but got %d", len(events))

		assert.ElementsMatch(
			t,
			[]string{
				event + topicName1,
				event + topicName2,
			},
			[]interface{}{
				events[0].(map[string]interface{})["_raw"],
				events[1].(map[string]interface{})["_raw"],
			},
			"Expected event bodies do not match for topics %s and %s", topicName1, topicName2,
		)

		assert.ElementsMatch(
			t,
			[]string{
				source1,
				source2,
			},
			[]interface{}{
				events[0].(map[string]interface{})["source"],
				events[1].(map[string]interface{})["source"],
			},
			"Expected sources do not match for topics %s and %s", topicName1, topicName2,
		)

		return true
	}, testCaseDuration, testCaseTick, "Search query: \n\"%s\"\n returned NO events for topics %s, %s", searchQuery, topicName1, topicName2)
}

func testScenarioWithCustomHeaders(t *testing.T) {
	fmt.Println("Running tests for custom headers")
	topicName := "kafka-custom-headers-test"
	event := "This event should have extra headers!"
	index := "kafka"
	sourcetype := "otel-custom-headers-test"
	source := "otel"
	configFileTemplate := "custom_headers_test.yaml.tmpl"

	// Custom headers to be used in the test
	headerKey := "custom-header"
	headerVal := "test-header-value"
	indexHeaderVal := "kafka-header-index"
	sourceHeaderVal := "source-value-from-header"
	sourcetypeHeaderVal := "sourcetype-value-from-header"
	hostHeaderVal := "host-value-from-header"

	replacements := map[string]any{
		"KafkaBrokerAddress": GetConfigVariable("KAFKA_BROKER_ADDRESS"),
		"KafkaTopicName":     topicName,
		"SplunkHECToken":     GetConfigVariable("HEC_TOKEN"),
		"SplunkHECEndpoint":  fmt.Sprintf("https://%s:8088/services/collector", GetConfigVariable("HOST")),
		"Source":             source,
		"Index":              index,
		"Sourcetype":         sourcetype,
		"CustomHeader":       headerKey,
	}

	configFileName := prepareConfigFile(t, configFileTemplate, replacements, configFilesDir)
	connectorHandler := startOTelKafkaConnector(t, configFileName, configFilesDir)

	addKafkaTopic(t, topicName, 1, 1)
	sendMessageToKafkaTopic(topicName, event,
		kafka.Header{
			Key:   headerKey,
			Value: []byte(headerVal),
		},
		kafka.Header{
			Key:   "index",
			Value: []byte(indexHeaderVal),
		},
		kafka.Header{
			Key:   "source",
			Value: []byte(sourceHeaderVal),
		},
		kafka.Header{
			Key:   "sourcetype",
			Value: []byte(sourcetypeHeaderVal),
		},
		kafka.Header{
			Key:   "host",
			Value: []byte(hostHeaderVal),
		})

	searchQuery := eventSearchQueryString + "index=" + indexHeaderVal + " sourcetype=" + sourcetypeHeaderVal + " source=" +
		sourceHeaderVal + " host=" + hostHeaderVal + " kafka.header." + headerKey + "=" + headerVal
	startTime := "-1m@m"
	require.Eventually(t, func() bool {
		events := checkEventsFromSplunk[any](t, searchQuery, startTime)
		if len(events) < 1 {
			return false
		}
		fmt.Println(" =========>  Events received: ", len(events))
		assert.Equal(t, 1, len(events), "Expected one event for topic %s, but got %d", topicName, len(events))
		assert.Equal(t, event, events[0].(map[string]interface{})["_raw"], "Expected event body does not match")
		assert.Equal(t, headerVal, events[0].(map[string]interface{})["kafka.header."+headerKey], "Expected header value does not match")
		return true
	}, testCaseDuration, testCaseTick, "Search query: \n\"%s\"\n returned NO events for topic %s", searchQuery, topicName)

	// Check if there are no events with the original source, index, and sourcetype
	searchQuery = eventSearchQueryString + "index=" + index + " sourcetype=" + sourcetype + " source=" + source
	events := checkEventsFromSplunk[any](t, searchQuery, startTime)
	assert.Equal(t, 0, len(events), "Expected zero events for topic %s but got %d", topicName, len(events))

	defer stopOTelKafkaConnector(t, connectorHandler)
}

func testScenarioTimestampExtraction(t *testing.T) {
	fmt.Println("Running tests for timestamp extraction")
	sourceTimestamp := time.Now().Format("20060102150405")
	topicName := "kafka-timestamp-extraction"
	index := "kafka"
	sourcetype := "otel-timestamp-extraction-test"
	source := "otel-" + sourceTimestamp
	configFileTemplate := "timestamp_extraction_test.yaml.tmpl"
	extractPattern := "(?P<timestamp>[0-9]{4}-[0-9]{2}-[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2})"
	formatStr := "2006-01-02 15:04:05"
	timestampStr := "2020-01-01 12:00:00"
	timestamp, err := time.Parse(formatStr, timestampStr)
	require.NoError(t, err, "Error parsing timestamp")
	event := "[" + timestampStr + "]" + " This event should have a custom timestamp!"
	replacements := map[string]any{
		"KafkaBrokerAddress": GetConfigVariable("KAFKA_BROKER_ADDRESS"),
		"KafkaTopicName":     topicName,
		"SplunkHECToken":     GetConfigVariable("HEC_TOKEN"),
		"SplunkHECEndpoint":  fmt.Sprintf("https://%s:8088/services/collector", GetConfigVariable("HOST")),
		"Source":             source,
		"Index":              index,
		"Sourcetype":         sourcetype,
		"ExtractPattern":     "\\\\[" + extractPattern + "\\\\]",
		"FormatStr":          formatStr,
	}

	configFileName := prepareConfigFile(t, configFileTemplate, replacements, configFilesDir)
	connectorHandler := startOTelKafkaConnector(t, configFileName, configFilesDir)

	addKafkaTopic(t, topicName, 1, 1)
	sendMessageToKafkaTopic(topicName, event)

	searchQuery := eventSearchQueryString + "index=" + index + " sourcetype=" + sourcetype + " source=" +
		source
	startTime := "2020-01-01T11:55:00"
	endTime := "2020-01-01T12:05:00"
	require.Eventually(t, func() bool {
		events := checkEventsFromSplunk[any](t, searchQuery, startTime, endTime)
		if len(events) < 1 {
			return false
		}
		fmt.Println(" =========>  Events received: ", len(events))
		assert.Equal(t, 1, len(events), "Expected one event for topic %s, but got %d", topicName, len(events))
		rawEvent := events[0].(map[string]interface{})["_raw"].(string)
		assert.Equal(t, event, rawEvent, "Expected event body does not match")

		eventTimeStr := events[0].(map[string]interface{})["_time"].(string)
		eventTime, err := time.ParseInLocation(time.RFC3339, eventTimeStr, time.UTC)
		require.NoError(t, err, "Error parsing event time from event")
		assert.Equal(t, timestamp, eventTime, "Event time does not match timestamp in the event body")

		return true
	}, testCaseDuration, testCaseTick, "Search query: \n\"%s\"\n returned NO events for topic %s", searchQuery, topicName)
	defer stopOTelKafkaConnector(t, connectorHandler)
}
