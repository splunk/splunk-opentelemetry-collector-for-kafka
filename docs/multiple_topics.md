## Collecting events from multiple topics

The configuration is similar to the default one described [here](../README.md#create-a-minimal-config-template), with the addition of multiple receivers - one for each topic you want to monitor. Thanks to the flexibility of the OpenTelemetry Collector, the setup can be tailored to meet specific requirements. This modular approach allows you to treat the components as building blocks, enabling you to create a pipeline that aligns perfectly with your use case. Depending on your needs, you can either use a single exporter for all receivers or configure a separate exporter for each receiver.

![SOC4Kafka multiple topics](images/kafka-multiple-topics.png)

### Example config

```yaml
receivers:
 kafka/example_topic:
   brokers: ["localhost:9092"]
   topic: "example-topic"
   encoding: "text"
 kafka/example_topic_2:
   brokers: ["localhost:9092"]
   topic: "example-topic-2"
   encoding: "text"

processors:
 batch:

exporters:
 splunk_hec/1:
   token: "your-splunk-hec-token"
   endpoint: "https://splunk-hec-endpoint:8088/services/collector"
   source: my-kafka
   sourcetype: kafka-otel
   index: kafka_otel
   headers:
     "__splunk_app_name": "soc4kafka"
 splunk_hec/2:
   token: "your-splunk-hec-token"
   endpoint: "https://splunk-hec-endpoint:8088/services/collector"
   source: my-kafka
   sourcetype: kafka-otel
   index: kafka_otel_another_index
   headers:
     "__splunk_app_name": "soc4kafka"

service:
 pipelines:
   logs/1:
     receivers: [kafka/example_topic]
     processors: [batch]
     exporters: [splunk_hec/1]
   logs/2:
     receivers: [kafka/example_topic_2]
     processors: [batch]
     exporters: [splunk_hec/2]
```