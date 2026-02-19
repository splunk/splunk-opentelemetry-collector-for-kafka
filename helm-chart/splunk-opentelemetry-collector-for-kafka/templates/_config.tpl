{{/*
Generate the complete OTel collector configuration
*/}}
{{- define "soc4kafka.generatedConfig" -}}
extensions:
{{- toYaml .Values.defaults.extensions | nindent 2 }}
{{- if and .Values.collectorLogs.enabled .Values.collectorLogs.forwardToSplunk.enabled }}
  file_storage:
    directory: {{ .Values.collectorLogs.fileStorage.directory }}
    create_directory: {{ .Values.collectorLogs.fileStorage.createDirectory }}
{{- end }}

receivers:
  {{- range .Values.kafkaReceivers }}
  {{- $receiverName := printf "kafka/%s" .name }}
  {{- $defaults := deepCopy $.Values.defaults.receivers.kafka }}
  {{- $receiverInput := omit . "name" }}
  {{- $receiverConfig := mustMergeOverwrite $defaults $receiverInput }}
  {{- if $receiverConfig.auth }}
    {{- if $receiverConfig.auth.plain_text }}
      {{- if $receiverConfig.auth.plain_text.secret }}
        {{- $envVarName := printf "KAFKA_%s_PLAIN_TEXT_PASSWORD" ($receiverName | upper | replace "/" "_" | replace "-" "_") }}
        {{- $plainTextWithoutSecret := omit $receiverConfig.auth.plain_text "secret" }}
        {{- $_ := set $plainTextWithoutSecret "password" (printf "${%s}" $envVarName) }}
        {{- $_ := set $receiverConfig.auth "plain_text" $plainTextWithoutSecret }}
      {{- end }}
    {{- end }}
    {{- if $receiverConfig.auth.sasl }}
      {{- if $receiverConfig.auth.sasl.secret }}
        {{- $envVarName := printf "KAFKA_%s_SASL_PASSWORD" ($receiverName | upper | replace "/" "_" | replace "-" "_") }}
        {{- $saslWithoutSecret := omit $receiverConfig.auth.sasl "secret" }}
        {{- $_ := set $saslWithoutSecret "password" (printf "${%s}" $envVarName) }}
        {{- $_ := set $receiverConfig.auth "sasl" $saslWithoutSecret }}
      {{- end }}
    {{- end }}
    {{- if $receiverConfig.auth.kerberos }}
      {{- if $receiverConfig.auth.kerberos.secret }}
        {{- $envVarName := printf "KAFKA_%s_KERBEROS_PASSWORD" ($receiverName | upper | replace "/" "_" | replace "-" "_") }}
        {{- $kerberosWithoutSecret := omit $receiverConfig.auth.kerberos "secret" }}
        {{- $_ := set $kerberosWithoutSecret "password" (printf "${%s}" $envVarName) }}
        {{- $_ := set $receiverConfig.auth "kerberos" $kerberosWithoutSecret }}
      {{- end }}
    {{- end }}
  {{- end }}
  {{ $receiverName }}:
    {{- toYaml $receiverConfig | nindent 4 }}
  {{- end }}
  {{- if and .Values.collectorLogs.enabled .Values.collectorLogs.forwardToSplunk.enabled }}
  filelog:
    include:
      - /var/log/otelcol/*.log
    start_at: beginning
    storage: file_storage
  {{- end }}
  {{- if .Values.enableMetrics }}
  prometheus:
    config:
      scrape_configs:
        - job_name: 'otel-collector'
          scrape_interval: 1s
          static_configs:
            - targets: ['0.0.0.0:8888']
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
          names: ["otelcol"]
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
  {{- end }}

processors:
{{- toYaml .Values.defaults.processors | nindent 2 }}

exporters:
  {{- range .Values.splunkExporters }}
  {{- $exporterName := ternary "splunk_hec" (printf "splunk_hec/%s" .name) (eq .name "primary") }}
  {{- $tokenValue := printf "${SPLUNK_HEC_TOKEN_%s}" (.name | upper | replace "-" "_") }}
  {{- $exporterConfig := omit . "name" "token" "secret" | mustMergeOverwrite (deepCopy $.Values.defaults.exporters.splunk_hec) }}
  {{- $_ := set $exporterConfig "token" $tokenValue }}
  {{ $exporterName }}:
    {{- toYaml $exporterConfig | nindent 4 }}
  {{- end }}
service:
  extensions:
    {{- range $name, $_ := .Values.defaults.extensions }}
    - {{ $name }}
    {{- end }}
    {{- if and .Values.collectorLogs.enabled .Values.collectorLogs.forwardToSplunk.enabled }}
    - file_storage
    {{- end }}
  {{- if or .Values.collectorLogs.enabled .Values.enableMetrics }}
  telemetry:
    {{- if .Values.collectorLogs.enabled }}
    logs:
      level: {{ .Values.collectorLogs.level }}
      output_paths:
        {{- toYaml .Values.collectorLogs.outputPaths | nindent 8 }}
      error_output_paths:
        {{- toYaml .Values.collectorLogs.errorOutputPaths | nindent 8 }}
    {{- end }}
    {{- if .Values.enableMetrics }}
    metrics:
      level: "detailed"
      readers:
        - pull:
            exporter:
              prometheus:
                host: '0.0.0.0'
                port: 8888
    {{- end }}
  {{- end }}
  pipelines:
    {{- range .Values.pipelines }}
    {{ .type }}/{{ .name }}:
      receivers:
        {{- range .receivers }}
        {{- $receiverName := printf "kafka/%s" . }}
        - {{ $receiverName }}
        {{- end }}
      processors:
        {{- toYaml .processors | nindent 8 }}
      exporters:
        {{- range .exporters }}
        {{- $exporterName := ternary "splunk_hec" (printf "splunk_hec/%s" .) (eq . "primary") }}
        - {{ $exporterName }}
        {{- end }}
    {{- end }}
    {{- if and .Values.collectorLogs.enabled .Values.collectorLogs.forwardToSplunk.enabled }}
    {{- $referencedExporterName := .Values.collectorLogs.forwardToSplunk.exporter | default (index .Values.splunkExporters 0).name }}
    {{- $referencedExporterName = ternary "splunk_hec" (printf "splunk_hec/%s" $referencedExporterName) (eq $referencedExporterName "primary") }}
    logs/internal:
      receivers:
        - filelog
      processors:
        - batch
        - resourcedetection
      exporters:
        - {{ $referencedExporterName }}
    {{- end }}
    {{- if .Values.enableMetrics }}
    {{- $metricsExporterName := ternary "splunk_hec" (printf "splunk_hec/%s" (index .Values.splunkExporters 0).name) (eq (index .Values.splunkExporters 0).name "primary") }}
    metrics:
      receivers:
        - prometheus
        - hostmetrics
      processors:
        - resourcedetection
      exporters:
        - {{ $metricsExporterName }}
    {{- end }}
{{- end }}

{{/*
Merge generated config with user overrides
*/}}
{{- define "soc4kafka.finalConfig" -}}
{{- $generated := fromYaml (include "soc4kafka.generatedConfig" .) -}}
{{- $override := .Values.configOverride | default dict -}}
{{- mustMergeOverwrite $generated $override | toYaml -}}
{{- end }}