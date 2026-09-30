# Troubleshooting

## Check pod status

```bash
kubectl get pods -l app.kubernetes.io/name=splunk-opentelemetry-collector-for-kafka
```

## View logs

```bash
kubectl logs -l app.kubernetes.io/name=splunk-opentelemetry-collector-for-kafka -f
```

## Check configuration

```bash
# View the generated OpenTelemetry config
kubectl get configmap -l app.kubernetes.io/name=splunk-opentelemetry-collector-for-kafka -o yaml
```

## Health check

```bash
# Port forward to health endpoint
kubectl port-forward svc/<release-name>-splunk-opentelemetry-collector-for-kafka 13133:13133

# Check health
curl http://localhost:13133/

# Or directly via pod:
kubectl port-forward -l app.kubernetes.io/name=splunk-opentelemetry-collector-for-kafka 13133:13133
```

## Common issues

### Pods not starting

- Check if secrets exist and have correct keys
- Verify Kafka brokers are reachable
- Check resource limits
- Review pod events: `kubectl describe pod <pod-name>`

### No data in Splunk

- Verify HTTP Event Collector (HEC) token is correct
- Check Splunk HEC endpoint is accessible
- Review collector logs for errors
- Verify pipeline configuration matches receiver/exporter names
- Check network connectivity from cluster to Splunk HEC endpoint

### Authentication failures

- Ensure secrets exist with key `password` for Kafka auth
- Verify secret names match configuration
- Check username/password are correct
- Review Kafka broker authentication requirements

### Configuration errors

- Verify receiver/exporter names in pipelines match actual names
- Check that at least one receiver, exporter, and pipeline is configured
- Review the generated ConfigMap for syntax errors
- Check Helm template rendering: `helm template . -f values.yaml`
