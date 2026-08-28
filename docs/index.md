# SOC4Kafka collector

The new SOC4Kafka collector, built on OpenTelemetry, enables the collection of Kafka messages and forwards these events to Splunk. It serves as a replacement for the existing
Splunk Connect for Kafka [(kafka-connect-splunk)](https://github.com/splunk/kafka-connect-splunk).

## Requirements

1. Kafka version 3.7.0 and above.
   - Tested with following versions: 3.7.0, 3.8.0, 3.9.0, 4.0.0
2. A Splunk environment of version 9.x and above, configured with valid [HTTP Event Collector (HEC)](https://dev.splunk.com/enterprise/docs/devtools/httpeventcollector/) token.


!!! info

    HEC Acknowledgements are not supported in SOC4Kafka

## Support technologies

Splunk OTel Collector for Kafka lets you subscribe to a Kafka topic and stream the data to the Splunk HTTP event collector on the following technologies:

- Apache Kafka
- Amazon Managed Streaming for Apache Kafka (Amazon MSK)
- Confluent Platform

## Key differences to Splunk Connect for Kafka

Not supported features which are available in previous version of Splunk Connect for Kafka but are not available in SOC4Kafka collector:

- Acknowledgment support
- Protobuf encoding

## How to start with SOC4Kafka?

See the [Overview](getting_started.md) guide to choose an installation method — Kubernetes, Ansible, or manual — and walk through downloading the package, building a config file, and running the collector.

## Advanced configuration

Thanks to the flexibility of the OpenTelemetry Collector, the setup can be tailored to meet specific requirements:

- [Collecting events from multiple topics](multiple_topics.md): Easily gather data from several Kafka topics at once.
- [Subscribing to topics using regex](regex_topics.md): Dynamically subscribe to topics that match specific patterns using regular expressions.
- [Extracting data from headers and timestamps](extracting_additional_data.md): Access and make use of metadata, like headers and timestamps, for more detailed insights.

## Scaling

SOC4Kafka supports horizontal scaling, allowing you to run multiple collector instances to handle increased Kafka message throughput. For more details check the [Scaling](scaling.md) guide.

## Load Balancing

SOC4Kafka supports load balancing across multiple collector instances, distributing the Kafka message processing workload evenly to improve reliability and performance. For more details check the [Load Balancing](loadbalancing.md) guide.

## Migration

Migration from Splunk Connect for Kafka to SOC4Kafka is described in the [Migration](migration.md) guide.

## Splunk Dashboard

A preconfigured health dashboard is available — see [Splunk Dashboard](splunk-dashboard.md).

## Troubleshooting

For troubleshooting check the [Troubleshooting](helm/troubleshooting.md) guide.
