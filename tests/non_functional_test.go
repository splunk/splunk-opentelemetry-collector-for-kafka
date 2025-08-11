package tests

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	splunkWaitTime     = time.Second * 60
	minIngestionRateMb = 1
)

type Statistic struct {
	EarliestTime string `json:"earliest_time"`
	LatestTime   string `json:"latest_time"`
	TotalEvents  string `json:"total_events"`
}

func TestNonFunctional(t *testing.T) {
	topicName := "kafka-perf-test"
	index := "kafka1"
	sourcetype := "otel"
	source := "otel"
	configFileTemplate := "basic_test.yaml.tmpl"
	numOfMsg := 1_000_000
	msgSize := 100

	replacements := map[string]any{
		"KafkaBrokerAddress": GetConfigVariable("KAFKA_BROKER_ADDRESS"),
		"KafkaTopicName":     topicName,
		"SplunkHECToken":     GetConfigVariable("HEC_TOKEN"),
		"SplunkHECEndpoint":  fmt.Sprintf("http://%s:8088/services/collector", GetConfigVariable("HOST")),
		"Source":             source,
		"Index":              index,
		"Sourcetype":         sourcetype,
	}

	configFileName := prepareConfigFile(t, configFileTemplate, replacements, configFilesDir)
	connectorHandler := startOTelKafkaConnector(t, configFileName, configFilesDir)
	defer stopOTelKafkaConnector(t, connectorHandler)

	addKafkaTopic(t, topicName, 1, 1)
	firstMsgSendTime, lstMsgSendTime := sendRandomizedMessages(topicName, numOfMsg, msgSize)
	fmt.Printf("It took kafka client %f seconds to send all events", lstMsgSendTime.Sub(firstMsgSendTime).Seconds())

	fmt.Printf("Waiting %f seconds for splunk to ingest all data", splunkWaitTime.Seconds())
	time.Sleep(splunkWaitTime)
	fmt.Printf("Finished waiting")

	searchQuery := eventSearchQueryString + "index=*" +
		" | stats earliest(_time) as earliest_time, latest(_time) as latest_time, count as total_events"
	startTime := "-2m@m"
	require.Eventually(t, func() bool {
		statistics := checkStatisticsFromSplunk[Statistic](t, searchQuery, startTime)
		if len(statistics) != 1 {
			return false
		}
		stats := statistics[0]
		totalEvents, _ := strconv.Atoi(stats.TotalEvents)
		fmt.Printf("\n %s %s %s \n", statistics[0].EarliestTime, statistics[0].LatestTime, statistics[0].TotalEvents)
		assert.Equal(t, numOfMsg, totalEvents, "Expected %d events, but got %d", topicName, totalEvents)

		// time
		earliestTime, _ := strconv.ParseFloat(stats.EarliestTime, 64)

		latestTime, _ := strconv.ParseFloat(stats.LatestTime, 64)

		ingestionTime := latestTime - earliestTime
		dataVolumeMb := float64((numOfMsg * msgSize) / (1024 * 1024))
		ingestionRateMb := dataVolumeMb / ingestionTime
		fmt.Printf("Splunk igested %d events of size %d in %f seconds. Which results in %f mb/s ingestion rate", totalEvents, msgSize, ingestionTime, ingestionRateMb)
		assert.Equal(t, minIngestionRateMb, ingestionRateMb, "Splunk ingestion rate of %f didn't satisfied minimum requirement of %f", ingestionRateMb, minIngestionRateMb)

		fmt.Printf("\n Ingest time was: %f \n", latestTime-earliestTime)
		ingestLag := earliestTime - float64(firstMsgSendTime.Unix())
		fmt.Printf("Ingest lag %f", ingestLag)
		return true
	}, testCaseDuration, testCaseTick, "Search query: \n\"%s\"\n returned NO events for topic %s", searchQuery, topicName)

	fmt.Printf("End of test")
}
