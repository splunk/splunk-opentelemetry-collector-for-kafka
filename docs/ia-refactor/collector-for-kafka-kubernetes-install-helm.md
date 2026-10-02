# Install the Collector for Kafka with Helm

## Install the Helm chart

1. Create a `values.yaml` file with your configuration:

```yaml
kafkaReceivers:
  - name: main
    brokers:
      - "kafka-broker:9092"
    logs:
      topics:
        - "application-logs"
      encoding: text
    group_id: "soc4kafka-main"

splunkExporters:
  - name: primary
    endpoint: "https://splunk-hec:8088/services/collector"
    token: "your-splunk-hec-token"
    source: "soc4kafka"
    sourcetype: "otel:logs"
    index: "main"
    sending_queue:
      enabled: true
      num_consumers: 10
      queue_size: 10000
      block_on_overflow: true
      sizer: items
      batch:
        min_size: 1000

pipelines:
  - name: main-logs
    type: logs
    receivers:
      - main
    exporters:
      - primary
    # processors optional; defaults to ["resourcedetection"] (defaults.pipelineProcessors in values.yaml)
```

2. Add the Helm repository:

```bash
helm repo add splunk-opentelemetry-collector-for-kafka https://splunk.github.io/splunk-opentelemetry-collector-for-kafka
```

3. Install the chart:

```bash
helm upgrade --install soc4kafka splunk-opentelemetry-collector-for-kafka/splunk-opentelemetry-collector-for-kafka -f values.yaml
```

!!! note

    For information about managing secrets (auto-created or existing Kubernetes secrets), see [Secret management](collector-for-kafka-configure-secrets.md).

## Deploy to OCI Streaming with MicroK8s and Helm

Run all commands from an SSH session on the OCI VM.

!!! info
    **Kubernetes distribution:** This procedure uses **MicroK8s** as an example of a single-node Kubernetes setup. The Collector for Kafka Helm chart runs on any conformant Kubernetes cluster, including EKS, GKE, AKS, K3s, and kubeadm. If you use another distribution, replace the `microk8s kubectl` and `microk8s helm3` commands with the commands for your cluster. The DNS configuration in step B.1, firewall configuration in step B.2, and OCI-specific CIDR ranges apply to MicroK8s on an OCI Ubuntu VM. They differ on other distributions and cloud providers.

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
