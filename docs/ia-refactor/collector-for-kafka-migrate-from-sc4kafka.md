# Migrate from Splunk Connect for Kafka

Naming: 

- SC4Kafka - the [Splunk Connect for Kafka](https://github.com/splunk/kafka-connect-splunk)
- SOC4Kafka - the Splunk OTel Collector for Kafka (the current project)

The main differences between SC4Kafka and SOC4Kafka include:

| **Field**                  | **SC4Kafka**                                                                          | **SOC4Kafka**                                                                                                                                                                                                    |
|----------------------------|---------------------------------------------------------------------------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| **Type**                   | Connector based on **Kafka Connect**, installed as an add-on for Kafka.               | Standalone product that works independently of Kafka.                                                                                                                                                            |
| **Message Retrieval**      | Retrieves events directly from Kafka.                                                 | Consumes messages from Kafka via the Kafka OpenTelemetry Receiver (native Kafka consumer protocol).                                                                                                             |
| **Processing**             | Sends events directly to Splunk using the **Splunk HTTP Event Collector (HEC) exporter**.                    | Processes messages internally and supports customization using **transform processors** before sending them to the **Splunk HEC exporter**.                                                                      |
| **Integration with Kafka** | Tightly integrated as part of the Kafka ecosystem.                                    | Can run independently and be deployed on an external server, separate from the Kafka cluster.                                                                                                                    |
| **Scaling**                | Scaling is managed using the `tasks.max` setting and supports multiple HEC endpoints. | Scaling is achieved by deploying multiple SOC4Kafka instances with the same `group_id`. Multiple HEC endpoints are not supported, but you can create multiple Splunk HEC exporters and add them to the pipeline. |

!!! note
    - The timestamp behavior differs between the two solutions. SOC4Kafka assigns a timestamp to the event based on when it is indexed, whereas Splunk Connect for Kafka uses the timestamp from when the event was originally produced. 
    - Additionally messages from SOC4Kafka appear in Splunk first, as it forwards events to Splunk immediately. In contrast, Splunk Connect for Kafka processes and forwards events in batches, typically every configured number of seconds.

## Migration process

### Prepare and migrate the configuration

Before migrating please get familiar with the [migration strategies](collector-for-kafka-migrate-strategy.md).

--- 
### Important notes:
- **Migration from the old SC4Kafka connector to SOC4Kafka collector is a manual process.** There is no automated tool available for this migration.
- Begin with a simple configuration, then gradually add more settings. This approach helps in isolating and troubleshooting potential issues during the migration.

---

Migrating from SC4Kafka to SOC4Kafka involves several steps to ensure a smooth transition. Below are the key steps to follow during the migration process:

1. **Review Current SC4Kafka Configuration**: Start by thoroughly reviewing your existing SC4Kafka configuration. Document all the settings, including topics, indexes, sourcetypes, and any custom configurations you have in place. 
    In order to read the existing SC4Kafka configuration you can use REST API calls as described in the [Reading the existing Splunk Connect for Kafka configuration](#reading-sc4kafka-existing-configuration).
2. **Map Configuration Parameters**: Use the [configuration mapping table](collector-for-kafka-migrate-config.md) to identify equivalent settings in SOC4Kafka. This will help you understand how to translate your SC4Kafka configuration into SOC4Kafka format.
3. **Create SOC4Kafka Configuration**: Based on the mapped parameters, create a new configuration file for SOC4Kafka. Make sure to include all relevant settings, such as Kafka brokers, topics, Splunk HEC endpoint, and token.
4. **Set Up SOC4Kafka**: [Install SOC4Kafka](collector-for-kafka-kubernetes-install-helm.md) on your desired server. Ensure that you have the necessary permissions and access to both Kafka and Splunk.
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
