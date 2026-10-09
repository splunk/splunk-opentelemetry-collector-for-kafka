# Install the Splunk Distribution of OpenTelemetry Collector for Kafka manually

Run the Splunk Distribution of OpenTelemetry Collector for Kafka from a downloaded package and a configuration file. For a command-by-command walkthrough that uses OCI Streaming as the source, see [Deploy the Splunk Distribution of OpenTelemetry Collector for Kafka for OCI Streaming](collector-for-kafka-oci-streaming.md).

## Download the Splunk OpenTelemetry Collector package

The Splunk Distribution of OpenTelemetry Collector for Kafka runs in the Splunk OpenTelemetry Collector package. Download the package for your platform from the [Splunk OpenTelemetry Collector releases](https://github.com/signalfx/splunk-otel-collector/releases). Release numbers start with `v`.

For example, on Linux with an AMD64 architecture, run the following `wget` command:

```commandline
wget https://github.com/signalfx/splunk-otel-collector/releases/download/v0.158.0/otelcol_linux_amd64
```

## Run the package with a configuration file

Run the package with a completed Splunk Distribution of OpenTelemetry Collector for Kafka configuration template.
For a minimal starting point, see [Create a minimal configuration template](../configure/collector-for-kafka-configure.md#create-a-minimal-configuration-template).

```commandline
./<otel_package> --config <config_file>
```

!!! note

    Make sure that the file has executable permissions before you run the command. On Linux, use the following command to add executable permissions:

```commandline
chmod a+x <otel_package>
```

For example, on Linux with an AMD64 architecture, run:

```commandline
chmod a+x otelcol_linux_amd64
./otelcol_linux_amd64 --config config.yaml
```

For information about the pipeline, see [Understand the Splunk Distribution of OpenTelemetry Collector for Kafka design](../collector-for-kafka-design.md).
