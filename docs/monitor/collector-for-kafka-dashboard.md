# Monitor the Splunk Distribution of OpenTelemetry Collector for Kafka with a dashboard

Use the Splunk Distribution of OpenTelemetry Collector for Kafka dashboard to monitor Kafka metrics. The dashboard combines data from three sources to show the environment's performance and health.

## Dashboard controls

The dashboard has seven tabs. Each tab contains graphs, gauges, and inputs for the data it displays. All tabs share these three controls:

- **Time Range**: Select the date and time range for metrics. Some gauges, such as Active Collectors, always show the latest information.
- **Log Indexes**: Choose the indexes that contain dashboard event data and metrics. Select both an event index and a metrics index to display all dashboard features.
- **Time Span**: Select the interval for aggregating data in time-based graphs.

![Dashboard filters for time range, log indexes, and time span](../images/dashboard/global-inputs.png)

### General

The **General** tab shows Splunk Distribution of OpenTelemetry Collector for Kafka health, including a table of active collector instances and a gauge for active Kafka brokers. The Active Collectors gauge shows data from the past 5 minutes.

![General dashboard tab showing active collector instances and Kafka brokers](../images/dashboard/global-tab-active-instances.png)

Graphs show the number of messages each instance receives and exports. They aggregate values across the instance's receivers and exporters. Select an instance to view details.

![Messages received and exported by each collector instance](../images/dashboard/global-tab-receiver-exporter.png)
![Receiver and exporter metrics for a selected collector instance](../images/dashboard/global-tab-receiver-exporter-per-instance.png)

The tab also shows exporter queue information.

![Exporter queue metrics in the General dashboard tab](../images/dashboard/global-tab-queue.png)

### Kafka

The **Kafka** tab shows topics, topic replicas, consumer groups, Kafka offsets, and consumer lag.

![Kafka dashboard tab showing topic and consumer group metrics](../images/dashboard/kafka-tab.png)

For the last two charts, select the appropriate topic for each consumer group. Otherwise, the charts show no results.

### CPU, memory, disk, network

The next four tabs show system metrics for the machines that run instances of the Splunk Distribution of OpenTelemetry Collector for Kafka.

The **CPU** tab shows the number of logical CPU cores, process CPU utilization, and system CPU utilization. Select the task types to include in the statistics. By default, the dashboard includes system and CPU modes.

![CPU metrics in the Splunk Distribution of OpenTelemetry Collector for Kafka dashboard](../images/dashboard/cpu-tab.png)
![Additional CPU metrics in the Splunk Distribution of OpenTelemetry Collector for Kafka dashboard](../images/dashboard/cpu-tab-2.png)

The **Memory** tab shows system memory utilization, total available memory, and system and memory usage. Select the memory metrics to include in the graphs.

![Memory metrics in the Splunk Distribution of OpenTelemetry Collector for Kafka dashboard](../images/dashboard/memory-tab.png)

The **Disk** tab shows disk usage. The filesystem utilization gauge shows the selected filesystem's space usage. Select a filesystem from the list.

![Disk usage metrics in the Splunk Distribution of OpenTelemetry Collector for Kafka dashboard](../images/dashboard/disk-tab.png)

The **Network** tab shows network traffic.

![Network metrics in the Splunk Distribution of OpenTelemetry Collector for Kafka dashboard](../images/dashboard/network-tab.png)

### Events

The **Events** tab shows data about events received by the Splunk platform. Use the lists to filter events by hostname, source, and sourcetype, and to view the distribution of those values among indexed events. In environments with high data ingress, this tab can take longer to load.

![Events dashboard tab with filters for hostname, source, and sourcetype](../images/dashboard/events-tab.png)

## Install the dashboard

### Configure the Splunk Distribution of OpenTelemetry Collector for Kafka

The dashboard uses these three telemetry data sources:

1. The [hostmetrics receiver](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/hostmetricsreceiver)
2. The [kafkametrics receiver](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/kafkametricsreceiver)
3. [Internal telemetry](https://opentelemetry.io/docs/collector/internal-telemetry/#lists-of-internal-metrics) from the OpenTelemetry Collector

Configure all three sources for the dashboard to display data.

1. Configure receivers:

```yaml
receivers:
  prometheus:
    config:
      scrape_configs:
        - job_name: 'otel-collector'
          scrape_interval: 1s
          static_configs:
            - targets: ['0.0.0.0:8888']
  kafkametrics:
    brokers: [<kafka-broker-address>:<port>]
    protocol_version: 2.0.0
    scrapers:
      - brokers
      - topics
      - consumers
  hostmetrics:
    collection_interval: 1s
    scrapers:
      cpu:
        metrics:
          system.cpu.utilization:
            enabled: true
          system.cpu.logical.count:
            enabled: true
      memory:
        metrics:
          system.memory.utilization:
            enabled: true
          system.memory.limit:
            enabled: true
      process:
        mute_process_all_errors: true
        include:
          names: [ "otelcol_linux_amd64" ] # name of the collector binary
          match_type: "regexp"
        metrics:
          process.memory.utilization:
            enabled: true
          process.cpu.utilization:
            enabled: true
      filesystem:
        metrics:
          system.filesystem.utilization:
            enabled: true
      disk:
        metrics:
          system.disk.io:
            enabled: true
      network:
```

2. Add the `resourcedetection` processor to filter data by host.

```yaml
processors:
  resourcedetection:
    detectors: ["system"]
    system:
      hostname_sources: ["os"]
```

3. Add an exporter for metrics.

```yaml
exporters:
  splunk_hec/metrics:
    token: <hec-token>
    endpoint: http://<splunk-address>:8088/services/collector
    source: <source>
    sourcetype: <sourcetype>
    index: <metrics-index>
    splunk_app_name: "soc4kafka"
    sending_queue:
      enabled: true
      num_consumers: 10
      queue_size: 10000
      block_on_overflow: true
      sizer: items
      batch:
        min_size: 1000
```

Create a metrics index:

![Splunk metric index configuration](../images/dashboard/metric-index.png)

4. Create the `telemetry` service.
```yaml
service:
  telemetry:
    metrics:
      level: "detailed"
      readers:
        - pull:
            exporter:
              prometheus:
                host: '0.0.0.0'
                port: 8888
```

5. Add a pipeline for metrics.

```yaml
service:
  pipelines:
    metrics:
      receivers: [ prometheus, hostmetrics, kafkametrics ]
      processors: [resourcedetection]
      exporters: [ splunk_hec/metrics ]
```

   Configure other receivers and exporters as usual.

### Create the dashboard in Splunk Web

!!! note

    This dashboard is available only in Splunk version 9.4.0 or later.

1. In Splunk Web, select **Search & Reporting > Dashboards**.
2. Select **Create New Dashboard**. Select **Dashboard Studio** and **Grid**.
3. In the edit view, open the **source code editor** and replace the initial configuration with the contents of this [JSON file](../../dashboards/SOC4Kafka-health-dashboard.json).

4. Save your changes.
