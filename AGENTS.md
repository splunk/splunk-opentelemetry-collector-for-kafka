# Agent instructions for splunk-opentelemetry-collector-for-kafka

This file gives AI agents (and humans) quick context for working in this repository.

## Project Overview
- **SOC4Kafka**: Splunk OpenTelemetry Collector for Kafka. It consumes Kafka messages and forwards them to Splunk via the HTTP Event Collector (HEC).
- **Replacement** for [kafka-connect-splunk](https://github.com/splunk/kafka-connect-splunk); built on the [Splunk OpenTelemetry Collector](https://github.com/signalfx/splunk-otel-collector).
- This repo does **not** contain the collector binary; it provides **configuration patterns**, a **Helm chart**, **documentation**, and **automation** (e.g. quickstart/Ansible).
- **Integration tests** in `tests/` are **written in Go** (functional and performance tests, plus shared helpers in `tests/common/`).

## Repository layout

| Path                                                             | Purpose                                                                       |
|------------------------------------------------------------------|-------------------------------------------------------------------------------|
| `helm-chart/splunk-opentelemetry-collector-for-kafka/`           | Helm chart (templates, values, schema, chart docs)                            |
| `helm-chart/splunk-opentelemetry-collector-for-kafka/templates/` | Kubernetes manifests (deployment, configmap, secret, service, etc.)           |
| `helm-chart/splunk-opentelemetry-collector-for-kafka/docs/`      | Chart-specific docs (installation, configuration, TLS, secrets, examples)     |
| `docs/`                                                          | Product/docs at repo root (migration, scaling, quickstart, otel design, etc.) |
| `quickstart/`                                                    | Ansible/automation for installing and configuring the collector               |
| `tests/`                                                         | Functional and performance tests (configs and test data under `testdata/`)    |
| `rendered/`                                                      | Output of `render_manifests.sh` (generated; used for tests/comparison)        |
| `.github/workflows/`                                             | CI (release, release-drafter, etc.)                                           |

## Key files

- **Root**: `README.md` – product overview, requirements, minimal config, run instructions.
- **Helm**: `helm-chart/splunk-opentelemetry-collector-for-kafka/values.yaml` (defaults), `values.schema.json` (validation), `Chart.yaml` (metadata, version).
- **Config generation**: `helm-chart/.../templates/_config.tpl` – OpenTelemetry collector config built from Helm values.
- **Deployment**: `helm-chart/.../templates/deployment.yaml` – main workload (image, env, volumes, probes).

## GitHub workflows

| Workflow | Trigger | Purpose |
|----------|---------|---------|
| `release-drafter.yaml` | Schedule, `workflow_dispatch` | Checks [splunk-otel-collector](https://github.com/signalfx/splunk-otel-collector) for new releases; when a newer version exists, updates version strings across the repo and opens a PR. Optional manual version override. |
| `release.yaml` | Push to `main` | If the commit message starts with `Prepare release v`, creates a GitHub release with that version and notes linking to the matching splunk-otel-collector release. |
| `functional_tests.yaml` | Push to `main`, pull requests | Runs functional and performance Go tests; uses `.github/actions/setup_env` to start Kafka and Splunk. |

## Documentation and technical writing

When editing or creating documentation, follow the guidelines in **[TECHNICAL_WRITER.md](TECHNICAL_WRITER.md)**. That file is written for any AI tool and defines the technical writer role, audience, tone, structure, and project-specific doc conventions.

## Conventions and practices

- **Secrets**: Never hardcode tokens, passwords, or keys. Use Kubernetes secrets, Vault, or env from a secret store; see `helm-chart/.../docs/secrets.md` and chart options for token/secret references.
- **Config**: Collector behavior is driven by YAML (receivers, processors, exporters, pipelines). Chart values map into that via `_config.tpl`; when adding features, update both `values.yaml` / `values.schema.json` and the template.
- **Versions**: Chart and app version live in `Chart.yaml`; keep docs and examples in sync with supported Kafka and Splunk versions (see root `README.md`). This project releases with the **same tag** as [splunk-otel-collector](https://github.com/signalfx/splunk-otel-collector) releases (e.g. when upstream releases `v0.147.1`, this repo releases the same version).
- **Compatibility**: Chart targets standard Kubernetes (no vendor CRDs); `kubeVersion` in `Chart.yaml` defines minimum cluster version.

## How to validate changes

- **Helm**: From repo root, template with a values file, e.g.  
  `helm template release-name helm-chart/splunk-opentelemetry-collector-for-kafka -f <values-file>`
- **Render script**: `./render_manifests.sh` uses `rendered/values_*.yaml` and writes to `rendered/manifests/tests_*/`; useful for diffing and test fixtures. Requires `helm` and uses `sed -i ''` (macOS-style).
- **Schema**: `values.schema.json` is used by Helm for `helm install/upgrade`; when adding or changing values, update the schema and keep defaults in `values.yaml` aligned.

## How to run integration tests

Integration tests are written in **Go**. Run them from the `tests/` directory. They require a running Kafka broker and Splunk instance (and related env vars); CI uses `.github/actions/setup_env` to bring these up.

### Local execution

1. **Set environment variables.** Copy `tests/local_execution/set_env.sh`, fill in the required values (Splunk host, password, HEC token, Kafka broker address), and source it before running tests. Set `CI_OTEL_BINARY_FILE` to the collector binary for your platform (e.g. `otelcol_darwin_arm64`, `otelcol_linux_amd64`). Optionally set `CI_SPLUNK_INDEX` (default in CI is `kafka`).
2. **Prerequisites.** Have Kafka and Splunk running and reachable at the host/ports you set. Ensure the collector binary named in `CI_OTEL_BINARY_FILE` exists on `PATH` or in the expected location.
3. **Run from `tests/`:**

- **Functional tests** (single topic, multiple topics, custom headers, timestamp extraction, regex topics):
  ```bash
  cd tests && go test ./functional_tests/functional_test.go -v -timeout 12m
  ```
- **Performance tests**:
  ```bash
  cd tests && go test ./performance_tests/performance_test.go -v -timeout 12m
  ```

See `.github/workflows/functional_tests.yaml` for the exact env (e.g. `CI_SPLUNK_*`, `CI_KAFKA_BROKER_ADDRESS`, `CI_OTEL_BINARY_FILE`) and Go version used in CI.

## Where to look for specific tasks

- **Add or change a Helm value**: `values.yaml`, `values.schema.json`, and `templates/_config.tpl` (and possibly `deployment.yaml` or `secret.yaml` if it affects env or mounts).
- **TLS / auth**: Chart docs under `helm-chart/.../docs/tls.md`, `docs/secrets.md`; templates for volumes and env that pass certs or flags.
- **Pipeline/receiver/exporter config**: Largely built in `_config.tpl` from structured values (e.g. `kafkaReceivers`, `splunkExporters`, `pipelines`).
- **User-facing behavior**: Root `docs/` and `helm-chart/.../README.md` plus `helm-chart/.../docs/`; keep examples and tables in README/docs in sync with the chart.

## External references

- [Splunk OTel Collector](https://github.com/signalfx/splunk-otel-collector) – base collector and releases.
- [OpenTelemetry Collector Kafka receiver](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/kafkareceiver) – receiver config.
- [Splunk HEC exporter](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/exporter/splunkhecexporter) – exporter config.

Use this file as the starting point for navigation and conventions when editing this repo.
