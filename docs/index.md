# Splunk Distribution of OpenTelemetry Collector for Kafka

The Splunk Distribution of OpenTelemetry Collector for Kafka subscribes to Kafka topics and sends their data to Splunk Observability Cloud or to the Splunk platform through a Splunk HTTP Event Collector (HEC) exporter. It replaces Splunk Connect for Kafka ([kafka-connect-splunk](https://github.com/splunk/kafka-connect-splunk)).

## Features

The Splunk Distribution of OpenTelemetry Collector for Kafka supports horizontal scaling and load balancing. See [Scale the Collector](operate/collector-for-kafka-scale.md) and [Load balance HEC traffic](operate/collector-for-kafka-load-balance.md).

## Requirements

1. Kafka versions 3.7.0, 3.8.0, 3.9.0, and 4.0.0.
2. Splunk Enterprise 9.x or later with a valid HEC token.

## Supported platforms

- Apache Kafka
- Amazon Managed Streaming for Apache Kafka (Amazon MSK)
- Confluent Platform

## Differences from Splunk Connect for Kafka

The Splunk Distribution of OpenTelemetry Collector for Kafka does not support acknowledgment support or Protobuf encoding. It also does not support HEC acknowledgments.

## Deploy the Collector

Choose a deployment method in [Deploy the Splunk Distribution of OpenTelemetry Collector for Kafka](deploy/collector-for-kafka-deploy.md).

## Configure the Collector

For receiver, processor, exporter, and pipeline guidance, see [Configure the Splunk Distribution of OpenTelemetry Collector for Kafka](configure/collector-for-kafka-configure.md).

## Advanced configuration

- [Collect events from multiple topics](configure/collector-for-kafka-configure-multiple-topics.md)
- [Subscribe to topics using regular expressions](configure/collector-for-kafka-configure-regex-topics.md)
- [Extract data from headers and timestamps](configure/collector-for-kafka-configure-extract-data.md)

## Migration

To migrate from Splunk Connect for Kafka, see [Migrate from Splunk Connect for Kafka](migrate/collector-for-kafka-migrate-from-sc4kafka.md).

## Monitor the Collector

Monitor the Splunk Distribution of OpenTelemetry Collector for Kafka with the [health dashboard](monitor/collector-for-kafka-dashboard.md) or [Collector logs](monitor/collector-for-kafka-collector-logs.md).
