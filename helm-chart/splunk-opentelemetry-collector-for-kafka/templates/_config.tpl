{{/*
Generate the complete OTel collector configuration
*/}}
{{- define "soc4kafka.generatedConfig" -}}
extensions:
{{- toYaml .Values.defaults.extensions | nindent 2 }}

receivers:
  {{- range .Values.kafkaReceivers }}
  {{- $receiverName := ternary "kafka" (printf "kafka/%s" .name) (eq .name "main") }}
  {{- $defaults := deepCopy $.Values.defaults.receivers.kafka }}
  {{- $receiverConfig := mustMergeOverwrite $defaults (omit . "name") }}
  {{- if $receiverConfig.auth }}
    {{- if $receiverConfig.auth.plain_text }}
      {{- if and $receiverConfig.auth.plain_text.password (kindIs "map" $receiverConfig.auth.plain_text.password) }}
        {{- if $receiverConfig.auth.plain_text.password.secret }}
          {{- $envVarName := printf "KAFKA_%s_PLAIN_TEXT_PASSWORD" ($receiverName | upper | replace "/" "_" | replace "-" "_") }}
          {{- $_ := set $receiverConfig.auth.plain_text "password" (printf "${%s}" $envVarName) }}
        {{- end }}
      {{- end }}
    {{- end }}
    {{- if $receiverConfig.auth.sasl }}
      {{- if and $receiverConfig.auth.sasl.password (kindIs "map" $receiverConfig.auth.sasl.password) }}
        {{- if $receiverConfig.auth.sasl.password.secret }}
          {{- $envVarName := printf "KAFKA_%s_SASL_PASSWORD" ($receiverName | upper | replace "/" "_" | replace "-" "_") }}
          {{- $_ := set $receiverConfig.auth.sasl "password" (printf "${%s}" $envVarName) }}
        {{- end }}
      {{- end }}
    {{- end }}
    {{- if $receiverConfig.auth.kerberos }}
      {{- if and $receiverConfig.auth.kerberos.password (kindIs "map" $receiverConfig.auth.kerberos.password) }}
        {{- if $receiverConfig.auth.kerberos.password.secret }}
          {{- $envVarName := printf "KAFKA_%s_KERBEROS_PASSWORD" ($receiverName | upper | replace "/" "_" | replace "-" "_") }}
          {{- $_ := set $receiverConfig.auth.kerberos "password" (printf "${%s}" $envVarName) }}
        {{- end }}
      {{- end }}
    {{- end }}
  {{- end }}
  {{ $receiverName }}:
    {{- toYaml $receiverConfig | nindent 4 }}
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