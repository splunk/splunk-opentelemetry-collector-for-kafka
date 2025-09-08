## Quickstart Guide

Welcome to the Quickstart Guide! This guide will help you get up and running with SOC4Kafka in just a few simple steps.

### Prerequisites
Before you begin, ensure you have the following prerequisites in place:
- A running instance of Splunk
  - with a valid HTTP Event Collector (HEC) token from your Splunk instance
  - index created for Kafka logs (e.g., `kafka_otel`)
- A running instance of Kafka
- Network connectivity between your Kafka instance and Splunk and the VM where SOC4Kafka will be installed
- Ansible installed on VM where SOC4Kafka will be installed

### Quickstart Steps
1. Download Ansible script: [install_soc4kafka_collector.yaml](../quickstart/install_soc4kafka_collector.yaml)
``wget https://raw.githubusercontent.com/signalfx/splunk-otel-collector/main/quickstart/install_soc4kafka_collector.yaml``

2. Fill in the variables in the Ansible script:
More information about the variables can be found in the [Variables Description](#variables-description) section below.

3. Run the Ansible playbook:
``ansible-playbook install_soc4kafka_collector.yaml``

4. Verify Ansible script ran successfully, you should see the command which needs to be run to start the collector, something like:
``./<otelcol_binary_file_name> --config values.yaml"``

5. Run the above command to start the collector.

Once the collector is running, you should start seeing logs in your Splunk instance. 
Now you are ready to explore more advanced configurations and features of SOC4Kafka!

### Variables Description

| Variable              | Type    | Description                                                                                     | Allowed Values                  | Default                  | Example                                                |
|-----------------------|---------|-------------------------------------------------------------------------------------------------|---------------------------------|--------------------------|--------------------------------------------------------|
| **Upgrade_SOC4Kafka** | Boolean | Set to `true` to upgrade the SOC4Kafka binary if it already exists.                             | `true`, `false`                | `true`                   | -                                                      |
| **Operating_System**  | String  | Specifies the operating system.                                                                | `linux`, `windows`, `darwin`   | `"linux"`                | -                                                      |
| **Architecture**      | String  | Specifies the system architecture.                                                             | `amd64`, `arm64`               | `"amd64"`                | -                                                      |
| **Brokers**           | String  | Comma-separated list of Kafka brokers in the format `broker:port`.                             | -                               | -                        | `"broker1:port1"` or`"broker1:port1,broker2:port2"`    |
| **Topic**             | String  | The Kafka topic to consume messages from.                                                     | -                               | -                        | `"example-topic"`                                      |
| **Encoding**          | String  | Specifies the message encoding format.                                                        | `text`, `json`                 | `"text"`                 | -                                                      |
| **Insecure_Skip_Verify** | Boolean | Set to `true` to skip TLS certificate verification. Not recommended for production.           | `true`, `false`                | `false`                  | -                                                      |
| **Splunk_HEC_Token**  | String  | The HTTP Event Collector (HEC) token for Splunk.                                              | -                               | -                        | `"your-splunk-hec-token"`                              |
| **Splunk_HEC_Endpoint** | String | The Splunk HEC endpoint URL.                                                                  | -                               | -                        | `"https://splunk-hec-endpoint:8088/services/collector"` |
| **Source**            | String  | The source field value to assign to events sent to Splunk.                                    | -                               | -                        | `"example-source"`                                     |
| **Sourcetype**        | String  | The sourcetype field value to assign to events sent to Splunk.                                | -                               | -                        | `"example-sourcetype"`                                 |
| **Splunk_Index**      | String  | The Splunk index where events will be stored.                                                | -                               | -                        | `"example-index"`                                      |