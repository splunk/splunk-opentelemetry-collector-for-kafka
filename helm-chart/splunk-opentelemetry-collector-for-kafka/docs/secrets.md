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
  --from-literal=password=your-kafka-password
```

## Important Notes

- All secrets are automatically mounted as environment variables and referenced in the OpenTelemetry configuration
- Splunk HEC token secrets use key `splunk-hec-token`
- Kafka authentication secrets use key `password`

## Using Vault for Secrets

The chart works seamlessly with HashiCorp Vault when using the [External Secrets Operator](https://external-secrets.io/). The External Secrets Operator syncs Vault secrets to Kubernetes Secrets, which the chart can then reference normally.

### Example with External Secrets Operator

1. Create an ExternalSecret that syncs from Vault:

```yaml
apiVersion: external-secrets.io/v1beta1
kind: ExternalSecret
metadata:
  name: splunk-hec-secret
spec:
  secretStoreRef:
    name: vault-backend
    kind: SecretStore
  target:
    name: splunk-hec-secret
    creationPolicy: Owner
  data:
  - secretKey: splunk-hec-token
    remoteRef:
      key: secret/data/splunk
      property: hec-token
```

2. Reference the synced Kubernetes Secret in your values.yaml:

```yaml
splunkExporters:
  - name: primary
    secret: "splunk-hec-secret"  # Created by External Secrets Operator
```

The chart will automatically mount and use the secret. The External Secrets Operator keeps the Kubernetes Secret in sync with Vault.
