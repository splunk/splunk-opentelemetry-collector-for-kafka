# Installation Guide

## Quick Start

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

pipelines:
  - name: main-logs
    type: logs
    receivers:
      - main
    exporters:
      - primary
    # processors optional; defaults to ["batch", "resourcedetection"] (defaults.pipelineProcessors in values.yaml)
```

2. Add helm repository (TBD)

3. Install the chart:

```bash
helm upgrade --install soc4kafka splunk-opentelemetry-collector-for-kafka/splunk-opentelemetry-collector-for-kafka -f values.yaml
```

**Note:** For information about managing secrets (auto-created or existing Kubernetes secrets), see [Secret Management](secrets.md).

## Upgrading

```bash
# Update your values.yaml file with new configuration, then upgrade
helm upgrade soc4kafka splunk-opentelemetry-collector-for-kafka/splunk-opentelemetry-collector-for-kafka -f values.yaml

# Or use multiple values files (useful for environment-specific overrides)
helm upgrade soc4kafka splunk-opentelemetry-collector-for-kafka/splunk-opentelemetry-collector-for-kafka -f values.yaml -f values-prod.yaml
```

**Best Practice:** Always use values files (`-f values.yaml`) instead of `--set` flags. This makes your configuration version-controlled, easier to maintain, and reusable across environments.

### Rolling updates (default behaviour)

By default, the chart uses a **rolling update** strategy (`maxSurge: 25%`, `maxUnavailable: 25%`). Pods are updated in waves so that the majority stay running during an upgrade.

**Important:** When you change collector configuration (for example, index, pipeline, or Splunk HEC settings) and run `helm upgrade`, only a subset of pods receive the new config at a time. Until the rollout finishes, some pods still run with the old config. As a result, events from different Kafka partitions can be indexed or processed differently during the rollout (e.g. different index or sourcetype). With 25%, fewer partitions are affected in each wave. After all pods are updated, behaviour is consistent again.

If you need strictly sequential or consistent indexing during config changes, you can set `strategy.type: Recreate` in your values. That restarts all pods at once; expect a short period with no ingestion until the new pods are ready.

## Uninstallation

```bash
helm uninstall soc4kafka
```

**Note:** This will delete the deployment, but secrets created outside the chart will remain. Auto-created secrets will be deleted.
