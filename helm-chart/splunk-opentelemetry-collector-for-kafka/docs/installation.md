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
    secret: "my-splunk-hec-secret"
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

2. Add helm repository <TBD>

TBD: Replace . to be the helm repository
3. Install the chart:

```bash
helm upgrade --install soc4kafka . -f values.yaml 
```

## Using Existing Secrets

If you prefer to manage secrets separately:

1. Create a secret for Splunk HEC token:
   ```bash
   kubectl create secret generic my-splunk-hec-secret \
     --from-literal=splunk-hec-token=YOUR_HEC_TOKEN
   ```

2. Reference the secret in your `values.yaml` (see example above).

3. Install the chart:
   ```bash
   helm install soc4kafka . -f values.yaml
   ```

## Upgrading

```bash
# Update your values.yaml file with new configuration, then upgrade
helm upgrade soc4kafka . -f values.yaml

# Or use multiple values files (useful for environment-specific overrides)
helm upgrade soc4kafka . -f values.yaml -f values-prod.yaml
```

**Best Practice:** Always use values files (`-f values.yaml`) instead of `--set` flags. This makes your configuration version-controlled, easier to maintain, and reusable across environments.

## Uninstallation

```bash
helm uninstall soc4kafka
```

**Note:** This will delete the deployment, but secrets created outside the chart will remain. Auto-created secrets will be deleted.
