## How to start with SOC4Kafka?

Choose an installation method that fits your environment:

- **Kubernetes (Helm):** Use the [Helm chart](helm/installation.md) to deploy SOC4Kafka on Kubernetes.
- **Automated (Ansible):** See the [Quickstart Guide](quickstart_guide.md) for automated installation.
- **Manual:** Follow the steps below to run the collector from a downloaded package and config file. For a full command-by-command walkthrough against a real source, see the [OCI Streaming on a VM](oci_installation.md) guide.

### Download Splunk OTel Collector package

The SOC4Kafka base package is the Splunk OpenTelemetry Collector, offering multiple installation methods to suit different needs.
Get the newest release (prefixed with `v`) using [this link](https://github.com/signalfx/splunk-otel-collector/releases), download
the package suited for your platform.

For instance, if you are using Linux on an AMD64 architecture, you can execute the following `wget` command:

```commandline
wget https://github.com/signalfx/splunk-otel-collector/releases/download/v0.158.0/otelcol_linux_amd64
```

### Create a minimal config template

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

#### Configuration Table

Mind that this is just a minimal configuration. You can customize it further based on your requirements by referring to the official documentation linked in the Component column.

| **Category**   | **Component**                                                                                                                         | **Parameter**               | **Description**                                                                            | **Required** | **Default Value** |
|----------------|---------------------------------------------------------------------------------------------------------------------------------------|-----------------------------|--------------------------------------------------------------------------------------------|--------------|-------------------|
| **Receivers**  | [kafka](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/kafkareceiver)                           | `brokers`                   | Kafka broker addresses for message consumption.                                            | Yes          | N/A               |
|                |                                                                                                                                       | `logs.topics`               | Kafka list of topics to subscribe to for receiving messages.                               | Yes          | N/A               |
|                |                                                                                                                                       | `logs.encoding`             | Encoding format of the Kafka messages.                                                     | No           | `"text"`          |
| **Processors** | [resourcedetection](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/processor/resourcedetectionprocessor) |                             | Sets a `host` field based on a machine's information.                                      | No           | N/A               |
| **Exporters**  | [splunk_hec](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/exporter/splunkhecexporter)                  | `token`                     | Splunk HEC token for authentication.                                                       | Yes          | N/A               |
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

#### Example configuration

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

Fill the file with your data and save it with a `.yaml` extension. For example `config.yaml`.

### Run Splunk OTel Collector package with config file

To run SOC4Kafka Connect, use the base package along with a completed configuration template.

```commandline
./<otel_package> --config <config_file>
```

!!! note

    Ensure the file has executable permissions before running the command. On Linux-based systems you can add executable permissions using the following command:

```commandline
chmod a+x <otel_package>
```

**Example**: For Linux on AMD64 architecture:

```commandline
chmod a+x otelcol_linux_amd64
./otelcol_linux_amd64 --config config.yaml
```

To understand the collector's pipeline design, refer to the [Design](otel_design.md) guide.