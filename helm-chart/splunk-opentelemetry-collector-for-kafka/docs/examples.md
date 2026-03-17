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
    # processors omitted: defaults to ["batch", "resourcedetection"]
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
```

## Kafka with TLS (e.g. port 9093)

When connecting to TLS-enabled Kafka brokers, add a `tls` block. Use `ca_pem` for a custom CA and set `insecure_skip_verify` only for development:

```yaml
kafkaReceivers:
  - name: third
    brokers:
      - "secure-kafka:9093"
    logs:
      topics:
        - "perf3"
    group_id: "soc4kafka-main3"
    tls:
      insecure_skip_verify: true   # Use false in production;
      
      # provide ca_pem for broker CA
      ca_pem: |
        -----BEGIN CERTIFICATE-----
        ...
        G8jotQpS1QbFzo8o3fRN/xQ=
        -----END CERTIFICATE-----

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
      - third
    exporters:
      - primary
```

See [TLS Configuration](tls.md) for all options and security recommendations.

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

## With Metrics Collection Enabled

Enable collection of collector internal metrics and system metrics (CPU, memory, disk, network):

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
  
  - name: metrics
    endpoint: "https://splunk:8088/services/collector"
    secret: "splunk-hec-secret"
    source: "otel-collector"
    sourcetype: "otel:metrics"
    index: "metrics"

pipelines:
  - name: logs
    type: logs
    receivers:
      - main
    exporters:
      - primary

# Enable metrics collection
collectorMetrics:
  enabled: true
  exporter: "metrics"  # Optional: use specific exporter for metrics (or omit to use first exporter)
```

When enabled, the chart automatically adds:
- **Prometheus receiver** - Scrapes the collector's internal telemetry endpoint (port 8888)
- **Hostmetrics receiver** - Collects system metrics (CPU, memory, disk, network, filesystem, process)
- **Telemetry service** - Exposes collector metrics via Prometheus endpoint
- **Metrics pipeline** - Forwards metrics to Splunk using the referenced `splunkExporter`

**Note:** Make sure you have a metrics-type index in Splunk for the metrics data. The service exposes port 8888 for Prometheus scraping if needed.
