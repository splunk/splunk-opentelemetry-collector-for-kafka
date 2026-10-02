# Configure the Helm chart for the Collector for Kafka

## Core configuration

### Kafka receivers

Define one or more Kafka receivers. The chart supports all standard Kafka receiver options. For details about the receiver, see [Understand the Collector for Kafka design](../collector-for-kafka-design.md).

```yaml
kafkaReceivers:
  - name: main
    brokers:
      - "kafka-broker-1:9092"
      - "kafka-broker-2:9092"
    logs:
      topics:
        - "application-logs"
        - "error-logs"
      encoding: text
    group_id: "soc4kafka-main"
```

**Chart-specific:** The `name` field is required. Use it to reference the receiver in pipelines.

!!! note

    For Kafka authentication with `plain_text`, SASL, or Kerberos, reference an existing Kubernetes Secret by using the `secret` field. See [Manage secrets for the Helm chart](../configure/collector-for-kafka-configure-secrets.md) for details.

**TLS:** For Kafka brokers that use TLS, add a `tls` block. For details and examples, see [Configure TLS](../configure/collector-for-kafka-configure-tls.md).

### Splunk HTTP Event Collector (HEC) exporters

Define one or more Splunk HEC exporters. The chart supports all standard Splunk HEC exporter options. For details about the exporter, see [Understand the Collector for Kafka design](../collector-for-kafka-design.md).

```yaml
splunkExporters:
  - name: primary
    endpoint: "https://splunk-hec:8088/services/collector"
    token: "your-splunk-hec-token"
    source: "soc4kafka"
    sourcetype: "otel:logs"
    index: "main"
    tls:
      insecure_skip_verify: false
    sending_queue:
      enabled: true
      num_consumers: 10
      queue_size: 10000
      block_on_overflow: true
      sizer: items
      batch:
        min_size: 1000
```

**Chart-specific:** The `name` field is required. Use it to reference the exporter in pipelines.

!!! note

    Instead of entering `token` directly, you can reference an existing Kubernetes Secret by using the `secret` field. See [Manage secrets for the Helm chart](../configure/collector-for-kafka-configure-secrets.md) for details.

**TLS:** Use an `https://` endpoint to connect with TLS. The same `tls` options as for Kafka apply. See [Configure TLS](../configure/collector-for-kafka-configure-tls.md).

**Queueing and batching:** The chart applies HEC exporter defaults from `defaults.exporters.splunk_hec`. The default queue uses `block_on_overflow: true` so the collector applies backpressure when Splunk HEC is slow instead of immediately rejecting data at queue capacity. Batching is configured under `sending_queue.batch`; the default pipelines no longer use the processor `batch`.

### Pipelines

Connect receivers to exporters. For details about pipelines, see [Understand the Collector for Kafka design](../collector-for-kafka-design.md).

**Chart-specific:** You can omit `processors`; the chart then uses `defaults.pipelineProcessors` (default: `["resourcedetection"]`). Override per pipeline or change the default in `values.yaml`.

```yaml
pipelines:
  - name: main-logs
    type: logs
    receivers:
      - main  # Must match a receiver name from kafkaReceivers
    exporters:
      - primary  # Must match an exporter name from splunkExporters
    # processors is optional; defaults to ["resourcedetection"] (see defaults.pipelineProcessors)
    processors:
      - resourcedetection
```

## Advanced configuration

See [values.yaml](https://github.com/splunk/splunk-opentelemetry-collector-for-kafka/blob/main/helm-chart/splunk-opentelemetry-collector-for-kafka/values.yaml) for all available configuration options. Key areas:

- **TLS** ([tls.md](../configure/collector-for-kafka-configure-tls.md)): Configure TLS for Kafka receivers and Splunk HEC exporters.
- **Component defaults** (`defaults`): Override default OpenTelemetry component settings.
- **Configuration override** (`configOverride`): Provide a complete OpenTelemetry configuration override.
- **Resources** (`resources`): Set CPU and memory limits and requests.
- **Autoscaling** (`autoscaling`): Configure the Horizontal Pod Autoscaler.
- **Pod disruption budget** (`podDisruptionBudget`): Configure a pod disruption budget for high availability.
- **Service account** (`serviceAccount`): Configure a service account with workload identity annotations for cloud environments, including AWS EKS, GCP GKE, and Azure AKS.
- **Collector logs** (`collectorLogs`): Collect the collector's internal logs to files and `stdout` or `stderr`.
- **Collector metrics** (`collectorMetrics`): Collect internal collector metrics and system metrics for CPU, memory, disk, and network.

### Collector logs

For information about configuring and collecting the Collector for Kafka logs, see [Collect logs from the Collector for Kafka](../monitor/collector-for-kafka-collector-logs.md).

### Metrics collection

Enable collection of internal collector metrics and system metrics. The chart configures the following components:

1. **Prometheus receiver**: Scrapes the collector's internal telemetry endpoint on port 8888.
2. **Hostmetrics receiver**: Collects system metrics for CPU, memory, disk, network, filesystems, and processes.
3. **Telemetry service**: Exposes collector metrics through a Prometheus endpoint.
4. **Metrics pipeline**: Forwards metrics to Splunk by using the first `splunkExporter`.

```yaml
collectorMetrics:
  enabled: true
  exporter: ""  # Optional: name of splunkExporter to use (e.g., "primary")
  # If not specified, uses the first splunkExporter
```

**Features:**

- Exposes internal collector metrics through a Prometheus endpoint on port 8888.
- Collects system metrics for CPU, memory, disk, network, filesystems, and processes by using the hostmetrics receiver.
- Forwards metrics to Splunk by using the referenced `splunkExporter`, or the first exporter if none is specified.
- Exposes a metrics port for Prometheus scraping.
- Sends all metrics through the `resourcedetection` processor for host filtering.

!!! note

    Create a metrics index in Splunk for the metrics data. For details, see the [Collector for Kafka dashboard](../monitor/collector-for-kafka-dashboard.md).

**Advanced configuration:** To customize metrics configuration, including exporters, scrapers, or intervals, use `configOverride` to override the generated configuration.

## Configuration precedence

The chart merges configuration in this order, from highest to lowest priority:

1. **`configOverride`**: Has the highest priority and completely overrides generated configuration. Use it for customizations that the other settings do not support.

2. **Explicit configuration** - Values specified directly in `kafkaReceivers` and `splunkExporters` override defaults. For example:
   
```yaml
kafkaReceivers:
 - name: main
   group_id: "custom-group"  # This overrides defaults.receivers.kafka.group_id
```

3. **`defaults`**: Has the lowest priority and provides default values for receivers, processors, and exporters that you do not configure explicitly.

**Merge behavior:**

- The chart uses `mustMergeOverwrite` to merge values recursively. Explicit values replace defaults at the same path.
- For nested objects, the chart replaces only the specified keys and preserves other default keys.
- The chart applies `configOverride` last, so it can override any part of the generated configuration.

**Example precedence:**
```yaml
defaults:
  exporters:
    splunk_hec:
      sending_queue:
        enabled: true
        num_consumers: 10
        queue_size: 10000
        block_on_overflow: true
        sizer: items
        batch:
          min_size: 1000

# If you specify in defaults but also in configOverride:
configOverride:
  exporters:
    splunk_hec:
      sending_queue:
        num_consumers: 20  # This wins - configOverride has highest priority
```

## Automatic pod restarts

The chart restarts pods automatically in these cases:

- **ConfigMap changes**: Pods restart when the OpenTelemetry configuration changes, based on the `checksum/config` annotation.
- **Secret changes**: Pods restart when token values in `values.yaml` change for chart-created Secrets or when Secret references change, based on the `checksum/secrets` annotation.

These restarts apply updated configuration and Secrets to the collector.
