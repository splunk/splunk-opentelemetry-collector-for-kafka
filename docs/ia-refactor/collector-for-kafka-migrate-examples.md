# Migration examples



The following examples show how to migrate common SC4Kafka configurations to SOC4Kafka.
The section contains examples for:

- Basic config for Kafka string messages
- Timestamp extraction
- Set host automatically
- Extract headers
- Send data from multiple kafka topics to multiple Splunk HEC endpoints
- Sending events that are already in HEC format

---

### The basic config for Kafka string messages

#### SC4Kafka config

```
curl localhost:8083/connectors -X POST -H "Content-Type: application/json" -d '{
    "name": "kafka-connect-splunk",
    "config": {
      "connector.class": "com.splunk.kafka.connect.SplunkSinkConnector",
      "tasks.max": "3",
      "splunk.indexes": "logs_index",
      "topics":"three-pat",
      "splunk.hec.uri": "https://splunk-hec-endpoint:8088",
      "splunk.hec.token": "your-splunk-hec-token"
    }
  }'
```

#### SOC4Kafka config

```yaml
receivers:
  kafka:
    brokers: ["kafka-broker:9092"]
    logs:
      topics:
        - "three-pat"
      encoding: "text"


exporters:
  splunk_hec:
    token: "your-splunk-hec-token"
    endpoint: "https://splunk-hec-endpoint:8088/services/collector"
    source: my-kafka
    sourcetype: kafka-otel
    index: "logs_index"
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

![Splunk search result for a message from the basic Kafka configuration](../images/migration/basic-message.png)

### Timestamp extraction

By default, SOC4Kafka assigns events the time when it collects them. To use a timestamp from the log body instead, extract it with a transform processor. For example, consider this event:

```
[2025-06-26 11:45:00]  the message with a timestamp
```

```yaml
receivers:
  kafka:
    brokers: ["kafka-broker:9092"]
    logs:
      topics:
        - "three-pat"
      encoding: "text"

processors:
  transform:
    error_mode: ignore
    log_statements:
      - set(log.attributes["extracted_ts"], ExtractPatterns(log.body, "\\[(?P<timestamp>[0-9]{4}-[0-9]{2}-[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2})\\]"))
      - set(log.time, Time(log.attributes["extracted_ts"]["timestamp"], "%Y-%m-%d %H:%M:%S", "UTC"))
      - delete_key(log.attributes, "extracted_ts")

exporters:
  splunk_hec:
    token: "your-splunk-hec-token"
    endpoint: "https://splunk-hec-endpoint:8088/services/collector"
    source: my-kafka
    sourcetype: kafka-otel
    index: "logs_index"
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
      processors: [transform]
      exporters: [splunk_hec]
```

and the event in Splunk would be:

![Splunk search result showing the timestamp extracted from a Kafka message](../images/migration/message-with-timestamp.png)

### Set host automatically

#### SOC4Kafka config

By default, events produced by SOC4Kafka may have the host field marked as `unknown`. This behavior can be adjusted using the [resource detection processor](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/processor/resourcedetectionprocessor).
The configuration example below demonstrates how to retrieve the hostname of the machine where SOC4Kafka is installed. Alternatively, the host value can be sourced from environmental variables or a specific API, depending on the client's requirements. The processor is flexible and can be tailored to meet specific use cases, as detailed in the [official documentation](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/processor/resourcedetectionprocessor).
```yaml
receivers:
  kafka:
    brokers: ["kafka-broker:9092"]
    logs:
      topics:
        - "three-pat"
      encoding: "text"

processors:
  resourcedetection:
    detectors: ["system"]
    system:
      hostname_sources: ["os"]

exporters:
  splunk_hec:
    token: "your-splunk-hec-token"
    endpoint: "https://splunk-hec-endpoint:8088/services/collector"
    source: my-kafka
    sourcetype: kafka-otel
    index: "logs_index"
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
  telemetry:
    logs:
      level: "debug"
  pipelines:
   logs:
     receivers: [kafka]
     processors: [resourcedetection]
     exporters: [splunk_hec]
```

![Splunk search result showing the host detected by SOC4Kafka](../images/migration/message-with-host.png)

### Extract headers

If incoming data includes additional headers, you can extract them as event attributes. The following examples extract these headers:

- index
- source
- sourcetype
- host
- myHeader1
- myHeader2

#### SC4Kafka config

```
curl localhost:8083/connectors -X POST -H "Content-Type: application/json" -d '{
    "name": "kafka-connect-splunk",
    "config": {
      "connector.class": "com.splunk.kafka.connect.SplunkSinkConnector",
      "tasks.max": "3",
      "splunk.indexes": "logs_index",
      "topics":"three-pat",
      "splunk.hec.uri": "https://splunk-hec-endpoint:8088",
      "splunk.hec.token": "your-splunk-hec-token",
      "splunk.header.index": "index",
      "splunk.header.source": "source",
      "splunk.header.sourcetype": "sourcetype",
      "splunk.header.host": "host",
      "splunk.header.custom": "myHeader1,myHeader2"
    }
  }'
```

#### SOC4Kafka config

```yaml
receivers:
  kafka:
    brokers: ["kafka-broker:9092"]
    logs:
      topics:
        - "three-pat"
      encoding: "text"
    header_extraction:
      extract_headers: true
      headers: ["index", "source", "sourcetype", "host", "myHeader1", "myHeader2"]


exporters:
  splunk_hec:
    token: "your-splunk-hec-token"
    endpoint: "https://splunk-hec-endpoint:8088/services/collector"
    splunk_app_name: "soc4kafka"
    sending_queue:
      enabled: true
      num_consumers: 10
      queue_size: 10000
      block_on_overflow: true
      sizer: items
      batch:
        min_size: 1000
    otel_attrs_to_hec_metadata:
      index: kafka.header.index
      host: kafka.header.host
      source: kafka.header.source
      sourcetype: kafka.header.sourcetype

service:
  pipelines:
    logs:
      receivers: [kafka]
      exporters: [splunk_hec]
```

This is how events generated by SC4Kafka are displayed in Splunk:

![Splunk event with headers extracted by Splunk Connect for Kafka](../images/migration/sc4kafka-headers.png)

Similarly, events generated by SOC4Kafka are presented in a comparable format:

![Splunk event with headers extracted by SOC4Kafka](../images/migration/soc4kafka-headers.png)

### Send data from multiple Kafka topics to multiple Splunk HEC endpoints

In SC4Kafka, you can provide a list of topics along with a corresponding list of indexes, where each topic's data is mapped to its respective index (e.g., the first topic maps to the first index, the second topic to the second index, and so on).


In SOC4Kafka, configure Kafka receivers and Splunk HEC exporters separately, then connect them in a pipeline. Each exporter can use different source and sourcetype values.

#### SC4Kafka config

```
curl localhost:8083/connectors -X POST -H "Content-Type: application/json" -d '{
    "name": "kafka-connect-splunk",
    "config": {
      "connector.class": "com.splunk.kafka.connect.SplunkSinkConnector",
      "tasks.max": "3",
      "splunk.indexes": "logs_index,kafka_otel",
      "topics":"three-pat,two-pat",
      "splunk.hec.uri": "https://splunk-hec-endpoint:8088",
      "splunk.hec.token": "your-splunk-hec-token",
    }
  }'
```

#### SOC4Kafka

```yaml
receivers:
  kafka/1:
    brokers: ["kafka-broker:9092"]
    logs:
      topics:
        - "three-pat"
      encoding: "text"

  kafka/2:
    brokers: ["kafka-broker:9092"]
    logs:
      topics:
        - "two-pat"
      encoding: "text"

processors:
  resourcedetection:
    detectors: ["system"]
    system:
      hostname_sources: ["os"]

exporters:
  splunk_hec/1:
    token: "your-splunk-hec-token"
    endpoint: "https://splunk-hec-endpoint:8088/services/collector"
    source: kafka-otel-three-pat
    sourcetype: kafka-otel
    index: "logs_index"
    splunk_app_name: "soc4kafka"
    sending_queue:
      enabled: true
      num_consumers: 10
      queue_size: 10000
      block_on_overflow: true
      sizer: items
      batch:
        min_size: 1000

  splunk_hec/2:
    token: "your-splunk-hec-token"
    endpoint: "https://splunk-hec-endpoint:8088/services/collector"
    source: kafka-otel-two-pat
    sourcetype: kafka-otel
    index: "kafka_otel"
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
   logs/1:
     receivers: [kafka/1]
     processors: [resourcedetection]
     exporters: [splunk_hec/1]
   logs/2:
     receivers: [kafka/2]
     processors: [resourcedetection]
     exporters: [splunk_hec/2]
```

The events generated by SC4Kafka are:

![Splunk Connect for Kafka event from the two-pat topic](../images/migration/sc4kafka-two-pat.png)
![Splunk Connect for Kafka event from the three-pat topic](../images/migration/sc4kafka-three-pat.png)

While the events from SOC4Kafka are:

![SOC4Kafka event from the two-pat topic](../images/migration/soc4kafka-two-pat.png)
![SOC4Kafka event from the three-pat topic](../images/migration/sock4kafka-three-pat.png)

SOC4Kafka lets you configure a unique source and sourcetype for each topic. Use these values to filter and organize events in Splunk.

### Sending events that are already in HEC format

To collect events that already use HEC format in SC4Kafka, set `splunk.hec.json.event.formatted` to `true`.

#### SC4Kafka config

```
curl localhost:8083/connectors -X POST -H "Content-Type: application/json" -d' {
    "name": "splunk-prod-financial",
      "config": {
        "connector.class": "com.splunk.kafka.connect.SplunkSinkConnector",
        "tasks.max": "20",
        "topics": "t1",
        "splunk.hec.uri": "https://idx1:8088,https://idx2:8088,https://idx3:8088",
        "splunk.hec.token": "your-splunk-hec-token",
        "splunk.hec.json.event.formatted": "true",
        "key.converter": "org.apache.kafka.connect.storage.StringConverter",
        "key.converter.schemas.enable": "false",
        "value.converter": "org.apache.kafka.connect.storage.StringConverter",
        "value.converter.schemas.enable": "false"
 }
 }'
```

#### SOC4Kafka

To achieve the same result in SOC4Kafka use `export_raw` option in exporter configuration:

```yaml
receivers:
  kafka:
    brokers: ["kafka-broker:9092"]
    logs:
      topics:
        - "topic"
      encoding: "text"

exporters:
  splunk_hec:
    token: "your-splunk-hec-token"
    endpoint: "https://splunk-hec-endpoint:8088/services/collector"
    source: otel
    sourcetype: otel
    index: test
    splunk_app_name: "soc4kafka"
    sending_queue:
      enabled: true
      num_consumers: 10
      queue_size: 10000
      block_on_overflow: true
      sizer: items
      batch:
        min_size: 1000
    export_raw: true

service:
  pipelines:
    logs:
      receivers: [kafka]
      exporters: [splunk_hec]
```

The following event uses this format:

```json
{
  "index":"test",
  "host":"localhost",
  "sourcetype":"sourcetype",
  "source":"source",
  "event":"This is already formatted event!",
  "fields":
  {
    "extra_field":"extra-field-1",
    "extra_fields_arr":
    [
      "extra-field-2",
      "extra-field-3"
    ]
  }
}
```

The example message appears like this in Splunk search results when properly configured:

![Splunk search result for an event already formatted for HEC](../images/migration/formatted-msg.png)
