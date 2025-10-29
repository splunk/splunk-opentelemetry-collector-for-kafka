## Subscribing to topics using regex

The Kafka receiver in the OpenTelemetry Collector supports subscribing to topics using regular expressions. This feature allows you to dynamically subscribe to multiple topics that match a specific pattern, making it easier to manage and collect logs from various sources.
To enable this feature, prefix your topic with the `^` character. This indicates that the topic value is a regex pattern.

### How Regex Topic Subscription Works
The Kafka receiver supports subscribing to topics using regular expressions. When you use a regex pattern (e.g., ^myPrefix.*), the receiver continuously monitors the Kafka cluster for new topics that match the pattern.
The receiver does not only subscribe to existing topics that match the regex, but also detects new matching topics as they are created.

Please Note: Ensure that your regex pattern is valid and correctly formatted to avoid any subscription issues.

### Example Config
```yaml
receivers:
  kafka:
    brokers: [<Brokers>]
    topic: ^<Regex-Topic-Pattern>
    encoding: <Encoding>

processors:
  batch:

exporters:
  splunk_hec:
    token: <Splunk HEC Token>
    endpoint: <Splunk HEC Endpoint>
    source: <Source>
    sourcetype: <Sourcetype>
    index: <Splunk index>
    headers:
      "__splunk_app_name": "soc4kafka"

service:
  pipelines:
    logs:
      receivers: [kafka]
      processors: [batch]
      exporters: [splunk_hec]
```