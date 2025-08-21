package non_functional_tests

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"testing"
	"tests/common"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	splunkWaitTime  = time.Second * 60
	minIngestRateMb = 1
	minIngestLag    = 1
)

type Statistic struct {
	EarliestTime string `json:"earliest_time"`
	LatestTime   string `json:"latest_time"`
	TotalEvents  string `json:"total_events"`
}

func TestNonFunctional(t *testing.T) {
	index := "kafka"
	sourcetype := "otel-perf-tests"
	source := "otel"
	configFileTemplate := "basic_test.yaml.tmpl"

	numMsg, err := strconv.Atoi(os.Getenv("NUM_MSG"))
	if err != nil {
		t.Fatalf("Couldn't read number of messages from env, NUM_MSG: %s", os.Getenv("NUM_MSG"))
	}
	recordSize, err := strconv.Atoi(os.Getenv("RECORD_SIZE"))
	if err != nil {
		t.Fatalf("Couldn't read size of the event record from env, RECORD_SIZE: %s", os.Getenv("RECORD_SIZE"))
	}
	topicName := os.Getenv("TOPIC_NAME")
	if len(topicName) == 0 {
		t.Fatalf("TOPIC_NAME env variable not set")
	}

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
	err = common.StartKafkaPerfScript(topicName, numMsg, recordSize)
	if err != nil {
		t.Fatalf("Couldn't start Kafka performance script: %s", err.Error())
	}
	defer common.StopOTelKafkaConnector(t, connectorHandler)

	searchQuery := common.EventSearchQueryString + "index=*" +
		" | stats earliest(_time) as earliest_time, latest(_time) as latest_time, count as total_events"
	startTime := "-2m@m"
	require.Eventually(t, func() bool {
		statistics := common.CheckStatisticsFromSplunk[Statistic](t, searchQuery, startTime)
		if len(statistics) != 1 {
			return false
		}
		stats := statistics[0]
		totalEvents, err := strconv.Atoi(stats.TotalEvents)
		if totalEvents != numMsg {
			return false
		}
		assert.Equal(t, numMsg, totalEvents, "Expected %d events, but got %d", numMsg, totalEvents)

		// time
		earliestTime, err := strconv.ParseFloat(stats.EarliestTime, 64)

		latestTime, err := strconv.ParseFloat(stats.LatestTime, 64)

		if err != nil {
			log.Fatalf("Couldn't parse data from job request: %s", err.Error())
		}

		ingestionTime := latestTime - earliestTime
		dataVolumeMb := float64((numMsg * recordSize) / (1024 * 1024))
		ingestionRateMb := dataVolumeMb / ingestionTime
		fmt.Printf("Splunk igested %d events of size %d in %f seconds. Which results in %f mb/s ingestion rate\n", totalEvents, recordSize, ingestionTime, ingestionRateMb)
		fmt.Printf("ingestionRateMb: %f", ingestionRateMb)
		fmt.Printf("min-ingestionRateMb: %f", float64(minIngestRateMb))
		compare := ingestionRateMb >= float64(minIngestRateMb)
		fmt.Printf("compare: %t", compare)
		assert.GreaterOrEqual(t, ingestionRateMb, float64(minIngestRateMb), "Splunk ingestion rate of %f didn't satisfied minimum requirement of %d", ingestionRateMb, minIngestRateMb)

		return true
	}, common.TestCaseDuration, common.TestCaseTick, "Search query: \n\"%s\"\n returned NO events", searchQuery)
}
