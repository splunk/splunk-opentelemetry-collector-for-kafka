## Quickstart Guide

This guide shows you how to install SOC4Kafka with Ansible.

!!! note
    This guide covers the setup of a simple, basic configuration to get you started quickly. Once the `values.yaml` file is generated, it can be further adjusted and customized to suit your specific needs. For more advanced configuration options, please refer to the documentation.

### Prerequisites

!!! note
    This guide is applicable for Linux and macOS systems. Windows is not supported.


Before you begin, confirm that you have the following prerequisites:

- A running instance of Splunk
    - with a valid HTTP Event Collector (HEC) token from your Splunk instance
    - index created for Kafka logs (e.g., `kafka_otel`)
- A running instance of Kafka
- Network connectivity among your Kafka instance, Splunk, and the VM where you will install SOC4Kafka
- Ansible installed on the VM where you will install SOC4Kafka

### Quickstart steps
1. Download Ansible script: [install_soc4kafka_collector.yaml](https://github.com/splunk/splunk-opentelemetry-collector-for-kafka/blob/main/quickstart/install_soc4kafka_collector.yaml)

```bash
wget https://raw.githubusercontent.com/splunk/splunk-opentelemetry-collector-for-kafka/refs/heads/main/quickstart/install_soc4kafka_collector.yaml
```
2. Fill in the variables in the Ansible script. More information about the variables can be found in the [Variables description](#variables-description) section below.

3. Run the Ansible playbook:
```bash
ansible-playbook install_soc4kafka_collector.yaml
```

4. After the playbook runs, look for a command to start the collector, such as:
```bash
./<otelcol_binary_file_name> --config values.yaml
```

5. Run the above command to start the collector.

Once the collector is running, you should start seeing logs in your Splunk instance. 
Continue with the advanced configuration topics to learn about more SOC4Kafka features.

### Variables description

| Variable              | Type    | Description                                                                                     | Allowed Values       | Default                  | Example                                                |
|-----------------------|---------|-------------------------------------------------------------------------------------------------|----------------------|--------------------------|--------------------------------------------------------|
| **Upgrade_SOC4Kafka** | Boolean | Set to `true` to upgrade the SOC4Kafka binary if it already exists.                             | `true`, `false`     | `true`                   | -                                                      |
| **Operating_System**  | String  | Specifies the operating system.                                                                | `linux`, `darwin`   | `"linux"`                | -                                                      |
| **Architecture**      | String  | Specifies the system architecture.                                                             | `amd64`, `arm64`    | `"amd64"`                | -                                                      |
| **Brokers**           | String  | Comma-separated list of Kafka brokers in the format `broker:port`.                             | -                    | -                        | `"broker1:port1"` or`"broker1:port1,broker2:port2"`    |
| **Topic**             | String  | The Kafka topic to consume messages from.                                                     | -                    | -                        | `"example-topic"`                                      |
| **Encoding**          | String  | Specifies the message encoding format.                                                        | `text`, `json`      | `"text"`                 | -                                                      |
| **Insecure_Skip_Verify** | Boolean | Set to `true` to skip TLS certificate verification. Not recommended for production.           | `true`, `false`     | `false`                  | -                                                      |
| **Splunk_HEC_Token**  | String  | The HTTP Event Collector (HEC) token for Splunk.                                              | -                    | -                        | `"your-splunk-hec-token"`                              |
| **Splunk_HEC_Endpoint** | String | The Splunk HEC endpoint URL.                                                                  | -                    | -                        | `"https://splunk-hec-endpoint:8088/services/collector"` |
| **Source**            | String  | The source field value to assign to events sent to Splunk.                                    | -                    | -                        | `"example-source"`                                     |
| **Sourcetype**        | String  | The sourcetype field value to assign to events sent to Splunk.                                | -                    | -                        | `"example-sourcetype"`                                 |
| **Splunk_Index**      | String  | The Splunk index where events will be stored.                                                | -                    | -                        | `"example-index"`                                      |