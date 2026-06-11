## Design

The SOC4Kafka collector is designed using the OpenTelemetry Collector framework and is composed of various classes of pipeline components. The key components of the Kafka OpenTelemetry (OTel) collector include:
- Receivers
- Processors
- Exporters

### Receivers

The Kafka receiver is responsible for fetching data from the Kafka cluster. Detailed configuration instructions for this receiver can be found [here](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/receiver/kafkareceiver/README.md).

### Processors

Processors are optional components within the data pipeline that transform data before it is exported. Each processor performs specific actions based on its configuration, such as filtering or dropping data, among others. SOC4Kafka configures Splunk HEC batching in the exporter `sending_queue.batch` instead of using a pipeline `batch` processor. More information about configuring processors is available [here](https://github.com/open-telemetry/opentelemetry-collector/tree/main/processor#general-information).

### Exporters

The Splunk HEC exporter is used to send data to a Splunk HEC index. Detailed configuration guidelines for this exporter can be found [here](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/exporter/splunkhecexporter/README.md).

![SOC4Kafka scheme](images/kafka-otel-scheme.png)
