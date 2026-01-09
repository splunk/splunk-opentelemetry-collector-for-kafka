## Migration from Splunk Connect for Kafka into Splunk OTel Collector for Kafka

Naming: 
- SC4Kafka - the [Splunk Connect for Kafka](https://github.com/splunk/kafka-connect-splunk)
- SOC4Kafka - the Splunk OTel Collector for Kafka (the current project)

The biggest difference between SC4Kafka and SOC4Kafka is that:

| **Field**                  | **SC4Kafka**                                                                          | **SOC4Kafka**                                                                                                                                                                                                    |
|----------------------------|---------------------------------------------------------------------------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| **Type**                   | Connector based on **Kafka Connect**, installed as an add-on for Kafka.               | Standalone product that works independently of Kafka.                                                                                                                                                            |
| **Message Retrieval**      | Retrieves events directly from Kafka.                                                 | Uses **REST calls** to Kafka (via the Kafka OpenTelemetry Receiver) to retrieve messages.                                                                                                                        |
| **Processing**             | Sends events directly to Splunk using the **Splunk HEC exporter**.                    | Processes messages internally and supports customization using **transform processors** before sending them to the **Splunk HEC exporter**.                                                                      |
| **Integration with Kafka** | Tightly integrated as part of the Kafka ecosystem.                                    | Can run independently and be deployed on an external server, separate from the Kafka cluster.                                                                                                                    |
| **Scaling**                | Scaling is managed using the `tasks.max` setting and supports multiple HEC endpoints. | Scaling is achieved by deploying multiple SOC4Kafka instances with the same `group_id`. Multiple HEC endpoints are not supported, but you can create multiple Splunk HEC exporters and add them to the pipeline. |

Note:
- The timestamp behavior differs between the two solutions. SOC4Kafka assigns a timestamp to the event based on when it is indexed, whereas Splunk Connect for Kafka uses the timestamp from when the event was originally produced. 
- Additionally messages from SOC4Kafka appear in Splunk first, as it forwards events to Splunk immediately. In contrast, Splunk Connect for Kafka processes and forwards events in batches, typically every configured number of seconds.

### SC4Kafka to SOC4Kafka mapping of configuration parameters

The configuration settings for SC4Kafka cannot be directly transferred to SOC4Kafka due to differences in their architecture and design. However, many configuration parameters have equivalent settings in SOC4Kafka.
A detailed description of the configuration parameter mappings can be found in the following [table](migration_config_values.md), which provides a comparison of the corresponding properties in both solutions.

## Migration process from SC4Kafka to SOC4Kafka

--- 
### Important Notes:
- **Migration from the old SC4Kafka connector to SOC4Kafka collector is a manual process.** There is no automated tool available for this migration.
- Begin with a simple configuration, then gradually add more settings. This approach helps in isolating and troubleshooting potential issues during the migration.

---

Migrating from SC4Kafka to SOC4Kafka involves several steps to ensure a smooth transition. Below are the key steps to follow during the migration process:
1. **Review Current SC4Kafka Configuration**: Start by thoroughly reviewing your existing SC4Kafka configuration. Document all the settings, including topics, indexes, sourcetypes, and any custom configurations you have in place. 
    In order to read the existing SC4Kafka configuration you can use REST API calls as described in the [section below](#reading-sc4kafka-existing-configuration).
2. **Map Configuration Parameters**: Use the [configuration mapping table](migration_config_values.md) to identify equivalent settings in SOC4Kafka. This will help you understand how to translate your SC4Kafka configuration into SOC4Kafka format.
3. **Create SOC4Kafka Configuration**: Based on the mapped parameters, create a new configuration file for SOC4Kafka. Make sure to include all relevant settings, such as Kafka brokers, topics, Splunk HEC endpoint, and token.
4. **Set Up SOC4Kafka**: [Install SOC4Kafka](../README.md/#how-to-start-with-soc4kafka) on your desired server. Ensure that you have the necessary permissions and access to both Kafka and Splunk.
5. **Test the Configuration**: Before fully switching over, test the SOC4Kafka configuration in a controlled environment. Verify that it can successfully connect to Kafka, retrieve messages, and send them to Splunk.
6. **Monitor and Validate**: Once you have deployed SOC4Kafka, closely monitor its performance and validate that all messages are being correctly forwarded to Splunk. Check for any discrepancies in data or performance issues.
7. **Decommission SC4Kafka**: After confirming that SOC4Kafka is functioning as expected, you can decommission your SC4Kafka setup. 


### Reading SC4Kafka existing configuration

#### Checking logs encoding format
When migrating from SC4Kafka to SOC4Kafka, it is important to consider the message format used in Kafka topics.
In case of SC4Kafka the default message format settings are stored in `connect-distributed.properties` file. The key
and value converter (`org.apache.kafka.connect.json.JsonConverter` or `org.apache.kafka.connect.storage.StringConverter`)
is a part of Kafka's java ecosystem, in SOC4Kafka it can be handled by setting `receivers.kafka.logs.encoding` to `json` or `text` depending on SC4Kafka configuration.

#### Reading SC4Kafka connector configuration
When migrating from SC4Kafka to SOC4Kafka following commands may be useful:

| Action | curl Command                                                   | Description |
|--------------------------------|----------------------------------------------------------------|----------------------------------------------|
| List active connectors | `curl http://localhost:8083/connectors`                        | Lists all active connectors |
| Get SC4Kafka connector info | `curl http://localhost:8083/connectors/<CONNECTOR_NAME>`       | Retrieves information about the specified SC4Kafka connector |
| Get SC4Kafka connector config | `curl http://localhost:8083/connectors/<CONNECTOR_NAME>/config` | Retrieves configuration details of the specified SC4Kafka connector |
| Get SC4Kafka connector task info | `curl http://localhost:8083/connectors/<CONNECTOR_NAME>/tasks`  | Retrieves task information for the specified SC4Kafka connector |


## Migration examples:

Following examples demonstrate how to migrate common SC4Kafka configurations to SOC4Kafka.
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
      topic: "three-pat"
      encoding: "text"

processors:
  batch:

exporters:
  splunk_hec:
    token: "your-splunk-hec-token"
    endpoint: "https://splunk-hec-endpoint:8088/services/collector"
    source: my-kafka
    sourcetype: kafka-otel
    index: "logs_index"
    headers:
      "__splunk_app_name": "soc4kafka"

service:
  pipelines:
    logs:
      receivers: [kafka]
      processors: [batch]
      exporters: [splunk_hec]
```

![1.png](images/migration/basic-message.png)

### Timestamp extraction

Even though by default kafka events from SOC4Kafka are marked with time of collecting data, if only we have a timestamp included as a part of the log body we can extract it. For example if the event is:

```
[2025-06-26 11:45:00]  the message with a timestamp
```

```yaml
receivers:
  kafka:
    brokers: ["kafka-broker:9092"]
    logs:
      topic: "three-pat"
      encoding: "text"

processors:
  transform:
    error_mode: ignore
    log_statements:
      - set(log.attributes["extracted_ts"], ExtractPatterns(log.body, "\\[(?P<timestamp>[0-9]{4}-[0-9]{2}-[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2})\\]"))
      - set(log.time, Time(log.attributes["extracted_ts"]["timestamp"], "2006-01-02 15:04:05", "UTC"))
      - delete_key(log.attributes, "extracted_ts")
  batch:

exporters:
  splunk_hec:
    token: "your-splunk-hec-token"
    endpoint: "https://splunk-hec-endpoint:8088/services/collector"
    source: my-kafka
    sourcetype: kafka-otel
    index: "logs_index"
    headers:
      "__splunk_app_name": "soc4kafka"

service:
  pipelines:
    logs:
      receivers: [kafka]
      processors: [batch, transform]
      exporters: [splunk_hec]
```

and the event in Splunk would be:

![3.png](images/migration/message-with-timestamp.png)

### Set host automatically

#### SOC4Kafka config

By default, events produced by SOC4Kafka may have the host field marked as `unknown`. This behavior can be adjusted using the [resource detection processor](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/processor/resourcedetectionprocessor).
The configuration example below demonstrates how to retrieve the hostname of the machine where SOC4Kafka is installed. Alternatively, the host value can be sourced from environmental variables or a specific API, depending on the client's requirements. The processor is flexible and can be tailored to meet specific use cases, as detailed in the [official documentation](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/processor/resourcedetectionprocessor\).
```yaml
receivers:
  kafka:
    brokers: ["kafka-broker:9092"]
    logs:
      topic: "three-pat"
      encoding: "text"

processors:
  resourcedetection:
    detectors: ["system"]
    system:
      hostname_sources: ["os"]
  batch:

exporters:
  splunk_hec:
    token: "your-splunk-hec-token"
    endpoint: "https://splunk-hec-endpoint:8088/services/collector"
    source: my-kafka
    sourcetype: kafka-otel
    index: "logs_index"
    headers:
      "__splunk_app_name": "soc4kafka"

service:
  telemetry:
    logs:
      level: "debug"
  pipelines:
   logs:
     receivers: [kafka]
     processors: [batch, resourcedetection]
     exporters: [splunk_hec]
```

![2.png](images/migration/message-with-host.png)

### Extract headers

If there are additional headers present in the incoming data, they can be extracted and added as event attributes. This allows for greater flexibility in customizing event metadata.
In the following examples, we will extract the following headers and include them as event attributes:
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
      topic: "three-pat"
      encoding: "text"
    header_extraction:
      extract_headers: true
      headers: ["index", "source", "sourcetype", "host", "myHeader1", "myHeader2"]

processors:
  batch:

exporters:
  splunk_hec:
    token: "your-splunk-hec-token"
    endpoint: "https://splunk-hec-endpoint:8088/services/collector"
    headers:
      "__splunk_app_name": "soc4kafka"
    otel_attrs_to_hec_metadata:
      index: kafka.header.index
      host: kafka.header.host
      source: kafka.header.source
      sourcetype: kafka.header.sourcetype

service:
  pipelines:
    logs:
      receivers: [kafka]
      processors: [batch]
      exporters: [splunk_hec]
```

This is how events generated by SC4Kafka are displayed in Splunk:

![sc4kafka-headers.png](images/migration/sc4kafka-headers.png)

Similarly, events generated by SOC4Kafka are presented in a comparable format:

![soc4kafka-headers.png](images/migration/soc4kafka-headers.png)

### Send data from multiple kafka topics to multiple Splunk HEC endpoints

In SC4Kafka, you can provide a list of topics along with a corresponding list of indexes, where each topic's data is mapped to its respective index (e.g., the first topic maps to the first index, the second topic to the second index, and so on).


In SOC4Kafka, the configuration is more flexible and modular. You define Kafka receivers and Splunk HEC exporters separately and connect them using a pipeline structure. Additionally, you can configure different source and sourcetype values directly within the settings of each Splunk HEC exporter, enabling greater customization for data routing and metadata assignment.

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
      topic: "three-pat"
      encoding: "text"

  kafka/2:
    brokers: ["kafka-broker:9092"]
    logs:
      topic: "two-pat"
      encoding: "text"

processors:
  resourcedetection:
    detectors: ["system"]
    system:
      hostname_sources: ["os"]
  batch:

exporters:
  splunk_hec/1:
    token: "your-splunk-hec-token"
    endpoint: "https://splunk-hec-endpoint:8088/services/collector"
    source: kafka-otel-three-pat
    sourcetype: kafka-otel
    index: "logs_index"
    headers:
      "__splunk_app_name": "soc4kafka"

  splunk_hec/2:
    token: "your-splunk-hec-token"
    endpoint: "https://splunk-hec-endpoint:8088/services/collector"
    source: kafka-otel-two-pat
    sourcetype: kafka-otel
    index: "kafka_otel"
    headers:
      "__splunk_app_name": "soc4kafka"

service:
  pipelines:
   logs/1:
     receivers: [kafka/1]
     processors: [batch, resourcedetection]
     exporters: [splunk_hec/1]
   logs/2:
     receivers: [kafka/2]
     processors: [batch, resourcedetection]
     exporters: [splunk_hec/2]
```

The events generated by SC4Kafka are:

![sc4kafka-two-pat.png](images/migration/sc4kafka-two-pat.png)
![sc4kafka-three-pat.png](images/migration/sc4kafka-three-pat.png)

While the events from SOC4Kafka are:

![soc4kafka-two-pat.png](images/migration/soc4kafka-two-pat.png)
![sock4kafka-three-pat.png](images/migration/sock4kafka-three-pat.png)

Mind that SOC4Kafka allows you to configure a unique sourcetype and source for each individual topic. This flexibility simplifies filtering and organizing data within Splunk, ensuring better control over your event categorization and search results.

### Sending events that are already in HEC format

In SC4Kafka you can collect events that are already formatted in HEC format, by setting `splunk.hec.json.event.formatted` option to `true`.

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
      topic: "topic"
      encoding: "text"

exporters:
  splunk_hec:
    token: "your-splunk-hec-token"
    endpoint: "https://splunk-hec-endpoint:8088/services/collector"
    source: otel
    sourcetype: otel
    index: test
    headers:
      "__splunk_app_name": "soc4kafka"
    export_raw: true

service:
  pipelines:
    logs:
      receivers: [kafka]
      exporters: [splunk_hec]
```

Example of event in this format:

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

![formatted-msg.png](images/migration/formatted-msg.png)
