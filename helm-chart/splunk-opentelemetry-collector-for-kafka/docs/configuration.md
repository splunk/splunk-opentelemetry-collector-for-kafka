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

## Automatic Pod Restarts

The chart includes automatic pod restart triggers:

- **ConfigMap changes**: Pods restart when the OpenTelemetry configuration changes (via `checksum/config` annotation)
- **Secret changes**: Pods restart when:
  - Token values in `values.yaml` change (for auto-created secrets)
  - Secret references change (via `checksum/secrets` annotation)

This ensures your collector always runs with the latest configuration and secrets.
