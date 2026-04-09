package integration_tests

import (
	"fmt"
	"testing"
	"tests/common"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	// These must match the values configured in ci_scripts/ci_values.yaml.
	index      = "kafka"
	sourcetype = "k8s-integration-test"
	source     = "soc4kafka-k8s-test"

	// Longer timeouts for K8s: the collector runs as a deployment and
	// may need time to subscribe to newly-created topics.
	testTimeout = 3 * time.Minute
	pollTick    = 5 * time.Second
)

func TestK8sIntegration(t *testing.T) {
	t.Run("basic single topic", testBasicSingleTopic)
	t.Run("multiple topics", testMultipleTopics)
}

func testBasicSingleTopic(t *testing.T) {
	topicName := "k8s-integration-test"
	event := "k8s-integration-basic-event"

	common.AddKafkaTopic(t, topicName, 1, 1)
	common.SendMessageToKafkaTopic(t, topicName, event)

	searchQuery := fmt.Sprintf("%sindex=%s sourcetype=%s source=%s \"%s\"",
		common.EventSearchQueryString, index, sourcetype, source, event)
	startTime := "-5m@m"

	require.Eventually(t, func() bool {
		events := common.GetEventsFromSplunk(t, searchQuery, startTime)
		if len(events) < 1 {
			return false
		}
		t.Logf("Events received: %d", len(events))
		assert.Equal(t, event, events[0].(map[string]interface{})["_raw"],
			"Event body does not match")
		return true
	}, testTimeout, pollTick,
		"Query %q returned no events for topic %s", searchQuery, topicName)
}

func testMultipleTopics(t *testing.T) {
	topic1 := "k8s-multi-topic-1"
	topic2 := "k8s-multi-topic-2"
	event1 := "k8s-multi-event-topic1"
	event2 := "k8s-multi-event-topic2"

	common.AddKafkaTopic(t, topic1, 1, 1)
	common.AddKafkaTopic(t, topic2, 1, 1)

	common.SendMessageToKafkaTopic(t, topic1, event1)
	common.SendMessageToKafkaTopic(t, topic2, event2)

	searchQuery := fmt.Sprintf("%sindex=%s sourcetype=%s source=%s",
		common.EventSearchQueryString, index, sourcetype, source)
	startTime := "-5m@m"

	require.Eventually(t, func() bool {
		events := common.GetEventsFromSplunk(t, searchQuery, startTime)
		if len(events) < 2 {
			t.Logf("Events received so far: %d (waiting for 2)", len(events))
			return false
		}
		t.Logf("Events received: %d", len(events))

		rawEvents := make([]interface{}, len(events))
		for i, e := range events {
			rawEvents[i] = e.(map[string]interface{})["_raw"]
		}
		assert.Contains(t, rawEvents, event1, "Missing event from %s", topic1)
		assert.Contains(t, rawEvents, event2, "Missing event from %s", topic2)
		return true
	}, testTimeout, pollTick,
		"Query %q did not return events from both topics", searchQuery)
}
