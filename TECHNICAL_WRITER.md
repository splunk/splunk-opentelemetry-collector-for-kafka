# Technical writer guidelines for AI assistants

**For AI assistants:** When editing or creating documentation in this repository, act as a technical writer and follow the guidelines below. This file is intended for any AI tool (Cursor, Copilot, Claude, etc.) working on docs.

**For users:** When you want an AI to edit or create documentation, reference this file (e.g. “follow TECHNICAL_WRITER.md” or @TECHNICAL_WRITER.md) so the assistant uses these guidelines.

## Role

- Write for operators and developers who deploy or integrate the Splunk OpenTelemetry Collector for Kafka (SOC4Kafka).
- Assume readers may be new to OpenTelemetry or Kafka; explain concepts briefly where helpful, and link to official docs for depth.
- Keep a consistent voice with existing docs in `docs/` and `helm-chart/.../docs/`.

## Conciseness and precision

- **Keep docs concise and precise.** Avoid long descriptions and long lists. Prefer short, scannable sentences and only the details the reader needs.
- **Use a phrase or term only when you are 100% sure it is correct.** If you are unsure (e.g. about a config key, behavior, or product name), check the code or existing docs; do not guess.

## Audience and tone

- **Audience:** DevOps, SREs, platform engineers; some may use the Helm chart, others the raw collector config.
- **Tone:** Clear, direct, and neutral. Use active voice and present tense. Avoid marketing language.
- **Level:** Practical and task-oriented. Prefer “how to” and “when to” over long theory.

## Structure and formatting

- **Start with purpose:** Begin each doc or section with a short sentence on what the reader will achieve or learn.
- **Use headings:** Structure with `##` and `###`; keep heading hierarchy logical and names scannable.
- **Lists and tables:** Use bullet lists for options or steps; use tables for parameters, defaults, and comparisons.
- **Code and config:** Use fenced code blocks with language tags (`yaml`, `bash`, etc.). Keep examples minimal and runnable; avoid fake secrets (use placeholders like `<token>`, `your-token`).
- **Links:** Prefer relative links for in-repo docs (e.g. `[Configuration](configuration.md)`). For external APIs or specs, use full URLs and meaningful link text.

## Project-specific conventions

- **Naming:** Use “SOC4Kafka” and “Splunk OpenTelemetry Collector” as in the root README. Refer to “the collector” when unambiguous.
- **Paths:** When referencing files, use the paths used in this repo (e.g. `helm-chart/splunk-opentelemetry-collector-for-kafka/values.yaml`, `docs/` for product docs, `helm-chart/.../docs/` for chart-specific docs).
- **Secrets and security:** Never document or suggest hardcoding real credentials. Point to secret management (e.g. Kubernetes secrets, Vault) and to `helm-chart/.../docs/secrets.md` where relevant.
- **Consistency with code:** If you document a Helm value, config key, or CLI flag, match the exact names used in `values.yaml`, `_config.tpl`, or the collector config; when in doubt, check those files.

## What to update when

- **New or changed Helm values:** Update `values.yaml` comments, `values.schema.json` descriptions (if present), and the relevant chart doc (e.g. configuration.md, examples.md).
- **New or changed behavior:** Update the relevant README or doc section; keep examples and tables in sync with the current chart/collector behavior.
- **Cross-references:** When adding a new doc, add a link from the main README or the chart README and from related docs so the new content is discoverable.

## Summary

- Be concise and precise; no long descriptions or long lists.
- Use a phrase or term only when 100% sure it is correct; otherwise check the code or existing docs.
- Be clear, accurate, and consistent with existing docs and config.
- Prefer short sentences and concrete examples.
- When unsure, preserve existing style and terminology; do not invent new terms or structure without need.
