# Splunk OpenTelemetry Collector for Kafka - Helm Chart

A Helm chart for deploying the Splunk OpenTelemetry Collector for Kafka (SOC4Kafka) on Kubernetes. This chart collects logs from Kafka topics and forwards them to Splunk via the HTTP Event Collector (HEC).

## Quick Start

```bash
helm install soc4kafka . -f values.yaml
```

## Examples

### Minimal Example

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
    secret: "splunk-hec-secret"
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

### With Custom Defaults and Memory Limiter

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
    secret: "splunk-hec-secret"
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
      - memory_limiter
      - batch
      - resourcedetection

defaults:
  processors:
    memory_limiter:
      limit_mib: 400
    batch:
      timeout: 10s
      send_batch_size: 2000
```

### With Config Override

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
    secret: "splunk-hec-secret"
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

configOverride:
  processors:
    batch:
      timeout: 15s
      send_batch_size: 5000
  exporters:
    splunk_hec/primary:
      compression: gzip
```

### Configuration Precedence

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

## Secrets

Using Kubernetes Secrets is recommended for storing sensitive credentials like Splunk HEC tokens and Kafka passwords. This keeps credentials out of your values.yaml files, allows for easier rotation, and integrates with external secret management systems like Vault.

### Splunk HEC Token

**Create secret:**
```bash
kubectl create secret generic splunk-hec-secret \
  --from-literal=splunk-hec-token=YOUR_HEC_TOKEN
```

**Reference in values.yaml:**
```yaml
splunkExporters:
  - name: primary
    secret: "splunk-hec-secret"  # Must have key "splunk-hec-token"
    endpoint: "https://splunk-hec:8088/services/collector"
```

**Alternative:** Provide `token` directly in values.yaml and the chart will auto-create the secret.

### Kafka Authentication

**Create secret:**
```bash
kubectl create secret generic kafka-auth-secret \
  --from-literal=password=your-kafka-password
```

**Reference in values.yaml:**
```yaml
kafkaReceivers:
  - name: main
    brokers:
      - "kafka-broker:9092"
    auth:
      plain_text:
        username: "kafka-user"
        secret: "kafka-auth-secret"  # Must have key "password"
      # Or for SASL/Kerberos:
      # sasl:
      #   username: "kafka-user"
      #   secret: "kafka-sasl-secret"
      # kerberos:
      #   secret: "kafka-kerberos-secret"
```

## Values Reference

See [values.yaml](values.yaml) for all available configuration options.
