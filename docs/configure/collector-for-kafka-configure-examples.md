# Splunk Distribution of OpenTelemetry Collector for Kafka configuration examples

## Configure a single receiver and exporter

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
    sending_queue:
      enabled: true
      num_consumers: 10
      queue_size: 10000
      block_on_overflow: true
      sizer: items
      batch:
        min_size: 1000

pipelines:
  - name: logs
    type: logs
    receivers:
      - main
    exporters:
      - primary
    # processors omitted: defaults to ["resourcedetection"]
```

## Send data from multiple topics to multiple indexes

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
      - resourcedetection
  
  - name: error-pipeline
    type: logs
    receivers:
      - error-logs
    exporters:
      - error-index
```

## Connect to Kafka brokers with TLS

To connect to Kafka brokers that use TLS, add a `tls` block. Use `ca_pem` for a custom CA certificate. Set `insecure_skip_verify` only in development environments:

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

For TLS options and security recommendations, see [Configure TLS](collector-for-kafka-configure-tls.md).

## Authenticate with Kafka by using secrets

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

## Collect internal logs

Enable collection of internal logs from the Splunk Distribution of OpenTelemetry Collector for Kafka to support debugging and monitoring:

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

When enabled, the chart writes collector logs to files in `/var/log/otelcol/` inside the container. The logs are also available in pod output and can be forwarded to Splunk.

Use `kubectl logs` to view the logs written to standard output and standard error. The chart can also forward logs to Splunk by using the referenced `splunkExporter`, which provides the endpoint, token, index, source, and sourcetype.

The chart tracks log file positions with the `file_storage` extension so that it does not read the same logs again after a restart.

The chart adds the following components:

- `filelog` receiver to read collector log files
- `file_storage` extension for checkpointing
- A `logs/internal` pipeline that connects the `filelog` receiver, processors, and referenced exporter.

## Collect metrics

Enable collection of internal collector metrics and system metrics, such as CPU, memory, disk, and network metrics:

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

When enabled, the chart adds the following components:

- **Prometheus receiver**: Scrapes the collector's internal telemetry endpoint on port 8888.
- **Hostmetrics receiver**: Collects system metrics for CPU, memory, disk, network, filesystems, and processes.
- **Telemetry service**: Exposes collector metrics through a Prometheus endpoint.
- **Metrics pipeline**: Forwards metrics to Splunk by using the referenced `splunkExporter`.

!!! note
    Create a metrics index in Splunk for the metrics data. The service exposes port 8888 for Prometheus scraping.
