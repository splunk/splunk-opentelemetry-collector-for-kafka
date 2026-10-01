# Collector for Kafka documentation IA draft

Working outline and source-to-target map for the IA refactor. Content has been moved into the pages listed below; copyediting and technical review remain. The hierarchy below is the current page tree.

`[provisional]` means confirm the topic and scope with engineering before treating the page as committed structure.

- `docs/index.md` — landing page (logical topic: `collector-for-kafka-intro`)
  - `collector-for-kafka-design.md`
- `collector-for-kafka-deploy.md` — choose a deployment method
  - `collector-for-kafka-kubernetes.md`
    - `collector-for-kafka-kubernetes-install-helm.md`
    - `collector-for-kafka-kubernetes-configure-helm.md` — Helm chart values and Kubernetes-specific configuration
    - `collector-for-kafka-kubernetes-troubleshoot-helm.md`
    - `collector-for-kafka-kubernetes-upgrade-helm.md`
    - `collector-for-kafka-kubernetes-uninstall-helm.md`
  - `collector-for-kafka-ansible.md` `[provisional]`
    - `collector-for-kafka-ansible-install.md`
    - `collector-for-kafka-ansible-configure.md`
    - `collector-for-kafka-ansible-troubleshoot.md` `[provisional]`
    - `collector-for-kafka-ansible-upgrade.md` `[provisional]`
    - `collector-for-kafka-ansible-uninstall.md` `[provisional]`
  - `collector-for-kafka-manual.md`
    - `collector-for-kafka-manual-install.md`
  - `collector-for-kafka-oci-streaming.md` `[provisional; current source is an OCI Ubuntu scenario]`
- `collector-for-kafka-configure.md` — shared Collector receiver, processor, exporter, pipeline, and configuration guidance
  - `collector-for-kafka-configure-secrets.md`
  - `collector-for-kafka-configure-tls.md`
  - `collector-for-kafka-configure-examples.md`
  - `collector-for-kafka-configure-advanced.md`
    - `collector-for-kafka-configure-multiple-topics.md`
    - `collector-for-kafka-configure-regex-topics.md`
    - `collector-for-kafka-configure-extract-data.md`
- `collector-for-kafka-operate.md`
  - `collector-for-kafka-scale.md`
  - `collector-for-kafka-load-balance.md`
  - `collector-for-kafka-troubleshoot.md`
- `collector-for-kafka-monitor.md`
  - `collector-for-kafka-dashboard.md`
  - `collector-for-kafka-collector-logs.md` — distinguish Collector internal logs from Kafka event logs
- `collector-for-kafka-migrate-from-sc4kafka.md`
  - `collector-for-kafka-migrate-config.md`
  - `collector-for-kafka-migrate-examples.md`
  - `collector-for-kafka-migrate-strategy.md`

## Configuration boundary

- Shared Collector configuration describes receivers, processors, exporters, pipelines, and options that apply across deployment methods.
- Helm configuration describes chart values and Kubernetes-specific behavior. It should link to shared Collector configuration rather than duplicate it.

## Source-to-target map

Status records content transfer only. Local Markdown links resolve; copyediting and technical review remain.

## Markdown-to-DITA filename exception

The landing page is stored at `docs/index.md` so Zensical serves it at the site root. Its logical topic name remains `collector-for-kafka-intro`; use `collector-for-kafka-intro.dita` for the DITA topic filename when converting the Markdown set. Other pages can keep matching Markdown and DITA basenames.

| Source file or section | Target page(s) | Action | Status |
| --- | --- | --- | --- |
| Original `docs/index.md` — overview, requirements, platforms, features, migration and monitoring links | `docs/index.md` (landing page; logical topic `collector-for-kafka-intro`) | Merge and split | Moved; stored at root for Zensical; DITA filename exception recorded above |
| `docs/otel_design.md` | `collector-for-kafka-design.md` | Move | Moved |
| `docs/getting_started.md` — deployment choices, manual-install introduction, package download/run, minimal config and table | `collector-for-kafka-deploy.md`, `collector-for-kafka-manual-install.md`, `collector-for-kafka-configure.md` | Split | Moved |
| `docs/quickstart_guide.md` — Ansible quickstart and variables | `collector-for-kafka-ansible-install.md` | Move | Moved; other Ansible pages remain provisional |
| `docs/oci_installation.md` — OCI Streaming prerequisites and method chooser | `collector-for-kafka-oci-streaming.md` | Split | Moved |
| `docs/oci_installation.md` — systemd procedure | `collector-for-kafka-manual-install.md` | Split | Moved |
| `docs/oci_installation.md` — MicroK8s/Helm procedure | `collector-for-kafka-kubernetes-install-helm.md` | Split | Moved |
| `docs/helm/installation.md` — install, upgrade, uninstall | `collector-for-kafka-kubernetes-install-helm.md`, `collector-for-kafka-kubernetes-upgrade-helm.md`, `collector-for-kafka-kubernetes-uninstall-helm.md` | Split | Moved |
| `docs/helm/configuration.md` — chart configuration, precedence, restarts, metrics | `collector-for-kafka-kubernetes-configure-helm.md` | Move | Moved |
| `docs/helm/configuration.md` — Collector log settings | `collector-for-kafka-collector-logs.md` | Split | Moved |
| `docs/helm/examples.md` | `collector-for-kafka-configure-examples.md` | Move | Moved; examples remain together and Helm-specific |
| `docs/helm/secrets.md` | `collector-for-kafka-configure-secrets.md` | Move | Moved |
| `docs/helm/tls.md` | `collector-for-kafka-configure-tls.md` | Move | Moved; Kubernetes secret mounting remains in its TLS context |
| `docs/helm/troubleshooting.md` | `collector-for-kafka-kubernetes-troubleshoot-helm.md` | Move | Moved; common issues remain in the Helm troubleshooting page |
| `docs/multiple_topics.md` | `collector-for-kafka-configure-multiple-topics.md` | Move | Moved |
| `docs/regex_topics.md` | `collector-for-kafka-configure-regex-topics.md` | Move | Moved |
| `docs/extracting_additional_data.md` | `collector-for-kafka-configure-extract-data.md` | Move | Moved |
| `docs/scaling.md` | `collector-for-kafka-scale.md` | Move | Moved |
| `docs/loadbalancing.md` | `collector-for-kafka-load-balance.md` | Move | Moved |
| `docs/splunk-dashboard.md` | `collector-for-kafka-dashboard.md` | Move | Moved |
| `docs/collecting_own_logs.md` | `collector-for-kafka-collector-logs.md` | Merge | Moved with Helm log settings |
| `docs/migration.md` — overview/process, config map, examples, strategy | `collector-for-kafka-migrate-from-sc4kafka.md`, `collector-for-kafka-migrate-config.md`, `collector-for-kafka-migrate-examples.md`, `collector-for-kafka-migrate-strategy.md` | Split | Moved |
| `docs/migration_config_values.md` | `collector-for-kafka-migrate-config.md` | Merge | Moved with mapping section |
