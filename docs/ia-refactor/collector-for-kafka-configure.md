# Configure the Collector for Kafka

## Create a minimal configuration template

```yaml
receivers:
  kafka:
    brokers: [<Brokers>]
    logs:
      topics:
        - <Topic>
      encoding: <Encoding>

processors:
  resourcedetection:
    detectors: ["system"]
    system:
      hostname_sources: ["os"]

exporters:
  splunk_hec:
    token: "<Splunk HEC Token>"
    endpoint: <Splunk HEC Endpoint>
    source: <Source>
    sourcetype: <Sourcetype>
    index: <Splunk index>
    tls:
      insecure_skip_verify: false
    splunk_app_name: "soc4kafka"
    sending_queue:
      enabled: true
      num_consumers: 10
      queue_size: 10000
      block_on_overflow: true
      sizer: items
      batch:
        min_size: 1000

service:
  pipelines:
    logs:
      receivers: [kafka]
      processors: [resourcedetection]
      exporters: [splunk_hec]
```

## Configuration reference

This table describes the minimal configuration. For customization options, see the linked component documentation.

| **Category**   | **Component**                                                                                                                         | **Parameter**               | **Description**                                                                            | **Required** | **Default Value** |
|----------------|---------------------------------------------------------------------------------------------------------------------------------------|-----------------------------|--------------------------------------------------------------------------------------------|--------------|-------------------|
| **Receivers**  | [kafka](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/kafkareceiver)                           | `brokers`                   | Kafka broker addresses for message consumption.                                            | Yes          | N/A               |
|                |                                                                                                                                       | `logs.topics`               | Kafka list of topics to subscribe to for receiving messages.                               | Yes          | N/A               |
|                |                                                                                                                                       | `logs.encoding`             | Encoding format of the Kafka messages.                                                     | No           | `"text"`          |
| **Processors** | [resourcedetection](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/processor/resourcedetectionprocessor) |                             | Sets a `host` field based on a machine's information.                                      | No           | N/A               |
| **Exporters**  | [splunk_hec](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/exporter/splunkhecexporter)                  | `token`                     | Splunk HTTP Event Collector (HEC) token for authentication.                                                       | Yes          | N/A               |
|                |                                                                                                                                       | `endpoint`                  | Splunk HEC endpoint URL for sending data.                                                  | Yes          | N/A               |
|                |                                                                                                                                       | `source`                    | Source metadata for events sent to Splunk.                                                 | No           | `"otel"`          |
|                |                                                                                                                                       | `sourcetype`                | Sourcetype metadata for events sent to Splunk.                                             | No           | `"otel"`          |
|                |                                                                                                                                       | `index`                     | Splunk index where the logs will be stored.                                                | Yes          | N/A               |
|                |                                                                                                                                       | `tls.insecure_skip_verify`  | Whether to skip checking the certificate of the HEC endpoint when sending data over HTTPS. | No           | false             |
|                |                                                                                                                                       | `sending_queue.queue_size`         | Maximum number of queued items waiting to be exported.                                     | No           | 10000             |
|                |                                                                                                                                       | `sending_queue.block_on_overflow`  | Applies backpressure instead of immediately rejecting data when the exporter queue is full. | No           | true              |
|                |                                                                                                                                       | `sending_queue.sizer`              | Counts queue capacity by items.                                                           | No           | items             |
|                |                                                                                                                                       | `sending_queue.batch`              | Enables exporter-level batching before requests are sent to Splunk HEC.                    | No           | enabled           |
|                |                                                                                                                                       | `sending_queue.batch.min_size`     | Minimum number of items to batch before sending a request.                                 | No           | 1000              |
| **Service**    |                                                                                                                                       | `pipelines.logs.receivers`  | Specifies the receiver(s) for the log pipeline.                                            | Yes          | N/A               |
|                |                                                                                                                                       | `pipelines.logs.processors` | Specifies the processor(s) for the log pipeline.                                           | No           | `[]` (empty)      |
|                |                                                                                                                                       | `pipelines.logs.exporters`  | Specifies the exporter(s) for the log pipeline.                                            | Yes          | N/A               |

## Example configuration

```yaml
receivers:
  kafka:
    brokers: ["kafka-broker-1:9092", "kafka-broker-2:9092", "kafka-broker-3:9092"]
    logs:
      topics:
       - "example-topic"
      encoding: "text"

processors:
  resourcedetection:
    detectors: ["system"]
    system:
      hostname_sources: ["os"]

exporters:
  splunk_hec:
    token: "your-splunk-hec-token"
    endpoint: "https://splunk-hec-endpoint:8088/services/collector"
    source: my-kafka
    sourcetype: kafka-otel
    index: kafka_otel
    tls:
      insecure_skip_verify: false
    splunk_app_name: "soc4kafka"
    sending_queue:
      enabled: true
      num_consumers: 10
      queue_size: 10000
      block_on_overflow: true
      sizer: items
      batch:
        min_size: 1000

service:
  pipelines:
    logs:
      receivers: [kafka]
      processors: [resourcedetection]
      exporters: [splunk_hec]
```

Enter your values in the file and save it with a `.yaml` extension, for example, `config.yaml`.

For information about component roles and pipeline design, see [Design the Collector for Kafka](collector-for-kafka-design.md). For deployment-specific chart values, see [Configure the Helm chart](collector-for-kafka-kubernetes-configure-helm.md).
