{{/*
Expand the name of the chart.
*/}}
{{- define "soc4kafka.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
*/}}
{{- define "soc4kafka.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "soc4kafka.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "soc4kafka.labels" -}}
helm.sh/chart: {{ include "soc4kafka.chart" . }}
{{ include "soc4kafka.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "soc4kafka.selectorLabels" -}}
app.kubernetes.io/name: {{ include "soc4kafka.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "soc4kafka.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "soc4kafka.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Calculate ConfigMap hash to trigger pod restarts on config changes
*/}}
{{- define "soc4kafka.configMapHash" -}}
{{- include "soc4kafka.finalConfig" . | sha256sum | trunc 8 }}
{{- end }}

{{/*
Calculate secrets hash to trigger pod restarts on secret changes
Includes:
- Token values for auto-created Splunk HEC secrets
- Secret name/key references for all secrets (Splunk HEC and Kafka auth)
*/}}
{{- define "soc4kafka.secretsHash" -}}
{{- $secretData := list }}
{{- range .Values.splunkExporters }}
  {{- $secretName := .secret | default (printf "%s-hec-%s" (include "soc4kafka.fullname" $) .name) }}
  {{- if .secret }}
    {{- $secretData = append $secretData (printf "splunk-hec:%s:splunk-hec-token" $secretName) }}
  {{- else if and .token (not (hasPrefix "${" (toString .token))) }}
    {{- $secretData = append $secretData (printf "splunk-hec:%s:splunk-hec-token:%s" $secretName (toString .token)) }}
  {{- end }}
{{- end }}
{{- range .Values.kafkaReceivers }}
  {{- if .auth }}
    {{- if .auth.plain_text }}
      {{- if .auth.plain_text.secret }}
        {{- $secretData = append $secretData (printf "kafka-plain-text:%s:password" .auth.plain_text.secret) }}
      {{- end }}
    {{- end }}
    {{- if .auth.sasl }}
      {{- if .auth.sasl.secret }}
        {{- $secretData = append $secretData (printf "kafka-sasl:%s:password" .auth.sasl.secret) }}
      {{- end }}
    {{- end }}
    {{- if .auth.kerberos }}
      {{- if .auth.kerberos.secret }}
        {{- $secretData = append $secretData (printf "kafka-kerberos:%s:password" .auth.kerberos.secret) }}
      {{- end }}
    {{- end }}
  {{- end }}
{{- end }}
{{- if and .Values.collectorLogs.enabled .Values.collectorLogs.forwardToSplunk.enabled }}
  {{- $internalSecretName := "" }}
  {{- if .Values.collectorLogs.forwardToSplunk.secret }}
    {{- $internalSecretName = .Values.collectorLogs.forwardToSplunk.secret }}
    {{- $secretData = append $secretData (printf "splunk-hec-internal:%s:splunk-hec-token" $internalSecretName) }}
  {{- else if and .Values.collectorLogs.forwardToSplunk.token (not (hasPrefix "${" .Values.collectorLogs.forwardToSplunk.token)) }}
    {{- $internalSecretName = printf "%s-hec-internal-logs" (include "soc4kafka.fullname" .) }}
    {{- $secretData = append $secretData (printf "splunk-hec-internal:%s:splunk-hec-token:%s" $internalSecretName (toString .Values.collectorLogs.forwardToSplunk.token)) }}
  {{- else }}
    {{- $firstExporter := index .Values.splunkExporters 0 }}
    {{- $internalSecretName = $firstExporter.secret | default (printf "%s-hec-%s" (include "soc4kafka.fullname" .) $firstExporter.name) }}
    {{- $secretData = append $secretData (printf "splunk-hec-internal:%s:splunk-hec-token" $internalSecretName) }}
  {{- end }}
{{- end }}
{{- $secretData | join "|" | sha256sum | trunc 8 }}
{{- end }}
