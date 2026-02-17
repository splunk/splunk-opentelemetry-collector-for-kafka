{{/*
Generate the complete OTel collector configuration
*/}}
{{- define "soc4kafka.generatedConfig" -}}
extensions:
{{- toYaml .Values.defaults.extensions | nindent 2 }}

receivers:
  {{- range .Values.kafkaReceivers }}
  {{- $receiverName := ternary "kafka" (printf "kafka/%s" .name) (eq .name "main") }}
  {{- $receiverConfig := dict "brokers" .brokers }}
  {{- $logsConfig := dict "topics" .topics }}
  {{- if .encoding }}
    {{- $_ := set $logsConfig "encoding" .encoding }}
  {{- end }}
  {{- $defaults := $.Values.defaults.receivers.kafka | default dict }}
  {{- $defaultsLogs := get $defaults "logs" | default dict }}
  {{- $mergedLogs := mustMergeOverwrite $defaultsLogs $logsConfig }}
  {{- $_ := set $receiverConfig "logs" $mergedLogs }}
  {{- $merged := mustMergeOverwrite $defaults $receiverConfig }}
  {{ $receiverName }}:
    {{- toYaml $merged | nindent 4 }}
  {{- end }}

processors:
{{- toYaml .Values.defaults.processors | nindent 2 }}

exporters:
  {{- range .Values.splunkExporters }}
  {{- $exporterName := ternary "splunk_hec" (printf "splunk_hec/%s" .name) (eq .name "primary") }}
  {{- $tokenValue := .token }}
  {{- if not (hasPrefix "${" .token) }}
    {{- $tokenValue = printf "${SPLUNK_HEC_TOKEN_%s}" (.name | upper | replace "-" "_") }}
  {{- end }}
  {{- $exporterConfig := dict "endpoint" .endpoint "token" $tokenValue "source" .source "sourcetype" .sourcetype "index" .index }}
  {{- if .tls }}
    {{- $_ := set $exporterConfig "tls" .tls }}
  {{- end }}
  {{- $defaults := $.Values.defaults.exporters.splunk_hec | default dict }}
  {{- $merged := mustMergeOverwrite $defaults $exporterConfig }}
  {{ $exporterName }}:
    {{- toYaml $merged | nindent 4 }}
  {{- end }}

service:
  extensions:
    {{- range $name, $_ := .Values.defaults.extensions }}
    - {{ $name }}
    {{- end }}
  pipelines:
    {{- range .Values.pipelines }}
    {{ .type }}/{{ .name }}:
      receivers:
        {{- range .receivers }}
        {{- $receiverName := ternary "kafka" (printf "kafka/%s" .) (eq . "main") }}
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
{{- end }}

{{/*
Merge generated config with user overrides
*/}}
{{- define "soc4kafka.finalConfig" -}}
{{- $generated := fromYaml (include "soc4kafka.generatedConfig" .) -}}
{{- $override := .Values.configOverride | default dict -}}
{{- mustMergeOverwrite $generated $override | toYaml -}}
{{- end }}