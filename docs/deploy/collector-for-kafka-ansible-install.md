# Install the Collector for Kafka with Ansible

Install the Collector for Kafka on a Linux or macOS host by using the Ansible playbook.

!!! note
    The playbook creates a basic configuration in `values.yaml`. Edit the file to meet your requirements. For more configuration options, see [Configure the Collector for Kafka](../configure/collector-for-kafka-configure.md).

## Prerequisites

!!! note
    The playbook supports Linux and macOS. It does not support Windows.


Before you begin, make sure that you have the following prerequisites:

- A running instance of Splunk
    - A valid HTTP Event Collector (HEC) token
    - An index for Kafka logs, such as `kafka_otel`
- A running instance of Kafka
- Network connectivity among Kafka, the Splunk platform, and the host where you will install the collector
- Ansible installed on the host where you will install the collector

## Install the Collector for Kafka

1. Download the Ansible script [install_soc4kafka_collector.yaml](https://github.com/splunk/splunk-opentelemetry-collector-for-kafka/blob/main/quickstart/install_soc4kafka_collector.yaml):

```bash
wget https://raw.githubusercontent.com/splunk/splunk-opentelemetry-collector-for-kafka/refs/heads/main/quickstart/install_soc4kafka_collector.yaml
```
2. Set the variables in the Ansible script. For descriptions, see [Ansible variables](#ansible-variables).

3. Run the Ansible playbook:
```bash
ansible-playbook install_soc4kafka_collector.yaml
```

4. After the playbook finishes, run the command it provides to start the collector. For example:

```bash
./<otelcol_binary_file_name> --config values.yaml
```


When the collector is running, its logs appear in the Splunk platform. For more configuration options, see [Configure the Collector for Kafka](../configure/collector-for-kafka-configure.md).


## Ansible variables

| Variable              | Type    | Description                                                                                     | Allowed values       | Default                  | Example                                                |
|-----------------------|---------|-------------------------------------------------------------------------------------------------|----------------------|--------------------------|--------------------------------------------------------|
| `Upgrade_SOC4Kafka` | Boolean | Set to `true` to upgrade the collector binary if it already exists.                             | `true`, `false`     | `true`                   | -                                                      |
| `Operating_System`  | String  | Operating system.                                                                             | `linux`, `darwin`   | `"linux"`                | -                                                      |
| `Architecture`      | String  | System architecture.                                                                          | `amd64`, `arm64`    | `"amd64"`                | -                                                      |
| `Brokers`           | String  | Comma-separated Kafka broker addresses in the format `broker:port`.                          | -                    | -                        | `"broker1:port1"` or `"broker1:port1,broker2:port2"` |
| **Topic**             | String  | Kafka topic from which to collect messages.                                                    | -                    | -                        | `"example-topic"`                                      |
| `Encoding`          | String  | Kafka message encoding format.                                                                | `text`, `json`      | `"text"`                 | -                                                      |
| `Insecure_Skip_Verify` | Boolean | Set to `true` to skip TLS certificate verification. Not recommended for production.           | `true`, `false`     | `false`                  | -                                                      |
| `Splunk_HEC_Token`  | String  | HTTP Event Collector (HEC) token for authentication.                                         | -                    | -                        | `"your-splunk-hec-token"`                              |
| `Splunk_HEC_Endpoint` | String | Splunk HEC endpoint URL.                                                                      | -                    | -                        | `"https://splunk-hec-endpoint:8088/services/collector"` |
| `Source`            | String  | Source value to assign to events sent to Splunk.                                               | -                    | -                        | `"example-source"`                                     |
| `Sourcetype`        | String  | Sourcetype value to assign to events sent to Splunk.                                           | -                    | -                        | `"example-sourcetype"`                                 |
| `Splunk_Index`      | String  | Splunk index where events are stored.                                                         | -                    | -                        | `"example-index"`                                      |
