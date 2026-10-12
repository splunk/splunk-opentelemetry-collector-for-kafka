# Uninstall the Splunk Distribution of OpenTelemetry Collector for Kafka Helm release

## Uninstall the Helm release

```bash
helm uninstall soc4kafka
```

!!! note

    This command deletes the deployment and any Secrets created by the chart. Secrets created outside the chart remain.
