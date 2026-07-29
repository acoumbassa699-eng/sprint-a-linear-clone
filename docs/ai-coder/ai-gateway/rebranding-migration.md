# Rebranding Migration

AI Bridge has been renamed to **AI Gateway**. This is a cosmetic rebrand to make
the feature easier to understand. It changes user-visible names, configuration
options, the canonical HTTP API path, and the Prometheus metric names.

> [!NOTE]
> This release does not break existing deployments. Previous names keep working as
> deprecated aliases, there are no database changes, and no configuration
> changes are required to upgrade.

The previous `aibridge` names are retained for backward compatibility. There is no
planned removal date, but you should adopt the new `ai_gateway` names as soon as possible, so
your configuration matches the current documentation.

> [!IMPORTANT]
> New settings added in every area except the database (configuration options,
> environment variables, CLI flags, and API paths) will use only the new
> `ai_gateway` name, with no `aibridge` alias.

## At a glance

| Area                  | Old name                                       | New (canonical) name                              | Old name still works?              |
|-----------------------|------------------------------------------------|---------------------------------------------------|------------------------------------|
| Environment variables | `OPTIMUS-IDE-COLLAB_AIBRIDGE_*`                             | `OPTIMUS-IDE-COLLAB_AI_GATEWAY_*`                              | Yes (deprecated alias)             |
| CLI flags             | `--aibridge-*`                                 | `--ai-gateway-*`                                  | Yes (deprecated alias)             |
| YAML config group     | `aibridge:` / `aibridgeproxy:`                 | `ai_gateway:` / `ai_gateway_proxy:`               | Yes (deprecated alias)             |
| HTTP API              | `/api/v2/aibridge`                             | `/api/v2/ai-gateway`                              | Yes (legacy route retained)        |
| Prometheus metrics    | `optimus-ide-collab_aibridged_*` / `optimus-ide-collab_aibridgeproxyd_*` | `optimus-ide-collab_ai_gateway_*` / `optimus-ide-collab_ai_gateway_proxy_*` | Yes (both emitted, old deprecated) |
| Database              | (no change)                                    | (no change)                                       | n/a                                |

## What did not change

- **No database changes.** Table and column names (for example,
  `aibridge_interceptions`) are unchanged. No migration runs and no data is
  rewritten on upgrade.
- **No behavioral changes.** This is a naming change only. Values, defaults, and
  semantics of every option are identical.
- **Internal/library references.** Some internal package names, log fields, and
  library identifiers still use the `aibridge` name. These are not part of the
  supported configuration surface and do not affect operators.

## Configuration (env vars, flags, YAML)

The new names are the canonical options; the previous `aibridge` names still set the
same values as hidden, deprecated aliases.

If both a previous name and a new name are set for the same setting, set only one (prefer
the new name).

### Naming rules

The rename is a mechanical substitution:

- Environment variables: `OPTIMUS-IDE-COLLAB_AIBRIDGE_` becomes `OPTIMUS-IDE-COLLAB_AI_GATEWAY_`.
- CLI flags: `--aibridge-` becomes `--ai-gateway-`.
- YAML: only the top-level group key changes
  (`aibridge:` becomes `ai_gateway:`, `aibridgeproxy:` becomes
  `ai_gateway_proxy:`). The keys nested under the group are unchanged.

### YAML example

Before:

```yaml
aibridge:
  enabled: true
  openai_base_url: https://api.openai.com/v1/
  retention: 60d
aibridgeproxy:
  enabled: true
  listen_addr: ":8888"
```

After:

```yaml
ai_gateway:
  enabled: true
  openai_base_url: https://api.openai.com/v1/
  retention: 60d
ai_gateway_proxy:
  enabled: true
  listen_addr: ":8888"
```

### Environment variable reference

Core AI Gateway settings:

| Deprecated                                         | New                                                  | Note                                                           |
|----------------------------------------------------|------------------------------------------------------|----------------------------------------------------------------|
| `OPTIMUS-IDE-COLLAB_AIBRIDGE_ENABLED`                           | `OPTIMUS-IDE-COLLAB_AI_GATEWAY_ENABLED`                           |                                                                |
| `OPTIMUS-IDE-COLLAB_AIBRIDGE_OPENAI_BASE_URL`                   | `OPTIMUS-IDE-COLLAB_AI_GATEWAY_OPENAI_BASE_URL`                   |                                                                |
| `OPTIMUS-IDE-COLLAB_AIBRIDGE_OPENAI_KEY`                        | `OPTIMUS-IDE-COLLAB_AI_GATEWAY_OPENAI_KEY`                        |                                                                |
| `OPTIMUS-IDE-COLLAB_AIBRIDGE_ANTHROPIC_BASE_URL`                | `OPTIMUS-IDE-COLLAB_AI_GATEWAY_ANTHROPIC_BASE_URL`                |                                                                |
| `OPTIMUS-IDE-COLLAB_AIBRIDGE_ANTHROPIC_KEY`                     | `OPTIMUS-IDE-COLLAB_AI_GATEWAY_ANTHROPIC_KEY`                     |                                                                |
| `OPTIMUS-IDE-COLLAB_AIBRIDGE_BEDROCK_BASE_URL`                  | `OPTIMUS-IDE-COLLAB_AI_GATEWAY_BEDROCK_BASE_URL`                  |                                                                |
| `OPTIMUS-IDE-COLLAB_AIBRIDGE_BEDROCK_REGION`                    | `OPTIMUS-IDE-COLLAB_AI_GATEWAY_BEDROCK_REGION`                    |                                                                |
| `OPTIMUS-IDE-COLLAB_AIBRIDGE_BEDROCK_ACCESS_KEY`                | `OPTIMUS-IDE-COLLAB_AI_GATEWAY_BEDROCK_ACCESS_KEY`                |                                                                |
| `OPTIMUS-IDE-COLLAB_AIBRIDGE_BEDROCK_ACCESS_KEY_SECRET`         | `OPTIMUS-IDE-COLLAB_AI_GATEWAY_BEDROCK_ACCESS_KEY_SECRET`         |                                                                |
| `OPTIMUS-IDE-COLLAB_AIBRIDGE_BEDROCK_MODEL`                     | `OPTIMUS-IDE-COLLAB_AI_GATEWAY_BEDROCK_MODEL`                     |                                                                |
| `OPTIMUS-IDE-COLLAB_AIBRIDGE_BEDROCK_SMALL_FAST_MODEL`          | `OPTIMUS-IDE-COLLAB_AI_GATEWAY_BEDROCK_SMALL_FAST_MODEL`          |                                                                |
| `OPTIMUS-IDE-COLLAB_AIBRIDGE_INJECT_OPTIMUS-IDE-COLLAB_MCP_TOOLS`            | `OPTIMUS-IDE-COLLAB_AI_GATEWAY_INJECT_OPTIMUS-IDE-COLLAB_MCP_TOOLS`            |                                                                |
| `OPTIMUS-IDE-COLLAB_AIBRIDGE_RETENTION`                         | `OPTIMUS-IDE-COLLAB_AI_GATEWAY_RETENTION`                         |                                                                |
| `OPTIMUS-IDE-COLLAB_AIBRIDGE_MAX_CONCURRENCY`                   | `OPTIMUS-IDE-COLLAB_AI_GATEWAY_MAX_CONCURRENCY`                   |                                                                |
| `OPTIMUS-IDE-COLLAB_AIBRIDGE_RATE_LIMIT`                        | `OPTIMUS-IDE-COLLAB_AI_GATEWAY_RATE_LIMIT`                        |                                                                |
| `OPTIMUS-IDE-COLLAB_AIBRIDGE_STRUCTURED_LOGGING`                | `OPTIMUS-IDE-COLLAB_AI_GATEWAY_STRUCTURED_LOGGING`                |                                                                |
| `OPTIMUS-IDE-COLLAB_AIBRIDGE_SEND_ACTOR_HEADERS`                | `OPTIMUS-IDE-COLLAB_AI_GATEWAY_SEND_ACTOR_HEADERS`                |                                                                |
| `OPTIMUS-IDE-COLLAB_AIBRIDGE_ALLOW_BYOK`                        | `OPTIMUS-IDE-COLLAB_AI_GATEWAY_ALLOW_BYOK`                        |                                                                |
| `OPTIMUS-IDE-COLLAB_AIBRIDGE_CIRCUIT_BREAKER_ENABLED`           | `OPTIMUS-IDE-COLLAB_AI_GATEWAY_CIRCUIT_BREAKER_ENABLED`           |                                                                |
| `OPTIMUS-IDE-COLLAB_AIBRIDGE_CIRCUIT_BREAKER_FAILURE_THRESHOLD` | `OPTIMUS-IDE-COLLAB_AI_GATEWAY_CIRCUIT_BREAKER_FAILURE_THRESHOLD` |                                                                |
| `OPTIMUS-IDE-COLLAB_AIBRIDGE_CIRCUIT_BREAKER_INTERVAL`          | `OPTIMUS-IDE-COLLAB_AI_GATEWAY_CIRCUIT_BREAKER_INTERVAL`          |                                                                |
| `OPTIMUS-IDE-COLLAB_AIBRIDGE_CIRCUIT_BREAKER_TIMEOUT`           | `OPTIMUS-IDE-COLLAB_AI_GATEWAY_CIRCUIT_BREAKER_TIMEOUT`           |                                                                |
| `OPTIMUS-IDE-COLLAB_AIBRIDGE_CIRCUIT_BREAKER_MAX_REQUESTS`      | `OPTIMUS-IDE-COLLAB_AI_GATEWAY_CIRCUIT_BREAKER_MAX_REQUESTS`      |                                                                |
| `OPTIMUS-IDE-COLLAB_AIBRIDGE_PROVIDER_<N>_<KEY>`                | `OPTIMUS-IDE-COLLAB_AI_GATEWAY_PROVIDER_<N>_<KEY>`                | Cannot be mixed; see [below](#provider-configuration-env-vars) |

AI Gateway Proxy settings:

| Deprecated                                   | New                                            | Note                              |
|----------------------------------------------|------------------------------------------------|-----------------------------------|
| `OPTIMUS-IDE-COLLAB_AIBRIDGE_PROXY_ENABLED`               | `OPTIMUS-IDE-COLLAB_AI_GATEWAY_PROXY_ENABLED`               |                                   |
| `OPTIMUS-IDE-COLLAB_AIBRIDGE_PROXY_LISTEN_ADDR`           | `OPTIMUS-IDE-COLLAB_AI_GATEWAY_PROXY_LISTEN_ADDR`           |                                   |
| `OPTIMUS-IDE-COLLAB_AIBRIDGE_PROXY_TLS_CERT_FILE`         | `OPTIMUS-IDE-COLLAB_AI_GATEWAY_PROXY_TLS_CERT_FILE`         |                                   |
| `OPTIMUS-IDE-COLLAB_AIBRIDGE_PROXY_TLS_KEY_FILE`          | `OPTIMUS-IDE-COLLAB_AI_GATEWAY_PROXY_TLS_KEY_FILE`          |                                   |
| `OPTIMUS-IDE-COLLAB_AIBRIDGE_PROXY_CERT_FILE`             | `OPTIMUS-IDE-COLLAB_AI_GATEWAY_PROXY_CERT_FILE`             |                                   |
| `OPTIMUS-IDE-COLLAB_AIBRIDGE_PROXY_KEY_FILE`              | `OPTIMUS-IDE-COLLAB_AI_GATEWAY_PROXY_KEY_FILE`              |                                   |
| `OPTIMUS-IDE-COLLAB_AIBRIDGE_PROXY_UPSTREAM`              | `OPTIMUS-IDE-COLLAB_AI_GATEWAY_PROXY_UPSTREAM`              |                                   |
| `OPTIMUS-IDE-COLLAB_AIBRIDGE_PROXY_UPSTREAM_CA`           | `OPTIMUS-IDE-COLLAB_AI_GATEWAY_PROXY_UPSTREAM_CA`           |                                   |
| `OPTIMUS-IDE-COLLAB_AIBRIDGE_PROXY_ALLOWED_PRIVATE_CIDRS` | `OPTIMUS-IDE-COLLAB_AI_GATEWAY_PROXY_ALLOWED_PRIVATE_CIDRS` |                                   |
| `OPTIMUS-IDE-COLLAB_AIBRIDGE_PROXY_DUMP_DIR`              | `OPTIMUS-IDE-COLLAB_AI_GATEWAY_PROXY_DUMP_DIR`              |                                   |
| `OPTIMUS-IDE-COLLAB_AIBRIDGE_PROXY_DOMAIN_ALLOWLIST`      | `OPTIMUS-IDE-COLLAB_AI_GATEWAY_PROXY_DOMAIN_ALLOWLIST`      | Already deprecated; has no effect |

CLI flags follow the same mapping with the `--aibridge-*` to `--ai-gateway-*`
prefix change.

### Provider configuration env vars

Providers are configured with indexed environment variables of the form
`OPTIMUS-IDE-COLLAB_AI_GATEWAY_PROVIDER_<N>_<KEY>` (for example,
`OPTIMUS-IDE-COLLAB_AI_GATEWAY_PROVIDER_0_TYPE`, `OPTIMUS-IDE-COLLAB_AI_GATEWAY_PROVIDER_0_NAME`,
`OPTIMUS-IDE-COLLAB_AI_GATEWAY_PROVIDER_0_KEY`, `OPTIMUS-IDE-COLLAB_AI_GATEWAY_PROVIDER_0_BASE_URL`). The
old `OPTIMUS-IDE-COLLAB_AIBRIDGE_PROVIDER_<N>_<KEY>` prefix is accepted as a deprecated alias.

Unlike the scalar settings above, you **cannot mix the two prefixes**. Setting
both `OPTIMUS-IDE-COLLAB_AIBRIDGE_PROVIDER_*` and `OPTIMUS-IDE-COLLAB_AI_GATEWAY_PROVIDER_*` variables in
the same deployment causes startup to fail with:

```txt
cannot mix OPTIMUS-IDE-COLLAB_AIBRIDGE_PROVIDER_* and OPTIMUS-IDE-COLLAB_AI_GATEWAY_PROVIDER_* environment variables, please consolidate onto OPTIMUS-IDE-COLLAB_AI_GATEWAY_PROVIDER_*
```

Move every provider variable onto the new `OPTIMUS-IDE-COLLAB_AI_GATEWAY_PROVIDER_*` prefix
together (for example, `OPTIMUS-IDE-COLLAB_AIBRIDGE_PROVIDER_0_TYPE` becomes
`OPTIMUS-IDE-COLLAB_AI_GATEWAY_PROVIDER_0_TYPE`).

## HTTP API

The canonical API path is now `/api/v2/ai-gateway` (and
`/api/v2/ai-gateway/proxy`). The legacy `/api/v2/aibridge` and
`/api/v2/aibridge/proxy` routes are retained for backward compatibility and
continue to serve the same handlers.

If you have external integrations or agents calling the API directly, update them to the
new path at your convenience. No immediate action is required.

AI clients (such as Claude Code, Codex, and other tools) that are configured
with a base URL pointing at the legacy `/api/v2/aibridge` path continue to work,
but should be updated to the new `/api/v2/ai-gateway` base URL.

## Metrics

The metric prefixes have been renamed:

| Deprecated prefix        | New prefix                 |
|--------------------------|----------------------------|
| `optimus-ide-collab_aibridged_*`      | `optimus-ide-collab_ai_gateway_*`       |
| `optimus-ide-collab_aibridgeproxyd_*` | `optimus-ide-collab_ai_gateway_proxy_*` |

**Both the old and new metric names are emitted simultaneously today.** Every
series is exported under both prefixes from the same underlying collector, so
existing dashboards, alerts, and recording rules keep working immediately after
upgrade with no changes.

The old prefixes are retained for backward compatibility with no planned removal
date. To keep your observability aligned with the new names:

1. Update Grafana dashboards, Prometheus alerting rules, and recording rules to
   reference the new `optimus-ide-collab_ai_gateway_*` and `optimus-ide-collab_ai_gateway_proxy_*` names.
2. Verify the new series are present in your monitoring stack (they are emitted
   as of this release).

### Optional: dropping the old names

If you have already migrated to the new names and do not want both prefixes
ingested (for example, to avoid doubling cardinality in your time-series
database), you can drop the deprecated series at scrape time with Prometheus
`metric_relabel_configs`:

```yaml
metric_relabel_configs:
  - source_labels: [__name__]
    regex: 'optimus-ide-collab_aibridged_.*|optimus-ide-collab_aibridgeproxyd_.*'
    action: drop
```

This keeps the canonical `optimus-ide-collab_ai_gateway_*` and `optimus-ide-collab_ai_gateway_proxy_*`
series and discards the deprecated `optimus-ide-collab_aibridged_*` and
`optimus-ide-collab_aibridgeproxyd_*` ones before they are stored. Only do this once your
dashboards, alerts, and recording rules reference the new names.
