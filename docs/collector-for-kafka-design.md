# Understand the Collector for Kafka design

The Collector for Kafka uses the OpenTelemetry Collector framework and includes three pipeline component types:

- Receivers
- Processors
- Exporters

![OpenTelemetry Collector pipeline with receivers, processors, and exporters](../images/kafka-otel-scheme.png)

## Receivers

The Kafka receiver collects data from a Kafka cluster. For configuration details, see the [Kafka receiver documentation](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/receiver/kafkareceiver/README.md).

## Processors

Processors are optional pipeline components that transform, filter, or drop data before export. The Collector for Kafka configures Splunk HTTP Event Collector (HEC) batching in the exporter `sending_queue.batch` setting instead of using a pipeline `batch` processor. For more information, see the [OpenTelemetry Collector processor documentation](https://github.com/open-telemetry/opentelemetry-collector/tree/main/processor#general-information).

## Exporters

Use the Splunk HEC exporter to send data to a Splunk index. For configuration details, see the [Splunk HEC exporter documentation](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/exporter/splunkhecexporter/README.md).
