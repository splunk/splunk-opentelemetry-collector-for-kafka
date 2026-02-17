{{/*
Generate the complete OTel collector configuration
*/}}
{{- define "soc4kafka.generatedConfig" -}}
extensions:
{{- toYaml .Values.defaults.extensions | nindent 2 }}

receivers:
  {{- range .Values.kafkaReceivers }}
  {{- $receiverName := ternary "kafka" (printf "kafka/%s" .name) (eq .name "main") }}
  {{ $receiverName }}:
    brokers:
      {{- toYaml .brokers | nindent 6 }}
    logs:
      topics:
      {{- range .topics }}
        - {{ . }}
      {{- end }}
      encoding: {{ .encoding | default $.Values.defaults.receivers.kafka.encoding }}
    {{- if .consumerGroup }}
    group_id: {{ .consumerGroup }}
    {{- end }}
    initial_offset: {{ .initialOffset | default $.Values.defaults.receivers.kafka.initial_offset }}
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
  {{ $exporterName }}:
    endpoint: {{ .endpoint | quote }}
    token: {{ $tokenValue | quote }}
    source: {{ .source | quote }}
    sourcetype: {{ .sourcetype | quote }}
    index: {{ .index | quote }}
    tls:
      insecure_skip_verify: {{ .insecureSkipVerify | default $.Values.defaults.exporters.splunk_hec.tls.insecure_skip_verify }}
    headers:
      {{- toYaml $.Values.defaults.exporters.splunk_hec.headers | nindent 6 }}
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