## Scaling SOC4Kafka

To handle higher throughput, you can deploy multiple instances of the SOC4Kafka. Kafka natively supports partition-based scaling, allowing multiple consumers within the same consumer group to share the workload.

### Steps to scale horizontally

#### Configure Kafka partitions

1. Ensure the Kafka topics you are consuming from are partitioned appropriately.
2. The number of partitions should match or exceed the number of collector instances to ensure even distribution.

#### Use the same consumer group

Configure all SOC4Kafka collectors to use the same `group_id`. Kafka ensures that each partition is consumed by only one collector instance within a consumer group.

```yaml
receivers:
  kafka:
    brokers: ["localhost:9092"]
    logs:
      topics:
        - "example-topic"
      encoding: "text"
    group_id: <GROUP ID>

processors:
 batch:

exporters:
  splunk_hec:
    token: "your-splunk-hec-token"
    endpoint: "https://splunk-hec-endpoint:8088/services/collector"
    source: my-kafka
    sourcetype: kafka-otel
    index: kafka_otel
    splunk_app_name: "soc4kafka"
    splunk_app_version: 0.145.0

service:
  pipelines:
    logs:
      receivers: [kafka]
      processors: [batch]
      exporters: [splunk_hec]
```

Note: Replace `<GROUP ID>` with a name that will be shared across all SOC4Kafka instances. This ensures that all instances are part of the same consumer group.