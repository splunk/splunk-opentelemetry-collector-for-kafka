# Collect logs from the Collector for Kafka

## Write logs to a file or standard output

To collect logs from the Collector for Kafka, add a `telemetry` section under the `service` block in its configuration file. Use this section to set the logging level and log output paths. You can write logs to a file, standard output (`stdout`), or both.

!!! note

    Make sure that the specified log file paths exist and that the collector process can write to them.

This example creates a `soc4kafka-otel` directory in the current working directory. Create the directory before starting the collector.

```bash
mkdir -p ./soc4kafka-otel
```

Add the following `telemetry` section to your configuration file:

```yaml
service:
  telemetry:
    logs:
      level: info
      output_paths:
        - ./soc4kafka-otel/otel-collector.log
        - stdout
      error_output_paths:
        - ./soc4kafka-otel/otel-collector-errors.log
        - stderr
```
The following configuration writes Collector for Kafka logs and includes the `telemetry` section:

```yaml
receivers:
  kafka:
    brokers: ["kafka-broker-1:9092", "kafka-broker-2:9092"]
    logs:
      topics: 
        - example-topic
      encoding: text

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
  telemetry:
    logs:
      level: info
      output_paths:
        - ./soc4kafka-otel/otel-collector.log
        - stdout
      error_output_paths:
        - ./soc4kafka-otel/otel-collector-errors.log
        - stderr
  pipelines:
    logs:
      receivers: [kafka]
      processors: [resourcedetection]
      exporters: [splunk_hec]
```

## Forward logs from multiple collector instances
To monitor collector instances on multiple hosts, add the `filelog` receiver to forward log files to Splunk. This lets you monitor logs centrally instead of connecting to each VM over SSH.

The `file_storage` extension records the read position in each log file. After a restart, the collector resumes reading from that position instead of sending the entire file again.

The following configuration shows how to configure this setup:

```yaml
receivers:
  kafka:
    brokers: ["kafka-broker-1:9092", "kafka-broker-2:9092"]
    logs:
      topics: 
        - example-topic
      encoding: text

  filelog:
    include:
      - "./soc4kafka-otel/*.log"
    start_at: beginning
    storage: file_storage

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
    splunk_app_name: "soc4kafka"
    sending_queue:
      enabled: true
      num_consumers: 10
      queue_size: 10000
      block_on_overflow: true
      sizer: items
      batch:
        min_size: 1000

  splunk_hec/internal_logs:
    token: "your-splunk-hec-token"
    endpoint: "https://splunk-hec-endpoint:8088/services/collector"
    index: kafka-logs
    splunk_app_name: "soc4kafka"
    sending_queue:
      enabled: true
      num_consumers: 10
      queue_size: 10000
      block_on_overflow: true
      sizer: items
      batch:
        min_size: 1000

extensions:
  file_storage:
    directory: "./soc4kafka-otel/checkpoint"
    create_directory: true

service:
  extensions: [file_storage]
  telemetry:
    logs:
      level: info
      output_paths:
        - ./soc4kafka-otel/otel-collector.log
        - stdout
      error_output_paths:
        - ./soc4kafka-otel/otel-collector-errors.log
        - stderr
  pipelines:
    logs:
      receivers: [kafka]
      processors: [resourcedetection]
      exporters: [splunk_hec]
    logs/internal:
      receivers: [filelog]
      processors: [resourcedetection]
      exporters: [splunk_hec/internal_logs]
```
The `logs/internal` pipeline sends logs collected by the `filelog` receiver to Splunk through the `splunk_hec/internal_logs` exporter. Configure this exporter to send internal logs to a separate Splunk index for analysis.

## Add timestamps to collector log file names
Add timestamps to log file names to identify files created by different collector instances or restarts. Use an environment variable in the configuration file, as shown in this example:

1. Prepare the `values_timestamp.yaml` configuration file. Include the `${TIMESTAMP}` variable in the log file names:
```yaml
receivers:
  kafka:
    brokers: ["kafka-broker-1:9092", "kafka-broker-2:9092"]
    logs:
      topics: 
        - test-v
      encoding: text

  filelog:
    include:
      - "./soc4kafka-otel/*.log"
    start_at: beginning
    storage: file_storage
    include_file_path: true
    include_file_name: false

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
    splunk_app_name: "soc4kafka"
    sending_queue:
      enabled: true
      num_consumers: 10
      queue_size: 10000
      block_on_overflow: true
      sizer: items
      batch:
        min_size: 1000

  splunk_hec/internal_logs:
    token: "your-splunk-hec-token"
    endpoint: "https://splunk-hec-endpoint:8088/services/collector"
    index: kafka-logs
    splunk_app_name: "soc4kafka"
    sending_queue:
      enabled: true
      num_consumers: 10
      queue_size: 10000
      block_on_overflow: true
      sizer: items
      batch:
        min_size: 1000

extensions:
  file_storage:
    directory: "./soc4kafka-otel/checkpoint"
    create_directory: true

service:
  extensions: [file_storage]
  telemetry:
    logs:
      level: info
      output_paths:
        - ./soc4kafka-otel/otel-collector-${TIMESTAMP}.log
        - stdout
      error_output_paths:
        - ./soc4kafka-otel/otel-collector-errors-${TIMESTAMP}.log
        - stderr
  pipelines:
    logs:
      receivers: [kafka]
      processors: [resourcedetection]
      exporters: [splunk_hec]
    logs/internal:
      receivers: [filelog]
      processors: [resourcedetection]
      exporters: [splunk_hec/internal_logs]
```

2. Use a Bash script to set the current timestamp as an environment variable, then start the collector with this configuration file.
```bash
#!/bin/bash
# Create timestamp and export it
export TIMESTAMP=$(date +"%Y%m%d_%H%M%S")

# Create logs directory if it doesn't exist
mkdir -p ./soc4kafka-otel

# Start the collector with a config file that uses the TIMESTAMP variable
echo "Starting OpenTelemetry Collector with timestamp: $TIMESTAMP"
./<otel_package> --config values_timestamp.yaml
```
The script creates log files with timestamps in their names, such as `otel-collector-20231005_143200.log` and `otel-collector-errors-20231005_143200.log`.

## Collect logs with Helm

### Configure log collection in the Helm chart

Enable collection of the collector's internal logs. The chart writes logs to files in `/var/log/otelcol` by using an `emptyDir` volume. It also writes logs to `stdout` or `stderr` for Kubernetes log aggregation.

By default, when `collectorLogs.enabled` is `true`, the chart forwards these logs to Splunk by using a `filelog` receiver. This lets you monitor logs from multiple collector instances in one place.

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

The chart:

- Writes logs to files and to `stdout` or `stderr`.
- Forwards logs to Splunk by using the `filelog` receiver.
- Uses the `file_storage` extension to track read positions and avoid reading the same logs again after a restart.
- Sends internal logs by using the referenced Splunk exporter. By default, it uses the first exporter. You can configure a dedicated exporter to send logs to a separate index.
- Uses the first `splunkExporter` endpoint and Secret by default. You can override these values.

!!! note

    Log files are stored in an `emptyDir` volume and are deleted when the pod is deleted. The chart forwards the logs to Splunk, where they remain available.
