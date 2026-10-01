# Uninstall the Collector for Kafka Helm release

## Uninstallation

```bash
helm uninstall soc4kafka
```

!!! note

    This will delete the deployment, but secrets created outside the chart will remain. Auto-created secrets will be deleted.
