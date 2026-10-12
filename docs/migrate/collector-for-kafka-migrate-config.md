# Map Splunk Connect for Kafka settings to the Splunk Distribution of OpenTelemetry Collector for Kafka

## Configuration mapping

You cannot transfer settings directly from Splunk Connect for Kafka to the Splunk Distribution of OpenTelemetry Collector for Kafka because the products use different architectures. However, many settings have equivalent options. Use the table below to map each setting.

In the Splunk Distribution of OpenTelemetry Collector for Kafka, configure settings on individual receivers (data sources) and exporters (data destinations). Then connect the components in a pipeline.

## Settings supported by the Splunk Distribution of OpenTelemetry Collector for Kafka

| Splunk Connect for Kafka field | Splunk Distribution of OpenTelemetry Collector for Kafka setting | Description |
|----------------------------------------------|--------------------------------------------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `topics` | `receivers.kafka.logs.topics` | Configure one topic per Kafka receiver. Add multiple receivers to a pipeline to collect from multiple topics. See [Collector design](../collector-for-kafka-design.md). |
| `topics.regex` | `receivers.kafka.logs.topics` | Prefix the topic pattern with `^` to use a regular expression. See [Subscribe to topics with regular expressions](../configure/collector-for-kafka-configure-regex-topics.md). |
| `splunk.indexes` | `exporters.splunk_hec.index` | Configure one `index` per Splunk HEC exporter. Add multiple exporters to a pipeline to send data to multiple indexes. |
| `splunk.sources` | `exporters.splunk_hec.source` | Configure one `source` per Splunk HEC exporter. Add multiple exporters to a pipeline to use multiple sources. |
| `splunk.sourcetypes` | `exporters.splunk_hec.sourcetype` | Configure one `sourcetype` per Splunk HEC exporter. Add multiple exporters to a pipeline to use multiple sourcetypes. |
| `splunk.hec.uri` | `exporters.splunk_hec.endpoint` | Splunk HEC endpoint. |
| `splunk.hec.token` | `exporters.splunk_hec.token` | Token for authenticating with Splunk HEC. |
| `splunk.hec.raw` | `exporters.splunk_hec.export_raw` | Sends only the log body to a Splunk HEC raw endpoint. |
| `splunk.hec.ssl.validate.certs` | `exporters.splunk_hec.tls.insecure_skip_verify` | Specifies whether to skip certificate validation for the HEC endpoint over HTTPS. The default is `false`. |
| `splunk.hec.http.keepalive` | `exporters.splunk_hec.health_check_enabled` | Specifies whether to check Splunk HEC health when the exporter starts. |
| `splunk.hec.max.http.connection.per.channel` | `exporters.splunk_hec.max_idle_conns` | Maximum number of simultaneous HTTP connections to Splunk HEC. The default is 100. |
| `splunk.hec.max.batch.size` | `splunk_hec.sending_queue.batch.min_size` | Number of spans, metric data points, or log records to batch before sending. The default is 1000. |
| `splunk.hec.event.timeout` | `splunk.timeout` | Timeout for Splunk exporter operations. |
| `splunk.hec.socket.timeout` | `splunk.socket.timeout` | Socket timeout for Splunk exporter operations. |
| `splunk.header.support` | `receivers.kafka.header_extraction.extract_headers` | When `true`, the Kafka receiver parses headers for use as metadata in Splunk events. See [Extract data from headers](../configure/collector-for-kafka-configure-extract-data.md#extract-data-from-headers). |
| `splunk.header.custom` | `receivers.kafka.header_extraction.headers` | Extract custom headers and use them with custom processors. |
| `splunk.header.index` | `exporters.splunk_hec.otel_attrs_to_hec_metadata.index` | Maps Kafka header values to Splunk index metadata by using custom processors. |
| `splunk.header.source` | `exporters.splunk_hec.otel_attrs_to_hec_metadata.source` | Maps Kafka header values to Splunk source metadata by using custom processors. |
| `splunk.header.sourcetype` | `exporters.splunk_hec.otel_attrs_to_hec_metadata.sourcetype` | Maps Kafka header values to Splunk sourcetype metadata by using custom processors. |
| `splunk.header.host` | `exporters.splunk_hec.otel_attrs_to_hec_metadata.host` | Maps Kafka header values to Splunk host metadata by using custom processors. |
| `enable.timestamp.extraction` | `processors.timestamp` | Configure timestamp extraction with processors. See [Extract timestamps](../configure/collector-for-kafka-configure-extract-data.md#extract-timestamps). |
| `timestamp.regex` | `processors.timestamp.regex` | Regular expression for extracting timestamps from log data. |
| `timestamp.format` | `processors.timestamp.format` | Format of extracted timestamps. |
| `timestamp.timezone` | `processors.timestamp.timezone` | Time zone for extracted timestamps. |

## Settings without an equivalent in the Splunk Distribution of OpenTelemetry Collector for Kafka

| Splunk Connect for Kafka field | Description |
|---------------------------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `connector.class` | The Splunk Distribution of OpenTelemetry Collector for Kafka uses receivers, processors, and exporters instead of a connector class. |
| `tasks.max` | Configure scaling differently. See [Scale the Splunk Distribution of OpenTelemetry Collector for Kafka](../operate/collector-for-kafka-scale.md). |
| `splunk.hec.raw.line.breaker` | Configure line breaking with custom processors. |
| `splunk.hec.json.event.enrichment` | Configure JSON enrichment with custom processors. |
| `splunk.hec.auto.extract.timestamp` | Configure timestamp extraction with processors. See [Extract timestamps](../configure/collector-for-kafka-configure-extract-data.md#extract-timestamps). |
| `value.converter` | Not supported by the Splunk Distribution of OpenTelemetry Collector for Kafka. |
| `value.converter.schema.registry.url` | Not supported by the Splunk Distribution of OpenTelemetry Collector for Kafka. |
| `value.converter.schemas.enable` | Not supported by the Splunk Distribution of OpenTelemetry Collector for Kafka. |
| `key.converter` | Not supported by the Splunk Distribution of OpenTelemetry Collector for Kafka. |
| `key.converter.schema.registry.url` | Not supported by the Splunk Distribution of OpenTelemetry Collector for Kafka. |
| `key.converter.schemas.enable` | Not supported by the Splunk Distribution of OpenTelemetry Collector for Kafka. |
| `splunk.hec.ack.enabled` | Not supported by the Splunk Distribution of OpenTelemetry Collector for Kafka. |
| `splunk.hec.ack.poll.interval` | Not supported by the Splunk Distribution of OpenTelemetry Collector for Kafka. |
| `splunk.hec.ack.poll.threads` | Not supported by the Splunk Distribution of OpenTelemetry Collector for Kafka. |
| `splunk.hec.total.channels` | The Splunk Distribution of OpenTelemetry Collector for Kafka does not use channels. |
| `splunk.hec.threads` | Threading is managed differently and does not require explicit configuration. |
| `splunk.hec.track.data` | Handle data tracking and debugging with custom processors or external monitoring tools. |
| `splunk.hec.json.event.formatted` | Send events already formatted for HEC with the `exporters.splunk_hec.export_raw` option. |
| `splunk.hec.ssl.trust.store.path` | Trust store configuration is not supported by the Splunk Distribution of OpenTelemetry Collector for Kafka. |
| `splunk.hec.ssl.trust.store.password` | |
| `kerberos.user.principal` | Kerberos authentication is supported by the Kafka receiver. For details, see the [Kafka receiver configuration](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/f1d708538c1038aacf60f6659ed23189481358e4/receiver/kafkareceiver/README.md?plain=1#L72). |
| `kerberos.keytab.path` | |
