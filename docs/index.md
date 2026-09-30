# Collector for Kafka

The Splunk Distribution of OpenTelemetry Collector for Kafka subscribes to a Kafka topic and streams its data to Splunk Observability Cloud or to the Splunk platform through a Splunk HTTP Event Collector. It replaces the Splunk Connect for Kafka [(kafka-connect-splunk)](https://github.com/splunk/kafka-connect-splunk).


## Features


### Scaling

The Splunk Distribution of OpenTelemetry Collector for Kafka supports horizontal scaling, allowing you to run multiple collector instances to handle increased Kafka message throughput. See [Scaling the collector](scaling-updated.md).


### Load balancing

The Splunk Distribution of OpenTelemetry Collector for Kafka supports load balancing across multiple collector instances, distributing the Kafka message processing workload evenly to improve reliability and performance. See [Load balancing](loadbalancing-updated.md).


## Requirements

1. Kafka versions 3.7.0, 3.8.0, 3.9.0, 4.0.0
2. A Splunk environment of version 9.x and above, configured with valid [HTTP Event Collector (HEC)](https://dev.splunk.com/enterprise/docs/devtools/httpeventcollector/) token.


## Supported platforms

- Apache Kafka
- Amazon Managed Streaming for Apache Kafka (Amazon MSK)
- Confluent Platform


## Key differences from Splunk Connect for Kafka

The following Splunk Connect for Kafka features are not available in the Splunk Distribution of OpenTelemetry Collector for Kafka:

- Acknowledgment support
- Protobuf encoding

!!! info

    Splunk Distribution of OpenTelemetry Collector for Kafka does not support HEC Acknowledgements.


## Get started

See the [Get started with SOC4Kafka](getting_started-updated.md) guide to choose an installation method — Kubernetes, Ansible, or manual — and walk through downloading the package, building a config file, and running the collector.

## Advanced configuration

To customize your collector configuration for your requirements, see these topics:

- [Collect events from multiple topics](multiple_topics-updated.md): Easily gather data from several Kafka topics at once.
- [Subscribe to topics using regex](regex_topics-updated.md): Dynamically subscribe to topics that match specific patterns using regular expressions.
- [Extract data from headers and timestamps](extracting_additional_data-updated.md): Access and make use of metadata, like headers and timestamps, for more detailed insights.


## Migration

To migrate from Splunk Connect for Kafka to the Splunk Distribution of OpenTelemetry Collector for Kafka, see [Migration from Splunk Connect for Kafka to Splunk OTel Collector for Kafka](migration-updated.md).


## Splunk dashboard

A preconfigured health dashboard is available. See [The Splunk Distribution of OpenTelemetry Collector for Kafka health dashboard](splunk-dashboard-updated.md).


## Log collection

To configure log collection, see [Collect logs from the Splunk Distribution of OpenTelemetry Collector for Kafka](collecting_own_logs-updated.md).
