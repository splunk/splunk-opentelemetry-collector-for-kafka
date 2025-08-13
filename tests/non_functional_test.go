package tests

import (
	"fmt"
	"log"
	"strconv"
	"testing"
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
	topicName := "kafka-perf-test"
	index := "kafka"
	sourcetype := "otel-perf-tests"
	source := "otel"
	configFileTemplate := "basic_test.yaml.tmpl"
	numOfMsg := 1_000_000
	msgSize := 100

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
	defer stopOTelKafkaConnector(t, connectorHandler)

	addKafkaTopic(t, topicName, 1, 1)
	firstMsgSendTime, lstMsgSendTime := sendRandomizedMessages(topicName, numOfMsg, msgSize)
	fmt.Printf("It took kafka client %f seconds to send all events\n", lstMsgSendTime.Sub(firstMsgSendTime).Seconds())

	fmt.Printf("Waiting %f seconds for splunk to ingest all data\n", splunkWaitTime.Seconds())
	time.Sleep(splunkWaitTime)
	fmt.Printf("Finished waiting\n")

	searchQuery := eventSearchQueryString + "index=" + index + " source=" + source + " sourcetype=" + sourcetype +
		" | stats earliest(_time) as earliest_time, latest(_time) as latest_time, count as total_events"
	startTime := "-2m@m"
	require.Eventually(t, func() bool {
		statistics := checkStatisticsFromSplunk[Statistic](t, searchQuery, startTime)
		if len(statistics) != 1 {
			return false
		}
		stats := statistics[0]
		totalEvents, err := strconv.Atoi(stats.TotalEvents)
		assert.Equal(t, numOfMsg, totalEvents, "Expected %d events, but got %d", numOfMsg, totalEvents)

		// time
		earliestTime, err := strconv.ParseFloat(stats.EarliestTime, 64)

		latestTime, err := strconv.ParseFloat(stats.LatestTime, 64)

		if err != nil {
			log.Fatalf("Couldn't parse data from job request: %s", err.Error())
		}

		ingestionTime := latestTime - earliestTime
		dataVolumeMb := float64((numOfMsg * msgSize) / (1024 * 1024))
		ingestionRateMb := dataVolumeMb / ingestionTime
		fmt.Printf("Splunk igested %d events of size %d in %f seconds. Which results in %f mb/s ingestion rate\n", totalEvents, msgSize, ingestionTime, ingestionRateMb)
		assert.GreaterOrEqual(t, ingestionRateMb, float64(minIngestRateMb), "Splunk ingestion rate of %f didn't satisfied minimum requirement of %d", ingestionRateMb, minIngestRateMb)

		ingestLag := earliestTime - float64(firstMsgSendTime.Unix())
		fmt.Printf("Ingest lag: %f\n", ingestLag)
		assert.LessOrEqual(t, ingestLag, float64(minIngestLag), "Ingest lag of %f seconds exceeded maximum value of %d second\n", ingestLag, minIngestLag)
		return true
	}, testCaseDuration, testCaseTick, "Search query: \n\"%s\"\n returned NO events", searchQuery)
}
