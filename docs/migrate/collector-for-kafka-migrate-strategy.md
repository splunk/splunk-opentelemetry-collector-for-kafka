# Choose a migration strategy



When you run Splunk Connect for Kafka and the Splunk Distribution of OpenTelemetry Collector for Kafka at the same time, Kafka's at-least-once delivery semantics can result in duplicate events in Splunk. This behavior is expected because Kafka prioritizes data durability over deduplication.

Choose a migration strategy based on how the two products use Kafka consumer groups.

## Check running consumer groups

Identify the running Splunk Connect for Kafka connectors:

```
curl http://localhost:8083/connectors
```

Each name in the output identifies a Kafka Connect connector. Kafka Connect derives its consumer group ID from the connector name by using this pattern: `connect-<CONNECTOR_NAME>`. For example, if the command returns:

```
["kafka-connect-splunk"]
```

The corresponding consumer group ID is `connect-kafka-connect-splunk`.

To check consumer group activity and partition assignments, run the `kafka-consumer-groups.sh` script in Kafka's `bin` directory:

```
./kafka-consumer-groups.sh \
  --bootstrap-server localhost:9092 \
  --describe \
  --group connect-kafka-connect-splunk
```

Example output:

```
./kafka-consumer-groups.sh   --bootstrap-server localhost:9092   --describe   --group connect-kafka-connect-splunk
GROUP                        TOPIC           PARTITION  CURRENT-OFFSET  LOG-END-OFFSET  LAG             CONSUMER-ID                                                                    HOST            CLIENT-ID
connect-kafka-connect-splunk topic1          0          3100135         3100135         0               connector-consumer-kafka-connect-splunk-0-a299967d-4ba2-4d8c-95f0-7f7db4f029ed /10.236.5.232   connector-consumer-kafka-connect-splunk-0
```

The output shows which partitions are assigned to Splunk Connect for Kafka and whether it is committing offsets.

## Use different consumer groups

If you do not configure `group_id`, the Splunk Distribution of OpenTelemetry Collector for Kafka uses `otel_collector` by default. Because this ID differs from the one used by Splunk Connect for Kafka, both products consume the same Kafka topic independently.

As a result:

* Both products process all events from Kafka.
* Splunk indexes each event twice.
* Duplicate events continue until you decommission Splunk Connect for Kafka.

Don't use this approach in production unless duplicate events are acceptable or you account for them during the migration.

## Use the same consumer group ID

Configure Splunk Connect for Kafka and the Splunk Distribution of OpenTelemetry Collector for Kafka to use the same Kafka consumer group ID. Kafka then assigns partitions across both products instead of having each product consume every message.

### Configure the consumer group ID

Set `group_id` in the Splunk Distribution of OpenTelemetry Collector for Kafka configuration to the ID used by Splunk Connect for Kafka. For example:

```yaml
receivers:
  kafka:
    brokers:
      - "kafka-broker:9092"
    logs:
      topics:
        - "topic1"
      encoding: "text"
    group_id: connect-kafka-connect-splunk
```

### Expected behavior

When both products use the same consumer group:

* Kafka assigns each partition to one consumer.
* The products share consumption according to their partition assignments.
* Kafka changes partition assignments during group rebalances.

Kafka consumer groups distribute work; they do not provide active-standby failover. When you decommission Splunk Connect for Kafka, the Splunk Distribution of OpenTelemetry Collector for Kafka takes over its uncommitted partitions. Splunk Connect for Kafka might replay offsets that it processed but had not committed, which can result in duplicate events in Splunk.

!!! note
    Kafka consumer groups support resilience and throughput, but they do not provide seamless connector replacement. This strategy reduces duplicate events compared with separate consumer groups, but it does not eliminate them.

Using the same consumer group ID is the recommended strategy when both products must temporarily coexist. It reduces duplicate events compared with separate consumer groups and supports a controlled transition. Coordinate connector shutdown to reduce the chance of replayed events.

## Use separate topics for a parallel migration

Create new Kafka topics for the Splunk Distribution of OpenTelemetry Collector for Kafka while Splunk Connect for Kafka continues to consume from the existing topics. Configure event producers to send data to the new topics. This lets both products run in parallel without sharing consumer groups or partitions.


**Configuration**

1. Create Kafka topics for the Splunk Distribution of OpenTelemetry Collector for Kafka. Use a naming convention that distinguishes them from existing topics. For example:

```yaml
kafka-topics.sh \
  --bootstrap-server localhost:9092 \
  --create \
  --topic topic2 \
  --partitions 10 \
  --replication-factor 1
```

2. Update Kafka producers to publish events to the new topics, such as `topic2`, instead of the original topics.
3. Configure the Splunk Distribution of OpenTelemetry Collector for Kafka to consume from the new topics:

```yaml
receivers:
  kafka:
    brokers:
      - "kafka-broker:9092"
    logs:
      topics:
        - "topic2"
      encoding: "text"
```

4. Keep Splunk Connect for Kafka configured to consume from the original topics, such as `topic1`.
5. After Splunk Connect for Kafka collects all messages, stop the connector.
