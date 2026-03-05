# Splunk OpenTelemetry Collector for Kafka - Helm Chart

Helm chart for deploying the Splunk OpenTelemetry Collector for Kafka (SOC4Kafka) on Kubernetes.

## Quick Start

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
    token: "my-splunk-hec-secret"
    source: "soc4kafka"
    sourcetype: "otel:logs"
    index: "main"

pipelines:
  - name: main-logs
    type: logs
    receivers:
      - main
    exporters:
      - primary
    processors:
      - batch
      - resourcedetection
```

```bash
helm upgrade --install soc4kafka splunk-opentelemetry-collector-for-kafka/splunk-opentelemetry-collector-for-kafka -f values.yaml
```

## Documentation

- [Installation Guide](docs/installation.md) - How to install and upgrade the chart
- [Configuration](docs/configuration.md) - Detailed configuration options
- [Examples](docs/examples.md) - Example configurations for common scenarios
- [Secret Management](docs/secrets.md) - Managing secrets for Splunk HEC tokens and Kafka authentication
- [Troubleshooting](docs/troubleshooting.md) - Common issues and solutions

## Features

- **Simple Configuration** - Define Kafka receivers, Splunk exporters, and pipelines in `values.yaml`
- **Secret Management** - Auto-create secrets or reference existing Kubernetes secrets
- **Collector Logs** - Optional collection and forwarding of collector's own logs
- **Metrics Collection** - Optional collection of collector internal and system metrics
- **Automatic Pod Restarts** - Pods restart automatically when configuration or secrets change
- **Production Ready** - Supports HPA, PDB, PVC, and resource limits
