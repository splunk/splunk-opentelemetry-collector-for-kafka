package performance_tests

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"strconv"
	"testing"
	"tests/common"
)

func TestPerformance(t *testing.T) {
	index := "kafka"
	sourcetype := "otel-perf-tests"
	source := "otel"
	configFileTemplate := "basic_perf_test.yaml.tmpl"

	numMsg, err := strconv.Atoi(os.Getenv("NUM_MSG"))
	require.NoError(t, err, "Couldn't read number of messages from env, NUM_MSG: %s", os.Getenv("NUM_MSG"))

	recordSize, err := strconv.Atoi(os.Getenv("RECORD_SIZE"))
	require.NoError(t, err, "Couldn't read size of the event record from env, RECORD_SIZE: %s", os.Getenv("RECORD_SIZE"))

	topicName := os.Getenv("TOPIC_NAME")
	require.Greater(t, len(topicName), 0, "TOPIC_NAME env variable is not set")

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
	require.NoError(t, err, "Couldn't start Kafka performance script")
	defer common.StopOTelKafkaConnector(t, connectorHandler)

	searchQuery := "| tstats earliest(_time) as earliest_time, latest(_time) as latest_time, count where index=" + index +
		" sourcetype=" + sourcetype + " source=" + source
	startTime := "-10m@m"

	require.EventuallyWithTf(t, func(c *assert.CollectT) {
		statistics := common.GetStatisticsFromSplunk(t, searchQuery, startTime)
		if len(statistics) != 1 {
			return
		}
		stats := statistics[0]
		totalEvents, err := strconv.Atoi(stats.TotalEvents)

		assert.NoError(c, err, "Couldn't parse total events from job request")
		assert.Equal(c, numMsg, totalEvents, "Expected %d events, but got %d", numMsg, totalEvents)
		if totalEvents != numMsg {
			return
		}

		earliestTime, err := strconv.ParseFloat(stats.EarliestTime, 64)
		require.NoError(t, err, "Couldn't parse earliest time from job request")
		latestTime, err := strconv.ParseFloat(stats.LatestTime, 64)
		require.NoError(t, err, "Couldn't parse latest time from job request")

		ingestionTime := latestTime - earliestTime
		dataVolumeMb := float64((numMsg * recordSize) / (1024 * 1024))
		ingestionRateMb := dataVolumeMb / ingestionTime
		minIngestRateMb, err := common.GetMinimumIngestRate()
		require.NoError(t, err, "Couldn't get min ingest rate for this test")

		t.Logf("Splunk igested %d events of size %d in %f seconds. Which results in %f mb/s ingestion rate\n", totalEvents, recordSize, ingestionTime, ingestionRateMb)
		t.Logf("Ingestion Rate [mb/s]: %f", ingestionRateMb)
		t.Logf("Minimum requried ingestion rate [mb/s]: %f", float64(minIngestRateMb))

		require.GreaterOrEqual(t, ingestionRateMb, float64(minIngestRateMb), "Splunk ingestion rate of %f didn't satisfied minimum requirement of %d", ingestionRateMb, minIngestRateMb)
	}, common.PerfTestCaseDuration, common.TestCaseTick, "Test with search query: \n%s\n failed", searchQuery)

}
