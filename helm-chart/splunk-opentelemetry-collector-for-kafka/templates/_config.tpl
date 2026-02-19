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
  {{- if and .Values.collectorLogs.enabled .Values.collectorLogs.forwardToSplunk.enabled }}
  {{- $firstExporter := index $.Values.splunkExporters 0 }}
  {{- $internalEndpoint := .Values.collectorLogs.forwardToSplunk.endpoint | default $firstExporter.endpoint }}
  {{- $internalTokenValue := printf "${SPLUNK_HEC_TOKEN_INTERNAL_LOGS}" }}
  splunk_hec/internal_logs:
    endpoint: {{ $internalEndpoint }}
    token: {{ $internalTokenValue }}
    index: {{ .Values.collectorLogs.forwardToSplunk.index }}
    source: {{ .Values.collectorLogs.forwardToSplunk.source }}
    sourcetype: {{ .Values.collectorLogs.forwardToSplunk.sourcetype }}
    splunk_app_name: "soc4kafka"
  {{- end }}

service:
  extensions:
    {{- range $name, $_ := .Values.defaults.extensions }}
    - {{ $name }}
    {{- end }}
    {{- if and .Values.collectorLogs.enabled .Values.collectorLogs.forwardToSplunk.enabled }}
    - file_storage
    {{- end }}
  {{- if .Values.collectorLogs.enabled }}
  telemetry:
    logs:
      level: {{ .Values.collectorLogs.level }}
      output_paths:
        {{- toYaml .Values.collectorLogs.outputPaths | nindent 8 }}
      error_output_paths:
        {{- toYaml .Values.collectorLogs.errorOutputPaths | nindent 8 }}
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
    logs/internal:
      receivers:
        - filelog
      processors:
        - batch
        - resourcedetection
      exporters:
        - splunk_hec/internal_logs
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