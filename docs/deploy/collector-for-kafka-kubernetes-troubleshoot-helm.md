# Troubleshoot the Helm deployment

Use these checks to investigate common issues with the Helm deployment.

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

## Check collector health

```bash
# Port forward to health endpoint
kubectl port-forward svc/<release-name>-splunk-opentelemetry-collector-for-kafka 13133:13133

# Check health
curl http://localhost:13133/

# Or connect directly to a pod:
kubectl port-forward -l app.kubernetes.io/name=splunk-opentelemetry-collector-for-kafka 13133:13133
```

## Common issues

### Pods not starting

- Confirm that the Secrets exist and use the required keys.
- Confirm that the Kafka brokers are reachable.
- Check the resource limits.
- Review pod events by running `kubectl describe pod <pod-name>`.

### No data in Splunk

- Confirm that the HTTP Event Collector (HEC) token is correct.
- Confirm that the Splunk HEC endpoint is accessible.
- Review collector logs for errors.
- Confirm that the pipeline configuration uses the correct receiver and exporter names.
- Check network connectivity from the cluster to the Splunk HEC endpoint.

### Authentication failures

- Confirm that the Secrets exist and include the `password` key for Kafka authentication.
- Confirm that the Secret names match the configuration.
- Check that the username and password are correct.
- Review the Kafka broker authentication requirements.

### Configuration errors

- Confirm that the receiver and exporter names in the pipelines match their configured names.
- Confirm that the configuration includes at least one receiver, exporter, and pipeline.
- Review the generated ConfigMap for syntax errors.
- Render the Helm templates by running `helm template . -f values.yaml`.
