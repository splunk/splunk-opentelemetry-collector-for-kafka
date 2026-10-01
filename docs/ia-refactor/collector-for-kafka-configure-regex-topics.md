## Subscribing to topics using regex

Use a regular expression to subscribe to Kafka topics that match a pattern.
Prefix the pattern with the `^` character to identify it as a regular expression.

### How regex topic subscription works
The Kafka receiver subscribes to existing topics that match the pattern and detects new matching topics as they are created. For example, `^myPrefix.*` matches topics that begin with `myPrefix`.

!!! note
    Ensure that your regex pattern is valid and correctly formatted to avoid any subscription issues.

### Excluding topics

You can exclude specific topics from being processed using the `kafka.logs.exclude_topics` field. This is useful when your regex pattern matches many topics, but you want to filter out certain ones from log collection.

The `kafka.logs.exclude_topics` field accepts a list of topic names or regex patterns that should be excluded from processing. Topics matching any pattern in the exclude list will be ignored, even if they match the subscription regex pattern. Learn more about regex topics [Kafka receiver regex topic exclusions documentation](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/kafkareceiver#regex-topic-patterns-with-exclusions).

#### Example config

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


exporters:
  splunk_hec:
    token: <Splunk HEC Token>
    endpoint: <Splunk HEC Endpoint>
    source: <Source>
    sourcetype: <Sourcetype>
    index: <Splunk index>
    splunk_app_name: "soc4kafka"
    sending_queue:
      enabled: true
      num_consumers: 10
      queue_size: 10000
      block_on_overflow: true
      sizer: items
      batch:
        min_size: 1000

service:
  pipelines:
    logs:
      receivers: [kafka]
      exporters: [splunk_hec]
```

### Configuration examples

#### Example 1: subscribe to all topics except system topics

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

#### Example 2: subscribe to application logs with pattern exclusions

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

#### Example 3: multiple topic patterns with exclusions

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
