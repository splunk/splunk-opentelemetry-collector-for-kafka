# Upgrade the Collector for Kafka Helm release

## Upgrading

```bash
# Update your values.yaml file with new configuration, then upgrade
helm upgrade soc4kafka splunk-opentelemetry-collector-for-kafka/splunk-opentelemetry-collector-for-kafka -f values.yaml

# Or use multiple values files (useful for environment-specific overrides)
helm upgrade soc4kafka splunk-opentelemetry-collector-for-kafka/splunk-opentelemetry-collector-for-kafka -f values.yaml -f values-prod.yaml
```

!!! note
    **Best Practice:** Always use values files (`-f values.yaml`) instead of `--set` flags. This makes your configuration version-controlled, easier to maintain, and reusable across environments.

### Rolling updates (default behaviour)

By default, the chart uses a **rolling update** strategy (`maxSurge: 25%`, `maxUnavailable: 25%`). Pods are updated in waves so that the majority stay running during an upgrade.

!!! warning

    When you change collector configuration (for example, index, pipeline, or Splunk HTTP Event Collector (HEC) settings) and run `helm upgrade`, only a subset of pods receive the new config at a time. Until the rollout finishes, some pods still run with the old config. As a result, events from different Kafka partitions can be indexed or processed differently during the rollout (e.g. different index or sourcetype). With 25%, fewer partitions are affected in each wave. After all pods are updated, behaviour is consistent again.

If you need strictly sequential or consistent indexing during config changes, you can set `strategy.type: Recreate` in your values. That restarts all pods at once; expect a short period with no ingestion until the new pods are ready.
