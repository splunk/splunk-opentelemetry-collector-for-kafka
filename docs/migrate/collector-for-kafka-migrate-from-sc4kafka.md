# Migrate from Splunk Connect for Kafka

In this topic, **Splunk Connect for Kafka** refers to the Kafka Connect product. **Collector for Kafka** refers to the Splunk Distribution of OpenTelemetry Collector for Kafka.

The following table summarizes the differences between the two products:

| Field                  | Splunk Connect for Kafka                                                                          | Collector for Kafka                                                                                                                                                                                                    |
|----------------------------|---------------------------------------------------------------------------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Type                   | A connector based on Kafka Connect that you install as an add-on for Kafka.               | A standalone product that runs independently of Kafka.                                                                                                                                                            |
| Message retrieval      | Retrieves events directly from Kafka.                                                 | Consumes messages from Kafka by using the Kafka receiver in the OpenTelemetry Collector.                                                                                                             |
| Processing             | Sends events directly to Splunk by using the Splunk HTTP Event Collector (HEC) exporter.                    | Processes messages and supports customization with transform processors before sending them to the Splunk HEC exporter.                                                                      |
| Integration with Kafka | Runs as part of the Kafka ecosystem.                                    | Runs independently and can be deployed on a server outside the Kafka cluster.                                                                                                                                                            |
| Scaling                | Uses the `tasks.max` setting and supports multiple HEC endpoints. | Scale by deploying multiple instances of the Collector for Kafka with the same `group_id`. Multiple HEC endpoints are not supported, but you can add multiple Splunk HEC exporters to the pipeline. |

!!! note
    - Timestamp behavior differs between the products. The Collector for Kafka assigns an event timestamp based on when it indexes the event. Splunk Connect for Kafka uses the timestamp from when the event was produced.
    - The Collector for Kafka forwards events to Splunk immediately. Splunk Connect for Kafka processes and forwards events in batches, typically at a configured interval.

## Migrate your deployment

### Prepare the configuration

Before you migrate, review the [migration strategies](collector-for-kafka-migrate-strategy.md).

!!! note
    You must migrate from Splunk Connect for Kafka to the Collector for Kafka manually. No automated migration tool is available. Start with a minimal configuration, then add settings incrementally to help isolate and troubleshoot issues.

Follow these steps to migrate from Splunk Connect for Kafka to the Collector for Kafka:

1. **Review the existing Splunk Connect for Kafka configuration.** Record the topics, indexes, sourcetypes, and any custom settings. To read the configuration, use the REST API commands in [Read the existing Splunk Connect for Kafka configuration](#read-the-existing-splunk-connect-for-kafka-configuration).
2. **Map the configuration settings.** Use the [configuration mapping table](collector-for-kafka-migrate-config.md) to find equivalent settings for the Collector for Kafka.
3. **Create a configuration for the Collector for Kafka.** Include the required Kafka brokers, topics, Splunk HEC endpoint, and token.
4. **Install the Collector for Kafka.** Follow the instructions for [installing the Collector for Kafka with Helm](../deploy/collector-for-kafka-kubernetes-install-helm.md). Make sure that it can access Kafka and Splunk.
5. **Test the configuration.** In a test environment, confirm that the Collector for Kafka connects to Kafka, collects messages, and sends them to Splunk.
6. **Monitor and validate the deployment.** Confirm that the Collector for Kafka forwards all messages. Check for data discrepancies and performance issues.
7. **Decommission Splunk Connect for Kafka.** After you confirm that the Collector for Kafka works as expected, decommission the Splunk Connect for Kafka deployment.


### Read the existing Splunk Connect for Kafka configuration

#### Check the message encoding
When you migrate, account for the message format used in each Kafka topic. Splunk Connect for Kafka stores default message format settings in the `connect-distributed.properties` file. Its key and value converters, such as `org.apache.kafka.connect.json.JsonConverter` and `org.apache.kafka.connect.storage.StringConverter`, are part of the Kafka Connect ecosystem. In the Collector for Kafka, set `receivers.kafka.logs.encoding` to `json` or `text` to match the Splunk Connect for Kafka configuration.

#### Use the REST API to read connector settings
Use the following REST API commands to read the Splunk Connect for Kafka configuration:

| Action | `curl` command | Description |
|--------------------------------|----------------------------------------------------------------|----------------------------------------------|
| List active connectors | `curl http://localhost:8083/connectors`                        | Lists all active connectors |
| Get Splunk Connect for Kafka connector info | `curl http://localhost:8083/connectors/<CONNECTOR_NAME>`       | Retrieves information about the specified Splunk Connect for Kafka connector |
| Get Splunk Connect for Kafka connector config | `curl http://localhost:8083/connectors/<CONNECTOR_NAME>/config` | Retrieves configuration details of the specified Splunk Connect for Kafka connector |
| Get Splunk Connect for Kafka connector task info | `curl http://localhost:8083/connectors/<CONNECTOR_NAME>/tasks`  | Retrieves task information for the specified Splunk Connect for Kafka connector |
