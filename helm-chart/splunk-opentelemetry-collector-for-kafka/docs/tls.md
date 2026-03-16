# TLS Configuration

This document describes how to configure TLS for **Kafka receivers** (e.g. port 9093) and **Splunk HEC exporters**. Both use the same `tls` options; the options table and patterns below apply to each.

## Kafka Receiver TLS

When your Kafka brokers use TLS (for example, port 9093 with SSL), configure the `tls` block under each Kafka receiver.

### Example: TLS with custom CA

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
| `ca_file` | string | Path to the CA cert. For a client this verifies the server certificate. Use with a mounted secret when you prefer file path over inline `ca_pem`. |
| `cert_file` | string | Path to the TLS cert to use for TLS required connections. Should only be used if `insecure` is set to false. |
| `cert_pem` | string | Alternative to `cert_file`. Provide the certificate contents as a string instead of a filepath. |
| `key_file` | string | Path to the TLS key to use for TLS required connections. Should only be used if `insecure` is set to false. |
| `key_pem` | string | Alternative to `key_file`. Provide the key contents as a string instead of a filepath. |

Additional TLS settings are supported by the collector and passed through to the config. For the full reference, see the [OpenTelemetry Collector TLS Configuration Settings](https://github.com/open-telemetry/opentelemetry-collector/blob/main/config/configtls/README.md).

## Splunk HEC Exporter TLS

The Splunk HEC exporter uses TLS when the `endpoint` URL uses `https://`. The **same `tls` options** as for Kafka receivers apply (see the [TLS options](#tls-options) table above).

Example with custom CA or relaxed verification:

```yaml
splunkExporters:
  - name: primary
    endpoint: "https://splunk-hec:8088/services/collector"
    token: "your-token"
    tls:
      # ca_pem: | ...               # Optional: PEM for private CA
      # ca_file: /etc/ssl/hec/ca.pem   # Optional: path if mounted via extraVolumes/extraVolumeMounts
```

## Using a CA from a Kubernetes secret (file path)

You can mount a Secret containing the CA certificate using `extraVolumes` and `extraVolumeMounts` in your values, then reference it with `tls.ca_file` in the Kafka receiver or Splunk HEC exporter. The same approach works for **`cert_file`** and **`key_file`**.

1. Create a Secret with the CA certificate (and optionally client cert/key; use a name that matches your use case, e.g. `kafka-ca` or `hec-ca`):

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
       # ... logs, group_id, etc.
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

   Secret keys are mounted as files; if your secret key is `ca.pem`, the path is `/<mountPath>/ca.pem`. The same volume can hold multiple files (e.g. `ca.pem`, `cert.pem`, `key.pem`); reference each in `tls` as `ca_file`, `cert_file`, and `key_file`.


## Security recommendations

- **Production:** Set `insecure_skip_verify: false` so the server certificate is verified. When the broker uses a private or corporate CA, provide that CA via `ca_pem` or `ca_file` so the collector only connects to brokers whose certificate is signed by your CA.
- **Development/testing:** You may set `insecure_skip_verify: true` for self-signed or internal brokers. Do not use this in production, as it is vulnerable to man-in-the-middle attacks.


## See also

- [Configuration](configuration.md) – Core configuration options
- [Secret Management](secrets.md) – Managing tokens and passwords securely
