# SOC4Kafka Installation Guide — OCI Streaming to Splunk

This guide walks you through installing and configuring the **Splunk OpenTelemetry Collector for
Kafka (SOC4Kafka)** on an OCI Ubuntu VM so that records published to an **OCI Streaming** stream are
forwarded to **Splunk** via HTTP Event Collector (HEC).

It covers two deployment forms — pick the one that fits your environment:

- **Option A — Bare metal / systemd**: the collector binary runs directly on the VM. No container
  runtime required. Good for a simple single-host setup.
- **Option B — Kubernetes**: the collector runs as a Kubernetes pod via the official Helm chart. Good
  if you want pod-level isolation and rolling updates.

---

## Before you start — values to have ready

Collect the following before touching the VM. Everything in this document is a
`<PLACEHOLDER>` — substitute your real values as you go.

### From OCI Console

| What you need | Where to find it | Placeholder |
|---|---|---|
| Kafka bootstrap endpoint | Streaming → Stream Pools → select pool → **Kafka Connection Settings** → Bootstrap Servers | `<KAFKA_BOOTSTRAP>` |
| Kafka SASL username | Same page → **Username** (fully formed, ready to copy) | `<SASL_USERNAME>` |
| OCI auth token (SASL password) | Profile → User Settings → **Auth Tokens** → Generate Token — copy immediately, shown once | `<OCI_AUTH_TOKEN>` |
| Stream / topic name | Streaming → **Streams** | `<TOPIC>` |

> **Auth token constraint:** the token must **not begin with** `#`, `&`, `*`, `!`, `>`, `|`, `@`,
> or `` ` ``. These characters have special meaning in YAML and will silently break the collector
> configuration. If your generated token starts with one of them, delete it and generate a new one.

### From Splunk

| What you need | Where to find it | Placeholder |
|---|---|---|
| HEC endpoint URL | Settings → Data inputs → HTTP Event Collector → host + port `8088`, path `/services/collector` | `<SPLUNK_HEC_ENDPOINT>` |
| HEC token | Settings → Data inputs → HTTP Event Collector → token's **Token Value** | `<SPLUNK_HEC_TOKEN>` |
| Target index | Settings → **Indexes** | `<SPLUNK_INDEX>` |

> **HEC prerequisites:** before installing the collector, make sure HEC is globally enabled
> (**Global Settings → Enabled**) and that **Indexer Acknowledgement is OFF** — SOC4Kafka does not
> implement HEC ACK and the connection will stall if it is on.

### Choose a consumer group name

Pick a short, unique string for `<CONSUMER_GROUP>` (e.g. `soc4kafka-v1`). This name identifies your
collector instance to the Kafka broker. **Use a fresh name** — reusing a group ID from a previous
failed install can cause the collector to loop indefinitely on startup.

---

## Option A — Bare metal / systemd

All commands run on the OCI VM over SSH.

### A.1 Install dependencies

These packages are used for connectivity testing and producing test messages. They are not required
for the collector itself to run.

```bash
sudo apt-get update
sudo apt-get install -y kafkacat curl netcat-openbsd
```

### A.2 Download the collector binary

SOC4Kafka releases are published on GitHub. Download the binary for your target version, make it
executable, and place it in a working directory.

```bash
mkdir -p ~/soc4kafka && cd ~/soc4kafka
wget https://github.com/signalfx/splunk-otel-collector/releases/download/v0.151.0/otelcol_linux_amd64
chmod +x otelcol_linux_amd64
```

> Check the [releases page](https://github.com/splunk/splunk-opentelemetry-collector-for-kafka/releases)
> for newer versions and substitute `v0.151.0` accordingly.

### A.3 Create the secrets file

Create `~/soc4kafka/collector.env` and restrict its permissions. This file holds all secrets so they
never appear in the config file or in process arguments.

```bash
cat > ~/soc4kafka/collector.env <<'EOF'
KAFKA_BOOTSTRAP=<KAFKA_BOOTSTRAP>
KAFKA_SASL_USER=<SASL_USERNAME>
KAFKA_SASL_PASS='<OCI_AUTH_TOKEN>'
SPLUNK_HEC_URL=<SPLUNK_HEC_ENDPOINT>
SPLUNK_HEC_TOKEN=<SPLUNK_HEC_TOKEN>
SPLUNK_INDEX=<SPLUNK_INDEX>
EOF
chmod 600 ~/soc4kafka/collector.env
```

> Wrap `KAFKA_SASL_PASS` in **single quotes** so the shell does not expand special characters in the
> token value.

### A.4 Create the collector config

Create `~/soc4kafka/config.yaml` with the content below. Substitute `<CONSUMER_GROUP>` and `<TOPIC>`
directly in the file — these are not secrets and do not need to be in the env file.

```yaml
receivers:
  kafka:
    brokers:
      - ${env:KAFKA_BOOTSTRAP}
    group_id: <CONSUMER_GROUP>
    client_id: <CONSUMER_GROUP>
    group_rebalance_strategy: range
    initial_offset: earliest
    tls:
      insecure_skip_verify: false
    auth:
      sasl:
        username: ${env:KAFKA_SASL_USER}
        password: ${env:KAFKA_SASL_PASS}
        mechanism: PLAIN
    logs:
      topics:
        - <TOPIC>
      encoding: text

processors:
  resourcedetection:
    detectors: [system]
    system:
      hostname_sources: ["os"]

exporters:
  splunk_hec:
    token: ${env:SPLUNK_HEC_TOKEN}
    endpoint: ${env:SPLUNK_HEC_URL}
    source: oci-streaming
    sourcetype: oci:streaming:text
    index: ${env:SPLUNK_INDEX}
    tls:
      insecure_skip_verify: true
    splunk_app_name: soc4kafka

service:
  pipelines:
    logs:
      receivers: [kafka]
      processors: [resourcedetection]
      exporters: [splunk_hec]
```

### A.5 Verify connectivity before starting the collector

Check that the VM can reach Splunk HEC:

```bash
nc -vz <SPLUNK_HEC_HOST> 8088
```

Check that the VM can reach the Kafka broker and authenticate (this is the single best end-to-end
connectivity test — success means DNS, routing, TLS, and SASL all work):

```bash
set -a; source ~/soc4kafka/collector.env; set +a
kafkacat -L \
  -b "$KAFKA_BOOTSTRAP" \
  -X security.protocol=SASL_SSL \
  -X sasl.mechanisms=PLAIN \
  -X sasl.username="$KAFKA_SASL_USER" \
  -X sasl.password="$KAFKA_SASL_PASS" | head -20
```

You should see `<TOPIC>` listed in the output. If it times out or returns an auth error, resolve that
before proceeding — the collector will exhibit the same failure.

### A.6 Start the collector

Run in the foreground first to watch the startup logs:

```bash
cd ~/soc4kafka
set -a; source ./collector.env; set +a
./otelcol_linux_amd64 --config config.yaml
```

A healthy startup looks like:

```
Everything is ready. Begin running and processing data.
...joined, balancing group   group: <CONSUMER_GROUP>
...synced                    assigned: <TOPIC>[0]
...beginning heartbeat loop
```

If you see `NOT_COORDINATOR` repeating, stop the collector, change `group_id` and `client_id` to a
new name in `config.yaml`, and restart.

### A.7 Install as a systemd service

Once the collector starts cleanly, promote it to a managed service so it restarts automatically and
its logs are captured by journald.

Copy files into place:

```bash
sudo mkdir -p /opt/soc4kafka /etc/soc4kafka
sudo cp ~/soc4kafka/otelcol_linux_amd64 /opt/soc4kafka/
sudo cp ~/soc4kafka/config.yaml /opt/soc4kafka/
sudo cp ~/soc4kafka/collector.env /etc/soc4kafka/collector.env
sudo chmod 600 /etc/soc4kafka/collector.env
```

Create the service user and set ownership:

```bash
sudo useradd --system --no-create-home --shell /usr/sbin/nologin otel
sudo chown otel:otel /opt/soc4kafka/otelcol_linux_amd64
sudo chown otel:otel /opt/soc4kafka/config.yaml
sudo chown otel:otel /etc/soc4kafka/collector.env
```

Create the unit file:

```bash
sudo tee /etc/systemd/system/soc4kafka.service > /dev/null <<'EOF'
[Unit]
Description=SOC4Kafka collector (OCI Streaming -> Splunk)
After=network-online.target
Wants=network-online.target

[Service]
User=otel
Group=otel
EnvironmentFile=/etc/soc4kafka/collector.env
ExecStart=/opt/soc4kafka/otelcol_linux_amd64 --config /opt/soc4kafka/config.yaml
Restart=on-failure
RestartSec=5s

[Install]
WantedBy=multi-user.target
EOF
```

Enable and start:

```bash
sudo systemctl daemon-reload
sudo systemctl enable --now soc4kafka
sudo journalctl -u soc4kafka -f
```

### A.8 Send a test message and confirm in Splunk

```bash
set -a; source ~/soc4kafka/collector.env; set +a
printf '{"hello":"splunk","ts":"%s"}\n' "$(date -u +%FT%TZ)" | \
  kafkacat -P \
    -b "$KAFKA_BOOTSTRAP" \
    -t <TOPIC> \
    -X security.protocol=SASL_SSL \
    -X sasl.mechanisms=PLAIN \
    -X sasl.username="$KAFKA_SASL_USER" \
    -X sasl.password="$KAFKA_SASL_PASS"
```

In Splunk search:

```
index=<SPLUNK_INDEX> sourcetype=oci:streaming:text
```

You can also monitor collector throughput from the VM:

```bash
curl -s http://127.0.0.1:8888/metrics | grep -E 'otelcol_(receiver_accepted|exporter_sent)'
```

`receiver_accepted_log_records_total` should increment when you produce; `exporter_sent_log_records_total`
should follow shortly after as the batch flushes.

---

## Option B — Kubernetes

All commands run on the OCI VM over SSH.

> **Kubernetes distribution note:** this guide uses **MicroK8s** as a representative example of a
> single-node Kubernetes setup. The SOC4Kafka Helm chart is distribution-agnostic and will run on
> any conformant Kubernetes cluster (EKS, GKE, AKS, K3s, vanilla kubeadm, etc.). If you are using
> a different distribution, substitute your cluster's `kubectl` and `helm` commands for the
> `microk8s kubectl` and `microk8s helm3` equivalents used below. The DNS configuration (step B.1),
> firewall fix (step B.2), and the OCI-specific CIDRs in the step B.2 callout are specific to
> MicroK8s on an OCI Ubuntu VM and will differ on other distributions or cloud providers.


>  MicroK8s ships its own bundled `helm3` and `kubectl`. The commands below use `microk8s helm3` and
> `microk8s kubectl` — not the system-level tools.

### B.1 Install MicroK8s

```bash
sudo snap install microk8s --classic --channel=1.33/stable
sudo usermod -a -G microk8s "$USER"
sudo chown -f -R "$USER" ~/.kube
newgrp microk8s
```

Enable the addons the chart needs:

```bash
microk8s enable hostpath-storage
microk8s enable rbac
microk8s enable metrics-server
```

Enable DNS pinned to the **OCI VCN resolver**. This resolver handles both private OCI names (your
broker's private endpoint) and public names (your Splunk HEC host) — using it as the single
upstream is important:

```bash
microk8s enable dns:169.254.169.254
```

> Do **not** add a public resolver like `8.8.8.8` alongside it. The OCI Streaming broker resolves
> to a private VCN IP, and a public resolver will return NXDOMAIN for it, causing intermittent
> connection failures.

### B.2 Fix the OCI host firewall

The OCI Ubuntu image ships a firewall rule that blocks forwarded traffic. This prevents pods from
reaching the Kubernetes API server, causing CoreDNS and Calico to crash-loop. Remove the rule:

```bash
sudo iptables -L FORWARD -n --line-numbers | head
sudo iptables -D FORWARD 1    # removes the REJECT rule (usually at position 1)
```

Pods recover within about 60 seconds. **Make the fix permanent** — the rule returns on reboot
otherwise:

```bash
# Edit the persisted ruleset and remove the REJECT line, then reload:
sudo grep -nE 'REJECT|icmp-host-prohibited' /etc/iptables/rules.v4
# Delete the matching line from the file, then:
sudo netfilter-persistent reload
```

On a test VM you can instead disable the OS firewall entirely — the OCI VCN security list still
controls ingress at the cloud layer:

```bash
sudo systemctl disable --now netfilter-persistent
```

> **If Calico still crash-loops** after removing the FORWARD rule, your image also has an INPUT-chain
> REJECT that blocks pod traffic to the Kubernetes API server VIP (`10.152.183.1`) and pod CIDR
> (`10.1.0.0/16`). These are standard MicroK8s defaults. Allow them:
>
> ```bash
> sudo iptables -I INPUT 4 -s 10.152.183.0/24 -j ACCEPT
> sudo iptables -I INPUT 4 -d 10.152.183.0/24 -j ACCEPT
> sudo iptables -I INPUT 4 -s 10.1.0.0/16    -j ACCEPT
> sudo iptables -I INPUT 4 -d 10.1.0.0/16    -j ACCEPT
> ```
>
> The `-I INPUT 4` inserts before the catch-all REJECT. Confirm position with
> `sudo iptables -L INPUT -n --line-numbers` first. If you customised MicroK8s CIDRs, replace the
> ranges with your actual service CIDR (`grep service-cluster-ip-range /var/snap/microk8s/current/args/*`)
> and pod CIDR (`grep cluster-cidr /var/snap/microk8s/current/args/*`).

### B.3 Create the namespace

```bash
microk8s kubectl create namespace soc4kafka
```

### B.4 Create the Kubernetes secrets

The collector reads credentials from Kubernetes Secrets injected as environment variables — they
never appear in the Helm values file.

```bash
# Kafka SASL password — the key name "password" is required by the chart
microk8s kubectl -n soc4kafka create secret generic kafka-sasl \
  --from-literal=password='<OCI_AUTH_TOKEN>'

# Splunk HEC token — the key name "splunk-hec-token" is required by the chart
microk8s kubectl -n soc4kafka create secret generic splunk-hec \
  --from-literal=splunk-hec-token='<SPLUNK_HEC_TOKEN>'
```

> Wrap values in **single quotes** to prevent the shell from interpreting special characters.

### B.5 Create `values.yaml`

Create this file on the VM (e.g. at `~/soc4kafka_microk8s/values.yaml`) before running the Helm
install. Substitute all `<PLACEHOLDERS>` with your real values.

```yaml
replicaCount: 1

kafkaReceivers:
  - name: main
    brokers:
      - <KAFKA_BOOTSTRAP>
    client_id: <CONSUMER_GROUP>       # e.g. soc4kafka-m8k-v1 — must be fresh
    group_id: <CONSUMER_GROUP>
    group_rebalance_strategy: range
    initial_offset: earliest
    logs:
      topics:
        - <TOPIC>
      encoding: text
    auth:
      sasl:
        username: <SASL_USERNAME>
        mechanism: PLAIN
        secret: kafka-sasl            # references the Secret created in step B.4
    tls:
      insecure_skip_verify: false     # OCI broker cert is publicly trusted (DigiCert)

splunkExporters:
  - name: primary
    endpoint: <SPLUNK_HEC_ENDPOINT>
    secret: splunk-hec                # references the Secret created in step B.4
    source: oci-streaming
    sourcetype: oci:streaming:text
    index: <SPLUNK_INDEX>
    splunk_app_name: soc4kafka
    tls:
      insecure_skip_verify: true      # Splunk default self-signed cert has no SAN

pipelines:
  - name: oci-to-splunk
    type: logs
    receivers: [main]
    exporters: [primary]
    processors: [resourcedetection]

extraEnv:
  - name: KAFKA_KAFKA_MAIN_SASL_PASSWORD
    valueFrom:
      secretKeyRef:
        name: kafka-sasl
        key: password

resources:
  limits:
    cpu: 500m
    memory: 512Mi
  requests:
    cpu: 100m
    memory: 256Mi

collectorLogs:
  enabled: false
collectorMetrics:
  enabled: false
```

### B.6 Install the chart

```bash
microk8s helm3 repo add splunk-opentelemetry-collector-for-kafka \
  https://splunk.github.io/splunk-opentelemetry-collector-for-kafka
microk8s helm3 repo update

microk8s helm3 upgrade --install soc4kafka \
  splunk-opentelemetry-collector-for-kafka/splunk-opentelemetry-collector-for-kafka \
  -n soc4kafka \
  -f ~/soc4kafka_microk8s/values.yaml
```

> Always include `-n soc4kafka`. Without it the release lands in the `default` namespace and will be
> difficult to find.

### B.7 Verify the deployment

Check that all pods are running:

```bash
microk8s kubectl get pods -A
```

Tail the collector logs and look for the healthy startup sequence:

```bash
microk8s kubectl -n soc4kafka logs -f \
  deploy/soc4kafka-splunk-opentelemetry-collector-for-kafka
```

Expected output:

```
Everything is ready. Begin running and processing data.
franz   joined, balancing group   group: <CONSUMER_GROUP>
franz   synced                    assigned: <TOPIC>[0]
franz   assigning partitions      ...
```

> If you see `NOT_COORDINATOR` repeating, change `client_id` and `group_id` to a new name in
> `values.yaml` and re-run the `helm3 upgrade` command from step B.6.

### B.8 Send a test message and confirm in Splunk

Produce a message from the VM (install `kafkacat` first if needed: `sudo apt-get install -y kafkacat`):

```bash
echo "hello-from-microk8s-$(date -Is)" | kafkacat -P \
  -b <KAFKA_BOOTSTRAP> \
  -t <TOPIC> \
  -X security.protocol=SASL_SSL \
  -X sasl.mechanisms=PLAIN \
  -X sasl.username='<SASL_USERNAME>' \
  -X sasl.password='<OCI_AUTH_TOKEN>' \
  -X ssl.ca.location=/etc/ssl/certs/ca-certificates.crt
```

Alternatively use the OCI Console: **Streaming → Streams → select stream → Produce Test Message**.

In Splunk search:

```
index=<SPLUNK_INDEX> sourcetype="oci:streaming:text" earliest=-5m
```
