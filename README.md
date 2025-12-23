# SOC4Kafka connector

The new SOC4Kafka connector, built on OpenTelemetry, enables the collection of Kafka messages and forwards these events to Splunk. It serves as a replacement for the existing Kafka Connector [(kafka-connect-splunk)](https://github.com/splunk/kafka-connect-splunk).

## Requirements

1. Kafka version 3.7.0 and above.
   - Tested with following versions: 3.7.0, 3.8.0, 3.9.0, 4.0.0
2. A Splunk environment of version 9.x and above, configured with valid [HTTP Event Collector (HEC)](https://dev.splunk.com/enterprise/docs/devtools/httpeventcollector/) token.

NOTE: HEC Acknowledgements are not supported in SOC4Kafka

## Support technologies

Splunk OTel Connector for Kafka lets you subscribe to a Kafka topic and stream the data to the Splunk HTTP event collector on the following technologies:

- Apache Kafka
- Amazon Managed Streaming for Apache Kafka (Amazon MSK)
- Confluent Platform

## Key differences to Splunk Connect for Kafka

Not supported features which are available in previous version of Splunk Connect for Kafka but are not available in SOC4Kafka connector:
- Acknowledgment support - Not supported
- Load balancing - Not supported
- Protobuf encoding - Not supported

## How to start with SOC4Kafka?
Follow the steps below to get started with SOC4Kafka. Or check our [Quickstart Guide](docs/quickstart_guide.md) for an automated installation using Ansible.

### Download Splunk OTel Collector package

The SOC4Kafka base package is the Splunk OpenTelemetry Connector, offering multiple installation methods to suit different needs.
Get the newest release (prefixed with `v`) using [this link](https://github.com/signalfx/splunk-otel-collector/releases), download 
the package suited for your platform.

For instance, if you are using Linux on an AMD64 architecture, you can execute the following `wget` command:

```commandline
wget https://github.com/signalfx/splunk-otel-collector/releases/download/v0.141.0/otelcol_linux_amd64
```

### Create a minimal config template

```yaml
receivers:
  kafka:
    brokers: [<Brokers>]
    logs:
      topic: <Topic>
      encoding: <Encoding>

processors:
  resourcedetection:
    detectors: ["system"]
    system:
      hostname_sources: ["os"]
  batch:

exporters:
  splunk_hec:
    token: <Splunk HEC Token>
    endpoint: <Splunk HEC Endpoint>
    source: <Source>
    sourcetype: <Sourcetype>
    index: <Splunk index>
    tls:
      insecure_skip_verify: false
    headers:
      "__splunk_app_name": "soc4kafka"

service:
  pipelines:
    logs:
      receivers: [kafka]
      processors: [batch, resourcedetection]
      exporters: [splunk_hec]
```

#### Configuration Table

Mind that this is just a minimal configuration. You can customize it further based on your requirements by referring to the official documentation linked in the Component column.

| **Category**   | **Component**                                                                                                                         | **Parameter**               | **Description**                                                                            | **Required** | **Default Value** |
|----------------|---------------------------------------------------------------------------------------------------------------------------------------|-----------------------------|--------------------------------------------------------------------------------------------|--------------|-------------------|
| **Receivers**  | [kafka](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/kafkareceiver)                           | `brokers`                   | Kafka broker addresses for message consumption.                                            | Yes          | N/A               |
|                |                                                                                                                                       | `logs.topic`                | Kafka topic to subscribe to for receiving messages.                                        | Yes          | N/A               |
|                |                                                                                                                                       | `logs.encoding`             | Encoding format of the Kafka messages.                                                     | No           | `"text"`          |
| **Processors** | [batch](https://github.com/open-telemetry/opentelemetry-collector/tree/main/processor/batchprocessor)                                 |                             | Groups messages into batches before exporting.                                             | No           | N/A               |
|                | [resourcedetection](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/processor/resourcedetectionprocessor) |                             | Sets a `host` field based on a machine's information.                                      | No           | N/A               |
| **Exporters**  | [splunk_hec](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/exporter/splunkhecexporter)                  | `token`                     | Splunk HEC token for authentication.                                                       | Yes          | N/A               |
|                |                                                                                                                                       | `endpoint`                  | Splunk HEC endpoint URL for sending data.                                                  | Yes          | N/A               |
|                |                                                                                                                                       | `source`                    | Source metadata for events sent to Splunk.                                                 | No           | `"otel"`          |
|                |                                                                                                                                       | `sourcetype`                | Sourcetype metadata for events sent to Splunk.                                             | No           | `"otel"`          |
|                |                                                                                                                                       | `index`                     | Splunk index where the logs will be stored.                                                | Yes          | N/A               |
|                |                                                                                                                                       | `tls.insecure_skip_verify`  | Whether to skip checking the certificate of the HEC endpoint when sending data over HTTPS. | No           | false             |
| **Service**    |                                                                                                                                       | `pipelines.logs.receivers`  | Specifies the receiver(s) for the log pipeline.                                            | Yes          | N/A               |
|                |                                                                                                                                       | `pipelines.logs.processors` | Specifies the processor(s) for the log pipeline.                                           | No           | `[]` (empty)      |
|                |                                                                                                                                       | `pipelines.logs.exporters`  | Specifies the exporter(s) for the log pipeline.                                            | Yes          | N/A               |

#### Example configuration

```yaml
receivers:
  kafka:
    brokers: ["kafka-broker-1:9092", "kafka-broker-2:9092", "kafka-broker-3:9092"]
    logs:
      topic: "example-topic"
      encoding: "text"

processors:
  resourcedetection:
    detectors: ["system"]
    system:
      hostname_sources: ["os"]
  batch:

exporters:
  splunk_hec:
    token: "your-splunk-hec-token"
    endpoint: "https://splunk-hec-endpoint:8088/services/collector"
    source: my-kafka
    sourcetype: kafka-otel
    index: kafka_otel
    tls:
      insecure_skip_verify: false
    headers:
       "__splunk_app_name": "soc4kafka"

service:
  pipelines:
    logs:
      receivers: [kafka]
      processors: [batch, resourcedetection]
      exporters: [splunk_hec]
```

Fill the file with your data and save it with a `.yaml` extension. For example `config.yaml`.

### Run Splunk OTel Collector package with config file

To run SOC4Kafka Connect, use the base package along with a completed configuration template.

```commandline
./<otel_package> --config <config_file>
```

**NOTE**: Ensure the file has executable permissions before running the command. On Linux-based systems you can add executable permissions using the following command:

```commandline
chmod a+x <otel_package>
```

**Example**: For Linux on AMD64 architecture:

```commandline
chmod a+x otelcol_linux_amd64
./otelcol_linux_amd64 --config config.yaml
```

## Advanced configuration

Thanks to the flexibility of the OpenTelemetry Collector, the setup can be tailored to meet specific requirements. This modular approach allows you to treat the components as building blocks, enabling you to create a pipeline that aligns perfectly with your use case.
To understand the design, refer to [this guide](docs/otel_design.md).

You can unlock a range of powerful features by adjusting the configuration, such as:
- [Collecting events from multiple topics](docs/multiple_topics.md): Easily gather data from several Kafka topics at once.
- [Subscribing to topics using regex](docs/regex_topics.md): Dynamically subscribe to topics that match specific patterns using regular expressions.
- [Extracting data from headers and timestamps](docs/extracting_additional_data.md): Access and make use of metadata, like headers and timestamps, for more detailed insights.

## Scaling 

For scaling check [this guide](docs/scaling.md).

## Migration 

Migration from Splunk Connect for Kafka to SOC4Kafka is [described here](docs/migration.md).

## Troubleshooting

For troubleshooting check [this guide](docs/troubleshooting.md).