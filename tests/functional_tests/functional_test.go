package functional_tests

import (
	"fmt"
	"testing"
	"tests/common"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		"KafkaBrokerAddress": common.GetConfigVariable("KAFKA_BROKER_ADDRESS"),
		"KafkaTopicName":     topicName,
		"SplunkHECToken":     common.GetConfigVariable("HEC_TOKEN"),
		"SplunkHECEndpoint":  fmt.Sprintf("https://%s:8088/services/collector", common.GetConfigVariable("HOST")),
		"Source":             source,
		"Index":              index,
		"Sourcetype":         sourcetype,
	}

	configFileName := common.PrepareConfigFile(t, configFileTemplate, replacements, common.ConfigFilesDir)
	connectorHandler := common.StartOTelKafkaConnector(t, configFileName, common.ConfigFilesDir)

	common.AddKafkaTopic(t, topicName, 1, 1)
	common.SendMessageToKafkaTopic(topicName, event)

	// check events in Splunk
	searchQuery := common.EventSearchQueryString + "index=" + index + " sourcetype=" + sourcetype + " source=" + source
	startTime := "-1m@m"
	require.Eventually(t, func() bool {
		events := common.CheckEventsFromSplunk[any](t, searchQuery, startTime)
		if len(events) < 1 {
			return false
		}
		fmt.Println(" =========>  Events received: ", len(events))
		assert.Equal(t, 1, len(events), "Expected one event for topic %s, but got %d", topicName, len(events))
		assert.Equal(t, event, events[0].(map[string]interface{})["_raw"], "Expected event body does not match")
		return true
	}, common.TestCaseDuration, common.TestCaseTick, "Search query: \n\"%s\"\n returned NO events for topic %s", searchQuery, topicName)
	defer common.StopOTelKafkaConnector(t, connectorHandler)
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
		"KafkaBrokerAddress": common.GetConfigVariable("KAFKA_BROKER_ADDRESS"),
		"KafkaTopicName1":    topicName1,
		"KafkaTopicName2":    topicName2,
		"SplunkHECToken":     common.GetConfigVariable("HEC_TOKEN"),
		"SplunkHECEndpoint":  fmt.Sprintf("https://%s:8088/services/collector", common.GetConfigVariable("HOST")),
		"Source1":            source1,
		"Source2":            source2,
		"Index":              index,
		"Sourcetype":         sourcetype,
	}

	configFileName := common.PrepareConfigFile(t, configFileTemplate, replacements, common.ConfigFilesDir)
	connectorHandler := common.StartOTelKafkaConnector(t, configFileName, common.ConfigFilesDir)
	defer common.StopOTelKafkaConnector(t, connectorHandler)

	common.AddKafkaTopic(t, topicName1, 1, 1)
	common.AddKafkaTopic(t, topicName2, 1, 1)
	common.SendMessageToKafkaTopic(topicName1, event+topicName1)
	common.SendMessageToKafkaTopic(topicName2, event+topicName2)

	searchQuery := common.EventSearchQueryString + "index=" + index + " sourcetype=" + sourcetype + " source=" + sourceSuf + "*"
	startTime := "-1m@m"
	require.Eventually(t, func() bool {
		events := common.CheckEventsFromSplunk[any](t, searchQuery, startTime)
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
	}, common.TestCaseDuration, common.TestCaseTick, "Search query: \n\"%s\"\n returned NO events for topics %s, %s", searchQuery, topicName1, topicName2)
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
		"KafkaBrokerAddress": common.GetConfigVariable("KAFKA_BROKER_ADDRESS"),
		"KafkaTopicName":     topicName,
		"SplunkHECToken":     common.GetConfigVariable("HEC_TOKEN"),
		"SplunkHECEndpoint":  fmt.Sprintf("https://%s:8088/services/collector", common.GetConfigVariable("HOST")),
		"Source":             source,
		"Index":              index,
		"Sourcetype":         sourcetype,
		"CustomHeader":       headerKey,
	}

	configFileName := common.PrepareConfigFile(t, configFileTemplate, replacements, common.ConfigFilesDir)
	connectorHandler := common.StartOTelKafkaConnector(t, configFileName, common.ConfigFilesDir)

	common.AddKafkaTopic(t, topicName, 1, 1)
	common.SendMessageToKafkaTopic(topicName, event,
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

	searchQuery := common.EventSearchQueryString + "index=" + indexHeaderVal + " sourcetype=" + sourcetypeHeaderVal + " source=" +
		sourceHeaderVal + " host=" + hostHeaderVal + " kafka.header." + headerKey + "=" + headerVal
	startTime := "-1m@m"
	require.Eventually(t, func() bool {
		events := common.CheckEventsFromSplunk[any](t, searchQuery, startTime)
		if len(events) < 1 {
			return false
		}
		fmt.Println(" =========>  Events received: ", len(events))
		assert.Equal(t, 1, len(events), "Expected one event for topic %s, but got %d", topicName, len(events))
		assert.Equal(t, event, events[0].(map[string]interface{})["_raw"], "Expected event body does not match")
		assert.Equal(t, headerVal, events[0].(map[string]interface{})["kafka.header."+headerKey], "Expected header value does not match")
		return true
	}, common.TestCaseDuration, common.TestCaseTick, "Search query: \n\"%s\"\n returned NO events for topic %s", searchQuery, topicName)

	// Check if there are no events with the original source, index, and sourcetype
	searchQuery = common.EventSearchQueryString + "index=" + index + " sourcetype=" + sourcetype + " source=" + source
	events := common.CheckEventsFromSplunk[any](t, searchQuery, startTime)
	assert.Equal(t, 0, len(events), "Expected zero events for topic %s but got %d", topicName, len(events))

	defer common.StopOTelKafkaConnector(t, connectorHandler)
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
		"KafkaBrokerAddress": common.GetConfigVariable("KAFKA_BROKER_ADDRESS"),
		"KafkaTopicName":     topicName,
		"SplunkHECToken":     common.GetConfigVariable("HEC_TOKEN"),
		"SplunkHECEndpoint":  fmt.Sprintf("https://%s:8088/services/collector", common.GetConfigVariable("HOST")),
		"Source":             source,
		"Index":              index,
		"Sourcetype":         sourcetype,
		"ExtractPattern":     "\\\\[" + extractPattern + "\\\\]",
		"FormatStr":          formatStr,
	}

	configFileName := common.PrepareConfigFile(t, configFileTemplate, replacements, common.ConfigFilesDir)
	connectorHandler := common.StartOTelKafkaConnector(t, configFileName, common.ConfigFilesDir)

	common.AddKafkaTopic(t, topicName, 1, 1)
	common.SendMessageToKafkaTopic(topicName, event)

	searchQuery := common.EventSearchQueryString + "index=" + index + " sourcetype=" + sourcetype + " source=" +
		source
	startTime := "2020-01-01T11:55:00"
	endTime := "2020-01-01T12:05:00"
	require.Eventually(t, func() bool {
		events := common.CheckEventsFromSplunk[any](t, searchQuery, startTime, endTime)
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
	}, common.TestCaseDuration, common.TestCaseTick, "Search query: \n\"%s\"\n returned NO events for topic %s", searchQuery, topicName)
	defer common.StopOTelKafkaConnector(t, connectorHandler)
}
