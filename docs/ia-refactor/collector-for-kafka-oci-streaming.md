# Deploy the Collector for Kafka for OCI Streaming

This guide walks you through installing and configuring the **Splunk OpenTelemetry Collector for
Kafka (SOC4Kafka)** on an OCI Ubuntu VM so that records published to an **OCI Streaming** stream are
forwarded to **Splunk** via HTTP Event Collector (HEC).

It covers two deployment forms - pick the one that fits your environment:

- **Option A - Bare metal / systemd**: the collector binary runs directly on the VM. No container
  runtime required. Good for a simple single-host setup.
- **Option B - Kubernetes**: the collector runs as a Kubernetes pod via the official Helm chart. Good
  if you want pod-level isolation and rolling updates.

---

## Before you start - values to have ready

Collect the following before touching the VM. Everything in this document is a
`<PLACEHOLDER>` - substitute your real values as you go.

### From OCI console

| What you need | Where to find it | Placeholder |
|---|---|---|
| Kafka bootstrap endpoint | Streaming → Stream Pools → select pool → **Kafka Connection Settings** → Bootstrap Servers | `<KAFKA_BOOTSTRAP>` |
| Kafka SASL username | Same page → **Username** (fully formed, ready to copy) | `<SASL_USERNAME>` |
| OCI auth token (SASL password) | Profile → User Settings → **Auth Tokens** → Generate Token - copy immediately, shown once | `<OCI_AUTH_TOKEN>` |
| Stream / topic name | Streaming → **Streams** | `<TOPIC>` |

!!! info 
    **Handling special characters:** the OCI auth token may contain characters like `&`, `|`, `>`, 
    `` ` ``, or `!` - no need to regenerate the token if it does. Just make sure to keep the single
    quotes shown around `KAFKA_SASL_PASS` and `--from-literal=...` below when you set it: without
    them, the **shell** interprets those characters itself and can silently truncate or empty out
    the value before it ever reaches the collector.

### From Splunk

| What you need | Where to find it | Placeholder |
|---|---|---|
| HEC endpoint URL | Settings → Data inputs → HTTP Event Collector → host + port `8088`, path `/services/collector` | `<SPLUNK_HEC_ENDPOINT>` |
| HEC token | Settings → Data inputs → HTTP Event Collector → token's **Token Value** | `<SPLUNK_HEC_TOKEN>` |
| Target index | Settings → **Indexes** | `<SPLUNK_INDEX>` |

!!! info 
    **HEC prerequisites:** before installing the collector, make sure HEC is globally enabled
    (**Global Settings → Enabled**) and that **Indexer Acknowledgement is OFF** - SOC4Kafka does not
    implement HEC ACK and the connection will stall if it is on.

### Choose a consumer group name

Pick a short, unique string for `<CONSUMER_GROUP>` (e.g. `soc4kafka-v1`). This name identifies your
collector instance to the Kafka broker. **Use a fresh name** - reusing a group ID from a previous
failed install can cause the collector to loop indefinitely on startup.

---

## Choose an installation method

The source guide provides an OCI Streaming setup for two environments. Follow [the systemd procedure](collector-for-kafka-manual-install.md#oci-streaming-on-ubuntu-with-systemd) on an Ubuntu VM, or [the MicroK8s and Helm procedure](collector-for-kafka-kubernetes-install-helm.md#oci-streaming-on-kubernetes-with-microk8s-and-helm) on a single-node Kubernetes VM.

