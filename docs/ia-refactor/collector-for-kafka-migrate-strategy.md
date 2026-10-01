# Choose a migration strategy



There are several considerations to take into account when migrating from SC4Kafka to SOC4Kafka. 
During the transitional phase—when both SC4Kafka and SOC4Kafka are enabled and SC4Kafka has not yet been 
decommissioned, Kafka’s at‑least‑once delivery semantics may result in duplicate events being ingested into Splunk. 
This behavior is expected and is a consequence of Kafka’s design, which prioritizes data durability and avoidance of 
data loss over deduplication.

There are multiple approaches to performing the migration. For most of them, it is important to understand how Kafka consumer groups are used by both connectors.

### Check running consumer groups

First, identify the running SC4Kafka connectors:

```
curl http://localhost:8083/connectors
```

Each connector name returned by this command corresponds to a Kafka Connect connector instance. Kafka Connect automatically derives a consumer group ID using the following pattern: `connect-<CONNECTOR_NAME>`.
For example, if the connector list is:

```
["kafka-connect-splunk"]
```

The corresponding consumer group ID will be: `connect-kafka-connect-splunk`.

You can verify consumer group activity and partition assignments from the Kafka perspective using the built‑in `kafka-consumer-groups.sh` script located in Kafka’s `bin` directory:

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

This output confirms which partitions are currently assigned to SC4Kafka and whether offsets are being actively committed.

### Strategy 1: use different consumer groups

If no `group_id` is explicitly configured in SOC4Kafka, the connector uses its default consumer group ID: `otel_collector`. Because this consumer group ID differs from the one used by SC4Kafka, both connectors will independently consume the same Kafka topic.

As a result:

* SC4Kafka and SOC4Kafka each process all events from Kafka
* All events are ingested twice into Splunk
* Duplicate events will continue to appear until SC4Kafka is fully decommissioned

This approach is not recommended for production use unless duplicate ingestion is explicitly acceptable or temporary 
duplication is accounted for during the migration window.

### Strategy 2: use the same consumer group ID

In this strategy, SC4Kafka and SOC4Kafka are configured to use the same Kafka consumer group ID. This causes both connectors to participate in the same consumer group and share partition assignments rather than independently consuming all messages.

**Configuration**

Configure SOC4Kafka with the same `group_id` used by SC4Kafka. For example:

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

Expected Behavior

When both connectors are running with the same consumer group:

* Kafka assigns partitions to only one consumer per partition
* SC4Kafka and SOC4Kafka share consumption based on partition assignments
* Events are not duplicated by design
* Kafka determines partition ownership dynamically during group rebalances

However, Kafka consumer groups are designed for parallel work sharing, not for active/standby failover. So when you decommission
SC4Kafka with the same `group_id` configured,  SOC4Kafka will take over uncommitted partitions. SC4Kafka might replay offsets that it processed but had not committed. This replay can result in duplicate events in Splunk.

!!! note
    Kafka consumer groups are optimized for resilience and throughput, not for seamless connector replacement. This strategy reduces duplication compared to using separate consumer groups but does not eliminate it entirely.

Using the same consumer group ID for SC4Kafka and SOC4Kafka is the recommended migration approach when both connectors must temporarily coexist. This strategy minimizes duplicate ingestion compared to using separate consumer groups and allows for a controlled transition, provided that connector shutdown is carefully coordinated.

### Strategy 3: use separate topics (parallel topics migration)


In this strategy, new Kafka topics are introduced specifically for SOC4Kafka, while SC4Kafka continues to consume from the existing topics. Event producers are reconfigured to send data to the new topics, allowing both connectors to operate in parallel without sharing consumer groups or partitions.


**Configuration**

1. Create new Kafka topics for SOC4Kafka, typically by using a clear naming convention to distinguish them from existing topics. For example:

```yaml
kafka-topics.sh \
  --bootstrap-server localhost:9092 \
  --create \
  --topic topic2 \
  --partitions 10 \
  --replication-factor 1
```

2. Update Kafka producers to publish events to the new topics (e.g., `topic2`) instead of the original topics.
3. Configure SOC4Kafka to consume from the new topics:

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

4. Leave SC4Kafka configured to consume from the original topics (e.g., `topic1`).
5. Once all messages are ingested by SC4Kafka, you can safely stop the connector.
