# Secret Management

The chart supports secrets for Splunk HEC tokens and Kafka authentication passwords. Secrets can be auto-created or referenced from existing Kubernetes secrets.

## Splunk HEC Tokens

### Auto-Created Secrets

If you provide a `token` value, the chart automatically creates a Kubernetes Secret:

```yaml
splunkExporters:
  - name: primary
    token: "your-hec-token-here"  # Secret will be auto-created
```

The secret will be named `{release-name}-hec-{exporter-name}` with key `splunk-hec-token`.

### Referenced Secrets

Reference an existing Kubernetes secret:

```yaml
splunkExporters:
  - name: primary
    secret: "my-existing-secret"  # Must have key "splunk-hec-token"
```

## Kafka Authentication Passwords

Kafka authentication passwords (plain_text, SASL, or Kerberos) must reference existing Kubernetes secrets:

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

## Creating Secrets Manually

```bash
# Splunk HEC token
kubectl create secret generic my-splunk-hec-secret \
  --from-literal=splunk-hec-token=YOUR_HEC_TOKEN

# Kafka authentication password
kubectl create secret generic kafka-auth-secret \
  --from-literal=password=YOUR_KAFKA_TOKEN
```

## Important Notes

- All secrets are automatically mounted as environment variables and referenced in the OpenTelemetry configuration
- Splunk HEC token secrets use key `splunk-hec-token`
- Kafka authentication secrets use key `password`

