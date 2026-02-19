# Configuration

## Core Configuration

### Kafka Receivers

Define one or more Kafka receivers. All standard Kafka receiver configuration options are supported. See [OpenTelemetry Kafka Receiver documentation](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/kafkareceiver) for details.

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
    auth:
      # See Secret Management documentation
```

**Chart-specific:** The `name` field is required and used to reference the receiver in pipelines.

### Splunk HEC Exporters

Define one or more Splunk HEC exporters. All standard Splunk HEC exporter configuration options are supported. See [OpenTelemetry Splunk HEC Exporter documentation](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/exporter/splunkhecexporter) for details.

```yaml
splunkExporters:
  - name: primary
    endpoint: "https://splunk-hec:8088/services/collector"
    secret: "my-splunk-hec-secret"
    source: "soc4kafka"
    sourcetype: "otel:logs"
    index: "main"
    tls:
      insecure_skip_verify: false
```

**Chart-specific:** The `name` field is required and used to reference the exporter in pipelines.

### Pipelines

Connect receivers to exporters. See [OpenTelemetry Service Pipelines documentation](https://opentelemetry.io/docs/collector/configuration/#service) for details.

```yaml
pipelines:
  - name: main-logs
    type: logs
    receivers:
      - main  # Must match a receiver name from kafkaReceivers
    exporters:
      - primary  # Must match an exporter name from splunkExporters
    processors:
      - batch
      - resourcedetection
```

## Advanced Configuration

See [values.yaml](../values.yaml) for all available configuration options. Key areas:

- **Component Defaults** (`defaults`): Override default OpenTelemetry component settings
- **Config Override** (`configOverride`): Provide complete OpenTelemetry config override
- **Resources** (`resources`): Set CPU and memory limits/requests
- **Autoscaling** (`autoscaling`): Configure Horizontal Pod Autoscaler
- **Pod Disruption Budget** (`podDisruptionBudget`): Configure PDB for high availability
- **Service Account** (`serviceAccount`): Configure service account with workload identity annotations for cloud environments (AWS EKS, GCP GKE, Azure AKS)
- **Collector Logs** (`collectorLogs`): Enable collection of the collector's own logs to files and stdout/stderr
- **Metrics** (`enableMetrics`): Enable collection of collector internal metrics and system metrics (CPU, memory, disk, network)

### Collector Logs

Enable collection of the OpenTelemetry Collector's own logs. When enabled, logs are written to files in `/var/log/otelcol` (using an emptyDir volume) and also to stdout/stderr for Kubernetes log aggregation.

By default, when `collectorLogs.enabled` is `true`, the chart also forwards these logs to Splunk using a `filelog` receiver. This allows you to centrally monitor collector logs from multiple instances.

```yaml
collectorLogs:
  enabled: true
  level: info  # Options: debug, info, warn, error
  outputPaths:
    - /var/log/otelcol/otel-collector.log
    - stdout
  errorOutputPaths:
    - /var/log/otelcol/otel-collector-errors.log
    - stderr
  # Forward collector logs to Splunk (enabled by default when collectorLogs.enabled is true)
  forwardToSplunk:
    enabled: true
    # Reference to an existing splunkExporter by name (uses it directly, no overrides)
    # If not specified, uses the first splunkExporter
    exporter: ""  # Optional: name of splunkExporter to use (e.g., "primary")
  # File storage extension for checkpointing (prevents re-reading logs on restart)
  fileStorage:
    directory: /var/log/otelcol/checkpoint
    createDirectory: true
```

**Features:**
- Logs are written to files and stdout/stderr
- Logs are automatically forwarded to Splunk via `filelog` receiver
- `file_storage` extension tracks read position to prevent re-reading logs on restart
- Internal logs are sent to a separate Splunk index for easier analysis
- Uses the first `splunkExporter`'s endpoint and secret by default (can be overridden)

**Note:** Log files are stored in an `emptyDir` volume, which means they are ephemeral and will be lost when the pod is deleted. However, logs are forwarded to Splunk, so they are preserved there.

### Metrics Collection

Enable collection of collector internal metrics and system metrics. When enabled, the chart automatically configures:

1. **Prometheus receiver** - Scrapes the collector's internal telemetry endpoint (exposed on port 8888)
2. **Hostmetrics receiver** - Collects system metrics (CPU, memory, disk, network, filesystem, process)
3. **Telemetry service** - Exposes collector metrics via Prometheus endpoint
4. **Metrics pipeline** - Forwards metrics to Splunk using the first `splunkExporter`

```yaml
enableMetrics: true
```

**Features:**
- Collector internal metrics exposed via Prometheus endpoint (port 8888)
- System metrics collected via hostmetrics receiver (CPU, memory, disk, network, filesystem, process)
- Metrics forwarded to Splunk using the first `splunkExporter`
- Service exposes metrics port for Prometheus scraping
- All metrics go through the `resourcedetection` processor for host filtering

**Note:** Make sure you have a metrics-type index in Splunk for the metrics data. See the [Splunk Dashboard documentation](../../../docs/splunk-dashboard.md) for details.

**Advanced Configuration:** For custom metrics configuration (different exporter, scrapers, intervals, etc.), use `configOverride` to override the generated configuration.

## Configuration Precedence

The chart merges configuration in the following order (highest to lowest priority):

1. **`configOverride`** - Highest priority. Completely overrides any generated configuration. Use this for advanced customizations that can't be achieved through other means.

2. **Explicit configuration** - Values specified directly in `kafkaReceivers` and `splunkExporters` override defaults. For example:
   ```yaml
   kafkaReceivers:
     - name: main
       group_id: "custom-group"  # This overrides defaults.receivers.kafka.group_id
   ```

3. **`defaults`** - Lowest priority. Provides default values for receivers, processors, and exporters. These are used when not explicitly specified.

**Merge behavior:**
- `mustMergeOverwrite` is used, which means values are deeply merged, and explicit values completely replace defaults at the same path
- For nested objects, only the specified keys are replaced; other keys from defaults are preserved
- `configOverride` is applied last and can override any part of the generated configuration

**Example precedence:**
```yaml
defaults:
  processors:
    batch:
      timeout: 5s
      send_batch_size: 1000

# If you specify in defaults but also in configOverride:
configOverride:
  processors:
    batch:
      timeout: 15s  # This wins - configOverride has highest priority
```

## Automatic Pod Restarts

The chart includes automatic pod restart triggers:

- **ConfigMap changes**: Pods restart when the OpenTelemetry configuration changes (via `checksum/config` annotation)
- **Secret changes**: Pods restart when:
  - Token values in `values.yaml` change (for auto-created secrets)
  - Secret references change (via `checksum/secrets` annotation)

This ensures your collector always runs with the latest configuration and secrets.
