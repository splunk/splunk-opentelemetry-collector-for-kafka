# Upgrade the Splunk Distribution of OpenTelemetry Collector for Kafka Helm release

## Upgrade the Helm release

```bash
# Update your values.yaml file with new configuration, then upgrade
helm upgrade soc4kafka splunk-opentelemetry-collector-for-kafka/splunk-opentelemetry-collector-for-kafka -f values.yaml

# Or use multiple values files (useful for environment-specific overrides)
helm upgrade soc4kafka splunk-opentelemetry-collector-for-kafka/splunk-opentelemetry-collector-for-kafka -f values.yaml -f values-prod.yaml
```

!!! note
    Use values files, such as `-f values.yaml`, instead of `--set` flags. This keeps your configuration under version control and lets you reuse it across environments.

### Rolling updates (default behavior)

By default, the chart uses a **rolling update** strategy (`maxSurge: 25%`, `maxUnavailable: 25%`). It updates pods in waves so that most pods keep running during an upgrade.

!!! warning

    When you change the collector configuration, such as the index, pipeline, or Splunk HTTP Event Collector (HEC) settings, and run `helm upgrade`, only some pods receive the new configuration at a time. Until the rollout finishes, other pods continue to use the old configuration. Events from different Kafka partitions can therefore be indexed or processed differently during the rollout, for example, with different indexes or sourcetypes. After all pods are updated, the configuration is consistent.

To apply configuration changes sequentially and keep indexing consistent, set `strategy.type: Recreate` in your values file. This restarts all pods at once. Data collection pauses until the new pods are ready.
