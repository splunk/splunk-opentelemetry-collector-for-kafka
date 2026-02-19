# Examples

## Basic Single Receiver/Exporter

```yaml
kafkaReceivers:
  - name: main
    brokers:
      - "kafka:9092"
    logs:
      topics:
        - "logs"
      encoding: text
    group_id: "soc4kafka"

splunkExporters:
  - name: primary
    endpoint: "https://splunk:8088/services/collector"
    secret: "splunk-hec-secret"
    source: "kafka"
    sourcetype: "otel:logs"
    index: "main"

pipelines:
  - name: logs
    type: logs
    receivers:
      - main
    exporters:
      - primary
    processors:
      - batch
      - resourcedetection
```

## Multiple Topics to Multiple Indexes

```yaml
kafkaReceivers:
  - name: app-logs
    brokers:
      - "kafka:9092"
    logs:
      topics:
        - "app-info"
        - "app-warn"
      encoding: json
    group_id: "soc4kafka-app"
  
  - name: error-logs
    brokers:
      - "kafka:9092"
    logs:
      topics:
        - "app-error"
        - "app-critical"
      encoding: json
    group_id: "soc4kafka-error"

splunkExporters:
  - name: main-index
    endpoint: "https://splunk:8088/services/collector"
    secret: "splunk-hec-main"
    source: "kafka"
    sourcetype: "otel:logs"
    index: "main"
  
  - name: error-index
    endpoint: "https://splunk:8088/services/collector"
    secret: "splunk-hec-error"
    source: "kafka"
    sourcetype: "otel:logs"
    index: "errors"

pipelines:
  - name: app-pipeline
    type: logs
    receivers:
      - app-logs
    exporters:
      - main-index
    processors:
      - batch
      - resourcedetection
  
  - name: error-pipeline
    type: logs
    receivers:
      - error-logs
    exporters:
      - error-index
    processors:
      - batch
      - resourcedetection
```

## Authenticated Kafka with Secret Management

```yaml
kafkaReceivers:
  - name: secure-main
    brokers:
      - "secure-kafka:9092"
    auth:
      plain_text:
        username: "kafka-user"
        secret: "kafka-auth-secret"
    logs:
      topics:
        - "secure-logs"
      encoding: json
    group_id: "soc4kafka-secure"

splunkExporters:
  - name: primary
    endpoint: "https://splunk:8088/services/collector"
    secret: "splunk-hec-secret"
    source: "kafka"
    sourcetype: "otel:logs"
    index: "main"

pipelines:
  - name: secure-logs
    type: logs
    receivers:
      - secure-main
    exporters:
      - primary
    processors:
      - batch
      - resourcedetection
```

## With Collector Logs Enabled

Enable collection of the collector's own logs for debugging and monitoring:

```yaml
kafkaReceivers:
  - name: main
    brokers:
      - "kafka:9092"
    logs:
      topics:
        - "logs"
      encoding: text
    group_id: "soc4kafka"

splunkExporters:
  - name: primary
    endpoint: "https://splunk:8088/services/collector"
    secret: "splunk-hec-secret"
    source: "kafka"
    sourcetype: "otel:logs"
    index: "main"

pipelines:
  - name: logs
    type: logs
    receivers:
      - main
    exporters:
      - primary
    processors:
      - batch
      - resourcedetection

# Enable collector's own logs and forward to Splunk
collectorLogs:
  enabled: true
  level: info
  outputPaths:
    - /var/log/otelcol/otel-collector.log
    - stdout
  errorOutputPaths:
    - /var/log/otelcol/otel-collector-errors.log
    - stderr
  # Automatically forwards logs to Splunk via filelog receiver
  forwardToSplunk:
    enabled: true
    exporter: ""  # Uses the "primary" exporter (or omit to use first exporter)
```

When enabled, collector logs will be:
- Written to files in `/var/log/otelcol/` inside the container
- Available in pod logs via `kubectl logs` (stdout/stderr)
- Automatically forwarded to Splunk using the referenced `splunkExporter` (uses its endpoint, token, index, source, sourcetype)
- Tracked with `file_storage` extension to prevent re-reading on restart

The chart automatically adds:
- `filelog` receiver to read collector log files
- `file_storage` extension for checkpointing
- `logs/internal` pipeline connecting filelog → processors → referenced exporter
