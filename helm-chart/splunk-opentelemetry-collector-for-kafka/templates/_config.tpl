{{/*
Generate the complete OTel collector configuration
*/}}
{{- define "soc4kafka.generatedConfig" -}}
extensions:
{{- toYaml .Values.defaults.extensions | nindent 2 }}

receivers:
  {{- range .Values.kafkaReceivers }}
  {{- $receiverName := ternary "kafka" (printf "kafka/%s" .name) (eq .name "main") }}
  {{- $receiverConfig := dict }}
  {{- range $key, $value := . }}
    {{- if ne $key "name" }}
      {{- if eq $key "logs" }}
        {{- $logsConfig := $value | default dict }}
        {{- $defaults := $.Values.defaults.receivers.kafka | default dict }}
        {{- $defaultsLogs := get $defaults "logs" | default dict }}
        {{- $mergedLogs := mustMergeOverwrite $defaultsLogs $logsConfig }}
        {{- $_ := set $receiverConfig "logs" $mergedLogs }}
      {{- else }}
        {{- $_ := set $receiverConfig $key $value }}
      {{- end }}
    {{- end }}
  {{- end }}
  {{- $defaults := $.Values.defaults.receivers.kafka | default dict }}
  {{- $merged := mustMergeOverwrite $defaults $receiverConfig }}
  {{ $receiverName }}:
    {{- toYaml $merged | nindent 4 }}
  {{- end }}

processors:
{{- toYaml .Values.defaults.processors | nindent 2 }}

exporters:
  {{- range .Values.splunkExporters }}
  {{- $exporterName := ternary "splunk_hec" (printf "splunk_hec/%s" .name) (eq .name "primary") }}
  {{- $exporterConfig := dict "endpoint" .endpoint "source" .source "sourcetype" .sourcetype "index" .index }}
  {{- $envVarName := printf "SPLUNK_HEC_TOKEN_%s" (.name | upper | replace "-" "_") }}
  {{- $tokenValue := printf "${%s}" $envVarName }}
  {{- if .token }}
    {{- $tokenStr := toString .token }}
    {{- if hasPrefix "${" $tokenStr }}
      {{- $tokenValue = $tokenStr }}
    {{- else }}
      {{- $tokenValue = printf "${%s}" $envVarName }}
    {{- end }}
  {{- end }}
  {{- $_ := set $exporterConfig "token" $tokenValue }}
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