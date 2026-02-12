## Subscribing to topics using regex

The Kafka receiver in the OpenTelemetry Collector supports subscribing to topics using regular expressions. This feature allows you to dynamically subscribe to multiple topics that match a specific pattern, making it easier to manage and collect logs from various sources.
To enable this feature, prefix your topic with the `^` character. This indicates that the topic value is a regex pattern.

### How Regex Topic Subscription Works
The Kafka receiver supports subscribing to topics using regular expressions. When you use a regex pattern (e.g., ^myPrefix.*), the receiver continuously monitors the Kafka cluster for new topics that match the pattern.
The receiver does not only subscribe to existing topics that match the regex, but also detects new matching topics as they are created.

Please Note: Ensure that your regex pattern is valid and correctly formatted to avoid any subscription issues.

### Excluding topics

You can exclude specific topics from being processed using the `kafka.logs.exclude_topics` field. This is useful when your regex pattern matches many topics, but you want to filter out certain ones from log collection.

The `kafka.logs.exclude_topics` field accepts a list of topic names or regex patterns that should be excluded from processing. Topics matching any pattern in the exclude list will be ignored, even if they match the subscription regex pattern. Learn more about regex topics [here](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/kafkareceiver#regex-topic-patterns-with-exclusions).

#### Example Config

```yaml
receivers:
  kafka:
    brokers: [<Brokers>]
    logs:
      topics: 
        - ^<Regex-Topic-Pattern>
      exclude_topics:
        - ^<Regex-Topic-To-Exclude-Pattern>
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
    splunk_app_name: "soc4kafka"
    splunk_app_version: 0.145.0

service:
  pipelines:
    logs:
      receivers: [kafka]
      processors: [batch]
      exporters: [splunk_hec]
```

### Configuration Examples

#### Example 1: Subscribe to All Topics Except System Topics

```yaml
receivers:
  kafka:
    brokers: ["localhost:9092"]
    logs:
      topics:
        - ^.*
      exclude_topics:
        - ^__.*
      encoding: text
```

This configuration subscribes to all topics but excludes any topics starting with __ (Kafka internal topics).

#### Example 2: Subscribe to Application Logs with Pattern Exclusions

```yaml
receivers:
  kafka:
    brokers: ["localhost:9092"]
    logs:
      topics:
        - ^app-.*
      exclude_topics:
        - ^app-test-.*
        - ^app-debug-.*
      encoding: text
```

This configuration subscribes to topics matching `app-*` pattern but excludes topics matching `app-test-*` and `app-debug-*` patterns.

#### Example 3: Multiple Topic Patterns with Exclusions

```yaml
receivers:
  kafka:
    brokers: ["localhost:9092"]
    logs:
      topics:
        - ^logs-.*
        - ^events-.*
      exclude_topics:
        - ^.*-archive$
        - ^.*-old$
      encoding: text
```

This configuration subscribes to multiple regex patterns while excluding topics ending with -archive or -old.

Note that `exclude_topics` doesn't have to be regex, it can be an exact name. When using exact names no `^` at the beginning is necessary.

