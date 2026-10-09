# Deploy the Splunk Distribution of OpenTelemetry Collector for Kafka for OCI Streaming

Install and configure the Splunk Distribution of OpenTelemetry Collector for Kafka on an OCI Ubuntu VM to forward records from an OCI Streaming stream to Splunk by using HTTP Event Collector (HEC).

Choose the deployment method that fits your environment:

- **Systemd**: Run the collector binary directly on the VM. Choose this method for a simple, single-host deployment when you do not want to install a container runtime or Kubernetes.
- **Kubernetes**: Run the collector in a Kubernetes pod by using the Helm chart. Choose this method when you want pod isolation and Kubernetes-managed rolling updates.

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
    (**Global Settings → Enabled**) and that **Indexer Acknowledgement is OFF**. The Splunk Distribution of OpenTelemetry Collector for Kafka
    does not implement HEC acknowledgments. If you enable them, the connection stalls.

### Choose a consumer group name

Set `<CONSUMER_GROUP>` to a short, unique string, such as `soc4kafka-v1`. Kafka uses this name to identify the collector instance. Use a new name for each installation. Reusing the group ID from a failed installation can cause the collector to loop during startup.

---

## Choose an installation method

Follow [the systemd procedure](#oci-streaming-on-ubuntu-with-systemd) on an Ubuntu VM, or [the MicroK8s and Helm procedure](#deploy-to-oci-streaming-with-microk8s-and-helm) on a single-node Kubernetes VM.

## OCI Streaming on Ubuntu with systemd

All commands run on the OCI VM over SSH.

### A.1 Install dependencies

These packages support connectivity tests and test message production. The Splunk Distribution of OpenTelemetry Collector for Kafka does not require them to run.

```bash
sudo apt-get update
sudo apt-get install -y kafkacat curl netcat-openbsd
```

### A.2 Download the collector binary

Splunk Distribution of OpenTelemetry Collector for Kafka releases are published on GitHub. Download the binary for your target version, make it executable, and place it in a working directory.

```bash
mkdir -p ~/soc4kafka && cd ~/soc4kafka
wget https://github.com/signalfx/splunk-otel-collector/releases/download/v0.158.0/otelcol_linux_amd64
chmod +x otelcol_linux_amd64
```

!!! note
    Check the [releases page](https://github.com/splunk/splunk-opentelemetry-collector-for-kafka/releases)
    for newer versions and replace `v0.158.0` with the version you want to install.

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

!!! note
    Wrap `KAFKA_SASL_PASS` in **single quotes** so the shell does not expand special characters in the
    token value.

### A.4 Create the collector configuration

Create `~/soc4kafka/config.yaml` with the content below. Substitute `<CONSUMER_GROUP>` and `<TOPIC>`
directly in the file. These values are not secrets and do not need to be in the environment file.

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

Check that the VM can reach and authenticate with the Kafka broker. This end-to-end test confirms that DNS, routing, TLS, and SASL work:

```bash
set -a; source ~/soc4kafka/collector.env; set +a
kafkacat -L \
  -b "$KAFKA_BOOTSTRAP" \
  -X security.protocol=SASL_SSL \
  -X sasl.mechanisms=PLAIN \
  -X sasl.username="$KAFKA_SASL_USER" \
  -X sasl.password="$KAFKA_SASL_PASS" | head -20
```

The output should list `<TOPIC>`. If the command times out or returns an authentication error, resolve the issue before continuing. The collector will encounter the same failure.

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

### A.7 Install the collector as a systemd service

After the collector starts without errors, install it as a managed service. The service restarts automatically, and journald captures its logs.

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

### A.8 Send a test message and confirm it appears in Splunk

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

In Splunk, search for the event:

```
index=<SPLUNK_INDEX> sourcetype=oci:streaming:text
```

You can also monitor collector throughput from the VM by running:

```bash
curl -s http://127.0.0.1:8888/metrics | grep -E 'otelcol_(receiver_accepted|exporter_sent)'
```

The `receiver_accepted_log_records_total` value increases when you produce a message. The `exporter_sent_log_records_total` value increases after the collector sends the batch.

## Deploy to OCI Streaming with MicroK8s and Helm

Run all commands from an SSH session on the OCI VM.

!!! info
    **Kubernetes distribution:** This procedure uses **MicroK8s** as an example of a single-node Kubernetes setup. The Splunk Distribution of OpenTelemetry Collector for Kafka Helm chart runs on any conformant Kubernetes cluster, including EKS, GKE, AKS, K3s, and kubeadm. If you use another distribution, replace the `microk8s kubectl` and `microk8s helm3` commands with the commands for your cluster. The DNS configuration in step B.1, firewall configuration in step B.2, and OCI-specific CIDR ranges apply to MicroK8s on an OCI Ubuntu VM. They differ on other distributions and cloud providers.

    MicroK8s includes its own `helm3` and `kubectl` commands. This procedure uses `microk8s helm3` and `microk8s kubectl` instead of system-level commands.

### B.1 Install MicroK8s

```bash
sudo snap install microk8s --classic --channel=1.33/stable
sudo usermod -a -G microk8s "$USER"
sudo chown -f -R "$USER" ~/.kube
newgrp microk8s
```

Enable the add-ons required by the chart:

```bash
microk8s enable hostpath-storage
microk8s enable rbac
microk8s enable metrics-server
```

Configure DNS to use the **OCI VCN resolver**. It resolves private OCI names, such as your broker's private endpoint, and public names, such as your Splunk HEC host. Use it as the only upstream resolver:

```bash
microk8s enable dns:169.254.169.254
```

!!! warning
    Do not add a public resolver such as `8.8.8.8`. The OCI Streaming broker resolves to a private VCN IP address, which a public resolver cannot resolve. Adding a public resolver can cause intermittent connection failures.

### B.2 Configure the OCI host firewall

The OCI Ubuntu image includes a firewall rule that blocks forwarded traffic. This prevents pods from reaching the Kubernetes API server and causes CoreDNS and Calico to restart repeatedly. Remove the rule:

```bash
sudo iptables -L FORWARD -n --line-numbers | head
sudo iptables -D FORWARD 1    # removes the REJECT rule (usually at position 1)
```

Pods recover within about 60 seconds. Make the change permanent; otherwise, the rule returns after a reboot:

```bash
# Edit the persisted ruleset and remove the REJECT line, then reload:
sudo grep -nE 'REJECT|icmp-host-prohibited' /etc/iptables/rules.v4
# Delete the matching line from the file, then:
sudo netfilter-persistent reload
```

On a test VM, you can turn off the OS firewall. The OCI VCN security list still controls ingress at the cloud layer:

```bash
sudo systemctl disable --now netfilter-persistent
```

!!! warning
    **If Calico continues to restart** after you remove the FORWARD rule, the image also has an INPUT-chain
    REJECT rule that blocks pod traffic to the Kubernetes API server VIP (`10.152.183.1`) and pod CIDR
    (`10.1.0.0/16`). These are standard MicroK8s defaults. Allow them:

    `sudo iptables -I INPUT 4 -s 10.152.183.0/24 -j ACCEPT`

    `sudo iptables -I INPUT 4 -d 10.152.183.0/24 -j ACCEPT`

    `sudo iptables -I INPUT 4 -s 10.1.0.0/16    -j ACCEPT`

    `sudo iptables -I INPUT 4 -d 10.1.0.0/16    -j ACCEPT`

    The `-I INPUT 4` option inserts the rule before the catch-all REJECT. Confirm the position by running
    `sudo iptables -L INPUT -n --line-numbers` first. If you customised MicroK8s CIDRs, replace the
    ranges with your service CIDR (`grep service-cluster-ip-range /var/snap/microk8s/current/args/*`)
    and pod CIDR (`grep cluster-cidr /var/snap/microk8s/current/args/*`).

### B.3 Create the namespace

```bash
microk8s kubectl create namespace soc4kafka
```

### B.4 Create the Kubernetes Secrets

The collector reads credentials from Kubernetes Secrets injected as environment variables. This keeps the credentials out of the Helm values file.

```bash
# Kafka SASL password - the key name "password" is required by the chart
microk8s kubectl -n soc4kafka create secret generic kafka-sasl \
  --from-literal=password='<OCI_AUTH_TOKEN>'

# Splunk HEC token - the key name "splunk-hec-token" is required by the chart
microk8s kubectl -n soc4kafka create secret generic splunk-hec \
  --from-literal=splunk-hec-token='<SPLUNK_HEC_TOKEN>'
```
!!! warning
    Enclose values in **single quotes** to prevent the shell from interpreting special characters.

### B.5 Create `values.yaml`

Create this file on the VM, for example, at `~/soc4kafka_microk8s/values.yaml`, before you install the chart. Replace each `<PLACEHOLDER>` with the corresponding value.

```yaml
replicaCount: 1

kafkaReceivers:
  - name: main
    brokers:
      - <KAFKA_BOOTSTRAP>
    client_id: <CONSUMER_GROUP>       # e.g. soc4kafka-m8k-v1 - must be fresh
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

!!! warning
    Include `-n soc4kafka` in the command. Otherwise, Helm installs the release in the `default` namespace.

### B.7 Verify the deployment

Confirm that all pods are running:

```bash
microk8s kubectl get pods -A
```

View the collector logs and confirm that the startup sequence completes:

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

!!! note
    If `NOT_COORDINATOR` repeats, change `client_id` and `group_id` to a new name in `values.yaml`. Then run the `helm3 upgrade` command from step B.6 again.

### B.8 Send a test message and confirm it appears in Splunk

Produce a message from the VM. If `kafkacat` is not installed, install it by running `sudo apt-get install -y kafkacat`:

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

Alternatively, use the OCI Console and select **Streaming → Streams → your stream → Produce Test Message**.

In Splunk, search for the event:

```
index=<SPLUNK_INDEX> sourcetype="oci:streaming:text" earliest=-5m
```
