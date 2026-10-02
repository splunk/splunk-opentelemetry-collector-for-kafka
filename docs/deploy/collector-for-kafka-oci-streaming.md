# Deploy the Collector for Kafka for OCI Streaming

Install and configure the Splunk Distribution of OpenTelemetry Collector for Kafka on an OCI Ubuntu VM to forward records from an OCI Streaming stream to Splunk by using HTTP Event Collector (HEC).

Choose one of these deployment methods:

- **Systemd**: Run the collector binary directly on the VM. This method does not require a container runtime.
- **Kubernetes**: Run the collector in a Kubernetes pod by using the Helm chart.

---

## Gather the required values

Collect the following values before you configure the VM. Replace each `<PLACEHOLDER>` with your value.

### From the OCI Console

| What you need | Where to find it | Placeholder |
|---|---|---|
| Kafka bootstrap endpoint | Streaming → Stream Pools → select pool → **Kafka Connection Settings** → Bootstrap Servers | `<KAFKA_BOOTSTRAP>` |
| Kafka SASL username | Same page → **Username** (fully formed, ready to copy) | `<SASL_USERNAME>` |
| OCI auth token (SASL password) | Profile → User Settings → **Auth Tokens** → Generate Token. Copy the token when it appears; OCI displays it only once. | `<OCI_AUTH_TOKEN>` |
| Stream / topic name | Streaming → **Streams** | `<TOPIC>` |

!!! info 
    **Preserve special characters:** The OCI auth token can contain characters such as `&`, `|`, `>`,
    `` ` ``, or `!`. Keep the single quotes shown around `KAFKA_SASL_PASS` and `--from-literal=...`
    when you enter the token. Without the quotes, the shell can interpret the characters and change
    the value before passing it to the collector.

### From the Splunk platform

| What you need | Where to find it | Placeholder |
|---|---|---|
| HEC endpoint URL | Settings → Data inputs → HTTP Event Collector → host + port `8088`, path `/services/collector` | `<SPLUNK_HEC_ENDPOINT>` |
| HEC token | Settings → Data inputs → HTTP Event Collector → token's **Token Value** | `<SPLUNK_HEC_TOKEN>` |
| Target index | Settings → **Indexes** | `<SPLUNK_INDEX>` |

!!! info 
    **HEC prerequisites:** Before you install the collector, make sure that HEC is globally enabled
    (**Global Settings → Enabled**) and that **Indexer Acknowledgement is OFF**. The Collector for Kafka
    does not implement HEC acknowledgments. If you enable them, the connection stalls.

### Choose a consumer group name

Set `<CONSUMER_GROUP>` to a short, unique string, such as `soc4kafka-v1`. Kafka uses this name to identify the collector instance. Use a new name for each installation. Reusing the group ID from a failed installation can cause the collector to loop during startup.

---

## Choose an installation method

Follow [the systemd procedure](collector-for-kafka-manual-install.md#oci-streaming-on-ubuntu-with-systemd) on an Ubuntu VM, or [the MicroK8s and Helm procedure](collector-for-kafka-kubernetes-install-helm.md#deploy-to-oci-streaming-with-microk8s-and-helm) on a single-node Kubernetes VM.
