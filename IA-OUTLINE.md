# Splunk Distribution of OpenTelemetry Collector for Kafka documentation information architecture draft

This outline records the source-to-target map for the information architecture refactor. Content has been moved into the pages listed below. The hierarchy reflects the current page tree.

`[provisional]` means confirm the topic and scope with engineering before treating the page as committed structure.

- `docs/index.md` — landing page (logical topic: `collector-for-kafka-intro`)
  - `collector-for-kafka-design.md`
- `deploy/collector-for-kafka-deploy.md` — choose a deployment method
  - `deploy/collector-for-kafka-kubernetes.md`
    - `deploy/collector-for-kafka-kubernetes-install-helm.md`
    - `deploy/collector-for-kafka-kubernetes-configure-helm.md` — Helm chart values and Kubernetes-specific configuration
    - `deploy/collector-for-kafka-kubernetes-troubleshoot-helm.md`
    - `deploy/collector-for-kafka-kubernetes-upgrade-helm.md`
    - `deploy/collector-for-kafka-kubernetes-uninstall-helm.md`
  - `deploy/collector-for-kafka-ansible.md` `[provisional]`
    - `deploy/collector-for-kafka-ansible-install.md`
    - Potential topics to confirm with engineering: configuration, troubleshooting, upgrading, and uninstalling. Create pages only if those topics need distinct Ansible-specific content.
  - `deploy/collector-for-kafka-manual.md`
    - `deploy/collector-for-kafka-manual-install.md`
  - `deploy/collector-for-kafka-oci-streaming.md` `[provisional; includes OCI Ubuntu systemd and MicroK8s/Helm procedures]`
- `configure/collector-for-kafka-configure.md` — shared Collector receiver, processor, exporter, pipeline, and configuration guidance
  - `configure/collector-for-kafka-configure-secrets.md`
  - `configure/collector-for-kafka-configure-tls.md`
  - `configure/collector-for-kafka-configure-examples.md`
  - `configure/collector-for-kafka-configure-advanced.md`
    - `configure/collector-for-kafka-configure-multiple-topics.md`
    - `configure/collector-for-kafka-configure-regex-topics.md`
    - `configure/collector-for-kafka-configure-extract-data.md`
- `operate/collector-for-kafka-operate.md`
  - `operate/collector-for-kafka-scale.md`
  - `operate/collector-for-kafka-load-balance.md`
- `monitor/collector-for-kafka-monitor.md`
  - `monitor/collector-for-kafka-dashboard.md`
  - `monitor/collector-for-kafka-collector-logs.md` — distinguish Collector internal logs from Kafka event logs
- `migrate/collector-for-kafka-migrate-from-sc4kafka.md`
  - `migrate/collector-for-kafka-migrate-config.md`
  - `migrate/collector-for-kafka-migrate-examples.md`
  - `migrate/collector-for-kafka-migrate-strategy.md`

## Configuration boundary

- Shared Collector configuration describes receivers, processors, exporters, pipelines, and options that apply across deployment methods.
- Helm configuration describes chart values and Kubernetes-specific behavior. It should link to shared Collector configuration rather than duplicate it.

## Source-to-target map

Status records content transfer only. Copyediting is complete; technical review remains.

## Markdown-to-DITA filename exception

The landing page is stored at `docs/index.md` so Zensical serves it at the site root. Its logical topic name remains `collector-for-kafka-intro`; use `collector-for-kafka-intro.dita` for the DITA topic filename when converting the Markdown set. Other pages can keep matching Markdown and DITA basenames.

| Source file or section | Target page(s) | Action | Status |
| --- | --- | --- | --- |
| Original `docs/index.md` — overview, requirements, platforms, features, migration and monitoring links | `docs/index.md` (landing page; logical topic `collector-for-kafka-intro`) | Merge and split | Moved; stored at root for Zensical; DITA filename exception recorded above |
| `docs/otel_design.md` | `collector-for-kafka-design.md` | Move | Moved |
| `docs/getting_started.md` — deployment choices, manual-install introduction, package download/run, minimal config and table | `deploy/collector-for-kafka-deploy.md`, `deploy/collector-for-kafka-manual-install.md`, `configure/collector-for-kafka-configure.md` | Split | Moved |
| `docs/quickstart_guide.md` — Ansible quickstart and variables | `deploy/collector-for-kafka-ansible-install.md` | Move | Moved; additional Ansible topics remain provisional and have no placeholder pages |
| `docs/oci_installation.md` — OCI Streaming prerequisites and method chooser | `deploy/collector-for-kafka-oci-streaming.md` | Split | Moved |
| `docs/oci_installation.md` — systemd procedure | `deploy/collector-for-kafka-oci-streaming.md` | Split | Moved; OCI-specific procedure grouped with the OCI Streaming deployment page |
| `docs/oci_installation.md` — MicroK8s/Helm procedure | `deploy/collector-for-kafka-oci-streaming.md` | Split | Moved; OCI-specific procedure grouped with the OCI Streaming deployment page |
| `docs/helm/installation.md` — install, upgrade, uninstall | `deploy/collector-for-kafka-kubernetes-install-helm.md`, `deploy/collector-for-kafka-kubernetes-upgrade-helm.md`, `deploy/collector-for-kafka-kubernetes-uninstall-helm.md` | Split | Moved |
| `docs/helm/configuration.md` — chart configuration, precedence, restarts, metrics | `deploy/collector-for-kafka-kubernetes-configure-helm.md` | Move | Moved |
| `docs/helm/configuration.md` — Collector log settings | `monitor/collector-for-kafka-collector-logs.md` | Split | Moved |
| `docs/helm/examples.md` | `configure/collector-for-kafka-configure-examples.md` | Move | Moved; examples remain together and Helm-specific |
| `docs/helm/secrets.md` | `configure/collector-for-kafka-configure-secrets.md` | Move | Moved |
| `docs/helm/tls.md` | `configure/collector-for-kafka-configure-tls.md` | Move | Moved; Kubernetes secret mounting remains in its TLS context |
| `docs/helm/troubleshooting.md` | `deploy/collector-for-kafka-kubernetes-troubleshoot-helm.md` | Move | Moved; common issues remain in the Helm troubleshooting page |
| `docs/multiple_topics.md` | `configure/collector-for-kafka-configure-multiple-topics.md` | Move | Moved |
| `docs/regex_topics.md` | `configure/collector-for-kafka-configure-regex-topics.md` | Move | Moved |
| `docs/extracting_additional_data.md` | `configure/collector-for-kafka-configure-extract-data.md` | Move | Moved |
| `docs/scaling.md` | `operate/collector-for-kafka-scale.md` | Move | Moved |
| `docs/loadbalancing.md` | `operate/collector-for-kafka-load-balance.md` | Move | Moved |
| `docs/splunk-dashboard.md` | `monitor/collector-for-kafka-dashboard.md` | Move | Moved |
| `docs/collecting_own_logs.md` | `monitor/collector-for-kafka-collector-logs.md` | Merge | Moved with Helm log settings |
| `docs/migration.md` — overview/process, config map, examples, strategy | `migrate/collector-for-kafka-migrate-from-sc4kafka.md`, `migrate/collector-for-kafka-migrate-config.md`, `migrate/collector-for-kafka-migrate-examples.md`, `migrate/collector-for-kafka-migrate-strategy.md` | Split | Moved |
| `docs/migration_config_values.md` | `migrate/collector-for-kafka-migrate-config.md` | Merge | Moved with mapping section |
