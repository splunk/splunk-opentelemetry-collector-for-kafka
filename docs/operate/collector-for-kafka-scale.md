# Scale the Collector for Kafka

To increase throughput, deploy multiple instances of the Collector for Kafka. Kafka distributes partitions among consumers in the same consumer group.

## Scale horizontally

### Configure Kafka partitions

1. Make sure that each Kafka topic has enough partitions for the expected number of collector instances.
2. Set the number of partitions to match or exceed the number of collector instances to distribute the workload across instances.

### Use the same consumer group

Configure each instance of the Collector for Kafka to use the same `group_id`. Kafka assigns each partition to one consumer in the group.

```yaml
receivers:
  kafka:
    brokers: ["localhost:9092"]
    logs:
      topics:
        - "example-topic"
      encoding: "text"
    group_id: <GROUP ID>

exporters:
  splunk_hec:
    token: "your-splunk-hec-token"
    endpoint: "https://splunk-hec-endpoint:8088/services/collector"
    source: my-kafka
    sourcetype: kafka-otel
    index: kafka_otel
    splunk_app_name: "soc4kafka"
    sending_queue:
      enabled: true
      num_consumers: 10
      queue_size: 10000
      block_on_overflow: true
      sizer: items
      batch:
        min_size: 1000

service:
  pipelines:
    logs:
      receivers: [kafka]
      exporters: [splunk_hec]
```

!!! note
    Replace `<GROUP ID>` with a name shared by all instances of the Collector for Kafka. This setting puts all instances in the same consumer group.
