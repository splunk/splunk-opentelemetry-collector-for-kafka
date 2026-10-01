# Install the Collector for Kafka manually

Run the Collector from a downloaded package and a configuration file. For a command-by-command walkthrough that uses OCI Streaming as the source, see [Deploy the Collector for Kafka for OCI Streaming](collector-for-kafka-oci-streaming.md).

### Download the Splunk OTel Collector package

SOC4Kafka uses the Splunk OpenTelemetry Collector package. Choose an installation method that suits your environment.
Download the package for your platform from the [Splunk OpenTelemetry Collector releases](https://github.com/signalfx/splunk-otel-collector/releases). Release numbers start with `v`.

For example, on Linux with an AMD64 architecture, run the following `wget` command:

```commandline
wget https://github.com/signalfx/splunk-otel-collector/releases/download/v0.158.0/otelcol_linux_amd64
```

### Run the Splunk OTel Collector package with a config file

To run SOC4Kafka, use the base package along with a completed configuration template.

```commandline
./<otel_package> --config <config_file>
```

!!! note

    Make sure the file has executable permissions before running the command. On Linux-based systems you can add executable permissions using the following command:

```commandline
chmod a+x <otel_package>
```

**Example**: For Linux on AMD64 architecture:

```commandline
chmod a+x otelcol_linux_amd64
./otelcol_linux_amd64 --config config.yaml
```

To understand the collector's pipeline design, refer to the [Design](collector-for-kafka-design.md) guide.

## OCI Streaming on Ubuntu with systemd

All commands run on the OCI VM over SSH.

### A.1 install dependencies

These packages are used for connectivity testing and producing test messages. They are not required
for the collector itself to run.

```bash
sudo apt-get update
sudo apt-get install -y kafkacat curl netcat-openbsd
```

### A.2 download the collector binary

SOC4Kafka releases are published on GitHub. Download the binary for your target version, make it
executable, and place it in a working directory.

```bash
mkdir -p ~/soc4kafka && cd ~/soc4kafka
wget https://github.com/signalfx/splunk-otel-collector/releases/download/v0.158.0/otelcol_linux_amd64
chmod +x otelcol_linux_amd64
```

!!! note 
    Check the [releases page](https://github.com/splunk/splunk-opentelemetry-collector-for-kafka/releases)
    for newer versions and substitute `v0.158.0` accordingly.

### A.3 create the secrets file

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

### A.4 create the collector config

Create `~/soc4kafka/config.yaml` with the content below. Substitute `<CONSUMER_GROUP>` and `<TOPIC>`
directly in the file - these are not secrets and do not need to be in the env file.

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

### A.5 verify connectivity before starting the collector

Check that the VM can reach Splunk HEC:

```bash
nc -vz <SPLUNK_HEC_HOST> 8088
```

Check that the VM can reach the Kafka broker and authenticate (this is the single best end-to-end
connectivity test - success means DNS, routing, TLS, and SASL all work):

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
before proceeding - the collector will exhibit the same failure.

### A.6 start the collector

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

### A.7 install as a systemd service

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

### A.8 send a test message and confirm in Splunk

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
