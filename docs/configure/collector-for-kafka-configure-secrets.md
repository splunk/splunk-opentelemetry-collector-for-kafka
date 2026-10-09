# Manage secrets for the Helm chart

The chart supports secrets for Splunk HTTP Event Collector (HEC) tokens and Kafka authentication passwords. You can have the chart create a Secret or reference an existing Kubernetes Secret.

## Splunk HEC tokens

### Secrets created by the chart

If you provide a `token` value, the chart automatically creates a Kubernetes Secret:

```yaml
splunkExporters:
  - name: primary
    token: "your-splunk-hec-token"  # Secret will be auto-created
```

The chart names the Secret `{release-name}-hec-{exporter-name}` and uses the `splunk-hec-token` key.

### Reference existing Secrets

Reference an existing Kubernetes Secret:

```yaml
splunkExporters:
  - name: primary
    secret: "my-existing-secret"  # Must have key "splunk-hec-token"
```

## Kafka authentication passwords

For Kafka authentication with `plain_text`, SASL, or Kerberos, reference an existing Kubernetes Secret:

```yaml
kafkaReceivers:
  - name: main
    brokers:
      - "kafka-broker:9092"
    auth:
      plain_text:
        username: "kafka-user"
        secret: "kafka-auth-secret"  # Secret must have key "password"
      # Or for SASL:
      # sasl:
      #   username: "kafka-user"
      #   secret: "kafka-sasl-secret"
      # Or for Kerberos:
      # kerberos:
      #   secret: "kafka-kerberos-secret"
    logs:
      topics:
        - "my-topic"
```

## Creating secrets manually

```bash
# Splunk HEC token
kubectl create secret generic my-splunk-hec-secret \
  --from-literal=splunk-hec-token=YOUR_HEC_TOKEN

# Kafka authentication password
kubectl create secret generic kafka-auth-secret \
  --from-literal=password=YOUR_KAFKA_TOKEN
```

## Secret requirements

- The chart mounts all secrets as environment variables and references them in the OpenTelemetry configuration.
- Use the `splunk-hec-token` key for Splunk HEC token Secrets.
- Use the `password` key for Kafka authentication Secrets.
