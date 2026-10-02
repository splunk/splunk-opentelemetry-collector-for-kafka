# Configure TLS

Configure TLS for Kafka receivers and Splunk HTTP Event Collector (HEC) exporters. Both components use the same `tls` options. The options and examples below apply to both.

## Configure TLS for Kafka receivers

When your Kafka brokers use TLS (for example, port 9093 with SSL), configure the `tls` block under each Kafka receiver.

### Configure TLS with a custom CA certificate

Use a custom CA certificate to verify the Kafka broker when the broker uses a private or corporate CA:

```yaml
kafkaReceivers:
  - name: third
    brokers:
      - "kafka-broker-1:9093"
    logs:
      topics:
        - "perf3"
    group_id: "soc4kafka-main3"
    tls:
      insecure_skip_verify: false
      ca_pem: |
        -----BEGIN CERTIFICATE-----
        ...
        G8jotQpS1QbFzo8o3fRN/xQ=
        -----END CERTIFICATE-----
```

### TLS options

| Option | Type | Description |
|--------|------|-------------|
| `insecure_skip_verify` | boolean | When `true`, skips verification of the broker's TLS certificate. Use only for development or testing. Default: `false`. |
| `ca_pem` | string | PEM-encoded CA certificate(s) used to verify the broker's certificate. Use for brokers signed by a private or corporate CA. |
| `ca_file` | string | Path to the CA certificate. For a client, this file verifies the server certificate. Use a mounted Secret when you prefer a file path to inline `ca_pem` content. |
| `cert_file` | string | Path to the TLS certificate for connections that require it. Use only when `insecure` is `false`. |
| `cert_pem` | string | Alternative to `cert_file`. Enter the certificate contents as a string instead of a file path. |
| `key_file` | string | Path to the TLS key for connections that require it. Use only when `insecure` is `false`. |
| `key_pem` | string | Alternative to `key_file`. Enter the key contents as a string instead of a file path. |

The collector supports additional TLS settings and passes them through to its configuration. For the complete reference, see [OpenTelemetry Collector TLS configuration settings](https://github.com/open-telemetry/opentelemetry-collector/blob/main/config/configtls/README.md).

## Configure TLS for Splunk HEC exporters

The Splunk HEC exporter uses TLS when its `endpoint` URL starts with `https://`. The [same `tls` options](#tls-options) apply to Kafka receivers and Splunk HEC exporters.

This example shows the custom CA and certificate verification options:

```yaml
splunkExporters:
  - name: primary
    endpoint: "https://splunk-hec:8088/services/collector"
    token: "your-token"
    tls:
      # ca_pem: | ...               # Optional: PEM for private CA
      # ca_file: /etc/ssl/hec/ca.pem   # Optional: path if mounted via extraVolumes/extraVolumeMounts
```

## Use a CA certificate from a Kubernetes Secret

Mount a Secret that contains the CA certificate by using `extraVolumes` and `extraVolumeMounts` in your values file. Then reference the certificate with `tls.ca_file` in the Kafka receiver or Splunk HEC exporter. You can use the same method for `cert_file` and `key_file`.

1. Create a Secret with the CA certificate and, optionally, the client certificate and key. Use a name that describes its purpose, such as `kafka-ca` or `hec-ca`:

```bash
kubectl create secret generic kafka-ca --from-file=ca.pem=/path/to/ca.pem
# or for HEC:
kubectl create secret generic hec-ca --from-file=ca.pem=/path/to/hec-ca.pem
```

2. In your Helm values, add the volume and mount, and set `tls.ca_file` to the path inside the container.

   **Kafka receiver example:**

```yaml
extraVolumes:
  - name: kafka-ca
    secret:
      secretName: kafka-ca
extraVolumeMounts:
  - name: kafka-ca
    mountPath: /etc/ssl/kafka
    readOnly: true

kafkaReceivers:
  - name: main
    brokers: ["kafka-broker-1:9093"]
    tls:
      insecure_skip_verify: false
      ca_file: /etc/ssl/kafka/ca.pem
```

   **Splunk HEC exporter example:**

```yaml
extraVolumes:
  - name: hec-ca
    secret:
      secretName: hec-ca
extraVolumeMounts:
  - name: hec-ca
    mountPath: /etc/ssl/hec
    readOnly: true

splunkExporters:
  - name: primary
    endpoint: "https://splunk-hec:8088/services/collector"
    token: "your-token"
    tls:
      insecure_skip_verify: false
      ca_file: /etc/ssl/hec/ca.pem
```

   Kubernetes mounts Secret keys as files. If the Secret key is `ca.pem`, the path is `/<mountPath>/ca.pem`. A volume can contain multiple files, such as `ca.pem`, `cert.pem`, and `key.pem`. Reference each file in the `tls` configuration by using `ca_file`, `cert_file`, or `key_file`.


## TLS security recommendations

You can set `insecure_skip_verify: true` for self-signed certificates or internal brokers. Don't use this setting in production because it leaves the connection vulnerable to man-in-the-middle attacks.


## Related topics

- [Configure the Helm chart](collector-for-kafka-kubernetes-configure-helm.md) for chart options.
- [Manage secrets for the Helm chart](collector-for-kafka-configure-secrets.md) for tokens and passwords.
