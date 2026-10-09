# Install the Splunk Distribution of OpenTelemetry Collector for Kafka with Helm

## Install the Helm chart

1. Create a `values.yaml` file with your configuration:

```yaml
kafkaReceivers:
  - name: main
    brokers:
      - "kafka-broker:9092"
    logs:
      topics:
        - "application-logs"
      encoding: text
    group_id: "soc4kafka-main"

splunkExporters:
  - name: primary
    endpoint: "https://splunk-hec:8088/services/collector"
    token: "your-splunk-hec-token"
    source: "soc4kafka"
    sourcetype: "otel:logs"
    index: "main"
    sending_queue:
      enabled: true
      num_consumers: 10
      queue_size: 10000
      block_on_overflow: true
      sizer: items
      batch:
        min_size: 1000

pipelines:
  - name: main-logs
    type: logs
    receivers:
      - main
    exporters:
      - primary
    # processors optional; defaults to ["resourcedetection"] (defaults.pipelineProcessors in values.yaml)
```

2. Add the Helm repository:

```bash
helm repo add splunk-opentelemetry-collector-for-kafka https://splunk.github.io/splunk-opentelemetry-collector-for-kafka
```

3. Install the chart:

```bash
helm upgrade --install soc4kafka splunk-opentelemetry-collector-for-kafka/splunk-opentelemetry-collector-for-kafka -f values.yaml
```

!!! note

    For information about managing secrets (auto-created or existing Kubernetes secrets), see [Secret management](../configure/collector-for-kafka-configure-secrets.md).
