# Splunk OpenTelemetry Collector for Kafka (SOC4Kafka) - Helm Chart

Deploy SOC4Kafka on Kubernetes or **MicroK8s** to collect messages from Kafka topics and forward them to Splunk via HTTP Event Collector (HEC).

## Requirements

- Kubernetes 1.19+ (including MicroK8s)
- Kafka 3.7+ reachable from the cluster
- Splunk with HEC enabled and a valid token

## Quick start (MicroK8s)

```bash
# From the chart directory
cd helm-chart/splunk-opentelemetry-collector-for-kafka

# 1. Create a secret with your Splunk HEC token (do not commit this)
kubectl create secret generic soc4kafka-hec --from-literal=splunk-hec-token=YOUR_HEC_TOKEN

# 2. Install with your Kafka and Splunk settings
helm install soc4kafka . \
  --set existingSecret=soc4kafka-hec \
  --set kafka.brokers="kafka.example.com:9092" \
  --set kafka.topics[0]=my-logs \
  --set splunk.hec.endpoint="https://splunk.example.com:8088/services/collector" \
  --set splunk.hec.index=main

# 3. Check status
kubectl get pods -l app.kubernetes.io/name=splunk-opentelemetry-collector-for-kafka
kubectl logs -l app.kubernetes.io/name=splunk-opentelemetry-collector-for-kafka -f
```

## MicroK8s notes

- **Enable Helm** (if needed): `microk8s enable helm3`
- **Kafka access**: Use a Kafka service address that pods can reach (e.g. in-cluster Kafka, or external hostname/IP). For Kafka in the same MicroK8s cluster, use the Kubernetes service DNS name (e.g. `kafka-broker.kafka.svc.cluster.local:9092`).
- **Image registry**: MicroK8s can pull from Quay; ensure no private registry auth is required or configure `imagePullSecrets` in values.
- **Storage**: For `persistence.enabled: true`, MicroK8s uses the default storage class; no extra setup required for basic PVCs.

## Configuration

| Parameter | Description | Default |
|-----------|-------------|---------|
| `image.repository` | Collector image | `quay.io/signalfx/splunk-otel-collector` |
| `image.tag` | Image tag | `0.143.0` |
| `kafka.brokers` | Kafka broker list (comma-separated) | `kafka-broker:9092` |
| `kafka.topics` | List of topics to consume | `[example-topic]` |
| `kafka.encoding` | Message encoding | `text` |
| `splunk.hec.token` | HEC token (prefer `existingSecret`) | `""` |
| `splunk.hec.endpoint` | HEC URL | (set on install) |
| `splunk.hec.source` | Event source | `soc4kafka` |
| `splunk.hec.sourcetype` | Sourcetype | `kafka-otel` |
| `splunk.hec.index` | Splunk index | `main` |
| `existingSecret` | Secret name containing `splunk-hec-token` | `""` |
| `replicaCount` | Number of collector replicas | `1` |
| `autoscaling.enabled` | Enable HPA | `false` |

**Security:** Do not store the Splunk HEC token in `values.yaml` or in version control. Use `existingSecret` or `--set splunk.hec.token=...` at install time.

## Examples

**Multiple topics:**

```yaml
# values.yaml or --set
kafka:
  brokers: "broker1:9092,broker2:9092"
  topics:
    - topic-a
    - topic-b
  encoding: "text"
```

**Use existing secret:**

```bash
helm install soc4kafka . \
  --set existingSecret=my-splunk-hec-secret \
  --set kafka.brokers="kafka:9092" \
  --set splunk.hec.endpoint="https://splunk:8088/services/collector"
```

**Custom values file (e.g. `microk8s-values.yaml`):**

```yaml
kafka:
  brokers: "kafka.kafka.svc.cluster.local:9092"
  topics:
    - app-logs
  encoding: "text"

splunk:
  hec:
    endpoint: "https://splunk.company.com:8088/services/collector"
    index: "kafka_otel"
    source: "microk8s-kafka"

existingSecret: "soc4kafka-hec"
```

Then: `helm install soc4kafka . -f microk8s-values.yaml`

## Uninstall

```bash
helm uninstall soc4kafka
kubectl delete secret soc4kafka-hec  # if you created it
```

## Links

- [SOC4Kafka project](https://github.com/splunk/splunk-opentelemetry-collector-for-kafka)
- [Quickstart guide](https://github.com/splunk/splunk-opentelemetry-collector-for-kafka/blob/main/docs/quickstart_guide.md)
