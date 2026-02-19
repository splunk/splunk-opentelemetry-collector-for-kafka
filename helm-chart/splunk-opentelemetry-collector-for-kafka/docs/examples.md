# Examples

## Basic Single Receiver/Exporter

```yaml
kafkaReceivers:
  - name: main
    brokers:
      - "kafka:9092"
    logs:
      topics:
        - "logs"
      encoding: text
    group_id: "soc4kafka"

splunkExporters:
  - name: primary
    endpoint: "https://splunk:8088/services/collector"
    secret: "splunk-hec-secret"
    source: "kafka"
    sourcetype: "otel:logs"
    index: "main"

pipelines:
  - name: logs
    type: logs
    receivers:
      - main
    exporters:
      - primary
    processors:
      - batch
      - resourcedetection
```

## Multiple Topics to Multiple Indexes

```yaml
kafkaReceivers:
  - name: app-logs
    brokers:
      - "kafka:9092"
    logs:
      topics:
        - "app-info"
        - "app-warn"
      encoding: json
    group_id: "soc4kafka-app"
  
  - name: error-logs
    brokers:
      - "kafka:9092"
    logs:
      topics:
        - "app-error"
        - "app-critical"
      encoding: json
    group_id: "soc4kafka-error"

splunkExporters:
  - name: main-index
    endpoint: "https://splunk:8088/services/collector"
    secret: "splunk-hec-main"
    source: "kafka"
    sourcetype: "otel:logs"
    index: "main"
  
  - name: error-index
    endpoint: "https://splunk:8088/services/collector"
    secret: "splunk-hec-error"
    source: "kafka"
    sourcetype: "otel:logs"
    index: "errors"

pipelines:
  - name: app-pipeline
    type: logs
    receivers:
      - app-logs
    exporters:
      - main-index
    processors:
      - batch
      - resourcedetection
  
  - name: error-pipeline
    type: logs
    receivers:
      - error-logs
    exporters:
      - error-index
    processors:
      - batch
      - resourcedetection
```

## Authenticated Kafka with Secret Management

```yaml
kafkaReceivers:
  - name: secure-main
    brokers:
      - "secure-kafka:9092"
    auth:
      plain_text:
        username: "kafka-user"
        secret: "kafka-auth-secret"
    logs:
      topics:
        - "secure-logs"
      encoding: json
    group_id: "soc4kafka-secure"

splunkExporters:
  - name: primary
    endpoint: "https://splunk:8088/services/collector"
    secret: "splunk-hec-secret"
    source: "kafka"
    sourcetype: "otel:logs"
    index: "main"

pipelines:
  - name: secure-logs
    type: logs
    receivers:
      - secure-main
    exporters:
      - primary
    processors:
      - batch
      - resourcedetection
```
