# Schemas

## agentsdk.AWSInstanceIdentityToken

```json
{
  "agent_name": "string",
  "document": "string",
  "signature": "string"
}
```

### Properties

| Name         | Type   | Required | Restrictions | Description                                                                                                                                      |
|--------------|--------|----------|--------------|--------------------------------------------------------------------------------------------------------------------------------------------------|
| `agent_name` | string | false    |              | Agent name optionally selects a specific agent when multiple agents share the same instance identity. An empty string is treated as unspecified. |
| `document`   | string | true     |              |                                                                                                                                                  |
| `signature`  | string | true     |              |                                                                                                                                                  |

## agentsdk.AuthenticateResponse

```json
{
  "session_token": "string"
}
```

### Properties

| Name            | Type   | Required | Restrictions | Description |
|-----------------|--------|----------|--------------|-------------|
| `session_token` | string | false    |              |             |

## agentsdk.AzureInstanceIdentityToken

```json
{
  "agent_name": "string",
  "encoding": "string",
  "signature": "string"
}
```

### Properties

| Name         | Type   | Required | Restrictions | Description                                                                                                                                      |
|--------------|--------|----------|--------------|--------------------------------------------------------------------------------------------------------------------------------------------------|
| `agent_name` | string | false    |              | Agent name optionally selects a specific agent when multiple agents share the same instance identity. An empty string is treated as unspecified. |
| `encoding`   | string | true     |              |                                                                                                                                                  |
| `signature`  | string | true     |              |                                                                                                                                                  |

## agentsdk.ExternalAuthResponse

```json
{
  "access_token": "string",
  "expires_at": "string",
  "password": "string",
  "token_extra": {},
  "type": "string",
  "url": "string",
  "username": "string"
}
```

### Properties

| Name           | Type   | Required | Restrictions | Description                                                                                                                    |
|----------------|--------|----------|--------------|--------------------------------------------------------------------------------------------------------------------------------|
| `access_token` | string | false    |              |                                                                                                                                |
| `expires_at`   | string | false    |              | Expires at is the time the token expires, normalized to UTC (for example, "2024-06-01T15:04:05Z"). Zero value means no expiry. |
| `password`     | string | false    |              |                                                                                                                                |
| `token_extra`  | object | false    |              |                                                                                                                                |
| `type`         | string | false    |              |                                                                                                                                |
| `url`          | string | false    |              |                                                                                                                                |
| `username`     | string | false    |              | Deprecated: Only supported on `/workspaceagents/me/gitauth` for backwards compatibility.                                       |

## agentsdk.GitSSHKey

```json
{
  "private_key": "string",
  "public_key": "string"
}
```

### Properties

| Name          | Type   | Required | Restrictions | Description |
|---------------|--------|----------|--------------|-------------|
| `private_key` | string | false    |              |             |
| `public_key`  | string | false    |              |             |

## agentsdk.GoogleInstanceIdentityToken

```json
{
  "agent_name": "string",
  "json_web_token": "string"
}
```

### Properties

| Name             | Type   | Required | Restrictions | Description                                                                                                                                      |
|------------------|--------|----------|--------------|--------------------------------------------------------------------------------------------------------------------------------------------------|
| `agent_name`     | string | false    |              | Agent name optionally selects a specific agent when multiple agents share the same instance identity. An empty string is treated as unspecified. |
| `json_web_token` | string | true     |              |                                                                                                                                                  |

## agentsdk.Log

```json
{
  "created_at": "string",
  "level": "trace",
  "output": "string"
}
```

### Properties

| Name         | Type                                   | Required | Restrictions | Description |
|--------------|----------------------------------------|----------|--------------|-------------|
| `created_at` | string                                 | false    |              |             |
| `level`      | [optimus-ide-collabsdk.LogLevel](#optimus-ide-collabsdkloglevel) | false    |              |             |
| `output`     | string                                 | false    |              |             |

## agentsdk.PatchAppStatus

```json
{
  "app_slug": "string",
  "icon": "string",
  "message": "string",
  "needs_user_attention": true,
  "state": "working",
  "uri": "string"
}
```

### Properties

| Name                   | Type                                                                 | Required | Restrictions | Description                                                               |
|------------------------|----------------------------------------------------------------------|----------|--------------|---------------------------------------------------------------------------|
| `app_slug`             | string                                                               | false    |              |                                                                           |
| `icon`                 | string                                                               | false    |              | Deprecated: this field is unused and will be removed in a future version. |
| `message`              | string                                                               | false    |              |                                                                           |
| `needs_user_attention` | boolean                                                              | false    |              | Deprecated: this field is unused and will be removed in a future version. |
| `state`                | [optimus-ide-collabsdk.WorkspaceAppStatusState](#optimus-ide-collabsdkworkspaceappstatusstate) | false    |              |                                                                           |
| `uri`                  | string                                                               | false    |              |                                                                           |

## agentsdk.PatchLogs

```json
{
  "log_source_id": "string",
  "logs": [
    {
      "created_at": "string",
      "level": "trace",
      "output": "string"
    }
  ]
}
```

### Properties

| Name            | Type                                  | Required | Restrictions | Description |
|-----------------|---------------------------------------|----------|--------------|-------------|
| `log_source_id` | string                                | false    |              |             |
| `logs`          | array of [agentsdk.Log](#agentsdklog) | false    |              |             |

## agentsdk.PostLogSourceRequest

```json
{
  "display_name": "string",
  "icon": "string",
  "id": "string"
}
```

### Properties

| Name           | Type   | Required | Restrictions | Description                                                                                                                                                                                    |
|----------------|--------|----------|--------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `display_name` | string | false    |              |                                                                                                                                                                                                |
| `icon`         | string | false    |              |                                                                                                                                                                                                |
| `id`           | string | false    |              | ID is a unique identifier for the log source. It is scoped to a workspace agent, and can be statically defined inside code to prevent duplicate sources from being created for the same agent. |

## agentsdk.ReinitializationEvent

```json
{
  "owner_id": "8826ee2e-7933-4665-aef2-2393f84a0d05",
  "reason": "prebuild_claimed",
  "workspace_id": "0967198e-ec7b-4c6b-b4d3-f71244cadbe9"
}
```

### Properties

| Name           | Type                                                               | Required | Restrictions | Description |
|----------------|--------------------------------------------------------------------|----------|--------------|-------------|
| `owner_id`     | string                                                             | false    |              |             |
| `reason`       | [agentsdk.ReinitializationReason](#agentsdkreinitializationreason) | false    |              |             |
| `workspace_id` | string                                                             | false    |              |             |

## agentsdk.ReinitializationReason

```json
"prebuild_claimed"
```

### Properties

#### Enumerated Values

| Value(s)           |
|--------------------|
| `prebuild_claimed` |

## optimus-ide-collabd.cspViolation

```json
{
  "csp-report": {}
}
```

### Properties

| Name         | Type   | Required | Restrictions | Description |
|--------------|--------|----------|--------------|-------------|
| `csp-report` | object | false    |              |             |

## optimus-ide-collabsdk.ACLAvailable

```json
{
  "groups": [
    {
      "avatar_url": "http://example.com",
      "display_name": "string",
      "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
      "members": [
        {
          "avatar_url": "http://example.com",
          "created_at": "2019-08-24T14:15:22Z",
          "email": "user@example.com",
          "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
          "is_service_account": true,
          "last_seen_at": "2019-08-24T14:15:22Z",
          "login_type": "",
          "name": "string",
          "status": "active",
          "theme_preference": "string",
          "updated_at": "2019-08-24T14:15:22Z",
          "username": "string"
        }
      ],
      "name": "string",
      "organization_display_name": "string",
      "organization_id": "7c60d51f-b44e-4682-87d6-449835ea4de6",
      "organization_name": "string",
      "quota_allowance": 0,
      "source": "user",
      "total_member_count": 0
    }
  ],
  "users": [
    {
      "avatar_url": "http://example.com",
      "created_at": "2019-08-24T14:15:22Z",
      "email": "user@example.com",
      "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
      "is_service_account": true,
      "last_seen_at": "2019-08-24T14:15:22Z",
      "login_type": "",
      "name": "string",
      "status": "active",
      "theme_preference": "string",
      "updated_at": "2019-08-24T14:15:22Z",
      "username": "string"
    }
  ]
}
```

### Properties

| Name     | Type                                                  | Required | Restrictions | Description |
|----------|-------------------------------------------------------|----------|--------------|-------------|
| `groups` | array of [optimus-ide-collabsdk.Group](#optimus-ide-collabsdkgroup)             | false    |              |             |
| `users`  | array of [optimus-ide-collabsdk.ReducedUser](#optimus-ide-collabsdkreduceduser) | false    |              |             |

## optimus-ide-collabsdk.AIBridgeAgenticAction

```json
{
  "model": "string",
  "thinking": [
    {
      "text": "string"
    }
  ],
  "token_usage": {
    "cache_read_input_tokens": 0,
    "cache_write_input_tokens": 0,
    "input_tokens": 0,
    "metadata": {
      "property1": null,
      "property2": null
    },
    "output_tokens": 0
  },
  "tool_calls": [
    {
      "created_at": "2019-08-24T14:15:22Z",
      "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
      "injected": true,
      "input": "string",
      "interception_id": "34d9b688-63ad-46f4-88b5-665c1e7f7824",
      "metadata": {
        "property1": null,
        "property2": null
      },
      "provider_response_id": "string",
      "server_url": "string",
      "tool": "string"
    }
  ]
}
```

### Properties

| Name          | Type                                                                                   | Required | Restrictions | Description |
|---------------|----------------------------------------------------------------------------------------|----------|--------------|-------------|
| `model`       | string                                                                                 | false    |              |             |
| `thinking`    | array of [optimus-ide-collabsdk.AIBridgeModelThought](#optimus-ide-collabsdkaibridgemodelthought)                | false    |              |             |
| `token_usage` | [optimus-ide-collabsdk.AIBridgeSessionThreadsTokenUsage](#optimus-ide-collabsdkaibridgesessionthreadstokenusage) | false    |              |             |
| `tool_calls`  | array of [optimus-ide-collabsdk.AIBridgeToolCall](#optimus-ide-collabsdkaibridgetoolcall)                        | false    |              |             |

## optimus-ide-collabsdk.AIBridgeAnthropicConfig

```json
{
  "base_url": "string",
  "key": "string"
}
```

### Properties

| Name       | Type   | Required | Restrictions | Description |
|------------|--------|----------|--------------|-------------|
| `base_url` | string | false    |              |             |
| `key`      | string | false    |              |             |

## optimus-ide-collabsdk.AIBridgeBedrockConfig

```json
{
  "access_key": "string",
  "access_key_secret": "string",
  "base_url": "string",
  "model": "string",
  "region": "string",
  "small_fast_model": "string"
}
```

### Properties

| Name                | Type   | Required | Restrictions | Description |
|---------------------|--------|----------|--------------|-------------|
| `access_key`        | string | false    |              |             |
| `access_key_secret` | string | false    |              |             |
| `base_url`          | string | false    |              |             |
| `model`             | string | false    |              |             |
| `region`            | string | false    |              |             |
| `small_fast_model`  | string | false    |              |             |

## optimus-ide-collabsdk.AIBridgeConfig

```json
{
  "allow_byok": true,
  "anthropic": {
    "base_url": "string",
    "key": "string"
  },
  "api_dump_dir": "string",
  "bedrock": {
    "access_key": "string",
    "access_key_secret": "string",
    "base_url": "string",
    "model": "string",
    "region": "string",
    "small_fast_model": "string"
  },
  "budget_period": "string",
  "budget_policy": "string",
  "circuit_breaker_enabled": true,
  "circuit_breaker_failure_threshold": 0,
  "circuit_breaker_interval": 0,
  "circuit_breaker_max_requests": 0,
  "circuit_breaker_timeout": 0,
  "enabled": true,
  "inject_optimus-ide-collab_mcp_tools": true,
  "max_concurrency": 0,
  "openai": {
    "base_url": "string",
    "key": "string"
  },
  "providers": [
    {
      "base_url": "string",
      "bedrock_model": "string",
      "bedrock_region": "string",
      "bedrock_small_fast_model": "string",
      "name": "string",
      "type": "string"
    }
  ],
  "rate_limit": 0,
  "retention": 0,
  "send_actor_headers": true,
  "structured_logging": true
}
```

### Properties

| Name                                | Type                                                                 | Required | Restrictions | Description                                                                                                                                                                     |
|-------------------------------------|----------------------------------------------------------------------|----------|--------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `allow_byok`                        | boolean                                                              | false    |              |                                                                                                                                                                                 |
| `anthropic`                         | [optimus-ide-collabsdk.AIBridgeAnthropicConfig](#optimus-ide-collabsdkaibridgeanthropicconfig) | false    |              | Deprecated: Use Providers with indexed `OPTIMUS-IDE-COLLAB_AI_GATEWAY_PROVIDER_<N>_*` env vars instead.                                                                                      |
| `api_dump_dir`                      | string                                                               | false    |              | Api dump dir is the base directory under which each provider's request/response dumps are written, in a subdirectory named after the provider. Empty disables dumping.          |
| `bedrock`                           | [optimus-ide-collabsdk.AIBridgeBedrockConfig](#optimus-ide-collabsdkaibridgebedrockconfig)     | false    |              | Deprecated: Use Providers with indexed `OPTIMUS-IDE-COLLAB_AI_GATEWAY_PROVIDER_<N>_*` env vars instead.                                                                                      |
| `budget_period`                     | string                                                               | false    |              |                                                                                                                                                                                 |
| `budget_policy`                     | string                                                               | false    |              | Budget settings for AI Governance cost controls.                                                                                                                                |
| `circuit_breaker_enabled`           | boolean                                                              | false    |              | Circuit breaker protects against cascading failures from upstream AI provider overload (503, 529).                                                                              |
| `circuit_breaker_failure_threshold` | integer                                                              | false    |              |                                                                                                                                                                                 |
| `circuit_breaker_interval`          | integer                                                              | false    |              |                                                                                                                                                                                 |
| `circuit_breaker_max_requests`      | integer                                                              | false    |              |                                                                                                                                                                                 |
| `circuit_breaker_timeout`           | integer                                                              | false    |              |                                                                                                                                                                                 |
| `enabled`                           | boolean                                                              | false    |              |                                                                                                                                                                                 |
| `inject_optimus-ide-collab_mcp_tools`            | boolean                                                              | false    |              | Deprecated: Injected MCP in AI Bridge is deprecated and will be removed in a future release.                                                                                    |
| `max_concurrency`                   | integer                                                              | false    |              |                                                                                                                                                                                 |
| `openai`                            | [optimus-ide-collabsdk.AIBridgeOpenAIConfig](#optimus-ide-collabsdkaibridgeopenaiconfig)       | false    |              | Deprecated: Use Providers with indexed `OPTIMUS-IDE-COLLAB_AI_GATEWAY_PROVIDER_<N>_*` env vars instead.                                                                                      |
| `providers`                         | array of [optimus-ide-collabsdk.AIProviderConfig](#optimus-ide-collabsdkaiproviderconfig)      | false    |              | Providers holds provider instances populated from `OPTIMUS-IDE-COLLAB_AI_GATEWAY_PROVIDER_<N>_<KEY>` env vars and/or the deprecated LegacyOpenAI/LegacyAnthropic/LegacyBedrock fields above. |
| `rate_limit`                        | integer                                                              | false    |              |                                                                                                                                                                                 |
| `retention`                         | integer                                                              | false    |              |                                                                                                                                                                                 |
| `send_actor_headers`                | boolean                                                              | false    |              |                                                                                                                                                                                 |
| `structured_logging`                | boolean                                                              | false    |              |                                                                                                                                                                                 |

## optimus-ide-collabsdk.AIBridgeListSessionsResponse

```json
{
  "count": 0,
  "sessions": [
    {
      "client": "string",
      "ended_at": "2019-08-24T14:15:22Z",
      "id": "string",
      "initiator": {
        "avatar_url": "http://example.com",
        "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
        "name": "string",
        "username": "string"
      },
      "last_active_at": "2019-08-24T14:15:22Z",
      "last_prompt": "string",
      "metadata": {
        "property1": null,
        "property2": null
      },
      "models": [
        "string"
      ],
      "network_calls": {
        "blocked": 0,
        "total": 0
      },
      "providers": [
        "string"
      ],
      "started_at": "2019-08-24T14:15:22Z",
      "threads": 0,
      "token_usage_summary": {
        "cache_read_input_tokens": 0,
        "cache_write_input_tokens": 0,
        "input_tokens": 0,
        "output_tokens": 0
      }
    }
  ]
}
```

### Properties

| Name       | Type                                                          | Required | Restrictions | Description |
|------------|---------------------------------------------------------------|----------|--------------|-------------|
| `count`    | integer                                                       | false    |              |             |
| `sessions` | array of [optimus-ide-collabsdk.AIBridgeSession](#optimus-ide-collabsdkaibridgesession) | false    |              |             |

## optimus-ide-collabsdk.AIBridgeModelThought

```json
{
  "text": "string"
}
```

### Properties

| Name   | Type   | Required | Restrictions | Description |
|--------|--------|----------|--------------|-------------|
| `text` | string | false    |              |             |

## optimus-ide-collabsdk.AIBridgeOpenAIConfig

```json
{
  "base_url": "string",
  "key": "string"
}
```

### Properties

| Name       | Type   | Required | Restrictions | Description |
|------------|--------|----------|--------------|-------------|
| `base_url` | string | false    |              |             |
| `key`      | string | false    |              |             |

## optimus-ide-collabsdk.AIBridgeProxyConfig

```json
{
  "allowed_private_cidrs": [
    "string"
  ],
  "api_dump_dir": "string",
  "cert_file": "string",
  "domain_allowlist": [
    "string"
  ],
  "enabled": true,
  "key_file": "string",
  "listen_addr": "string",
  "target": "string",
  "tls_cert_file": "string",
  "tls_key_file": "string",
  "upstream_proxy": "string",
  "upstream_proxy_ca": "string"
}
```

### Properties

| Name                    | Type            | Required | Restrictions | Description |
|-------------------------|-----------------|----------|--------------|-------------|
| `allowed_private_cidrs` | array of string | false    |              |             |
| `api_dump_dir`          | string          | false    |              |             |
| `cert_file`             | string          | false    |              |             |
| `domain_allowlist`      | array of string | false    |              |             |
| `enabled`               | boolean         | false    |              |             |
| `key_file`              | string          | false    |              |             |
| `listen_addr`           | string          | false    |              |             |
| `target`                | string          | false    |              |             |
| `tls_cert_file`         | string          | false    |              |             |
| `tls_key_file`          | string          | false    |              |             |
| `upstream_proxy`        | string          | false    |              |             |
| `upstream_proxy_ca`     | string          | false    |              |             |

## optimus-ide-collabsdk.AIBridgeSession

```json
{
  "client": "string",
  "ended_at": "2019-08-24T14:15:22Z",
  "id": "string",
  "initiator": {
    "avatar_url": "http://example.com",
    "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
    "name": "string",
    "username": "string"
  },
  "last_active_at": "2019-08-24T14:15:22Z",
  "last_prompt": "string",
  "metadata": {
    "property1": null,
    "property2": null
  },
  "models": [
    "string"
  ],
  "network_calls": {
    "blocked": 0,
    "total": 0
  },
  "providers": [
    "string"
  ],
  "started_at": "2019-08-24T14:15:22Z",
  "threads": 0,
  "token_usage_summary": {
    "cache_read_input_tokens": 0,
    "cache_write_input_tokens": 0,
    "input_tokens": 0,
    "output_tokens": 0
  }
}
```

### Properties

| Name                  | Type                                                                                     | Required | Restrictions | Description                                                                                                                                                                                                                           |
|-----------------------|------------------------------------------------------------------------------------------|----------|--------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `client`              | string                                                                                   | false    |              |                                                                                                                                                                                                                                       |
| `ended_at`            | string                                                                                   | false    |              |                                                                                                                                                                                                                                       |
| `id`                  | string                                                                                   | false    |              |                                                                                                                                                                                                                                       |
| `initiator`           | [optimus-ide-collabsdk.MinimalUser](#optimus-ide-collabsdkminimaluser)                                             | false    |              |                                                                                                                                                                                                                                       |
| `last_active_at`      | string                                                                                   | false    |              |                                                                                                                                                                                                                                       |
| `last_prompt`         | string                                                                                   | false    |              |                                                                                                                                                                                                                                       |
| `metadata`            | object                                                                                   | false    |              |                                                                                                                                                                                                                                       |
| » `[any property]`    | any                                                                                      | false    |              |                                                                                                                                                                                                                                       |
| `models`              | array of string                                                                          | false    |              |                                                                                                                                                                                                                                       |
| `network_calls`       | [optimus-ide-collabsdk.AIBridgeSessionNetworkCallSummary](#optimus-ide-collabsdkaibridgesessionnetworkcallsummary) | false    |              | Network calls summarizes the Agent Firewall network calls made during the session. A nil value means the session did not pass through Agent Firewall, so network call monitoring was not active, which the UI surfaces as "Disabled". |
| `providers`           | array of string                                                                          | false    |              |                                                                                                                                                                                                                                       |
| `started_at`          | string                                                                                   | false    |              |                                                                                                                                                                                                                                       |
| `threads`             | integer                                                                                  | false    |              |                                                                                                                                                                                                                                       |
| `token_usage_summary` | [optimus-ide-collabsdk.AIBridgeSessionTokenUsageSummary](#optimus-ide-collabsdkaibridgesessiontokenusagesummary)   | false    |              |                                                                                                                                                                                                                                       |

## optimus-ide-collabsdk.AIBridgeSessionNetworkCallSummary

```json
{
  "blocked": 0,
  "total": 0
}
```

### Properties

| Name      | Type    | Required | Restrictions | Description |
|-----------|---------|----------|--------------|-------------|
| `blocked` | integer | false    |              |             |
| `total`   | integer | false    |              |             |

## optimus-ide-collabsdk.AIBridgeSessionThreadsResponse

```json
{
  "client": "string",
  "ended_at": "2019-08-24T14:15:22Z",
  "id": "string",
  "initiator": {
    "avatar_url": "http://example.com",
    "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
    "name": "string",
    "username": "string"
  },
  "metadata": {
    "property1": null,
    "property2": null
  },
  "models": [
    "string"
  ],
  "page_ended_at": "2019-08-24T14:15:22Z",
  "page_started_at": "2019-08-24T14:15:22Z",
  "providers": [
    "string"
  ],
  "started_at": "2019-08-24T14:15:22Z",
  "threads": [
    {
      "agent_firewall_sequence_number": 0,
      "agent_firewall_session_id": "3735294f-18b1-4e7a-a269-99c30f0b30e7",
      "agentic_actions": [
        {
          "model": "string",
          "thinking": [
            {
              "text": "string"
            }
          ],
          "token_usage": {
            "cache_read_input_tokens": 0,
            "cache_write_input_tokens": 0,
            "input_tokens": 0,
            "metadata": {
              "property1": null,
              "property2": null
            },
            "output_tokens": 0
          },
          "tool_calls": [
            {
              "created_at": "2019-08-24T14:15:22Z",
              "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
              "injected": true,
              "input": "string",
              "interception_id": "34d9b688-63ad-46f4-88b5-665c1e7f7824",
              "metadata": {
                "property1": null,
                "property2": null
              },
              "provider_response_id": "string",
              "server_url": "string",
              "tool": "string"
            }
          ]
        }
      ],
      "credential_hint": "string",
      "credential_kind": "string",
      "ended_at": "2019-08-24T14:15:22Z",
      "error_message": "string",
      "error_type": "string",
      "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
      "model": "string",
      "prompt": "string",
      "provider": "string",
      "started_at": "2019-08-24T14:15:22Z",
      "token_usage": {
        "cache_read_input_tokens": 0,
        "cache_write_input_tokens": 0,
        "input_tokens": 0,
        "metadata": {
          "property1": null,
          "property2": null
        },
        "output_tokens": 0
      }
    }
  ],
  "token_usage_summary": {
    "cache_read_input_tokens": 0,
    "cache_write_input_tokens": 0,
    "input_tokens": 0,
    "metadata": {
      "property1": null,
      "property2": null
    },
    "output_tokens": 0
  }
}
```

### Properties

| Name                  | Type                                                                                   | Required | Restrictions | Description |
|-----------------------|----------------------------------------------------------------------------------------|----------|--------------|-------------|
| `client`              | string                                                                                 | false    |              |             |
| `ended_at`            | string                                                                                 | false    |              |             |
| `id`                  | string                                                                                 | false    |              |             |
| `initiator`           | [optimus-ide-collabsdk.MinimalUser](#optimus-ide-collabsdkminimaluser)                                           | false    |              |             |
| `metadata`            | object                                                                                 | false    |              |             |
| » `[any property]`    | any                                                                                    | false    |              |             |
| `models`              | array of string                                                                        | false    |              |             |
| `page_ended_at`       | string                                                                                 | false    |              |             |
| `page_started_at`     | string                                                                                 | false    |              |             |
| `providers`           | array of string                                                                        | false    |              |             |
| `started_at`          | string                                                                                 | false    |              |             |
| `threads`             | array of [optimus-ide-collabsdk.AIBridgeThread](#optimus-ide-collabsdkaibridgethread)                            | false    |              |             |
| `token_usage_summary` | [optimus-ide-collabsdk.AIBridgeSessionThreadsTokenUsage](#optimus-ide-collabsdkaibridgesessionthreadstokenusage) | false    |              |             |

## optimus-ide-collabsdk.AIBridgeSessionThreadsTokenUsage

```json
{
  "cache_read_input_tokens": 0,
  "cache_write_input_tokens": 0,
  "input_tokens": 0,
  "metadata": {
    "property1": null,
    "property2": null
  },
  "output_tokens": 0
}
```

### Properties

| Name                       | Type    | Required | Restrictions | Description |
|----------------------------|---------|----------|--------------|-------------|
| `cache_read_input_tokens`  | integer | false    |              |             |
| `cache_write_input_tokens` | integer | false    |              |             |
| `input_tokens`             | integer | false    |              |             |
| `metadata`                 | object  | false    |              |             |
| » `[any property]`         | any     | false    |              |             |
| `output_tokens`            | integer | false    |              |             |

## optimus-ide-collabsdk.AIBridgeSessionTokenUsageSummary

```json
{
  "cache_read_input_tokens": 0,
  "cache_write_input_tokens": 0,
  "input_tokens": 0,
  "output_tokens": 0
}
```

### Properties

| Name                       | Type    | Required | Restrictions | Description |
|----------------------------|---------|----------|--------------|-------------|
| `cache_read_input_tokens`  | integer | false    |              |             |
| `cache_write_input_tokens` | integer | false    |              |             |
| `input_tokens`             | integer | false    |              |             |
| `output_tokens`            | integer | false    |              |             |

## optimus-ide-collabsdk.AIBridgeThread

```json
{
  "agent_firewall_sequence_number": 0,
  "agent_firewall_session_id": "3735294f-18b1-4e7a-a269-99c30f0b30e7",
  "agentic_actions": [
    {
      "model": "string",
      "thinking": [
        {
          "text": "string"
        }
      ],
      "token_usage": {
        "cache_read_input_tokens": 0,
        "cache_write_input_tokens": 0,
        "input_tokens": 0,
        "metadata": {
          "property1": null,
          "property2": null
        },
        "output_tokens": 0
      },
      "tool_calls": [
        {
          "created_at": "2019-08-24T14:15:22Z",
          "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
          "injected": true,
          "input": "string",
          "interception_id": "34d9b688-63ad-46f4-88b5-665c1e7f7824",
          "metadata": {
            "property1": null,
            "property2": null
          },
          "provider_response_id": "string",
          "server_url": "string",
          "tool": "string"
        }
      ]
    }
  ],
  "credential_hint": "string",
  "credential_kind": "string",
  "ended_at": "2019-08-24T14:15:22Z",
  "error_message": "string",
  "error_type": "string",
  "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
  "model": "string",
  "prompt": "string",
  "provider": "string",
  "started_at": "2019-08-24T14:15:22Z",
  "token_usage": {
    "cache_read_input_tokens": 0,
    "cache_write_input_tokens": 0,
    "input_tokens": 0,
    "metadata": {
      "property1": null,
      "property2": null
    },
    "output_tokens": 0
  }
}
```

### Properties

| Name                             | Type                                                                                   | Required | Restrictions | Description                                                                                                                                                                                                                               |
|----------------------------------|----------------------------------------------------------------------------------------|----------|--------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `agent_firewall_sequence_number` | integer                                                                                | false    |              | Agent firewall sequence number is the firewall sequence number from the root interception. Used to determine the position of this LLM request in the firewall event stream. Nil when the request did not pass through the agent firewall. |
| `agent_firewall_session_id`      | string                                                                                 | false    |              | Agent firewall session ID links this thread to an agent firewall confinement session. Nil when the request did not pass through the agent firewall.                                                                                       |
| `agentic_actions`                | array of [optimus-ide-collabsdk.AIBridgeAgenticAction](#optimus-ide-collabsdkaibridgeagenticaction)              | false    |              |                                                                                                                                                                                                                                           |
| `credential_hint`                | string                                                                                 | false    |              |                                                                                                                                                                                                                                           |
| `credential_kind`                | string                                                                                 | false    |              |                                                                                                                                                                                                                                           |
| `ended_at`                       | string                                                                                 | false    |              |                                                                                                                                                                                                                                           |
| `error_message`                  | string                                                                                 | false    |              | Error message is the raw terminal upstream error message from the root interception. Nil when the interception succeeded.                                                                                                                 |
| `error_type`                     | string                                                                                 | false    |              | Error type is the categorized terminal upstream error from the root interception, or nil when the interception succeeded. See the aibridge_interception_error_type enum for possible values.                                              |
| `id`                             | string                                                                                 | false    |              |                                                                                                                                                                                                                                           |
| `model`                          | string                                                                                 | false    |              |                                                                                                                                                                                                                                           |
| `prompt`                         | string                                                                                 | false    |              |                                                                                                                                                                                                                                           |
| `provider`                       | string                                                                                 | false    |              |                                                                                                                                                                                                                                           |
| `started_at`                     | string                                                                                 | false    |              |                                                                                                                                                                                                                                           |
| `token_usage`                    | [optimus-ide-collabsdk.AIBridgeSessionThreadsTokenUsage](#optimus-ide-collabsdkaibridgesessionthreadstokenusage) | false    |              |                                                                                                                                                                                                                                           |

## optimus-ide-collabsdk.AIBridgeToolCall

```json
{
  "created_at": "2019-08-24T14:15:22Z",
  "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
  "injected": true,
  "input": "string",
  "interception_id": "34d9b688-63ad-46f4-88b5-665c1e7f7824",
  "metadata": {
    "property1": null,
    "property2": null
  },
  "provider_response_id": "string",
  "server_url": "string",
  "tool": "string"
}
```

### Properties

| Name                   | Type    | Required | Restrictions | Description |
|------------------------|---------|----------|--------------|-------------|
| `created_at`           | string  | false    |              |             |
| `id`                   | string  | false    |              |             |
| `injected`             | boolean | false    |              |             |
| `input`                | string  | false    |              |             |
| `interception_id`      | string  | false    |              |             |
| `metadata`             | object  | false    |              |             |
| » `[any property]`     | any     | false    |              |             |
| `provider_response_id` | string  | false    |              |             |
| `server_url`           | string  | false    |              |             |
| `tool`                 | string  | false    |              |             |

## optimus-ide-collabsdk.AIBudgetLimitSource

```json
"user_override"
```

### Properties

#### Enumerated Values

| Value(s)                 |
|--------------------------|
| `group`, `user_override` |

## optimus-ide-collabsdk.AIConfig

```json
{
  "aibridge_proxy": {
    "allowed_private_cidrs": [
      "string"
    ],
    "api_dump_dir": "string",
    "cert_file": "string",
    "domain_allowlist": [
      "string"
    ],
    "enabled": true,
    "key_file": "string",
    "listen_addr": "string",
    "target": "string",
    "tls_cert_file": "string",
    "tls_key_file": "string",
    "upstream_proxy": "string",
    "upstream_proxy_ca": "string"
  },
  "bridge": {
    "allow_byok": true,
    "anthropic": {
      "base_url": "string",
      "key": "string"
    },
    "api_dump_dir": "string",
    "bedrock": {
      "access_key": "string",
      "access_key_secret": "string",
      "base_url": "string",
      "model": "string",
      "region": "string",
      "small_fast_model": "string"
    },
    "budget_period": "string",
    "budget_policy": "string",
    "circuit_breaker_enabled": true,
    "circuit_breaker_failure_threshold": 0,
    "circuit_breaker_interval": 0,
    "circuit_breaker_max_requests": 0,
    "circuit_breaker_timeout": 0,
    "enabled": true,
    "inject_optimus-ide-collab_mcp_tools": true,
    "max_concurrency": 0,
    "openai": {
      "base_url": "string",
      "key": "string"
    },
    "providers": [
      {
        "base_url": "string",
        "bedrock_model": "string",
        "bedrock_region": "string",
        "bedrock_small_fast_model": "string",
        "name": "string",
        "type": "string"
      }
    ],
    "rate_limit": 0,
    "retention": 0,
    "send_actor_headers": true,
    "structured_logging": true
  },
  "chat": {
    "acquire_batch_size": 0,
    "debug_logging_enabled": true
  }
}
```

### Properties

| Name             | Type                                                         | Required | Restrictions | Description |
|------------------|--------------------------------------------------------------|----------|--------------|-------------|
| `aibridge_proxy` | [optimus-ide-collabsdk.AIBridgeProxyConfig](#optimus-ide-collabsdkaibridgeproxyconfig) | false    |              |             |
| `bridge`         | [optimus-ide-collabsdk.AIBridgeConfig](#optimus-ide-collabsdkaibridgeconfig)           | false    |              |             |
| `chat`           | [optimus-ide-collabsdk.ChatConfig](#optimus-ide-collabsdkchatconfig)                   | false    |              |             |

## optimus-ide-collabsdk.AIGatewayKey

```json
{
  "created_at": "2019-08-24T14:15:22Z",
  "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
  "key_prefix": "string",
  "last_heartbeat_at": "2019-08-24T14:15:22Z",
  "name": "string"
}
```

### Properties

| Name                | Type   | Required | Restrictions | Description |
|---------------------|--------|----------|--------------|-------------|
| `created_at`        | string | false    |              |             |
| `id`                | string | false    |              |             |
| `key_prefix`        | string | false    |              |             |
| `last_heartbeat_at` | string | false    |              |             |
| `name`              | string | false    |              |             |

## optimus-ide-collabsdk.AIGroupBudget

```json
{
  "limit_source": "user_override",
  "spend_limit_micros": 0
}
```

### Properties

| Name                 | Type                                                         | Required | Restrictions | Description |
|----------------------|--------------------------------------------------------------|----------|--------------|-------------|
| `limit_source`       | [optimus-ide-collabsdk.AIBudgetLimitSource](#optimus-ide-collabsdkaibudgetlimitsource) | false    |              |             |
| `spend_limit_micros` | integer                                                      | false    |              |             |

## optimus-ide-collabsdk.AIProvider

```json
{
  "api_keys": [
    {
      "created_at": "2019-08-24T14:15:22Z",
      "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
      "masked": "string"
    }
  ],
  "base_url": "string",
  "created_at": "2019-08-24T14:15:22Z",
  "display_name": "string",
  "enabled": true,
  "icon": "string",
  "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
  "name": "string",
  "settings": {},
  "type": "openai",
  "updated_at": "2019-08-24T14:15:22Z"
}
```

### Properties

| Name           | Type                                                       | Required | Restrictions | Description |
|----------------|------------------------------------------------------------|----------|--------------|-------------|
| `api_keys`     | array of [optimus-ide-collabsdk.AIProviderKey](#optimus-ide-collabsdkaiproviderkey)  | false    |              |             |
| `base_url`     | string                                                     | false    |              |             |
| `created_at`   | string                                                     | false    |              |             |
| `display_name` | string                                                     | false    |              |             |
| `enabled`      | boolean                                                    | false    |              |             |
| `icon`         | string                                                     | false    |              |             |
| `id`           | string                                                     | false    |              |             |
| `name`         | string                                                     | false    |              |             |
| `settings`     | [optimus-ide-collabsdk.AIProviderSettings](#optimus-ide-collabsdkaiprovidersettings) | false    |              |             |
| `type`         | [optimus-ide-collabsdk.AIProviderType](#optimus-ide-collabsdkaiprovidertype)         | false    |              |             |
| `updated_at`   | string                                                     | false    |              |             |

## optimus-ide-collabsdk.AIProviderConfig

```json
{
  "base_url": "string",
  "bedrock_model": "string",
  "bedrock_region": "string",
  "bedrock_small_fast_model": "string",
  "name": "string",
  "type": "string"
}
```

### Properties

| Name                       | Type   | Required | Restrictions | Description                                                                                                                                           |
|----------------------------|--------|----------|--------------|-------------------------------------------------------------------------------------------------------------------------------------------------------|
| `base_url`                 | string | false    |              | Base URL is the base URL of the upstream provider API.                                                                                                |
| `bedrock_model`            | string | false    |              |                                                                                                                                                       |
| `bedrock_region`           | string | false    |              |                                                                                                                                                       |
| `bedrock_small_fast_model` | string | false    |              |                                                                                                                                                       |
| `name`                     | string | false    |              | Name is the unique instance identifier used for routing. Defaults to Type if not provided.                                                            |
| `type`                     | string | false    |              | Type is the provider type. Valid values are: "openai", "anthropic", "azure", "bedrock", "google", "openai-compat", "openrouter", "vercel", "copilot". |

## optimus-ide-collabsdk.AIProviderKey

```json
{
  "created_at": "2019-08-24T14:15:22Z",
  "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
  "masked": "string"
}
```

### Properties

| Name         | Type   | Required | Restrictions | Description |
|--------------|--------|----------|--------------|-------------|
| `created_at` | string | false    |              |             |
| `id`         | string | false    |              |             |
| `masked`     | string | false    |              |             |

## optimus-ide-collabsdk.AIProviderKeyMutation

```json
{
  "api_key": "string",
  "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08"
}
```

### Properties

| Name      | Type   | Required | Restrictions | Description |
|-----------|--------|----------|--------------|-------------|
| `api_key` | string | false    |              |             |
| `id`      | string | false    |              |             |

## optimus-ide-collabsdk.AIProviderSettings

```json
{}
```

### Properties

None

## optimus-ide-collabsdk.AIProviderType

```json
"openai"
```

### Properties

#### Enumerated Values

| Value(s)                                                                                                |
|---------------------------------------------------------------------------------------------------------|
| `anthropic`, `azure`, `bedrock`, `copilot`, `google`, `openai`, `openai-compat`, `openrouter`, `vercel` |

## optimus-ide-collabsdk.APIAllowListTarget

```json
{
  "id": "string",
  "type": "*"
}
```

### Properties

| Name   | Type                                           | Required | Restrictions | Description |
|--------|------------------------------------------------|----------|--------------|-------------|
| `id`   | string                                         | false    |              |             |
| `type` | [optimus-ide-collabsdk.RBACResource](#optimus-ide-collabsdkrbacresource) | false    |              |             |

## optimus-ide-collabsdk.APIKey

```json
{
  "allow_list": [
    {
      "id": "string",
      "type": "*"
    }
  ],
  "created_at": "2019-08-24T14:15:22Z",
  "expires_at": "2019-08-24T14:15:22Z",
  "id": "string",
  "last_used": "2019-08-24T14:15:22Z",
  "lifetime_seconds": 0,
  "login_type": "password",
  "scope": "all",
  "scopes": [
    "all"
  ],
  "token_name": "string",
  "updated_at": "2019-08-24T14:15:22Z",
  "user_id": "a169451c-8525-4352-b8ca-070dd449a1a5"
}
```

### Properties

| Name               | Type                                                                | Required | Restrictions | Description                     |
|--------------------|---------------------------------------------------------------------|----------|--------------|---------------------------------|
| `allow_list`       | array of [optimus-ide-collabsdk.APIAllowListTarget](#optimus-ide-collabsdkapiallowlisttarget) | false    |              |                                 |
| `created_at`       | string                                                              | true     |              |                                 |
| `expires_at`       | string                                                              | true     |              |                                 |
| `id`               | string                                                              | true     |              |                                 |
| `last_used`        | string                                                              | true     |              |                                 |
| `lifetime_seconds` | integer                                                             | true     |              |                                 |
| `login_type`       | [optimus-ide-collabsdk.LoginType](#optimus-ide-collabsdklogintype)                            | true     |              |                                 |
| `scope`            | [optimus-ide-collabsdk.APIKeyScope](#optimus-ide-collabsdkapikeyscope)                        | false    |              | Deprecated: use Scopes instead. |
| `scopes`           | array of [optimus-ide-collabsdk.APIKeyScope](#optimus-ide-collabsdkapikeyscope)               | false    |              |                                 |
| `token_name`       | string                                                              | true     |              |                                 |
| `updated_at`       | string                                                              | true     |              |                                 |
| `user_id`          | string                                                              | true     |              |                                 |

#### Enumerated Values

| Property     | Value(s)                              |
|--------------|---------------------------------------|
| `login_type` | `github`, `oidc`, `password`, `token` |
| `scope`      | `all`, `application_connect`          |

## optimus-ide-collabsdk.APIKeyScope

```json
"all"
```

### Properties

#### Enumerated Values

| Value(s)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                 |
|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `ai_gateway_key:*`, `ai_gateway_key:create`, `ai_gateway_key:delete`, `ai_gateway_key:read`, `ai_gateway_key:update`, `ai_model_price:*`, `ai_model_price:read`, `ai_model_price:update`, `ai_provider:*`, `ai_provider:create`, `ai_provider:delete`, `ai_provider:read`, `ai_provider:update`, `ai_seat:*`, `ai_seat:create`, `ai_seat:read`, `aibridge_interception:*`, `aibridge_interception:create`, `aibridge_interception:read`, `aibridge_interception:update`, `all`, `api_key:*`, `api_key:create`, `api_key:delete`, `api_key:read`, `api_key:update`, `application_connect`, `assign_org_role:*`, `assign_org_role:assign`, `assign_org_role:create`, `assign_org_role:delete`, `assign_org_role:read`, `assign_org_role:unassign`, `assign_org_role:update`, `assign_role:*`, `assign_role:assign`, `assign_role:read`, `assign_role:unassign`, `audit_log:*`, `audit_log:create`, `audit_log:read`, `boundary_log:*`, `boundary_log:create`, `boundary_log:delete`, `boundary_log:read`, `boundary_usage:*`, `boundary_usage:delete`, `boundary_usage:read`, `boundary_usage:update`, `chat:*`, `chat:create`, `chat:delete`, `chat:read`, `chat:share`, `chat:update`, `optimus-ide-collab:all`, `optimus-ide-collab:apikeys.manage_self`, `optimus-ide-collab:application_connect`, `optimus-ide-collab:templates.author`, `optimus-ide-collab:templates.build`, `optimus-ide-collab:workspaces.access`, `optimus-ide-collab:workspaces.create`, `optimus-ide-collab:workspaces.delete`, `optimus-ide-collab:workspaces.operate`, `connection_log:*`, `connection_log:read`, `connection_log:update`, `crypto_key:*`, `crypto_key:create`, `crypto_key:delete`, `crypto_key:read`, `crypto_key:update`, `debug_info:*`, `debug_info:read`, `deployment_config:*`, `deployment_config:read`, `deployment_config:update`, `deployment_stats:*`, `deployment_stats:read`, `file:*`, `file:create`, `file:read`, `group:*`, `group:create`, `group:delete`, `group:read`, `group:update`, `group_member:*`, `group_member:read`, `idpsync_settings:*`, `idpsync_settings:read`, `idpsync_settings:update`, `inbox_notification:*`, `inbox_notification:create`, `inbox_notification:read`, `inbox_notification:update`, `license:*`, `license:create`, `license:delete`, `license:read`, `notification_message:*`, `notification_message:create`, `notification_message:delete`, `notification_message:read`, `notification_message:update`, `notification_preference:*`, `notification_preference:read`, `notification_preference:update`, `notification_template:*`, `notification_template:read`, `notification_template:update`, `oauth2_app:*`, `oauth2_app:create`, `oauth2_app:delete`, `oauth2_app:read`, `oauth2_app:update`, `oauth2_app_code_token:*`, `oauth2_app_code_token:create`, `oauth2_app_code_token:delete`, `oauth2_app_code_token:read`, `oauth2_app_secret:*`, `oauth2_app_secret:create`, `oauth2_app_secret:delete`, `oauth2_app_secret:read`, `oauth2_app_secret:update`, `organization:*`, `organization:create`, `organization:delete`, `organization:read`, `organization:update`, `organization_member:*`, `organization_member:create`, `organization_member:delete`, `organization_member:read`, `organization_member:update`, `prebuilt_workspace:*`, `prebuilt_workspace:delete`, `prebuilt_workspace:update`, `provisioner_daemon:*`, `provisioner_daemon:create`, `provisioner_daemon:delete`, `provisioner_daemon:read`, `provisioner_daemon:update`, `provisioner_jobs:*`, `provisioner_jobs:create`, `provisioner_jobs:read`, `provisioner_jobs:update`, `replicas:*`, `replicas:read`, `system:*`, `system:create`, `system:delete`, `system:read`, `system:update`, `tailnet_coordinator:*`, `tailnet_coordinator:create`, `tailnet_coordinator:delete`, `tailnet_coordinator:read`, `tailnet_coordinator:update`, `task:*`, `task:create`, `task:delete`, `task:read`, `task:update`, `template:*`, `template:create`, `template:delete`, `template:read`, `template:update`, `template:use`, `template:view_insights`, `usage_event:*`, `usage_event:create`, `usage_event:read`, `usage_event:update`, `user:*`, `user:create`, `user:delete`, `user:read`, `user:read_personal`, `user:update`, `user:update_personal`, `user_secret:*`, `user_secret:create`, `user_secret:delete`, `user_secret:read`, `user_secret:update`, `user_skill:*`, `user_skill:create`, `user_skill:delete`, `user_skill:read`, `user_skill:update`, `webpush_subscription:*`, `webpush_subscription:create`, `webpush_subscription:delete`, `webpush_subscription:read`, `workspace:*`, `workspace:application_connect`, `workspace:create`, `workspace:create_agent`, `workspace:delete`, `workspace:delete_agent`, `workspace:read`, `workspace:share`, `workspace:ssh`, `workspace:start`, `workspace:stop`, `workspace:update`, `workspace:update_agent`, `workspace_agent_devcontainers:*`, `workspace_agent_devcontainers:create`, `workspace_agent_resource_monitor:*`, `workspace_agent_resource_monitor:create`, `workspace_agent_resource_monitor:read`, `workspace_agent_resource_monitor:update`, `workspace_build_orchestration:*`, `workspace_build_orchestration:create`, `workspace_build_orchestration:delete`, `workspace_build_orchestration:read`, `workspace_build_orchestration:update`, `workspace_dormant:*`, `workspace_dormant:application_connect`, `workspace_dormant:create`, `workspace_dormant:create_agent`, `workspace_dormant:delete`, `workspace_dormant:delete_agent`, `workspace_dormant:read`, `workspace_dormant:share`, `workspace_dormant:ssh`, `workspace_dormant:start`, `workspace_dormant:stop`, `workspace_dormant:update`, `workspace_dormant:update_agent`, `workspace_proxy:*`, `workspace_proxy:create`, `workspace_proxy:delete`, `workspace_proxy:read`, `workspace_proxy:update` |

## optimus-ide-collabsdk.AddLicenseRequest

```json
{
  "license": "string"
}
```

### Properties

| Name      | Type   | Required | Restrictions | Description |
|-----------|--------|----------|--------------|-------------|
| `license` | string | true     |              |             |

## optimus-ide-collabsdk.AgentChatSendShortcut

```json
"enter"
```

### Properties

#### Enumerated Values

| Value(s)                  |
|---------------------------|
| `enter`, `modifier_enter` |

## optimus-ide-collabsdk.AgentConnectionTiming

```json
{
  "ended_at": "2019-08-24T14:15:22Z",
  "stage": "init",
  "started_at": "2019-08-24T14:15:22Z",
  "workspace_agent_id": "string",
  "workspace_agent_name": "string"
}
```

### Properties

| Name                   | Type                                         | Required | Restrictions | Description |
|------------------------|----------------------------------------------|----------|--------------|-------------|
| `ended_at`             | string                                       | false    |              |             |
| `stage`                | [optimus-ide-collabsdk.TimingStage](#optimus-ide-collabsdktimingstage) | false    |              |             |
| `started_at`           | string                                       | false    |              |             |
| `workspace_agent_id`   | string                                       | false    |              |             |
| `workspace_agent_name` | string                                       | false    |              |             |

## optimus-ide-collabsdk.AgentDisplayMode

```json
"auto"
```

### Properties

#### Enumerated Values

| Value(s)                                      |
|-----------------------------------------------|
| `always_collapsed`, `always_expanded`, `auto` |

## optimus-ide-collabsdk.AgentFirewallLog

```json
{
  "allowed": true,
  "captured_at": "2019-08-24T14:15:22Z",
  "created_at": "2019-08-24T14:15:22Z",
  "detail": "string",
  "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
  "matched_rule": "string",
  "method": "string",
  "proto": "string",
  "sequence_number": 0,
  "session_id": "1ffd059c-17ea-40a8-8aef-70fd0307db82"
}
```

### Properties

| Name              | Type    | Required | Restrictions | Description |
|-------------------|---------|----------|--------------|-------------|
| `allowed`         | boolean | false    |              |             |
| `captured_at`     | string  | false    |              |             |
| `created_at`      | string  | false    |              |             |
| `detail`          | string  | false    |              |             |
| `id`              | string  | false    |              |             |
| `matched_rule`    | string  | false    |              |             |
| `method`          | string  | false    |              |             |
| `proto`           | string  | false    |              |             |
| `sequence_number` | integer | false    |              |             |
| `session_id`      | string  | false    |              |             |

## optimus-ide-collabsdk.AgentFirewallSession

```json
{
  "confined_process": "string",
  "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
  "owner_id": "8826ee2e-7933-4665-aef2-2393f84a0d05",
  "started_at": "2019-08-24T14:15:22Z",
  "workspace_id": "0967198e-ec7b-4c6b-b4d3-f71244cadbe9"
}
```

### Properties

| Name               | Type   | Required | Restrictions | Description |
|--------------------|--------|----------|--------------|-------------|
| `confined_process` | string | false    |              |             |
| `id`               | string | false    |              |             |
| `owner_id`         | string | false    |              |             |
| `started_at`       | string | false    |              |             |
| `workspace_id`     | string | false    |              |             |

## optimus-ide-collabsdk.AgentFirewallSessionLogsResponse

```json
{
  "results": [
    {
      "allowed": true,
      "captured_at": "2019-08-24T14:15:22Z",
      "created_at": "2019-08-24T14:15:22Z",
      "detail": "string",
      "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
      "matched_rule": "string",
      "method": "string",
      "proto": "string",
      "sequence_number": 0,
      "session_id": "1ffd059c-17ea-40a8-8aef-70fd0307db82"
    }
  ]
}
```

### Properties

| Name      | Type                                                            | Required | Restrictions | Description |
|-----------|-----------------------------------------------------------------|----------|--------------|-------------|
| `results` | array of [optimus-ide-collabsdk.AgentFirewallLog](#optimus-ide-collabsdkagentfirewalllog) | false    |              |             |

## optimus-ide-collabsdk.AgentScriptTiming

```json
{
  "display_name": "string",
  "ended_at": "2019-08-24T14:15:22Z",
  "exit_code": 0,
  "stage": "init",
  "started_at": "2019-08-24T14:15:22Z",
  "status": "string",
  "workspace_agent_id": "string",
  "workspace_agent_name": "string"
}
```

### Properties

| Name                   | Type                                         | Required | Restrictions | Description |
|------------------------|----------------------------------------------|----------|--------------|-------------|
| `display_name`         | string                                       | false    |              |             |
| `ended_at`             | string                                       | false    |              |             |
| `exit_code`            | integer                                      | false    |              |             |
| `stage`                | [optimus-ide-collabsdk.TimingStage](#optimus-ide-collabsdktimingstage) | false    |              |             |
| `started_at`           | string                                       | false    |              |             |
| `status`               | string                                       | false    |              |             |
| `workspace_agent_id`   | string                                       | false    |              |             |
| `workspace_agent_name` | string                                       | false    |              |             |

## optimus-ide-collabsdk.AgentSubsystem

```json
"envbox"
```

### Properties

#### Enumerated Values

| Value(s)                            |
|-------------------------------------|
| `envbox`, `envbuilder`, `exectrace` |

## optimus-ide-collabsdk.AppHostResponse

```json
{
  "host": "string"
}
```

### Properties

| Name   | Type   | Required | Restrictions | Description                                                   |
|--------|--------|----------|--------------|---------------------------------------------------------------|
| `host` | string | false    |              | Host is the externally accessible URL for the Optimus-IDE-Collab instance. |

## optimus-ide-collabsdk.AppearanceConfig

```json
{
  "announcement_banners": [
    {
      "background_color": "string",
      "enabled": true,
      "message": "string"
    }
  ],
  "application_name": "string",
  "docs_url": "string",
  "logo_url": "string",
  "service_banner": {
    "background_color": "string",
    "enabled": true,
    "message": "string"
  },
  "support_links": [
    {
      "icon": "bug",
      "location": "navbar",
      "name": "string",
      "target": "string"
    }
  ]
}
```

### Properties

| Name                   | Type                                                    | Required | Restrictions | Description                                                         |
|------------------------|---------------------------------------------------------|----------|--------------|---------------------------------------------------------------------|
| `announcement_banners` | array of [optimus-ide-collabsdk.BannerConfig](#optimus-ide-collabsdkbannerconfig) | false    |              |                                                                     |
| `application_name`     | string                                                  | false    |              |                                                                     |
| `docs_url`             | string                                                  | false    |              |                                                                     |
| `logo_url`             | string                                                  | false    |              |                                                                     |
| `service_banner`       | [optimus-ide-collabsdk.BannerConfig](#optimus-ide-collabsdkbannerconfig)          | false    |              | Deprecated: ServiceBanner has been replaced by AnnouncementBanners. |
| `support_links`        | array of [optimus-ide-collabsdk.LinkConfig](#optimus-ide-collabsdklinkconfig)     | false    |              |                                                                     |

## optimus-ide-collabsdk.ArchiveTemplateVersionsRequest

```json
{
  "all": true
}
```

### Properties

| Name  | Type    | Required | Restrictions | Description                                                                                                              |
|-------|---------|----------|--------------|--------------------------------------------------------------------------------------------------------------------------|
| `all` | boolean | false    |              | By default, only failed versions are archived. Set this to true to archive all unused versions regardless of job status. |

## optimus-ide-collabsdk.AssignableRoles

```json
{
  "assignable": true,
  "built_in": true,
  "display_name": "string",
  "name": "string",
  "organization_id": "7c60d51f-b44e-4682-87d6-449835ea4de6",
  "organization_member_permissions": [
    {
      "action": "application_connect",
      "negate": true,
      "resource_type": "*"
    }
  ],
  "organization_permissions": [
    {
      "action": "application_connect",
      "negate": true,
      "resource_type": "*"
    }
  ],
  "site_permissions": [
    {
      "action": "application_connect",
      "negate": true,
      "resource_type": "*"
    }
  ],
  "user_permissions": [
    {
      "action": "application_connect",
      "negate": true,
      "resource_type": "*"
    }
  ]
}
```

### Properties

| Name                              | Type                                                | Required | Restrictions | Description                                                                                            |
|-----------------------------------|-----------------------------------------------------|----------|--------------|--------------------------------------------------------------------------------------------------------|
| `assignable`                      | boolean                                             | false    |              |                                                                                                        |
| `built_in`                        | boolean                                             | false    |              | Built in roles are immutable                                                                           |
| `display_name`                    | string                                              | false    |              |                                                                                                        |
| `name`                            | string                                              | false    |              |                                                                                                        |
| `organization_id`                 | string                                              | false    |              |                                                                                                        |
| `organization_member_permissions` | array of [optimus-ide-collabsdk.Permission](#optimus-ide-collabsdkpermission) | false    |              | Organization member permissions are specific for the organization in the field 'OrganizationID' above. |
| `organization_permissions`        | array of [optimus-ide-collabsdk.Permission](#optimus-ide-collabsdkpermission) | false    |              | Organization permissions are specific for the organization in the field 'OrganizationID' above.        |
| `site_permissions`                | array of [optimus-ide-collabsdk.Permission](#optimus-ide-collabsdkpermission) | false    |              |                                                                                                        |
| `user_permissions`                | array of [optimus-ide-collabsdk.Permission](#optimus-ide-collabsdkpermission) | false    |              |                                                                                                        |

## optimus-ide-collabsdk.AuditAction

```json
"create"
```

### Properties

#### Enumerated Values

| Value(s)                                                                                                                                        |
|-------------------------------------------------------------------------------------------------------------------------------------------------|
| `close`, `connect`, `create`, `delete`, `disconnect`, `login`, `logout`, `open`, `register`, `request_password_reset`, `start`, `stop`, `write` |

## optimus-ide-collabsdk.AuditDiff

```json
{
  "property1": {
    "new": null,
    "old": null,
    "secret": true
  },
  "property2": {
    "new": null,
    "old": null,
    "secret": true
  }
}
```

### Properties

| Name             | Type                                               | Required | Restrictions | Description |
|------------------|----------------------------------------------------|----------|--------------|-------------|
| `[any property]` | [optimus-ide-collabsdk.AuditDiffField](#optimus-ide-collabsdkauditdifffield) | false    |              |             |

## optimus-ide-collabsdk.AuditDiffField

```json
{
  "new": null,
  "old": null,
  "secret": true
}
```

### Properties

| Name     | Type    | Required | Restrictions | Description |
|----------|---------|----------|--------------|-------------|
| `new`    | any     | false    |              |             |
| `old`    | any     | false    |              |             |
| `secret` | boolean | false    |              |             |

## optimus-ide-collabsdk.AuditLog

```json
{
  "action": "create",
  "additional_fields": {},
  "description": "string",
  "diff": {
    "property1": {
      "new": null,
      "old": null,
      "secret": true
    },
    "property2": {
      "new": null,
      "old": null,
      "secret": true
    }
  },
  "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
  "ip": "string",
  "is_deleted": true,
  "organization": {
    "display_name": "string",
    "icon": "string",
    "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
    "name": "string"
  },
  "organization_id": "7c60d51f-b44e-4682-87d6-449835ea4de6",
  "request_id": "266ea41d-adf5-480b-af50-15b940c2b846",
  "resource_icon": "string",
  "resource_id": "4d5215ed-38bb-48ed-879a-fdb9ca58522f",
  "resource_link": "string",
  "resource_target": "string",
  "resource_type": "template",
  "status_code": 0,
  "time": "2019-08-24T14:15:22Z",
  "user": {
    "avatar_url": "http://example.com",
    "created_at": "2019-08-24T14:15:22Z",
    "email": "user@example.com",
    "has_ai_seat": true,
    "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
    "is_service_account": true,
    "last_seen_at": "2019-08-24T14:15:22Z",
    "login_type": "",
    "name": "string",
    "organization_ids": [
      "497f6eca-6276-4993-bfeb-53cbbbba6f08"
    ],
    "roles": [
      {
        "display_name": "string",
        "name": "string",
        "organization_id": "string"
      }
    ],
    "status": "active",
    "theme_preference": "string",
    "updated_at": "2019-08-24T14:15:22Z",
    "username": "string"
  },
  "user_agent": "string"
}
```

### Properties

| Name                | Type                                                         | Required | Restrictions | Description                                  |
|---------------------|--------------------------------------------------------------|----------|--------------|----------------------------------------------|
| `action`            | [optimus-ide-collabsdk.AuditAction](#optimus-ide-collabsdkauditaction)                 | false    |              |                                              |
| `additional_fields` | object                                                       | false    |              |                                              |
| `description`       | string                                                       | false    |              |                                              |
| `diff`              | [optimus-ide-collabsdk.AuditDiff](#optimus-ide-collabsdkauditdiff)                     | false    |              |                                              |
| `id`                | string                                                       | false    |              |                                              |
| `ip`                | string                                                       | false    |              |                                              |
| `is_deleted`        | boolean                                                      | false    |              |                                              |
| `organization`      | [optimus-ide-collabsdk.MinimalOrganization](#optimus-ide-collabsdkminimalorganization) | false    |              |                                              |
| `organization_id`   | string                                                       | false    |              | Deprecated: Use 'organization.id' instead.   |
| `request_id`        | string                                                       | false    |              |                                              |
| `resource_icon`     | string                                                       | false    |              |                                              |
| `resource_id`       | string                                                       | false    |              |                                              |
| `resource_link`     | string                                                       | false    |              |                                              |
| `resource_target`   | string                                                       | false    |              | Resource target is the name of the resource. |
| `resource_type`     | [optimus-ide-collabsdk.ResourceType](#optimus-ide-collabsdkresourcetype)               | false    |              |                                              |
| `status_code`       | integer                                                      | false    |              |                                              |
| `time`              | string                                                       | false    |              |                                              |
| `user`              | [optimus-ide-collabsdk.User](#optimus-ide-collabsdkuser)                               | false    |              |                                              |
| `user_agent`        | string                                                       | false    |              |                                              |

## optimus-ide-collabsdk.AuditLogResponse

```json
{
  "audit_logs": [
    {
      "action": "create",
      "additional_fields": {},
      "description": "string",
      "diff": {
        "property1": {
          "new": null,
          "old": null,
          "secret": true
        },
        "property2": {
          "new": null,
          "old": null,
          "secret": true
        }
      },
      "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
      "ip": "string",
      "is_deleted": true,
      "organization": {
        "display_name": "string",
        "icon": "string",
        "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
        "name": "string"
      },
      "organization_id": "7c60d51f-b44e-4682-87d6-449835ea4de6",
      "request_id": "266ea41d-adf5-480b-af50-15b940c2b846",
      "resource_icon": "string",
      "resource_id": "4d5215ed-38bb-48ed-879a-fdb9ca58522f",
      "resource_link": "string",
      "resource_target": "string",
      "resource_type": "template",
      "status_code": 0,
      "time": "2019-08-24T14:15:22Z",
      "user": {
        "avatar_url": "http://example.com",
        "created_at": "2019-08-24T14:15:22Z",
        "email": "user@example.com",
        "has_ai_seat": true,
        "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
        "is_service_account": true,
        "last_seen_at": "2019-08-24T14:15:22Z",
        "login_type": "",
        "name": "string",
        "organization_ids": [
          "497f6eca-6276-4993-bfeb-53cbbbba6f08"
        ],
        "roles": [
          {
            "display_name": "string",
            "name": "string",
            "organization_id": "string"
          }
        ],
        "status": "active",
        "theme_preference": "string",
        "updated_at": "2019-08-24T14:15:22Z",
        "username": "string"
      },
      "user_agent": "string"
    }
  ],
  "count": 0,
  "count_cap": 0
}
```

### Properties

| Name         | Type                                            | Required | Restrictions | Description |
|--------------|-------------------------------------------------|----------|--------------|-------------|
| `audit_logs` | array of [optimus-ide-collabsdk.AuditLog](#optimus-ide-collabsdkauditlog) | false    |              |             |
| `count`      | integer                                         | false    |              |             |
| `count_cap`  | integer                                         | false    |              |             |

## optimus-ide-collabsdk.AuthMethod

```json
{
  "enabled": true
}
```

### Properties

| Name      | Type    | Required | Restrictions | Description |
|-----------|---------|----------|--------------|-------------|
| `enabled` | boolean | false    |              |             |

## optimus-ide-collabsdk.AuthMethods

```json
{
  "github": {
    "default_provider_configured": true,
    "enabled": true
  },
  "oidc": {
    "enabled": true,
    "iconUrl": "string",
    "signInText": "string"
  },
  "password": {
    "enabled": true
  },
  "terms_of_service_url": "string"
}
```

### Properties

| Name                   | Type                                                   | Required | Restrictions | Description |
|------------------------|--------------------------------------------------------|----------|--------------|-------------|
| `github`               | [optimus-ide-collabsdk.GithubAuthMethod](#optimus-ide-collabsdkgithubauthmethod) | false    |              |             |
| `oidc`                 | [optimus-ide-collabsdk.OIDCAuthMethod](#optimus-ide-collabsdkoidcauthmethod)     | false    |              |             |
| `password`             | [optimus-ide-collabsdk.AuthMethod](#optimus-ide-collabsdkauthmethod)             | false    |              |             |
| `terms_of_service_url` | string                                                 | false    |              |             |

## optimus-ide-collabsdk.AuthorizationCheck

```json
{
  "action": "create",
  "object": {
    "any_org": true,
    "organization_id": "string",
    "owner_id": "string",
    "resource_id": "string",
    "resource_type": "*"
  }
}
```

AuthorizationCheck is used to check if the currently authenticated user (or the specified user) can do a given action to a given set of objects.

### Properties

| Name     | Type                                                         | Required | Restrictions | Description                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                           |
|----------|--------------------------------------------------------------|----------|--------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `action` | [optimus-ide-collabsdk.RBACAction](#optimus-ide-collabsdkrbacaction)                   | false    |              |                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                       |
| `object` | [optimus-ide-collabsdk.AuthorizationObject](#optimus-ide-collabsdkauthorizationobject) | false    |              | Object can represent a "set" of objects, such as: all workspaces in an organization, all workspaces owned by me, and all workspaces across the entire product. When defining an object, use the most specific language when possible to produce the smallest set. Meaning to set as many fields on 'Object' as you can. Example, if you want to check if you can update all workspaces owned by 'me', try to also add an 'OrganizationID' to the settings. Omitting the 'OrganizationID' could produce the incorrect value, as workspaces have both `user` and `organization` owners. |

#### Enumerated Values

| Property | Value(s)                             |
|----------|--------------------------------------|
| `action` | `create`, `delete`, `read`, `update` |

## optimus-ide-collabsdk.AuthorizationObject

```json
{
  "any_org": true,
  "organization_id": "string",
  "owner_id": "string",
  "resource_id": "string",
  "resource_type": "*"
}
```

AuthorizationObject can represent a "set" of objects, such as: all workspaces in an organization, all workspaces owned by me, all workspaces across the entire product.

### Properties

| Name              | Type                                           | Required | Restrictions | Description                                                                                                                                                                                                                                                                                                                                                          |
|-------------------|------------------------------------------------|----------|--------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `any_org`         | boolean                                        | false    |              | Any org (optional) will disregard the org_owner when checking for permissions. This cannot be set to true if the OrganizationID is set.                                                                                                                                                                                                                              |
| `organization_id` | string                                         | false    |              | Organization ID (optional) adds the set constraint to all resources owned by a given organization.                                                                                                                                                                                                                                                                   |
| `owner_id`        | string                                         | false    |              | Owner ID (optional) adds the set constraint to all resources owned by a given user.                                                                                                                                                                                                                                                                                  |
| `resource_id`     | string                                         | false    |              | Resource ID (optional) reduces the set to a singular resource. This assigns a resource ID to the resource type, eg: a single workspace. The rbac library will not fetch the resource from the database, so if you are using this option, you should also set the owner ID and organization ID if possible. Be as specific as possible using all the fields relevant. |
| `resource_type`   | [optimus-ide-collabsdk.RBACResource](#optimus-ide-collabsdkrbacresource) | false    |              | Resource type is the name of the resource. `./optimus-ide-collabd/rbac/object.go` has the list of valid resource types.                                                                                                                                                                                                                                                           |

## optimus-ide-collabsdk.AuthorizationRequest

```json
{
  "checks": {
    "property1": {
      "action": "create",
      "object": {
        "any_org": true,
        "organization_id": "string",
        "owner_id": "string",
        "resource_id": "string",
        "resource_type": "*"
      }
    },
    "property2": {
      "action": "create",
      "object": {
        "any_org": true,
        "organization_id": "string",
        "owner_id": "string",
        "resource_id": "string",
        "resource_type": "*"
      }
    }
  }
}
```

### Properties

| Name               | Type                                                       | Required | Restrictions | Description                                                                                                                                                                                                                                                                      |
|--------------------|------------------------------------------------------------|----------|--------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `checks`           | object                                                     | false    |              | Checks is a map keyed with an arbitrary string to a permission check. The key can be any string that is helpful to the caller, and allows multiple permission checks to be run in a single request. The key ensures that each permission check has the same key in the response. |
| » `[any property]` | [optimus-ide-collabsdk.AuthorizationCheck](#optimus-ide-collabsdkauthorizationcheck) | false    |              | It is used to check if the currently authenticated user (or the specified user) can do a given action to a given set of objects.                                                                                                                                                 |

## optimus-ide-collabsdk.AuthorizationResponse

```json
{
  "property1": true,
  "property2": true
}
```

### Properties

| Name             | Type    | Required | Restrictions | Description |
|------------------|---------|----------|--------------|-------------|
| `[any property]` | boolean | false    |              |             |

## optimus-ide-collabsdk.AutomaticUpdates

```json
"always"
```

### Properties

#### Enumerated Values

| Value(s)          |
|-------------------|
| `always`, `never` |

## optimus-ide-collabsdk.BannerConfig

```json
{
  "background_color": "string",
  "enabled": true,
  "message": "string"
}
```

### Properties

| Name               | Type    | Required | Restrictions | Description |
|--------------------|---------|----------|--------------|-------------|
| `background_color` | string  | false    |              |             |
| `enabled`          | boolean | false    |              |             |
| `message`          | string  | false    |              |             |

## optimus-ide-collabsdk.BuildInfoResponse

```json
{
  "agent_api_version": "string",
  "dashboard_url": "string",
  "deployment_id": "string",
  "external_url": "string",
  "provisioner_api_version": "string",
  "telemetry": true,
  "upgrade_message": "string",
  "version": "string",
  "webpush_public_key": "string",
  "workspace_proxy": true
}
```

### Properties

| Name                      | Type    | Required | Restrictions | Description                                                                                                                                                         |
|---------------------------|---------|----------|--------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `agent_api_version`       | string  | false    |              | Agent api version is the current version of the Agent API (back versions MAY still be supported).                                                                   |
| `dashboard_url`           | string  | false    |              | Dashboard URL is the URL to hit the deployment's dashboard. For external workspace proxies, this is the optimus-ide-collabd they are connected to.                               |
| `deployment_id`           | string  | false    |              | Deployment ID is the unique identifier for this deployment.                                                                                                         |
| `external_url`            | string  | false    |              | External URL references the current Optimus-IDE-Collab version. For production builds, this will link directly to a release. For development builds, this will link to a commit. |
| `provisioner_api_version` | string  | false    |              | Provisioner api version is the current version of the Provisioner API                                                                                               |
| `telemetry`               | boolean | false    |              | Telemetry is a boolean that indicates whether telemetry is enabled.                                                                                                 |
| `upgrade_message`         | string  | false    |              | Upgrade message is the message displayed to users when an outdated client is detected.                                                                              |
| `version`                 | string  | false    |              | Version returns the semantic version of the build.                                                                                                                  |
| `webpush_public_key`      | string  | false    |              | Webpush public key is the public key for push notifications via Web Push.                                                                                           |
| `workspace_proxy`         | boolean | false    |              |                                                                                                                                                                     |

## optimus-ide-collabsdk.BuildReason

```json
"initiator"
```

### Properties

#### Enumerated Values

| Value(s)                                                                                                                                                                                   |
|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `autostart`, `autostop`, `cli`, `dashboard`, `dormancy`, `initiator`, `jetbrains_connection`, `ssh_connection`, `task_auto_pause`, `task_manual_pause`, `task_resume`, `vscode_connection` |

## optimus-ide-collabsdk.CORSBehavior

```json
"simple"
```

### Properties

#### Enumerated Values

| Value(s)             |
|----------------------|
| `passthru`, `simple` |

## optimus-ide-collabsdk.ChangePasswordWithOneTimePasscodeRequest

```json
{
  "email": "user@example.com",
  "one_time_passcode": "string",
  "password": "string"
}
```

### Properties

| Name                | Type   | Required | Restrictions | Description |
|---------------------|--------|----------|--------------|-------------|
| `email`             | string | true     |              |             |
| `one_time_passcode` | string | true     |              |             |
| `password`          | string | true     |              |             |

## optimus-ide-collabsdk.Chat

```json
{
  "agent_id": "2b1e3b65-2c04-4fa2-a2d7-467901e98978",
  "archived": true,
  "build_id": "bfb1f3fa-bf7b-43a5-9e0b-26cc050e44cb",
  "children": [
    {
      "agent_id": "2b1e3b65-2c04-4fa2-a2d7-467901e98978",
      "archived": true,
      "build_id": "bfb1f3fa-bf7b-43a5-9e0b-26cc050e44cb",
      "children": [],
      "client_type": "ui",
      "context": {
        "dirty": true,
        "dirty_since": "2019-08-24T14:15:22Z",
        "error": "string",
        "resources": [
          {
            "error": "string",
            "kind": "instruction_file",
            "size_bytes": 0,
            "skill_description": "string",
            "skill_name": "string",
            "source": "string",
            "status": "ok",
            "tools": [
              {
                "description": "string",
                "name": "string"
              }
            ]
          }
        ]
      },
      "created_at": "2019-08-24T14:15:22Z",
      "diff_status": {
        "additions": 0,
        "approved": true,
        "author_avatar_url": "string",
        "author_login": "string",
        "base_branch": "string",
        "changed_files": 0,
        "changes_requested": true,
        "chat_id": "efc9fe20-a1e5-4a8c-9c48-f1b30c1e4f86",
        "commits": 0,
        "deletions": 0,
        "head_branch": "string",
        "pr_number": 0,
        "pull_request_draft": true,
        "pull_request_state": "string",
        "pull_request_title": "string",
        "refreshed_at": "2019-08-24T14:15:22Z",
        "reviewer_count": 0,
        "stale_at": "2019-08-24T14:15:22Z",
        "url": "string"
      },
      "files": [
        {
          "created_at": "2019-08-24T14:15:22Z",
          "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
          "mime_type": "string",
          "name": "string",
          "organization_id": "7c60d51f-b44e-4682-87d6-449835ea4de6",
          "owner_id": "8826ee2e-7933-4665-aef2-2393f84a0d05"
        }
      ],
      "has_unread": true,
      "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
      "labels": {
        "property1": "string",
        "property2": "string"
      },
      "last_error": {
        "detail": "string",
        "kind": "generic",
        "message": "string",
        "provider": "string",
        "retryable": true,
        "status_code": 0
      },
      "last_model_config_id": "30ebb95f-c255-4759-9429-89aa4ec1554c",
      "last_reasoning_effort": "string",
      "last_turn_summary": "string",
      "mcp_server_ids": [
        "497f6eca-6276-4993-bfeb-53cbbbba6f08"
      ],
      "organization_id": "7c60d51f-b44e-4682-87d6-449835ea4de6",
      "owner_id": "8826ee2e-7933-4665-aef2-2393f84a0d05",
      "owner_name": "string",
      "owner_username": "string",
      "parent_chat_id": "c3609ee6-3b11-4a93-b9ae-e4fabcc99359",
      "pin_order": 0,
      "plan_mode": "plan",
      "root_chat_id": "2898031c-fdce-4e3e-8c53-4481dd42fcd7",
      "shared": true,
      "status": "waiting",
      "summary": "string",
      "title": "string",
      "updated_at": "2019-08-24T14:15:22Z",
      "warnings": [
        "string"
      ],
      "workspace_id": "0967198e-ec7b-4c6b-b4d3-f71244cadbe9"
    }
  ],
  "client_type": "ui",
  "context": {
    "dirty": true,
    "dirty_since": "2019-08-24T14:15:22Z",
    "error": "string",
    "resources": [
      {
        "error": "string",
        "kind": "instruction_file",
        "size_bytes": 0,
        "skill_description": "string",
        "skill_name": "string",
        "source": "string",
        "status": "ok",
        "tools": [
          {
            "description": "string",
            "name": "string"
          }
        ]
      }
    ]
  },
  "created_at": "2019-08-24T14:15:22Z",
  "diff_status": {
    "additions": 0,
    "approved": true,
    "author_avatar_url": "string",
    "author_login": "string",
    "base_branch": "string",
    "changed_files": 0,
    "changes_requested": true,
    "chat_id": "efc9fe20-a1e5-4a8c-9c48-f1b30c1e4f86",
    "commits": 0,
    "deletions": 0,
    "head_branch": "string",
    "pr_number": 0,
    "pull_request_draft": true,
    "pull_request_state": "string",
    "pull_request_title": "string",
    "refreshed_at": "2019-08-24T14:15:22Z",
    "reviewer_count": 0,
    "stale_at": "2019-08-24T14:15:22Z",
    "url": "string"
  },
  "files": [
    {
      "created_at": "2019-08-24T14:15:22Z",
      "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
      "mime_type": "string",
      "name": "string",
      "organization_id": "7c60d51f-b44e-4682-87d6-449835ea4de6",
      "owner_id": "8826ee2e-7933-4665-aef2-2393f84a0d05"
    }
  ],
  "has_unread": true,
  "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
  "labels": {
    "property1": "string",
    "property2": "string"
  },
  "last_error": {
    "detail": "string",
    "kind": "generic",
    "message": "string",
    "provider": "string",
    "retryable": true,
    "status_code": 0
  },
  "last_model_config_id": "30ebb95f-c255-4759-9429-89aa4ec1554c",
  "last_reasoning_effort": "string",
  "last_turn_summary": "string",
  "mcp_server_ids": [
    "497f6eca-6276-4993-bfeb-53cbbbba6f08"
  ],
  "organization_id": "7c60d51f-b44e-4682-87d6-449835ea4de6",
  "owner_id": "8826ee2e-7933-4665-aef2-2393f84a0d05",
  "owner_name": "string",
  "owner_username": "string",
  "parent_chat_id": "c3609ee6-3b11-4a93-b9ae-e4fabcc99359",
  "pin_order": 0,
  "plan_mode": "plan",
  "root_chat_id": "2898031c-fdce-4e3e-8c53-4481dd42fcd7",
  "shared": true,
  "status": "waiting",
  "summary": "string",
  "title": "string",
  "updated_at": "2019-08-24T14:15:22Z",
  "warnings": [
    "string"
  ],
  "workspace_id": "0967198e-ec7b-4c6b-b4d3-f71244cadbe9"
}
```

### Properties

| Name                    | Type                                                            | Required | Restrictions | Description                                                                                                                                                                                                                                                                |
|-------------------------|-----------------------------------------------------------------|----------|--------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `agent_id`              | string                                                          | false    |              |                                                                                                                                                                                                                                                                            |
| `archived`              | boolean                                                         | false    |              |                                                                                                                                                                                                                                                                            |
| `build_id`              | string                                                          | false    |              |                                                                                                                                                                                                                                                                            |
| `children`              | array of [optimus-ide-collabsdk.Chat](#optimus-ide-collabsdkchat)                         | false    |              | Children holds child (subagent) chats nested under this root chat. Always initialized to an empty slice so the JSON field is present as []. Child chats cannot create their own subagents, so nesting depth is capped at 1 and this slice is always empty for child chats. |
| `client_type`           | [optimus-ide-collabsdk.ChatClientType](#optimus-ide-collabsdkchatclienttype)              | false    |              |                                                                                                                                                                                                                                                                            |
| `context`               | [optimus-ide-collabsdk.ChatContext](#optimus-ide-collabsdkchatcontext)                    | false    |              | Context reports the chat's pinned workspace-context state and whether it has drifted from the agent's latest pushed snapshot. Nil when the chat has no pinned context yet.                                                                                                 |
| `created_at`            | string                                                          | false    |              |                                                                                                                                                                                                                                                                            |
| `diff_status`           | [optimus-ide-collabsdk.ChatDiffStatus](#optimus-ide-collabsdkchatdiffstatus)              | false    |              |                                                                                                                                                                                                                                                                            |
| `files`                 | array of [optimus-ide-collabsdk.ChatFileMetadata](#optimus-ide-collabsdkchatfilemetadata) | false    |              |                                                                                                                                                                                                                                                                            |
| `has_unread`            | boolean                                                         | false    |              | Has unread is true when assistant messages exist beyond the owner's read cursor, which updates on stream connect and disconnect.                                                                                                                                           |
| `id`                    | string                                                          | false    |              |                                                                                                                                                                                                                                                                            |
| `labels`                | object                                                          | false    |              |                                                                                                                                                                                                                                                                            |
| » `[any property]`      | string                                                          | false    |              |                                                                                                                                                                                                                                                                            |
| `last_error`            | [optimus-ide-collabsdk.ChatError](#optimus-ide-collabsdkchaterror)                        | false    |              |                                                                                                                                                                                                                                                                            |
| `last_model_config_id`  | string                                                          | false    |              |                                                                                                                                                                                                                                                                            |
| `last_reasoning_effort` | string                                                          | false    |              |                                                                                                                                                                                                                                                                            |
| `last_turn_summary`     | string                                                          | false    |              |                                                                                                                                                                                                                                                                            |
| `mcp_server_ids`        | array of string                                                 | false    |              |                                                                                                                                                                                                                                                                            |
| `organization_id`       | string                                                          | false    |              |                                                                                                                                                                                                                                                                            |
| `owner_id`              | string                                                          | false    |              |                                                                                                                                                                                                                                                                            |
| `owner_name`            | string                                                          | false    |              |                                                                                                                                                                                                                                                                            |
| `owner_username`        | string                                                          | false    |              |                                                                                                                                                                                                                                                                            |
| `parent_chat_id`        | string                                                          | false    |              |                                                                                                                                                                                                                                                                            |
| `pin_order`             | integer                                                         | false    |              |                                                                                                                                                                                                                                                                            |
| `plan_mode`             | [optimus-ide-collabsdk.ChatPlanMode](#optimus-ide-collabsdkchatplanmode)                  | false    |              |                                                                                                                                                                                                                                                                            |
| `root_chat_id`          | string                                                          | false    |              |                                                                                                                                                                                                                                                                            |
| `shared`                | boolean                                                         | false    |              | Shared is true when this chat's root chat has explicit user or group ACL entries.                                                                                                                                                                                          |
| `status`                | [optimus-ide-collabsdk.ChatStatus](#optimus-ide-collabsdkchatstatus)                      | false    |              |                                                                                                                                                                                                                                                                            |
| `summary`               | string                                                          | false    |              | Summary is the persisted whole-chat summary, generated in the background. It is nil until the first summary has been produced.                                                                                                                                             |
| `title`                 | string                                                          | false    |              |                                                                                                                                                                                                                                                                            |
| `updated_at`            | string                                                          | false    |              |                                                                                                                                                                                                                                                                            |
| `warnings`              | array of string                                                 | false    |              |                                                                                                                                                                                                                                                                            |
| `workspace_id`          | string                                                          | false    |              |                                                                                                                                                                                                                                                                            |

## optimus-ide-collabsdk.ChatACL

```json
{
  "groups": [
    {
      "avatar_url": "http://example.com",
      "display_name": "string",
      "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
      "members": [
        {
          "avatar_url": "http://example.com",
          "created_at": "2019-08-24T14:15:22Z",
          "email": "user@example.com",
          "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
          "is_service_account": true,
          "last_seen_at": "2019-08-24T14:15:22Z",
          "login_type": "",
          "name": "string",
          "status": "active",
          "theme_preference": "string",
          "updated_at": "2019-08-24T14:15:22Z",
          "username": "string"
        }
      ],
      "name": "string",
      "organization_display_name": "string",
      "organization_id": "7c60d51f-b44e-4682-87d6-449835ea4de6",
      "organization_name": "string",
      "quota_allowance": 0,
      "role": "read",
      "source": "user",
      "total_member_count": 0
    }
  ],
  "users": [
    {
      "avatar_url": "http://example.com",
      "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
      "name": "string",
      "role": "read",
      "username": "string"
    }
  ]
}
```

### Properties

| Name     | Type                                              | Required | Restrictions | Description |
|----------|---------------------------------------------------|----------|--------------|-------------|
| `groups` | array of [optimus-ide-collabsdk.ChatGroup](#optimus-ide-collabsdkchatgroup) | false    |              |             |
| `users`  | array of [optimus-ide-collabsdk.ChatUser](#optimus-ide-collabsdkchatuser)   | false    |              |             |

## optimus-ide-collabsdk.ChatBusyBehavior

```json
"queue"
```

### Properties

#### Enumerated Values

| Value(s)             |
|----------------------|
| `interrupt`, `queue` |

## optimus-ide-collabsdk.ChatClientType

```json
"ui"
```

### Properties

#### Enumerated Values

| Value(s)    |
|-------------|
| `api`, `ui` |

## optimus-ide-collabsdk.ChatConfig

```json
{
  "acquire_batch_size": 0,
  "debug_logging_enabled": true
}
```

### Properties

| Name                    | Type    | Required | Restrictions | Description |
|-------------------------|---------|----------|--------------|-------------|
| `acquire_batch_size`    | integer | false    |              |             |
| `debug_logging_enabled` | boolean | false    |              |             |

## optimus-ide-collabsdk.ChatContext

```json
{
  "dirty": true,
  "dirty_since": "2019-08-24T14:15:22Z",
  "error": "string",
  "resources": [
    {
      "error": "string",
      "kind": "instruction_file",
      "size_bytes": 0,
      "skill_description": "string",
      "skill_name": "string",
      "source": "string",
      "status": "ok",
      "tools": [
        {
          "description": "string",
          "name": "string"
        }
      ]
    }
  ]
}
```

### Properties

| Name          | Type                                                                  | Required | Restrictions | Description                                                                                                                                                                                                                                |
|---------------|-----------------------------------------------------------------------|----------|--------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `dirty`       | boolean                                                               | false    |              | Dirty is true when the agent's latest snapshot hash differs from the chat's pinned hash.                                                                                                                                                   |
| `dirty_since` | string                                                                | false    |              | Dirty since is when drift was first detected; nil when not dirty.                                                                                                                                                                          |
| `error`       | string                                                                | false    |              | Error is the snapshot-level error copied from the pinned snapshot (empty when healthy).                                                                                                                                                    |
| `resources`   | array of [optimus-ide-collabsdk.ChatContextResource](#optimus-ide-collabsdkchatcontextresource) | false    |              | Resources is the chat's pinned context (instruction files and skills) the prompt is built from, metadata only (no bodies). It is populated only on the single-chat GET response; list and watch payloads leave it nil to stay lightweight. |

## optimus-ide-collabsdk.ChatContextResource

```json
{
  "error": "string",
  "kind": "instruction_file",
  "size_bytes": 0,
  "skill_description": "string",
  "skill_name": "string",
  "source": "string",
  "status": "ok",
  "tools": [
    {
      "description": "string",
      "name": "string"
    }
  ]
}
```

### Properties

| Name                | Type                                                                     | Required | Restrictions | Description                                                                                                                                                                                                                                                                |
|---------------------|--------------------------------------------------------------------------|----------|--------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `error`             | string                                                                   | false    |              | Error explains a non-ok Status; empty when healthy. May also carry a non-fatal warning when Status is ok.                                                                                                                                                                  |
| `kind`              | [optimus-ide-collabsdk.ChatContextResourceKind](#optimus-ide-collabsdkchatcontextresourcekind)     | false    |              |                                                                                                                                                                                                                                                                            |
| `size_bytes`        | integer                                                                  | false    |              | Size bytes is the original payload size in bytes.                                                                                                                                                                                                                          |
| `skill_description` | string                                                                   | false    |              |                                                                                                                                                                                                                                                                            |
| `skill_name`        | string                                                                   | false    |              | Skill name and SkillDescription are populated only for skill kinds.                                                                                                                                                                                                        |
| `source`            | string                                                                   | false    |              | Source is the resource locator: the canonical file path for an instruction file, the skill directory for a skill, the file path for an MCP config, or the server name for an MCP server.                                                                                   |
| `status`            | [optimus-ide-collabsdk.ChatContextResourceStatus](#optimus-ide-collabsdkchatcontextresourcestatus) | false    |              | Status is the resource's health. Non-ok resources (invalid, unreadable, oversize, excluded) are still reported so the UI can surface why a resource was dropped from the prompt instead of silently omitting it; their body-specific fields (skill name, tools) are empty. |
| `tools`             | array of [optimus-ide-collabsdk.ChatContextTool](#optimus-ide-collabsdkchatcontexttool)            | false    |              | Tools lists the tools exposed by an MCP server. Populated only for the mcp_server kind; nil otherwise.                                                                                                                                                                     |

## optimus-ide-collabsdk.ChatContextResourceKind

```json
"instruction_file"
```

### Properties

#### Enumerated Values

| Value(s)                                                |
|---------------------------------------------------------|
| `instruction_file`, `mcp_config`, `mcp_server`, `skill` |

## optimus-ide-collabsdk.ChatContextResourceStatus

```json
"ok"
```

### Properties

#### Enumerated Values

| Value(s)                                              |
|-------------------------------------------------------|
| `excluded`, `invalid`, `ok`, `oversize`, `unreadable` |

## optimus-ide-collabsdk.ChatContextTool

```json
{
  "description": "string",
  "name": "string"
}
```

### Properties

| Name          | Type   | Required | Restrictions | Description                                                                                                       |
|---------------|--------|----------|--------------|-------------------------------------------------------------------------------------------------------------------|
| `description` | string | false    |              | Description is the tool's human-readable summary; may be empty.                                                   |
| `name`        | string | false    |              | Name is the tool name with the "<server>__" prefix the agent adds stripped, so it reads as the server exposes it. |

## optimus-ide-collabsdk.ChatCost

```json
{
  "chat_id": "efc9fe20-a1e5-4a8c-9c48-f1b30c1e4f86",
  "priced_message_count": 0,
  "total_cost_micros": 0,
  "unpriced_messages_having_usage_count": 0
}
```

### Properties

| Name                                   | Type    | Required | Restrictions | Description |
|----------------------------------------|---------|----------|--------------|-------------|
| `chat_id`                              | string  | false    |              |             |
| `priced_message_count`                 | integer | false    |              |             |
| `total_cost_micros`                    | integer | false    |              |             |
| `unpriced_messages_having_usage_count` | integer | false    |              |             |

## optimus-ide-collabsdk.ChatDiffContents

```json
{
  "branch": "string",
  "chat_id": "efc9fe20-a1e5-4a8c-9c48-f1b30c1e4f86",
  "diff": "string",
  "provider": "string",
  "pull_request_url": "string",
  "remote_origin": "string"
}
```

### Properties

| Name               | Type   | Required | Restrictions | Description |
|--------------------|--------|----------|--------------|-------------|
| `branch`           | string | false    |              |             |
| `chat_id`          | string | false    |              |             |
| `diff`             | string | false    |              |             |
| `provider`         | string | false    |              |             |
| `pull_request_url` | string | false    |              |             |
| `remote_origin`    | string | false    |              |             |

## optimus-ide-collabsdk.ChatDiffStatus

```json
{
  "additions": 0,
  "approved": true,
  "author_avatar_url": "string",
  "author_login": "string",
  "base_branch": "string",
  "changed_files": 0,
  "changes_requested": true,
  "chat_id": "efc9fe20-a1e5-4a8c-9c48-f1b30c1e4f86",
  "commits": 0,
  "deletions": 0,
  "head_branch": "string",
  "pr_number": 0,
  "pull_request_draft": true,
  "pull_request_state": "string",
  "pull_request_title": "string",
  "refreshed_at": "2019-08-24T14:15:22Z",
  "reviewer_count": 0,
  "stale_at": "2019-08-24T14:15:22Z",
  "url": "string"
}
```

### Properties

| Name                 | Type    | Required | Restrictions | Description |
|----------------------|---------|----------|--------------|-------------|
| `additions`          | integer | false    |              |             |
| `approved`           | boolean | false    |              |             |
| `author_avatar_url`  | string  | false    |              |             |
| `author_login`       | string  | false    |              |             |
| `base_branch`        | string  | false    |              |             |
| `changed_files`      | integer | false    |              |             |
| `changes_requested`  | boolean | false    |              |             |
| `chat_id`            | string  | false    |              |             |
| `commits`            | integer | false    |              |             |
| `deletions`          | integer | false    |              |             |
| `head_branch`        | string  | false    |              |             |
| `pr_number`          | integer | false    |              |             |
| `pull_request_draft` | boolean | false    |              |             |
| `pull_request_state` | string  | false    |              |             |
| `pull_request_title` | string  | false    |              |             |
| `refreshed_at`       | string  | false    |              |             |
| `reviewer_count`     | integer | false    |              |             |
| `stale_at`           | string  | false    |              |             |
| `url`                | string  | false    |              |             |

## optimus-ide-collabsdk.ChatError

```json
{
  "detail": "string",
  "kind": "generic",
  "message": "string",
  "provider": "string",
  "retryable": true,
  "status_code": 0
}
```

### Properties

| Name          | Type                                             | Required | Restrictions | Description                                                                                               |
|---------------|--------------------------------------------------|----------|--------------|-----------------------------------------------------------------------------------------------------------|
| `detail`      | string                                           | false    |              | Detail is optional provider-specific context shown alongside the normalized error message when available. |
| `kind`        | [optimus-ide-collabsdk.ChatErrorKind](#optimus-ide-collabsdkchaterrorkind) | false    |              | Kind classifies the error for consistent client rendering.                                                |
| `message`     | string                                           | false    |              | Message is the normalized, user-facing error message.                                                     |
| `provider`    | string                                           | false    |              | Provider identifies the upstream model provider when known.                                               |
| `retryable`   | boolean                                          | false    |              | Retryable reports whether the underlying error is transient.                                              |
| `status_code` | integer                                          | false    |              | Status code is the best-effort upstream HTTP status code.                                                 |

## optimus-ide-collabsdk.ChatErrorKind

```json
"generic"
```

### Properties

#### Enumerated Values

| Value(s)                                                                                                                                                          |
|-------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `auth`, `config`, `content_filter`, `generic`, `missing_key`, `overloaded`, `provider_disabled`, `rate_limit`, `stream_silence_timeout`, `timeout`, `usage_limit` |

## optimus-ide-collabsdk.ChatFileMetadata

```json
{
  "created_at": "2019-08-24T14:15:22Z",
  "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
  "mime_type": "string",
  "name": "string",
  "organization_id": "7c60d51f-b44e-4682-87d6-449835ea4de6",
  "owner_id": "8826ee2e-7933-4665-aef2-2393f84a0d05"
}
```

### Properties

| Name              | Type   | Required | Restrictions | Description |
|-------------------|--------|----------|--------------|-------------|
| `created_at`      | string | false    |              |             |
| `id`              | string | false    |              |             |
| `mime_type`       | string | false    |              |             |
| `name`            | string | false    |              |             |
| `organization_id` | string | false    |              |             |
| `owner_id`        | string | false    |              |             |

## optimus-ide-collabsdk.ChatGroup

```json
{
  "avatar_url": "http://example.com",
  "display_name": "string",
  "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
  "members": [
    {
      "avatar_url": "http://example.com",
      "created_at": "2019-08-24T14:15:22Z",
      "email": "user@example.com",
      "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
      "is_service_account": true,
      "last_seen_at": "2019-08-24T14:15:22Z",
      "login_type": "",
      "name": "string",
      "status": "active",
      "theme_preference": "string",
      "updated_at": "2019-08-24T14:15:22Z",
      "username": "string"
    }
  ],
  "name": "string",
  "organization_display_name": "string",
  "organization_id": "7c60d51f-b44e-4682-87d6-449835ea4de6",
  "organization_name": "string",
  "quota_allowance": 0,
  "role": "read",
  "source": "user",
  "total_member_count": 0
}
```

### Properties

| Name                        | Type                                                  | Required | Restrictions | Description                                                                                                                                                           |
|-----------------------------|-------------------------------------------------------|----------|--------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `avatar_url`                | string                                                | false    |              |                                                                                                                                                                       |
| `display_name`              | string                                                | false    |              |                                                                                                                                                                       |
| `id`                        | string                                                | false    |              |                                                                                                                                                                       |
| `members`                   | array of [optimus-ide-collabsdk.ReducedUser](#optimus-ide-collabsdkreduceduser) | false    |              |                                                                                                                                                                       |
| `name`                      | string                                                | false    |              |                                                                                                                                                                       |
| `organization_display_name` | string                                                | false    |              |                                                                                                                                                                       |
| `organization_id`           | string                                                | false    |              |                                                                                                                                                                       |
| `organization_name`         | string                                                | false    |              |                                                                                                                                                                       |
| `quota_allowance`           | integer                                               | false    |              |                                                                                                                                                                       |
| `role`                      | [optimus-ide-collabsdk.ChatRole](#optimus-ide-collabsdkchatrole)                | false    |              |                                                                                                                                                                       |
| `source`                    | [optimus-ide-collabsdk.GroupSource](#optimus-ide-collabsdkgroupsource)          | false    |              |                                                                                                                                                                       |
| `total_member_count`        | integer                                               | false    |              | How many members are in this group. Shows the total count, even if the user is not authorized to read group member details. May be greater than `len(Group.Members)`. |

#### Enumerated Values

| Property | Value(s) |
|----------|----------|
| `role`   | `read`   |

## optimus-ide-collabsdk.ChatInputPart

```json
{
  "content": "string",
  "end_line": 0,
  "file_id": "8a0cfb4f-ddc9-436d-91bb-75133c583767",
  "file_name": "string",
  "start_line": 0,
  "text": "string",
  "type": "text"
}
```

### Properties

| Name         | Type                                                     | Required | Restrictions | Description                                                                    |
|--------------|----------------------------------------------------------|----------|--------------|--------------------------------------------------------------------------------|
| `content`    | string                                                   | false    |              | The code content from the diff that was commented on.                          |
| `end_line`   | integer                                                  | false    |              |                                                                                |
| `file_id`    | string                                                   | false    |              |                                                                                |
| `file_name`  | string                                                   | false    |              | The following fields are only set when Type is ChatInputPartTypeFileReference. |
| `start_line` | integer                                                  | false    |              |                                                                                |
| `text`       | string                                                   | false    |              |                                                                                |
| `type`       | [optimus-ide-collabsdk.ChatInputPartType](#optimus-ide-collabsdkchatinputparttype) | false    |              |                                                                                |

## optimus-ide-collabsdk.ChatInputPartType

```json
"text"
```

### Properties

#### Enumerated Values

| Value(s)                         |
|----------------------------------|
| `file`, `file-reference`, `text` |

## optimus-ide-collabsdk.ChatMessage

```json
{
  "chat_id": "efc9fe20-a1e5-4a8c-9c48-f1b30c1e4f86",
  "content": [
    {
      "args": [
        0
      ],
      "args_delta": "string",
      "completed_at": "2019-08-24T14:15:22Z",
      "content": "string",
      "context_file_agent_id": {
        "uuid": "string",
        "valid": true
      },
      "context_file_content": "string",
      "context_file_directory": "string",
      "context_file_os": "string",
      "context_file_path": "string",
      "context_file_skill_meta_file": "string",
      "context_file_truncated": true,
      "created_at": "2019-08-24T14:15:22Z",
      "data": [
        0
      ],
      "end_line": 0,
      "file_id": {
        "uuid": "string",
        "valid": true
      },
      "file_name": "string",
      "is_error": true,
      "is_media": true,
      "mcp_server_config_id": {
        "uuid": "string",
        "valid": true
      },
      "media_type": "string",
      "name": "string",
      "parsed_commands": [
        [
          "string"
        ]
      ],
      "provider_executed": true,
      "provider_metadata": [
        0
      ],
      "result": [
        0
      ],
      "result_delta": "string",
      "result_reset": true,
      "skill_description": "string",
      "skill_dir": "string",
      "skill_name": "string",
      "source_id": "string",
      "start_line": 0,
      "text": "string",
      "title": "string",
      "tool_call_id": "string",
      "tool_name": "string",
      "type": "text",
      "url": "string"
    }
  ],
  "created_at": "2019-08-24T14:15:22Z",
  "created_by": "ee824cad-d7a6-4f48-87dc-e8461a9201c4",
  "id": 0,
  "model_config_id": "f5fb4d91-62ca-4377-9ee6-5d43ba00d205",
  "role": "system",
  "usage": {
    "cache_creation_tokens": 0,
    "cache_read_tokens": 0,
    "context_limit": 0,
    "input_tokens": 0,
    "output_tokens": 0,
    "reasoning_tokens": 0,
    "total_tokens": 0
  }
}
```

### Properties

| Name              | Type                                                          | Required | Restrictions | Description |
|-------------------|---------------------------------------------------------------|----------|--------------|-------------|
| `chat_id`         | string                                                        | false    |              |             |
| `content`         | array of [optimus-ide-collabsdk.ChatMessagePart](#optimus-ide-collabsdkchatmessagepart) | false    |              |             |
| `created_at`      | string                                                        | false    |              |             |
| `created_by`      | string                                                        | false    |              |             |
| `id`              | integer                                                       | false    |              |             |
| `model_config_id` | string                                                        | false    |              |             |
| `role`            | [optimus-ide-collabsdk.ChatMessageRole](#optimus-ide-collabsdkchatmessagerole)          | false    |              |             |
| `usage`           | [optimus-ide-collabsdk.ChatMessageUsage](#optimus-ide-collabsdkchatmessageusage)        | false    |              |             |

## optimus-ide-collabsdk.ChatMessagePart

```json
{
  "args": [
    0
  ],
  "args_delta": "string",
  "completed_at": "2019-08-24T14:15:22Z",
  "content": "string",
  "context_file_agent_id": {
    "uuid": "string",
    "valid": true
  },
  "context_file_content": "string",
  "context_file_directory": "string",
  "context_file_os": "string",
  "context_file_path": "string",
  "context_file_skill_meta_file": "string",
  "context_file_truncated": true,
  "created_at": "2019-08-24T14:15:22Z",
  "data": [
    0
  ],
  "end_line": 0,
  "file_id": {
    "uuid": "string",
    "valid": true
  },
  "file_name": "string",
  "is_error": true,
  "is_media": true,
  "mcp_server_config_id": {
    "uuid": "string",
    "valid": true
  },
  "media_type": "string",
  "name": "string",
  "parsed_commands": [
    [
      "string"
    ]
  ],
  "provider_executed": true,
  "provider_metadata": [
    0
  ],
  "result": [
    0
  ],
  "result_delta": "string",
  "result_reset": true,
  "skill_description": "string",
  "skill_dir": "string",
  "skill_name": "string",
  "source_id": "string",
  "start_line": 0,
  "text": "string",
  "title": "string",
  "tool_call_id": "string",
  "tool_name": "string",
  "type": "text",
  "url": "string"
}
```

### Properties

| Name                           | Type                                                         | Required | Restrictions | Description                                                                                                                                                                                                                                                                                                                                                                                                |
|--------------------------------|--------------------------------------------------------------|----------|--------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `args`                         | array of integer                                             | false    |              |                                                                                                                                                                                                                                                                                                                                                                                                            |
| `args_delta`                   | string                                                       | false    |              |                                                                                                                                                                                                                                                                                                                                                                                                            |
| `completed_at`                 | string                                                       | false    |              | Completed at is the time a reasoning part finished streaming, so reasoning duration can be computed as completed_at minus created_at. For interrupted reasoning, this is the interruption time. Absent when reasoning timestamp data was not recorded (e.g. messages persisted before this feature was added).                                                                                             |
| `content`                      | string                                                       | false    |              | The code content from the diff that was commented on.                                                                                                                                                                                                                                                                                                                                                      |
| `context_file_agent_id`        | [uuid.NullUUID](#uuidnulluuid)                               | false    |              | Context file agent ID is the workspace agent that provided this context file. Used to detect when the agent changes (e.g. workspace rebuilt) so instruction files can be re-persisted with fresh content.                                                                                                                                                                                                  |
| `context_file_content`         | string                                                       | false    |              | Context file content holds the file content sent to the LLM. Internal only: stripped before API responses to keep payloads small. The backend reads it when building the prompt via partsToMessageParts.                                                                                                                                                                                                   |
| `context_file_directory`       | string                                                       | false    |              | Context file directory is the working directory of the workspace agent. Internal only: same purpose as ContextFileOS.                                                                                                                                                                                                                                                                                      |
| `context_file_os`              | string                                                       | false    |              | Context file os is the operating system of the workspace agent. Internal only: used during prompt expansion so the LLM knows the OS even on turns where InsertSystem is not called.                                                                                                                                                                                                                        |
| `context_file_path`            | string                                                       | false    |              | Context file path is the absolute path of a file loaded into the LLM context (e.g. an AGENTS.md instruction file).                                                                                                                                                                                                                                                                                         |
| `context_file_skill_meta_file` | string                                                       | false    |              | Context file skill meta file is the basename of the skill meta file (e.g. "SKILL.md") at the time of persistence. Internal only: restored on subsequent turns so the read_skill tool uses the correct filename even when the agent configured a non-default value.                                                                                                                                         |
| `context_file_truncated`       | boolean                                                      | false    |              | Context file truncated indicates the file exceeded the 64KiB instruction file limit and was truncated.                                                                                                                                                                                                                                                                                                     |
| `created_at`                   | string                                                       | false    |              | Created at is the timestamp this part carries. The semantics depend on the part type: for tool-call and tool-result parts it is the time the call was emitted or the result was produced (tool duration is the result's created_at minus the call's created_at); for reasoning parts it is the time reasoning started streaming.                                                                           |
| `data`                         | array of integer                                             | false    |              |                                                                                                                                                                                                                                                                                                                                                                                                            |
| `end_line`                     | integer                                                      | false    |              |                                                                                                                                                                                                                                                                                                                                                                                                            |
| `file_id`                      | [uuid.NullUUID](#uuidnulluuid)                               | false    |              |                                                                                                                                                                                                                                                                                                                                                                                                            |
| `file_name`                    | string                                                       | false    |              |                                                                                                                                                                                                                                                                                                                                                                                                            |
| `is_error`                     | boolean                                                      | false    |              |                                                                                                                                                                                                                                                                                                                                                                                                            |
| `is_media`                     | boolean                                                      | false    |              |                                                                                                                                                                                                                                                                                                                                                                                                            |
| `mcp_server_config_id`         | [uuid.NullUUID](#uuidnulluuid)                               | false    |              |                                                                                                                                                                                                                                                                                                                                                                                                            |
| `media_type`                   | string                                                       | false    |              |                                                                                                                                                                                                                                                                                                                                                                                                            |
| `name`                         | string                                                       | false    |              |                                                                                                                                                                                                                                                                                                                                                                                                            |
| `parsed_commands`              | array of array                                               | false    |              | Parsed commands holds parsed programs from an execute tool call's shell command, one entry per simple command in source order. Each entry is [program] or [program, arg] where arg is the first non-flag positional argument. Program names are normalized to their base name (e.g. /usr/bin/go becomes go). Only populated when ToolName is "execute" and the command parses successfully; nil otherwise. |
| `provider_executed`            | boolean                                                      | false    |              | Provider executed indicates the tool call was executed by the provider (e.g. Anthropic computer use).                                                                                                                                                                                                                                                                                                      |
| `provider_metadata`            | array of integer                                             | false    |              | Provider metadata holds provider-specific response metadata (e.g. Anthropic cache control hints) as raw JSON. Internal only: stripped by db2sdk before API responses.                                                                                                                                                                                                                                      |
| `result`                       | array of integer                                             | false    |              |                                                                                                                                                                                                                                                                                                                                                                                                            |
| `result_delta`                 | string                                                       | false    |              |                                                                                                                                                                                                                                                                                                                                                                                                            |
| `result_reset`                 | boolean                                                      | false    |              |                                                                                                                                                                                                                                                                                                                                                                                                            |
| `skill_description`            | string                                                       | false    |              | Skill description is the short description from the skill's SKILL.md frontmatter.                                                                                                                                                                                                                                                                                                                          |
| `skill_dir`                    | string                                                       | false    |              | Skill dir is the absolute path to the skill directory inside the workspace filesystem. Internal only: used by read_skill/read_skill_file tools to locate skill files.                                                                                                                                                                                                                                      |
| `skill_name`                   | string                                                       | false    |              | Skill name is the kebab-case name of a discovered skill from the workspace's .agents/skills/ directory.                                                                                                                                                                                                                                                                                                    |
| `source_id`                    | string                                                       | false    |              |                                                                                                                                                                                                                                                                                                                                                                                                            |
| `start_line`                   | integer                                                      | false    |              |                                                                                                                                                                                                                                                                                                                                                                                                            |
| `text`                         | string                                                       | false    |              |                                                                                                                                                                                                                                                                                                                                                                                                            |
| `title`                        | string                                                       | false    |              |                                                                                                                                                                                                                                                                                                                                                                                                            |
| `tool_call_id`                 | string                                                       | false    |              |                                                                                                                                                                                                                                                                                                                                                                                                            |
| `tool_name`                    | string                                                       | false    |              |                                                                                                                                                                                                                                                                                                                                                                                                            |
| `type`                         | [optimus-ide-collabsdk.ChatMessagePartType](#optimus-ide-collabsdkchatmessageparttype) | false    |              |                                                                                                                                                                                                                                                                                                                                                                                                            |
| `url`                          | string                                                       | false    |              |                                                                                                                                                                                                                                                                                                                                                                                                            |

## optimus-ide-collabsdk.ChatMessagePartType

```json
"text"
```

### Properties

#### Enumerated Values

| Value(s)                                                                                                     |
|--------------------------------------------------------------------------------------------------------------|
| `context-file`, `file`, `file-reference`, `reasoning`, `skill`, `source`, `text`, `tool-call`, `tool-result` |

## optimus-ide-collabsdk.ChatMessageRole

```json
"system"
```

### Properties

#### Enumerated Values

| Value(s)                              |
|---------------------------------------|
| `assistant`, `system`, `tool`, `user` |

## optimus-ide-collabsdk.ChatMessageUsage

```json
{
  "cache_creation_tokens": 0,
  "cache_read_tokens": 0,
  "context_limit": 0,
  "input_tokens": 0,
  "output_tokens": 0,
  "reasoning_tokens": 0,
  "total_tokens": 0
}
```

### Properties

| Name                    | Type    | Required | Restrictions | Description |
|-------------------------|---------|----------|--------------|-------------|
| `cache_creation_tokens` | integer | false    |              |             |
| `cache_read_tokens`     | integer | false    |              |             |
| `context_limit`         | integer | false    |              |             |
| `input_tokens`          | integer | false    |              |             |
| `output_tokens`         | integer | false    |              |             |
| `reasoning_tokens`      | integer | false    |              |             |
| `total_tokens`          | integer | false    |              |             |

## optimus-ide-collabsdk.ChatMessagesResponse

```json
{
  "has_more": true,
  "messages": [
    {
      "chat_id": "efc9fe20-a1e5-4a8c-9c48-f1b30c1e4f86",
      "content": [
        {
          "args": [
            0
          ],
          "args_delta": "string",
          "completed_at": "2019-08-24T14:15:22Z",
          "content": "string",
          "context_file_agent_id": {
            "uuid": "string",
            "valid": true
          },
          "context_file_content": "string",
          "context_file_directory": "string",
          "context_file_os": "string",
          "context_file_path": "string",
          "context_file_skill_meta_file": "string",
          "context_file_truncated": true,
          "created_at": "2019-08-24T14:15:22Z",
          "data": [
            0
          ],
          "end_line": 0,
          "file_id": {
            "uuid": "string",
            "valid": true
          },
          "file_name": "string",
          "is_error": true,
          "is_media": true,
          "mcp_server_config_id": {
            "uuid": "string",
            "valid": true
          },
          "media_type": "string",
          "name": "string",
          "parsed_commands": [
            [
              "string"
            ]
          ],
          "provider_executed": true,
          "provider_metadata": [
            0
          ],
          "result": [
            0
          ],
          "result_delta": "string",
          "result_reset": true,
          "skill_description": "string",
          "skill_dir": "string",
          "skill_name": "string",
          "source_id": "string",
          "start_line": 0,
          "text": "string",
          "title": "string",
          "tool_call_id": "string",
          "tool_name": "string",
          "type": "text",
          "url": "string"
        }
      ],
      "created_at": "2019-08-24T14:15:22Z",
      "created_by": "ee824cad-d7a6-4f48-87dc-e8461a9201c4",
      "id": 0,
      "model_config_id": "f5fb4d91-62ca-4377-9ee6-5d43ba00d205",
      "role": "system",
      "usage": {
        "cache_creation_tokens": 0,
        "cache_read_tokens": 0,
        "context_limit": 0,
        "input_tokens": 0,
        "output_tokens": 0,
        "reasoning_tokens": 0,
        "total_tokens": 0
      }
    }
  ],
  "queued_messages": [
    {
      "chat_id": "efc9fe20-a1e5-4a8c-9c48-f1b30c1e4f86",
      "content": [
        {
          "args": [
            0
          ],
          "args_delta": "string",
          "completed_at": "2019-08-24T14:15:22Z",
          "content": "string",
          "context_file_agent_id": {
            "uuid": "string",
            "valid": true
          },
          "context_file_content": "string",
          "context_file_directory": "string",
          "context_file_os": "string",
          "context_file_path": "string",
          "context_file_skill_meta_file": "string",
          "context_file_truncated": true,
          "created_at": "2019-08-24T14:15:22Z",
          "data": [
            0
          ],
          "end_line": 0,
          "file_id": {
            "uuid": "string",
            "valid": true
          },
          "file_name": "string",
          "is_error": true,
          "is_media": true,
          "mcp_server_config_id": {
            "uuid": "string",
            "valid": true
          },
          "media_type": "string",
          "name": "string",
          "parsed_commands": [
            [
              "string"
            ]
          ],
          "provider_executed": true,
          "provider_metadata": [
            0
          ],
          "result": [
            0
          ],
          "result_delta": "string",
          "result_reset": true,
          "skill_description": "string",
          "skill_dir": "string",
          "skill_name": "string",
          "source_id": "string",
          "start_line": 0,
          "text": "string",
          "title": "string",
          "tool_call_id": "string",
          "tool_name": "string",
          "type": "text",
          "url": "string"
        }
      ],
      "created_at": "2019-08-24T14:15:22Z",
      "id": 0,
      "model_config_id": "f5fb4d91-62ca-4377-9ee6-5d43ba00d205"
    }
  ]
}
```

### Properties

| Name              | Type                                                              | Required | Restrictions | Description |
|-------------------|-------------------------------------------------------------------|----------|--------------|-------------|
| `has_more`        | boolean                                                           | false    |              |             |
| `messages`        | array of [optimus-ide-collabsdk.ChatMessage](#optimus-ide-collabsdkchatmessage)             | false    |              |             |
| `queued_messages` | array of [optimus-ide-collabsdk.ChatQueuedMessage](#optimus-ide-collabsdkchatqueuedmessage) | false    |              |             |

## optimus-ide-collabsdk.ChatModel

```json
{
  "display_name": "string",
  "id": "string",
  "model": "string",
  "provider": "string"
}
```

### Properties

| Name           | Type   | Required | Restrictions | Description |
|----------------|--------|----------|--------------|-------------|
| `display_name` | string | false    |              |             |
| `id`           | string | false    |              |             |
| `model`        | string | false    |              |             |
| `provider`     | string | false    |              |             |

## optimus-ide-collabsdk.ChatModelProvider

```json
{
  "available": true,
  "models": [
    {
      "display_name": "string",
      "id": "string",
      "model": "string",
      "provider": "string"
    }
  ],
  "provider": "string",
  "unavailable_reason": "missing_api_key"
}
```

### Properties

| Name                 | Type                                                                                       | Required | Restrictions | Description |
|----------------------|--------------------------------------------------------------------------------------------|----------|--------------|-------------|
| `available`          | boolean                                                                                    | false    |              |             |
| `models`             | array of [optimus-ide-collabsdk.ChatModel](#optimus-ide-collabsdkchatmodel)                                          | false    |              |             |
| `provider`           | string                                                                                     | false    |              |             |
| `unavailable_reason` | [optimus-ide-collabsdk.ChatModelProviderUnavailableReason](#optimus-ide-collabsdkchatmodelproviderunavailablereason) | false    |              |             |

## optimus-ide-collabsdk.ChatModelProviderUnavailableReason

```json
"missing_api_key"
```

### Properties

#### Enumerated Values

| Value(s)                                                   |
|------------------------------------------------------------|
| `fetch_failed`, `missing_api_key`, `user_api_key_required` |

## optimus-ide-collabsdk.ChatModelsResponse

```json
{
  "providers": [
    {
      "available": true,
      "models": [
        {
          "display_name": "string",
          "id": "string",
          "model": "string",
          "provider": "string"
        }
      ],
      "provider": "string",
      "unavailable_reason": "missing_api_key"
    }
  ],
  "unsupported_providers": [
    {
      "display_name": "string",
      "provider": "string"
    }
  ]
}
```

### Properties

| Name                    | Type                                                                          | Required | Restrictions | Description                                                                                                            |
|-------------------------|-------------------------------------------------------------------------------|----------|--------------|------------------------------------------------------------------------------------------------------------------------|
| `providers`             | array of [optimus-ide-collabsdk.ChatModelProvider](#optimus-ide-collabsdkchatmodelprovider)             | false    |              |                                                                                                                        |
| `unsupported_providers` | array of [optimus-ide-collabsdk.ChatUnsupportedProvider](#optimus-ide-collabsdkchatunsupportedprovider) | false    |              | Unsupported providers lists configured providers the Agents harness cannot use, so the UI can explain the empty state. |

## optimus-ide-collabsdk.ChatPlanMode

```json
"plan"
```

### Properties

#### Enumerated Values

| Value(s) |
|----------|
| `plan`   |

## optimus-ide-collabsdk.ChatPrompt

```json
{
  "id": 0,
  "text": "string"
}
```

### Properties

| Name   | Type    | Required | Restrictions | Description |
|--------|---------|----------|--------------|-------------|
| `id`   | integer | false    |              |             |
| `text` | string  | false    |              |             |

## optimus-ide-collabsdk.ChatPromptsResponse

```json
{
  "prompts": [
    {
      "id": 0,
      "text": "string"
    }
  ]
}
```

### Properties

| Name      | Type                                                | Required | Restrictions | Description |
|-----------|-----------------------------------------------------|----------|--------------|-------------|
| `prompts` | array of [optimus-ide-collabsdk.ChatPrompt](#optimus-ide-collabsdkchatprompt) | false    |              |             |

## optimus-ide-collabsdk.ChatQueuedMessage

```json
{
  "chat_id": "efc9fe20-a1e5-4a8c-9c48-f1b30c1e4f86",
  "content": [
    {
      "args": [
        0
      ],
      "args_delta": "string",
      "completed_at": "2019-08-24T14:15:22Z",
      "content": "string",
      "context_file_agent_id": {
        "uuid": "string",
        "valid": true
      },
      "context_file_content": "string",
      "context_file_directory": "string",
      "context_file_os": "string",
      "context_file_path": "string",
      "context_file_skill_meta_file": "string",
      "context_file_truncated": true,
      "created_at": "2019-08-24T14:15:22Z",
      "data": [
        0
      ],
      "end_line": 0,
      "file_id": {
        "uuid": "string",
        "valid": true
      },
      "file_name": "string",
      "is_error": true,
      "is_media": true,
      "mcp_server_config_id": {
        "uuid": "string",
        "valid": true
      },
      "media_type": "string",
      "name": "string",
      "parsed_commands": [
        [
          "string"
        ]
      ],
      "provider_executed": true,
      "provider_metadata": [
        0
      ],
      "result": [
        0
      ],
      "result_delta": "string",
      "result_reset": true,
      "skill_description": "string",
      "skill_dir": "string",
      "skill_name": "string",
      "source_id": "string",
      "start_line": 0,
      "text": "string",
      "title": "string",
      "tool_call_id": "string",
      "tool_name": "string",
      "type": "text",
      "url": "string"
    }
  ],
  "created_at": "2019-08-24T14:15:22Z",
  "id": 0,
  "model_config_id": "f5fb4d91-62ca-4377-9ee6-5d43ba00d205"
}
```

### Properties

| Name              | Type                                                          | Required | Restrictions | Description |
|-------------------|---------------------------------------------------------------|----------|--------------|-------------|
| `chat_id`         | string                                                        | false    |              |             |
| `content`         | array of [optimus-ide-collabsdk.ChatMessagePart](#optimus-ide-collabsdkchatmessagepart) | false    |              |             |
| `created_at`      | string                                                        | false    |              |             |
| `id`              | integer                                                       | false    |              |             |
| `model_config_id` | string                                                        | false    |              |             |

## optimus-ide-collabsdk.ChatRetentionDaysResponse

```json
{
  "retention_days": 0
}
```

### Properties

| Name             | Type    | Required | Restrictions | Description |
|------------------|---------|----------|--------------|-------------|
| `retention_days` | integer | false    |              |             |

## optimus-ide-collabsdk.ChatRole

```json
"read"
```

### Properties

#### Enumerated Values

| Value(s)   |
|------------|
| ``, `read` |

## optimus-ide-collabsdk.ChatStatus

```json
"waiting"
```

### Properties

#### Enumerated Values

| Value(s)                                                         |
|------------------------------------------------------------------|
| `error`, `interrupting`, `requires_action`, `running`, `waiting` |

## optimus-ide-collabsdk.ChatStreamActionRequired

```json
{
  "tool_calls": [
    {
      "args": "string",
      "tool_call_id": "string",
      "tool_name": "string"
    }
  ]
}
```

### Properties

| Name         | Type                                                                | Required | Restrictions | Description |
|--------------|---------------------------------------------------------------------|----------|--------------|-------------|
| `tool_calls` | array of [optimus-ide-collabsdk.ChatStreamToolCall](#optimus-ide-collabsdkchatstreamtoolcall) | false    |              |             |

## optimus-ide-collabsdk.ChatStreamEvent

```json
{
  "action_required": {
    "tool_calls": [
      {
        "args": "string",
        "tool_call_id": "string",
        "tool_name": "string"
      }
    ]
  },
  "chat_id": "efc9fe20-a1e5-4a8c-9c48-f1b30c1e4f86",
  "error": {
    "detail": "string",
    "kind": "generic",
    "message": "string",
    "provider": "string",
    "retryable": true,
    "status_code": 0
  },
  "message": {
    "chat_id": "efc9fe20-a1e5-4a8c-9c48-f1b30c1e4f86",
    "content": [
      {
        "args": [
          0
        ],
        "args_delta": "string",
        "completed_at": "2019-08-24T14:15:22Z",
        "content": "string",
        "context_file_agent_id": {
          "uuid": "string",
          "valid": true
        },
        "context_file_content": "string",
        "context_file_directory": "string",
        "context_file_os": "string",
        "context_file_path": "string",
        "context_file_skill_meta_file": "string",
        "context_file_truncated": true,
        "created_at": "2019-08-24T14:15:22Z",
        "data": [
          0
        ],
        "end_line": 0,
        "file_id": {
          "uuid": "string",
          "valid": true
        },
        "file_name": "string",
        "is_error": true,
        "is_media": true,
        "mcp_server_config_id": {
          "uuid": "string",
          "valid": true
        },
        "media_type": "string",
        "name": "string",
        "parsed_commands": [
          [
            "string"
          ]
        ],
        "provider_executed": true,
        "provider_metadata": [
          0
        ],
        "result": [
          0
        ],
        "result_delta": "string",
        "result_reset": true,
        "skill_description": "string",
        "skill_dir": "string",
        "skill_name": "string",
        "source_id": "string",
        "start_line": 0,
        "text": "string",
        "title": "string",
        "tool_call_id": "string",
        "tool_name": "string",
        "type": "text",
        "url": "string"
      }
    ],
    "created_at": "2019-08-24T14:15:22Z",
    "created_by": "ee824cad-d7a6-4f48-87dc-e8461a9201c4",
    "id": 0,
    "model_config_id": "f5fb4d91-62ca-4377-9ee6-5d43ba00d205",
    "role": "system",
    "usage": {
      "cache_creation_tokens": 0,
      "cache_read_tokens": 0,
      "context_limit": 0,
      "input_tokens": 0,
      "output_tokens": 0,
      "reasoning_tokens": 0,
      "total_tokens": 0
    }
  },
  "message_part": {
    "generation_attempt": 0,
    "history_version": 0,
    "part": {
      "args": [
        0
      ],
      "args_delta": "string",
      "completed_at": "2019-08-24T14:15:22Z",
      "content": "string",
      "context_file_agent_id": {
        "uuid": "string",
        "valid": true
      },
      "context_file_content": "string",
      "context_file_directory": "string",
      "context_file_os": "string",
      "context_file_path": "string",
      "context_file_skill_meta_file": "string",
      "context_file_truncated": true,
      "created_at": "2019-08-24T14:15:22Z",
      "data": [
        0
      ],
      "end_line": 0,
      "file_id": {
        "uuid": "string",
        "valid": true
      },
      "file_name": "string",
      "is_error": true,
      "is_media": true,
      "mcp_server_config_id": {
        "uuid": "string",
        "valid": true
      },
      "media_type": "string",
      "name": "string",
      "parsed_commands": [
        [
          "string"
        ]
      ],
      "provider_executed": true,
      "provider_metadata": [
        0
      ],
      "result": [
        0
      ],
      "result_delta": "string",
      "result_reset": true,
      "skill_description": "string",
      "skill_dir": "string",
      "skill_name": "string",
      "source_id": "string",
      "start_line": 0,
      "text": "string",
      "title": "string",
      "tool_call_id": "string",
      "tool_name": "string",
      "type": "text",
      "url": "string"
    },
    "role": "system",
    "seq": 0
  },
  "queued_messages": [
    {
      "chat_id": "efc9fe20-a1e5-4a8c-9c48-f1b30c1e4f86",
      "content": [
        {
          "args": [
            0
          ],
          "args_delta": "string",
          "completed_at": "2019-08-24T14:15:22Z",
          "content": "string",
          "context_file_agent_id": {
            "uuid": "string",
            "valid": true
          },
          "context_file_content": "string",
          "context_file_directory": "string",
          "context_file_os": "string",
          "context_file_path": "string",
          "context_file_skill_meta_file": "string",
          "context_file_truncated": true,
          "created_at": "2019-08-24T14:15:22Z",
          "data": [
            0
          ],
          "end_line": 0,
          "file_id": {
            "uuid": "string",
            "valid": true
          },
          "file_name": "string",
          "is_error": true,
          "is_media": true,
          "mcp_server_config_id": {
            "uuid": "string",
            "valid": true
          },
          "media_type": "string",
          "name": "string",
          "parsed_commands": [
            [
              "string"
            ]
          ],
          "provider_executed": true,
          "provider_metadata": [
            0
          ],
          "result": [
            0
          ],
          "result_delta": "string",
          "result_reset": true,
          "skill_description": "string",
          "skill_dir": "string",
          "skill_name": "string",
          "source_id": "string",
          "start_line": 0,
          "text": "string",
          "title": "string",
          "tool_call_id": "string",
          "tool_name": "string",
          "type": "text",
          "url": "string"
        }
      ],
      "created_at": "2019-08-24T14:15:22Z",
      "id": 0,
      "model_config_id": "f5fb4d91-62ca-4377-9ee6-5d43ba00d205"
    }
  ],
  "retry": {
    "attempt": 0,
    "delay_ms": 0,
    "error": "string",
    "kind": "generic",
    "provider": "string",
    "retrying_at": "2019-08-24T14:15:22Z",
    "status_code": 0
  },
  "status": {
    "status": "waiting"
  },
  "type": "message_part"
}
```

### Properties

| Name              | Type                                                                   | Required | Restrictions | Description |
|-------------------|------------------------------------------------------------------------|----------|--------------|-------------|
| `action_required` | [optimus-ide-collabsdk.ChatStreamActionRequired](#optimus-ide-collabsdkchatstreamactionrequired) | false    |              |             |
| `chat_id`         | string                                                                 | false    |              |             |
| `error`           | [optimus-ide-collabsdk.ChatError](#optimus-ide-collabsdkchaterror)                               | false    |              |             |
| `message`         | [optimus-ide-collabsdk.ChatMessage](#optimus-ide-collabsdkchatmessage)                           | false    |              |             |
| `message_part`    | [optimus-ide-collabsdk.ChatStreamMessagePart](#optimus-ide-collabsdkchatstreammessagepart)       | false    |              |             |
| `queued_messages` | array of [optimus-ide-collabsdk.ChatQueuedMessage](#optimus-ide-collabsdkchatqueuedmessage)      | false    |              |             |
| `retry`           | [optimus-ide-collabsdk.ChatStreamRetry](#optimus-ide-collabsdkchatstreamretry)                   | false    |              |             |
| `status`          | [optimus-ide-collabsdk.ChatStreamStatus](#optimus-ide-collabsdkchatstreamstatus)                 | false    |              |             |
| `type`            | [optimus-ide-collabsdk.ChatStreamEventType](#optimus-ide-collabsdkchatstreameventtype)           | false    |              |             |

## optimus-ide-collabsdk.ChatStreamEventType

```json
"message_part"
```

### Properties

#### Enumerated Values

| Value(s)                                                                                                                   |
|----------------------------------------------------------------------------------------------------------------------------|
| `action_required`, `error`, `history_reset`, `message`, `message_part`, `preview_reset`, `queue_update`, `retry`, `status` |

## optimus-ide-collabsdk.ChatStreamMessagePart

```json
{
  "generation_attempt": 0,
  "history_version": 0,
  "part": {
    "args": [
      0
    ],
    "args_delta": "string",
    "completed_at": "2019-08-24T14:15:22Z",
    "content": "string",
    "context_file_agent_id": {
      "uuid": "string",
      "valid": true
    },
    "context_file_content": "string",
    "context_file_directory": "string",
    "context_file_os": "string",
    "context_file_path": "string",
    "context_file_skill_meta_file": "string",
    "context_file_truncated": true,
    "created_at": "2019-08-24T14:15:22Z",
    "data": [
      0
    ],
    "end_line": 0,
    "file_id": {
      "uuid": "string",
      "valid": true
    },
    "file_name": "string",
    "is_error": true,
    "is_media": true,
    "mcp_server_config_id": {
      "uuid": "string",
      "valid": true
    },
    "media_type": "string",
    "name": "string",
    "parsed_commands": [
      [
        "string"
      ]
    ],
    "provider_executed": true,
    "provider_metadata": [
      0
    ],
    "result": [
      0
    ],
    "result_delta": "string",
    "result_reset": true,
    "skill_description": "string",
    "skill_dir": "string",
    "skill_name": "string",
    "source_id": "string",
    "start_line": 0,
    "text": "string",
    "title": "string",
    "tool_call_id": "string",
    "tool_name": "string",
    "type": "text",
    "url": "string"
  },
  "role": "system",
  "seq": 0
}
```

### Properties

| Name                 | Type                                                 | Required | Restrictions | Description |
|----------------------|------------------------------------------------------|----------|--------------|-------------|
| `generation_attempt` | integer                                              | false    |              |             |
| `history_version`    | integer                                              | false    |              |             |
| `part`               | [optimus-ide-collabsdk.ChatMessagePart](#optimus-ide-collabsdkchatmessagepart) | false    |              |             |
| `role`               | [optimus-ide-collabsdk.ChatMessageRole](#optimus-ide-collabsdkchatmessagerole) | false    |              |             |
| `seq`                | integer                                              | false    |              |             |

## optimus-ide-collabsdk.ChatStreamRetry

```json
{
  "attempt": 0,
  "delay_ms": 0,
  "error": "string",
  "kind": "generic",
  "provider": "string",
  "retrying_at": "2019-08-24T14:15:22Z",
  "status_code": 0
}
```

### Properties

| Name          | Type                                             | Required | Restrictions | Description                                                       |
|---------------|--------------------------------------------------|----------|--------------|-------------------------------------------------------------------|
| `attempt`     | integer                                          | false    |              | Attempt is the 1-indexed retry attempt number.                    |
| `delay_ms`    | integer                                          | false    |              | Delay ms is the backoff delay in milliseconds before the retry.   |
| `error`       | string                                           | false    |              | Error is the normalized error message from the failed attempt.    |
| `kind`        | [optimus-ide-collabsdk.ChatErrorKind](#optimus-ide-collabsdkchaterrorkind) | false    |              | Kind classifies the retry reason for consistent client rendering. |
| `provider`    | string                                           | false    |              | Provider identifies the upstream model provider when known.       |
| `retrying_at` | string                                           | false    |              | Retrying at is the timestamp when the retry will be attempted.    |
| `status_code` | integer                                          | false    |              | Status code is the best-effort upstream HTTP status code.         |

## optimus-ide-collabsdk.ChatStreamStatus

```json
{
  "status": "waiting"
}
```

### Properties

| Name     | Type                                       | Required | Restrictions | Description |
|----------|--------------------------------------------|----------|--------------|-------------|
| `status` | [optimus-ide-collabsdk.ChatStatus](#optimus-ide-collabsdkchatstatus) | false    |              |             |

## optimus-ide-collabsdk.ChatStreamToolCall

```json
{
  "args": "string",
  "tool_call_id": "string",
  "tool_name": "string"
}
```

### Properties

| Name           | Type   | Required | Restrictions | Description |
|----------------|--------|----------|--------------|-------------|
| `args`         | string | false    |              |             |
| `tool_call_id` | string | false    |              |             |
| `tool_name`    | string | false    |              |             |

## optimus-ide-collabsdk.ChatUnsupportedProvider

```json
{
  "display_name": "string",
  "provider": "string"
}
```

### Properties

| Name           | Type   | Required | Restrictions | Description                                    |
|----------------|--------|----------|--------------|------------------------------------------------|
| `display_name` | string | false    |              |                                                |
| `provider`     | string | false    |              | Provider is the provider type, e.g. "copilot". |

## optimus-ide-collabsdk.ChatUser

```json
{
  "avatar_url": "http://example.com",
  "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
  "name": "string",
  "role": "read",
  "username": "string"
}
```

### Properties

| Name         | Type                                   | Required | Restrictions | Description |
|--------------|----------------------------------------|----------|--------------|-------------|
| `avatar_url` | string                                 | false    |              |             |
| `id`         | string                                 | true     |              |             |
| `name`       | string                                 | false    |              |             |
| `role`       | [optimus-ide-collabsdk.ChatRole](#optimus-ide-collabsdkchatrole) | false    |              |             |
| `username`   | string                                 | true     |              |             |

#### Enumerated Values

| Property | Value(s) |
|----------|----------|
| `role`   | `read`   |

## optimus-ide-collabsdk.ChatWatchEvent

```json
{
  "chat": {
    "agent_id": "2b1e3b65-2c04-4fa2-a2d7-467901e98978",
    "archived": true,
    "build_id": "bfb1f3fa-bf7b-43a5-9e0b-26cc050e44cb",
    "children": [
      {}
    ],
    "client_type": "ui",
    "context": {
      "dirty": true,
      "dirty_since": "2019-08-24T14:15:22Z",
      "error": "string",
      "resources": [
        {
          "error": "string",
          "kind": "instruction_file",
          "size_bytes": 0,
          "skill_description": "string",
          "skill_name": "string",
          "source": "string",
          "status": "ok",
          "tools": [
            {
              "description": "string",
              "name": "string"
            }
          ]
        }
      ]
    },
    "created_at": "2019-08-24T14:15:22Z",
    "diff_status": {
      "additions": 0,
      "approved": true,
      "author_avatar_url": "string",
      "author_login": "string",
      "base_branch": "string",
      "changed_files": 0,
      "changes_requested": true,
      "chat_id": "efc9fe20-a1e5-4a8c-9c48-f1b30c1e4f86",
      "commits": 0,
      "deletions": 0,
      "head_branch": "string",
      "pr_number": 0,
      "pull_request_draft": true,
      "pull_request_state": "string",
      "pull_request_title": "string",
      "refreshed_at": "2019-08-24T14:15:22Z",
      "reviewer_count": 0,
      "stale_at": "2019-08-24T14:15:22Z",
      "url": "string"
    },
    "files": [
      {
        "created_at": "2019-08-24T14:15:22Z",
        "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
        "mime_type": "string",
        "name": "string",
        "organization_id": "7c60d51f-b44e-4682-87d6-449835ea4de6",
        "owner_id": "8826ee2e-7933-4665-aef2-2393f84a0d05"
      }
    ],
    "has_unread": true,
    "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
    "labels": {
      "property1": "string",
      "property2": "string"
    },
    "last_error": {
      "detail": "string",
      "kind": "generic",
      "message": "string",
      "provider": "string",
      "retryable": true,
      "status_code": 0
    },
    "last_model_config_id": "30ebb95f-c255-4759-9429-89aa4ec1554c",
    "last_reasoning_effort": "string",
    "last_turn_summary": "string",
    "mcp_server_ids": [
      "497f6eca-6276-4993-bfeb-53cbbbba6f08"
    ],
    "organization_id": "7c60d51f-b44e-4682-87d6-449835ea4de6",
    "owner_id": "8826ee2e-7933-4665-aef2-2393f84a0d05",
    "owner_name": "string",
    "owner_username": "string",
    "parent_chat_id": "c3609ee6-3b11-4a93-b9ae-e4fabcc99359",
    "pin_order": 0,
    "plan_mode": "plan",
    "root_chat_id": "2898031c-fdce-4e3e-8c53-4481dd42fcd7",
    "shared": true,
    "status": "waiting",
    "summary": "string",
    "title": "string",
    "updated_at": "2019-08-24T14:15:22Z",
    "warnings": [
      "string"
    ],
    "workspace_id": "0967198e-ec7b-4c6b-b4d3-f71244cadbe9"
  },
  "kind": "status_change",
  "tool_calls": [
    {
      "args": "string",
      "tool_call_id": "string",
      "tool_name": "string"
    }
  ]
}
```

### Properties

| Name         | Type                                                                | Required | Restrictions | Description |
|--------------|---------------------------------------------------------------------|----------|--------------|-------------|
| `chat`       | [optimus-ide-collabsdk.Chat](#optimus-ide-collabsdkchat)                                      | false    |              |             |
| `kind`       | [optimus-ide-collabsdk.ChatWatchEventKind](#optimus-ide-collabsdkchatwatcheventkind)          | false    |              |             |
| `tool_calls` | array of [optimus-ide-collabsdk.ChatStreamToolCall](#optimus-ide-collabsdkchatstreamtoolcall) | false    |              |             |

## optimus-ide-collabsdk.ChatWatchEventKind

```json
"status_change"
```

### Properties

#### Enumerated Values

| Value(s)                                                                                                                                                 |
|----------------------------------------------------------------------------------------------------------------------------------------------------------|
| `action_required`, `chat_summary_change`, `context_dirty`, `created`, `deleted`, `diff_status_change`, `status_change`, `summary_change`, `title_change` |

## optimus-ide-collabsdk.ClusterConfig

```json
{
  "host": "string"
}
```

### Properties

| Name   | Type   | Required | Restrictions | Description |
|--------|--------|----------|--------------|-------------|
| `host` | string | false    |              |             |

## optimus-ide-collabsdk.ConnectionLatency

```json
{
  "p50": 31.312,
  "p95": 119.832
}
```

### Properties

| Name  | Type   | Required | Restrictions | Description |
|-------|--------|----------|--------------|-------------|
| `p50` | number | false    |              |             |
| `p95` | number | false    |              |             |

## optimus-ide-collabsdk.ConnectionLog

```json
{
  "agent_name": "string",
  "connect_time": "2019-08-24T14:15:22Z",
  "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
  "ip": "string",
  "organization": {
    "display_name": "string",
    "icon": "string",
    "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
    "name": "string"
  },
  "ssh_info": {
    "connection_id": "d3547de1-d1f2-4344-b4c2-17169b7526f9",
    "disconnect_reason": "string",
    "disconnect_time": "2019-08-24T14:15:22Z",
    "exit_code": 0
  },
  "type": "ssh",
  "web_info": {
    "slug_or_port": "string",
    "status_code": 0,
    "user": {
      "avatar_url": "http://example.com",
      "created_at": "2019-08-24T14:15:22Z",
      "email": "user@example.com",
      "has_ai_seat": true,
      "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
      "is_service_account": true,
      "last_seen_at": "2019-08-24T14:15:22Z",
      "login_type": "",
      "name": "string",
      "organization_ids": [
        "497f6eca-6276-4993-bfeb-53cbbbba6f08"
      ],
      "roles": [
        {
          "display_name": "string",
          "name": "string",
          "organization_id": "string"
        }
      ],
      "status": "active",
      "theme_preference": "string",
      "updated_at": "2019-08-24T14:15:22Z",
      "username": "string"
    },
    "user_agent": "string"
  },
  "workspace_id": "0967198e-ec7b-4c6b-b4d3-f71244cadbe9",
  "workspace_name": "string",
  "workspace_owner_id": "e7078695-5279-4c86-8774-3ac2367a2fc7",
  "workspace_owner_username": "string"
}
```

### Properties

| Name                       | Type                                                           | Required | Restrictions | Description                                                                                                                                              |
|----------------------------|----------------------------------------------------------------|----------|--------------|----------------------------------------------------------------------------------------------------------------------------------------------------------|
| `agent_name`               | string                                                         | false    |              |                                                                                                                                                          |
| `connect_time`             | string                                                         | false    |              |                                                                                                                                                          |
| `id`                       | string                                                         | false    |              |                                                                                                                                                          |
| `ip`                       | string                                                         | false    |              |                                                                                                                                                          |
| `organization`             | [optimus-ide-collabsdk.MinimalOrganization](#optimus-ide-collabsdkminimalorganization)   | false    |              |                                                                                                                                                          |
| `ssh_info`                 | [optimus-ide-collabsdk.ConnectionLogSSHInfo](#optimus-ide-collabsdkconnectionlogsshinfo) | false    |              | Ssh info is only set when `type` is one of: - `ConnectionTypeSSH` - `ConnectionTypeReconnectingPTY` - `ConnectionTypeVSCode` - `ConnectionTypeJetBrains` |
| `type`                     | [optimus-ide-collabsdk.ConnectionType](#optimus-ide-collabsdkconnectiontype)             | false    |              |                                                                                                                                                          |
| `web_info`                 | [optimus-ide-collabsdk.ConnectionLogWebInfo](#optimus-ide-collabsdkconnectionlogwebinfo) | false    |              | Web info is only set when `type` is one of: - `ConnectionTypePortForwarding` - `ConnectionTypeWorkspaceApp` - `ConnectionTypeTunnel`                     |
| `workspace_id`             | string                                                         | false    |              |                                                                                                                                                          |
| `workspace_name`           | string                                                         | false    |              |                                                                                                                                                          |
| `workspace_owner_id`       | string                                                         | false    |              |                                                                                                                                                          |
| `workspace_owner_username` | string                                                         | false    |              |                                                                                                                                                          |

## optimus-ide-collabsdk.ConnectionLogResponse

```json
{
  "connection_logs": [
    {
      "agent_name": "string",
      "connect_time": "2019-08-24T14:15:22Z",
      "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
      "ip": "string",
      "organization": {
        "display_name": "string",
        "icon": "string",
        "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
        "name": "string"
      },
      "ssh_info": {
        "connection_id": "d3547de1-d1f2-4344-b4c2-17169b7526f9",
        "disconnect_reason": "string",
        "disconnect_time": "2019-08-24T14:15:22Z",
        "exit_code": 0
      },
      "type": "ssh",
      "web_info": {
        "slug_or_port": "string",
        "status_code": 0,
        "user": {
          "avatar_url": "http://example.com",
          "created_at": "2019-08-24T14:15:22Z",
          "email": "user@example.com",
          "has_ai_seat": true,
          "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
          "is_service_account": true,
          "last_seen_at": "2019-08-24T14:15:22Z",
          "login_type": "",
          "name": "string",
          "organization_ids": [
            "497f6eca-6276-4993-bfeb-53cbbbba6f08"
          ],
          "roles": [
            {
              "display_name": "string",
              "name": "string",
              "organization_id": "string"
            }
          ],
          "status": "active",
          "theme_preference": "string",
          "updated_at": "2019-08-24T14:15:22Z",
          "username": "string"
        },
        "user_agent": "string"
      },
      "workspace_id": "0967198e-ec7b-4c6b-b4d3-f71244cadbe9",
      "workspace_name": "string",
      "workspace_owner_id": "e7078695-5279-4c86-8774-3ac2367a2fc7",
      "workspace_owner_username": "string"
    }
  ],
  "count": 0,
  "count_cap": 0
}
```

### Properties

| Name              | Type                                                      | Required | Restrictions | Description |
|-------------------|-----------------------------------------------------------|----------|--------------|-------------|
| `connection_logs` | array of [optimus-ide-collabsdk.ConnectionLog](#optimus-ide-collabsdkconnectionlog) | false    |              |             |
| `count`           | integer                                                   | false    |              |             |
| `count_cap`       | integer                                                   | false    |              |             |

## optimus-ide-collabsdk.ConnectionLogSSHInfo

```json
{
  "connection_id": "d3547de1-d1f2-4344-b4c2-17169b7526f9",
  "disconnect_reason": "string",
  "disconnect_time": "2019-08-24T14:15:22Z",
  "exit_code": 0
}
```

### Properties

| Name                | Type    | Required | Restrictions | Description                                                                                                                           |
|---------------------|---------|----------|--------------|---------------------------------------------------------------------------------------------------------------------------------------|
| `connection_id`     | string  | false    |              |                                                                                                                                       |
| `disconnect_reason` | string  | false    |              | Disconnect reason is omitted if a disconnect event with the same connection ID has not yet been seen.                                 |
| `disconnect_time`   | string  | false    |              | Disconnect time is omitted if a disconnect event with the same connection ID has not yet been seen.                                   |
| `exit_code`         | integer | false    |              | Exit code is the exit code of the SSH session. It is omitted if a disconnect event with the same connection ID has not yet been seen. |

## optimus-ide-collabsdk.ConnectionLogWebInfo

```json
{
  "slug_or_port": "string",
  "status_code": 0,
  "user": {
    "avatar_url": "http://example.com",
    "created_at": "2019-08-24T14:15:22Z",
    "email": "user@example.com",
    "has_ai_seat": true,
    "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
    "is_service_account": true,
    "last_seen_at": "2019-08-24T14:15:22Z",
    "login_type": "",
    "name": "string",
    "organization_ids": [
      "497f6eca-6276-4993-bfeb-53cbbbba6f08"
    ],
    "roles": [
      {
        "display_name": "string",
        "name": "string",
        "organization_id": "string"
      }
    ],
    "status": "active",
    "theme_preference": "string",
    "updated_at": "2019-08-24T14:15:22Z",
    "username": "string"
  },
  "user_agent": "string"
}
```

### Properties

| Name           | Type                           | Required | Restrictions | Description                                                               |
|----------------|--------------------------------|----------|--------------|---------------------------------------------------------------------------|
| `slug_or_port` | string                         | false    |              |                                                                           |
| `status_code`  | integer                        | false    |              | Status code is the HTTP status code of the request.                       |
| `user`         | [optimus-ide-collabsdk.User](#optimus-ide-collabsdkuser) | false    |              | User is omitted if the connection event was from an unauthenticated user. |
| `user_agent`   | string                         | false    |              |                                                                           |

## optimus-ide-collabsdk.ConnectionType

```json
"ssh"
```

### Properties

#### Enumerated Values

| Value(s)                                                                                       |
|------------------------------------------------------------------------------------------------|
| `jetbrains`, `port_forwarding`, `reconnecting_pty`, `ssh`, `tunnel`, `vscode`, `workspace_app` |

## optimus-ide-collabsdk.ConvertLoginRequest

```json
{
  "password": "string",
  "to_type": ""
}
```

### Properties

| Name       | Type                                     | Required | Restrictions | Description                              |
|------------|------------------------------------------|----------|--------------|------------------------------------------|
| `password` | string                                   | true     |              |                                          |
| `to_type`  | [optimus-ide-collabsdk.LoginType](#optimus-ide-collabsdklogintype) | true     |              | To type is the login type to convert to. |

## optimus-ide-collabsdk.CreateAIGatewayKeyRequest

```json
{
  "name": "string"
}
```

### Properties

| Name   | Type   | Required | Restrictions | Description |
|--------|--------|----------|--------------|-------------|
| `name` | string | true     |              |             |

## optimus-ide-collabsdk.CreateAIGatewayKeyResponse

```json
{
  "created_at": "2019-08-24T14:15:22Z",
  "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
  "key": "string",
  "key_prefix": "string",
  "name": "string"
}
```

### Properties

| Name         | Type   | Required | Restrictions | Description |
|--------------|--------|----------|--------------|-------------|
| `created_at` | string | false    |              |             |
| `id`         | string | false    |              |             |
| `key`        | string | false    |              |             |
| `key_prefix` | string | false    |              |             |
| `name`       | string | false    |              |             |

## optimus-ide-collabsdk.CreateAIProviderRequest

```json
{
  "api_keys": [
    "string"
  ],
  "base_url": "string",
  "display_name": "string",
  "enabled": true,
  "icon": "string",
  "name": "string",
  "settings": {},
  "type": "openai"
}
```

### Properties

| Name           | Type                                                       | Required | Restrictions | Description |
|----------------|------------------------------------------------------------|----------|--------------|-------------|
| `api_keys`     | array of string                                            | false    |              |             |
| `base_url`     | string                                                     | false    |              |             |
| `display_name` | string                                                     | false    |              |             |
| `enabled`      | boolean                                                    | false    |              |             |
| `icon`         | string                                                     | false    |              |             |
| `name`         | string                                                     | false    |              |             |
| `settings`     | [optimus-ide-collabsdk.AIProviderSettings](#optimus-ide-collabsdkaiprovidersettings) | false    |              |             |
| `type`         | [optimus-ide-collabsdk.AIProviderType](#optimus-ide-collabsdkaiprovidertype)         | false    |              |             |

## optimus-ide-collabsdk.CreateChatMessageRequest

```json
{
  "busy_behavior": "queue",
  "content": [
    {
      "content": "string",
      "end_line": 0,
      "file_id": "8a0cfb4f-ddc9-436d-91bb-75133c583767",
      "file_name": "string",
      "start_line": 0,
      "text": "string",
      "type": "text"
    }
  ],
  "mcp_server_ids": [
    "497f6eca-6276-4993-bfeb-53cbbbba6f08"
  ],
  "model_config_id": "f5fb4d91-62ca-4377-9ee6-5d43ba00d205",
  "plan_mode": "plan",
  "reasoning_effort": "string"
}
```

### Properties

| Name               | Type                                                      | Required | Restrictions | Description                                                                                                  |
|--------------------|-----------------------------------------------------------|----------|--------------|--------------------------------------------------------------------------------------------------------------|
| `busy_behavior`    | [optimus-ide-collabsdk.ChatBusyBehavior](#optimus-ide-collabsdkchatbusybehavior)    | false    |              |                                                                                                              |
| `content`          | array of [optimus-ide-collabsdk.ChatInputPart](#optimus-ide-collabsdkchatinputpart) | false    |              |                                                                                                              |
| `mcp_server_ids`   | array of string                                           | false    |              |                                                                                                              |
| `model_config_id`  | string                                                    | false    |              |                                                                                                              |
| `plan_mode`        | [optimus-ide-collabsdk.ChatPlanMode](#optimus-ide-collabsdkchatplanmode)            | false    |              | Plan mode switches the chat's persistent plan mode. nil: no change, ptr to "plan": enable, ptr to "": clear. |
| `reasoning_effort` | string                                                    | false    |              |                                                                                                              |

#### Enumerated Values

| Property        | Value(s)             |
|-----------------|----------------------|
| `busy_behavior` | `interrupt`, `queue` |

## optimus-ide-collabsdk.CreateChatMessageResponse

```json
{
  "message": {
    "chat_id": "efc9fe20-a1e5-4a8c-9c48-f1b30c1e4f86",
    "content": [
      {
        "args": [
          0
        ],
        "args_delta": "string",
        "completed_at": "2019-08-24T14:15:22Z",
        "content": "string",
        "context_file_agent_id": {
          "uuid": "string",
          "valid": true
        },
        "context_file_content": "string",
        "context_file_directory": "string",
        "context_file_os": "string",
        "context_file_path": "string",
        "context_file_skill_meta_file": "string",
        "context_file_truncated": true,
        "created_at": "2019-08-24T14:15:22Z",
        "data": [
          0
        ],
        "end_line": 0,
        "file_id": {
          "uuid": "string",
          "valid": true
        },
        "file_name": "string",
        "is_error": true,
        "is_media": true,
        "mcp_server_config_id": {
          "uuid": "string",
          "valid": true
        },
        "media_type": "string",
        "name": "string",
        "parsed_commands": [
          [
            "string"
          ]
        ],
        "provider_executed": true,
        "provider_metadata": [
          0
        ],
        "result": [
          0
        ],
        "result_delta": "string",
        "result_reset": true,
        "skill_description": "string",
        "skill_dir": "string",
        "skill_name": "string",
        "source_id": "string",
        "start_line": 0,
        "text": "string",
        "title": "string",
        "tool_call_id": "string",
        "tool_name": "string",
        "type": "text",
        "url": "string"
      }
    ],
    "created_at": "2019-08-24T14:15:22Z",
    "created_by": "ee824cad-d7a6-4f48-87dc-e8461a9201c4",
    "id": 0,
    "model_config_id": "f5fb4d91-62ca-4377-9ee6-5d43ba00d205",
    "role": "system",
    "usage": {
      "cache_creation_tokens": 0,
      "cache_read_tokens": 0,
      "context_limit": 0,
      "input_tokens": 0,
      "output_tokens": 0,
      "reasoning_tokens": 0,
      "total_tokens": 0
    }
  },
  "queued": true,
  "queued_message": {
    "chat_id": "efc9fe20-a1e5-4a8c-9c48-f1b30c1e4f86",
    "content": [
      {
        "args": [
          0
        ],
        "args_delta": "string",
        "completed_at": "2019-08-24T14:15:22Z",
        "content": "string",
        "context_file_agent_id": {
          "uuid": "string",
          "valid": true
        },
        "context_file_content": "string",
        "context_file_directory": "string",
        "context_file_os": "string",
        "context_file_path": "string",
        "context_file_skill_meta_file": "string",
        "context_file_truncated": true,
        "created_at": "2019-08-24T14:15:22Z",
        "data": [
          0
        ],
        "end_line": 0,
        "file_id": {
          "uuid": "string",
          "valid": true
        },
        "file_name": "string",
        "is_error": true,
        "is_media": true,
        "mcp_server_config_id": {
          "uuid": "string",
          "valid": true
        },
        "media_type": "string",
        "name": "string",
        "parsed_commands": [
          [
            "string"
          ]
        ],
        "provider_executed": true,
        "provider_metadata": [
          0
        ],
        "result": [
          0
        ],
        "result_delta": "string",
        "result_reset": true,
        "skill_description": "string",
        "skill_dir": "string",
        "skill_name": "string",
        "source_id": "string",
        "start_line": 0,
        "text": "string",
        "title": "string",
        "tool_call_id": "string",
        "tool_name": "string",
        "type": "text",
        "url": "string"
      }
    ],
    "created_at": "2019-08-24T14:15:22Z",
    "id": 0,
    "model_config_id": "f5fb4d91-62ca-4377-9ee6-5d43ba00d205"
  },
  "warnings": [
    "string"
  ]
}
```

### Properties

| Name             | Type                                                     | Required | Restrictions | Description |
|------------------|----------------------------------------------------------|----------|--------------|-------------|
| `message`        | [optimus-ide-collabsdk.ChatMessage](#optimus-ide-collabsdkchatmessage)             | false    |              |             |
| `queued`         | boolean                                                  | false    |              |             |
| `queued_message` | [optimus-ide-collabsdk.ChatQueuedMessage](#optimus-ide-collabsdkchatqueuedmessage) | false    |              |             |
| `warnings`       | array of string                                          | false    |              |             |

## optimus-ide-collabsdk.CreateChatRequest

```json
{
  "client_type": "ui",
  "content": [
    {
      "content": "string",
      "end_line": 0,
      "file_id": "8a0cfb4f-ddc9-436d-91bb-75133c583767",
      "file_name": "string",
      "start_line": 0,
      "text": "string",
      "type": "text"
    }
  ],
  "labels": {
    "property1": "string",
    "property2": "string"
  },
  "mcp_server_ids": [
    "497f6eca-6276-4993-bfeb-53cbbbba6f08"
  ],
  "model_config_id": "f5fb4d91-62ca-4377-9ee6-5d43ba00d205",
  "organization_id": "7c60d51f-b44e-4682-87d6-449835ea4de6",
  "plan_mode": "plan",
  "reasoning_effort": "string",
  "system_prompt": "string",
  "unsafe_dynamic_tools": [
    {
      "description": "string",
      "input_schema": [
        0
      ],
      "name": "string"
    }
  ],
  "workspace_id": "0967198e-ec7b-4c6b-b4d3-f71244cadbe9"
}
```

### Properties

| Name                   | Type                                                      | Required | Restrictions | Description                                                                                                                                |
|------------------------|-----------------------------------------------------------|----------|--------------|--------------------------------------------------------------------------------------------------------------------------------------------|
| `client_type`          | [optimus-ide-collabsdk.ChatClientType](#optimus-ide-collabsdkchatclienttype)        | false    |              |                                                                                                                                            |
| `content`              | array of [optimus-ide-collabsdk.ChatInputPart](#optimus-ide-collabsdkchatinputpart) | false    |              |                                                                                                                                            |
| `labels`               | object                                                    | false    |              |                                                                                                                                            |
| » `[any property]`     | string                                                    | false    |              |                                                                                                                                            |
| `mcp_server_ids`       | array of string                                           | false    |              |                                                                                                                                            |
| `model_config_id`      | string                                                    | false    |              |                                                                                                                                            |
| `organization_id`      | string                                                    | false    |              |                                                                                                                                            |
| `plan_mode`            | [optimus-ide-collabsdk.ChatPlanMode](#optimus-ide-collabsdkchatplanmode)            | false    |              |                                                                                                                                            |
| `reasoning_effort`     | string                                                    | false    |              |                                                                                                                                            |
| `system_prompt`        | string                                                    | false    |              |                                                                                                                                            |
| `unsafe_dynamic_tools` | array of [optimus-ide-collabsdk.DynamicTool](#optimus-ide-collabsdkdynamictool)     | false    |              | Unsafe dynamic tools declares client-executed tools that the LLM can invoke. This API is highly experimental and highly subject to change. |
| `workspace_id`         | string                                                    | false    |              |                                                                                                                                            |

## optimus-ide-collabsdk.CreateFirstUserOnboardingInfo

```json
{
  "newsletter_marketing": true,
  "newsletter_releases": true
}
```

### Properties

| Name                   | Type    | Required | Restrictions | Description |
|------------------------|---------|----------|--------------|-------------|
| `newsletter_marketing` | boolean | false    |              |             |
| `newsletter_releases`  | boolean | false    |              |             |

## optimus-ide-collabsdk.CreateFirstUserRequest

```json
{
  "email": "string",
  "name": "string",
  "onboarding_info": {
    "newsletter_marketing": true,
    "newsletter_releases": true
  },
  "password": "string",
  "trial": true,
  "trial_info": {
    "company_name": "string",
    "country": "string",
    "developers": "string",
    "first_name": "string",
    "job_title": "string",
    "last_name": "string",
    "phone_number": "string"
  },
  "username": "string"
}
```

### Properties

| Name              | Type                                                                             | Required | Restrictions | Description |
|-------------------|----------------------------------------------------------------------------------|----------|--------------|-------------|
| `email`           | string                                                                           | true     |              |             |
| `name`            | string                                                                           | false    |              |             |
| `onboarding_info` | [optimus-ide-collabsdk.CreateFirstUserOnboardingInfo](#optimus-ide-collabsdkcreatefirstuseronboardinginfo) | false    |              |             |
| `password`        | string                                                                           | true     |              |             |
| `trial`           | boolean                                                                          | false    |              |             |
| `trial_info`      | [optimus-ide-collabsdk.CreateFirstUserTrialInfo](#optimus-ide-collabsdkcreatefirstusertrialinfo)           | false    |              |             |
| `username`        | string                                                                           | true     |              |             |

## optimus-ide-collabsdk.CreateFirstUserResponse

```json
{
  "organization_id": "7c60d51f-b44e-4682-87d6-449835ea4de6",
  "user_id": "a169451c-8525-4352-b8ca-070dd449a1a5"
}
```

### Properties

| Name              | Type   | Required | Restrictions | Description |
|-------------------|--------|----------|--------------|-------------|
| `organization_id` | string | false    |              |             |
| `user_id`         | string | false    |              |             |

## optimus-ide-collabsdk.CreateFirstUserTrialInfo

```json
{
  "company_name": "string",
  "country": "string",
  "developers": "string",
  "first_name": "string",
  "job_title": "string",
  "last_name": "string",
  "phone_number": "string"
}
```

### Properties

| Name           | Type   | Required | Restrictions | Description |
|----------------|--------|----------|--------------|-------------|
| `company_name` | string | false    |              |             |
| `country`      | string | false    |              |             |
| `developers`   | string | false    |              |             |
| `first_name`   | string | false    |              |             |
| `job_title`    | string | false    |              |             |
| `last_name`    | string | false    |              |             |
| `phone_number` | string | false    |              |             |

## optimus-ide-collabsdk.CreateGroupRequest

```json
{
  "avatar_url": "string",
  "display_name": "string",
  "name": "string",
  "quota_allowance": 0
}
```

### Properties

| Name              | Type    | Required | Restrictions | Description |
|-------------------|---------|----------|--------------|-------------|
| `avatar_url`      | string  | false    |              |             |
| `display_name`    | string  | false    |              |             |
| `name`            | string  | true     |              |             |
| `quota_allowance` | integer | false    |              |             |

## optimus-ide-collabsdk.CreateOrganizationRequest

```json
{
  "description": "string",
  "display_name": "string",
  "icon": "string",
  "name": "string"
}
```

### Properties

| Name           | Type   | Required | Restrictions | Description                                                            |
|----------------|--------|----------|--------------|------------------------------------------------------------------------|
| `description`  | string | false    |              |                                                                        |
| `display_name` | string | false    |              | Display name will default to the same value as `Name` if not provided. |
| `icon`         | string | false    |              |                                                                        |
| `name`         | string | true     |              |                                                                        |

## optimus-ide-collabsdk.CreateProvisionerKeyResponse

```json
{
  "key": "string"
}
```

### Properties

| Name  | Type   | Required | Restrictions | Description |
|-------|--------|----------|--------------|-------------|
| `key` | string | false    |              |             |

## optimus-ide-collabsdk.CreateTaskRequest

```json
{
  "display_name": "string",
  "input": "string",
  "name": "string",
  "template_version_id": "0ba39c92-1f1b-4c32-aa3e-9925d7713eb1",
  "template_version_preset_id": "512a53a7-30da-446e-a1fc-713c630baff1"
}
```

### Properties

| Name                         | Type   | Required | Restrictions | Description |
|------------------------------|--------|----------|--------------|-------------|
| `display_name`               | string | false    |              |             |
| `input`                      | string | false    |              |             |
| `name`                       | string | false    |              |             |
| `template_version_id`        | string | false    |              |             |
| `template_version_preset_id` | string | false    |              |             |

## optimus-ide-collabsdk.CreateTemplateRequest

```json
{
  "activity_bump_ms": 0,
  "allow_user_autostart": true,
  "allow_user_autostop": true,
  "allow_user_cancel_workspace_jobs": true,
  "autostart_requirement": {
    "days_of_week": [
      "monday"
    ]
  },
  "autostop_requirement": {
    "days_of_week": [
      "monday"
    ],
    "weeks": 0
  },
  "cors_behavior": "simple",
  "default_ttl_ms": 0,
  "delete_ttl_ms": 0,
  "description": "string",
  "disable_everyone_group_access": true,
  "display_name": "string",
  "dormant_ttl_ms": 0,
  "failure_ttl_ms": 0,
  "icon": "string",
  "max_port_share_level": "owner",
  "name": "string",
  "require_active_version": true,
  "template_use_classic_parameter_flow": true,
  "template_version_id": "0ba39c92-1f1b-4c32-aa3e-9925d7713eb1",
  "time_til_autostop_notify_ms": 0
}
```

### Properties

| Name                                  | Type                                                                           | Required | Restrictions | Description                                                                                                                                                                                                                                                                                                         |
|---------------------------------------|--------------------------------------------------------------------------------|----------|--------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `activity_bump_ms`                    | integer                                                                        | false    |              | Activity bump ms allows optionally specifying the activity bump duration for all workspaces created from this template. Defaults to 1h but can be set to 0 to disable activity bumping.                                                                                                                             |
| `allow_user_autostart`                | boolean                                                                        | false    |              | Allow user autostart allows users to set a schedule for autostarting their workspace. By default this is true. This can only be disabled when using an enterprise license.                                                                                                                                          |
| `allow_user_autostop`                 | boolean                                                                        | false    |              | Allow user autostop allows users to set a custom workspace TTL to use in place of the template's DefaultTTL field. By default this is true. If false, the DefaultTTL will always be used. This can only be disabled when using an enterprise license.                                                               |
| `allow_user_cancel_workspace_jobs`    | boolean                                                                        | false    |              | Allow users to cancel in-progress workspace jobs. *bool as the default value is "true".                                                                                                                                                                                                                             |
| `autostart_requirement`               | [optimus-ide-collabsdk.TemplateAutostartRequirement](#optimus-ide-collabsdktemplateautostartrequirement) | false    |              | Autostart requirement allows optionally specifying the autostart allowed days for workspaces created from this template. This is an enterprise feature.                                                                                                                                                             |
| `autostop_requirement`                | [optimus-ide-collabsdk.TemplateAutostopRequirement](#optimus-ide-collabsdktemplateautostoprequirement)   | false    |              | Autostop requirement allows optionally specifying the autostop requirement for workspaces created from this template. This is an enterprise feature.                                                                                                                                                                |
| `cors_behavior`                       | [optimus-ide-collabsdk.CORSBehavior](#optimus-ide-collabsdkcorsbehavior)                                 | false    |              | Cors behavior allows optionally specifying the CORS behavior for all shared ports.                                                                                                                                                                                                                                  |
| `default_ttl_ms`                      | integer                                                                        | false    |              | Default ttl ms allows optionally specifying the default TTL for all workspaces created from this template.                                                                                                                                                                                                          |
| `delete_ttl_ms`                       | integer                                                                        | false    |              | Delete ttl ms allows optionally specifying the max lifetime before Optimus-IDE-Collab permanently deletes dormant workspaces created from this template.                                                                                                                                                                         |
| `description`                         | string                                                                         | false    |              | Description is a description of what the template contains. It must be less than 128 bytes.                                                                                                                                                                                                                         |
| `disable_everyone_group_access`       | boolean                                                                        | false    |              | Disable everyone group access allows optionally disabling the default behavior of granting the 'everyone' group access to use the template. If this is set to true, the template will not be available to all users, and must be explicitly granted to users or groups in the permissions settings of the template. |
| `display_name`                        | string                                                                         | false    |              | Display name is the displayed name of the template.                                                                                                                                                                                                                                                                 |
| `dormant_ttl_ms`                      | integer                                                                        | false    |              | Dormant ttl ms allows optionally specifying the max lifetime before Optimus-IDE-Collab locks inactive workspaces created from this template.                                                                                                                                                                                     |
| `failure_ttl_ms`                      | integer                                                                        | false    |              | Failure ttl ms allows optionally specifying the max lifetime before Optimus-IDE-Collab stops all resources for failed workspaces created from this template.                                                                                                                                                                     |
| `icon`                                | string                                                                         | false    |              | Icon is a relative path or external URL that specifies an icon to be displayed in the dashboard.                                                                                                                                                                                                                    |
| `max_port_share_level`                | [optimus-ide-collabsdk.WorkspaceAgentPortShareLevel](#optimus-ide-collabsdkworkspaceagentportsharelevel) | false    |              | Max port share level allows optionally specifying the maximum port share level for workspaces created from the template.                                                                                                                                                                                            |
| `name`                                | string                                                                         | true     |              | Name is the name of the template.                                                                                                                                                                                                                                                                                   |
| `require_active_version`              | boolean                                                                        | false    |              | Require active version mandates that workspaces are built with the active template version.                                                                                                                                                                                                                         |
| `template_use_classic_parameter_flow` | boolean                                                                        | false    |              | Template use classic parameter flow allows optionally specifying whether the template should use the classic parameter flow. The default if unset is true, and is why `*bool` is used here. When dynamic parameters becomes the default, this will default to false.                                                |
|`template_version_id`|string|true||Template version ID is an in-progress or completed job to use as an initial version of the template.
This is required on creation to enable a user-flow of validating a template works. There is no reason the data-model cannot support empty templates, but it doesn't make sense for users.|
|`time_til_autostop_notify_ms`|integer|false||Time til autostop notify ms allows optionally specifying the duration before the autostop deadline at which a reminder notification is sent for workspaces created from this template. Defaults to 0 (disabled).|

## optimus-ide-collabsdk.CreateTemplateVersionDryRunRequest

```json
{
  "rich_parameter_values": [
    {
      "name": "string",
      "value": "string"
    }
  ],
  "user_variable_values": [
    {
      "name": "string",
      "value": "string"
    }
  ],
  "workspace_name": "string"
}
```

### Properties

| Name                    | Type                                                                          | Required | Restrictions | Description |
|-------------------------|-------------------------------------------------------------------------------|----------|--------------|-------------|
| `rich_parameter_values` | array of [optimus-ide-collabsdk.WorkspaceBuildParameter](#optimus-ide-collabsdkworkspacebuildparameter) | false    |              |             |
| `user_variable_values`  | array of [optimus-ide-collabsdk.VariableValue](#optimus-ide-collabsdkvariablevalue)                     | false    |              |             |
| `workspace_name`        | string                                                                        | false    |              |             |

## optimus-ide-collabsdk.CreateTemplateVersionRequest

```json
{
  "example_id": "string",
  "file_id": "8a0cfb4f-ddc9-436d-91bb-75133c583767",
  "message": "string",
  "name": "string",
  "provisioner": "terraform",
  "storage_method": "file",
  "tags": {
    "property1": "string",
    "property2": "string"
  },
  "template_id": "c6d67e98-83ea-49f0-8812-e4abae2b68bc",
  "user_variable_values": [
    {
      "name": "string",
      "value": "string"
    }
  ]
}
```

### Properties

| Name                   | Type                                                                   | Required | Restrictions | Description                                                  |
|------------------------|------------------------------------------------------------------------|----------|--------------|--------------------------------------------------------------|
| `example_id`           | string                                                                 | false    |              |                                                              |
| `file_id`              | string                                                                 | false    |              |                                                              |
| `message`              | string                                                                 | false    |              |                                                              |
| `name`                 | string                                                                 | false    |              |                                                              |
| `provisioner`          | string                                                                 | true     |              |                                                              |
| `storage_method`       | [optimus-ide-collabsdk.ProvisionerStorageMethod](#optimus-ide-collabsdkprovisionerstoragemethod) | true     |              |                                                              |
| `tags`                 | object                                                                 | false    |              |                                                              |
| » `[any property]`     | string                                                                 | false    |              |                                                              |
| `template_id`          | string                                                                 | false    |              | Template ID optionally associates a version with a template. |
| `user_variable_values` | array of [optimus-ide-collabsdk.VariableValue](#optimus-ide-collabsdkvariablevalue)              | false    |              |                                                              |

#### Enumerated Values

| Property         | Value(s)            |
|------------------|---------------------|
| `provisioner`    | `echo`, `terraform` |
| `storage_method` | `file`              |

## optimus-ide-collabsdk.CreateTestAuditLogRequest

```json
{
  "action": "create",
  "additional_fields": [
    0
  ],
  "build_reason": "autostart",
  "organization_id": "7c60d51f-b44e-4682-87d6-449835ea4de6",
  "request_id": "266ea41d-adf5-480b-af50-15b940c2b846",
  "resource_id": "4d5215ed-38bb-48ed-879a-fdb9ca58522f",
  "resource_type": "template",
  "time": "2019-08-24T14:15:22Z"
}
```

### Properties

| Name                | Type                                           | Required | Restrictions | Description |
|---------------------|------------------------------------------------|----------|--------------|-------------|
| `action`            | [optimus-ide-collabsdk.AuditAction](#optimus-ide-collabsdkauditaction)   | false    |              |             |
| `additional_fields` | array of integer                               | false    |              |             |
| `build_reason`      | [optimus-ide-collabsdk.BuildReason](#optimus-ide-collabsdkbuildreason)   | false    |              |             |
| `organization_id`   | string                                         | false    |              |             |
| `request_id`        | string                                         | false    |              |             |
| `resource_id`       | string                                         | false    |              |             |
| `resource_type`     | [optimus-ide-collabsdk.ResourceType](#optimus-ide-collabsdkresourcetype) | false    |              |             |
| `time`              | string                                         | false    |              |             |

#### Enumerated Values

| Property        | Value(s)                                                                                                 |
|-----------------|----------------------------------------------------------------------------------------------------------|
| `action`        | `create`, `delete`, `start`, `stop`, `write`                                                             |
| `build_reason`  | `autostart`, `autostop`, `initiator`                                                                     |
| `resource_type` | `auditable_group`, `git_ssh_key`, `template`, `template_version`, `user`, `workspace`, `workspace_build` |

## optimus-ide-collabsdk.CreateTokenRequest

```json
{
  "allow_list": [
    {
      "id": "string",
      "type": "*"
    }
  ],
  "lifetime": 0,
  "scope": "all",
  "scopes": [
    "all"
  ],
  "token_name": "string"
}
```

### Properties

| Name         | Type                                                                | Required | Restrictions | Description                     |
|--------------|---------------------------------------------------------------------|----------|--------------|---------------------------------|
| `allow_list` | array of [optimus-ide-collabsdk.APIAllowListTarget](#optimus-ide-collabsdkapiallowlisttarget) | false    |              |                                 |
| `lifetime`   | integer                                                             | false    |              |                                 |
| `scope`      | [optimus-ide-collabsdk.APIKeyScope](#optimus-ide-collabsdkapikeyscope)                        | false    |              | Deprecated: use Scopes instead. |
| `scopes`     | array of [optimus-ide-collabsdk.APIKeyScope](#optimus-ide-collabsdkapikeyscope)               | false    |              |                                 |
| `token_name` | string                                                              | false    |              |                                 |

## optimus-ide-collabsdk.CreateUserRequestWithOrgs

```json
{
  "email": "user@example.com",
  "login_type": "",
  "name": "string",
  "organization_ids": [
    "497f6eca-6276-4993-bfeb-53cbbbba6f08"
  ],
  "password": "string",
  "roles": [
    "string"
  ],
  "service_account": true,
  "user_status": "active",
  "username": "string"
}
```

### Properties

| Name               | Type                                       | Required | Restrictions | Description                                                                         |
|--------------------|--------------------------------------------|----------|--------------|-------------------------------------------------------------------------------------|
| `email`            | string                                     | false    |              |                                                                                     |
| `login_type`       | [optimus-ide-collabsdk.LoginType](#optimus-ide-collabsdklogintype)   | false    |              | Login type defaults to LoginTypePassword.                                           |
| `name`             | string                                     | false    |              |                                                                                     |
| `organization_ids` | array of string                            | false    |              | Organization ids is a list of organization IDs that the user should be a member of. |
| `password`         | string                                     | false    |              |                                                                                     |
| `roles`            | array of string                            | false    |              | Roles is an optional list of site-level roles to assign at creation.                |
| `service_account`  | boolean                                    | false    |              | Service accounts are admin-managed accounts that cannot login.                      |
| `user_status`      | [optimus-ide-collabsdk.UserStatus](#optimus-ide-collabsdkuserstatus) | false    |              | User status defaults to UserStatusDormant.                                          |
| `username`         | string                                     | true     |              |                                                                                     |

## optimus-ide-collabsdk.CreateUserSecretRequest

```json
{
  "description": "string",
  "enabled": true,
  "env_name": "string",
  "file_path": "string",
  "name": "string",
  "value": "string"
}
```

### Properties

| Name          | Type    | Required | Restrictions | Description |
|---------------|---------|----------|--------------|-------------|
| `description` | string  | false    |              |             |
| `enabled`     | boolean | false    |              |             |
| `env_name`    | string  | false    |              |             |
| `file_path`   | string  | false    |              |             |
| `name`        | string  | false    |              |             |
| `value`       | string  | false    |              |             |

## optimus-ide-collabsdk.CreateUserSkillRequest

```json
{
  "content": "string"
}
```

### Properties

| Name      | Type   | Required | Restrictions | Description                                                                                                                                                           |
|-----------|--------|----------|--------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `content` | string | false    |              | Content must be SKILL.md-format Markdown with YAML frontmatter. The frontmatter must include name, may include description, and must be followed by a non-empty body. |

## optimus-ide-collabsdk.CreateWorkspaceBuildOnSuccessRequest

```json
{
  "rich_parameter_values": [
    {
      "name": "string",
      "value": "string"
    }
  ],
  "template_version_id": "0ba39c92-1f1b-4c32-aa3e-9925d7713eb1",
  "template_version_preset_id": "512a53a7-30da-446e-a1fc-713c630baff1",
  "transition": "start"
}
```

### Properties

| Name                         | Type                                                                          | Required | Restrictions | Description                                                                                                                                                                                                                                                                       |
|------------------------------|-------------------------------------------------------------------------------|----------|--------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `rich_parameter_values`      | array of [optimus-ide-collabsdk.WorkspaceBuildParameter](#optimus-ide-collabsdkworkspacebuildparameter) | false    |              | Rich parameter values are applied to the child build. Parameters not listed here fall back to their values from the previous build, matching normal build behavior.                                                                                                               |
| `template_version_id`        | string                                                                        | false    |              | Template version ID pins the child build to a specific template version. Pinning requires permission to update the template, since the active version may change before the child build runs. When empty, the child build uses the template's active version at the time it runs. |
| `template_version_preset_id` | string                                                                        | false    |              | Template version preset ID selects a preset for the child build. It requires TemplateVersionID to also be set.                                                                                                                                                                    |
| `transition`                 | [optimus-ide-collabsdk.WorkspaceTransition](#optimus-ide-collabsdkworkspacetransition)                  | true     |              | Transition must be "start". The parent build's transition must be "stop".                                                                                                                                                                                                         |

#### Enumerated Values

| Property     | Value(s) |
|--------------|----------|
| `transition` | `start`  |

## optimus-ide-collabsdk.CreateWorkspaceBuildReason

```json
"dashboard"
```

### Properties

#### Enumerated Values

| Value(s)                                                                                                              |
|-----------------------------------------------------------------------------------------------------------------------|
| `cli`, `dashboard`, `jetbrains_connection`, `ssh_connection`, `task_manual_pause`, `task_resume`, `vscode_connection` |

## optimus-ide-collabsdk.CreateWorkspaceBuildRequest

```json
{
  "dry_run": true,
  "log_level": "debug",
  "on_success": {
    "rich_parameter_values": [
      {
        "name": "string",
        "value": "string"
      }
    ],
    "template_version_id": "0ba39c92-1f1b-4c32-aa3e-9925d7713eb1",
    "template_version_preset_id": "512a53a7-30da-446e-a1fc-713c630baff1",
    "transition": "start"
  },
  "orphan": true,
  "reason": "dashboard",
  "rich_parameter_values": [
    {
      "name": "string",
      "value": "string"
    }
  ],
  "state": [
    0
  ],
  "template_version_id": "0ba39c92-1f1b-4c32-aa3e-9925d7713eb1",
  "template_version_preset_id": "512a53a7-30da-446e-a1fc-713c630baff1",
  "transition": "start"
}
```

### Properties

| Name                         | Type                                                                                           | Required | Restrictions | Description                                                                                                                                                                                                   |
|------------------------------|------------------------------------------------------------------------------------------------|----------|--------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `dry_run`                    | boolean                                                                                        | false    |              |                                                                                                                                                                                                               |
| `log_level`                  | [optimus-ide-collabsdk.ProvisionerLogLevel](#optimus-ide-collabsdkprovisionerloglevel)                                   | false    |              | Log level changes the default logging verbosity of a provider ("info" if empty).                                                                                                                              |
| `on_success`                 | [optimus-ide-collabsdk.CreateWorkspaceBuildOnSuccessRequest](#optimus-ide-collabsdkcreateworkspacebuildonsuccessrequest) | false    |              | On success queues a follow-up workspace build after this build succeeds. It currently supports restarting a workspace by starting it after a successful stop build.                                           |
| `orphan`                     | boolean                                                                                        | false    |              | Orphan may be set for the Destroy transition.                                                                                                                                                                 |
| `reason`                     | [optimus-ide-collabsdk.CreateWorkspaceBuildReason](#optimus-ide-collabsdkcreateworkspacebuildreason)                     | false    |              | Reason sets the reason for the workspace build.                                                                                                                                                               |
| `rich_parameter_values`      | array of [optimus-ide-collabsdk.WorkspaceBuildParameter](#optimus-ide-collabsdkworkspacebuildparameter)                  | false    |              | Rich parameter values are optional. It will write params to the 'workspace' scope. This will overwrite any existing parameters with the same name. This will not delete old params not included in this list. |
| `state`                      | array of integer                                                                               | false    |              |                                                                                                                                                                                                               |
| `template_version_id`        | string                                                                                         | false    |              |                                                                                                                                                                                                               |
| `template_version_preset_id` | string                                                                                         | false    |              | Template version preset ID is the ID of the template version preset to use for the build.                                                                                                                     |
| `transition`                 | [optimus-ide-collabsdk.WorkspaceTransition](#optimus-ide-collabsdkworkspacetransition)                                   | true     |              |                                                                                                                                                                                                               |

#### Enumerated Values

| Property     | Value(s)                                                                                               |
|--------------|--------------------------------------------------------------------------------------------------------|
| `log_level`  | `debug`                                                                                                |
| `reason`     | `cli`, `dashboard`, `jetbrains_connection`, `ssh_connection`, `task_manual_pause`, `vscode_connection` |
| `transition` | `delete`, `start`, `stop`                                                                              |

## optimus-ide-collabsdk.CreateWorkspaceProxyRequest

```json
{
  "display_name": "string",
  "icon": "string",
  "name": "string"
}
```

### Properties

| Name           | Type   | Required | Restrictions | Description |
|----------------|--------|----------|--------------|-------------|
| `display_name` | string | false    |              |             |
| `icon`         | string | false    |              |             |
| `name`         | string | true     |              |             |

## optimus-ide-collabsdk.CreateWorkspaceRequest

```json
{
  "automatic_updates": "always",
  "autostart_schedule": "string",
  "name": "string",
  "rich_parameter_values": [
    {
      "name": "string",
      "value": "string"
    }
  ],
  "template_id": "c6d67e98-83ea-49f0-8812-e4abae2b68bc",
  "template_version_id": "0ba39c92-1f1b-4c32-aa3e-9925d7713eb1",
  "template_version_preset_id": "512a53a7-30da-446e-a1fc-713c630baff1",
  "ttl_ms": 0
}
```

CreateWorkspaceRequest provides options for creating a new workspace. Only one of TemplateID or TemplateVersionID can be specified, not both. If TemplateID is specified, the active version of the template will be used. Workspace names: - Must start with a letter or number - Can only contain letters, numbers, and hyphens - Cannot contain spaces or special characters - Cannot be named `new` or `create` - Must be unique within your workspaces - Maximum length of 32 characters

### Properties

| Name                         | Type                                                                          | Required | Restrictions | Description                                                                                             |
|------------------------------|-------------------------------------------------------------------------------|----------|--------------|---------------------------------------------------------------------------------------------------------|
| `automatic_updates`          | [optimus-ide-collabsdk.AutomaticUpdates](#optimus-ide-collabsdkautomaticupdates)                        | false    |              |                                                                                                         |
| `autostart_schedule`         | string                                                                        | false    |              |                                                                                                         |
| `name`                       | string                                                                        | true     |              |                                                                                                         |
| `rich_parameter_values`      | array of [optimus-ide-collabsdk.WorkspaceBuildParameter](#optimus-ide-collabsdkworkspacebuildparameter) | false    |              | Rich parameter values allows for additional parameters to be provided during the initial provision.     |
| `template_id`                | string                                                                        | false    |              | Template ID specifies which template should be used for creating the workspace.                         |
| `template_version_id`        | string                                                                        | false    |              | Template version ID can be used to specify a specific version of a template for creating the workspace. |
| `template_version_preset_id` | string                                                                        | false    |              |                                                                                                         |
| `ttl_ms`                     | integer                                                                       | false    |              |                                                                                                         |

## optimus-ide-collabsdk.CryptoKey

```json
{
  "deletes_at": "2019-08-24T14:15:22Z",
  "feature": "workspace_apps_api_key",
  "secret": "string",
  "sequence": 0,
  "starts_at": "2019-08-24T14:15:22Z"
}
```

### Properties

| Name         | Type                                                   | Required | Restrictions | Description |
|--------------|--------------------------------------------------------|----------|--------------|-------------|
| `deletes_at` | string                                                 | false    |              |             |
| `feature`    | [optimus-ide-collabsdk.CryptoKeyFeature](#optimus-ide-collabsdkcryptokeyfeature) | false    |              |             |
| `secret`     | string                                                 | false    |              |             |
| `sequence`   | integer                                                | false    |              |             |
| `starts_at`  | string                                                 | false    |              |             |

## optimus-ide-collabsdk.CryptoKeyFeature

```json
"workspace_apps_api_key"
```

### Properties

#### Enumerated Values

| Value(s)                                                                                      |
|-----------------------------------------------------------------------------------------------|
| `nats_ca`, `oidc_convert`, `tailnet_resume`, `workspace_apps_api_key`, `workspace_apps_token` |

## optimus-ide-collabsdk.CustomNotificationContent

```json
{
  "message": "string",
  "title": "string"
}
```

### Properties

| Name      | Type   | Required | Restrictions | Description |
|-----------|--------|----------|--------------|-------------|
| `message` | string | false    |              |             |
| `title`   | string | false    |              |             |

## optimus-ide-collabsdk.CustomNotificationRequest

```json
{
  "content": {
    "message": "string",
    "title": "string"
  }
}
```

### Properties

| Name      | Type                                                                     | Required | Restrictions | Description |
|-----------|--------------------------------------------------------------------------|----------|--------------|-------------|
| `content` | [optimus-ide-collabsdk.CustomNotificationContent](#optimus-ide-collabsdkcustomnotificationcontent) | false    |              |             |

## optimus-ide-collabsdk.CustomRoleRequest

```json
{
  "display_name": "string",
  "name": "string",
  "organization_member_permissions": [
    {
      "action": "application_connect",
      "negate": true,
      "resource_type": "*"
    }
  ],
  "organization_permissions": [
    {
      "action": "application_connect",
      "negate": true,
      "resource_type": "*"
    }
  ],
  "site_permissions": [
    {
      "action": "application_connect",
      "negate": true,
      "resource_type": "*"
    }
  ],
  "user_permissions": [
    {
      "action": "application_connect",
      "negate": true,
      "resource_type": "*"
    }
  ]
}
```

### Properties

| Name                              | Type                                                | Required | Restrictions | Description                                                                           |
|-----------------------------------|-----------------------------------------------------|----------|--------------|---------------------------------------------------------------------------------------|
| `display_name`                    | string                                              | false    |              |                                                                                       |
| `name`                            | string                                              | false    |              |                                                                                       |
| `organization_member_permissions` | array of [optimus-ide-collabsdk.Permission](#optimus-ide-collabsdkpermission) | false    |              | Organization member permissions are specific to the organization the role belongs to. |
| `organization_permissions`        | array of [optimus-ide-collabsdk.Permission](#optimus-ide-collabsdkpermission) | false    |              | Organization permissions are specific to the organization the role belongs to.        |
| `site_permissions`                | array of [optimus-ide-collabsdk.Permission](#optimus-ide-collabsdkpermission) | false    |              |                                                                                       |
| `user_permissions`                | array of [optimus-ide-collabsdk.Permission](#optimus-ide-collabsdkpermission) | false    |              |                                                                                       |

## optimus-ide-collabsdk.DAUEntry

```json
{
  "amount": 0,
  "date": "string"
}
```

### Properties

| Name     | Type    | Required | Restrictions | Description                                                                              |
|----------|---------|----------|--------------|------------------------------------------------------------------------------------------|
| `amount` | integer | false    |              |                                                                                          |
| `date`   | string  | false    |              | Date is a string formatted as 2024-01-31. Timezone and time information is not included. |

## optimus-ide-collabsdk.DAUsResponse

```json
{
  "entries": [
    {
      "amount": 0,
      "date": "string"
    }
  ],
  "tz_hour_offset": 0
}
```

### Properties

| Name             | Type                                            | Required | Restrictions | Description |
|------------------|-------------------------------------------------|----------|--------------|-------------|
| `entries`        | array of [optimus-ide-collabsdk.DAUEntry](#optimus-ide-collabsdkdauentry) | false    |              |             |
| `tz_hour_offset` | integer                                         | false    |              |             |

## optimus-ide-collabsdk.DERP

```json
{
  "config": {
    "block_direct": true,
    "force_websockets": true,
    "path": "string",
    "url": "string"
  },
  "server": {
    "enable": true,
    "region_code": "string",
    "region_id": 0,
    "region_name": "string",
    "relay_url": {
      "forceQuery": true,
      "fragment": "string",
      "host": "string",
      "omitHost": true,
      "opaque": "string",
      "path": "string",
      "rawFragment": "string",
      "rawPath": "string",
      "rawQuery": "string",
      "scheme": "string",
      "user": {}
    },
    "stun_addresses": [
      "string"
    ]
  }
}
```

### Properties

| Name     | Type                                                   | Required | Restrictions | Description |
|----------|--------------------------------------------------------|----------|--------------|-------------|
| `config` | [optimus-ide-collabsdk.DERPConfig](#optimus-ide-collabsdkderpconfig)             | false    |              |             |
| `server` | [optimus-ide-collabsdk.DERPServerConfig](#optimus-ide-collabsdkderpserverconfig) | false    |              |             |

## optimus-ide-collabsdk.DERPConfig

```json
{
  "block_direct": true,
  "force_websockets": true,
  "path": "string",
  "url": "string"
}
```

### Properties

| Name               | Type    | Required | Restrictions | Description |
|--------------------|---------|----------|--------------|-------------|
| `block_direct`     | boolean | false    |              |             |
| `force_websockets` | boolean | false    |              |             |
| `path`             | string  | false    |              |             |
| `url`              | string  | false    |              |             |

## optimus-ide-collabsdk.DERPRegion

```json
{
  "latency_ms": 0,
  "preferred": true
}
```

### Properties

| Name         | Type    | Required | Restrictions | Description |
|--------------|---------|----------|--------------|-------------|
| `latency_ms` | number  | false    |              |             |
| `preferred`  | boolean | false    |              |             |

## optimus-ide-collabsdk.DERPServerConfig

```json
{
  "enable": true,
  "region_code": "string",
  "region_id": 0,
  "region_name": "string",
  "relay_url": {
    "forceQuery": true,
    "fragment": "string",
    "host": "string",
    "omitHost": true,
    "opaque": "string",
    "path": "string",
    "rawFragment": "string",
    "rawPath": "string",
    "rawQuery": "string",
    "scheme": "string",
    "user": {}
  },
  "stun_addresses": [
    "string"
  ]
}
```

### Properties

| Name             | Type                       | Required | Restrictions | Description |
|------------------|----------------------------|----------|--------------|-------------|
| `enable`         | boolean                    | false    |              |             |
| `region_code`    | string                     | false    |              |             |
| `region_id`      | integer                    | false    |              |             |
| `region_name`    | string                     | false    |              |             |
| `relay_url`      | [serpent.URL](#serpenturl) | false    |              |             |
| `stun_addresses` | array of string            | false    |              |             |

## optimus-ide-collabsdk.DangerousConfig

```json
{
  "allow_all_cors": true,
  "allow_path_app_sharing": true,
  "allow_path_app_site_owner_access": true
}
```

### Properties

| Name                               | Type    | Required | Restrictions | Description |
|------------------------------------|---------|----------|--------------|-------------|
| `allow_all_cors`                   | boolean | false    |              |             |
| `allow_path_app_sharing`           | boolean | false    |              |             |
| `allow_path_app_site_owner_access` | boolean | false    |              |             |

## optimus-ide-collabsdk.DeleteExternalAuthByIDResponse

```json
{
  "token_revocation_error": "string",
  "token_revoked": true
}
```

### Properties

| Name                     | Type    | Required | Restrictions | Description                                                                    |
|--------------------------|---------|----------|--------------|--------------------------------------------------------------------------------|
| `token_revocation_error` | string  | false    |              |                                                                                |
| `token_revoked`          | boolean | false    |              | Token revoked set to true if token revocation was attempted and was successful |

## optimus-ide-collabsdk.DeleteWebpushSubscription

```json
{
  "endpoint": "string"
}
```

### Properties

| Name       | Type   | Required | Restrictions | Description |
|------------|--------|----------|--------------|-------------|
| `endpoint` | string | false    |              |             |

## optimus-ide-collabsdk.DeleteWorkspaceAgentPortShareRequest

```json
{
  "agent_name": "string",
  "port": 0
}
```

### Properties

| Name         | Type    | Required | Restrictions | Description |
|--------------|---------|----------|--------------|-------------|
| `agent_name` | string  | false    |              |             |
| `port`       | integer | false    |              |             |

## optimus-ide-collabsdk.DeploymentConfig

```json
{
  "config": {
    "access_url": {
      "forceQuery": true,
      "fragment": "string",
      "host": "string",
      "omitHost": true,
      "opaque": "string",
      "path": "string",
      "rawFragment": "string",
      "rawPath": "string",
      "rawQuery": "string",
      "scheme": "string",
      "user": {}
    },
    "additional_csp_policy": [
      "string"
    ],
    "address": {
      "host": "string",
      "port": "string"
    },
    "agent_fallback_troubleshooting_url": {
      "forceQuery": true,
      "fragment": "string",
      "host": "string",
      "omitHost": true,
      "opaque": "string",
      "path": "string",
      "rawFragment": "string",
      "rawPath": "string",
      "rawQuery": "string",
      "scheme": "string",
      "user": {}
    },
    "agent_stat_refresh_interval": 0,
    "ai": {
      "aibridge_proxy": {
        "allowed_private_cidrs": [
          "string"
        ],
        "api_dump_dir": "string",
        "cert_file": "string",
        "domain_allowlist": [
          "string"
        ],
        "enabled": true,
        "key_file": "string",
        "listen_addr": "string",
        "target": "string",
        "tls_cert_file": "string",
        "tls_key_file": "string",
        "upstream_proxy": "string",
        "upstream_proxy_ca": "string"
      },
      "bridge": {
        "allow_byok": true,
        "anthropic": {
          "base_url": "string",
          "key": "string"
        },
        "api_dump_dir": "string",
        "bedrock": {
          "access_key": "string",
          "access_key_secret": "string",
          "base_url": "string",
          "model": "string",
          "region": "string",
          "small_fast_model": "string"
        },
        "budget_period": "string",
        "budget_policy": "string",
        "circuit_breaker_enabled": true,
        "circuit_breaker_failure_threshold": 0,
        "circuit_breaker_interval": 0,
        "circuit_breaker_max_requests": 0,
        "circuit_breaker_timeout": 0,
        "enabled": true,
        "inject_optimus-ide-collab_mcp_tools": true,
        "max_concurrency": 0,
        "openai": {
          "base_url": "string",
          "key": "string"
        },
        "providers": [
          {
            "base_url": "string",
            "bedrock_model": "string",
            "bedrock_region": "string",
            "bedrock_small_fast_model": "string",
            "name": "string",
            "type": "string"
          }
        ],
        "rate_limit": 0,
        "retention": 0,
        "send_actor_headers": true,
        "structured_logging": true
      },
      "chat": {
        "acquire_batch_size": 0,
        "debug_logging_enabled": true
      }
    },
    "allow_workspace_renames": true,
    "autobuild_poll_interval": 0,
    "browser_only": true,
    "cache_directory": "string",
    "cli_upgrade_message": "string",
    "cluster": {
      "host": "string"
    },
    "config": "string",
    "config_ssh": {
      "deploymentName": "string",
      "sshconfigOptions": [
        "string"
      ]
    },
    "dangerous": {
      "allow_all_cors": true,
      "allow_path_app_sharing": true,
      "allow_path_app_site_owner_access": true
    },
    "derp": {
      "config": {
        "block_direct": true,
        "force_websockets": true,
        "path": "string",
        "url": "string"
      },
      "server": {
        "enable": true,
        "region_code": "string",
        "region_id": 0,
        "region_name": "string",
        "relay_url": {
          "forceQuery": true,
          "fragment": "string",
          "host": "string",
          "omitHost": true,
          "opaque": "string",
          "path": "string",
          "rawFragment": "string",
          "rawPath": "string",
          "rawQuery": "string",
          "scheme": "string",
          "user": {}
        },
        "stun_addresses": [
          "string"
        ]
      }
    },
    "disable_chat_sharing": true,
    "disable_owner_workspace_exec": true,
    "disable_password_auth": true,
    "disable_path_apps": true,
    "disable_workspace_sharing": true,
    "docs_url": {
      "forceQuery": true,
      "fragment": "string",
      "host": "string",
      "omitHost": true,
      "opaque": "string",
      "path": "string",
      "rawFragment": "string",
      "rawPath": "string",
      "rawQuery": "string",
      "scheme": "string",
      "user": {}
    },
    "enable_authz_recording": true,
    "enable_terraform_debug_mode": true,
    "ephemeral_deployment": true,
    "experiments": [
      "string"
    ],
    "external_auth": {
      "value": [
        {
          "api_base_url": "string",
          "app_install_url": "string",
          "app_installations_url": "string",
          "auth_url": "string",
          "client_id": "string",
          "code_challenge_methods_supported": [
            "string"
          ],
          "device_code_url": "string",
          "device_flow": true,
          "display_icon": "string",
          "display_name": "string",
          "id": "string",
          "mcp_tool_allow_regex": "string",
          "mcp_tool_deny_regex": "string",
          "mcp_url": "string",
          "no_refresh": true,
          "regex": "string",
          "revoke_url": "string",
          "scopes": [
            "string"
          ],
          "token_url": "string",
          "type": "string",
          "validate_url": "string"
        }
      ]
    },
    "external_auth_github_default_provider_enable": true,
    "external_token_encryption_keys": [
      "string"
    ],
    "healthcheck": {
      "refresh": 0,
      "threshold_database": 0
    },
    "hide_ai_tasks": true,
    "http_address": "string",
    "http_cookies": {
      "host_prefix": true,
      "same_site": "string",
      "secure_auth_cookie": true
    },
    "job_hang_detector_interval": 0,
    "logging": {
      "human": "string",
      "json": "string",
      "log_filter": [
        "string"
      ],
      "stackdriver": "string"
    },
    "metrics_cache_refresh_interval": 0,
    "notifications": {
      "dispatch_timeout": 0,
      "email": {
        "auth": {
          "identity": "string",
          "password": "string",
          "password_file": "string",
          "username": "string"
        },
        "force_tls": true,
        "from": "string",
        "hello": "string",
        "smarthost": "string",
        "tls": {
          "ca_file": "string",
          "cert_file": "string",
          "insecure_skip_verify": true,
          "key_file": "string",
          "server_name": "string",
          "start_tls": true
        }
      },
      "fetch_interval": 0,
      "inbox": {
        "enabled": true
      },
      "lease_count": 0,
      "lease_period": 0,
      "max_send_attempts": 0,
      "method": "string",
      "retry_interval": 0,
      "sync_buffer_size": 0,
      "sync_interval": 0,
      "webhook": {
        "endpoint": {
          "forceQuery": true,
          "fragment": "string",
          "host": "string",
          "omitHost": true,
          "opaque": "string",
          "path": "string",
          "rawFragment": "string",
          "rawPath": "string",
          "rawQuery": "string",
          "scheme": "string",
          "user": {}
        }
      }
    },
    "oauth2": {
      "github": {
        "allow_everyone": true,
        "allow_signups": true,
        "allowed_orgs": [
          "string"
        ],
        "allowed_teams": [
          "string"
        ],
        "client_id": "string",
        "client_secret": "string",
        "default_provider_enable": true,
        "device_flow": true,
        "enterprise_base_url": "string"
      }
    },
    "oidc": {
      "allow_signups": true,
      "auth_url_params": {},
      "auto_repair_links": true,
      "client_cert_file": "string",
      "client_id": "string",
      "client_key_file": "string",
      "client_secret": "string",
      "email_domain": [
        "string"
      ],
      "email_fallback": true,
      "email_field": "string",
      "group_allow_list": [
        "string"
      ],
      "group_auto_create": true,
      "group_mapping": {},
      "group_regex_filter": {},
      "groups_field": "string",
      "icon_url": {
        "forceQuery": true,
        "fragment": "string",
        "host": "string",
        "omitHost": true,
        "opaque": "string",
        "path": "string",
        "rawFragment": "string",
        "rawPath": "string",
        "rawQuery": "string",
        "scheme": "string",
        "user": {}
      },
      "ignore_email_verified": true,
      "ignore_user_info": true,
      "issuer_url": "string",
      "name_field": "string",
      "organization_assign_default": true,
      "organization_field": "string",
      "organization_mapping": {},
      "redirect_allowed_hosts": [
        "string"
      ],
      "redirect_url": {
        "forceQuery": true,
        "fragment": "string",
        "host": "string",
        "omitHost": true,
        "opaque": "string",
        "path": "string",
        "rawFragment": "string",
        "rawPath": "string",
        "rawQuery": "string",
        "scheme": "string",
        "user": {}
      },
      "scopes": [
        "string"
      ],
      "sign_in_text": "string",
      "signups_disabled_text": "string",
      "skip_issuer_checks": true,
      "source_user_info_from_access_token": true,
      "user_role_field": "string",
      "user_role_mapping": {},
      "user_roles_default": [
        "string"
      ],
      "username_field": "string"
    },
    "pg_auth": "string",
    "pg_conn_max_idle": "string",
    "pg_conn_max_open": 0,
    "pg_connection_url": "string",
    "pprof": {
      "address": {
        "host": "string",
        "port": "string"
      },
      "enable": true
    },
    "prometheus": {
      "address": {
        "host": "string",
        "port": "string"
      },
      "aggregate_agent_stats_by": [
        "string"
      ],
      "collect_agent_stats": true,
      "collect_db_metrics": true,
      "enable": true
    },
    "provisioner": {
      "daemon_poll_interval": 0,
      "daemon_poll_jitter": 0,
      "daemon_psk": "string",
      "daemon_types": [
        "string"
      ],
      "daemons": 0,
      "force_cancel_interval": 0
    },
    "proxy_health_status_interval": 0,
    "proxy_trusted_headers": [
      "string"
    ],
    "proxy_trusted_origins": [
      "string"
    ],
    "rate_limit": {
      "api": 0,
      "disable_all": true
    },
    "redirect_to_access_url": true,
    "retention": {
      "api_keys": 0,
      "audit_logs": 0,
      "boundary_logs": 0,
      "connection_logs": 0,
      "workspace_agent_logs": 0
    },
    "scim_api_key": "string",
    "scim_use_legacy": true,
    "session_lifetime": {
      "default_duration": 0,
      "default_token_lifetime": 0,
      "disable_expiry_refresh": true,
      "max_admin_token_lifetime": 0,
      "max_token_lifetime": 0,
      "refresh_default_duration": 0
    },
    "ssh_keygen_algorithm": "string",
    "stats_collection": {
      "usage_stats": {
        "enable": true
      }
    },
    "strict_transport_security": 0,
    "strict_transport_security_options": [
      "string"
    ],
    "support": {
      "links": {
        "value": [
          {
            "icon": "bug",
            "location": "navbar",
            "name": "string",
            "target": "string"
          }
        ]
      }
    },
    "swagger": {
      "enable": true
    },
    "telemetry": {
      "enable": true,
      "trace": true,
      "url": {
        "forceQuery": true,
        "fragment": "string",
        "host": "string",
        "omitHost": true,
        "opaque": "string",
        "path": "string",
        "rawFragment": "string",
        "rawPath": "string",
        "rawQuery": "string",
        "scheme": "string",
        "user": {}
      }
    },
    "template_builder": {
      "disabled": true,
      "registry_url": "string"
    },
    "terms_of_service_url": "string",
    "tls": {
      "address": {
        "host": "string",
        "port": "string"
      },
      "allow_insecure_ciphers": true,
      "cert_file": [
        "string"
      ],
      "client_auth": "string",
      "client_ca_file": "string",
      "client_cert_file": "string",
      "client_key_file": "string",
      "enable": true,
      "key_file": [
        "string"
      ],
      "min_version": "string",
      "redirect_http": true,
      "supported_ciphers": [
        "string"
      ]
    },
    "trace": {
      "capture_logs": true,
      "data_dog": true,
      "enable": true,
      "honeycomb_api_key": "string"
    },
    "update_check": true,
    "user_quiet_hours_schedule": {
      "allow_user_custom": true,
      "default_schedule": "string"
    },
    "verbose": true,
    "web_terminal_renderer": "string",
    "wgtunnel_host": "string",
    "wildcard_access_url": "string",
    "workspace_hostname_suffix": "string",
    "workspace_prebuilds": {
      "failure_hard_limit": 0,
      "reconciliation_backoff_interval": 0,
      "reconciliation_backoff_lookback": 0,
      "reconciliation_interval": 0
    },
    "write_config": true
  },
  "options": [
    {
      "annotations": {
        "property1": "string",
        "property2": "string"
      },
      "default": "string",
      "description": "string",
      "env": "string",
      "flag": "string",
      "flag_shorthand": "string",
      "group": {
        "description": "string",
        "name": "string",
        "parent": {
          "description": "string",
          "name": "string",
          "parent": {},
          "yaml": "string"
        },
        "yaml": "string"
      },
      "hidden": true,
      "name": "string",
      "required": true,
      "use_instead": [
        {}
      ],
      "value": null,
      "value_source": "",
      "yaml": "string"
    }
  ]
}
```

### Properties

| Name      | Type                                                   | Required | Restrictions | Description |
|-----------|--------------------------------------------------------|----------|--------------|-------------|
| `config`  | [optimus-ide-collabsdk.DeploymentValues](#optimus-ide-collabsdkdeploymentvalues) | false    |              |             |
| `options` | array of [serpent.Option](#serpentoption)              | false    |              |             |

## optimus-ide-collabsdk.DeploymentStats

```json
{
  "aggregated_from": "2019-08-24T14:15:22Z",
  "collected_at": "2019-08-24T14:15:22Z",
  "next_update_at": "2019-08-24T14:15:22Z",
  "session_count": {
    "jetbrains": 0,
    "reconnecting_pty": 0,
    "ssh": 0,
    "vscode": 0
  },
  "workspaces": {
    "building": 0,
    "connection_latency_ms": {
      "p50": 0,
      "p95": 0
    },
    "failed": 0,
    "pending": 0,
    "running": 0,
    "rx_bytes": 0,
    "stopped": 0,
    "tx_bytes": 0
  }
}
```

### Properties

| Name              | Type                                                                         | Required | Restrictions | Description                                                                                                                 |
|-------------------|------------------------------------------------------------------------------|----------|--------------|-----------------------------------------------------------------------------------------------------------------------------|
| `aggregated_from` | string                                                                       | false    |              | Aggregated from is the time in which stats are aggregated from. This might be back in time a specific duration or interval. |
| `collected_at`    | string                                                                       | false    |              | Collected at is the time in which stats are collected at.                                                                   |
| `next_update_at`  | string                                                                       | false    |              | Next update at is the time when the next batch of stats will be updated.                                                    |
| `session_count`   | [optimus-ide-collabsdk.SessionCountDeploymentStats](#optimus-ide-collabsdksessioncountdeploymentstats) | false    |              |                                                                                                                             |
| `workspaces`      | [optimus-ide-collabsdk.WorkspaceDeploymentStats](#optimus-ide-collabsdkworkspacedeploymentstats)       | false    |              |                                                                                                                             |

## optimus-ide-collabsdk.DeploymentValues

```json
{
  "access_url": {
    "forceQuery": true,
    "fragment": "string",
    "host": "string",
    "omitHost": true,
    "opaque": "string",
    "path": "string",
    "rawFragment": "string",
    "rawPath": "string",
    "rawQuery": "string",
    "scheme": "string",
    "user": {}
  },
  "additional_csp_policy": [
    "string"
  ],
  "address": {
    "host": "string",
    "port": "string"
  },
  "agent_fallback_troubleshooting_url": {
    "forceQuery": true,
    "fragment": "string",
    "host": "string",
    "omitHost": true,
    "opaque": "string",
    "path": "string",
    "rawFragment": "string",
    "rawPath": "string",
    "rawQuery": "string",
    "scheme": "string",
    "user": {}
  },
  "agent_stat_refresh_interval": 0,
  "ai": {
    "aibridge_proxy": {
      "allowed_private_cidrs": [
        "string"
      ],
      "api_dump_dir": "string",
      "cert_file": "string",
      "domain_allowlist": [
        "string"
      ],
      "enabled": true,
      "key_file": "string",
      "listen_addr": "string",
      "target": "string",
      "tls_cert_file": "string",
      "tls_key_file": "string",
      "upstream_proxy": "string",
      "upstream_proxy_ca": "string"
    },
    "bridge": {
      "allow_byok": true,
      "anthropic": {
        "base_url": "string",
        "key": "string"
      },
      "api_dump_dir": "string",
      "bedrock": {
        "access_key": "string",
        "access_key_secret": "string",
        "base_url": "string",
        "model": "string",
        "region": "string",
        "small_fast_model": "string"
      },
      "budget_period": "string",
      "budget_policy": "string",
      "circuit_breaker_enabled": true,
      "circuit_breaker_failure_threshold": 0,
      "circuit_breaker_interval": 0,
      "circuit_breaker_max_requests": 0,
      "circuit_breaker_timeout": 0,
      "enabled": true,
      "inject_optimus-ide-collab_mcp_tools": true,
      "max_concurrency": 0,
      "openai": {
        "base_url": "string",
        "key": "string"
      },
      "providers": [
        {
          "base_url": "string",
          "bedrock_model": "string",
          "bedrock_region": "string",
          "bedrock_small_fast_model": "string",
          "name": "string",
          "type": "string"
        }
      ],
      "rate_limit": 0,
      "retention": 0,
      "send_actor_headers": true,
      "structured_logging": true
    },
    "chat": {
      "acquire_batch_size": 0,
      "debug_logging_enabled": true
    }
  },
  "allow_workspace_renames": true,
  "autobuild_poll_interval": 0,
  "browser_only": true,
  "cache_directory": "string",
  "cli_upgrade_message": "string",
  "cluster": {
    "host": "string"
  },
  "config": "string",
  "config_ssh": {
    "deploymentName": "string",
    "sshconfigOptions": [
      "string"
    ]
  },
  "dangerous": {
    "allow_all_cors": true,
    "allow_path_app_sharing": true,
    "allow_path_app_site_owner_access": true
  },
  "derp": {
    "config": {
      "block_direct": true,
      "force_websockets": true,
      "path": "string",
      "url": "string"
    },
    "server": {
      "enable": true,
      "region_code": "string",
      "region_id": 0,
      "region_name": "string",
      "relay_url": {
        "forceQuery": true,
        "fragment": "string",
        "host": "string",
        "omitHost": true,
        "opaque": "string",
        "path": "string",
        "rawFragment": "string",
        "rawPath": "string",
        "rawQuery": "string",
        "scheme": "string",
        "user": {}
      },
      "stun_addresses": [
        "string"
      ]
    }
  },
  "disable_chat_sharing": true,
  "disable_owner_workspace_exec": true,
  "disable_password_auth": true,
  "disable_path_apps": true,
  "disable_workspace_sharing": true,
  "docs_url": {
    "forceQuery": true,
    "fragment": "string",
    "host": "string",
    "omitHost": true,
    "opaque": "string",
    "path": "string",
    "rawFragment": "string",
    "rawPath": "string",
    "rawQuery": "string",
    "scheme": "string",
    "user": {}
  },
  "enable_authz_recording": true,
  "enable_terraform_debug_mode": true,
  "ephemeral_deployment": true,
  "experiments": [
    "string"
  ],
  "external_auth": {
    "value": [
      {
        "api_base_url": "string",
        "app_install_url": "string",
        "app_installations_url": "string",
        "auth_url": "string",
        "client_id": "string",
        "code_challenge_methods_supported": [
          "string"
        ],
        "device_code_url": "string",
        "device_flow": true,
        "display_icon": "string",
        "display_name": "string",
        "id": "string",
        "mcp_tool_allow_regex": "string",
        "mcp_tool_deny_regex": "string",
        "mcp_url": "string",
        "no_refresh": true,
        "regex": "string",
        "revoke_url": "string",
        "scopes": [
          "string"
        ],
        "token_url": "string",
        "type": "string",
        "validate_url": "string"
      }
    ]
  },
  "external_auth_github_default_provider_enable": true,
  "external_token_encryption_keys": [
    "string"
  ],
  "healthcheck": {
    "refresh": 0,
    "threshold_database": 0
  },
  "hide_ai_tasks": true,
  "http_address": "string",
  "http_cookies": {
    "host_prefix": true,
    "same_site": "string",
    "secure_auth_cookie": true
  },
  "job_hang_detector_interval": 0,
  "logging": {
    "human": "string",
    "json": "string",
    "log_filter": [
      "string"
    ],
    "stackdriver": "string"
  },
  "metrics_cache_refresh_interval": 0,
  "notifications": {
    "dispatch_timeout": 0,
    "email": {
      "auth": {
        "identity": "string",
        "password": "string",
        "password_file": "string",
        "username": "string"
      },
      "force_tls": true,
      "from": "string",
      "hello": "string",
      "smarthost": "string",
      "tls": {
        "ca_file": "string",
        "cert_file": "string",
        "insecure_skip_verify": true,
        "key_file": "string",
        "server_name": "string",
        "start_tls": true
      }
    },
    "fetch_interval": 0,
    "inbox": {
      "enabled": true
    },
    "lease_count": 0,
    "lease_period": 0,
    "max_send_attempts": 0,
    "method": "string",
    "retry_interval": 0,
    "sync_buffer_size": 0,
    "sync_interval": 0,
    "webhook": {
      "endpoint": {
        "forceQuery": true,
        "fragment": "string",
        "host": "string",
        "omitHost": true,
        "opaque": "string",
        "path": "string",
        "rawFragment": "string",
        "rawPath": "string",
        "rawQuery": "string",
        "scheme": "string",
        "user": {}
      }
    }
  },
  "oauth2": {
    "github": {
      "allow_everyone": true,
      "allow_signups": true,
      "allowed_orgs": [
        "string"
      ],
      "allowed_teams": [
        "string"
      ],
      "client_id": "string",
      "client_secret": "string",
      "default_provider_enable": true,
      "device_flow": true,
      "enterprise_base_url": "string"
    }
  },
  "oidc": {
    "allow_signups": true,
    "auth_url_params": {},
    "auto_repair_links": true,
    "client_cert_file": "string",
    "client_id": "string",
    "client_key_file": "string",
    "client_secret": "string",
    "email_domain": [
      "string"
    ],
    "email_fallback": true,
    "email_field": "string",
    "group_allow_list": [
      "string"
    ],
    "group_auto_create": true,
    "group_mapping": {},
    "group_regex_filter": {},
    "groups_field": "string",
    "icon_url": {
      "forceQuery": true,
      "fragment": "string",
      "host": "string",
      "omitHost": true,
      "opaque": "string",
      "path": "string",
      "rawFragment": "string",
      "rawPath": "string",
      "rawQuery": "string",
      "scheme": "string",
      "user": {}
    },
    "ignore_email_verified": true,
    "ignore_user_info": true,
    "issuer_url": "string",
    "name_field": "string",
    "organization_assign_default": true,
    "organization_field": "string",
    "organization_mapping": {},
    "redirect_allowed_hosts": [
      "string"
    ],
    "redirect_url": {
      "forceQuery": true,
      "fragment": "string",
      "host": "string",
      "omitHost": true,
      "opaque": "string",
      "path": "string",
      "rawFragment": "string",
      "rawPath": "string",
      "rawQuery": "string",
      "scheme": "string",
      "user": {}
    },
    "scopes": [
      "string"
    ],
    "sign_in_text": "string",
    "signups_disabled_text": "string",
    "skip_issuer_checks": true,
    "source_user_info_from_access_token": true,
    "user_role_field": "string",
    "user_role_mapping": {},
    "user_roles_default": [
      "string"
    ],
    "username_field": "string"
  },
  "pg_auth": "string",
  "pg_conn_max_idle": "string",
  "pg_conn_max_open": 0,
  "pg_connection_url": "string",
  "pprof": {
    "address": {
      "host": "string",
      "port": "string"
    },
    "enable": true
  },
  "prometheus": {
    "address": {
      "host": "string",
      "port": "string"
    },
    "aggregate_agent_stats_by": [
      "string"
    ],
    "collect_agent_stats": true,
    "collect_db_metrics": true,
    "enable": true
  },
  "provisioner": {
    "daemon_poll_interval": 0,
    "daemon_poll_jitter": 0,
    "daemon_psk": "string",
    "daemon_types": [
      "string"
    ],
    "daemons": 0,
    "force_cancel_interval": 0
  },
  "proxy_health_status_interval": 0,
  "proxy_trusted_headers": [
    "string"
  ],
  "proxy_trusted_origins": [
    "string"
  ],
  "rate_limit": {
    "api": 0,
    "disable_all": true
  },
  "redirect_to_access_url": true,
  "retention": {
    "api_keys": 0,
    "audit_logs": 0,
    "boundary_logs": 0,
    "connection_logs": 0,
    "workspace_agent_logs": 0
  },
  "scim_api_key": "string",
  "scim_use_legacy": true,
  "session_lifetime": {
    "default_duration": 0,
    "default_token_lifetime": 0,
    "disable_expiry_refresh": true,
    "max_admin_token_lifetime": 0,
    "max_token_lifetime": 0,
    "refresh_default_duration": 0
  },
  "ssh_keygen_algorithm": "string",
  "stats_collection": {
    "usage_stats": {
      "enable": true
    }
  },
  "strict_transport_security": 0,
  "strict_transport_security_options": [
    "string"
  ],
  "support": {
    "links": {
      "value": [
        {
          "icon": "bug",
          "location": "navbar",
          "name": "string",
          "target": "string"
        }
      ]
    }
  },
  "swagger": {
    "enable": true
  },
  "telemetry": {
    "enable": true,
    "trace": true,
    "url": {
      "forceQuery": true,
      "fragment": "string",
      "host": "string",
      "omitHost": true,
      "opaque": "string",
      "path": "string",
      "rawFragment": "string",
      "rawPath": "string",
      "rawQuery": "string",
      "scheme": "string",
      "user": {}
    }
  },
  "template_builder": {
    "disabled": true,
    "registry_url": "string"
  },
  "terms_of_service_url": "string",
  "tls": {
    "address": {
      "host": "string",
      "port": "string"
    },
    "allow_insecure_ciphers": true,
    "cert_file": [
      "string"
    ],
    "client_auth": "string",
    "client_ca_file": "string",
    "client_cert_file": "string",
    "client_key_file": "string",
    "enable": true,
    "key_file": [
      "string"
    ],
    "min_version": "string",
    "redirect_http": true,
    "supported_ciphers": [
      "string"
    ]
  },
  "trace": {
    "capture_logs": true,
    "data_dog": true,
    "enable": true,
    "honeycomb_api_key": "string"
  },
  "update_check": true,
  "user_quiet_hours_schedule": {
    "allow_user_custom": true,
    "default_schedule": "string"
  },
  "verbose": true,
  "web_terminal_renderer": "string",
  "wgtunnel_host": "string",
  "wildcard_access_url": "string",
  "workspace_hostname_suffix": "string",
  "workspace_prebuilds": {
    "failure_hard_limit": 0,
    "reconciliation_backoff_interval": 0,
    "reconciliation_backoff_lookback": 0,
    "reconciliation_interval": 0
  },
  "write_config": true
}
```

### Properties

| Name                                           | Type                                                                                                 | Required | Restrictions | Description                                                        |
|------------------------------------------------|------------------------------------------------------------------------------------------------------|----------|--------------|--------------------------------------------------------------------|
| `access_url`                                   | [serpent.URL](#serpenturl)                                                                           | false    |              |                                                                    |
| `additional_csp_policy`                        | array of string                                                                                      | false    |              |                                                                    |
| `address`                                      | [serpent.HostPort](#serpenthostport)                                                                 | false    |              | Deprecated: Use HTTPAddress or TLS.Address instead.                |
| `agent_fallback_troubleshooting_url`           | [serpent.URL](#serpenturl)                                                                           | false    |              |                                                                    |
| `agent_stat_refresh_interval`                  | integer                                                                                              | false    |              |                                                                    |
| `ai`                                           | [optimus-ide-collabsdk.AIConfig](#optimus-ide-collabsdkaiconfig)                                                               | false    |              |                                                                    |
| `allow_workspace_renames`                      | boolean                                                                                              | false    |              |                                                                    |
| `autobuild_poll_interval`                      | integer                                                                                              | false    |              |                                                                    |
| `browser_only`                                 | boolean                                                                                              | false    |              |                                                                    |
| `cache_directory`                              | string                                                                                               | false    |              |                                                                    |
| `cli_upgrade_message`                          | string                                                                                               | false    |              |                                                                    |
| `cluster`                                      | [optimus-ide-collabsdk.ClusterConfig](#optimus-ide-collabsdkclusterconfig)                                                     | false    |              |                                                                    |
| `config`                                       | string                                                                                               | false    |              |                                                                    |
| `config_ssh`                                   | [optimus-ide-collabsdk.SSHConfig](#optimus-ide-collabsdksshconfig)                                                             | false    |              |                                                                    |
| `dangerous`                                    | [optimus-ide-collabsdk.DangerousConfig](#optimus-ide-collabsdkdangerousconfig)                                                 | false    |              |                                                                    |
| `derp`                                         | [optimus-ide-collabsdk.DERP](#optimus-ide-collabsdkderp)                                                                       | false    |              |                                                                    |
| `disable_chat_sharing`                         | boolean                                                                                              | false    |              |                                                                    |
| `disable_owner_workspace_exec`                 | boolean                                                                                              | false    |              |                                                                    |
| `disable_password_auth`                        | boolean                                                                                              | false    |              |                                                                    |
| `disable_path_apps`                            | boolean                                                                                              | false    |              |                                                                    |
| `disable_workspace_sharing`                    | boolean                                                                                              | false    |              |                                                                    |
| `docs_url`                                     | [serpent.URL](#serpenturl)                                                                           | false    |              |                                                                    |
| `enable_authz_recording`                       | boolean                                                                                              | false    |              |                                                                    |
| `enable_terraform_debug_mode`                  | boolean                                                                                              | false    |              |                                                                    |
| `ephemeral_deployment`                         | boolean                                                                                              | false    |              |                                                                    |
| `experiments`                                  | array of string                                                                                      | false    |              |                                                                    |
| `external_auth`                                | [serpent.Struct-array_optimus-ide-collabsdk_ExternalAuthConfig](#serpentstruct-array_optimus-ide-collabsdk_externalauthconfig) | false    |              |                                                                    |
| `external_auth_github_default_provider_enable` | boolean                                                                                              | false    |              |                                                                    |
| `external_token_encryption_keys`               | array of string                                                                                      | false    |              |                                                                    |
| `healthcheck`                                  | [optimus-ide-collabsdk.HealthcheckConfig](#optimus-ide-collabsdkhealthcheckconfig)                                             | false    |              |                                                                    |
| `hide_ai_tasks`                                | boolean                                                                                              | false    |              |                                                                    |
| `http_address`                                 | string                                                                                               | false    |              | Http address is a string because it may be set to zero to disable. |
| `http_cookies`                                 | [optimus-ide-collabsdk.HTTPCookieConfig](#optimus-ide-collabsdkhttpcookieconfig)                                               | false    |              |                                                                    |
| `job_hang_detector_interval`                   | integer                                                                                              | false    |              |                                                                    |
| `logging`                                      | [optimus-ide-collabsdk.LoggingConfig](#optimus-ide-collabsdkloggingconfig)                                                     | false    |              |                                                                    |
| `metrics_cache_refresh_interval`               | integer                                                                                              | false    |              |                                                                    |
| `notifications`                                | [optimus-ide-collabsdk.NotificationsConfig](#optimus-ide-collabsdknotificationsconfig)                                         | false    |              |                                                                    |
| `oauth2`                                       | [optimus-ide-collabsdk.OAuth2Config](#optimus-ide-collabsdkoauth2config)                                                       | false    |              |                                                                    |
| `oidc`                                         | [optimus-ide-collabsdk.OIDCConfig](#optimus-ide-collabsdkoidcconfig)                                                           | false    |              |                                                                    |
| `pg_auth`                                      | string                                                                                               | false    |              |                                                                    |
| `pg_conn_max_idle`                             | string                                                                                               | false    |              |                                                                    |
| `pg_conn_max_open`                             | integer                                                                                              | false    |              |                                                                    |
| `pg_connection_url`                            | string                                                                                               | false    |              |                                                                    |
| `pprof`                                        | [optimus-ide-collabsdk.PprofConfig](#optimus-ide-collabsdkpprofconfig)                                                         | false    |              |                                                                    |
| `prometheus`                                   | [optimus-ide-collabsdk.PrometheusConfig](#optimus-ide-collabsdkprometheusconfig)                                               | false    |              |                                                                    |
| `provisioner`                                  | [optimus-ide-collabsdk.ProvisionerConfig](#optimus-ide-collabsdkprovisionerconfig)                                             | false    |              |                                                                    |
| `proxy_health_status_interval`                 | integer                                                                                              | false    |              |                                                                    |
| `proxy_trusted_headers`                        | array of string                                                                                      | false    |              |                                                                    |
| `proxy_trusted_origins`                        | array of string                                                                                      | false    |              |                                                                    |
| `rate_limit`                                   | [optimus-ide-collabsdk.RateLimitConfig](#optimus-ide-collabsdkratelimitconfig)                                                 | false    |              |                                                                    |
| `redirect_to_access_url`                       | boolean                                                                                              | false    |              |                                                                    |
| `retention`                                    | [optimus-ide-collabsdk.RetentionConfig](#optimus-ide-collabsdkretentionconfig)                                                 | false    |              |                                                                    |
| `scim_api_key`                                 | string                                                                                               | false    |              |                                                                    |
| `scim_use_legacy`                              | boolean                                                                                              | false    |              |                                                                    |
| `session_lifetime`                             | [optimus-ide-collabsdk.SessionLifetime](#optimus-ide-collabsdksessionlifetime)                                                 | false    |              |                                                                    |
| `ssh_keygen_algorithm`                         | string                                                                                               | false    |              |                                                                    |
| `stats_collection`                             | [optimus-ide-collabsdk.StatsCollectionConfig](#optimus-ide-collabsdkstatscollectionconfig)                                     | false    |              |                                                                    |
| `strict_transport_security`                    | integer                                                                                              | false    |              |                                                                    |
| `strict_transport_security_options`            | array of string                                                                                      | false    |              |                                                                    |
| `support`                                      | [optimus-ide-collabsdk.SupportConfig](#optimus-ide-collabsdksupportconfig)                                                     | false    |              |                                                                    |
| `swagger`                                      | [optimus-ide-collabsdk.SwaggerConfig](#optimus-ide-collabsdkswaggerconfig)                                                     | false    |              |                                                                    |
| `telemetry`                                    | [optimus-ide-collabsdk.TelemetryConfig](#optimus-ide-collabsdktelemetryconfig)                                                 | false    |              |                                                                    |
| `template_builder`                             | [optimus-ide-collabsdk.TemplateBuilderConfig](#optimus-ide-collabsdktemplatebuilderconfig)                                     | false    |              |                                                                    |
| `terms_of_service_url`                         | string                                                                                               | false    |              |                                                                    |
| `tls`                                          | [optimus-ide-collabsdk.TLSConfig](#optimus-ide-collabsdktlsconfig)                                                             | false    |              |                                                                    |
| `trace`                                        | [optimus-ide-collabsdk.TraceConfig](#optimus-ide-collabsdktraceconfig)                                                         | false    |              |                                                                    |
| `update_check`                                 | boolean                                                                                              | false    |              |                                                                    |
| `user_quiet_hours_schedule`                    | [optimus-ide-collabsdk.UserQuietHoursScheduleConfig](#optimus-ide-collabsdkuserquiethoursscheduleconfig)                       | false    |              |                                                                    |
| `verbose`                                      | boolean                                                                                              | false    |              |                                                                    |
| `web_terminal_renderer`                        | string                                                                                               | false    |              |                                                                    |
| `wgtunnel_host`                                | string                                                                                               | false    |              |                                                                    |
| `wildcard_access_url`                          | string                                                                                               | false    |              |                                                                    |
| `workspace_hostname_suffix`                    | string                                                                                               | false    |              |                                                                    |
| `workspace_prebuilds`                          | [optimus-ide-collabsdk.PrebuildsConfig](#optimus-ide-collabsdkprebuildsconfig)                                                 | false    |              |                                                                    |
| `write_config`                                 | boolean                                                                                              | false    |              |                                                                    |

## optimus-ide-collabsdk.DiagnosticExtra

```json
{
  "code": "string"
}
```

### Properties

| Name   | Type   | Required | Restrictions | Description |
|--------|--------|----------|--------------|-------------|
| `code` | string | false    |              |             |

## optimus-ide-collabsdk.DiagnosticSeverityString

```json
"error"
```

### Properties

#### Enumerated Values

| Value(s)           |
|--------------------|
| `error`, `warning` |

## optimus-ide-collabsdk.DisplayApp

```json
"vscode"
```

### Properties

#### Enumerated Values

| Value(s)                                                                            |
|-------------------------------------------------------------------------------------|
| `port_forwarding_helper`, `ssh_helper`, `vscode`, `vscode_insiders`, `web_terminal` |

## optimus-ide-collabsdk.DynamicParametersRequest

```json
{
  "id": 0,
  "inputs": {
    "property1": "string",
    "property2": "string"
  },
  "owner_id": "8826ee2e-7933-4665-aef2-2393f84a0d05"
}
```

### Properties

| Name               | Type    | Required | Restrictions | Description                                                                                                  |
|--------------------|---------|----------|--------------|--------------------------------------------------------------------------------------------------------------|
| `id`               | integer | false    |              | ID identifies the request. The response contains the same ID so that the client can match it to the request. |
| `inputs`           | object  | false    |              |                                                                                                              |
| » `[any property]` | string  | false    |              |                                                                                                              |
| `owner_id`         | string  | false    |              | Owner ID if uuid.Nil, it defaults to `optimus-ide-collabsdk.Me`                                                           |

## optimus-ide-collabsdk.DynamicParametersResponse

```json
{
  "diagnostics": [
    {
      "detail": "string",
      "extra": {
        "code": "string"
      },
      "severity": "error",
      "summary": "string"
    }
  ],
  "id": 0,
  "parameters": [
    {
      "default_value": {
        "valid": true,
        "value": "string"
      },
      "description": "string",
      "diagnostics": [
        {
          "detail": "string",
          "extra": {
            "code": "string"
          },
          "severity": "error",
          "summary": "string"
        }
      ],
      "display_name": "string",
      "ephemeral": true,
      "form_type": "",
      "icon": "string",
      "mutable": true,
      "name": "string",
      "options": [
        {
          "description": "string",
          "icon": "string",
          "name": "string",
          "value": {
            "valid": true,
            "value": "string"
          }
        }
      ],
      "order": 0,
      "required": true,
      "styling": {
        "disabled": true,
        "label": "string",
        "mask_input": true,
        "placeholder": "string"
      },
      "type": "string",
      "validations": [
        {
          "validation_error": "string",
          "validation_max": 0,
          "validation_min": 0,
          "validation_monotonic": "string",
          "validation_regex": "string"
        }
      ],
      "value": {
        "valid": true,
        "value": "string"
      }
    }
  ]
}
```

### Properties

| Name          | Type                                                                | Required | Restrictions | Description |
|---------------|---------------------------------------------------------------------|----------|--------------|-------------|
| `diagnostics` | array of [optimus-ide-collabsdk.FriendlyDiagnostic](#optimus-ide-collabsdkfriendlydiagnostic) | false    |              |             |
| `id`          | integer                                                             | false    |              |             |
| `parameters`  | array of [optimus-ide-collabsdk.PreviewParameter](#optimus-ide-collabsdkpreviewparameter)     | false    |              |             |

## optimus-ide-collabsdk.DynamicTool

```json
{
  "description": "string",
  "input_schema": [
    0
  ],
  "name": "string"
}
```

### Properties

| Name           | Type             | Required | Restrictions | Description                                                                                                                                  |
|----------------|------------------|----------|--------------|----------------------------------------------------------------------------------------------------------------------------------------------|
| `description`  | string           | false    |              |                                                                                                                                              |
| `input_schema` | array of integer | false    |              | Input schema JSON key "input_schema" uses snake_case for SDK consistency, deviating from the camelCase "inputSchema" convention used by MCP. |
| `name`         | string           | false    |              |                                                                                                                                              |

## optimus-ide-collabsdk.EditChatMessageRequest

```json
{
  "content": [
    {
      "content": "string",
      "end_line": 0,
      "file_id": "8a0cfb4f-ddc9-436d-91bb-75133c583767",
      "file_name": "string",
      "start_line": 0,
      "text": "string",
      "type": "text"
    }
  ],
  "model_config_id": "f5fb4d91-62ca-4377-9ee6-5d43ba00d205",
  "reasoning_effort": "string"
}
```

### Properties

| Name               | Type                                                      | Required | Restrictions | Description                                                                                                                                                                  |
|--------------------|-----------------------------------------------------------|----------|--------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `content`          | array of [optimus-ide-collabsdk.ChatInputPart](#optimus-ide-collabsdkchatinputpart) | false    |              |                                                                                                                                                                              |
| `model_config_id`  | string                                                    | false    |              | Model config ID when set, overrides the model used for the replacement user message and the assistant turn that follows. When nil the original message's model is preserved. |
| `reasoning_effort` | string                                                    | false    |              |                                                                                                                                                                              |

## optimus-ide-collabsdk.EditChatMessageResponse

```json
{
  "message": {
    "chat_id": "efc9fe20-a1e5-4a8c-9c48-f1b30c1e4f86",
    "content": [
      {
        "args": [
          0
        ],
        "args_delta": "string",
        "completed_at": "2019-08-24T14:15:22Z",
        "content": "string",
        "context_file_agent_id": {
          "uuid": "string",
          "valid": true
        },
        "context_file_content": "string",
        "context_file_directory": "string",
        "context_file_os": "string",
        "context_file_path": "string",
        "context_file_skill_meta_file": "string",
        "context_file_truncated": true,
        "created_at": "2019-08-24T14:15:22Z",
        "data": [
          0
        ],
        "end_line": 0,
        "file_id": {
          "uuid": "string",
          "valid": true
        },
        "file_name": "string",
        "is_error": true,
        "is_media": true,
        "mcp_server_config_id": {
          "uuid": "string",
          "valid": true
        },
        "media_type": "string",
        "name": "string",
        "parsed_commands": [
          [
            "string"
          ]
        ],
        "provider_executed": true,
        "provider_metadata": [
          0
        ],
        "result": [
          0
        ],
        "result_delta": "string",
        "result_reset": true,
        "skill_description": "string",
        "skill_dir": "string",
        "skill_name": "string",
        "source_id": "string",
        "start_line": 0,
        "text": "string",
        "title": "string",
        "tool_call_id": "string",
        "tool_name": "string",
        "type": "text",
        "url": "string"
      }
    ],
    "created_at": "2019-08-24T14:15:22Z",
    "created_by": "ee824cad-d7a6-4f48-87dc-e8461a9201c4",
    "id": 0,
    "model_config_id": "f5fb4d91-62ca-4377-9ee6-5d43ba00d205",
    "role": "system",
    "usage": {
      "cache_creation_tokens": 0,
      "cache_read_tokens": 0,
      "context_limit": 0,
      "input_tokens": 0,
      "output_tokens": 0,
      "reasoning_tokens": 0,
      "total_tokens": 0
    }
  },
  "warnings": [
    "string"
  ]
}
```

### Properties

| Name       | Type                                         | Required | Restrictions | Description |
|------------|----------------------------------------------|----------|--------------|-------------|
| `message`  | [optimus-ide-collabsdk.ChatMessage](#optimus-ide-collabsdkchatmessage) | false    |              |             |
| `warnings` | array of string                              | false    |              |             |

## optimus-ide-collabsdk.Entitlement

```json
"entitled"
```

### Properties

#### Enumerated Values

| Value(s)                                   |
|--------------------------------------------|
| `entitled`, `grace_period`, `not_entitled` |

## optimus-ide-collabsdk.Entitlements

```json
{
  "errors": [
    "string"
  ],
  "features": {
    "property1": {
      "actual": 0,
      "enabled": true,
      "entitlement": "entitled",
      "limit": 0,
      "usage_period": {
        "end": "2019-08-24T14:15:22Z",
        "issued_at": "2019-08-24T14:15:22Z",
        "start": "2019-08-24T14:15:22Z"
      }
    },
    "property2": {
      "actual": 0,
      "enabled": true,
      "entitlement": "entitled",
      "limit": 0,
      "usage_period": {
        "end": "2019-08-24T14:15:22Z",
        "issued_at": "2019-08-24T14:15:22Z",
        "start": "2019-08-24T14:15:22Z"
      }
    }
  },
  "has_license": true,
  "refreshed_at": "2019-08-24T14:15:22Z",
  "require_telemetry": true,
  "trial": true,
  "warnings": [
    "string"
  ]
}
```

### Properties

| Name                | Type                                 | Required | Restrictions | Description |
|---------------------|--------------------------------------|----------|--------------|-------------|
| `errors`            | array of string                      | false    |              |             |
| `features`          | object                               | false    |              |             |
| » `[any property]`  | [optimus-ide-collabsdk.Feature](#optimus-ide-collabsdkfeature) | false    |              |             |
| `has_license`       | boolean                              | false    |              |             |
| `refreshed_at`      | string                               | false    |              |             |
| `require_telemetry` | boolean                              | false    |              |             |
| `trial`             | boolean                              | false    |              |             |
| `warnings`          | array of string                      | false    |              |             |

## optimus-ide-collabsdk.Experiment

```json
"example"
```

### Properties

#### Enumerated Values

| Value(s)                                                                                                                                                                                                                                                                                               |
|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `ai-gateway-cost-control`, `ai-gateway-seat-exclusion`, `auto-fill-parameters`, `chat-advisor`, `chat-virtual-desktop`, `example`, `mcp-server-http`, `minimum-implicit-member`, `nats_pubsub`, `notifications`, `oauth2`, `workspace-build-updates`, `workspace-capable-licensing`, `workspace-usage` |

## optimus-ide-collabsdk.ExternalAPIKeyScopes

```json
{
  "external": [
    "all"
  ]
}
```

### Properties

| Name       | Type                                                  | Required | Restrictions | Description |
|------------|-------------------------------------------------------|----------|--------------|-------------|
| `external` | array of [optimus-ide-collabsdk.APIKeyScope](#optimus-ide-collabsdkapikeyscope) | false    |              |             |

## optimus-ide-collabsdk.ExternalAgentCredentials

```json
{
  "agent_token": "string",
  "command": "string"
}
```

### Properties

| Name          | Type   | Required | Restrictions | Description |
|---------------|--------|----------|--------------|-------------|
| `agent_token` | string | false    |              |             |
| `command`     | string | false    |              |             |

## optimus-ide-collabsdk.ExternalAuth

```json
{
  "app_install_url": "string",
  "app_installable": true,
  "authenticated": true,
  "device": true,
  "display_name": "string",
  "installations": [
    {
      "account": {
        "avatar_url": "string",
        "id": 0,
        "login": "string",
        "name": "string",
        "profile_url": "string"
      },
      "configure_url": "string",
      "id": 0
    }
  ],
  "supports_revocation": true,
  "user": {
    "avatar_url": "string",
    "id": 0,
    "login": "string",
    "name": "string",
    "profile_url": "string"
  }
}
```

### Properties

| Name                  | Type                                                                                  | Required | Restrictions | Description                                                             |
|-----------------------|---------------------------------------------------------------------------------------|----------|--------------|-------------------------------------------------------------------------|
| `app_install_url`     | string                                                                                | false    |              | App install URL is the URL to install the app.                          |
| `app_installable`     | boolean                                                                               | false    |              | App installable is true if the request for app installs was successful. |
| `authenticated`       | boolean                                                                               | false    |              |                                                                         |
| `device`              | boolean                                                                               | false    |              |                                                                         |
| `display_name`        | string                                                                                | false    |              |                                                                         |
| `installations`       | array of [optimus-ide-collabsdk.ExternalAuthAppInstallation](#optimus-ide-collabsdkexternalauthappinstallation) | false    |              | Installations are the installations that the user has access to.        |
| `supports_revocation` | boolean                                                                               | false    |              |                                                                         |
| `user`                | [optimus-ide-collabsdk.ExternalAuthUser](#optimus-ide-collabsdkexternalauthuser)                                | false    |              | User is the user that authenticated with the provider.                  |

## optimus-ide-collabsdk.ExternalAuthAppInstallation

```json
{
  "account": {
    "avatar_url": "string",
    "id": 0,
    "login": "string",
    "name": "string",
    "profile_url": "string"
  },
  "configure_url": "string",
  "id": 0
}
```

### Properties

| Name            | Type                                                   | Required | Restrictions | Description |
|-----------------|--------------------------------------------------------|----------|--------------|-------------|
| `account`       | [optimus-ide-collabsdk.ExternalAuthUser](#optimus-ide-collabsdkexternalauthuser) | false    |              |             |
| `configure_url` | string                                                 | false    |              |             |
| `id`            | integer                                                | false    |              |             |

## optimus-ide-collabsdk.ExternalAuthConfig

```json
{
  "api_base_url": "string",
  "app_install_url": "string",
  "app_installations_url": "string",
  "auth_url": "string",
  "client_id": "string",
  "code_challenge_methods_supported": [
    "string"
  ],
  "device_code_url": "string",
  "device_flow": true,
  "display_icon": "string",
  "display_name": "string",
  "id": "string",
  "mcp_tool_allow_regex": "string",
  "mcp_tool_deny_regex": "string",
  "mcp_url": "string",
  "no_refresh": true,
  "regex": "string",
  "revoke_url": "string",
  "scopes": [
    "string"
  ],
  "token_url": "string",
  "type": "string",
  "validate_url": "string"
}
```

### Properties

| Name                               | Type            | Required | Restrictions | Description                                                                                                                                                 |
|------------------------------------|-----------------|----------|--------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `api_base_url`                     | string          | false    |              | Api base URL is the base URL for provider REST API calls (e.g., "https://api.github.com" for GitHub). Derived from defaults when not explicitly configured. |
| `app_install_url`                  | string          | false    |              |                                                                                                                                                             |
| `app_installations_url`            | string          | false    |              |                                                                                                                                                             |
| `auth_url`                         | string          | false    |              |                                                                                                                                                             |
| `client_id`                        | string          | false    |              |                                                                                                                                                             |
| `code_challenge_methods_supported` | array of string | false    |              | Code challenge methods supported lists the PKCE code challenge methods The only one supported by Optimus-IDE-Collab is "S256".                                           |
| `device_code_url`                  | string          | false    |              |                                                                                                                                                             |
| `device_flow`                      | boolean         | false    |              |                                                                                                                                                             |
| `display_icon`                     | string          | false    |              | Display icon is a URL to an icon to display in the UI.                                                                                                      |
| `display_name`                     | string          | false    |              | Display name is shown in the UI to identify the auth config.                                                                                                |
| `id`                               | string          | false    |              | ID is a unique identifier for the auth config. It defaults to `type` when not provided.                                                                     |
| `mcp_tool_allow_regex`             | string          | false    |              | Deprecated: Injected MCP in AI Bridge is deprecated and will be removed in a future release.                                                                |
| `mcp_tool_deny_regex`              | string          | false    |              | Deprecated: Injected MCP in AI Bridge is deprecated and will be removed in a future release.                                                                |
| `mcp_url`                          | string          | false    |              | Deprecated: Injected MCP in AI Bridge is deprecated and will be removed in a future release.                                                                |
| `no_refresh`                       | boolean         | false    |              |                                                                                                                                                             |
|`regex`|string|false||Regex allows API requesters to match an auth config by a string (e.g. optimus-ide-collab.com) instead of by it's type.
Git clone makes use of this by parsing the URL from: 'Username for "https://github.com":' And sending it to the Optimus-IDE-Collab server to match against the Regex.|
|`revoke_url`|string|false|||
|`scopes`|array of string|false|||
|`token_url`|string|false|||
|`type`|string|false||Type is the type of external auth config.|
|`validate_url`|string|false|||

## optimus-ide-collabsdk.ExternalAuthDevice

```json
{
  "device_code": "string",
  "expires_in": 0,
  "interval": 0,
  "user_code": "string",
  "verification_uri": "string"
}
```

### Properties

| Name               | Type    | Required | Restrictions | Description |
|--------------------|---------|----------|--------------|-------------|
| `device_code`      | string  | false    |              |             |
| `expires_in`       | integer | false    |              |             |
| `interval`         | integer | false    |              |             |
| `user_code`        | string  | false    |              |             |
| `verification_uri` | string  | false    |              |             |

## optimus-ide-collabsdk.ExternalAuthLink

```json
{
  "authenticated": true,
  "created_at": "2019-08-24T14:15:22Z",
  "expires": "2019-08-24T14:15:22Z",
  "has_refresh_token": true,
  "provider_id": "string",
  "updated_at": "2019-08-24T14:15:22Z",
  "validate_error": "string"
}
```

### Properties

| Name                | Type    | Required | Restrictions | Description |
|---------------------|---------|----------|--------------|-------------|
| `authenticated`     | boolean | false    |              |             |
| `created_at`        | string  | false    |              |             |
| `expires`           | string  | false    |              |             |
| `has_refresh_token` | boolean | false    |              |             |
| `provider_id`       | string  | false    |              |             |
| `updated_at`        | string  | false    |              |             |
| `validate_error`    | string  | false    |              |             |

## optimus-ide-collabsdk.ExternalAuthUser

```json
{
  "avatar_url": "string",
  "id": 0,
  "login": "string",
  "name": "string",
  "profile_url": "string"
}
```

### Properties

| Name          | Type    | Required | Restrictions | Description |
|---------------|---------|----------|--------------|-------------|
| `avatar_url`  | string  | false    |              |             |
| `id`          | integer | false    |              |             |
| `login`       | string  | false    |              |             |
| `name`        | string  | false    |              |             |
| `profile_url` | string  | false    |              |             |

## optimus-ide-collabsdk.Feature

```json
{
  "actual": 0,
  "enabled": true,
  "entitlement": "entitled",
  "limit": 0,
  "usage_period": {
    "end": "2019-08-24T14:15:22Z",
    "issued_at": "2019-08-24T14:15:22Z",
    "start": "2019-08-24T14:15:22Z"
  }
}
```

### Properties

| Name          | Type                                         | Required | Restrictions | Description |
|---------------|----------------------------------------------|----------|--------------|-------------|
| `actual`      | integer                                      | false    |              |             |
| `enabled`     | boolean                                      | false    |              |             |
| `entitlement` | [optimus-ide-collabsdk.Entitlement](#optimus-ide-collabsdkentitlement) | false    |              |             |
| `limit`       | integer                                      | false    |              |             |
|`usage_period`|[optimus-ide-collabsdk.UsagePeriod](#optimus-ide-collabsdkusageperiod)|false||Usage period denotes that the usage is a counter that accumulates over this period (and most likely resets with the issuance of the next license).
These dates are determined from the license that this entitlement comes from, see enterprise/optimus-ide-collabd/license/license.go.
Only certain features set these fields: - FeatureManagedAgentLimit|

## optimus-ide-collabsdk.FriendlyDiagnostic

```json
{
  "detail": "string",
  "extra": {
    "code": "string"
  },
  "severity": "error",
  "summary": "string"
}
```

### Properties

| Name       | Type                                                                   | Required | Restrictions | Description |
|------------|------------------------------------------------------------------------|----------|--------------|-------------|
| `detail`   | string                                                                 | false    |              |             |
| `extra`    | [optimus-ide-collabsdk.DiagnosticExtra](#optimus-ide-collabsdkdiagnosticextra)                   | false    |              |             |
| `severity` | [optimus-ide-collabsdk.DiagnosticSeverityString](#optimus-ide-collabsdkdiagnosticseveritystring) | false    |              |             |
| `summary`  | string                                                                 | false    |              |             |

## optimus-ide-collabsdk.GenerateAPIKeyResponse

```json
{
  "key": "string"
}
```

### Properties

| Name  | Type   | Required | Restrictions | Description |
|-------|--------|----------|--------------|-------------|
| `key` | string | false    |              |             |

## optimus-ide-collabsdk.GetInboxNotificationResponse

```json
{
  "notification": {
    "actions": [
      {
        "label": "string",
        "url": "string"
      }
    ],
    "content": "string",
    "created_at": "2019-08-24T14:15:22Z",
    "icon": "string",
    "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
    "read_at": "string",
    "targets": [
      "497f6eca-6276-4993-bfeb-53cbbbba6f08"
    ],
    "template_id": "c6d67e98-83ea-49f0-8812-e4abae2b68bc",
    "title": "string",
    "user_id": "a169451c-8525-4352-b8ca-070dd449a1a5"
  },
  "unread_count": 0
}
```

### Properties

| Name           | Type                                                     | Required | Restrictions | Description |
|----------------|----------------------------------------------------------|----------|--------------|-------------|
| `notification` | [optimus-ide-collabsdk.InboxNotification](#optimus-ide-collabsdkinboxnotification) | false    |              |             |
| `unread_count` | integer                                                  | false    |              |             |

## optimus-ide-collabsdk.GetUserStatusCountsResponse

```json
{
  "status_counts": {
    "property1": [
      {
        "count": 10,
        "date": "2019-08-24T14:15:22Z"
      }
    ],
    "property2": [
      {
        "count": 10,
        "date": "2019-08-24T14:15:22Z"
      }
    ]
  }
}
```

### Properties

| Name               | Type                                                                      | Required | Restrictions | Description |
|--------------------|---------------------------------------------------------------------------|----------|--------------|-------------|
| `status_counts`    | object                                                                    | false    |              |             |
| » `[any property]` | array of [optimus-ide-collabsdk.UserStatusChangeCount](#optimus-ide-collabsdkuserstatuschangecount) | false    |              |             |

## optimus-ide-collabsdk.GetUsersResponse

```json
{
  "count": 0,
  "users": [
    {
      "avatar_url": "http://example.com",
      "created_at": "2019-08-24T14:15:22Z",
      "email": "user@example.com",
      "has_ai_seat": true,
      "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
      "is_service_account": true,
      "last_seen_at": "2019-08-24T14:15:22Z",
      "login_type": "",
      "name": "string",
      "organization_ids": [
        "497f6eca-6276-4993-bfeb-53cbbbba6f08"
      ],
      "roles": [
        {
          "display_name": "string",
          "name": "string",
          "organization_id": "string"
        }
      ],
      "status": "active",
      "theme_preference": "string",
      "updated_at": "2019-08-24T14:15:22Z",
      "username": "string"
    }
  ]
}
```

### Properties

| Name    | Type                                    | Required | Restrictions | Description |
|---------|-----------------------------------------|----------|--------------|-------------|
| `count` | integer                                 | false    |              |             |
| `users` | array of [optimus-ide-collabsdk.User](#optimus-ide-collabsdkuser) | false    |              |             |

## optimus-ide-collabsdk.GitSSHKey

```json
{
  "created_at": "2019-08-24T14:15:22Z",
  "public_key": "string",
  "updated_at": "2019-08-24T14:15:22Z",
  "user_id": "a169451c-8525-4352-b8ca-070dd449a1a5"
}
```

### Properties

| Name         | Type   | Required | Restrictions | Description                                                                                                                                                                                       |
|--------------|--------|----------|--------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `created_at` | string | false    |              |                                                                                                                                                                                                   |
| `public_key` | string | false    |              | Public key is the SSH public key in OpenSSH format. Example: "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAID3OmYJvT7q1cF1azbybYy0OZ9yrXfA+M6Lr4vzX5zlp\n" Note: The key includes a trailing newline (\n). |
| `updated_at` | string | false    |              |                                                                                                                                                                                                   |
| `user_id`    | string | false    |              |                                                                                                                                                                                                   |

## optimus-ide-collabsdk.GithubAuthMethod

```json
{
  "default_provider_configured": true,
  "enabled": true
}
```

### Properties

| Name                          | Type    | Required | Restrictions | Description |
|-------------------------------|---------|----------|--------------|-------------|
| `default_provider_configured` | boolean | false    |              |             |
| `enabled`                     | boolean | false    |              |             |

## optimus-ide-collabsdk.Group

```json
{
  "avatar_url": "http://example.com",
  "display_name": "string",
  "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
  "members": [
    {
      "avatar_url": "http://example.com",
      "created_at": "2019-08-24T14:15:22Z",
      "email": "user@example.com",
      "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
      "is_service_account": true,
      "last_seen_at": "2019-08-24T14:15:22Z",
      "login_type": "",
      "name": "string",
      "status": "active",
      "theme_preference": "string",
      "updated_at": "2019-08-24T14:15:22Z",
      "username": "string"
    }
  ],
  "name": "string",
  "organization_display_name": "string",
  "organization_id": "7c60d51f-b44e-4682-87d6-449835ea4de6",
  "organization_name": "string",
  "quota_allowance": 0,
  "source": "user",
  "total_member_count": 0
}
```

### Properties

| Name                        | Type                                                  | Required | Restrictions | Description                                                                                                                                                           |
|-----------------------------|-------------------------------------------------------|----------|--------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `avatar_url`                | string                                                | false    |              |                                                                                                                                                                       |
| `display_name`              | string                                                | false    |              |                                                                                                                                                                       |
| `id`                        | string                                                | false    |              |                                                                                                                                                                       |
| `members`                   | array of [optimus-ide-collabsdk.ReducedUser](#optimus-ide-collabsdkreduceduser) | false    |              |                                                                                                                                                                       |
| `name`                      | string                                                | false    |              |                                                                                                                                                                       |
| `organization_display_name` | string                                                | false    |              |                                                                                                                                                                       |
| `organization_id`           | string                                                | false    |              |                                                                                                                                                                       |
| `organization_name`         | string                                                | false    |              |                                                                                                                                                                       |
| `quota_allowance`           | integer                                               | false    |              |                                                                                                                                                                       |
| `source`                    | [optimus-ide-collabsdk.GroupSource](#optimus-ide-collabsdkgroupsource)          | false    |              |                                                                                                                                                                       |
| `total_member_count`        | integer                                               | false    |              | How many members are in this group. Shows the total count, even if the user is not authorized to read group member details. May be greater than `len(Group.Members)`. |

## optimus-ide-collabsdk.GroupAIBudget

```json
{
  "created_at": "2019-08-24T14:15:22Z",
  "group_id": "306db4e0-7449-4501-b76f-075576fe2d8f",
  "spend_limit_micros": 0,
  "updated_at": "2019-08-24T14:15:22Z"
}
```

### Properties

| Name                 | Type    | Required | Restrictions | Description |
|----------------------|---------|----------|--------------|-------------|
| `created_at`         | string  | false    |              |             |
| `group_id`           | string  | false    |              |             |
| `spend_limit_micros` | integer | false    |              |             |
| `updated_at`         | string  | false    |              |             |

## optimus-ide-collabsdk.GroupAISpend

```json
{
  "current_spend_micros": 0,
  "group_id": "306db4e0-7449-4501-b76f-075576fe2d8f",
  "period_end": "2019-08-24T14:15:22Z",
  "period_start": "2019-08-24T14:15:22Z",
  "spend_limit_micros": 0
}
```

### Properties

| Name                   | Type    | Required | Restrictions | Description                                                                                                |
|------------------------|---------|----------|--------------|------------------------------------------------------------------------------------------------------------|
| `current_spend_micros` | integer | false    |              | Current spend micros is the group's spend over the current budget period.                                  |
| `group_id`             | string  | false    |              |                                                                                                            |
| `period_end`           | string  | false    |              | Period end is the exclusive upper bound of the current budget period.                                      |
| `period_start`         | string  | false    |              | Period start is the inclusive lower bound of the current budget period.                                    |
| `spend_limit_micros`   | integer | false    |              | Spend limit micros is the group's configured AI spend limit. Null when the group has no configured budget. |

## optimus-ide-collabsdk.GroupMemberAISpend

```json
{
  "effective_group_id": "85e2b926-ddfb-4c66-b68e-b66e5acec6c0",
  "group_budget": {
    "limit_source": "user_override",
    "spend_limit_micros": 0
  },
  "group_spend_micros": 0,
  "user_id": "a169451c-8525-4352-b8ca-070dd449a1a5"
}
```

### Properties

| Name                 | Type                                             | Required | Restrictions | Description                                                                                                                                                                                                                                           |
|----------------------|--------------------------------------------------|----------|--------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `effective_group_id` | string                                           | false    |              | Effective group ID is the user's effective budget group within the queried group's organization, falling back to the Everyone group when no budget applies. Null when the effective group belongs to a different organization than the queried group. |
| `group_budget`       | [optimus-ide-collabsdk.AIGroupBudget](#optimus-ide-collabsdkaigroupbudget) | false    |              | Group budget is the budget when the queried group is this user's effective budget source. Null when the user's budget resolves to another group or no budget applies to the user.                                                                     |
| `group_spend_micros` | integer                                          | false    |              | Group spend micros is the user's spend attributed to the queried group over the current budget period.                                                                                                                                                |
| `user_id`            | string                                           | false    |              |                                                                                                                                                                                                                                                       |

## optimus-ide-collabsdk.GroupMembersAISpend

```json
{
  "members": [
    {
      "effective_group_id": "85e2b926-ddfb-4c66-b68e-b66e5acec6c0",
      "group_budget": {
        "limit_source": "user_override",
        "spend_limit_micros": 0
      },
      "group_spend_micros": 0,
      "user_id": "a169451c-8525-4352-b8ca-070dd449a1a5"
    }
  ],
  "period_end": "2019-08-24T14:15:22Z",
  "period_start": "2019-08-24T14:15:22Z"
}
```

### Properties

| Name           | Type                                                                | Required | Restrictions | Description                                                             |
|----------------|---------------------------------------------------------------------|----------|--------------|-------------------------------------------------------------------------|
| `members`      | array of [optimus-ide-collabsdk.GroupMemberAISpend](#optimus-ide-collabsdkgroupmemberaispend) | false    |              |                                                                         |
| `period_end`   | string                                                              | false    |              | Period end is the exclusive upper bound of the current budget period.   |
| `period_start` | string                                                              | false    |              | Period start is the inclusive lower bound of the current budget period. |

## optimus-ide-collabsdk.GroupMembersResponse

```json
{
  "count": 0,
  "users": [
    {
      "avatar_url": "http://example.com",
      "created_at": "2019-08-24T14:15:22Z",
      "email": "user@example.com",
      "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
      "is_service_account": true,
      "last_seen_at": "2019-08-24T14:15:22Z",
      "login_type": "",
      "name": "string",
      "status": "active",
      "theme_preference": "string",
      "updated_at": "2019-08-24T14:15:22Z",
      "username": "string"
    }
  ]
}
```

### Properties

| Name    | Type                                                  | Required | Restrictions | Description |
|---------|-------------------------------------------------------|----------|--------------|-------------|
| `count` | integer                                               | false    |              |             |
| `users` | array of [optimus-ide-collabsdk.ReducedUser](#optimus-ide-collabsdkreduceduser) | false    |              |             |

## optimus-ide-collabsdk.GroupSource

```json
"user"
```

### Properties

#### Enumerated Values

| Value(s)       |
|----------------|
| `oidc`, `user` |

## optimus-ide-collabsdk.GroupSyncSettings

```json
{
  "auto_create_missing_groups": true,
  "field": "string",
  "legacy_group_name_mapping": {
    "property1": "string",
    "property2": "string"
  },
  "mapping": {
    "property1": [
      "string"
    ],
    "property2": [
      "string"
    ]
  },
  "regex_filter": {}
}
```

### Properties

| Name                         | Type                           | Required | Restrictions | Description                                                                                                                                                                                                                                                                            |
|------------------------------|--------------------------------|----------|--------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `auto_create_missing_groups` | boolean                        | false    |              | Auto create missing groups controls whether groups returned by the OIDC provider are automatically created in Optimus-IDE-Collab if they are missing.                                                                                                                                               |
| `field`                      | string                         | false    |              | Field is the name of the claim field that specifies what groups a user should be in. If empty, no groups will be synced.                                                                                                                                                               |
| `legacy_group_name_mapping`  | object                         | false    |              | Legacy group name mapping is deprecated. It remaps an IDP group name to a Optimus-IDE-Collab group name. Since configuration is now done at runtime, group IDs are used to account for group renames. For legacy configurations, this config option has to remain. Deprecated: Use Mapping instead. |
| » `[any property]`           | string                         | false    |              |                                                                                                                                                                                                                                                                                        |
| `mapping`                    | object                         | false    |              | Mapping is a map from OIDC groups to Optimus-IDE-Collab group IDs                                                                                                                                                                                                                                   |
| » `[any property]`           | array of string                | false    |              |                                                                                                                                                                                                                                                                                        |
| `regex_filter`               | [regexp.Regexp](#regexpregexp) | false    |              | Regex filter is a regular expression that filters the groups returned by the OIDC provider. Any group not matched by this regex will be ignored. If the group filter is nil, then no group filtering will occur.                                                                       |

## optimus-ide-collabsdk.HTTPCookieConfig

```json
{
  "host_prefix": true,
  "same_site": "string",
  "secure_auth_cookie": true
}
```

### Properties

| Name                 | Type    | Required | Restrictions | Description |
|----------------------|---------|----------|--------------|-------------|
| `host_prefix`        | boolean | false    |              |             |
| `same_site`          | string  | false    |              |             |
| `secure_auth_cookie` | boolean | false    |              |             |

## optimus-ide-collabsdk.Healthcheck

```json
{
  "interval": 0,
  "threshold": 0,
  "url": "string"
}
```

### Properties

| Name        | Type    | Required | Restrictions | Description                                                                                      |
|-------------|---------|----------|--------------|--------------------------------------------------------------------------------------------------|
| `interval`  | integer | false    |              | Interval specifies the seconds between each health check.                                        |
| `threshold` | integer | false    |              | Threshold specifies the number of consecutive failed health checks before returning "unhealthy". |
| `url`       | string  | false    |              | URL specifies the endpoint to check for the app health.                                          |

## optimus-ide-collabsdk.HealthcheckConfig

```json
{
  "refresh": 0,
  "threshold_database": 0
}
```

### Properties

| Name                 | Type    | Required | Restrictions | Description |
|----------------------|---------|----------|--------------|-------------|
| `refresh`            | integer | false    |              |             |
| `threshold_database` | integer | false    |              |             |

## optimus-ide-collabsdk.ImportUserSecretsRequest

```json
{
  "content": "string",
  "format": "env"
}
```

### Properties

| Name      | Type                                                     | Required | Restrictions | Description |
|-----------|----------------------------------------------------------|----------|--------------|-------------|
| `content` | string                                                   | true     |              |             |
| `format`  | [optimus-ide-collabsdk.SecretsFileFormat](#optimus-ide-collabsdksecretsfileformat) | true     |              |             |

## optimus-ide-collabsdk.InboxNotification

```json
{
  "actions": [
    {
      "label": "string",
      "url": "string"
    }
  ],
  "content": "string",
  "created_at": "2019-08-24T14:15:22Z",
  "icon": "string",
  "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
  "read_at": "string",
  "targets": [
    "497f6eca-6276-4993-bfeb-53cbbbba6f08"
  ],
  "template_id": "c6d67e98-83ea-49f0-8812-e4abae2b68bc",
  "title": "string",
  "user_id": "a169451c-8525-4352-b8ca-070dd449a1a5"
}
```

### Properties

| Name          | Type                                                                          | Required | Restrictions | Description |
|---------------|-------------------------------------------------------------------------------|----------|--------------|-------------|
| `actions`     | array of [optimus-ide-collabsdk.InboxNotificationAction](#optimus-ide-collabsdkinboxnotificationaction) | false    |              |             |
| `content`     | string                                                                        | false    |              |             |
| `created_at`  | string                                                                        | false    |              |             |
| `icon`        | string                                                                        | false    |              |             |
| `id`          | string                                                                        | false    |              |             |
| `read_at`     | string                                                                        | false    |              |             |
| `targets`     | array of string                                                               | false    |              |             |
| `template_id` | string                                                                        | false    |              |             |
| `title`       | string                                                                        | false    |              |             |
| `user_id`     | string                                                                        | false    |              |             |

## optimus-ide-collabsdk.InboxNotificationAction

```json
{
  "label": "string",
  "url": "string"
}
```

### Properties

| Name    | Type   | Required | Restrictions | Description |
|---------|--------|----------|--------------|-------------|
| `label` | string | false    |              |             |
| `url`   | string | false    |              |             |

## optimus-ide-collabsdk.InsightsReportInterval

```json
"day"
```

### Properties

#### Enumerated Values

| Value(s)      |
|---------------|
| `day`, `week` |

## optimus-ide-collabsdk.InvalidatePresetsResponse

```json
{
  "invalidated": [
    {
      "preset_name": "string",
      "template_name": "string",
      "template_version_name": "string"
    }
  ]
}
```

### Properties

| Name          | Type                                                              | Required | Restrictions | Description |
|---------------|-------------------------------------------------------------------|----------|--------------|-------------|
| `invalidated` | array of [optimus-ide-collabsdk.InvalidatedPreset](#optimus-ide-collabsdkinvalidatedpreset) | false    |              |             |

## optimus-ide-collabsdk.InvalidatedPreset

```json
{
  "preset_name": "string",
  "template_name": "string",
  "template_version_name": "string"
}
```

### Properties

| Name                    | Type   | Required | Restrictions | Description |
|-------------------------|--------|----------|--------------|-------------|
| `preset_name`           | string | false    |              |             |
| `template_name`         | string | false    |              |             |
| `template_version_name` | string | false    |              |             |

## optimus-ide-collabsdk.IssueReconnectingPTYSignedTokenRequest

```json
{
  "agentID": "bc282582-04f9-45ce-b904-3e3bfab66958",
  "url": "string"
}
```

### Properties

| Name      | Type   | Required | Restrictions | Description                                                            |
|-----------|--------|----------|--------------|------------------------------------------------------------------------|
| `agentID` | string | true     |              |                                                                        |
| `url`     | string | true     |              | URL is the URL of the reconnecting-pty endpoint you are connecting to. |

## optimus-ide-collabsdk.IssueReconnectingPTYSignedTokenResponse

```json
{
  "signed_token": "string"
}
```

### Properties

| Name           | Type   | Required | Restrictions | Description |
|----------------|--------|----------|--------------|-------------|
| `signed_token` | string | false    |              |             |

## optimus-ide-collabsdk.JobErrorCode

```json
"REQUIRED_TEMPLATE_VARIABLES"
```

### Properties

#### Enumerated Values

| Value(s)                                            |
|-----------------------------------------------------|
| `INSUFFICIENT_QUOTA`, `REQUIRED_TEMPLATE_VARIABLES` |

## optimus-ide-collabsdk.License

```json
{
  "claims": {},
  "id": 0,
  "uploaded_at": "2019-08-24T14:15:22Z",
  "uuid": "095be615-a8ad-4c33-8e9c-c7612fbf6c9f"
}
```

### Properties

| Name          | Type    | Required | Restrictions | Description                                                                                                                                                                                             |
|---------------|---------|----------|--------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `claims`      | object  | false    |              | Claims are the JWT claims asserted by the license.  Here we use a generic string map to ensure that all data from the server is parsed verbatim, not just the fields this version of Optimus-IDE-Collab understands. |
| `id`          | integer | false    |              |                                                                                                                                                                                                         |
| `uploaded_at` | string  | false    |              |                                                                                                                                                                                                         |
| `uuid`        | string  | false    |              |                                                                                                                                                                                                         |

## optimus-ide-collabsdk.LinkConfig

```json
{
  "icon": "bug",
  "location": "navbar",
  "name": "string",
  "target": "string"
}
```

### Properties

| Name       | Type   | Required | Restrictions | Description |
|------------|--------|----------|--------------|-------------|
| `icon`     | string | false    |              |             |
| `location` | string | false    |              |             |
| `name`     | string | false    |              |             |
| `target`   | string | false    |              |             |

#### Enumerated Values

| Property   | Value(s)                      |
|------------|-------------------------------|
| `icon`     | `bug`, `chat`, `docs`, `star` |
| `location` | `dropdown`, `navbar`          |

## optimus-ide-collabsdk.ListInboxNotificationsResponse

```json
{
  "notifications": [
    {
      "actions": [
        {
          "label": "string",
          "url": "string"
        }
      ],
      "content": "string",
      "created_at": "2019-08-24T14:15:22Z",
      "icon": "string",
      "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
      "read_at": "string",
      "targets": [
        "497f6eca-6276-4993-bfeb-53cbbbba6f08"
      ],
      "template_id": "c6d67e98-83ea-49f0-8812-e4abae2b68bc",
      "title": "string",
      "user_id": "a169451c-8525-4352-b8ca-070dd449a1a5"
    }
  ],
  "unread_count": 0
}
```

### Properties

| Name            | Type                                                              | Required | Restrictions | Description |
|-----------------|-------------------------------------------------------------------|----------|--------------|-------------|
| `notifications` | array of [optimus-ide-collabsdk.InboxNotification](#optimus-ide-collabsdkinboxnotification) | false    |              |             |
| `unread_count`  | integer                                                           | false    |              |             |

## optimus-ide-collabsdk.LogLevel

```json
"trace"
```

### Properties

#### Enumerated Values

| Value(s)                                  |
|-------------------------------------------|
| `debug`, `error`, `info`, `trace`, `warn` |

## optimus-ide-collabsdk.LogSource

```json
"provisioner_daemon"
```

### Properties

#### Enumerated Values

| Value(s)                            |
|-------------------------------------|
| `provisioner`, `provisioner_daemon` |

## optimus-ide-collabsdk.LoggingConfig

```json
{
  "human": "string",
  "json": "string",
  "log_filter": [
    "string"
  ],
  "stackdriver": "string"
}
```

### Properties

| Name          | Type            | Required | Restrictions | Description |
|---------------|-----------------|----------|--------------|-------------|
| `human`       | string          | false    |              |             |
| `json`        | string          | false    |              |             |
| `log_filter`  | array of string | false    |              |             |
| `stackdriver` | string          | false    |              |             |

## optimus-ide-collabsdk.LoginType

```json
""
```

### Properties

#### Enumerated Values

| Value(s)                                          |
|---------------------------------------------------|
| ``, `github`, `none`, `oidc`, `password`, `token` |

## optimus-ide-collabsdk.LoginWithPasswordRequest

```json
{
  "email": "user@example.com",
  "password": "string"
}
```

### Properties

| Name       | Type   | Required | Restrictions | Description |
|------------|--------|----------|--------------|-------------|
| `email`    | string | true     |              |             |
| `password` | string | true     |              |             |

## optimus-ide-collabsdk.LoginWithPasswordResponse

```json
{
  "session_token": "string"
}
```

### Properties

| Name            | Type   | Required | Restrictions | Description |
|-----------------|--------|----------|--------------|-------------|
| `session_token` | string | true     |              |             |

## optimus-ide-collabsdk.MatchedProvisioners

```json
{
  "available": 0,
  "count": 0,
  "most_recently_seen": "2019-08-24T14:15:22Z"
}
```

### Properties

| Name                 | Type    | Required | Restrictions | Description                                                                                                                                                         |
|----------------------|---------|----------|--------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `available`          | integer | false    |              | Available is the number of provisioner daemons that are available to take jobs. This may be less than the count if some provisioners are busy or have been stopped. |
| `count`              | integer | false    |              | Count is the number of provisioner daemons that matched the given tags. If the count is 0, it means no provisioner daemons matched the requested tags.              |
| `most_recently_seen` | string  | false    |              | Most recently seen is the most recently seen time of the set of matched provisioners. If no provisioners matched, this field will be null.                          |

## optimus-ide-collabsdk.MinimalOrganization

```json
{
  "display_name": "string",
  "icon": "string",
  "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
  "name": "string"
}
```

### Properties

| Name           | Type   | Required | Restrictions | Description |
|----------------|--------|----------|--------------|-------------|
| `display_name` | string | false    |              |             |
| `icon`         | string | false    |              |             |
| `id`           | string | true     |              |             |
| `name`         | string | false    |              |             |

## optimus-ide-collabsdk.MinimalUser

```json
{
  "avatar_url": "http://example.com",
  "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
  "name": "string",
  "username": "string"
}
```

### Properties

| Name         | Type   | Required | Restrictions | Description |
|--------------|--------|----------|--------------|-------------|
| `avatar_url` | string | false    |              |             |
| `id`         | string | true     |              |             |
| `name`       | string | false    |              |             |
| `username`   | string | true     |              |             |

## optimus-ide-collabsdk.NotificationMethodsResponse

```json
{
  "available": [
    "string"
  ],
  "default": "string"
}
```

### Properties

| Name        | Type            | Required | Restrictions | Description |
|-------------|-----------------|----------|--------------|-------------|
| `available` | array of string | false    |              |             |
| `default`   | string          | false    |              |             |

## optimus-ide-collabsdk.NotificationPreference

```json
{
  "disabled": true,
  "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
  "updated_at": "2019-08-24T14:15:22Z"
}
```

### Properties

| Name         | Type    | Required | Restrictions | Description |
|--------------|---------|----------|--------------|-------------|
| `disabled`   | boolean | false    |              |             |
| `id`         | string  | false    |              |             |
| `updated_at` | string  | false    |              |             |

## optimus-ide-collabsdk.NotificationTemplate

```json
{
  "actions": "string",
  "body_template": "string",
  "enabled_by_default": true,
  "group": "string",
  "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
  "kind": "string",
  "method": "string",
  "name": "string",
  "title_template": "string"
}
```

### Properties

| Name                 | Type    | Required | Restrictions | Description |
|----------------------|---------|----------|--------------|-------------|
| `actions`            | string  | false    |              |             |
| `body_template`      | string  | false    |              |             |
| `enabled_by_default` | boolean | false    |              |             |
| `group`              | string  | false    |              |             |
| `id`                 | string  | false    |              |             |
| `kind`               | string  | false    |              |             |
| `method`             | string  | false    |              |             |
| `name`               | string  | false    |              |             |
| `title_template`     | string  | false    |              |             |

## optimus-ide-collabsdk.NotificationsConfig

```json
{
  "dispatch_timeout": 0,
  "email": {
    "auth": {
      "identity": "string",
      "password": "string",
      "password_file": "string",
      "username": "string"
    },
    "force_tls": true,
    "from": "string",
    "hello": "string",
    "smarthost": "string",
    "tls": {
      "ca_file": "string",
      "cert_file": "string",
      "insecure_skip_verify": true,
      "key_file": "string",
      "server_name": "string",
      "start_tls": true
    }
  },
  "fetch_interval": 0,
  "inbox": {
    "enabled": true
  },
  "lease_count": 0,
  "lease_period": 0,
  "max_send_attempts": 0,
  "method": "string",
  "retry_interval": 0,
  "sync_buffer_size": 0,
  "sync_interval": 0,
  "webhook": {
    "endpoint": {
      "forceQuery": true,
      "fragment": "string",
      "host": "string",
      "omitHost": true,
      "opaque": "string",
      "path": "string",
      "rawFragment": "string",
      "rawPath": "string",
      "rawQuery": "string",
      "scheme": "string",
      "user": {}
    }
  }
}
```

### Properties

| Name                | Type                                                                       | Required | Restrictions | Description                                                                                                                                                                                                                                                                                                                                                                                                                                         |
|---------------------|----------------------------------------------------------------------------|----------|--------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `dispatch_timeout`  | integer                                                                    | false    |              | How long to wait while a notification is being sent before giving up.                                                                                                                                                                                                                                                                                                                                                                               |
| `email`             | [optimus-ide-collabsdk.NotificationsEmailConfig](#optimus-ide-collabsdknotificationsemailconfig)     | false    |              | Email settings.                                                                                                                                                                                                                                                                                                                                                                                                                                     |
| `fetch_interval`    | integer                                                                    | false    |              | How often to query the database for queued notifications.                                                                                                                                                                                                                                                                                                                                                                                           |
| `inbox`             | [optimus-ide-collabsdk.NotificationsInboxConfig](#optimus-ide-collabsdknotificationsinboxconfig)     | false    |              | Inbox settings.                                                                                                                                                                                                                                                                                                                                                                                                                                     |
| `lease_count`       | integer                                                                    | false    |              | How many notifications a notifier should lease per fetch interval.                                                                                                                                                                                                                                                                                                                                                                                  |
| `lease_period`      | integer                                                                    | false    |              | How long a notifier should lease a message. This is effectively how long a notification is 'owned' by a notifier, and once this period expires it will be available for lease by another notifier. Leasing is important in order for multiple running notifiers to not pick the same messages to deliver concurrently. This lease period will only expire if a notifier shuts down ungracefully; a dispatch of the notification releases the lease. |
| `max_send_attempts` | integer                                                                    | false    |              | The upper limit of attempts to send a notification.                                                                                                                                                                                                                                                                                                                                                                                                 |
| `method`            | string                                                                     | false    |              | Which delivery method to use (available options: 'smtp', 'webhook').                                                                                                                                                                                                                                                                                                                                                                                |
| `retry_interval`    | integer                                                                    | false    |              | The minimum time between retries.                                                                                                                                                                                                                                                                                                                                                                                                                   |
| `sync_buffer_size`  | integer                                                                    | false    |              | The notifications system buffers message updates in memory to ease pressure on the database. This option controls how many updates are kept in memory. The lower this value the lower the change of state inconsistency in a non-graceful shutdown - but it also increases load on the database. It is recommended to keep this option at its default value.                                                                                        |
| `sync_interval`     | integer                                                                    | false    |              | The notifications system buffers message updates in memory to ease pressure on the database. This option controls how often it synchronizes its state with the database. The shorter this value the lower the change of state inconsistency in a non-graceful shutdown - but it also increases load on the database. It is recommended to keep this option at its default value.                                                                    |
| `webhook`           | [optimus-ide-collabsdk.NotificationsWebhookConfig](#optimus-ide-collabsdknotificationswebhookconfig) | false    |              | Webhook settings.                                                                                                                                                                                                                                                                                                                                                                                                                                   |

## optimus-ide-collabsdk.NotificationsEmailAuthConfig

```json
{
  "identity": "string",
  "password": "string",
  "password_file": "string",
  "username": "string"
}
```

### Properties

| Name            | Type   | Required | Restrictions | Description                                                |
|-----------------|--------|----------|--------------|------------------------------------------------------------|
| `identity`      | string | false    |              | Identity for PLAIN auth.                                   |
| `password`      | string | false    |              | Password for LOGIN/PLAIN auth.                             |
| `password_file` | string | false    |              | File from which to load the password for LOGIN/PLAIN auth. |
| `username`      | string | false    |              | Username for LOGIN/PLAIN auth.                             |

## optimus-ide-collabsdk.NotificationsEmailConfig

```json
{
  "auth": {
    "identity": "string",
    "password": "string",
    "password_file": "string",
    "username": "string"
  },
  "force_tls": true,
  "from": "string",
  "hello": "string",
  "smarthost": "string",
  "tls": {
    "ca_file": "string",
    "cert_file": "string",
    "insecure_skip_verify": true,
    "key_file": "string",
    "server_name": "string",
    "start_tls": true
  }
}
```

### Properties

| Name        | Type                                                                           | Required | Restrictions | Description                                                           |
|-------------|--------------------------------------------------------------------------------|----------|--------------|-----------------------------------------------------------------------|
| `auth`      | [optimus-ide-collabsdk.NotificationsEmailAuthConfig](#optimus-ide-collabsdknotificationsemailauthconfig) | false    |              | Authentication details.                                               |
| `force_tls` | boolean                                                                        | false    |              | Force tls causes a TLS connection to be attempted.                    |
| `from`      | string                                                                         | false    |              | The sender's address.                                                 |
| `hello`     | string                                                                         | false    |              | The hostname identifying the SMTP server.                             |
| `smarthost` | string                                                                         | false    |              | The intermediary SMTP host through which emails are sent (host:port). |
| `tls`       | [optimus-ide-collabsdk.NotificationsEmailTLSConfig](#optimus-ide-collabsdknotificationsemailtlsconfig)   | false    |              | Tls details.                                                          |

## optimus-ide-collabsdk.NotificationsEmailTLSConfig

```json
{
  "ca_file": "string",
  "cert_file": "string",
  "insecure_skip_verify": true,
  "key_file": "string",
  "server_name": "string",
  "start_tls": true
}
```

### Properties

| Name                   | Type    | Required | Restrictions | Description                                                  |
|------------------------|---------|----------|--------------|--------------------------------------------------------------|
| `ca_file`              | string  | false    |              | Ca file specifies the location of the CA certificate to use. |
| `cert_file`            | string  | false    |              | Cert file specifies the location of the certificate to use.  |
| `insecure_skip_verify` | boolean | false    |              | Insecure skip verify skips target certificate validation.    |
| `key_file`             | string  | false    |              | Key file specifies the location of the key to use.           |
| `server_name`          | string  | false    |              | Server name to verify the hostname for the targets.          |
| `start_tls`            | boolean | false    |              | Start tls attempts to upgrade plain connections to TLS.      |

## optimus-ide-collabsdk.NotificationsInboxConfig

```json
{
  "enabled": true
}
```

### Properties

| Name      | Type    | Required | Restrictions | Description |
|-----------|---------|----------|--------------|-------------|
| `enabled` | boolean | false    |              |             |

## optimus-ide-collabsdk.NotificationsSettings

```json
{
  "notifier_paused": true
}
```

### Properties

| Name              | Type    | Required | Restrictions | Description |
|-------------------|---------|----------|--------------|-------------|
| `notifier_paused` | boolean | false    |              |             |

## optimus-ide-collabsdk.NotificationsWebhookConfig

```json
{
  "endpoint": {
    "forceQuery": true,
    "fragment": "string",
    "host": "string",
    "omitHost": true,
    "opaque": "string",
    "path": "string",
    "rawFragment": "string",
    "rawPath": "string",
    "rawQuery": "string",
    "scheme": "string",
    "user": {}
  }
}
```

### Properties

| Name       | Type                       | Required | Restrictions | Description                                                          |
|------------|----------------------------|----------|--------------|----------------------------------------------------------------------|
| `endpoint` | [serpent.URL](#serpenturl) | false    |              | The URL to which the payload will be sent with an HTTP POST request. |

## optimus-ide-collabsdk.NullHCLString

```json
{
  "valid": true,
  "value": "string"
}
```

### Properties

| Name    | Type    | Required | Restrictions | Description |
|---------|---------|----------|--------------|-------------|
| `valid` | boolean | false    |              |             |
| `value` | string  | false    |              |             |

## optimus-ide-collabsdk.OAuth2AppEndpoints

```json
{
  "authorization": "string",
  "device_authorization": "string",
  "token": "string",
  "token_revoke": "string"
}
```

### Properties

| Name                   | Type   | Required | Restrictions | Description                       |
|------------------------|--------|----------|--------------|-----------------------------------|
| `authorization`        | string | false    |              |                                   |
| `device_authorization` | string | false    |              | Device authorization is optional. |
| `token`                | string | false    |              |                                   |
| `token_revoke`         | string | false    |              |                                   |

## optimus-ide-collabsdk.OAuth2AuthorizationServerMetadata

```json
{
  "authorization_endpoint": "string",
  "code_challenge_methods_supported": [
    "S256"
  ],
  "grant_types_supported": [
    "authorization_code"
  ],
  "issuer": "string",
  "registration_endpoint": "string",
  "response_types_supported": [
    "code"
  ],
  "revocation_endpoint": "string",
  "scopes_supported": [
    "string"
  ],
  "token_endpoint": "string",
  "token_endpoint_auth_methods_supported": [
    "client_secret_basic"
  ]
}
```

### Properties

| Name                                    | Type                                                                                      | Required | Restrictions | Description |
|-----------------------------------------|-------------------------------------------------------------------------------------------|----------|--------------|-------------|
| `authorization_endpoint`                | string                                                                                    | false    |              |             |
| `code_challenge_methods_supported`      | array of [optimus-ide-collabsdk.OAuth2PKCECodeChallengeMethod](#optimus-ide-collabsdkoauth2pkcecodechallengemethod) | false    |              |             |
| `grant_types_supported`                 | array of [optimus-ide-collabsdk.OAuth2ProviderGrantType](#optimus-ide-collabsdkoauth2providergranttype)             | false    |              |             |
| `issuer`                                | string                                                                                    | false    |              |             |
| `registration_endpoint`                 | string                                                                                    | false    |              |             |
| `response_types_supported`              | array of [optimus-ide-collabsdk.OAuth2ProviderResponseType](#optimus-ide-collabsdkoauth2providerresponsetype)       | false    |              |             |
| `revocation_endpoint`                   | string                                                                                    | false    |              |             |
| `scopes_supported`                      | array of string                                                                           | false    |              |             |
| `token_endpoint`                        | string                                                                                    | false    |              |             |
| `token_endpoint_auth_methods_supported` | array of [optimus-ide-collabsdk.OAuth2TokenEndpointAuthMethod](#optimus-ide-collabsdkoauth2tokenendpointauthmethod) | false    |              |             |

## optimus-ide-collabsdk.OAuth2ClientConfiguration

```json
{
  "client_id": "string",
  "client_id_issued_at": 0,
  "client_name": "string",
  "client_secret_expires_at": 0,
  "client_uri": "string",
  "contacts": [
    "string"
  ],
  "grant_types": [
    "authorization_code"
  ],
  "jwks": {},
  "jwks_uri": "string",
  "logo_uri": "string",
  "policy_uri": "string",
  "redirect_uris": [
    "string"
  ],
  "registration_access_token": "string",
  "registration_client_uri": "string",
  "response_types": [
    "code"
  ],
  "scope": "string",
  "software_id": "string",
  "software_version": "string",
  "token_endpoint_auth_method": "client_secret_basic",
  "tos_uri": "string"
}
```

### Properties

| Name                         | Type                                                                                | Required | Restrictions | Description |
|------------------------------|-------------------------------------------------------------------------------------|----------|--------------|-------------|
| `client_id`                  | string                                                                              | false    |              |             |
| `client_id_issued_at`        | integer                                                                             | false    |              |             |
| `client_name`                | string                                                                              | false    |              |             |
| `client_secret_expires_at`   | integer                                                                             | false    |              |             |
| `client_uri`                 | string                                                                              | false    |              |             |
| `contacts`                   | array of string                                                                     | false    |              |             |
| `grant_types`                | array of [optimus-ide-collabsdk.OAuth2ProviderGrantType](#optimus-ide-collabsdkoauth2providergranttype)       | false    |              |             |
| `jwks`                       | object                                                                              | false    |              |             |
| `jwks_uri`                   | string                                                                              | false    |              |             |
| `logo_uri`                   | string                                                                              | false    |              |             |
| `policy_uri`                 | string                                                                              | false    |              |             |
| `redirect_uris`              | array of string                                                                     | false    |              |             |
| `registration_access_token`  | string                                                                              | false    |              |             |
| `registration_client_uri`    | string                                                                              | false    |              |             |
| `response_types`             | array of [optimus-ide-collabsdk.OAuth2ProviderResponseType](#optimus-ide-collabsdkoauth2providerresponsetype) | false    |              |             |
| `scope`                      | string                                                                              | false    |              |             |
| `software_id`                | string                                                                              | false    |              |             |
| `software_version`           | string                                                                              | false    |              |             |
| `token_endpoint_auth_method` | [optimus-ide-collabsdk.OAuth2TokenEndpointAuthMethod](#optimus-ide-collabsdkoauth2tokenendpointauthmethod)    | false    |              |             |
| `tos_uri`                    | string                                                                              | false    |              |             |

## optimus-ide-collabsdk.OAuth2ClientRegistrationRequest

```json
{
  "client_name": "string",
  "client_uri": "string",
  "contacts": [
    "string"
  ],
  "grant_types": [
    "authorization_code"
  ],
  "jwks": {},
  "jwks_uri": "string",
  "logo_uri": "string",
  "policy_uri": "string",
  "redirect_uris": [
    "string"
  ],
  "response_types": [
    "code"
  ],
  "scope": "string",
  "software_id": "string",
  "software_statement": "string",
  "software_version": "string",
  "token_endpoint_auth_method": "client_secret_basic",
  "tos_uri": "string"
}
```

### Properties

| Name                         | Type                                                                                | Required | Restrictions | Description |
|------------------------------|-------------------------------------------------------------------------------------|----------|--------------|-------------|
| `client_name`                | string                                                                              | false    |              |             |
| `client_uri`                 | string                                                                              | false    |              |             |
| `contacts`                   | array of string                                                                     | false    |              |             |
| `grant_types`                | array of [optimus-ide-collabsdk.OAuth2ProviderGrantType](#optimus-ide-collabsdkoauth2providergranttype)       | false    |              |             |
| `jwks`                       | object                                                                              | false    |              |             |
| `jwks_uri`                   | string                                                                              | false    |              |             |
| `logo_uri`                   | string                                                                              | false    |              |             |
| `policy_uri`                 | string                                                                              | false    |              |             |
| `redirect_uris`              | array of string                                                                     | false    |              |             |
| `response_types`             | array of [optimus-ide-collabsdk.OAuth2ProviderResponseType](#optimus-ide-collabsdkoauth2providerresponsetype) | false    |              |             |
| `scope`                      | string                                                                              | false    |              |             |
| `software_id`                | string                                                                              | false    |              |             |
| `software_statement`         | string                                                                              | false    |              |             |
| `software_version`           | string                                                                              | false    |              |             |
| `token_endpoint_auth_method` | [optimus-ide-collabsdk.OAuth2TokenEndpointAuthMethod](#optimus-ide-collabsdkoauth2tokenendpointauthmethod)    | false    |              |             |
| `tos_uri`                    | string                                                                              | false    |              |             |

## optimus-ide-collabsdk.OAuth2ClientRegistrationResponse

```json
{
  "client_id": "string",
  "client_id_issued_at": 0,
  "client_name": "string",
  "client_secret": "string",
  "client_secret_expires_at": 0,
  "client_uri": "string",
  "contacts": [
    "string"
  ],
  "grant_types": [
    "authorization_code"
  ],
  "jwks": {},
  "jwks_uri": "string",
  "logo_uri": "string",
  "policy_uri": "string",
  "redirect_uris": [
    "string"
  ],
  "registration_access_token": "string",
  "registration_client_uri": "string",
  "response_types": [
    "code"
  ],
  "scope": "string",
  "software_id": "string",
  "software_version": "string",
  "token_endpoint_auth_method": "client_secret_basic",
  "tos_uri": "string"
}
```

### Properties

| Name                         | Type                                                                                | Required | Restrictions | Description |
|------------------------------|-------------------------------------------------------------------------------------|----------|--------------|-------------|
| `client_id`                  | string                                                                              | false    |              |             |
| `client_id_issued_at`        | integer                                                                             | false    |              |             |
| `client_name`                | string                                                                              | false    |              |             |
| `client_secret`              | string                                                                              | false    |              |             |
| `client_secret_expires_at`   | integer                                                                             | false    |              |             |
| `client_uri`                 | string                                                                              | false    |              |             |
| `contacts`                   | array of string                                                                     | false    |              |             |
| `grant_types`                | array of [optimus-ide-collabsdk.OAuth2ProviderGrantType](#optimus-ide-collabsdkoauth2providergranttype)       | false    |              |             |
| `jwks`                       | object                                                                              | false    |              |             |
| `jwks_uri`                   | string                                                                              | false    |              |             |
| `logo_uri`                   | string                                                                              | false    |              |             |
| `policy_uri`                 | string                                                                              | false    |              |             |
| `redirect_uris`              | array of string                                                                     | false    |              |             |
| `registration_access_token`  | string                                                                              | false    |              |             |
| `registration_client_uri`    | string                                                                              | false    |              |             |
| `response_types`             | array of [optimus-ide-collabsdk.OAuth2ProviderResponseType](#optimus-ide-collabsdkoauth2providerresponsetype) | false    |              |             |
| `scope`                      | string                                                                              | false    |              |             |
| `software_id`                | string                                                                              | false    |              |             |
| `software_version`           | string                                                                              | false    |              |             |
| `token_endpoint_auth_method` | [optimus-ide-collabsdk.OAuth2TokenEndpointAuthMethod](#optimus-ide-collabsdkoauth2tokenendpointauthmethod)    | false    |              |             |
| `tos_uri`                    | string                                                                              | false    |              |             |

## optimus-ide-collabsdk.OAuth2Config

```json
{
  "github": {
    "allow_everyone": true,
    "allow_signups": true,
    "allowed_orgs": [
      "string"
    ],
    "allowed_teams": [
      "string"
    ],
    "client_id": "string",
    "client_secret": "string",
    "default_provider_enable": true,
    "device_flow": true,
    "enterprise_base_url": "string"
  }
}
```

### Properties

| Name     | Type                                                       | Required | Restrictions | Description |
|----------|------------------------------------------------------------|----------|--------------|-------------|
| `github` | [optimus-ide-collabsdk.OAuth2GithubConfig](#optimus-ide-collabsdkoauth2githubconfig) | false    |              |             |

## optimus-ide-collabsdk.OAuth2GithubConfig

```json
{
  "allow_everyone": true,
  "allow_signups": true,
  "allowed_orgs": [
    "string"
  ],
  "allowed_teams": [
    "string"
  ],
  "client_id": "string",
  "client_secret": "string",
  "default_provider_enable": true,
  "device_flow": true,
  "enterprise_base_url": "string"
}
```

### Properties

| Name                      | Type            | Required | Restrictions | Description |
|---------------------------|-----------------|----------|--------------|-------------|
| `allow_everyone`          | boolean         | false    |              |             |
| `allow_signups`           | boolean         | false    |              |             |
| `allowed_orgs`            | array of string | false    |              |             |
| `allowed_teams`           | array of string | false    |              |             |
| `client_id`               | string          | false    |              |             |
| `client_secret`           | string          | false    |              |             |
| `default_provider_enable` | boolean         | false    |              |             |
| `device_flow`             | boolean         | false    |              |             |
| `enterprise_base_url`     | string          | false    |              |             |

## optimus-ide-collabsdk.OAuth2PKCECodeChallengeMethod

```json
"S256"
```

### Properties

#### Enumerated Values

| Value(s)        |
|-----------------|
| `S256`, `plain` |

## optimus-ide-collabsdk.OAuth2ProtectedResourceMetadata

```json
{
  "authorization_servers": [
    "string"
  ],
  "bearer_methods_supported": [
    "string"
  ],
  "resource": "string",
  "scopes_supported": [
    "string"
  ]
}
```

### Properties

| Name                       | Type            | Required | Restrictions | Description |
|----------------------------|-----------------|----------|--------------|-------------|
| `authorization_servers`    | array of string | false    |              |             |
| `bearer_methods_supported` | array of string | false    |              |             |
| `resource`                 | string          | false    |              |             |
| `scopes_supported`         | array of string | false    |              |             |

## optimus-ide-collabsdk.OAuth2ProviderApp

```json
{
  "callback_url": "string",
  "endpoints": {
    "authorization": "string",
    "device_authorization": "string",
    "token": "string",
    "token_revoke": "string"
  },
  "icon": "string",
  "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
  "name": "string"
}
```

### Properties

| Name           | Type                                                       | Required | Restrictions | Description                                                                                                                                                                                             |
|----------------|------------------------------------------------------------|----------|--------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `callback_url` | string                                                     | false    |              |                                                                                                                                                                                                         |
| `endpoints`    | [optimus-ide-collabsdk.OAuth2AppEndpoints](#optimus-ide-collabsdkoauth2appendpoints) | false    |              | Endpoints are included in the app response for easier discovery. The OAuth2 spec does not have a defined place to find these (for comparison, OIDC has a '/.well-known/openid-configuration' endpoint). |
| `icon`         | string                                                     | false    |              |                                                                                                                                                                                                         |
| `id`           | string                                                     | false    |              |                                                                                                                                                                                                         |
| `name`         | string                                                     | false    |              |                                                                                                                                                                                                         |

## optimus-ide-collabsdk.OAuth2ProviderAppSecret

```json
{
  "client_secret_truncated": "string",
  "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
  "last_used_at": "string"
}
```

### Properties

| Name                      | Type   | Required | Restrictions | Description |
|---------------------------|--------|----------|--------------|-------------|
| `client_secret_truncated` | string | false    |              |             |
| `id`                      | string | false    |              |             |
| `last_used_at`            | string | false    |              |             |

## optimus-ide-collabsdk.OAuth2ProviderAppSecretFull

```json
{
  "client_secret_full": "string",
  "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08"
}
```

### Properties

| Name                 | Type   | Required | Restrictions | Description |
|----------------------|--------|----------|--------------|-------------|
| `client_secret_full` | string | false    |              |             |
| `id`                 | string | false    |              |             |

## optimus-ide-collabsdk.OAuth2ProviderGrantType

```json
"authorization_code"
```

### Properties

#### Enumerated Values

| Value(s)                                                                            |
|-------------------------------------------------------------------------------------|
| `authorization_code`, `client_credentials`, `implicit`, `password`, `refresh_token` |

## optimus-ide-collabsdk.OAuth2ProviderResponseType

```json
"code"
```

### Properties

#### Enumerated Values

| Value(s)        |
|-----------------|
| `code`, `token` |

## optimus-ide-collabsdk.OAuth2ProviderSettings

```json
{
  "dynamic_client_registration_enabled": true
}
```

### Properties

| Name                                  | Type    | Required | Restrictions | Description |
|---------------------------------------|---------|----------|--------------|-------------|
| `dynamic_client_registration_enabled` | boolean | false    |              |             |

## optimus-ide-collabsdk.OAuth2TokenEndpointAuthMethod

```json
"client_secret_basic"
```

### Properties

#### Enumerated Values

| Value(s)                                            |
|-----------------------------------------------------|
| `client_secret_basic`, `client_secret_post`, `none` |

## optimus-ide-collabsdk.OAuthConversionResponse

```json
{
  "expires_at": "2019-08-24T14:15:22Z",
  "state_string": "string",
  "to_type": "",
  "user_id": "a169451c-8525-4352-b8ca-070dd449a1a5"
}
```

### Properties

| Name           | Type                                     | Required | Restrictions | Description |
|----------------|------------------------------------------|----------|--------------|-------------|
| `expires_at`   | string                                   | false    |              |             |
| `state_string` | string                                   | false    |              |             |
| `to_type`      | [optimus-ide-collabsdk.LoginType](#optimus-ide-collabsdklogintype) | false    |              |             |
| `user_id`      | string                                   | false    |              |             |

## optimus-ide-collabsdk.OIDCAuthMethod

```json
{
  "enabled": true,
  "iconUrl": "string",
  "signInText": "string"
}
```

### Properties

| Name         | Type    | Required | Restrictions | Description |
|--------------|---------|----------|--------------|-------------|
| `enabled`    | boolean | false    |              |             |
| `iconUrl`    | string  | false    |              |             |
| `signInText` | string  | false    |              |             |

## optimus-ide-collabsdk.OIDCClaimsResponse

```json
{
  "claims": {}
}
```

### Properties

| Name     | Type   | Required | Restrictions | Description                                                                                                                                                                 |
|----------|--------|----------|--------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `claims` | object | false    |              | Claims are the merged claims from the OIDC provider. These are the union of the ID token claims and the userinfo claims, where userinfo claims take precedence on conflict. |

## optimus-ide-collabsdk.OIDCConfig

```json
{
  "allow_signups": true,
  "auth_url_params": {},
  "auto_repair_links": true,
  "client_cert_file": "string",
  "client_id": "string",
  "client_key_file": "string",
  "client_secret": "string",
  "email_domain": [
    "string"
  ],
  "email_fallback": true,
  "email_field": "string",
  "group_allow_list": [
    "string"
  ],
  "group_auto_create": true,
  "group_mapping": {},
  "group_regex_filter": {},
  "groups_field": "string",
  "icon_url": {
    "forceQuery": true,
    "fragment": "string",
    "host": "string",
    "omitHost": true,
    "opaque": "string",
    "path": "string",
    "rawFragment": "string",
    "rawPath": "string",
    "rawQuery": "string",
    "scheme": "string",
    "user": {}
  },
  "ignore_email_verified": true,
  "ignore_user_info": true,
  "issuer_url": "string",
  "name_field": "string",
  "organization_assign_default": true,
  "organization_field": "string",
  "organization_mapping": {},
  "redirect_allowed_hosts": [
    "string"
  ],
  "redirect_url": {
    "forceQuery": true,
    "fragment": "string",
    "host": "string",
    "omitHost": true,
    "opaque": "string",
    "path": "string",
    "rawFragment": "string",
    "rawPath": "string",
    "rawQuery": "string",
    "scheme": "string",
    "user": {}
  },
  "scopes": [
    "string"
  ],
  "sign_in_text": "string",
  "signups_disabled_text": "string",
  "skip_issuer_checks": true,
  "source_user_info_from_access_token": true,
  "user_role_field": "string",
  "user_role_mapping": {},
  "user_roles_default": [
    "string"
  ],
  "username_field": "string"
}
```

### Properties

| Name                                 | Type                             | Required | Restrictions | Description                                                                                                                                                                                                                                                                                                                                                                                                                                                |
|--------------------------------------|----------------------------------|----------|--------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `allow_signups`                      | boolean                          | false    |              |                                                                                                                                                                                                                                                                                                                                                                                                                                                            |
| `auth_url_params`                    | object                           | false    |              |                                                                                                                                                                                                                                                                                                                                                                                                                                                            |
| `auto_repair_links`                  | boolean                          | false    |              |                                                                                                                                                                                                                                                                                                                                                                                                                                                            |
| `client_cert_file`                   | string                           | false    |              |                                                                                                                                                                                                                                                                                                                                                                                                                                                            |
| `client_id`                          | string                           | false    |              |                                                                                                                                                                                                                                                                                                                                                                                                                                                            |
| `client_key_file`                    | string                           | false    |              | Client key file & ClientCertFile are used in place of ClientSecret for PKI auth.                                                                                                                                                                                                                                                                                                                                                                           |
| `client_secret`                      | string                           | false    |              |                                                                                                                                                                                                                                                                                                                                                                                                                                                            |
| `email_domain`                       | array of string                  | false    |              |                                                                                                                                                                                                                                                                                                                                                                                                                                                            |
| `email_fallback`                     | boolean                          | false    |              | Email fallback allows OIDC logins to fall back to email-based matching when the `linked_id` (issuer+subject) does not match an existing user link. INSECURE: weakens the linked_id check. It exists for IdP brokers that do not issue a stable `sub` for the same user across connections.                                                                                                                                                                 |
| `email_field`                        | string                           | false    |              |                                                                                                                                                                                                                                                                                                                                                                                                                                                            |
| `group_allow_list`                   | array of string                  | false    |              |                                                                                                                                                                                                                                                                                                                                                                                                                                                            |
| `group_auto_create`                  | boolean                          | false    |              |                                                                                                                                                                                                                                                                                                                                                                                                                                                            |
| `group_mapping`                      | object                           | false    |              |                                                                                                                                                                                                                                                                                                                                                                                                                                                            |
| `group_regex_filter`                 | [serpent.Regexp](#serpentregexp) | false    |              |                                                                                                                                                                                                                                                                                                                                                                                                                                                            |
| `groups_field`                       | string                           | false    |              |                                                                                                                                                                                                                                                                                                                                                                                                                                                            |
| `icon_url`                           | [serpent.URL](#serpenturl)       | false    |              |                                                                                                                                                                                                                                                                                                                                                                                                                                                            |
| `ignore_email_verified`              | boolean                          | false    |              |                                                                                                                                                                                                                                                                                                                                                                                                                                                            |
| `ignore_user_info`                   | boolean                          | false    |              | Ignore user info & UserInfoFromAccessToken are mutually exclusive. Only 1 can be set to true. Ideally this would be an enum with 3 states, ['none', 'userinfo', 'access_token']. However, for backward compatibility, `ignore_user_info` must remain. And `access_token` is a niche, non-spec compliant edge case. So it's use is rare, and should not be advised.                                                                                         |
| `issuer_url`                         | string                           | false    |              |                                                                                                                                                                                                                                                                                                                                                                                                                                                            |
| `name_field`                         | string                           | false    |              |                                                                                                                                                                                                                                                                                                                                                                                                                                                            |
| `organization_assign_default`        | boolean                          | false    |              |                                                                                                                                                                                                                                                                                                                                                                                                                                                            |
| `organization_field`                 | string                           | false    |              |                                                                                                                                                                                                                                                                                                                                                                                                                                                            |
| `organization_mapping`               | object                           | false    |              |                                                                                                                                                                                                                                                                                                                                                                                                                                                            |
| `redirect_allowed_hosts`             | array of string                  | false    |              | Redirect allowed hosts is an allowlist of hostnames that may be used as the host of the OIDC redirect_uri. When non-empty, the redirect_uri is constructed from the incoming request's Host header (validated against this list) instead of from AccessURL. Every listed host must also be registered as a valid redirect URI in the OIDC provider. This setting is mutually exclusive with RedirectURL: if RedirectURL is set, this allowlist is ignored. |
| `redirect_url`                       | [serpent.URL](#serpenturl)       | false    |              | Redirect URL is optional, defaulting to 'ACCESS_URL'. Only useful in niche situations where the OIDC callback domain is different from the ACCESS_URL domain.                                                                                                                                                                                                                                                                                              |
| `scopes`                             | array of string                  | false    |              |                                                                                                                                                                                                                                                                                                                                                                                                                                                            |
| `sign_in_text`                       | string                           | false    |              |                                                                                                                                                                                                                                                                                                                                                                                                                                                            |
| `signups_disabled_text`              | string                           | false    |              |                                                                                                                                                                                                                                                                                                                                                                                                                                                            |
| `skip_issuer_checks`                 | boolean                          | false    |              |                                                                                                                                                                                                                                                                                                                                                                                                                                                            |
| `source_user_info_from_access_token` | boolean                          | false    |              | Source user info from access token as mentioned above is an edge case. This allows sourcing the user_info from the access token itself instead of a user_info endpoint. This assumes the access token is a valid JWT with a set of claims to be merged with the id_token.                                                                                                                                                                                  |
| `user_role_field`                    | string                           | false    |              |                                                                                                                                                                                                                                                                                                                                                                                                                                                            |
| `user_role_mapping`                  | object                           | false    |              |                                                                                                                                                                                                                                                                                                                                                                                                                                                            |
| `user_roles_default`                 | array of string                  | false    |              |                                                                                                                                                                                                                                                                                                                                                                                                                                                            |
| `username_field`                     | string                           | false    |              |                                                                                                                                                                                                                                                                                                                                                                                                                                                            |

## optimus-ide-collabsdk.OptionType

```json
"string"
```

### Properties

#### Enumerated Values

| Value(s)                                   |
|--------------------------------------------|
| `bool`, `list(string)`, `number`, `string` |

## optimus-ide-collabsdk.Organization

```json
{
  "created_at": "2019-08-24T14:15:22Z",
  "default_org_member_roles": [
    "string"
  ],
  "description": "string",
  "display_name": "string",
  "icon": "string",
  "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
  "is_default": true,
  "name": "string",
  "updated_at": "2019-08-24T14:15:22Z"
}
```

### Properties

| Name                       | Type            | Required | Restrictions | Description                                                                                                                                     |
|----------------------------|-----------------|----------|--------------|-------------------------------------------------------------------------------------------------------------------------------------------------|
| `created_at`               | string          | true     |              |                                                                                                                                                 |
| `default_org_member_roles` | array of string | false    |              | Default org member roles are unioned into every member's effective roles at request time. Changes propagate to all members on the next request. |
| `description`              | string          | false    |              |                                                                                                                                                 |
| `display_name`             | string          | false    |              |                                                                                                                                                 |
| `icon`                     | string          | false    |              |                                                                                                                                                 |
| `id`                       | string          | true     |              |                                                                                                                                                 |
| `is_default`               | boolean         | true     |              |                                                                                                                                                 |
| `name`                     | string          | false    |              |                                                                                                                                                 |
| `updated_at`               | string          | true     |              |                                                                                                                                                 |

## optimus-ide-collabsdk.OrganizationGroupAISpend

```json
{
  "current_spend_micros": 0,
  "group_id": "306db4e0-7449-4501-b76f-075576fe2d8f",
  "spend_limit_micros": 0
}
```

### Properties

| Name                   | Type    | Required | Restrictions | Description                                                                                                |
|------------------------|---------|----------|--------------|------------------------------------------------------------------------------------------------------------|
| `current_spend_micros` | integer | false    |              | Current spend micros is the group's spend over the current budget period.                                  |
| `group_id`             | string  | false    |              |                                                                                                            |
| `spend_limit_micros`   | integer | false    |              | Spend limit micros is the group's configured AI spend limit. Null when the group has no configured budget. |

## optimus-ide-collabsdk.OrganizationGroupsAISpend

```json
{
  "groups": [
    {
      "current_spend_micros": 0,
      "group_id": "306db4e0-7449-4501-b76f-075576fe2d8f",
      "spend_limit_micros": 0
    }
  ],
  "period_end": "2019-08-24T14:15:22Z",
  "period_start": "2019-08-24T14:15:22Z"
}
```

### Properties

| Name           | Type                                                                            | Required | Restrictions | Description                                                             |
|----------------|---------------------------------------------------------------------------------|----------|--------------|-------------------------------------------------------------------------|
| `groups`       | array of [optimus-ide-collabsdk.OrganizationGroupAISpend](#optimus-ide-collabsdkorganizationgroupaispend) | false    |              |                                                                         |
| `period_end`   | string                                                                          | false    |              | Period end is the exclusive upper bound of the current budget period.   |
| `period_start` | string                                                                          | false    |              | Period start is the inclusive lower bound of the current budget period. |

## optimus-ide-collabsdk.OrganizationMember

```json
{
  "created_at": "2019-08-24T14:15:22Z",
  "organization_id": "7c60d51f-b44e-4682-87d6-449835ea4de6",
  "roles": [
    {
      "display_name": "string",
      "name": "string",
      "organization_id": "string"
    }
  ],
  "updated_at": "2019-08-24T14:15:22Z",
  "user_id": "a169451c-8525-4352-b8ca-070dd449a1a5"
}
```

### Properties

| Name              | Type                                            | Required | Restrictions | Description |
|-------------------|-------------------------------------------------|----------|--------------|-------------|
| `created_at`      | string                                          | false    |              |             |
| `organization_id` | string                                          | false    |              |             |
| `roles`           | array of [optimus-ide-collabsdk.SlimRole](#optimus-ide-collabsdkslimrole) | false    |              |             |
| `updated_at`      | string                                          | false    |              |             |
| `user_id`         | string                                          | false    |              |             |

## optimus-ide-collabsdk.OrganizationMemberWithUserData

```json
{
  "avatar_url": "string",
  "created_at": "2019-08-24T14:15:22Z",
  "email": "string",
  "global_roles": [
    {
      "display_name": "string",
      "name": "string",
      "organization_id": "string"
    }
  ],
  "has_ai_seat": true,
  "is_service_account": true,
  "last_seen_at": "2019-08-24T14:15:22Z",
  "login_type": "",
  "name": "string",
  "organization_id": "7c60d51f-b44e-4682-87d6-449835ea4de6",
  "roles": [
    {
      "display_name": "string",
      "name": "string",
      "organization_id": "string"
    }
  ],
  "status": "active",
  "updated_at": "2019-08-24T14:15:22Z",
  "user_created_at": "2019-08-24T14:15:22Z",
  "user_id": "a169451c-8525-4352-b8ca-070dd449a1a5",
  "user_updated_at": "2019-08-24T14:15:22Z",
  "username": "string"
}
```

### Properties

| Name                 | Type                                            | Required | Restrictions | Description                                                                                      |
|----------------------|-------------------------------------------------|----------|--------------|--------------------------------------------------------------------------------------------------|
| `avatar_url`         | string                                          | false    |              |                                                                                                  |
| `created_at`         | string                                          | false    |              |                                                                                                  |
| `email`              | string                                          | false    |              |                                                                                                  |
| `global_roles`       | array of [optimus-ide-collabsdk.SlimRole](#optimus-ide-collabsdkslimrole) | false    |              |                                                                                                  |
| `has_ai_seat`        | boolean                                         | false    |              | Has ai seat intentionally omits omitempty so the API always includes the field, even when false. |
| `is_service_account` | boolean                                         | false    |              |                                                                                                  |
| `last_seen_at`       | string                                          | false    |              |                                                                                                  |
| `login_type`         | [optimus-ide-collabsdk.LoginType](#optimus-ide-collabsdklogintype)        | false    |              |                                                                                                  |
| `name`               | string                                          | false    |              |                                                                                                  |
| `organization_id`    | string                                          | false    |              |                                                                                                  |
| `roles`              | array of [optimus-ide-collabsdk.SlimRole](#optimus-ide-collabsdkslimrole) | false    |              |                                                                                                  |
| `status`             | [optimus-ide-collabsdk.UserStatus](#optimus-ide-collabsdkuserstatus)      | false    |              |                                                                                                  |
| `updated_at`         | string                                          | false    |              |                                                                                                  |
| `user_created_at`    | string                                          | false    |              |                                                                                                  |
| `user_id`            | string                                          | false    |              |                                                                                                  |
| `user_updated_at`    | string                                          | false    |              |                                                                                                  |
| `username`           | string                                          | false    |              |                                                                                                  |

#### Enumerated Values

| Property | Value(s)              |
|----------|-----------------------|
| `status` | `active`, `suspended` |

## optimus-ide-collabsdk.OrganizationSyncSettings

```json
{
  "field": "string",
  "mapping": {
    "property1": [
      "string"
    ],
    "property2": [
      "string"
    ]
  },
  "organization_assign_default": true
}
```

### Properties

| Name                          | Type            | Required | Restrictions | Description                                                                                                                                                                         |
|-------------------------------|-----------------|----------|--------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `field`                       | string          | false    |              | Field selects the claim field to be used as the created user's organizations. If the field is the empty string, then no organization updates will ever come from the OIDC provider. |
| `mapping`                     | object          | false    |              | Mapping maps from an OIDC claim --> Optimus-IDE-Collab organization uuid                                                                                                                         |
| » `[any property]`            | array of string | false    |              |                                                                                                                                                                                     |
| `organization_assign_default` | boolean         | false    |              | Organization assign default will ensure the default org is always included for every user, regardless of their claims. This preserves legacy behavior.                              |

## optimus-ide-collabsdk.PaginatedMembersResponse

```json
{
  "count": 0,
  "members": [
    {
      "avatar_url": "string",
      "created_at": "2019-08-24T14:15:22Z",
      "email": "string",
      "global_roles": [
        {
          "display_name": "string",
          "name": "string",
          "organization_id": "string"
        }
      ],
      "has_ai_seat": true,
      "is_service_account": true,
      "last_seen_at": "2019-08-24T14:15:22Z",
      "login_type": "",
      "name": "string",
      "organization_id": "7c60d51f-b44e-4682-87d6-449835ea4de6",
      "roles": [
        {
          "display_name": "string",
          "name": "string",
          "organization_id": "string"
        }
      ],
      "status": "active",
      "updated_at": "2019-08-24T14:15:22Z",
      "user_created_at": "2019-08-24T14:15:22Z",
      "user_id": "a169451c-8525-4352-b8ca-070dd449a1a5",
      "user_updated_at": "2019-08-24T14:15:22Z",
      "username": "string"
    }
  ]
}
```

### Properties

| Name      | Type                                                                                        | Required | Restrictions | Description |
|-----------|---------------------------------------------------------------------------------------------|----------|--------------|-------------|
| `count`   | integer                                                                                     | false    |              |             |
| `members` | array of [optimus-ide-collabsdk.OrganizationMemberWithUserData](#optimus-ide-collabsdkorganizationmemberwithuserdata) | false    |              |             |

## optimus-ide-collabsdk.ParameterFormType

```json
""
```

### Properties

#### Enumerated Values

| Value(s)                                                                                                            |
|---------------------------------------------------------------------------------------------------------------------|
| ``, `checkbox`, `dropdown`, `error`, `input`, `multi-select`, `radio`, `slider`, `switch`, `tag-select`, `textarea` |

## optimus-ide-collabsdk.PatchGroupIDPSyncConfigRequest

```json
{
  "auto_create_missing_groups": true,
  "field": "string",
  "regex_filter": {}
}
```

### Properties

| Name                         | Type                           | Required | Restrictions | Description |
|------------------------------|--------------------------------|----------|--------------|-------------|
| `auto_create_missing_groups` | boolean                        | false    |              |             |
| `field`                      | string                         | false    |              |             |
| `regex_filter`               | [regexp.Regexp](#regexpregexp) | false    |              |             |

## optimus-ide-collabsdk.PatchGroupIDPSyncMappingRequest

```json
{
  "add": [
    {
      "gets": "string",
      "given": "string"
    }
  ],
  "remove": [
    {
      "gets": "string",
      "given": "string"
    }
  ]
}
```

### Properties

| Name      | Type            | Required | Restrictions | Description                                              |
|-----------|-----------------|----------|--------------|----------------------------------------------------------|
| `add`     | array of object | false    |              |                                                          |
| `» gets`  | string          | false    |              | The ID of the Optimus-IDE-Collab resource the user should be added to |
| `» given` | string          | false    |              | The IdP claim the user has                               |
| `remove`  | array of object | false    |              |                                                          |
| `» gets`  | string          | false    |              | The ID of the Optimus-IDE-Collab resource the user should be added to |
| `» given` | string          | false    |              | The IdP claim the user has                               |

## optimus-ide-collabsdk.PatchGroupRequest

```json
{
  "add_users": [
    "string"
  ],
  "avatar_url": "string",
  "display_name": "string",
  "name": "string",
  "quota_allowance": 0,
  "remove_users": [
    "string"
  ]
}
```

### Properties

| Name              | Type            | Required | Restrictions | Description |
|-------------------|-----------------|----------|--------------|-------------|
| `add_users`       | array of string | false    |              |             |
| `avatar_url`      | string          | false    |              |             |
| `display_name`    | string          | false    |              |             |
| `name`            | string          | false    |              |             |
| `quota_allowance` | integer         | false    |              |             |
| `remove_users`    | array of string | false    |              |             |

## optimus-ide-collabsdk.PatchOrganizationIDPSyncConfigRequest

```json
{
  "assign_default": true,
  "field": "string"
}
```

### Properties

| Name             | Type    | Required | Restrictions | Description |
|------------------|---------|----------|--------------|-------------|
| `assign_default` | boolean | false    |              |             |
| `field`          | string  | false    |              |             |

## optimus-ide-collabsdk.PatchOrganizationIDPSyncMappingRequest

```json
{
  "add": [
    {
      "gets": "string",
      "given": "string"
    }
  ],
  "remove": [
    {
      "gets": "string",
      "given": "string"
    }
  ]
}
```

### Properties

| Name      | Type            | Required | Restrictions | Description                                              |
|-----------|-----------------|----------|--------------|----------------------------------------------------------|
| `add`     | array of object | false    |              |                                                          |
| `» gets`  | string          | false    |              | The ID of the Optimus-IDE-Collab resource the user should be added to |
| `» given` | string          | false    |              | The IdP claim the user has                               |
| `remove`  | array of object | false    |              |                                                          |
| `» gets`  | string          | false    |              | The ID of the Optimus-IDE-Collab resource the user should be added to |
| `» given` | string          | false    |              | The IdP claim the user has                               |

## optimus-ide-collabsdk.PatchRoleIDPSyncConfigRequest

```json
{
  "field": "string"
}
```

### Properties

| Name    | Type   | Required | Restrictions | Description |
|---------|--------|----------|--------------|-------------|
| `field` | string | false    |              |             |

## optimus-ide-collabsdk.PatchRoleIDPSyncMappingRequest

```json
{
  "add": [
    {
      "gets": "string",
      "given": "string"
    }
  ],
  "remove": [
    {
      "gets": "string",
      "given": "string"
    }
  ]
}
```

### Properties

| Name      | Type            | Required | Restrictions | Description                                              |
|-----------|-----------------|----------|--------------|----------------------------------------------------------|
| `add`     | array of object | false    |              |                                                          |
| `» gets`  | string          | false    |              | The ID of the Optimus-IDE-Collab resource the user should be added to |
| `» given` | string          | false    |              | The IdP claim the user has                               |
| `remove`  | array of object | false    |              |                                                          |
| `» gets`  | string          | false    |              | The ID of the Optimus-IDE-Collab resource the user should be added to |
| `» given` | string          | false    |              | The IdP claim the user has                               |

## optimus-ide-collabsdk.PatchTemplateVersionRequest

```json
{
  "message": "string",
  "name": "string"
}
```

### Properties

| Name      | Type   | Required | Restrictions | Description |
|-----------|--------|----------|--------------|-------------|
| `message` | string | false    |              |             |
| `name`    | string | false    |              |             |

## optimus-ide-collabsdk.PatchWorkspaceProxy

```json
{
  "display_name": "string",
  "icon": "string",
  "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
  "name": "string",
  "regenerate_token": true
}
```

### Properties

| Name               | Type    | Required | Restrictions | Description |
|--------------------|---------|----------|--------------|-------------|
| `display_name`     | string  | true     |              |             |
| `icon`             | string  | true     |              |             |
| `id`               | string  | true     |              |             |
| `name`             | string  | true     |              |             |
| `regenerate_token` | boolean | false    |              |             |

## optimus-ide-collabsdk.PauseTaskResponse

```json
{
  "workspace_build": {
    "build_number": 0,
    "created_at": "2019-08-24T14:15:22Z",
    "daily_cost": 0,
    "deadline": "2019-08-24T14:15:22Z",
    "has_ai_task": true,
    "has_external_agent": true,
    "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
    "initiator_id": "06588898-9a84-4b35-ba8f-f9cbd64946f3",
    "initiator_name": "string",
    "job": {
      "available_workers": [
        "497f6eca-6276-4993-bfeb-53cbbbba6f08"
      ],
      "canceled_at": "2019-08-24T14:15:22Z",
      "completed_at": "2019-08-24T14:15:22Z",
      "created_at": "2019-08-24T14:15:22Z",
      "error": "string",
      "error_code": "REQUIRED_TEMPLATE_VARIABLES",
      "file_id": "8a0cfb4f-ddc9-436d-91bb-75133c583767",
      "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
      "initiator_id": "06588898-9a84-4b35-ba8f-f9cbd64946f3",
      "input": {
        "error": "string",
        "template_version_id": "0ba39c92-1f1b-4c32-aa3e-9925d7713eb1",
        "workspace_build_id": "badaf2eb-96c5-4050-9f1d-db2d39ca5478"
      },
      "logs_overflowed": true,
      "metadata": {
        "template_display_name": "string",
        "template_icon": "string",
        "template_id": "c6d67e98-83ea-49f0-8812-e4abae2b68bc",
        "template_name": "string",
        "template_version_name": "string",
        "workspace_build_transition": "start",
        "workspace_id": "0967198e-ec7b-4c6b-b4d3-f71244cadbe9",
        "workspace_name": "string"
      },
      "organization_id": "7c60d51f-b44e-4682-87d6-449835ea4de6",
      "queue_position": 0,
      "queue_size": 0,
      "started_at": "2019-08-24T14:15:22Z",
      "status": "pending",
      "tags": {
        "property1": "string",
        "property2": "string"
      },
      "type": "template_version_import",
      "worker_id": "ae5fa6f7-c55b-40c1-b40a-b36ac467652b",
      "worker_name": "string"
    },
    "matched_provisioners": {
      "available": 0,
      "count": 0,
      "most_recently_seen": "2019-08-24T14:15:22Z"
    },
    "max_deadline": "2019-08-24T14:15:22Z",
    "reason": "initiator",
    "resources": [
      {
        "agents": [
          {
            "api_version": "string",
            "apps": [
              {
                "command": "string",
                "display_name": "string",
                "external": true,
                "group": "string",
                "health": "disabled",
                "healthcheck": {
                  "interval": 0,
                  "threshold": 0,
                  "url": "string"
                },
                "hidden": true,
                "icon": "string",
                "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
                "open_in": "slim-window",
                "sharing_level": "owner",
                "slug": "string",
                "statuses": [
                  {
                    "agent_id": "2b1e3b65-2c04-4fa2-a2d7-467901e98978",
                    "app_id": "affd1d10-9538-4fc8-9e0b-4594a28c1335",
                    "created_at": "2019-08-24T14:15:22Z",
                    "icon": "string",
                    "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
                    "message": "string",
                    "needs_user_attention": true,
                    "state": "working",
                    "uri": "string",
                    "workspace_id": "0967198e-ec7b-4c6b-b4d3-f71244cadbe9"
                  }
                ],
                "subdomain": true,
                "subdomain_name": "string",
                "tooltip": "string",
                "url": "string"
              }
            ],
            "architecture": "string",
            "connection_timeout_seconds": 0,
            "created_at": "2019-08-24T14:15:22Z",
            "directory": "string",
            "disconnected_at": "2019-08-24T14:15:22Z",
            "display_apps": [
              "vscode"
            ],
            "environment_variables": {
              "property1": "string",
              "property2": "string"
            },
            "expanded_directory": "string",
            "first_connected_at": "2019-08-24T14:15:22Z",
            "health": {
              "healthy": false,
              "reason": "agent has lost connection"
            },
            "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
            "instance_id": "string",
            "last_connected_at": "2019-08-24T14:15:22Z",
            "latency": {
              "property1": {
                "latency_ms": 0,
                "preferred": true
              },
              "property2": {
                "latency_ms": 0,
                "preferred": true
              }
            },
            "lifecycle_state": "created",
            "log_sources": [
              {
                "created_at": "2019-08-24T14:15:22Z",
                "display_name": "string",
                "icon": "string",
                "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
                "workspace_agent_id": "7ad2e618-fea7-4c1a-b70a-f501566a72f1"
              }
            ],
            "logs_length": 0,
            "logs_overflowed": true,
            "name": "string",
            "operating_system": "string",
            "parent_id": {
              "uuid": "string",
              "valid": true
            },
            "ready_at": "2019-08-24T14:15:22Z",
            "resource_id": "4d5215ed-38bb-48ed-879a-fdb9ca58522f",
            "scripts": [
              {
                "cron": "string",
                "display_name": "string",
                "exit_code": 0,
                "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
                "log_path": "string",
                "log_source_id": "4197ab25-95cf-4b91-9c78-f7f2af5d353a",
                "run_on_start": true,
                "run_on_stop": true,
                "script": "string",
                "start_blocks_login": true,
                "status": "ok",
                "timeout": 0
              }
            ],
            "started_at": "2019-08-24T14:15:22Z",
            "startup_script_behavior": "blocking",
            "status": "connecting",
            "subsystems": [
              "envbox"
            ],
            "troubleshooting_url": "string",
            "updated_at": "2019-08-24T14:15:22Z",
            "version": "string"
          }
        ],
        "created_at": "2019-08-24T14:15:22Z",
        "daily_cost": 0,
        "hide": true,
        "icon": "string",
        "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
        "job_id": "453bd7d7-5355-4d6d-a38e-d9e7eb218c3f",
        "metadata": [
          {
            "key": "string",
            "sensitive": true,
            "value": "string"
          }
        ],
        "name": "string",
        "type": "string",
        "workspace_transition": "start"
      }
    ],
    "status": "pending",
    "template_version_id": "0ba39c92-1f1b-4c32-aa3e-9925d7713eb1",
    "template_version_name": "string",
    "template_version_preset_id": "512a53a7-30da-446e-a1fc-713c630baff1",
    "transition": "start",
    "updated_at": "2019-08-24T14:15:22Z",
    "workspace_id": "0967198e-ec7b-4c6b-b4d3-f71244cadbe9",
    "workspace_name": "string",
    "workspace_owner_avatar_url": "string",
    "workspace_owner_id": "e7078695-5279-4c86-8774-3ac2367a2fc7",
    "workspace_owner_name": "string"
  }
}
```

### Properties

| Name              | Type                                               | Required | Restrictions | Description |
|-------------------|----------------------------------------------------|----------|--------------|-------------|
| `workspace_build` | [optimus-ide-collabsdk.WorkspaceBuild](#optimus-ide-collabsdkworkspacebuild) | false    |              |             |

## optimus-ide-collabsdk.Permission

```json
{
  "action": "application_connect",
  "negate": true,
  "resource_type": "*"
}
```

### Properties

| Name            | Type                                           | Required | Restrictions | Description                             |
|-----------------|------------------------------------------------|----------|--------------|-----------------------------------------|
| `action`        | [optimus-ide-collabsdk.RBACAction](#optimus-ide-collabsdkrbacaction)     | false    |              |                                         |
| `negate`        | boolean                                        | false    |              | Negate makes this a negative permission |
| `resource_type` | [optimus-ide-collabsdk.RBACResource](#optimus-ide-collabsdkrbacresource) | false    |              |                                         |

## optimus-ide-collabsdk.PostOAuth2ProviderAppRequest

```json
{
  "callback_url": "string",
  "icon": "string",
  "name": "string"
}
```

### Properties

| Name           | Type   | Required | Restrictions | Description |
|----------------|--------|----------|--------------|-------------|
| `callback_url` | string | true     |              |             |
| `icon`         | string | false    |              |             |
| `name`         | string | true     |              |             |

## optimus-ide-collabsdk.PostWorkspaceUsageRequest

```json
{
  "agent_id": "2b1e3b65-2c04-4fa2-a2d7-467901e98978",
  "app_name": "vscode"
}
```

### Properties

| Name       | Type                                           | Required | Restrictions | Description |
|------------|------------------------------------------------|----------|--------------|-------------|
| `agent_id` | string                                         | false    |              |             |
| `app_name` | [optimus-ide-collabsdk.UsageAppName](#optimus-ide-collabsdkusageappname) | false    |              |             |

## optimus-ide-collabsdk.PprofConfig

```json
{
  "address": {
    "host": "string",
    "port": "string"
  },
  "enable": true
}
```

### Properties

| Name      | Type                                 | Required | Restrictions | Description |
|-----------|--------------------------------------|----------|--------------|-------------|
| `address` | [serpent.HostPort](#serpenthostport) | false    |              |             |
| `enable`  | boolean                              | false    |              |             |

## optimus-ide-collabsdk.PrebuildsConfig

```json
{
  "failure_hard_limit": 0,
  "reconciliation_backoff_interval": 0,
  "reconciliation_backoff_lookback": 0,
  "reconciliation_interval": 0
}
```

### Properties

| Name                              | Type    | Required | Restrictions | Description                                                                                                                                                                                                                                                                                       |
|-----------------------------------|---------|----------|--------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `failure_hard_limit`              | integer | false    |              | Failure hard limit defines the maximum number of consecutive failed prebuild attempts allowed before a preset is considered to be in a hard limit state. When a preset hits this limit, no new prebuilds will be created until the limit is reset. FailureHardLimit is disabled when set to zero. |
| `reconciliation_backoff_interval` | integer | false    |              | Reconciliation backoff interval specifies the amount of time to increase the backoff interval when errors occur during reconciliation.                                                                                                                                                            |
| `reconciliation_backoff_lookback` | integer | false    |              | Reconciliation backoff lookback determines the time window to look back when calculating the number of failed prebuilds, which influences the backoff strategy.                                                                                                                                   |
| `reconciliation_interval`         | integer | false    |              | Reconciliation interval defines how often the workspace prebuilds state should be reconciled.                                                                                                                                                                                                     |

## optimus-ide-collabsdk.PrebuildsSettings

```json
{
  "reconciliation_paused": true
}
```

### Properties

| Name                    | Type    | Required | Restrictions | Description |
|-------------------------|---------|----------|--------------|-------------|
| `reconciliation_paused` | boolean | false    |              |             |

## optimus-ide-collabsdk.Preset

```json
{
  "default": true,
  "description": "string",
  "desiredPrebuildInstances": 0,
  "icon": "string",
  "id": "string",
  "name": "string",
  "parameters": [
    {
      "name": "string",
      "value": "string"
    }
  ]
}
```

### Properties

| Name                       | Type                                                          | Required | Restrictions | Description |
|----------------------------|---------------------------------------------------------------|----------|--------------|-------------|
| `default`                  | boolean                                                       | false    |              |             |
| `description`              | string                                                        | false    |              |             |
| `desiredPrebuildInstances` | integer                                                       | false    |              |             |
| `icon`                     | string                                                        | false    |              |             |
| `id`                       | string                                                        | false    |              |             |
| `name`                     | string                                                        | false    |              |             |
| `parameters`               | array of [optimus-ide-collabsdk.PresetParameter](#optimus-ide-collabsdkpresetparameter) | false    |              |             |

## optimus-ide-collabsdk.PresetParameter

```json
{
  "name": "string",
  "value": "string"
}
```

### Properties

| Name    | Type   | Required | Restrictions | Description |
|---------|--------|----------|--------------|-------------|
| `name`  | string | false    |              |             |
| `value` | string | false    |              |             |

## optimus-ide-collabsdk.PreviewParameter

```json
{
  "default_value": {
    "valid": true,
    "value": "string"
  },
  "description": "string",
  "diagnostics": [
    {
      "detail": "string",
      "extra": {
        "code": "string"
      },
      "severity": "error",
      "summary": "string"
    }
  ],
  "display_name": "string",
  "ephemeral": true,
  "form_type": "",
  "icon": "string",
  "mutable": true,
  "name": "string",
  "options": [
    {
      "description": "string",
      "icon": "string",
      "name": "string",
      "value": {
        "valid": true,
        "value": "string"
      }
    }
  ],
  "order": 0,
  "required": true,
  "styling": {
    "disabled": true,
    "label": "string",
    "mask_input": true,
    "placeholder": "string"
  },
  "type": "string",
  "validations": [
    {
      "validation_error": "string",
      "validation_max": 0,
      "validation_min": 0,
      "validation_monotonic": "string",
      "validation_regex": "string"
    }
  ],
  "value": {
    "valid": true,
    "value": "string"
  }
}
```

### Properties

| Name            | Type                                                                                | Required | Restrictions | Description                             |
|-----------------|-------------------------------------------------------------------------------------|----------|--------------|-----------------------------------------|
| `default_value` | [optimus-ide-collabsdk.NullHCLString](#optimus-ide-collabsdknullhclstring)                                    | false    |              |                                         |
| `description`   | string                                                                              | false    |              |                                         |
| `diagnostics`   | array of [optimus-ide-collabsdk.FriendlyDiagnostic](#optimus-ide-collabsdkfriendlydiagnostic)                 | false    |              |                                         |
| `display_name`  | string                                                                              | false    |              |                                         |
| `ephemeral`     | boolean                                                                             | false    |              |                                         |
| `form_type`     | [optimus-ide-collabsdk.ParameterFormType](#optimus-ide-collabsdkparameterformtype)                            | false    |              |                                         |
| `icon`          | string                                                                              | false    |              |                                         |
| `mutable`       | boolean                                                                             | false    |              |                                         |
| `name`          | string                                                                              | false    |              |                                         |
| `options`       | array of [optimus-ide-collabsdk.PreviewParameterOption](#optimus-ide-collabsdkpreviewparameteroption)         | false    |              |                                         |
| `order`         | integer                                                                             | false    |              | legacy_variable_name was removed (= 14) |
| `required`      | boolean                                                                             | false    |              |                                         |
| `styling`       | [optimus-ide-collabsdk.PreviewParameterStyling](#optimus-ide-collabsdkpreviewparameterstyling)                | false    |              |                                         |
| `type`          | [optimus-ide-collabsdk.OptionType](#optimus-ide-collabsdkoptiontype)                                          | false    |              |                                         |
| `validations`   | array of [optimus-ide-collabsdk.PreviewParameterValidation](#optimus-ide-collabsdkpreviewparametervalidation) | false    |              |                                         |
| `value`         | [optimus-ide-collabsdk.NullHCLString](#optimus-ide-collabsdknullhclstring)                                    | false    |              |                                         |

## optimus-ide-collabsdk.PreviewParameterOption

```json
{
  "description": "string",
  "icon": "string",
  "name": "string",
  "value": {
    "valid": true,
    "value": "string"
  }
}
```

### Properties

| Name          | Type                                             | Required | Restrictions | Description |
|---------------|--------------------------------------------------|----------|--------------|-------------|
| `description` | string                                           | false    |              |             |
| `icon`        | string                                           | false    |              |             |
| `name`        | string                                           | false    |              |             |
| `value`       | [optimus-ide-collabsdk.NullHCLString](#optimus-ide-collabsdknullhclstring) | false    |              |             |

## optimus-ide-collabsdk.PreviewParameterStyling

```json
{
  "disabled": true,
  "label": "string",
  "mask_input": true,
  "placeholder": "string"
}
```

### Properties

| Name          | Type    | Required | Restrictions | Description |
|---------------|---------|----------|--------------|-------------|
| `disabled`    | boolean | false    |              |             |
| `label`       | string  | false    |              |             |
| `mask_input`  | boolean | false    |              |             |
| `placeholder` | string  | false    |              |             |

## optimus-ide-collabsdk.PreviewParameterValidation

```json
{
  "validation_error": "string",
  "validation_max": 0,
  "validation_min": 0,
  "validation_monotonic": "string",
  "validation_regex": "string"
}
```

### Properties

| Name                   | Type    | Required | Restrictions | Description                             |
|------------------------|---------|----------|--------------|-----------------------------------------|
| `validation_error`     | string  | false    |              |                                         |
| `validation_max`       | integer | false    |              |                                         |
| `validation_min`       | integer | false    |              |                                         |
| `validation_monotonic` | string  | false    |              |                                         |
| `validation_regex`     | string  | false    |              | All validation attributes are optional. |

## optimus-ide-collabsdk.PrometheusConfig

```json
{
  "address": {
    "host": "string",
    "port": "string"
  },
  "aggregate_agent_stats_by": [
    "string"
  ],
  "collect_agent_stats": true,
  "collect_db_metrics": true,
  "enable": true
}
```

### Properties

| Name                       | Type                                 | Required | Restrictions | Description |
|----------------------------|--------------------------------------|----------|--------------|-------------|
| `address`                  | [serpent.HostPort](#serpenthostport) | false    |              |             |
| `aggregate_agent_stats_by` | array of string                      | false    |              |             |
| `collect_agent_stats`      | boolean                              | false    |              |             |
| `collect_db_metrics`       | boolean                              | false    |              |             |
| `enable`                   | boolean                              | false    |              |             |

## optimus-ide-collabsdk.ProvisionerConfig

```json
{
  "daemon_poll_interval": 0,
  "daemon_poll_jitter": 0,
  "daemon_psk": "string",
  "daemon_types": [
    "string"
  ],
  "daemons": 0,
  "force_cancel_interval": 0
}
```

### Properties

| Name                    | Type            | Required | Restrictions | Description                                               |
|-------------------------|-----------------|----------|--------------|-----------------------------------------------------------|
| `daemon_poll_interval`  | integer         | false    |              |                                                           |
| `daemon_poll_jitter`    | integer         | false    |              |                                                           |
| `daemon_psk`            | string          | false    |              |                                                           |
| `daemon_types`          | array of string | false    |              |                                                           |
| `daemons`               | integer         | false    |              | Daemons is the number of built-in terraform provisioners. |
| `force_cancel_interval` | integer         | false    |              |                                                           |

## optimus-ide-collabsdk.ProvisionerDaemon

```json
{
  "api_version": "string",
  "created_at": "2019-08-24T14:15:22Z",
  "current_job": {
    "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
    "status": "pending",
    "template_display_name": "string",
    "template_icon": "string",
    "template_name": "string"
  },
  "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
  "key_id": "1e779c8a-6786-4c89-b7c3-a6666f5fd6b5",
  "key_name": "string",
  "last_seen_at": "2019-08-24T14:15:22Z",
  "name": "string",
  "organization_id": "7c60d51f-b44e-4682-87d6-449835ea4de6",
  "previous_job": {
    "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
    "status": "pending",
    "template_display_name": "string",
    "template_icon": "string",
    "template_name": "string"
  },
  "provisioners": [
    "string"
  ],
  "status": "offline",
  "tags": {
    "property1": "string",
    "property2": "string"
  },
  "version": "string"
}
```

### Properties

| Name               | Type                                                                 | Required | Restrictions | Description      |
|--------------------|----------------------------------------------------------------------|----------|--------------|------------------|
| `api_version`      | string                                                               | false    |              |                  |
| `created_at`       | string                                                               | false    |              |                  |
| `current_job`      | [optimus-ide-collabsdk.ProvisionerDaemonJob](#optimus-ide-collabsdkprovisionerdaemonjob)       | false    |              |                  |
| `id`               | string                                                               | false    |              |                  |
| `key_id`           | string                                                               | false    |              |                  |
| `key_name`         | string                                                               | false    |              | Optional fields. |
| `last_seen_at`     | string                                                               | false    |              |                  |
| `name`             | string                                                               | false    |              |                  |
| `organization_id`  | string                                                               | false    |              |                  |
| `previous_job`     | [optimus-ide-collabsdk.ProvisionerDaemonJob](#optimus-ide-collabsdkprovisionerdaemonjob)       | false    |              |                  |
| `provisioners`     | array of string                                                      | false    |              |                  |
| `status`           | [optimus-ide-collabsdk.ProvisionerDaemonStatus](#optimus-ide-collabsdkprovisionerdaemonstatus) | false    |              |                  |
| `tags`             | object                                                               | false    |              |                  |
| » `[any property]` | string                                                               | false    |              |                  |
| `version`          | string                                                               | false    |              |                  |

#### Enumerated Values

| Property | Value(s)                  |
|----------|---------------------------|
| `status` | `busy`, `idle`, `offline` |

## optimus-ide-collabsdk.ProvisionerDaemonJob

```json
{
  "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
  "status": "pending",
  "template_display_name": "string",
  "template_icon": "string",
  "template_name": "string"
}
```

### Properties

| Name                    | Type                                                           | Required | Restrictions | Description |
|-------------------------|----------------------------------------------------------------|----------|--------------|-------------|
| `id`                    | string                                                         | false    |              |             |
| `status`                | [optimus-ide-collabsdk.ProvisionerJobStatus](#optimus-ide-collabsdkprovisionerjobstatus) | false    |              |             |
| `template_display_name` | string                                                         | false    |              |             |
| `template_icon`         | string                                                         | false    |              |             |
| `template_name`         | string                                                         | false    |              |             |

#### Enumerated Values

| Property | Value(s)                                                             |
|----------|----------------------------------------------------------------------|
| `status` | `canceled`, `canceling`, `failed`, `pending`, `running`, `succeeded` |

## optimus-ide-collabsdk.ProvisionerDaemonStatus

```json
"offline"
```

### Properties

#### Enumerated Values

| Value(s)                  |
|---------------------------|
| `busy`, `idle`, `offline` |

## optimus-ide-collabsdk.ProvisionerJob

```json
{
  "available_workers": [
    "497f6eca-6276-4993-bfeb-53cbbbba6f08"
  ],
  "canceled_at": "2019-08-24T14:15:22Z",
  "completed_at": "2019-08-24T14:15:22Z",
  "created_at": "2019-08-24T14:15:22Z",
  "error": "string",
  "error_code": "REQUIRED_TEMPLATE_VARIABLES",
  "file_id": "8a0cfb4f-ddc9-436d-91bb-75133c583767",
  "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
  "initiator_id": "06588898-9a84-4b35-ba8f-f9cbd64946f3",
  "input": {
    "error": "string",
    "template_version_id": "0ba39c92-1f1b-4c32-aa3e-9925d7713eb1",
    "workspace_build_id": "badaf2eb-96c5-4050-9f1d-db2d39ca5478"
  },
  "logs_overflowed": true,
  "metadata": {
    "template_display_name": "string",
    "template_icon": "string",
    "template_id": "c6d67e98-83ea-49f0-8812-e4abae2b68bc",
    "template_name": "string",
    "template_version_name": "string",
    "workspace_build_transition": "start",
    "workspace_id": "0967198e-ec7b-4c6b-b4d3-f71244cadbe9",
    "workspace_name": "string"
  },
  "organization_id": "7c60d51f-b44e-4682-87d6-449835ea4de6",
  "queue_position": 0,
  "queue_size": 0,
  "started_at": "2019-08-24T14:15:22Z",
  "status": "pending",
  "tags": {
    "property1": "string",
    "property2": "string"
  },
  "type": "template_version_import",
  "worker_id": "ae5fa6f7-c55b-40c1-b40a-b36ac467652b",
  "worker_name": "string"
}
```

### Properties

| Name                | Type                                                               | Required | Restrictions | Description |
|---------------------|--------------------------------------------------------------------|----------|--------------|-------------|
| `available_workers` | array of string                                                    | false    |              |             |
| `canceled_at`       | string                                                             | false    |              |             |
| `completed_at`      | string                                                             | false    |              |             |
| `created_at`        | string                                                             | false    |              |             |
| `error`             | string                                                             | false    |              |             |
| `error_code`        | [optimus-ide-collabsdk.JobErrorCode](#optimus-ide-collabsdkjoberrorcode)                     | false    |              |             |
| `file_id`           | string                                                             | false    |              |             |
| `id`                | string                                                             | false    |              |             |
| `initiator_id`      | string                                                             | false    |              |             |
| `input`             | [optimus-ide-collabsdk.ProvisionerJobInput](#optimus-ide-collabsdkprovisionerjobinput)       | false    |              |             |
| `logs_overflowed`   | boolean                                                            | false    |              |             |
| `metadata`          | [optimus-ide-collabsdk.ProvisionerJobMetadata](#optimus-ide-collabsdkprovisionerjobmetadata) | false    |              |             |
| `organization_id`   | string                                                             | false    |              |             |
| `queue_position`    | integer                                                            | false    |              |             |
| `queue_size`        | integer                                                            | false    |              |             |
| `started_at`        | string                                                             | false    |              |             |
| `status`            | [optimus-ide-collabsdk.ProvisionerJobStatus](#optimus-ide-collabsdkprovisionerjobstatus)     | false    |              |             |
| `tags`              | object                                                             | false    |              |             |
| » `[any property]`  | string                                                             | false    |              |             |
| `type`              | [optimus-ide-collabsdk.ProvisionerJobType](#optimus-ide-collabsdkprovisionerjobtype)         | false    |              |             |
| `worker_id`         | string                                                             | false    |              |             |
| `worker_name`       | string                                                             | false    |              |             |

#### Enumerated Values

| Property     | Value(s)                                                             |
|--------------|----------------------------------------------------------------------|
| `error_code` | `INSUFFICIENT_QUOTA`, `REQUIRED_TEMPLATE_VARIABLES`                  |
| `status`     | `canceled`, `canceling`, `failed`, `pending`, `running`, `succeeded` |

## optimus-ide-collabsdk.ProvisionerJobInput

```json
{
  "error": "string",
  "template_version_id": "0ba39c92-1f1b-4c32-aa3e-9925d7713eb1",
  "workspace_build_id": "badaf2eb-96c5-4050-9f1d-db2d39ca5478"
}
```

### Properties

| Name                  | Type   | Required | Restrictions | Description |
|-----------------------|--------|----------|--------------|-------------|
| `error`               | string | false    |              |             |
| `template_version_id` | string | false    |              |             |
| `workspace_build_id`  | string | false    |              |             |

## optimus-ide-collabsdk.ProvisionerJobLog

```json
{
  "created_at": "2019-08-24T14:15:22Z",
  "id": 0,
  "log_level": "trace",
  "log_source": "provisioner_daemon",
  "output": "string",
  "stage": "string"
}
```

### Properties

| Name         | Type                                     | Required | Restrictions | Description |
|--------------|------------------------------------------|----------|--------------|-------------|
| `created_at` | string                                   | false    |              |             |
| `id`         | integer                                  | false    |              |             |
| `log_level`  | [optimus-ide-collabsdk.LogLevel](#optimus-ide-collabsdkloglevel)   | false    |              |             |
| `log_source` | [optimus-ide-collabsdk.LogSource](#optimus-ide-collabsdklogsource) | false    |              |             |
| `output`     | string                                   | false    |              |             |
| `stage`      | string                                   | false    |              |             |

#### Enumerated Values

| Property    | Value(s)                                  |
|-------------|-------------------------------------------|
| `log_level` | `debug`, `error`, `info`, `trace`, `warn` |

## optimus-ide-collabsdk.ProvisionerJobMetadata

```json
{
  "template_display_name": "string",
  "template_icon": "string",
  "template_id": "c6d67e98-83ea-49f0-8812-e4abae2b68bc",
  "template_name": "string",
  "template_version_name": "string",
  "workspace_build_transition": "start",
  "workspace_id": "0967198e-ec7b-4c6b-b4d3-f71244cadbe9",
  "workspace_name": "string"
}
```

### Properties

| Name                         | Type                                                         | Required | Restrictions | Description |
|------------------------------|--------------------------------------------------------------|----------|--------------|-------------|
| `template_display_name`      | string                                                       | false    |              |             |
| `template_icon`              | string                                                       | false    |              |             |
| `template_id`                | string                                                       | false    |              |             |
| `template_name`              | string                                                       | false    |              |             |
| `template_version_name`      | string                                                       | false    |              |             |
| `workspace_build_transition` | [optimus-ide-collabsdk.WorkspaceTransition](#optimus-ide-collabsdkworkspacetransition) | false    |              |             |
| `workspace_id`               | string                                                       | false    |              |             |
| `workspace_name`             | string                                                       | false    |              |             |

## optimus-ide-collabsdk.ProvisionerJobStatus

```json
"pending"
```

### Properties

#### Enumerated Values

| Value(s)                                                                        |
|---------------------------------------------------------------------------------|
| `canceled`, `canceling`, `failed`, `pending`, `running`, `succeeded`, `unknown` |

## optimus-ide-collabsdk.ProvisionerJobType

```json
"template_version_import"
```

### Properties

#### Enumerated Values

| Value(s)                                                                 |
|--------------------------------------------------------------------------|
| `template_version_dry_run`, `template_version_import`, `workspace_build` |

## optimus-ide-collabsdk.ProvisionerKey

```json
{
  "created_at": "2019-08-24T14:15:22Z",
  "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
  "name": "string",
  "organization": "452c1a86-a0af-475b-b03f-724878b0f387",
  "tags": {
    "property1": "string",
    "property2": "string"
  }
}
```

### Properties

| Name           | Type                                                       | Required | Restrictions | Description |
|----------------|------------------------------------------------------------|----------|--------------|-------------|
| `created_at`   | string                                                     | false    |              |             |
| `id`           | string                                                     | false    |              |             |
| `name`         | string                                                     | false    |              |             |
| `organization` | string                                                     | false    |              |             |
| `tags`         | [optimus-ide-collabsdk.ProvisionerKeyTags](#optimus-ide-collabsdkprovisionerkeytags) | false    |              |             |

## optimus-ide-collabsdk.ProvisionerKeyDaemons

```json
{
  "daemons": [
    {
      "api_version": "string",
      "created_at": "2019-08-24T14:15:22Z",
      "current_job": {
        "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
        "status": "pending",
        "template_display_name": "string",
        "template_icon": "string",
        "template_name": "string"
      },
      "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
      "key_id": "1e779c8a-6786-4c89-b7c3-a6666f5fd6b5",
      "key_name": "string",
      "last_seen_at": "2019-08-24T14:15:22Z",
      "name": "string",
      "organization_id": "7c60d51f-b44e-4682-87d6-449835ea4de6",
      "previous_job": {
        "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
        "status": "pending",
        "template_display_name": "string",
        "template_icon": "string",
        "template_name": "string"
      },
      "provisioners": [
        "string"
      ],
      "status": "offline",
      "tags": {
        "property1": "string",
        "property2": "string"
      },
      "version": "string"
    }
  ],
  "key": {
    "created_at": "2019-08-24T14:15:22Z",
    "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
    "name": "string",
    "organization": "452c1a86-a0af-475b-b03f-724878b0f387",
    "tags": {
      "property1": "string",
      "property2": "string"
    }
  }
}
```

### Properties

| Name      | Type                                                              | Required | Restrictions | Description |
|-----------|-------------------------------------------------------------------|----------|--------------|-------------|
| `daemons` | array of [optimus-ide-collabsdk.ProvisionerDaemon](#optimus-ide-collabsdkprovisionerdaemon) | false    |              |             |
| `key`     | [optimus-ide-collabsdk.ProvisionerKey](#optimus-ide-collabsdkprovisionerkey)                | false    |              |             |

## optimus-ide-collabsdk.ProvisionerKeyTags

```json
{
  "property1": "string",
  "property2": "string"
}
```

### Properties

| Name             | Type   | Required | Restrictions | Description |
|------------------|--------|----------|--------------|-------------|
| `[any property]` | string | false    |              |             |

## optimus-ide-collabsdk.ProvisionerLogLevel

```json
"debug"
```

### Properties

#### Enumerated Values

| Value(s) |
|----------|
| `debug`  |

## optimus-ide-collabsdk.ProvisionerStorageMethod

```json
"file"
```

### Properties

#### Enumerated Values

| Value(s) |
|----------|
| `file`   |

## optimus-ide-collabsdk.ProvisionerTiming

```json
{
  "action": "string",
  "ended_at": "2019-08-24T14:15:22Z",
  "job_id": "453bd7d7-5355-4d6d-a38e-d9e7eb218c3f",
  "resource": "string",
  "source": "string",
  "stage": "init",
  "started_at": "2019-08-24T14:15:22Z"
}
```

### Properties

| Name         | Type                                         | Required | Restrictions | Description |
|--------------|----------------------------------------------|----------|--------------|-------------|
| `action`     | string                                       | false    |              |             |
| `ended_at`   | string                                       | false    |              |             |
| `job_id`     | string                                       | false    |              |             |
| `resource`   | string                                       | false    |              |             |
| `source`     | string                                       | false    |              |             |
| `stage`      | [optimus-ide-collabsdk.TimingStage](#optimus-ide-collabsdktimingstage) | false    |              |             |
| `started_at` | string                                       | false    |              |             |

## optimus-ide-collabsdk.ProxyHealthReport

```json
{
  "errors": [
    "string"
  ],
  "warnings": [
    "string"
  ]
}
```

### Properties

| Name       | Type            | Required | Restrictions | Description                                                                              |
|------------|-----------------|----------|--------------|------------------------------------------------------------------------------------------|
| `errors`   | array of string | false    |              | Errors are problems that prevent the workspace proxy from being healthy                  |
| `warnings` | array of string | false    |              | Warnings do not prevent the workspace proxy from being healthy, but should be addressed. |

## optimus-ide-collabsdk.ProxyHealthStatus

```json
"ok"
```

### Properties

#### Enumerated Values

| Value(s)                                         |
|--------------------------------------------------|
| `ok`, `unhealthy`, `unreachable`, `unregistered` |

## optimus-ide-collabsdk.PutExtendWorkspaceRequest

```json
{
  "deadline": "2019-08-24T14:15:22Z"
}
```

### Properties

| Name       | Type   | Required | Restrictions | Description |
|------------|--------|----------|--------------|-------------|
| `deadline` | string | true     |              |             |

## optimus-ide-collabsdk.PutOAuth2ProviderAppRequest

```json
{
  "callback_url": "string",
  "icon": "string",
  "name": "string"
}
```

### Properties

| Name           | Type   | Required | Restrictions | Description |
|----------------|--------|----------|--------------|-------------|
| `callback_url` | string | true     |              |             |
| `icon`         | string | false    |              |             |
| `name`         | string | true     |              |             |

## optimus-ide-collabsdk.RBACAction

```json
"application_connect"
```

### Properties

#### Enumerated Values

| Value(s)                                                                                                                                                                                                                       |
|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `application_connect`, `assign`, `create`, `create_agent`, `delete`, `delete_agent`, `read`, `read_personal`, `share`, `ssh`, `start`, `stop`, `unassign`, `update`, `update_agent`, `update_personal`, `use`, `view_insights` |

## optimus-ide-collabsdk.RBACResource

```json
"*"
```

### Properties

#### Enumerated Values

| Value(s)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                           |
|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `*`, `ai_gateway_key`, `ai_model_price`, `ai_provider`, `ai_seat`, `aibridge_interception`, `api_key`, `assign_org_role`, `assign_role`, `audit_log`, `boundary_log`, `boundary_usage`, `chat`, `connection_log`, `crypto_key`, `debug_info`, `deployment_config`, `deployment_stats`, `file`, `group`, `group_member`, `idpsync_settings`, `inbox_notification`, `license`, `notification_message`, `notification_preference`, `notification_template`, `oauth2_app`, `oauth2_app_code_token`, `oauth2_app_secret`, `organization`, `organization_member`, `prebuilt_workspace`, `provisioner_daemon`, `provisioner_jobs`, `replicas`, `system`, `tailnet_coordinator`, `task`, `template`, `usage_event`, `user`, `user_secret`, `user_skill`, `webpush_subscription`, `workspace`, `workspace_agent_devcontainers`, `workspace_agent_resource_monitor`, `workspace_build_orchestration`, `workspace_dormant`, `workspace_proxy` |

## optimus-ide-collabsdk.RateLimitConfig

```json
{
  "api": 0,
  "disable_all": true
}
```

### Properties

| Name          | Type    | Required | Restrictions | Description |
|---------------|---------|----------|--------------|-------------|
| `api`         | integer | false    |              |             |
| `disable_all` | boolean | false    |              |             |

## optimus-ide-collabsdk.ReducedUser

```json
{
  "avatar_url": "http://example.com",
  "created_at": "2019-08-24T14:15:22Z",
  "email": "user@example.com",
  "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
  "is_service_account": true,
  "last_seen_at": "2019-08-24T14:15:22Z",
  "login_type": "",
  "name": "string",
  "status": "active",
  "theme_preference": "string",
  "updated_at": "2019-08-24T14:15:22Z",
  "username": "string"
}
```

### Properties

| Name                 | Type                                       | Required | Restrictions | Description                                                                                |
|----------------------|--------------------------------------------|----------|--------------|--------------------------------------------------------------------------------------------|
| `avatar_url`         | string                                     | false    |              |                                                                                            |
| `created_at`         | string                                     | true     |              |                                                                                            |
| `email`              | string                                     | true     |              |                                                                                            |
| `id`                 | string                                     | true     |              |                                                                                            |
| `is_service_account` | boolean                                    | false    |              |                                                                                            |
| `last_seen_at`       | string                                     | false    |              |                                                                                            |
| `login_type`         | [optimus-ide-collabsdk.LoginType](#optimus-ide-collabsdklogintype)   | false    |              |                                                                                            |
| `name`               | string                                     | false    |              |                                                                                            |
| `status`             | [optimus-ide-collabsdk.UserStatus](#optimus-ide-collabsdkuserstatus) | false    |              |                                                                                            |
| `theme_preference`   | string                                     | false    |              | Deprecated: this value should be retrieved from `optimus-ide-collabsdk.UserPreferenceSettings` instead. |
| `updated_at`         | string                                     | false    |              |                                                                                            |
| `username`           | string                                     | true     |              |                                                                                            |

#### Enumerated Values

| Property | Value(s)              |
|----------|-----------------------|
| `status` | `active`, `suspended` |

## optimus-ide-collabsdk.Region

```json
{
  "display_name": "string",
  "healthy": true,
  "icon_url": "string",
  "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
  "name": "string",
  "path_app_url": "string",
  "wildcard_hostname": "string"
}
```

### Properties

| Name                | Type    | Required | Restrictions | Description                                                                                                                                                                       |
|---------------------|---------|----------|--------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `display_name`      | string  | false    |              |                                                                                                                                                                                   |
| `healthy`           | boolean | false    |              |                                                                                                                                                                                   |
| `icon_url`          | string  | false    |              |                                                                                                                                                                                   |
| `id`                | string  | false    |              |                                                                                                                                                                                   |
| `name`              | string  | false    |              |                                                                                                                                                                                   |
| `path_app_url`      | string  | false    |              | Path app URL is the URL to the base path for path apps. Optional unless wildcard_hostname is set. E.g. https://us.example.com                                                     |
| `wildcard_hostname` | string  | false    |              | Wildcard hostname is the wildcard hostname for subdomain apps. E.g. *.us.example.com E.g.*--suffix.au.example.com Optional. Does not need to be on the same domain as PathAppURL. |

## optimus-ide-collabsdk.RegionsResponse-optimus-ide-collabsdk_Region

```json
{
  "regions": [
    {
      "display_name": "string",
      "healthy": true,
      "icon_url": "string",
      "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
      "name": "string",
      "path_app_url": "string",
      "wildcard_hostname": "string"
    }
  ]
}
```

### Properties

| Name      | Type                                        | Required | Restrictions | Description |
|-----------|---------------------------------------------|----------|--------------|-------------|
| `regions` | array of [optimus-ide-collabsdk.Region](#optimus-ide-collabsdkregion) | false    |              |             |

## optimus-ide-collabsdk.RegionsResponse-optimus-ide-collabsdk_WorkspaceProxy

```json
{
  "regions": [
    {
      "created_at": "2019-08-24T14:15:22Z",
      "deleted": true,
      "derp_enabled": true,
      "derp_only": true,
      "display_name": "string",
      "healthy": true,
      "icon_url": "string",
      "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
      "name": "string",
      "path_app_url": "string",
      "status": {
        "checked_at": "2019-08-24T14:15:22Z",
        "report": {
          "errors": [
            "string"
          ],
          "warnings": [
            "string"
          ]
        },
        "status": "ok"
      },
      "updated_at": "2019-08-24T14:15:22Z",
      "version": "string",
      "wildcard_hostname": "string"
    }
  ]
}
```

### Properties

| Name      | Type                                                        | Required | Restrictions | Description |
|-----------|-------------------------------------------------------------|----------|--------------|-------------|
| `regions` | array of [optimus-ide-collabsdk.WorkspaceProxy](#optimus-ide-collabsdkworkspaceproxy) | false    |              |             |

## optimus-ide-collabsdk.Replica

```json
{
  "created_at": "2019-08-24T14:15:22Z",
  "database_latency": 0,
  "error": "string",
  "hostname": "string",
  "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
  "region_id": 0,
  "relay_address": "string"
}
```

### Properties

| Name               | Type    | Required | Restrictions | Description                                                        |
|--------------------|---------|----------|--------------|--------------------------------------------------------------------|
| `created_at`       | string  | false    |              | Created at is the timestamp when the replica was first seen.       |
| `database_latency` | integer | false    |              | Database latency is the latency in microseconds to the database.   |
| `error`            | string  | false    |              | Error is the replica error.                                        |
| `hostname`         | string  | false    |              | Hostname is the hostname of the replica.                           |
| `id`               | string  | false    |              | ID is the unique identifier for the replica.                       |
| `region_id`        | integer | false    |              | Region ID is the region of the replica.                            |
| `relay_address`    | string  | false    |              | Relay address is the accessible address to relay DERP connections. |

## optimus-ide-collabsdk.RequestOneTimePasscodeRequest

```json
{
  "email": "user@example.com"
}
```

### Properties

| Name    | Type   | Required | Restrictions | Description |
|---------|--------|----------|--------------|-------------|
| `email` | string | true     |              |             |

## optimus-ide-collabsdk.ResolveAutostartResponse

```json
{
  "parameter_mismatch": true
}
```

### Properties

| Name                 | Type    | Required | Restrictions | Description |
|----------------------|---------|----------|--------------|-------------|
| `parameter_mismatch` | boolean | false    |              |             |

## optimus-ide-collabsdk.ResourceType

```json
"template"
```

### Properties

#### Enumerated Values

| Value(s)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                         |
|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `ai_gateway_key`, `ai_provider`, `ai_provider_key`, `ai_seat`, `api_key`, `chat`, `convert_login`, `custom_role`, `git_ssh_key`, `group`, `group_ai_budget`, `health_settings`, `idp_sync_settings_group`, `idp_sync_settings_organization`, `idp_sync_settings_role`, `license`, `notification_template`, `notifications_settings`, `oauth2_provider_app`, `oauth2_provider_app_secret`, `oauth2_provider_settings`, `organization`, `organization_member`, `prebuilds_settings`, `task`, `template`, `template_version`, `user`, `user_ai_budget_override`, `user_secret`, `user_skill`, `workspace`, `workspace_agent`, `workspace_app`, `workspace_build`, `workspace_proxy` |

## optimus-ide-collabsdk.Response

```json
{
  "detail": "string",
  "message": "string",
  "validations": [
    {
      "detail": "string",
      "field": "string"
    }
  ]
}
```

### Properties

| Name          | Type                                                          | Required | Restrictions | Description                                                                                                                                                                                                                        |
|---------------|---------------------------------------------------------------|----------|--------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `detail`      | string                                                        | false    |              | Detail is a debug message that provides further insight into why the action failed. This information can be technical and a regular golang err.Error() text. - "database: too many open connections" - "stat: too many open files" |
| `message`     | string                                                        | false    |              | Message is an actionable message that depicts actions the request took. These messages should be fully formed sentences with proper punctuation. Examples: - "A user has been created." - "Failed to create a user."               |
| `validations` | array of [optimus-ide-collabsdk.ValidationError](#optimus-ide-collabsdkvalidationerror) | false    |              | Validations are form field-specific friendly error messages. They will be shown on a form field in the UI. These can also be used to add additional context if there is a set of errors in the primary 'Message'.                  |

## optimus-ide-collabsdk.ResumeTaskResponse

```json
{
  "workspace_build": {
    "build_number": 0,
    "created_at": "2019-08-24T14:15:22Z",
    "daily_cost": 0,
    "deadline": "2019-08-24T14:15:22Z",
    "has_ai_task": true,
    "has_external_agent": true,
    "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
    "initiator_id": "06588898-9a84-4b35-ba8f-f9cbd64946f3",
    "initiator_name": "string",
    "job": {
      "available_workers": [
        "497f6eca-6276-4993-bfeb-53cbbbba6f08"
      ],
      "canceled_at": "2019-08-24T14:15:22Z",
      "completed_at": "2019-08-24T14:15:22Z",
      "created_at": "2019-08-24T14:15:22Z",
      "error": "string",
      "error_code": "REQUIRED_TEMPLATE_VARIABLES",
      "file_id": "8a0cfb4f-ddc9-436d-91bb-75133c583767",
      "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
      "initiator_id": "06588898-9a84-4b35-ba8f-f9cbd64946f3",
      "input": {
        "error": "string",
        "template_version_id": "0ba39c92-1f1b-4c32-aa3e-9925d7713eb1",
        "workspace_build_id": "badaf2eb-96c5-4050-9f1d-db2d39ca5478"
      },
      "logs_overflowed": true,
      "metadata": {
        "template_display_name": "string",
        "template_icon": "string",
        "template_id": "c6d67e98-83ea-49f0-8812-e4abae2b68bc",
        "template_name": "string",
        "template_version_name": "string",
        "workspace_build_transition": "start",
        "workspace_id": "0967198e-ec7b-4c6b-b4d3-f71244cadbe9",
        "workspace_name": "string"
      },
      "organization_id": "7c60d51f-b44e-4682-87d6-449835ea4de6",
      "queue_position": 0,
      "queue_size": 0,
      "started_at": "2019-08-24T14:15:22Z",
      "status": "pending",
      "tags": {
        "property1": "string",
        "property2": "string"
      },
      "type": "template_version_import",
      "worker_id": "ae5fa6f7-c55b-40c1-b40a-b36ac467652b",
      "worker_name": "string"
    },
    "matched_provisioners": {
      "available": 0,
      "count": 0,
      "most_recently_seen": "2019-08-24T14:15:22Z"
    },
    "max_deadline": "2019-08-24T14:15:22Z",
    "reason": "initiator",
    "resources": [
      {
        "agents": [
          {
            "api_version": "string",
            "apps": [
              {
                "command": "string",
                "display_name": "string",
                "external": true,
                "group": "string",
                "health": "disabled",
                "healthcheck": {
                  "interval": 0,
                  "threshold": 0,
                  "url": "string"
                },
                "hidden": true,
                "icon": "string",
                "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
                "open_in": "slim-window",
                "sharing_level": "owner",
                "slug": "string",
                "statuses": [
                  {
                    "agent_id": "2b1e3b65-2c04-4fa2-a2d7-467901e98978",
                    "app_id": "affd1d10-9538-4fc8-9e0b-4594a28c1335",
                    "created_at": "2019-08-24T14:15:22Z",
                    "icon": "string",
                    "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
                    "message": "string",
                    "needs_user_attention": true,
                    "state": "working",
                    "uri": "string",
                    "workspace_id": "0967198e-ec7b-4c6b-b4d3-f71244cadbe9"
                  }
                ],
                "subdomain": true,
                "subdomain_name": "string",
                "tooltip": "string",
                "url": "string"
              }
            ],
            "architecture": "string",
            "connection_timeout_seconds": 0,
            "created_at": "2019-08-24T14:15:22Z",
            "directory": "string",
            "disconnected_at": "2019-08-24T14:15:22Z",
            "display_apps": [
              "vscode"
            ],
            "environment_variables": {
              "property1": "string",
              "property2": "string"
            },
            "expanded_directory": "string",
            "first_connected_at": "2019-08-24T14:15:22Z",
            "health": {
              "healthy": false,
              "reason": "agent has lost connection"
            },
            "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
            "instance_id": "string",
            "last_connected_at": "2019-08-24T14:15:22Z",
            "latency": {
              "property1": {
                "latency_ms": 0,
                "preferred": true
              },
              "property2": {
                "latency_ms": 0,
                "preferred": true
              }
            },
            "lifecycle_state": "created",
            "log_sources": [
              {
                "created_at": "2019-08-24T14:15:22Z",
                "display_name": "string",
                "icon": "string",
                "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
                "workspace_agent_id": "7ad2e618-fea7-4c1a-b70a-f501566a72f1"
              }
            ],
            "logs_length": 0,
            "logs_overflowed": true,
            "name": "string",
            "operating_system": "string",
            "parent_id": {
              "uuid": "string",
              "valid": true
            },
            "ready_at": "2019-08-24T14:15:22Z",
            "resource_id": "4d5215ed-38bb-48ed-879a-fdb9ca58522f",
            "scripts": [
              {
                "cron": "string",
                "display_name": "string",
                "exit_code": 0,
                "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
                "log_path": "string",
                "log_source_id": "4197ab25-95cf-4b91-9c78-f7f2af5d353a",
                "run_on_start": true,
                "run_on_stop": true,
                "script": "string",
                "start_blocks_login": true,
                "status": "ok",
                "timeout": 0
              }
            ],
            "started_at": "2019-08-24T14:15:22Z",
            "startup_script_behavior": "blocking",
            "status": "connecting",
            "subsystems": [
              "envbox"
            ],
            "troubleshooting_url": "string",
            "updated_at": "2019-08-24T14:15:22Z",
            "version": "string"
          }
        ],
        "created_at": "2019-08-24T14:15:22Z",
        "daily_cost": 0,
        "hide": true,
        "icon": "string",
        "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
        "job_id": "453bd7d7-5355-4d6d-a38e-d9e7eb218c3f",
        "metadata": [
          {
            "key": "string",
            "sensitive": true,
            "value": "string"
          }
        ],
        "name": "string",
        "type": "string",
        "workspace_transition": "start"
      }
    ],
    "status": "pending",
    "template_version_id": "0ba39c92-1f1b-4c32-aa3e-9925d7713eb1",
    "template_version_name": "string",
    "template_version_preset_id": "512a53a7-30da-446e-a1fc-713c630baff1",
    "transition": "start",
    "updated_at": "2019-08-24T14:15:22Z",
    "workspace_id": "0967198e-ec7b-4c6b-b4d3-f71244cadbe9",
    "workspace_name": "string",
    "workspace_owner_avatar_url": "string",
    "workspace_owner_id": "e7078695-5279-4c86-8774-3ac2367a2fc7",
    "workspace_owner_name": "string"
  }
}
```

### Properties

| Name              | Type                                               | Required | Restrictions | Description |
|-------------------|----------------------------------------------------|----------|--------------|-------------|
| `workspace_build` | [optimus-ide-collabsdk.WorkspaceBuild](#optimus-ide-collabsdkworkspacebuild) | false    |              |             |

## optimus-ide-collabsdk.RetentionConfig

```json
{
  "api_keys": 0,
  "audit_logs": 0,
  "boundary_logs": 0,
  "connection_logs": 0,
  "workspace_agent_logs": 0
}
```

### Properties

| Name                   | Type    | Required | Restrictions | Description                                                                                                                                                                                                                                                                          |
|------------------------|---------|----------|--------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `api_keys`             | integer | false    |              | Api keys controls how long expired API keys are retained before being deleted. Keys are only deleted if they have been expired for at least this duration. Defaults to 7 days to preserve existing behavior.                                                                         |
| `audit_logs`           | integer | false    |              | Audit logs controls how long audit log entries are retained. Set to 0 to disable (keep indefinitely).                                                                                                                                                                                |
| `boundary_logs`        | integer | false    |              | Boundary logs controls how long boundary audit log entries are retained. Boundary logs record every HTTP request processed by a Boundary confinement proxy. Set to 0 to disable automatic deletion (keep indefinitely). Adjust to match your organization's regulatory requirements. |
| `connection_logs`      | integer | false    |              | Connection logs controls how long connection log entries are retained. Set to 0 to disable (keep indefinitely).                                                                                                                                                                      |
| `workspace_agent_logs` | integer | false    |              | Workspace agent logs controls how long workspace agent logs are retained. Logs are deleted if the agent hasn't connected within this period. Logs from the latest build are always retained regardless of age. Defaults to 7 days to preserve existing behavior.                     |

## optimus-ide-collabsdk.Role

```json
{
  "display_name": "string",
  "name": "string",
  "organization_id": "7c60d51f-b44e-4682-87d6-449835ea4de6",
  "organization_member_permissions": [
    {
      "action": "application_connect",
      "negate": true,
      "resource_type": "*"
    }
  ],
  "organization_permissions": [
    {
      "action": "application_connect",
      "negate": true,
      "resource_type": "*"
    }
  ],
  "site_permissions": [
    {
      "action": "application_connect",
      "negate": true,
      "resource_type": "*"
    }
  ],
  "user_permissions": [
    {
      "action": "application_connect",
      "negate": true,
      "resource_type": "*"
    }
  ]
}
```

### Properties

| Name                              | Type                                                | Required | Restrictions | Description                                                                                            |
|-----------------------------------|-----------------------------------------------------|----------|--------------|--------------------------------------------------------------------------------------------------------|
| `display_name`                    | string                                              | false    |              |                                                                                                        |
| `name`                            | string                                              | false    |              |                                                                                                        |
| `organization_id`                 | string                                              | false    |              |                                                                                                        |
| `organization_member_permissions` | array of [optimus-ide-collabsdk.Permission](#optimus-ide-collabsdkpermission) | false    |              | Organization member permissions are specific for the organization in the field 'OrganizationID' above. |
| `organization_permissions`        | array of [optimus-ide-collabsdk.Permission](#optimus-ide-collabsdkpermission) | false    |              | Organization permissions are specific for the organization in the field 'OrganizationID' above.        |
| `site_permissions`                | array of [optimus-ide-collabsdk.Permission](#optimus-ide-collabsdkpermission) | false    |              |                                                                                                        |
| `user_permissions`                | array of [optimus-ide-collabsdk.Permission](#optimus-ide-collabsdkpermission) | false    |              |                                                                                                        |

## optimus-ide-collabsdk.RoleSyncSettings

```json
{
  "field": "string",
  "mapping": {
    "property1": [
      "string"
    ],
    "property2": [
      "string"
    ]
  }
}
```

### Properties

| Name               | Type            | Required | Restrictions | Description                                                                                                                            |
|--------------------|-----------------|----------|--------------|----------------------------------------------------------------------------------------------------------------------------------------|
| `field`            | string          | false    |              | Field is the name of the claim field that specifies what organization roles a user should be given. If empty, no roles will be synced. |
| `mapping`          | object          | false    |              | Mapping is a map from OIDC groups to Optimus-IDE-Collab organization roles.                                                                         |
| » `[any property]` | array of string | false    |              |                                                                                                                                        |

## optimus-ide-collabsdk.SSHConfig

```json
{
  "deploymentName": "string",
  "sshconfigOptions": [
    "string"
  ]
}
```

### Properties

| Name               | Type            | Required | Restrictions | Description                                                                                         |
|--------------------|-----------------|----------|--------------|-----------------------------------------------------------------------------------------------------|
| `deploymentName`   | string          | false    |              | Deploymentname is the config-ssh Hostname prefix                                                    |
| `sshconfigOptions` | array of string | false    |              | Sshconfigoptions are additional options to add to the ssh config file. This will override defaults. |

## optimus-ide-collabsdk.SSHConfigResponse

```json
{
  "hostname_prefix": "string",
  "hostname_suffix": "string",
  "ssh_config_options": {
    "property1": "string",
    "property2": "string"
  }
}
```

### Properties

| Name                 | Type   | Required | Restrictions | Description                                                                                                           |
|----------------------|--------|----------|--------------|-----------------------------------------------------------------------------------------------------------------------|
| `hostname_prefix`    | string | false    |              | Hostname prefix is the prefix we append to workspace names for SSH hostnames. Deprecated: use HostnameSuffix instead. |
| `hostname_suffix`    | string | false    |              | Hostname suffix is the suffix to append to workspace names for SSH hostnames.                                         |
| `ssh_config_options` | object | false    |              |                                                                                                                       |
| » `[any property]`   | string | false    |              |                                                                                                                       |

## optimus-ide-collabsdk.SecretsFileFormat

```json
"env"
```

### Properties

#### Enumerated Values

| Value(s)              |
|-----------------------|
| `env`, `json`, `yaml` |

## optimus-ide-collabsdk.ServerSentEvent

```json
{
  "data": null,
  "type": "ping"
}
```

### Properties

| Name   | Type                                                         | Required | Restrictions | Description |
|--------|--------------------------------------------------------------|----------|--------------|-------------|
| `data` | any                                                          | false    |              |             |
| `type` | [optimus-ide-collabsdk.ServerSentEventType](#optimus-ide-collabsdkserversenteventtype) | false    |              |             |

## optimus-ide-collabsdk.ServerSentEventType

```json
"ping"
```

### Properties

#### Enumerated Values

| Value(s)                |
|-------------------------|
| `data`, `error`, `ping` |

## optimus-ide-collabsdk.SessionCountDeploymentStats

```json
{
  "jetbrains": 0,
  "reconnecting_pty": 0,
  "ssh": 0,
  "vscode": 0
}
```

### Properties

| Name               | Type    | Required | Restrictions | Description |
|--------------------|---------|----------|--------------|-------------|
| `jetbrains`        | integer | false    |              |             |
| `reconnecting_pty` | integer | false    |              |             |
| `ssh`              | integer | false    |              |             |
| `vscode`           | integer | false    |              |             |

## optimus-ide-collabsdk.SessionLifetime

```json
{
  "default_duration": 0,
  "default_token_lifetime": 0,
  "disable_expiry_refresh": true,
  "max_admin_token_lifetime": 0,
  "max_token_lifetime": 0,
  "refresh_default_duration": 0
}
```

### Properties

| Name                       | Type    | Required | Restrictions | Description                                                                                                                                                                            |
|----------------------------|---------|----------|--------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `default_duration`         | integer | false    |              | Default duration is only for browser, workspace app and oauth sessions.                                                                                                                |
| `default_token_lifetime`   | integer | false    |              |                                                                                                                                                                                        |
| `disable_expiry_refresh`   | boolean | false    |              | Disable expiry refresh will disable automatically refreshing api keys when they are used from the api. This means the api key lifetime at creation is the lifetime of the api key.     |
| `max_admin_token_lifetime` | integer | false    |              |                                                                                                                                                                                        |
| `max_token_lifetime`       | integer | false    |              |                                                                                                                                                                                        |
| `refresh_default_duration` | integer | false    |              | Refresh default duration is the default lifetime for OAuth2 refresh tokens. This should generally be longer than access token lifetimes to allow refreshing after access token expiry. |

## optimus-ide-collabsdk.ShareableWorkspaceOwners

```json
"none"
```

### Properties

#### Enumerated Values

| Value(s)                               |
|----------------------------------------|
| `everyone`, `none`, `service_accounts` |

## optimus-ide-collabsdk.SharedWorkspaceActor

```json
{
  "actor_type": "group",
  "avatar_url": "http://example.com",
  "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
  "name": "string",
  "roles": [
    "admin"
  ]
}
```

### Properties

| Name         | Type                                                                   | Required | Restrictions | Description |
|--------------|------------------------------------------------------------------------|----------|--------------|-------------|
| `actor_type` | [optimus-ide-collabsdk.SharedWorkspaceActorType](#optimus-ide-collabsdksharedworkspaceactortype) | false    |              |             |
| `avatar_url` | string                                                                 | false    |              |             |
| `id`         | string                                                                 | false    |              |             |
| `name`       | string                                                                 | false    |              |             |
| `roles`      | array of [optimus-ide-collabsdk.WorkspaceRole](#optimus-ide-collabsdkworkspacerole)              | false    |              |             |

#### Enumerated Values

| Property     | Value(s)        |
|--------------|-----------------|
| `actor_type` | `group`, `user` |

## optimus-ide-collabsdk.SharedWorkspaceActorType

```json
"group"
```

### Properties

#### Enumerated Values

| Value(s)        |
|-----------------|
| `group`, `user` |

## optimus-ide-collabsdk.SlimRole

```json
{
  "display_name": "string",
  "name": "string",
  "organization_id": "string"
}
```

### Properties

| Name              | Type   | Required | Restrictions | Description |
|-------------------|--------|----------|--------------|-------------|
| `display_name`    | string | false    |              |             |
| `name`            | string | false    |              |             |
| `organization_id` | string | false    |              |             |

## optimus-ide-collabsdk.StatsCollectionConfig

```json
{
  "usage_stats": {
    "enable": true
  }
}
```

### Properties

| Name          | Type                                                   | Required | Restrictions | Description |
|---------------|--------------------------------------------------------|----------|--------------|-------------|
| `usage_stats` | [optimus-ide-collabsdk.UsageStatsConfig](#optimus-ide-collabsdkusagestatsconfig) | false    |              |             |

## optimus-ide-collabsdk.SupportConfig

```json
{
  "links": {
    "value": [
      {
        "icon": "bug",
        "location": "navbar",
        "name": "string",
        "target": "string"
      }
    ]
  }
}
```

### Properties

| Name    | Type                                                                                 | Required | Restrictions | Description |
|---------|--------------------------------------------------------------------------------------|----------|--------------|-------------|
| `links` | [serpent.Struct-array_optimus-ide-collabsdk_LinkConfig](#serpentstruct-array_optimus-ide-collabsdk_linkconfig) | false    |              |             |

## optimus-ide-collabsdk.SwaggerConfig

```json
{
  "enable": true
}
```

### Properties

| Name     | Type    | Required | Restrictions | Description |
|----------|---------|----------|--------------|-------------|
| `enable` | boolean | false    |              |             |

## optimus-ide-collabsdk.TLSConfig

```json
{
  "address": {
    "host": "string",
    "port": "string"
  },
  "allow_insecure_ciphers": true,
  "cert_file": [
    "string"
  ],
  "client_auth": "string",
  "client_ca_file": "string",
  "client_cert_file": "string",
  "client_key_file": "string",
  "enable": true,
  "key_file": [
    "string"
  ],
  "min_version": "string",
  "redirect_http": true,
  "supported_ciphers": [
    "string"
  ]
}
```

### Properties

| Name                     | Type                                 | Required | Restrictions | Description |
|--------------------------|--------------------------------------|----------|--------------|-------------|
| `address`                | [serpent.HostPort](#serpenthostport) | false    |              |             |
| `allow_insecure_ciphers` | boolean                              | false    |              |             |
| `cert_file`              | array of string                      | false    |              |             |
| `client_auth`            | string                               | false    |              |             |
| `client_ca_file`         | string                               | false    |              |             |
| `client_cert_file`       | string                               | false    |              |             |
| `client_key_file`        | string                               | false    |              |             |
| `enable`                 | boolean                              | false    |              |             |
| `key_file`               | array of string                      | false    |              |             |
| `min_version`            | string                               | false    |              |             |
| `redirect_http`          | boolean                              | false    |              |             |
| `supported_ciphers`      | array of string                      | false    |              |             |

## optimus-ide-collabsdk.Task

```json
{
  "created_at": "2019-08-24T14:15:22Z",
  "current_state": {
    "message": "string",
    "state": "working",
    "timestamp": "2019-08-24T14:15:22Z",
    "uri": "string"
  },
  "display_name": "string",
  "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
  "initial_prompt": "string",
  "name": "string",
  "organization_id": "7c60d51f-b44e-4682-87d6-449835ea4de6",
  "owner_avatar_url": "string",
  "owner_id": "8826ee2e-7933-4665-aef2-2393f84a0d05",
  "owner_name": "string",
  "status": "pending",
  "template_display_name": "string",
  "template_icon": "string",
  "template_id": "c6d67e98-83ea-49f0-8812-e4abae2b68bc",
  "template_name": "string",
  "template_version_id": "0ba39c92-1f1b-4c32-aa3e-9925d7713eb1",
  "updated_at": "2019-08-24T14:15:22Z",
  "workspace_agent_health": {
    "healthy": false,
    "reason": "agent has lost connection"
  },
  "workspace_agent_id": {
    "uuid": "string",
    "valid": true
  },
  "workspace_agent_lifecycle": "created",
  "workspace_app_id": {
    "uuid": "string",
    "valid": true
  },
  "workspace_build_number": 0,
  "workspace_id": {
    "uuid": "string",
    "valid": true
  },
  "workspace_name": "string",
  "workspace_status": "pending"
}
```

### Properties

| Name                        | Type                                                                 | Required | Restrictions | Description |
|-----------------------------|----------------------------------------------------------------------|----------|--------------|-------------|
| `created_at`                | string                                                               | false    |              |             |
| `current_state`             | [optimus-ide-collabsdk.TaskStateEntry](#optimus-ide-collabsdktaskstateentry)                   | false    |              |             |
| `display_name`              | string                                                               | false    |              |             |
| `id`                        | string                                                               | false    |              |             |
| `initial_prompt`            | string                                                               | false    |              |             |
| `name`                      | string                                                               | false    |              |             |
| `organization_id`           | string                                                               | false    |              |             |
| `owner_avatar_url`          | string                                                               | false    |              |             |
| `owner_id`                  | string                                                               | false    |              |             |
| `owner_name`                | string                                                               | false    |              |             |
| `status`                    | [optimus-ide-collabsdk.TaskStatus](#optimus-ide-collabsdktaskstatus)                           | false    |              |             |
| `template_display_name`     | string                                                               | false    |              |             |
| `template_icon`             | string                                                               | false    |              |             |
| `template_id`               | string                                                               | false    |              |             |
| `template_name`             | string                                                               | false    |              |             |
| `template_version_id`       | string                                                               | false    |              |             |
| `updated_at`                | string                                                               | false    |              |             |
| `workspace_agent_health`    | [optimus-ide-collabsdk.WorkspaceAgentHealth](#optimus-ide-collabsdkworkspaceagenthealth)       | false    |              |             |
| `workspace_agent_id`        | [uuid.NullUUID](#uuidnulluuid)                                       | false    |              |             |
| `workspace_agent_lifecycle` | [optimus-ide-collabsdk.WorkspaceAgentLifecycle](#optimus-ide-collabsdkworkspaceagentlifecycle) | false    |              |             |
| `workspace_app_id`          | [uuid.NullUUID](#uuidnulluuid)                                       | false    |              |             |
| `workspace_build_number`    | integer                                                              | false    |              |             |
| `workspace_id`              | [uuid.NullUUID](#uuidnulluuid)                                       | false    |              |             |
| `workspace_name`            | string                                                               | false    |              |             |
| `workspace_status`          | [optimus-ide-collabsdk.WorkspaceStatus](#optimus-ide-collabsdkworkspacestatus)                 | false    |              |             |

#### Enumerated Values

| Property           | Value(s)                                                                                                          |
|--------------------|-------------------------------------------------------------------------------------------------------------------|
| `status`           | `active`, `error`, `initializing`, `paused`, `pending`, `unknown`                                                 |
| `workspace_status` | `canceled`, `canceling`, `deleted`, `deleting`, `failed`, `pending`, `running`, `starting`, `stopped`, `stopping` |

## optimus-ide-collabsdk.TaskLogEntry

```json
{
  "content": "string",
  "id": 0,
  "time": "2019-08-24T14:15:22Z",
  "type": "input"
}
```

### Properties

| Name      | Type                                         | Required | Restrictions | Description |
|-----------|----------------------------------------------|----------|--------------|-------------|
| `content` | string                                       | false    |              |             |
| `id`      | integer                                      | false    |              |             |
| `time`    | string                                       | false    |              |             |
| `type`    | [optimus-ide-collabsdk.TaskLogType](#optimus-ide-collabsdktasklogtype) | false    |              |             |

## optimus-ide-collabsdk.TaskLogType

```json
"input"
```

### Properties

#### Enumerated Values

| Value(s)          |
|-------------------|
| `input`, `output` |

## optimus-ide-collabsdk.TaskLogsResponse

```json
{
  "logs": [
    {
      "content": "string",
      "id": 0,
      "time": "2019-08-24T14:15:22Z",
      "type": "input"
    }
  ],
  "snapshot": true,
  "snapshot_at": "string"
}
```

### Properties

| Name          | Type                                                    | Required | Restrictions | Description |
|---------------|---------------------------------------------------------|----------|--------------|-------------|
| `logs`        | array of [optimus-ide-collabsdk.TaskLogEntry](#optimus-ide-collabsdktasklogentry) | false    |              |             |
| `snapshot`    | boolean                                                 | false    |              |             |
| `snapshot_at` | string                                                  | false    |              |             |

## optimus-ide-collabsdk.TaskSendRequest

```json
{
  "input": "string"
}
```

### Properties

| Name    | Type   | Required | Restrictions | Description |
|---------|--------|----------|--------------|-------------|
| `input` | string | false    |              |             |

## optimus-ide-collabsdk.TaskState

```json
"working"
```

### Properties

#### Enumerated Values

| Value(s)                                |
|-----------------------------------------|
| `complete`, `failed`, `idle`, `working` |

## optimus-ide-collabsdk.TaskStateEntry

```json
{
  "message": "string",
  "state": "working",
  "timestamp": "2019-08-24T14:15:22Z",
  "uri": "string"
}
```

### Properties

| Name        | Type                                     | Required | Restrictions | Description |
|-------------|------------------------------------------|----------|--------------|-------------|
| `message`   | string                                   | false    |              |             |
| `state`     | [optimus-ide-collabsdk.TaskState](#optimus-ide-collabsdktaskstate) | false    |              |             |
| `timestamp` | string                                   | false    |              |             |
| `uri`       | string                                   | false    |              |             |

## optimus-ide-collabsdk.TaskStatus

```json
"pending"
```

### Properties

#### Enumerated Values

| Value(s)                                                          |
|-------------------------------------------------------------------|
| `active`, `error`, `initializing`, `paused`, `pending`, `unknown` |

## optimus-ide-collabsdk.TasksListResponse

```json
{
  "count": 0,
  "tasks": [
    {
      "created_at": "2019-08-24T14:15:22Z",
      "current_state": {
        "message": "string",
        "state": "working",
        "timestamp": "2019-08-24T14:15:22Z",
        "uri": "string"
      },
      "display_name": "string",
      "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
      "initial_prompt": "string",
      "name": "string",
      "organization_id": "7c60d51f-b44e-4682-87d6-449835ea4de6",
      "owner_avatar_url": "string",
      "owner_id": "8826ee2e-7933-4665-aef2-2393f84a0d05",
      "owner_name": "string",
      "status": "pending",
      "template_display_name": "string",
      "template_icon": "string",
      "template_id": "c6d67e98-83ea-49f0-8812-e4abae2b68bc",
      "template_name": "string",
      "template_version_id": "0ba39c92-1f1b-4c32-aa3e-9925d7713eb1",
      "updated_at": "2019-08-24T14:15:22Z",
      "workspace_agent_health": {
        "healthy": false,
        "reason": "agent has lost connection"
      },
      "workspace_agent_id": {
        "uuid": "string",
        "valid": true
      },
      "workspace_agent_lifecycle": "created",
      "workspace_app_id": {
        "uuid": "string",
        "valid": true
      },
      "workspace_build_number": 0,
      "workspace_id": {
        "uuid": "string",
        "valid": true
      },
      "workspace_name": "string",
      "workspace_status": "pending"
    }
  ]
}
```

### Properties

| Name    | Type                                    | Required | Restrictions | Description |
|---------|-----------------------------------------|----------|--------------|-------------|
| `count` | integer                                 | false    |              |             |
| `tasks` | array of [optimus-ide-collabsdk.Task](#optimus-ide-collabsdktask) | false    |              |             |

## optimus-ide-collabsdk.TelemetryConfig

```json
{
  "enable": true,
  "trace": true,
  "url": {
    "forceQuery": true,
    "fragment": "string",
    "host": "string",
    "omitHost": true,
    "opaque": "string",
    "path": "string",
    "rawFragment": "string",
    "rawPath": "string",
    "rawQuery": "string",
    "scheme": "string",
    "user": {}
  }
}
```

### Properties

| Name     | Type                       | Required | Restrictions | Description |
|----------|----------------------------|----------|--------------|-------------|
| `enable` | boolean                    | false    |              |             |
| `trace`  | boolean                    | false    |              |             |
| `url`    | [serpent.URL](#serpenturl) | false    |              |             |

## optimus-ide-collabsdk.Template

```json
{
  "active_user_count": 0,
  "active_version_id": "eae64611-bd53-4a80-bb77-df1e432c0fbc",
  "activity_bump_ms": 0,
  "allow_user_autostart": true,
  "allow_user_autostop": true,
  "allow_user_cancel_workspace_jobs": true,
  "autostart_requirement": {
    "days_of_week": [
      "monday"
    ]
  },
  "autostop_requirement": {
    "days_of_week": [
      "monday"
    ],
    "weeks": 0
  },
  "build_time_stats": {
    "property1": {
      "p50": 123,
      "p95": 146
    },
    "property2": {
      "p50": 123,
      "p95": 146
    }
  },
  "cors_behavior": "simple",
  "created_at": "2019-08-24T14:15:22Z",
  "created_by_id": "9377d689-01fb-4abf-8450-3368d2c1924f",
  "created_by_name": "string",
  "default_ttl_ms": 0,
  "deleted": true,
  "deprecated": true,
  "deprecation_message": "string",
  "description": "string",
  "disable_module_cache": true,
  "display_name": "string",
  "failure_ttl_ms": 0,
  "icon": "string",
  "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
  "max_port_share_level": "owner",
  "name": "string",
  "organization_display_name": "string",
  "organization_icon": "string",
  "organization_id": "7c60d51f-b44e-4682-87d6-449835ea4de6",
  "organization_name": "string",
  "provisioner": "terraform",
  "require_active_version": true,
  "time_til_autostop_notify_ms": 0,
  "time_til_dormant_autodelete_ms": 0,
  "time_til_dormant_ms": 0,
  "updated_at": "2019-08-24T14:15:22Z",
  "use_classic_parameter_flow": true
}
```

### Properties

| Name                               | Type                                                                           | Required | Restrictions | Description                                                                                                                                                                                     |
|------------------------------------|--------------------------------------------------------------------------------|----------|--------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `active_user_count`                | integer                                                                        | false    |              | Active user count is set to -1 when loading.                                                                                                                                                    |
| `active_version_id`                | string                                                                         | false    |              |                                                                                                                                                                                                 |
| `activity_bump_ms`                 | integer                                                                        | false    |              |                                                                                                                                                                                                 |
| `allow_user_autostart`             | boolean                                                                        | false    |              | Allow user autostart and AllowUserAutostop are enterprise-only. Their values are only used if your license is entitled to use the advanced template scheduling feature.                         |
| `allow_user_autostop`              | boolean                                                                        | false    |              |                                                                                                                                                                                                 |
| `allow_user_cancel_workspace_jobs` | boolean                                                                        | false    |              |                                                                                                                                                                                                 |
| `autostart_requirement`            | [optimus-ide-collabsdk.TemplateAutostartRequirement](#optimus-ide-collabsdktemplateautostartrequirement) | false    |              |                                                                                                                                                                                                 |
| `autostop_requirement`             | [optimus-ide-collabsdk.TemplateAutostopRequirement](#optimus-ide-collabsdktemplateautostoprequirement)   | false    |              | Autostop requirement and AutostartRequirement are enterprise features. Its value is only used if your license is entitled to use the advanced template scheduling feature.                      |
| `build_time_stats`                 | [optimus-ide-collabsdk.TemplateBuildTimeStats](#optimus-ide-collabsdktemplatebuildtimestats)             | false    |              |                                                                                                                                                                                                 |
| `cors_behavior`                    | [optimus-ide-collabsdk.CORSBehavior](#optimus-ide-collabsdkcorsbehavior)                                 | false    |              |                                                                                                                                                                                                 |
| `created_at`                       | string                                                                         | false    |              |                                                                                                                                                                                                 |
| `created_by_id`                    | string                                                                         | false    |              |                                                                                                                                                                                                 |
| `created_by_name`                  | string                                                                         | false    |              |                                                                                                                                                                                                 |
| `default_ttl_ms`                   | integer                                                                        | false    |              |                                                                                                                                                                                                 |
| `deleted`                          | boolean                                                                        | false    |              |                                                                                                                                                                                                 |
| `deprecated`                       | boolean                                                                        | false    |              |                                                                                                                                                                                                 |
| `deprecation_message`              | string                                                                         | false    |              |                                                                                                                                                                                                 |
| `description`                      | string                                                                         | false    |              |                                                                                                                                                                                                 |
| `disable_module_cache`             | boolean                                                                        | false    |              | Disable module cache disables the use of cached Terraform modules during provisioning.                                                                                                          |
| `display_name`                     | string                                                                         | false    |              |                                                                                                                                                                                                 |
| `failure_ttl_ms`                   | integer                                                                        | false    |              | Failure ttl ms TimeTilDormantMillis, and TimeTilDormantAutoDeleteMillis are enterprise-only. Their values are used if your license is entitled to use the advanced template scheduling feature. |
| `icon`                             | string                                                                         | false    |              |                                                                                                                                                                                                 |
| `id`                               | string                                                                         | false    |              |                                                                                                                                                                                                 |
| `max_port_share_level`             | [optimus-ide-collabsdk.WorkspaceAgentPortShareLevel](#optimus-ide-collabsdkworkspaceagentportsharelevel) | false    |              |                                                                                                                                                                                                 |
| `name`                             | string                                                                         | false    |              |                                                                                                                                                                                                 |
| `organization_display_name`        | string                                                                         | false    |              |                                                                                                                                                                                                 |
| `organization_icon`                | string                                                                         | false    |              |                                                                                                                                                                                                 |
| `organization_id`                  | string                                                                         | false    |              |                                                                                                                                                                                                 |
| `organization_name`                | string                                                                         | false    |              |                                                                                                                                                                                                 |
| `provisioner`                      | string                                                                         | false    |              |                                                                                                                                                                                                 |
| `require_active_version`           | boolean                                                                        | false    |              | Require active version mandates that workspaces are built with the active template version.                                                                                                     |
| `time_til_autostop_notify_ms`      | integer                                                                        | false    |              | Time til autostop notify ms is the duration before the workspace's autostop deadline at which a reminder notification is sent. 0 disables the notification.                                     |
| `time_til_dormant_autodelete_ms`   | integer                                                                        | false    |              |                                                                                                                                                                                                 |
| `time_til_dormant_ms`              | integer                                                                        | false    |              |                                                                                                                                                                                                 |
| `updated_at`                       | string                                                                         | false    |              |                                                                                                                                                                                                 |
| `use_classic_parameter_flow`       | boolean                                                                        | false    |              |                                                                                                                                                                                                 |

#### Enumerated Values

| Property      | Value(s)    |
|---------------|-------------|
| `provisioner` | `terraform` |

## optimus-ide-collabsdk.TemplateACL

```json
{
  "group": [
    {
      "avatar_url": "http://example.com",
      "display_name": "string",
      "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
      "members": [
        {
          "avatar_url": "http://example.com",
          "created_at": "2019-08-24T14:15:22Z",
          "email": "user@example.com",
          "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
          "is_service_account": true,
          "last_seen_at": "2019-08-24T14:15:22Z",
          "login_type": "",
          "name": "string",
          "status": "active",
          "theme_preference": "string",
          "updated_at": "2019-08-24T14:15:22Z",
          "username": "string"
        }
      ],
      "name": "string",
      "organization_display_name": "string",
      "organization_id": "7c60d51f-b44e-4682-87d6-449835ea4de6",
      "organization_name": "string",
      "quota_allowance": 0,
      "role": "admin",
      "source": "user",
      "total_member_count": 0
    }
  ],
  "users": [
    {
      "avatar_url": "http://example.com",
      "created_at": "2019-08-24T14:15:22Z",
      "email": "user@example.com",
      "has_ai_seat": true,
      "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
      "is_service_account": true,
      "last_seen_at": "2019-08-24T14:15:22Z",
      "login_type": "",
      "name": "string",
      "organization_ids": [
        "497f6eca-6276-4993-bfeb-53cbbbba6f08"
      ],
      "role": "admin",
      "roles": [
        {
          "display_name": "string",
          "name": "string",
          "organization_id": "string"
        }
      ],
      "status": "active",
      "theme_preference": "string",
      "updated_at": "2019-08-24T14:15:22Z",
      "username": "string"
    }
  ]
}
```

### Properties

| Name    | Type                                                      | Required | Restrictions | Description |
|---------|-----------------------------------------------------------|----------|--------------|-------------|
| `group` | array of [optimus-ide-collabsdk.TemplateGroup](#optimus-ide-collabsdktemplategroup) | false    |              |             |
| `users` | array of [optimus-ide-collabsdk.TemplateUser](#optimus-ide-collabsdktemplateuser)   | false    |              |             |

## optimus-ide-collabsdk.TemplateAppUsage

```json
{
  "display_name": "Visual Studio Code",
  "icon": "string",
  "seconds": 80500,
  "slug": "vscode",
  "template_ids": [
    "497f6eca-6276-4993-bfeb-53cbbbba6f08"
  ],
  "times_used": 2,
  "type": "builtin"
}
```

### Properties

| Name           | Type                                                   | Required | Restrictions | Description |
|----------------|--------------------------------------------------------|----------|--------------|-------------|
| `display_name` | string                                                 | false    |              |             |
| `icon`         | string                                                 | false    |              |             |
| `seconds`      | integer                                                | false    |              |             |
| `slug`         | string                                                 | false    |              |             |
| `template_ids` | array of string                                        | false    |              |             |
| `times_used`   | integer                                                | false    |              |             |
| `type`         | [optimus-ide-collabsdk.TemplateAppsType](#optimus-ide-collabsdktemplateappstype) | false    |              |             |

## optimus-ide-collabsdk.TemplateAppsType

```json
"builtin"
```

### Properties

#### Enumerated Values

| Value(s)         |
|------------------|
| `app`, `builtin` |

## optimus-ide-collabsdk.TemplateAutostartRequirement

```json
{
  "days_of_week": [
    "monday"
  ]
}
```

### Properties

| Name           | Type            | Required | Restrictions | Description                                                                                                                             |
|----------------|-----------------|----------|--------------|-----------------------------------------------------------------------------------------------------------------------------------------|
| `days_of_week` | array of string | false    |              | Days of week is a list of days of the week in which autostart is allowed to happen. If no days are specified, autostart is not allowed. |

## optimus-ide-collabsdk.TemplateAutostopRequirement

```json
{
  "days_of_week": [
    "monday"
  ],
  "weeks": 0
}
```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|`days_of_week`|array of string|false||Days of week is a list of days of the week on which restarts are required. Restarts happen within the user's quiet hours (in their configured timezone). If no days are specified, restarts are not required. Weekdays cannot be specified twice.
Restarts will only happen on weekdays in this list on weeks which line up with Weeks.|
|`weeks`|integer|false||Weeks is the number of weeks between required restarts. Weeks are synced across all workspaces (and Optimus-IDE-Collab deployments) using modulo math on a hardcoded epoch week of January 2nd, 2023 (the first Monday of 2023). Values of 0 or 1 indicate weekly restarts. Values of 2 indicate fortnightly restarts, etc.|

## optimus-ide-collabsdk.TemplateBuildTimeStats

```json
{
  "property1": {
    "p50": 123,
    "p95": 146
  },
  "property2": {
    "p50": 123,
    "p95": 146
  }
}
```

### Properties

| Name             | Type                                                 | Required | Restrictions | Description |
|------------------|------------------------------------------------------|----------|--------------|-------------|
| `[any property]` | [optimus-ide-collabsdk.TransitionStats](#optimus-ide-collabsdktransitionstats) | false    |              |             |

## optimus-ide-collabsdk.TemplateBuilderBase

```json
{
  "description": "string",
  "icon": "string",
  "id": "string",
  "name": "string",
  "os": "string",
  "prerequisites": "string",
  "variables": [
    {
      "default": [
        0
      ],
      "description": "string",
      "name": "string",
      "required": true,
      "sensitive": true,
      "type": "string"
    }
  ]
}
```

### Properties

| Name            | Type                                                                                      | Required | Restrictions | Description |
|-----------------|-------------------------------------------------------------------------------------------|----------|--------------|-------------|
| `description`   | string                                                                                    | false    |              |             |
| `icon`          | string                                                                                    | false    |              |             |
| `id`            | string                                                                                    | false    |              |             |
| `name`          | string                                                                                    | false    |              |             |
| `os`            | string                                                                                    | false    |              |             |
| `prerequisites` | string                                                                                    | false    |              |             |
| `variables`     | array of [optimus-ide-collabsdk.TemplateBuilderModuleVariable](#optimus-ide-collabsdktemplatebuildermodulevariable) | false    |              |             |

## optimus-ide-collabsdk.TemplateBuilderBasesResponse

```json
{
  "bases": [
    {
      "description": "string",
      "icon": "string",
      "id": "string",
      "name": "string",
      "os": "string",
      "prerequisites": "string",
      "variables": [
        {
          "default": [
            0
          ],
          "description": "string",
          "name": "string",
          "required": true,
          "sensitive": true,
          "type": "string"
        }
      ]
    }
  ]
}
```

### Properties

| Name    | Type                                                                  | Required | Restrictions | Description |
|---------|-----------------------------------------------------------------------|----------|--------------|-------------|
| `bases` | array of [optimus-ide-collabsdk.TemplateBuilderBase](#optimus-ide-collabsdktemplatebuilderbase) | false    |              |             |

## optimus-ide-collabsdk.TemplateBuilderComposeModule

```json
{
  "id": "string",
  "variables": {
    "property1": "string",
    "property2": "string"
  }
}
```

### Properties

| Name               | Type   | Required | Restrictions | Description |
|--------------------|--------|----------|--------------|-------------|
| `id`               | string | false    |              |             |
| `variables`        | object | false    |              |             |
| » `[any property]` | string | false    |              |             |

## optimus-ide-collabsdk.TemplateBuilderComposeRequest

```json
{
  "base_template_id": "string",
  "base_variable_values": {
    "property1": "string",
    "property2": "string"
  },
  "modules": [
    {
      "id": "string",
      "variables": {
        "property1": "string",
        "property2": "string"
      }
    }
  ]
}
```

### Properties

| Name                   | Type                                                                                    | Required | Restrictions | Description |
|------------------------|-----------------------------------------------------------------------------------------|----------|--------------|-------------|
| `base_template_id`     | string                                                                                  | false    |              |             |
| `base_variable_values` | object                                                                                  | false    |              |             |
| » `[any property]`     | string                                                                                  | false    |              |             |
| `modules`              | array of [optimus-ide-collabsdk.TemplateBuilderComposeModule](#optimus-ide-collabsdktemplatebuildercomposemodule) | false    |              |             |

## optimus-ide-collabsdk.TemplateBuilderConfig

```json
{
  "disabled": true,
  "registry_url": "string"
}
```

### Properties

| Name           | Type    | Required | Restrictions | Description |
|----------------|---------|----------|--------------|-------------|
| `disabled`     | boolean | false    |              |             |
| `registry_url` | string  | false    |              |             |

## optimus-ide-collabsdk.TemplateBuilderCreateTemplateRequest

```json
{
  "base_template_id": "string",
  "base_variable_values": {
    "property1": "string",
    "property2": "string"
  },
  "description": "string",
  "display_name": "string",
  "icon": "string",
  "modules": [
    {
      "id": "string",
      "variables": {
        "property1": "string",
        "property2": "string"
      }
    }
  ],
  "name": "string",
  "organization_id": "7c60d51f-b44e-4682-87d6-449835ea4de6",
  "provisioner_tags": {
    "property1": "string",
    "property2": "string"
  }
}
```

### Properties

| Name                   | Type                                                                                    | Required | Restrictions | Description |
|------------------------|-----------------------------------------------------------------------------------------|----------|--------------|-------------|
| `base_template_id`     | string                                                                                  | false    |              |             |
| `base_variable_values` | object                                                                                  | false    |              |             |
| » `[any property]`     | string                                                                                  | false    |              |             |
| `description`          | string                                                                                  | false    |              |             |
| `display_name`         | string                                                                                  | false    |              |             |
| `icon`                 | string                                                                                  | false    |              |             |
| `modules`              | array of [optimus-ide-collabsdk.TemplateBuilderComposeModule](#optimus-ide-collabsdktemplatebuildercomposemodule) | false    |              |             |
| `name`                 | string                                                                                  | true     |              |             |
| `organization_id`      | string                                                                                  | true     |              |             |
| `provisioner_tags`     | object                                                                                  | false    |              |             |
| » `[any property]`     | string                                                                                  | false    |              |             |

## optimus-ide-collabsdk.TemplateBuilderCreateTemplateResponse

```json
{
  "template": {
    "active_user_count": 0,
    "active_version_id": "eae64611-bd53-4a80-bb77-df1e432c0fbc",
    "activity_bump_ms": 0,
    "allow_user_autostart": true,
    "allow_user_autostop": true,
    "allow_user_cancel_workspace_jobs": true,
    "autostart_requirement": {
      "days_of_week": [
        "monday"
      ]
    },
    "autostop_requirement": {
      "days_of_week": [
        "monday"
      ],
      "weeks": 0
    },
    "build_time_stats": {
      "property1": {
        "p50": 123,
        "p95": 146
      },
      "property2": {
        "p50": 123,
        "p95": 146
      }
    },
    "cors_behavior": "simple",
    "created_at": "2019-08-24T14:15:22Z",
    "created_by_id": "9377d689-01fb-4abf-8450-3368d2c1924f",
    "created_by_name": "string",
    "default_ttl_ms": 0,
    "deleted": true,
    "deprecated": true,
    "deprecation_message": "string",
    "description": "string",
    "disable_module_cache": true,
    "display_name": "string",
    "failure_ttl_ms": 0,
    "icon": "string",
    "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
    "max_port_share_level": "owner",
    "name": "string",
    "organization_display_name": "string",
    "organization_icon": "string",
    "organization_id": "7c60d51f-b44e-4682-87d6-449835ea4de6",
    "organization_name": "string",
    "provisioner": "terraform",
    "require_active_version": true,
    "time_til_autostop_notify_ms": 0,
    "time_til_dormant_autodelete_ms": 0,
    "time_til_dormant_ms": 0,
    "updated_at": "2019-08-24T14:15:22Z",
    "use_classic_parameter_flow": true
  }
}
```

### Properties

| Name       | Type                                   | Required | Restrictions | Description |
|------------|----------------------------------------|----------|--------------|-------------|
| `template` | [optimus-ide-collabsdk.Template](#optimus-ide-collabsdktemplate) | false    |              |             |

## optimus-ide-collabsdk.TemplateBuilderModule

```json
{
  "category": "string",
  "compatible_os": [
    "string"
  ],
  "conflicts_with": [
    "string"
  ],
  "description": "string",
  "display_name": "string",
  "icon": "string",
  "id": "string",
  "variables": [
    {
      "default": [
        0
      ],
      "description": "string",
      "name": "string",
      "required": true,
      "sensitive": true,
      "type": "string"
    }
  ],
  "version": "string"
}
```

### Properties

| Name             | Type                                                                                      | Required | Restrictions | Description |
|------------------|-------------------------------------------------------------------------------------------|----------|--------------|-------------|
| `category`       | string                                                                                    | false    |              |             |
| `compatible_os`  | array of string                                                                           | false    |              |             |
| `conflicts_with` | array of string                                                                           | false    |              |             |
| `description`    | string                                                                                    | false    |              |             |
| `display_name`   | string                                                                                    | false    |              |             |
| `icon`           | string                                                                                    | false    |              |             |
| `id`             | string                                                                                    | false    |              |             |
| `variables`      | array of [optimus-ide-collabsdk.TemplateBuilderModuleVariable](#optimus-ide-collabsdktemplatebuildermodulevariable) | false    |              |             |
| `version`        | string                                                                                    | false    |              |             |

## optimus-ide-collabsdk.TemplateBuilderModuleVariable

```json
{
  "default": [
    0
  ],
  "description": "string",
  "name": "string",
  "required": true,
  "sensitive": true,
  "type": "string"
}
```

### Properties

| Name          | Type                                                                         | Required | Restrictions | Description |
|---------------|------------------------------------------------------------------------------|----------|--------------|-------------|
| `default`     | array of integer                                                             | false    |              |             |
| `description` | string                                                                       | false    |              |             |
| `name`        | string                                                                       | false    |              |             |
| `required`    | boolean                                                                      | false    |              |             |
| `sensitive`   | boolean                                                                      | false    |              |             |
| `type`        | [optimus-ide-collabsdk.TemplateBuilderVariableType](#optimus-ide-collabsdktemplatebuildervariabletype) | false    |              |             |

## optimus-ide-collabsdk.TemplateBuilderModulesResponse

```json
{
  "modules": [
    {
      "category": "string",
      "compatible_os": [
        "string"
      ],
      "conflicts_with": [
        "string"
      ],
      "description": "string",
      "display_name": "string",
      "icon": "string",
      "id": "string",
      "variables": [
        {
          "default": [
            0
          ],
          "description": "string",
          "name": "string",
          "required": true,
          "sensitive": true,
          "type": "string"
        }
      ],
      "version": "string"
    }
  ]
}
```

### Properties

| Name      | Type                                                                      | Required | Restrictions | Description |
|-----------|---------------------------------------------------------------------------|----------|--------------|-------------|
| `modules` | array of [optimus-ide-collabsdk.TemplateBuilderModule](#optimus-ide-collabsdktemplatebuildermodule) | false    |              |             |

## optimus-ide-collabsdk.TemplateBuilderSessionEventType

```json
"wizard_entry"
```

### Properties

#### Enumerated Values

| Value(s)                             |
|--------------------------------------|
| `compose_completion`, `wizard_entry` |

## optimus-ide-collabsdk.TemplateBuilderSessionRequest

```json
{
  "base_template_id": "string",
  "duration_seconds": 0,
  "event_type": "wizard_entry",
  "module_ids": [
    "string"
  ],
  "session_id": "1ffd059c-17ea-40a8-8aef-70fd0307db82",
  "success": true
}
```

### Properties

| Name               | Type                                                                                 | Required | Restrictions | Description |
|--------------------|--------------------------------------------------------------------------------------|----------|--------------|-------------|
| `base_template_id` | string                                                                               | false    |              |             |
| `duration_seconds` | number                                                                               | false    |              |             |
| `event_type`       | [optimus-ide-collabsdk.TemplateBuilderSessionEventType](#optimus-ide-collabsdktemplatebuildersessioneventtype) | true     |              |             |
| `module_ids`       | array of string                                                                      | false    |              |             |
| `session_id`       | string                                                                               | true     |              |             |
| `success`          | boolean                                                                              | false    |              |             |

#### Enumerated Values

| Property     | Value(s)                             |
|--------------|--------------------------------------|
| `event_type` | `compose_completion`, `wizard_entry` |

## optimus-ide-collabsdk.TemplateBuilderVariableType

```json
"string"
```

### Properties

#### Enumerated Values

| Value(s)                   |
|----------------------------|
| `bool`, `number`, `string` |

## optimus-ide-collabsdk.TemplateExample

```json
{
  "description": "string",
  "icon": "string",
  "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
  "markdown": "string",
  "name": "string",
  "tags": [
    "string"
  ],
  "url": "string"
}
```

### Properties

| Name          | Type            | Required | Restrictions | Description |
|---------------|-----------------|----------|--------------|-------------|
| `description` | string          | false    |              |             |
| `icon`        | string          | false    |              |             |
| `id`          | string          | false    |              |             |
| `markdown`    | string          | false    |              |             |
| `name`        | string          | false    |              |             |
| `tags`        | array of string | false    |              |             |
| `url`         | string          | false    |              |             |

## optimus-ide-collabsdk.TemplateGroup

```json
{
  "avatar_url": "http://example.com",
  "display_name": "string",
  "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
  "members": [
    {
      "avatar_url": "http://example.com",
      "created_at": "2019-08-24T14:15:22Z",
      "email": "user@example.com",
      "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
      "is_service_account": true,
      "last_seen_at": "2019-08-24T14:15:22Z",
      "login_type": "",
      "name": "string",
      "status": "active",
      "theme_preference": "string",
      "updated_at": "2019-08-24T14:15:22Z",
      "username": "string"
    }
  ],
  "name": "string",
  "organization_display_name": "string",
  "organization_id": "7c60d51f-b44e-4682-87d6-449835ea4de6",
  "organization_name": "string",
  "quota_allowance": 0,
  "role": "admin",
  "source": "user",
  "total_member_count": 0
}
```

### Properties

| Name                        | Type                                                  | Required | Restrictions | Description                                                                                                                                                           |
|-----------------------------|-------------------------------------------------------|----------|--------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `avatar_url`                | string                                                | false    |              |                                                                                                                                                                       |
| `display_name`              | string                                                | false    |              |                                                                                                                                                                       |
| `id`                        | string                                                | false    |              |                                                                                                                                                                       |
| `members`                   | array of [optimus-ide-collabsdk.ReducedUser](#optimus-ide-collabsdkreduceduser) | false    |              |                                                                                                                                                                       |
| `name`                      | string                                                | false    |              |                                                                                                                                                                       |
| `organization_display_name` | string                                                | false    |              |                                                                                                                                                                       |
| `organization_id`           | string                                                | false    |              |                                                                                                                                                                       |
| `organization_name`         | string                                                | false    |              |                                                                                                                                                                       |
| `quota_allowance`           | integer                                               | false    |              |                                                                                                                                                                       |
| `role`                      | [optimus-ide-collabsdk.TemplateRole](#optimus-ide-collabsdktemplaterole)        | false    |              |                                                                                                                                                                       |
| `source`                    | [optimus-ide-collabsdk.GroupSource](#optimus-ide-collabsdkgroupsource)          | false    |              |                                                                                                                                                                       |
| `total_member_count`        | integer                                               | false    |              | How many members are in this group. Shows the total count, even if the user is not authorized to read group member details. May be greater than `len(Group.Members)`. |

#### Enumerated Values

| Property | Value(s)       |
|----------|----------------|
| `role`   | `admin`, `use` |

## optimus-ide-collabsdk.TemplateInsightsIntervalReport

```json
{
  "active_users": 14,
  "end_time": "2019-08-24T14:15:22Z",
  "interval": "week",
  "start_time": "2019-08-24T14:15:22Z",
  "template_ids": [
    "497f6eca-6276-4993-bfeb-53cbbbba6f08"
  ]
}
```

### Properties

| Name           | Type                                                               | Required | Restrictions | Description |
|----------------|--------------------------------------------------------------------|----------|--------------|-------------|
| `active_users` | integer                                                            | false    |              |             |
| `end_time`     | string                                                             | false    |              |             |
| `interval`     | [optimus-ide-collabsdk.InsightsReportInterval](#optimus-ide-collabsdkinsightsreportinterval) | false    |              |             |
| `start_time`   | string                                                             | false    |              |             |
| `template_ids` | array of string                                                    | false    |              |             |

## optimus-ide-collabsdk.TemplateInsightsReport

```json
{
  "active_users": 22,
  "apps_usage": [
    {
      "display_name": "Visual Studio Code",
      "icon": "string",
      "seconds": 80500,
      "slug": "vscode",
      "template_ids": [
        "497f6eca-6276-4993-bfeb-53cbbbba6f08"
      ],
      "times_used": 2,
      "type": "builtin"
    }
  ],
  "end_time": "2019-08-24T14:15:22Z",
  "parameters_usage": [
    {
      "description": "string",
      "display_name": "string",
      "name": "string",
      "options": [
        {
          "description": "string",
          "icon": "string",
          "name": "string",
          "value": "string"
        }
      ],
      "template_ids": [
        "497f6eca-6276-4993-bfeb-53cbbbba6f08"
      ],
      "type": "string",
      "values": [
        {
          "count": 0,
          "value": "string"
        }
      ]
    }
  ],
  "start_time": "2019-08-24T14:15:22Z",
  "template_ids": [
    "497f6eca-6276-4993-bfeb-53cbbbba6f08"
  ]
}
```

### Properties

| Name               | Type                                                                        | Required | Restrictions | Description |
|--------------------|-----------------------------------------------------------------------------|----------|--------------|-------------|
| `active_users`     | integer                                                                     | false    |              |             |
| `apps_usage`       | array of [optimus-ide-collabsdk.TemplateAppUsage](#optimus-ide-collabsdktemplateappusage)             | false    |              |             |
| `end_time`         | string                                                                      | false    |              |             |
| `parameters_usage` | array of [optimus-ide-collabsdk.TemplateParameterUsage](#optimus-ide-collabsdktemplateparameterusage) | false    |              |             |
| `start_time`       | string                                                                      | false    |              |             |
| `template_ids`     | array of string                                                             | false    |              |             |

## optimus-ide-collabsdk.TemplateInsightsResponse

```json
{
  "interval_reports": [
    {
      "active_users": 14,
      "end_time": "2019-08-24T14:15:22Z",
      "interval": "week",
      "start_time": "2019-08-24T14:15:22Z",
      "template_ids": [
        "497f6eca-6276-4993-bfeb-53cbbbba6f08"
      ]
    }
  ],
  "report": {
    "active_users": 22,
    "apps_usage": [
      {
        "display_name": "Visual Studio Code",
        "icon": "string",
        "seconds": 80500,
        "slug": "vscode",
        "template_ids": [
          "497f6eca-6276-4993-bfeb-53cbbbba6f08"
        ],
        "times_used": 2,
        "type": "builtin"
      }
    ],
    "end_time": "2019-08-24T14:15:22Z",
    "parameters_usage": [
      {
        "description": "string",
        "display_name": "string",
        "name": "string",
        "options": [
          {
            "description": "string",
            "icon": "string",
            "name": "string",
            "value": "string"
          }
        ],
        "template_ids": [
          "497f6eca-6276-4993-bfeb-53cbbbba6f08"
        ],
        "type": "string",
        "values": [
          {
            "count": 0,
            "value": "string"
          }
        ]
      }
    ],
    "start_time": "2019-08-24T14:15:22Z",
    "template_ids": [
      "497f6eca-6276-4993-bfeb-53cbbbba6f08"
    ]
  }
}
```

### Properties

| Name               | Type                                                                                        | Required | Restrictions | Description |
|--------------------|---------------------------------------------------------------------------------------------|----------|--------------|-------------|
| `interval_reports` | array of [optimus-ide-collabsdk.TemplateInsightsIntervalReport](#optimus-ide-collabsdktemplateinsightsintervalreport) | false    |              |             |
| `report`           | [optimus-ide-collabsdk.TemplateInsightsReport](#optimus-ide-collabsdktemplateinsightsreport)                          | false    |              |             |

## optimus-ide-collabsdk.TemplateParameterUsage

```json
{
  "description": "string",
  "display_name": "string",
  "name": "string",
  "options": [
    {
      "description": "string",
      "icon": "string",
      "name": "string",
      "value": "string"
    }
  ],
  "template_ids": [
    "497f6eca-6276-4993-bfeb-53cbbbba6f08"
  ],
  "type": "string",
  "values": [
    {
      "count": 0,
      "value": "string"
    }
  ]
}
```

### Properties

| Name           | Type                                                                                        | Required | Restrictions | Description |
|----------------|---------------------------------------------------------------------------------------------|----------|--------------|-------------|
| `description`  | string                                                                                      | false    |              |             |
| `display_name` | string                                                                                      | false    |              |             |
| `name`         | string                                                                                      | false    |              |             |
| `options`      | array of [optimus-ide-collabsdk.TemplateVersionParameterOption](#optimus-ide-collabsdktemplateversionparameteroption) | false    |              |             |
| `template_ids` | array of string                                                                             | false    |              |             |
| `type`         | string                                                                                      | false    |              |             |
| `values`       | array of [optimus-ide-collabsdk.TemplateParameterValue](#optimus-ide-collabsdktemplateparametervalue)                 | false    |              |             |

## optimus-ide-collabsdk.TemplateParameterValue

```json
{
  "count": 0,
  "value": "string"
}
```

### Properties

| Name    | Type    | Required | Restrictions | Description |
|---------|---------|----------|--------------|-------------|
| `count` | integer | false    |              |             |
| `value` | string  | false    |              |             |

## optimus-ide-collabsdk.TemplateRole

```json
"admin"
```

### Properties

#### Enumerated Values

| Value(s)           |
|--------------------|
| ``, `admin`, `use` |

## optimus-ide-collabsdk.TemplateUser

```json
{
  "avatar_url": "http://example.com",
  "created_at": "2019-08-24T14:15:22Z",
  "email": "user@example.com",
  "has_ai_seat": true,
  "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
  "is_service_account": true,
  "last_seen_at": "2019-08-24T14:15:22Z",
  "login_type": "",
  "name": "string",
  "organization_ids": [
    "497f6eca-6276-4993-bfeb-53cbbbba6f08"
  ],
  "role": "admin",
  "roles": [
    {
      "display_name": "string",
      "name": "string",
      "organization_id": "string"
    }
  ],
  "status": "active",
  "theme_preference": "string",
  "updated_at": "2019-08-24T14:15:22Z",
  "username": "string"
}
```

### Properties

| Name                 | Type                                            | Required | Restrictions | Description                                                                                      |
|----------------------|-------------------------------------------------|----------|--------------|--------------------------------------------------------------------------------------------------|
| `avatar_url`         | string                                          | false    |              |                                                                                                  |
| `created_at`         | string                                          | true     |              |                                                                                                  |
| `email`              | string                                          | true     |              |                                                                                                  |
| `has_ai_seat`        | boolean                                         | false    |              | Has ai seat intentionally omits omitempty so the API always includes the field, even when false. |
| `id`                 | string                                          | true     |              |                                                                                                  |
| `is_service_account` | boolean                                         | false    |              |                                                                                                  |
| `last_seen_at`       | string                                          | false    |              |                                                                                                  |
| `login_type`         | [optimus-ide-collabsdk.LoginType](#optimus-ide-collabsdklogintype)        | false    |              |                                                                                                  |
| `name`               | string                                          | false    |              |                                                                                                  |
| `organization_ids`   | array of string                                 | false    |              |                                                                                                  |
| `role`               | [optimus-ide-collabsdk.TemplateRole](#optimus-ide-collabsdktemplaterole)  | false    |              |                                                                                                  |
| `roles`              | array of [optimus-ide-collabsdk.SlimRole](#optimus-ide-collabsdkslimrole) | false    |              |                                                                                                  |
| `status`             | [optimus-ide-collabsdk.UserStatus](#optimus-ide-collabsdkuserstatus)      | false    |              |                                                                                                  |
| `theme_preference`   | string                                          | false    |              | Deprecated: this value should be retrieved from `optimus-ide-collabsdk.UserPreferenceSettings` instead.       |
| `updated_at`         | string                                          | false    |              |                                                                                                  |
| `username`           | string                                          | true     |              |                                                                                                  |

#### Enumerated Values

| Property | Value(s)              |
|----------|-----------------------|
| `role`   | `admin`, `use`        |
| `status` | `active`, `suspended` |

## optimus-ide-collabsdk.TemplateVersion

```json
{
  "archived": true,
  "created_at": "2019-08-24T14:15:22Z",
  "created_by": {
    "avatar_url": "http://example.com",
    "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
    "name": "string",
    "username": "string"
  },
  "has_external_agent": true,
  "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
  "job": {
    "available_workers": [
      "497f6eca-6276-4993-bfeb-53cbbbba6f08"
    ],
    "canceled_at": "2019-08-24T14:15:22Z",
    "completed_at": "2019-08-24T14:15:22Z",
    "created_at": "2019-08-24T14:15:22Z",
    "error": "string",
    "error_code": "REQUIRED_TEMPLATE_VARIABLES",
    "file_id": "8a0cfb4f-ddc9-436d-91bb-75133c583767",
    "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
    "initiator_id": "06588898-9a84-4b35-ba8f-f9cbd64946f3",
    "input": {
      "error": "string",
      "template_version_id": "0ba39c92-1f1b-4c32-aa3e-9925d7713eb1",
      "workspace_build_id": "badaf2eb-96c5-4050-9f1d-db2d39ca5478"
    },
    "logs_overflowed": true,
    "metadata": {
      "template_display_name": "string",
      "template_icon": "string",
      "template_id": "c6d67e98-83ea-49f0-8812-e4abae2b68bc",
      "template_name": "string",
      "template_version_name": "string",
      "workspace_build_transition": "start",
      "workspace_id": "0967198e-ec7b-4c6b-b4d3-f71244cadbe9",
      "workspace_name": "string"
    },
    "organization_id": "7c60d51f-b44e-4682-87d6-449835ea4de6",
    "queue_position": 0,
    "queue_size": 0,
    "started_at": "2019-08-24T14:15:22Z",
    "status": "pending",
    "tags": {
      "property1": "string",
      "property2": "string"
    },
    "type": "template_version_import",
    "worker_id": "ae5fa6f7-c55b-40c1-b40a-b36ac467652b",
    "worker_name": "string"
  },
  "matched_provisioners": {
    "available": 0,
    "count": 0,
    "most_recently_seen": "2019-08-24T14:15:22Z"
  },
  "message": "string",
  "name": "string",
  "organization_id": "7c60d51f-b44e-4682-87d6-449835ea4de6",
  "readme": "string",
  "template_id": "c6d67e98-83ea-49f0-8812-e4abae2b68bc",
  "updated_at": "2019-08-24T14:15:22Z",
  "warnings": [
    "UNSUPPORTED_WORKSPACES"
  ]
}
```

### Properties

| Name                   | Type                                                                        | Required | Restrictions | Description |
|------------------------|-----------------------------------------------------------------------------|----------|--------------|-------------|
| `archived`             | boolean                                                                     | false    |              |             |
| `created_at`           | string                                                                      | false    |              |             |
| `created_by`           | [optimus-ide-collabsdk.MinimalUser](#optimus-ide-collabsdkminimaluser)                                | false    |              |             |
| `has_external_agent`   | boolean                                                                     | false    |              |             |
| `id`                   | string                                                                      | false    |              |             |
| `job`                  | [optimus-ide-collabsdk.ProvisionerJob](#optimus-ide-collabsdkprovisionerjob)                          | false    |              |             |
| `matched_provisioners` | [optimus-ide-collabsdk.MatchedProvisioners](#optimus-ide-collabsdkmatchedprovisioners)                | false    |              |             |
| `message`              | string                                                                      | false    |              |             |
| `name`                 | string                                                                      | false    |              |             |
| `organization_id`      | string                                                                      | false    |              |             |
| `readme`               | string                                                                      | false    |              |             |
| `template_id`          | string                                                                      | false    |              |             |
| `updated_at`           | string                                                                      | false    |              |             |
| `warnings`             | array of [optimus-ide-collabsdk.TemplateVersionWarning](#optimus-ide-collabsdktemplateversionwarning) | false    |              |             |

## optimus-ide-collabsdk.TemplateVersionExternalAuth

```json
{
  "authenticate_url": "string",
  "authenticated": true,
  "display_icon": "string",
  "display_name": "string",
  "id": "string",
  "optional": true,
  "type": "string"
}
```

### Properties

| Name               | Type    | Required | Restrictions | Description |
|--------------------|---------|----------|--------------|-------------|
| `authenticate_url` | string  | false    |              |             |
| `authenticated`    | boolean | false    |              |             |
| `display_icon`     | string  | false    |              |             |
| `display_name`     | string  | false    |              |             |
| `id`               | string  | false    |              |             |
| `optional`         | boolean | false    |              |             |
| `type`             | string  | false    |              |             |

## optimus-ide-collabsdk.TemplateVersionParameter

```json
{
  "default_value": "string",
  "description": "string",
  "description_plaintext": "string",
  "display_name": "string",
  "ephemeral": true,
  "form_type": "",
  "icon": "string",
  "mutable": true,
  "name": "string",
  "options": [
    {
      "description": "string",
      "icon": "string",
      "name": "string",
      "value": "string"
    }
  ],
  "required": true,
  "type": "string",
  "validation_error": "string",
  "validation_max": 0,
  "validation_min": 0,
  "validation_monotonic": "increasing",
  "validation_regex": "string"
}
```

### Properties

| Name                    | Type                                                                                        | Required | Restrictions | Description                                                                                        |
|-------------------------|---------------------------------------------------------------------------------------------|----------|--------------|----------------------------------------------------------------------------------------------------|
| `default_value`         | string                                                                                      | false    |              |                                                                                                    |
| `description`           | string                                                                                      | false    |              |                                                                                                    |
| `description_plaintext` | string                                                                                      | false    |              |                                                                                                    |
| `display_name`          | string                                                                                      | false    |              |                                                                                                    |
| `ephemeral`             | boolean                                                                                     | false    |              |                                                                                                    |
| `form_type`             | string                                                                                      | false    |              | Form type has an enum value of empty string, `""`. Keep the leading comma in the enums struct tag. |
| `icon`                  | string                                                                                      | false    |              |                                                                                                    |
| `mutable`               | boolean                                                                                     | false    |              |                                                                                                    |
| `name`                  | string                                                                                      | false    |              |                                                                                                    |
| `options`               | array of [optimus-ide-collabsdk.TemplateVersionParameterOption](#optimus-ide-collabsdktemplateversionparameteroption) | false    |              |                                                                                                    |
| `required`              | boolean                                                                                     | false    |              |                                                                                                    |
| `type`                  | string                                                                                      | false    |              |                                                                                                    |
| `validation_error`      | string                                                                                      | false    |              |                                                                                                    |
| `validation_max`        | integer                                                                                     | false    |              |                                                                                                    |
| `validation_min`        | integer                                                                                     | false    |              |                                                                                                    |
| `validation_monotonic`  | [optimus-ide-collabsdk.ValidationMonotonicOrder](#optimus-ide-collabsdkvalidationmonotonicorder)                      | false    |              |                                                                                                    |
| `validation_regex`      | string                                                                                      | false    |              |                                                                                                    |

#### Enumerated Values

| Property               | Value(s)                                                                                                            |
|------------------------|---------------------------------------------------------------------------------------------------------------------|
| `form_type`            | ``, `checkbox`, `dropdown`, `error`, `input`, `multi-select`, `radio`, `slider`, `switch`, `tag-select`, `textarea` |
| `type`                 | `bool`, `list(string)`, `number`, `string`                                                                          |
| `validation_monotonic` | `decreasing`, `increasing`                                                                                          |

## optimus-ide-collabsdk.TemplateVersionParameterOption

```json
{
  "description": "string",
  "icon": "string",
  "name": "string",
  "value": "string"
}
```

### Properties

| Name          | Type   | Required | Restrictions | Description |
|---------------|--------|----------|--------------|-------------|
| `description` | string | false    |              |             |
| `icon`        | string | false    |              |             |
| `name`        | string | false    |              |             |
| `value`       | string | false    |              |             |

## optimus-ide-collabsdk.TemplateVersionVariable

```json
{
  "default_value": "string",
  "description": "string",
  "name": "string",
  "required": true,
  "sensitive": true,
  "type": "string",
  "value": "string"
}
```

### Properties

| Name            | Type    | Required | Restrictions | Description |
|-----------------|---------|----------|--------------|-------------|
| `default_value` | string  | false    |              |             |
| `description`   | string  | false    |              |             |
| `name`          | string  | false    |              |             |
| `required`      | boolean | false    |              |             |
| `sensitive`     | boolean | false    |              |             |
| `type`          | string  | false    |              |             |
| `value`         | string  | false    |              |             |

#### Enumerated Values

| Property | Value(s)                   |
|----------|----------------------------|
| `type`   | `bool`, `number`, `string` |

## optimus-ide-collabsdk.TemplateVersionWarning

```json
"UNSUPPORTED_WORKSPACES"
```

### Properties

#### Enumerated Values

| Value(s)                 |
|--------------------------|
| `UNSUPPORTED_WORKSPACES` |

## optimus-ide-collabsdk.TerminalFontName

```json
""
```

### Properties

#### Enumerated Values

| Value(s)                                                                            |
|-------------------------------------------------------------------------------------|
| ``, `fira-code`, `geist-mono`, `ibm-plex-mono`, `jetbrains-mono`, `source-code-pro` |

## optimus-ide-collabsdk.ThemeMode

```json
""
```

### Properties

#### Enumerated Values

| Value(s)             |
|----------------------|
| ``, `single`, `sync` |

## optimus-ide-collabsdk.ThinkingDisplayMode

```json
"auto"
```

### Properties

#### Enumerated Values

| Value(s)                                                 |
|----------------------------------------------------------|
| `always_collapsed`, `always_expanded`, `auto`, `preview` |

## optimus-ide-collabsdk.TimingStage

```json
"init"
```

### Properties

#### Enumerated Values

| Value(s)                                                             |
|----------------------------------------------------------------------|
| `apply`, `connect`, `cron`, `graph`, `init`, `plan`, `start`, `stop` |

## optimus-ide-collabsdk.TokenConfig

```json
{
  "max_token_lifetime": 0
}
```

### Properties

| Name                 | Type    | Required | Restrictions | Description |
|----------------------|---------|----------|--------------|-------------|
| `max_token_lifetime` | integer | false    |              |             |

## optimus-ide-collabsdk.TraceConfig

```json
{
  "capture_logs": true,
  "data_dog": true,
  "enable": true,
  "honeycomb_api_key": "string"
}
```

### Properties

| Name                | Type    | Required | Restrictions | Description |
|---------------------|---------|----------|--------------|-------------|
| `capture_logs`      | boolean | false    |              |             |
| `data_dog`          | boolean | false    |              |             |
| `enable`            | boolean | false    |              |             |
| `honeycomb_api_key` | string  | false    |              |             |

## optimus-ide-collabsdk.TransitionStats

```json
{
  "p50": 123,
  "p95": 146
}
```

### Properties

| Name  | Type    | Required | Restrictions | Description |
|-------|---------|----------|--------------|-------------|
| `p50` | integer | false    |              |             |
| `p95` | integer | false    |              |             |

## optimus-ide-collabsdk.UpdateAIProviderRequest

```json
{
  "api_keys": [
    {
      "api_key": "string",
      "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08"
    }
  ],
  "base_url": "string",
  "display_name": "string",
  "enabled": true,
  "icon": "string",
  "settings": {}
}
```

### Properties

| Name           | Type                                                                      | Required | Restrictions | Description |
|----------------|---------------------------------------------------------------------------|----------|--------------|-------------|
| `api_keys`     | array of [optimus-ide-collabsdk.AIProviderKeyMutation](#optimus-ide-collabsdkaiproviderkeymutation) | false    |              |             |
| `base_url`     | string                                                                    | false    |              |             |
| `display_name` | string                                                                    | false    |              |             |
| `enabled`      | boolean                                                                   | false    |              |             |
| `icon`         | string                                                                    | false    |              |             |
| `settings`     | [optimus-ide-collabsdk.AIProviderSettings](#optimus-ide-collabsdkaiprovidersettings)                | false    |              |             |

## optimus-ide-collabsdk.UpdateActiveTemplateVersion

```json
{
  "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08"
}
```

### Properties

| Name | Type   | Required | Restrictions | Description |
|------|--------|----------|--------------|-------------|
| `id` | string | true     |              |             |

## optimus-ide-collabsdk.UpdateAppearanceConfig

```json
{
  "announcement_banners": [
    {
      "background_color": "string",
      "enabled": true,
      "message": "string"
    }
  ],
  "application_name": "string",
  "logo_url": "string",
  "service_banner": {
    "background_color": "string",
    "enabled": true,
    "message": "string"
  }
}
```

### Properties

| Name                   | Type                                                    | Required | Restrictions | Description                                                         |
|------------------------|---------------------------------------------------------|----------|--------------|---------------------------------------------------------------------|
| `announcement_banners` | array of [optimus-ide-collabsdk.BannerConfig](#optimus-ide-collabsdkbannerconfig) | false    |              |                                                                     |
| `application_name`     | string                                                  | false    |              |                                                                     |
| `logo_url`             | string                                                  | false    |              |                                                                     |
| `service_banner`       | [optimus-ide-collabsdk.BannerConfig](#optimus-ide-collabsdkbannerconfig)          | false    |              | Deprecated: ServiceBanner has been replaced by AnnouncementBanners. |

## optimus-ide-collabsdk.UpdateChatACL

```json
{
  "group_roles": {
    "property1": "read",
    "property2": "read"
  },
  "user_roles": {
    "property1": "read",
    "property2": "read"
  }
}
```

### Properties

| Name               | Type                                   | Required | Restrictions | Description |
|--------------------|----------------------------------------|----------|--------------|-------------|
| `group_roles`      | object                                 | false    |              |             |
| » `[any property]` | [optimus-ide-collabsdk.ChatRole](#optimus-ide-collabsdkchatrole) | false    |              |             |
| `user_roles`       | object                                 | false    |              |             |
| » `[any property]` | [optimus-ide-collabsdk.ChatRole](#optimus-ide-collabsdkchatrole) | false    |              |             |

## optimus-ide-collabsdk.UpdateChatRequest

```json
{
  "archived": true,
  "labels": {
    "property1": "string",
    "property2": "string"
  },
  "pin_order": 0,
  "plan_mode": "plan",
  "title": "string",
  "workspace_id": "0967198e-ec7b-4c6b-b4d3-f71244cadbe9"
}
```

### Properties

| Name               | Type                                           | Required | Restrictions | Description                                                                                                                                                                                                                                                                                                                                                                                                                             |
|--------------------|------------------------------------------------|----------|--------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `archived`         | boolean                                        | false    |              |                                                                                                                                                                                                                                                                                                                                                                                                                                         |
| `labels`           | object                                         | false    |              |                                                                                                                                                                                                                                                                                                                                                                                                                                         |
| » `[any property]` | string                                         | false    |              |                                                                                                                                                                                                                                                                                                                                                                                                                                         |
| `pin_order`        | integer                                        | false    |              | Pin order controls the chat's pinned state and position. - nil: no change to pin state. - 0: unpin the chat. - >0 (chat is unpinned): pin the chat, appending it to   the end of the pinned list. The specific value is   ignored; the server assigns the next available position. - >0 (chat is already pinned): move the chat to the   requested position, shifting neighbors as needed. The   value is clamped to [1, pinned_count]. |
| `plan_mode`        | [optimus-ide-collabsdk.ChatPlanMode](#optimus-ide-collabsdkchatplanmode) | false    |              | Plan mode switches the chat's persistent plan mode. nil: no change, ptr to "plan": enable, ptr to "": clear.                                                                                                                                                                                                                                                                                                                            |
| `title`            | string                                         | false    |              |                                                                                                                                                                                                                                                                                                                                                                                                                                         |
| `workspace_id`     | string                                         | false    |              |                                                                                                                                                                                                                                                                                                                                                                                                                                         |

## optimus-ide-collabsdk.UpdateChatRetentionDaysRequest

```json
{
  "retention_days": 0
}
```

### Properties

| Name             | Type    | Required | Restrictions | Description |
|------------------|---------|----------|--------------|-------------|
| `retention_days` | integer | false    |              |             |

## optimus-ide-collabsdk.UpdateCheckResponse

```json
{
  "current": true,
  "url": "string",
  "version": "string"
}
```

### Properties

| Name      | Type    | Required | Restrictions | Description                                                             |
|-----------|---------|----------|--------------|-------------------------------------------------------------------------|
| `current` | boolean | false    |              | Current indicates whether the server version is the same as the latest. |
| `url`     | string  | false    |              | URL to download the latest release of Optimus-IDE-Collab.                            |
| `version` | string  | false    |              | Version is the semantic version for the latest release of Optimus-IDE-Collab.        |

## optimus-ide-collabsdk.UpdateOrganizationRequest

```json
{
  "default_org_member_roles": [
    "string"
  ],
  "description": "string",
  "display_name": "string",
  "icon": "string",
  "name": "string"
}
```

### Properties

| Name                       | Type            | Required | Restrictions | Description                                                                     |
|----------------------------|-----------------|----------|--------------|---------------------------------------------------------------------------------|
| `default_org_member_roles` | array of string | false    |              | Default org member roles when non-nil, replaces the org's default member roles. |
| `description`              | string          | false    |              |                                                                                 |
| `display_name`             | string          | false    |              |                                                                                 |
| `icon`                     | string          | false    |              |                                                                                 |
| `name`                     | string          | false    |              |                                                                                 |

## optimus-ide-collabsdk.UpdateRoles

```json
{
  "roles": [
    "string"
  ]
}
```

### Properties

| Name    | Type            | Required | Restrictions | Description |
|---------|-----------------|----------|--------------|-------------|
| `roles` | array of string | false    |              |             |

## optimus-ide-collabsdk.UpdateTaskInputRequest

```json
{
  "input": "string"
}
```

### Properties

| Name    | Type   | Required | Restrictions | Description |
|---------|--------|----------|--------------|-------------|
| `input` | string | false    |              |             |

## optimus-ide-collabsdk.UpdateTemplateACL

```json
{
  "group_perms": {
    "8bd26b20-f3e8-48be-a903-46bb920cf671": "use",
    "<group_id>": "admin"
  },
  "user_perms": {
    "4df59e74-c027-470b-ab4d-cbba8963a5e9": "use",
    "<user_id>": "admin"
  }
}
```

### Properties

| Name               | Type                                           | Required | Restrictions | Description                                                                                                                                                                                                       |
|--------------------|------------------------------------------------|----------|--------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `group_perms`      | object                                         | false    |              | Group perms is a mapping from valid group UUIDs to the template role they should be granted. To remove a group from the template, use "" as the role (available as a constant named optimus-ide-collabsdk.TemplateRoleDeleted) |
| » `[any property]` | [optimus-ide-collabsdk.TemplateRole](#optimus-ide-collabsdktemplaterole) | false    |              |                                                                                                                                                                                                                   |
| `user_perms`       | object                                         | false    |              | User perms is a mapping from valid user UUIDs to the template role they should be granted. To remove a user from the template, use "" as the role (available as a constant named optimus-ide-collabsdk.TemplateRoleDeleted)    |
| » `[any property]` | [optimus-ide-collabsdk.TemplateRole](#optimus-ide-collabsdktemplaterole) | false    |              |                                                                                                                                                                                                                   |

## optimus-ide-collabsdk.UpdateTemplateMeta

```json
{
  "activity_bump_ms": 0,
  "allow_user_autostart": true,
  "allow_user_autostop": true,
  "allow_user_cancel_workspace_jobs": true,
  "autostart_requirement": {
    "days_of_week": [
      "monday"
    ]
  },
  "autostop_requirement": {
    "days_of_week": [
      "monday"
    ],
    "weeks": 0
  },
  "cors_behavior": "simple",
  "default_ttl_ms": 0,
  "deprecation_message": "string",
  "description": "string",
  "disable_everyone_group_access": true,
  "disable_module_cache": true,
  "display_name": "string",
  "failure_ttl_ms": 0,
  "icon": "string",
  "max_port_share_level": "owner",
  "name": "string",
  "require_active_version": true,
  "time_til_autostop_notify_ms": 0,
  "time_til_dormant_autodelete_ms": 0,
  "time_til_dormant_ms": 0,
  "update_workspace_dormant_at": true,
  "update_workspace_last_used_at": true,
  "use_classic_parameter_flow": true
}
```

### Properties

| Name                               | Type                                                                           | Required | Restrictions | Description                                                                                                                                                                                                                                                                                                                                                                        |
|------------------------------------|--------------------------------------------------------------------------------|----------|--------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `activity_bump_ms`                 | integer                                                                        | false    |              | Activity bump ms allows optionally specifying the activity bump duration for all workspaces created from this template. Defaults to 1h but can be set to 0 to disable activity bumping.                                                                                                                                                                                            |
| `allow_user_autostart`             | boolean                                                                        | false    |              |                                                                                                                                                                                                                                                                                                                                                                                    |
| `allow_user_autostop`              | boolean                                                                        | false    |              |                                                                                                                                                                                                                                                                                                                                                                                    |
| `allow_user_cancel_workspace_jobs` | boolean                                                                        | false    |              |                                                                                                                                                                                                                                                                                                                                                                                    |
| `autostart_requirement`            | [optimus-ide-collabsdk.TemplateAutostartRequirement](#optimus-ide-collabsdktemplateautostartrequirement) | false    |              |                                                                                                                                                                                                                                                                                                                                                                                    |
| `autostop_requirement`             | [optimus-ide-collabsdk.TemplateAutostopRequirement](#optimus-ide-collabsdktemplateautostoprequirement)   | false    |              | Autostop requirement and AutostartRequirement can only be set if your license includes the advanced template scheduling feature. If you attempt to set this value while unlicensed, it will be ignored.                                                                                                                                                                            |
| `cors_behavior`                    | [optimus-ide-collabsdk.CORSBehavior](#optimus-ide-collabsdkcorsbehavior)                                 | false    |              |                                                                                                                                                                                                                                                                                                                                                                                    |
| `default_ttl_ms`                   | integer                                                                        | false    |              |                                                                                                                                                                                                                                                                                                                                                                                    |
| `deprecation_message`              | string                                                                         | false    |              | Deprecation message if set, will mark the template as deprecated and block any new workspaces from using this template. If passed an empty string, will remove the deprecated message, making the template usable for new workspaces again.                                                                                                                                        |
| `description`                      | string                                                                         | false    |              |                                                                                                                                                                                                                                                                                                                                                                                    |
| `disable_everyone_group_access`    | boolean                                                                        | false    |              | Disable everyone group access allows optionally disabling the default behavior of granting the 'everyone' group access to use the template. If this is set to true, the template will not be available to all users, and must be explicitly granted to users or groups in the permissions settings of the template.                                                                |
| `disable_module_cache`             | boolean                                                                        | false    |              | Disable module cache disables the using of cached Terraform modules during provisioning. It is recommended not to disable this.                                                                                                                                                                                                                                                    |
| `display_name`                     | string                                                                         | false    |              |                                                                                                                                                                                                                                                                                                                                                                                    |
| `failure_ttl_ms`                   | integer                                                                        | false    |              |                                                                                                                                                                                                                                                                                                                                                                                    |
| `icon`                             | string                                                                         | false    |              |                                                                                                                                                                                                                                                                                                                                                                                    |
| `max_port_share_level`             | [optimus-ide-collabsdk.WorkspaceAgentPortShareLevel](#optimus-ide-collabsdkworkspaceagentportsharelevel) | false    |              |                                                                                                                                                                                                                                                                                                                                                                                    |
| `name`                             | string                                                                         | false    |              |                                                                                                                                                                                                                                                                                                                                                                                    |
| `require_active_version`           | boolean                                                                        | false    |              | Require active version mandates workspaces built using this template use the active version of the template. This option has no effect on template admins.                                                                                                                                                                                                                         |
| `time_til_autostop_notify_ms`      | integer                                                                        | false    |              | Time til autostop notify ms allows optionally specifying the duration before the autostop deadline at which a reminder notification is sent for workspaces created from this template. Defaults to 0 (disabled). Omitting the field keeps the existing value.                                                                                                                      |
| `time_til_dormant_autodelete_ms`   | integer                                                                        | false    |              |                                                                                                                                                                                                                                                                                                                                                                                    |
| `time_til_dormant_ms`              | integer                                                                        | false    |              |                                                                                                                                                                                                                                                                                                                                                                                    |
| `update_workspace_dormant_at`      | boolean                                                                        | false    |              | Update workspace dormant at updates the dormant_at field of workspaces spawned from the template. This is useful for preventing dormant workspaces being immediately deleted when updating the dormant_ttl field to a new, shorter value.                                                                                                                                          |
| `update_workspace_last_used_at`    | boolean                                                                        | false    |              | Update workspace last used at updates the last_used_at field of workspaces spawned from the template. This is useful for preventing workspaces being immediately locked when updating the inactivity_ttl field to a new, shorter value.                                                                                                                                            |
| `use_classic_parameter_flow`       | boolean                                                                        | false    |              | Use classic parameter flow is a flag that switches the default behavior to use the classic parameter flow when creating a workspace. This only affects deployments with the experiment "dynamic-parameters" enabled. This setting will live for a period after the experiment is made the default. An "opt-out" is present in case the new feature breaks some existing templates. |

## optimus-ide-collabsdk.UpdateUserAppearanceSettingsRequest

```json
{
  "terminal_font": "",
  "theme_dark": "light",
  "theme_light": "light",
  "theme_mode": "sync",
  "theme_preference": "string"
}
```

### Properties

| Name               | Type                                                   | Required | Restrictions | Description                                                                                                                                                                                                                                                                                                        |
|--------------------|--------------------------------------------------------|----------|--------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `terminal_font`    | [optimus-ide-collabsdk.TerminalFontName](#optimus-ide-collabsdkterminalfontname) | true     |              |                                                                                                                                                                                                                                                                                                                    |
| `theme_dark`       | string                                                 | false    |              | Theme dark is required when ThemeMode is "sync". In "single" mode an empty value means "preserve the previously persisted slot" rather than "clear the slot", so partial updates that send only one slot keep the other intact.                                                                                    |
| `theme_light`      | string                                                 | false    |              | Theme light is required when ThemeMode is "sync". In "single" mode an empty value means "preserve the previously persisted slot" rather than "clear the slot", so partial updates that send only one slot keep the other intact.                                                                                   |
| `theme_mode`       | [optimus-ide-collabsdk.ThemeMode](#optimus-ide-collabsdkthememode)               | false    |              | Theme mode is optional for backward compatibility. When empty, the server leaves theme_mode, theme_light, and theme_dark unchanged so older CLI clients do not erase sync-mode settings. Legacy auto preferences are the exception: they clear theme_mode so clients can migrate the old sync-with-system setting. |
| `theme_preference` | string                                                 | true     |              |                                                                                                                                                                                                                                                                                                                    |

#### Enumerated Values

| Property      | Value(s)                                                                                    |
|---------------|---------------------------------------------------------------------------------------------|
| `theme_dark`  | `dark`, `dark-protan-deuter`, `dark-tritan`, `light`, `light-protan-deuter`, `light-tritan` |
| `theme_light` | `dark`, `dark-protan-deuter`, `dark-tritan`, `light`, `light-protan-deuter`, `light-tritan` |
| `theme_mode`  | `single`, `sync`                                                                            |

## optimus-ide-collabsdk.UpdateUserNotificationPreferences

```json
{
  "template_disabled_map": {
    "property1": true,
    "property2": true
  }
}
```

### Properties

| Name                    | Type    | Required | Restrictions | Description |
|-------------------------|---------|----------|--------------|-------------|
| `template_disabled_map` | object  | false    |              |             |
| » `[any property]`      | boolean | false    |              |             |

## optimus-ide-collabsdk.UpdateUserPasswordRequest

```json
{
  "old_password": "string",
  "password": "string"
}
```

### Properties

| Name           | Type   | Required | Restrictions | Description |
|----------------|--------|----------|--------------|-------------|
| `old_password` | string | false    |              |             |
| `password`     | string | true     |              |             |

## optimus-ide-collabsdk.UpdateUserPreferenceSettingsRequest

```json
{
  "agent_chat_send_shortcut": "enter",
  "code_diff_display_mode": "auto",
  "shell_tool_display_mode": "auto",
  "task_notification_alert_dismissed": true,
  "thinking_display_mode": "auto"
}
```

### Properties

| Name                                | Type                                                             | Required | Restrictions | Description |
|-------------------------------------|------------------------------------------------------------------|----------|--------------|-------------|
| `agent_chat_send_shortcut`          | [optimus-ide-collabsdk.AgentChatSendShortcut](#optimus-ide-collabsdkagentchatsendshortcut) | false    |              |             |
| `code_diff_display_mode`            | [optimus-ide-collabsdk.AgentDisplayMode](#optimus-ide-collabsdkagentdisplaymode)           | false    |              |             |
| `shell_tool_display_mode`           | [optimus-ide-collabsdk.AgentDisplayMode](#optimus-ide-collabsdkagentdisplaymode)           | false    |              |             |
| `task_notification_alert_dismissed` | boolean                                                          | false    |              |             |
| `thinking_display_mode`             | [optimus-ide-collabsdk.ThinkingDisplayMode](#optimus-ide-collabsdkthinkingdisplaymode)     | false    |              |             |

## optimus-ide-collabsdk.UpdateUserProfileRequest

```json
{
  "avatar_url": "http://example.com",
  "name": "string",
  "username": "string"
}
```

### Properties

| Name         | Type   | Required | Restrictions | Description                                                                                                                                                                                 |
|--------------|--------|----------|--------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `avatar_url` | string | false    |              | Avatar URL is only applied for users whose login type is password or none. For other login types the avatar is synced from the identity provider on login, so a submitted value is ignored. |
| `name`       | string | false    |              |                                                                                                                                                                                             |
| `username`   | string | true     |              |                                                                                                                                                                                             |

## optimus-ide-collabsdk.UpdateUserQuietHoursScheduleRequest

```json
{
  "schedule": "string"
}
```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|`schedule`|string|true||Schedule is a cron expression that defines when the user's quiet hours window is. Schedule must not be empty. For new users, the schedule is set to 2am in their browser or computer's timezone. The schedule denotes the beginning of a 4 hour window where the workspace is allowed to automatically stop or restart due to maintenance or template schedule.
The schedule must be daily with a single time, and should have a timezone specified via a CRON_TZ prefix (otherwise UTC will be used).
If the schedule is empty, the user will be updated to use the default schedule.|

## optimus-ide-collabsdk.UpdateUserSecretRequest

```json
{
  "description": "string",
  "enabled": true,
  "env_name": "string",
  "file_path": "string",
  "value": "string"
}
```

### Properties

| Name          | Type    | Required | Restrictions | Description |
|---------------|---------|----------|--------------|-------------|
| `description` | string  | false    |              |             |
| `enabled`     | boolean | false    |              |             |
| `env_name`    | string  | false    |              |             |
| `file_path`   | string  | false    |              |             |
| `value`       | string  | false    |              |             |

## optimus-ide-collabsdk.UpdateUserSkillRequest

```json
{
  "content": "string"
}
```

### Properties

| Name      | Type   | Required | Restrictions | Description                                                                                                                                                           |
|-----------|--------|----------|--------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `content` | string | false    |              | Content must be SKILL.md-format Markdown with YAML frontmatter. The frontmatter must include name, may include description, and must be followed by a non-empty body. |

## optimus-ide-collabsdk.UpdateWorkspaceACL

```json
{
  "group_roles": {
    "property1": "admin",
    "property2": "admin"
  },
  "user_roles": {
    "property1": "admin",
    "property2": "admin"
  }
}
```

### Properties

| Name               | Type                                             | Required | Restrictions | Description                                                                                                                                                                                                          |
|--------------------|--------------------------------------------------|----------|--------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `group_roles`      | object                                           | false    |              | Group roles is a mapping from valid group UUIDs to the workspace role they should be granted. To remove a group from the workspace, use "" as the role (available as a constant named optimus-ide-collabsdk.WorkspaceRoleDeleted) |
| » `[any property]` | [optimus-ide-collabsdk.WorkspaceRole](#optimus-ide-collabsdkworkspacerole) | false    |              |                                                                                                                                                                                                                      |
| `user_roles`       | object                                           | false    |              | User roles is a mapping from valid user UUIDs to the workspace role they should be granted. To remove a user from the workspace, use "" as the role (available as a constant named optimus-ide-collabsdk.WorkspaceRoleDeleted)    |
| » `[any property]` | [optimus-ide-collabsdk.WorkspaceRole](#optimus-ide-collabsdkworkspacerole) | false    |              |                                                                                                                                                                                                                      |

## optimus-ide-collabsdk.UpdateWorkspaceAutomaticUpdatesRequest

```json
{
  "automatic_updates": "always"
}
```

### Properties

| Name                | Type                                                   | Required | Restrictions | Description |
|---------------------|--------------------------------------------------------|----------|--------------|-------------|
| `automatic_updates` | [optimus-ide-collabsdk.AutomaticUpdates](#optimus-ide-collabsdkautomaticupdates) | false    |              |             |

## optimus-ide-collabsdk.UpdateWorkspaceAutostartRequest

```json
{
  "schedule": "string"
}
```

### Properties

| Name       | Type   | Required | Restrictions | Description                                                                                                                                                                                                                                    |
|------------|--------|----------|--------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `schedule` | string | false    |              | Schedule is expected to be of the form `CRON_TZ=<IANA Timezone> <min> <hour> * * <dow>` Example: `CRON_TZ=US/Central 30 9 * * 1-5` represents 0930 in the timezone US/Central on weekdays (Mon-Fri). `CRON_TZ` defaults to UTC if not present. |

## optimus-ide-collabsdk.UpdateWorkspaceBuildStateRequest

```json
{
  "state": [
    0
  ]
}
```

### Properties

| Name    | Type             | Required | Restrictions | Description |
|---------|------------------|----------|--------------|-------------|
| `state` | array of integer | false    |              |             |

## optimus-ide-collabsdk.UpdateWorkspaceDormancy

```json
{
  "dormant": true
}
```

### Properties

| Name      | Type    | Required | Restrictions | Description |
|-----------|---------|----------|--------------|-------------|
| `dormant` | boolean | false    |              |             |

## optimus-ide-collabsdk.UpdateWorkspaceRequest

```json
{
  "name": "string"
}
```

### Properties

| Name   | Type   | Required | Restrictions | Description |
|--------|--------|----------|--------------|-------------|
| `name` | string | false    |              |             |

## optimus-ide-collabsdk.UpdateWorkspaceSharingSettingsRequest

```json
{
  "shareable_workspace_owners": "none",
  "sharing_disabled": true
}
```

### Properties

| Name                         | Type                                                                   | Required | Restrictions | Description                                                                                                                     |
|------------------------------|------------------------------------------------------------------------|----------|--------------|---------------------------------------------------------------------------------------------------------------------------------|
| `shareable_workspace_owners` | [optimus-ide-collabsdk.ShareableWorkspaceOwners](#optimus-ide-collabsdkshareableworkspaceowners) | false    |              | Shareable workspace owners controls whose workspaces can be shared within the organization.                                     |
| `sharing_disabled`           | boolean                                                                | false    |              | Sharing disabled is deprecated and left for backward compatibility purposes. Deprecated: use `ShareableWorkspaceOwners` instead |

#### Enumerated Values

| Property                     | Value(s)                               |
|------------------------------|----------------------------------------|
| `shareable_workspace_owners` | `everyone`, `none`, `service_accounts` |

## optimus-ide-collabsdk.UpdateWorkspaceTTLRequest

```json
{
  "ttl_ms": 0
}
```

### Properties

| Name     | Type    | Required | Restrictions | Description |
|----------|---------|----------|--------------|-------------|
| `ttl_ms` | integer | false    |              |             |

## optimus-ide-collabsdk.UploadChatFileResponse

```json
{
  "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08"
}
```

### Properties

| Name | Type   | Required | Restrictions | Description |
|------|--------|----------|--------------|-------------|
| `id` | string | false    |              |             |

## optimus-ide-collabsdk.UploadResponse

```json
{
  "hash": "19686d84-b10d-4f90-b18e-84fd3fa038fd"
}
```

### Properties

| Name   | Type   | Required | Restrictions | Description |
|--------|--------|----------|--------------|-------------|
| `hash` | string | false    |              |             |

## optimus-ide-collabsdk.UpsertGroupAIBudgetRequest

```json
{
  "spend_limit_micros": 0
}
```

### Properties

| Name                 | Type    | Required | Restrictions | Description |
|----------------------|---------|----------|--------------|-------------|
| `spend_limit_micros` | integer | false    |              |             |

## optimus-ide-collabsdk.UpsertUserAIBudgetOverrideRequest

```json
{
  "group_id": "306db4e0-7449-4501-b76f-075576fe2d8f",
  "spend_limit_micros": 0
}
```

### Properties

| Name                 | Type    | Required | Restrictions | Description                                                                                       |
|----------------------|---------|----------|--------------|---------------------------------------------------------------------------------------------------|
| `group_id`           | string  | true     |              | Group ID is the group the user's spend is attributed to. The user must be a member of this group. |
| `spend_limit_micros` | integer | false    |              |                                                                                                   |

## optimus-ide-collabsdk.UpsertWorkspaceAgentPortShareRequest

```json
{
  "agent_name": "string",
  "port": 0,
  "protocol": "http",
  "share_level": "owner"
}
```

### Properties

| Name          | Type                                                                                 | Required | Restrictions | Description |
|---------------|--------------------------------------------------------------------------------------|----------|--------------|-------------|
| `agent_name`  | string                                                                               | false    |              |             |
| `port`        | integer                                                                              | false    |              |             |
| `protocol`    | [optimus-ide-collabsdk.WorkspaceAgentPortShareProtocol](#optimus-ide-collabsdkworkspaceagentportshareprotocol) | false    |              |             |
| `share_level` | [optimus-ide-collabsdk.WorkspaceAgentPortShareLevel](#optimus-ide-collabsdkworkspaceagentportsharelevel)       | false    |              |             |

#### Enumerated Values

| Property      | Value(s)                                           |
|---------------|----------------------------------------------------|
| `protocol`    | `http`, `https`                                    |
| `share_level` | `authenticated`, `organization`, `owner`, `public` |

## optimus-ide-collabsdk.UsageAppName

```json
"vscode"
```

### Properties

#### Enumerated Values

| Value(s)                                         |
|--------------------------------------------------|
| `jetbrains`, `reconnecting-pty`, `ssh`, `vscode` |

## optimus-ide-collabsdk.UsagePeriod

```json
{
  "end": "2019-08-24T14:15:22Z",
  "issued_at": "2019-08-24T14:15:22Z",
  "start": "2019-08-24T14:15:22Z"
}
```

### Properties

| Name        | Type   | Required | Restrictions | Description |
|-------------|--------|----------|--------------|-------------|
| `end`       | string | false    |              |             |
| `issued_at` | string | false    |              |             |
| `start`     | string | false    |              |             |

## optimus-ide-collabsdk.UsageStatsConfig

```json
{
  "enable": true
}
```

### Properties

| Name     | Type    | Required | Restrictions | Description |
|----------|---------|----------|--------------|-------------|
| `enable` | boolean | false    |              |             |

## optimus-ide-collabsdk.User

```json
{
  "avatar_url": "http://example.com",
  "created_at": "2019-08-24T14:15:22Z",
  "email": "user@example.com",
  "has_ai_seat": true,
  "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
  "is_service_account": true,
  "last_seen_at": "2019-08-24T14:15:22Z",
  "login_type": "",
  "name": "string",
  "organization_ids": [
    "497f6eca-6276-4993-bfeb-53cbbbba6f08"
  ],
  "roles": [
    {
      "display_name": "string",
      "name": "string",
      "organization_id": "string"
    }
  ],
  "status": "active",
  "theme_preference": "string",
  "updated_at": "2019-08-24T14:15:22Z",
  "username": "string"
}
```

### Properties

| Name                 | Type                                            | Required | Restrictions | Description                                                                                      |
|----------------------|-------------------------------------------------|----------|--------------|--------------------------------------------------------------------------------------------------|
| `avatar_url`         | string                                          | false    |              |                                                                                                  |
| `created_at`         | string                                          | true     |              |                                                                                                  |
| `email`              | string                                          | true     |              |                                                                                                  |
| `has_ai_seat`        | boolean                                         | false    |              | Has ai seat intentionally omits omitempty so the API always includes the field, even when false. |
| `id`                 | string                                          | true     |              |                                                                                                  |
| `is_service_account` | boolean                                         | false    |              |                                                                                                  |
| `last_seen_at`       | string                                          | false    |              |                                                                                                  |
| `login_type`         | [optimus-ide-collabsdk.LoginType](#optimus-ide-collabsdklogintype)        | false    |              |                                                                                                  |
| `name`               | string                                          | false    |              |                                                                                                  |
| `organization_ids`   | array of string                                 | false    |              |                                                                                                  |
| `roles`              | array of [optimus-ide-collabsdk.SlimRole](#optimus-ide-collabsdkslimrole) | false    |              |                                                                                                  |
| `status`             | [optimus-ide-collabsdk.UserStatus](#optimus-ide-collabsdkuserstatus)      | false    |              |                                                                                                  |
| `theme_preference`   | string                                          | false    |              | Deprecated: this value should be retrieved from `optimus-ide-collabsdk.UserPreferenceSettings` instead.       |
| `updated_at`         | string                                          | false    |              |                                                                                                  |
| `username`           | string                                          | true     |              |                                                                                                  |

#### Enumerated Values

| Property | Value(s)              |
|----------|-----------------------|
| `status` | `active`, `suspended` |

## optimus-ide-collabsdk.UserAIBudgetOverride

```json
{
  "created_at": "2019-08-24T14:15:22Z",
  "group_id": "306db4e0-7449-4501-b76f-075576fe2d8f",
  "spend_limit_micros": 0,
  "updated_at": "2019-08-24T14:15:22Z",
  "user_id": "a169451c-8525-4352-b8ca-070dd449a1a5"
}
```

### Properties

| Name                 | Type    | Required | Restrictions | Description |
|----------------------|---------|----------|--------------|-------------|
| `created_at`         | string  | false    |              |             |
| `group_id`           | string  | false    |              |             |
| `spend_limit_micros` | integer | false    |              |             |
| `updated_at`         | string  | false    |              |             |
| `user_id`            | string  | false    |              |             |

## optimus-ide-collabsdk.UserAISpendStatus

```json
{
  "current_spend_micros": 0,
  "effective_group_id": "85e2b926-ddfb-4c66-b68e-b66e5acec6c0",
  "limit_source": "user_override",
  "period_end": "2019-08-24T14:15:22Z",
  "period_start": "2019-08-24T14:15:22Z",
  "spend_limit_micros": 0,
  "user_id": "a169451c-8525-4352-b8ca-070dd449a1a5"
}
```

### Properties

| Name                   | Type                                                         | Required | Restrictions | Description                                                                                                                                                                    |
|------------------------|--------------------------------------------------------------|----------|--------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `current_spend_micros` | integer                                                      | false    |              | Current spend micros is the user's spend on their effective group over the current budget period.                                                                              |
| `effective_group_id`   | string                                                       | false    |              | Effective group ID is the group the spend is attributed to, falling back to the Everyone group when no budget applies. Null only when the user has no organization membership. |
| `limit_source`         | [optimus-ide-collabsdk.AIBudgetLimitSource](#optimus-ide-collabsdkaibudgetlimitsource) | false    |              | Limit source identifies which tier produced the limit. Null when no budget applies.                                                                                            |
| `period_end`           | string                                                       | false    |              | Period end is the exclusive upper bound of the current budget period.                                                                                                          |
| `period_start`         | string                                                       | false    |              | Period start is the inclusive lower bound of the current budget period.                                                                                                        |
| `spend_limit_micros`   | integer                                                      | false    |              | Spend limit micros is the effective spend limit in micro-units. Null when no budget applies to the user (unlimited).                                                           |
| `user_id`              | string                                                       | false    |              |                                                                                                                                                                                |

## optimus-ide-collabsdk.UserActivity

```json
{
  "avatar_url": "http://example.com",
  "seconds": 80500,
  "template_ids": [
    "497f6eca-6276-4993-bfeb-53cbbbba6f08"
  ],
  "user_id": "a169451c-8525-4352-b8ca-070dd449a1a5",
  "username": "string"
}
```

### Properties

| Name           | Type            | Required | Restrictions | Description |
|----------------|-----------------|----------|--------------|-------------|
| `avatar_url`   | string          | false    |              |             |
| `seconds`      | integer         | false    |              |             |
| `template_ids` | array of string | false    |              |             |
| `user_id`      | string          | false    |              |             |
| `username`     | string          | false    |              |             |

## optimus-ide-collabsdk.UserActivityInsightsReport

```json
{
  "end_time": "2019-08-24T14:15:22Z",
  "start_time": "2019-08-24T14:15:22Z",
  "template_ids": [
    "497f6eca-6276-4993-bfeb-53cbbbba6f08"
  ],
  "users": [
    {
      "avatar_url": "http://example.com",
      "seconds": 80500,
      "template_ids": [
        "497f6eca-6276-4993-bfeb-53cbbbba6f08"
      ],
      "user_id": "a169451c-8525-4352-b8ca-070dd449a1a5",
      "username": "string"
    }
  ]
}
```

### Properties

| Name           | Type                                                    | Required | Restrictions | Description |
|----------------|---------------------------------------------------------|----------|--------------|-------------|
| `end_time`     | string                                                  | false    |              |             |
| `start_time`   | string                                                  | false    |              |             |
| `template_ids` | array of string                                         | false    |              |             |
| `users`        | array of [optimus-ide-collabsdk.UserActivity](#optimus-ide-collabsdkuseractivity) | false    |              |             |

## optimus-ide-collabsdk.UserActivityInsightsResponse

```json
{
  "report": {
    "end_time": "2019-08-24T14:15:22Z",
    "start_time": "2019-08-24T14:15:22Z",
    "template_ids": [
      "497f6eca-6276-4993-bfeb-53cbbbba6f08"
    ],
    "users": [
      {
        "avatar_url": "http://example.com",
        "seconds": 80500,
        "template_ids": [
          "497f6eca-6276-4993-bfeb-53cbbbba6f08"
        ],
        "user_id": "a169451c-8525-4352-b8ca-070dd449a1a5",
        "username": "string"
      }
    ]
  }
}
```

### Properties

| Name     | Type                                                                       | Required | Restrictions | Description |
|----------|----------------------------------------------------------------------------|----------|--------------|-------------|
| `report` | [optimus-ide-collabsdk.UserActivityInsightsReport](#optimus-ide-collabsdkuseractivityinsightsreport) | false    |              |             |

## optimus-ide-collabsdk.UserAppearanceSettings

```json
{
  "terminal_font": "",
  "theme_dark": "string",
  "theme_light": "string",
  "theme_mode": "",
  "theme_preference": "string"
}
```

### Properties

| Name               | Type                                                   | Required | Restrictions | Description                                                                                                                                                                                                                                                                                                                               |
|--------------------|--------------------------------------------------------|----------|--------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `terminal_font`    | [optimus-ide-collabsdk.TerminalFontName](#optimus-ide-collabsdkterminalfontname) | false    |              |                                                                                                                                                                                                                                                                                                                                           |
| `theme_dark`       | string                                                 | false    |              | Ignored when ThemeMode is "single"                                                                                                                                                                                                                                                                                                        |
| `theme_light`      | string                                                 | false    |              | Ignored when ThemeMode is "single"                                                                                                                                                                                                                                                                                                        |
| `theme_mode`       | [optimus-ide-collabsdk.ThemeMode](#optimus-ide-collabsdkthememode)               | false    |              |                                                                                                                                                                                                                                                                                                                                           |
| `theme_preference` | string                                                 | false    |              | Theme preference is the legacy single-field appearance setting. In "single" mode it mirrors the active theme. In "sync" mode modern clients normally mirror the active OS slot, but older clients can update only this field, so it may diverge from ThemeLight or ThemeDark until a modern client saves the full appearance state again. |

## optimus-ide-collabsdk.UserLatency

```json
{
  "avatar_url": "http://example.com",
  "latency_ms": {
    "p50": 31.312,
    "p95": 119.832
  },
  "template_ids": [
    "497f6eca-6276-4993-bfeb-53cbbbba6f08"
  ],
  "user_id": "a169451c-8525-4352-b8ca-070dd449a1a5",
  "username": "string"
}
```

### Properties

| Name           | Type                                                     | Required | Restrictions | Description |
|----------------|----------------------------------------------------------|----------|--------------|-------------|
| `avatar_url`   | string                                                   | false    |              |             |
| `latency_ms`   | [optimus-ide-collabsdk.ConnectionLatency](#optimus-ide-collabsdkconnectionlatency) | false    |              |             |
| `template_ids` | array of string                                          | false    |              |             |
| `user_id`      | string                                                   | false    |              |             |
| `username`     | string                                                   | false    |              |             |

## optimus-ide-collabsdk.UserLatencyInsightsReport

```json
{
  "end_time": "2019-08-24T14:15:22Z",
  "start_time": "2019-08-24T14:15:22Z",
  "template_ids": [
    "497f6eca-6276-4993-bfeb-53cbbbba6f08"
  ],
  "users": [
    {
      "avatar_url": "http://example.com",
      "latency_ms": {
        "p50": 31.312,
        "p95": 119.832
      },
      "template_ids": [
        "497f6eca-6276-4993-bfeb-53cbbbba6f08"
      ],
      "user_id": "a169451c-8525-4352-b8ca-070dd449a1a5",
      "username": "string"
    }
  ]
}
```

### Properties

| Name           | Type                                                  | Required | Restrictions | Description |
|----------------|-------------------------------------------------------|----------|--------------|-------------|
| `end_time`     | string                                                | false    |              |             |
| `start_time`   | string                                                | false    |              |             |
| `template_ids` | array of string                                       | false    |              |             |
| `users`        | array of [optimus-ide-collabsdk.UserLatency](#optimus-ide-collabsdkuserlatency) | false    |              |             |

## optimus-ide-collabsdk.UserLatencyInsightsResponse

```json
{
  "report": {
    "end_time": "2019-08-24T14:15:22Z",
    "start_time": "2019-08-24T14:15:22Z",
    "template_ids": [
      "497f6eca-6276-4993-bfeb-53cbbbba6f08"
    ],
    "users": [
      {
        "avatar_url": "http://example.com",
        "latency_ms": {
          "p50": 31.312,
          "p95": 119.832
        },
        "template_ids": [
          "497f6eca-6276-4993-bfeb-53cbbbba6f08"
        ],
        "user_id": "a169451c-8525-4352-b8ca-070dd449a1a5",
        "username": "string"
      }
    ]
  }
}
```

### Properties

| Name     | Type                                                                     | Required | Restrictions | Description |
|----------|--------------------------------------------------------------------------|----------|--------------|-------------|
| `report` | [optimus-ide-collabsdk.UserLatencyInsightsReport](#optimus-ide-collabsdkuserlatencyinsightsreport) | false    |              |             |

## optimus-ide-collabsdk.UserLoginType

```json
{
  "login_type": ""
}
```

### Properties

| Name         | Type                                     | Required | Restrictions | Description |
|--------------|------------------------------------------|----------|--------------|-------------|
| `login_type` | [optimus-ide-collabsdk.LoginType](#optimus-ide-collabsdklogintype) | false    |              |             |

## optimus-ide-collabsdk.UserParameter

```json
{
  "name": "string",
  "value": "string"
}
```

### Properties

| Name    | Type   | Required | Restrictions | Description |
|---------|--------|----------|--------------|-------------|
| `name`  | string | false    |              |             |
| `value` | string | false    |              |             |

## optimus-ide-collabsdk.UserPreferenceSettings

```json
{
  "agent_chat_send_shortcut": "enter",
  "code_diff_display_mode": "auto",
  "shell_tool_display_mode": "auto",
  "task_notification_alert_dismissed": true,
  "thinking_display_mode": "auto"
}
```

### Properties

| Name                                | Type                                                             | Required | Restrictions | Description |
|-------------------------------------|------------------------------------------------------------------|----------|--------------|-------------|
| `agent_chat_send_shortcut`          | [optimus-ide-collabsdk.AgentChatSendShortcut](#optimus-ide-collabsdkagentchatsendshortcut) | false    |              |             |
| `code_diff_display_mode`            | [optimus-ide-collabsdk.AgentDisplayMode](#optimus-ide-collabsdkagentdisplaymode)           | false    |              |             |
| `shell_tool_display_mode`           | [optimus-ide-collabsdk.AgentDisplayMode](#optimus-ide-collabsdkagentdisplaymode)           | false    |              |             |
| `task_notification_alert_dismissed` | boolean                                                          | false    |              |             |
| `thinking_display_mode`             | [optimus-ide-collabsdk.ThinkingDisplayMode](#optimus-ide-collabsdkthinkingdisplaymode)     | false    |              |             |

## optimus-ide-collabsdk.UserQuietHoursScheduleConfig

```json
{
  "allow_user_custom": true,
  "default_schedule": "string"
}
```

### Properties

| Name                | Type    | Required | Restrictions | Description |
|---------------------|---------|----------|--------------|-------------|
| `allow_user_custom` | boolean | false    |              |             |
| `default_schedule`  | string  | false    |              |             |

## optimus-ide-collabsdk.UserQuietHoursScheduleResponse

```json
{
  "next": "2019-08-24T14:15:22Z",
  "raw_schedule": "string",
  "time": "string",
  "timezone": "string",
  "user_can_set": true,
  "user_set": true
}
```

### Properties

| Name           | Type    | Required | Restrictions | Description                                                                                                                                                                      |
|----------------|---------|----------|--------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `next`         | string  | false    |              | Next is the next time that the quiet hours window will start.                                                                                                                    |
| `raw_schedule` | string  | false    |              |                                                                                                                                                                                  |
| `time`         | string  | false    |              | Time is the time of day that the quiet hours window starts in the given Timezone each day.                                                                                       |
| `timezone`     | string  | false    |              | raw format from the cron expression, UTC if unspecified                                                                                                                          |
| `user_can_set` | boolean | false    |              | User can set is true if the user is allowed to set their own quiet hours schedule. If false, the user cannot set a custom schedule and the default schedule will always be used. |
| `user_set`     | boolean | false    |              | User set is true if the user has set their own quiet hours schedule. If false, the user is using the default schedule.                                                           |

## optimus-ide-collabsdk.UserSecret

```json
{
  "created_at": "2019-08-24T14:15:22Z",
  "description": "string",
  "enabled": true,
  "env_name": "string",
  "file_path": "string",
  "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
  "name": "string",
  "updated_at": "2019-08-24T14:15:22Z"
}
```

### Properties

| Name          | Type    | Required | Restrictions | Description                                                                                                                                                                                                                          |
|---------------|---------|----------|--------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `created_at`  | string  | false    |              |                                                                                                                                                                                                                                      |
| `description` | string  | false    |              |                                                                                                                                                                                                                                      |
| `enabled`     | boolean | false    |              | Enabled controls whether the secret is injected into workspaces. Disabled secrets remain visible and editable, but are not added to the agent manifest, so they are not exposed as environment variables or written to secret files. |
| `env_name`    | string  | false    |              |                                                                                                                                                                                                                                      |
| `file_path`   | string  | false    |              |                                                                                                                                                                                                                                      |
| `id`          | string  | false    |              |                                                                                                                                                                                                                                      |
| `name`        | string  | false    |              |                                                                                                                                                                                                                                      |
| `updated_at`  | string  | false    |              |                                                                                                                                                                                                                                      |

## optimus-ide-collabsdk.UserSkill

```json
{
  "content": "string",
  "created_at": "2019-08-24T14:15:22Z",
  "description": "string",
  "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
  "name": "string",
  "updated_at": "2019-08-24T14:15:22Z"
}
```

### Properties

| Name          | Type   | Required | Restrictions | Description |
|---------------|--------|----------|--------------|-------------|
| `content`     | string | false    |              |             |
| `created_at`  | string | false    |              |             |
| `description` | string | false    |              |             |
| `id`          | string | false    |              |             |
| `name`        | string | false    |              |             |
| `updated_at`  | string | false    |              |             |

## optimus-ide-collabsdk.UserSkillMetadata

```json
{
  "created_at": "2019-08-24T14:15:22Z",
  "description": "string",
  "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
  "name": "string",
  "updated_at": "2019-08-24T14:15:22Z"
}
```

### Properties

| Name          | Type   | Required | Restrictions | Description |
|---------------|--------|----------|--------------|-------------|
| `created_at`  | string | false    |              |             |
| `description` | string | false    |              |             |
| `id`          | string | false    |              |             |
| `name`        | string | false    |              |             |
| `updated_at`  | string | false    |              |             |

## optimus-ide-collabsdk.UserStatus

```json
"active"
```

### Properties

#### Enumerated Values

| Value(s)                         |
|----------------------------------|
| `active`, `dormant`, `suspended` |

## optimus-ide-collabsdk.UserStatusChangeCount

```json
{
  "count": 10,
  "date": "2019-08-24T14:15:22Z"
}
```

### Properties

| Name    | Type    | Required | Restrictions | Description |
|---------|---------|----------|--------------|-------------|
| `count` | integer | false    |              |             |
| `date`  | string  | false    |              |             |

## optimus-ide-collabsdk.ValidateUserPasswordRequest

```json
{
  "password": "string"
}
```

### Properties

| Name       | Type   | Required | Restrictions | Description |
|------------|--------|----------|--------------|-------------|
| `password` | string | true     |              |             |

## optimus-ide-collabsdk.ValidateUserPasswordResponse

```json
{
  "details": "string",
  "valid": true
}
```

### Properties

| Name      | Type    | Required | Restrictions | Description |
|-----------|---------|----------|--------------|-------------|
| `details` | string  | false    |              |             |
| `valid`   | boolean | false    |              |             |

## optimus-ide-collabsdk.ValidationError

```json
{
  "detail": "string",
  "field": "string"
}
```

### Properties

| Name     | Type   | Required | Restrictions | Description |
|----------|--------|----------|--------------|-------------|
| `detail` | string | true     |              |             |
| `field`  | string | true     |              |             |

## optimus-ide-collabsdk.ValidationMonotonicOrder

```json
"increasing"
```

### Properties

#### Enumerated Values

| Value(s)                   |
|----------------------------|
| `decreasing`, `increasing` |

## optimus-ide-collabsdk.VariableValue

```json
{
  "name": "string",
  "value": "string"
}
```

### Properties

| Name    | Type   | Required | Restrictions | Description |
|---------|--------|----------|--------------|-------------|
| `name`  | string | false    |              |             |
| `value` | string | false    |              |             |

## optimus-ide-collabsdk.WebpushSubscription

```json
{
  "auth_key": "string",
  "endpoint": "string",
  "p256dh_key": "string"
}
```

### Properties

| Name         | Type   | Required | Restrictions | Description |
|--------------|--------|----------|--------------|-------------|
| `auth_key`   | string | false    |              |             |
| `endpoint`   | string | false    |              |             |
| `p256dh_key` | string | false    |              |             |

## optimus-ide-collabsdk.Workspace

```json
{
  "allow_renames": true,
  "automatic_updates": "always",
  "autostart_schedule": "string",
  "created_at": "2019-08-24T14:15:22Z",
  "deleting_at": "2019-08-24T14:15:22Z",
  "dormant_at": "2019-08-24T14:15:22Z",
  "favorite": true,
  "health": {
    "failing_agents": [
      "497f6eca-6276-4993-bfeb-53cbbbba6f08"
    ],
    "healthy": false
  },
  "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
  "is_prebuild": true,
  "last_used_at": "2019-08-24T14:15:22Z",
  "latest_app_status": {
    "agent_id": "2b1e3b65-2c04-4fa2-a2d7-467901e98978",
    "app_id": "affd1d10-9538-4fc8-9e0b-4594a28c1335",
    "created_at": "2019-08-24T14:15:22Z",
    "icon": "string",
    "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
    "message": "string",
    "needs_user_attention": true,
    "state": "working",
    "uri": "string",
    "workspace_id": "0967198e-ec7b-4c6b-b4d3-f71244cadbe9"
  },
  "latest_build": {
    "build_number": 0,
    "created_at": "2019-08-24T14:15:22Z",
    "daily_cost": 0,
    "deadline": "2019-08-24T14:15:22Z",
    "has_ai_task": true,
    "has_external_agent": true,
    "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
    "initiator_id": "06588898-9a84-4b35-ba8f-f9cbd64946f3",
    "initiator_name": "string",
    "job": {
      "available_workers": [
        "497f6eca-6276-4993-bfeb-53cbbbba6f08"
      ],
      "canceled_at": "2019-08-24T14:15:22Z",
      "completed_at": "2019-08-24T14:15:22Z",
      "created_at": "2019-08-24T14:15:22Z",
      "error": "string",
      "error_code": "REQUIRED_TEMPLATE_VARIABLES",
      "file_id": "8a0cfb4f-ddc9-436d-91bb-75133c583767",
      "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
      "initiator_id": "06588898-9a84-4b35-ba8f-f9cbd64946f3",
      "input": {
        "error": "string",
        "template_version_id": "0ba39c92-1f1b-4c32-aa3e-9925d7713eb1",
        "workspace_build_id": "badaf2eb-96c5-4050-9f1d-db2d39ca5478"
      },
      "logs_overflowed": true,
      "metadata": {
        "template_display_name": "string",
        "template_icon": "string",
        "template_id": "c6d67e98-83ea-49f0-8812-e4abae2b68bc",
        "template_name": "string",
        "template_version_name": "string",
        "workspace_build_transition": "start",
        "workspace_id": "0967198e-ec7b-4c6b-b4d3-f71244cadbe9",
        "workspace_name": "string"
      },
      "organization_id": "7c60d51f-b44e-4682-87d6-449835ea4de6",
      "queue_position": 0,
      "queue_size": 0,
      "started_at": "2019-08-24T14:15:22Z",
      "status": "pending",
      "tags": {
        "property1": "string",
        "property2": "string"
      },
      "type": "template_version_import",
      "worker_id": "ae5fa6f7-c55b-40c1-b40a-b36ac467652b",
      "worker_name": "string"
    },
    "matched_provisioners": {
      "available": 0,
      "count": 0,
      "most_recently_seen": "2019-08-24T14:15:22Z"
    },
    "max_deadline": "2019-08-24T14:15:22Z",
    "reason": "initiator",
    "resources": [
      {
        "agents": [
          {
            "api_version": "string",
            "apps": [
              {
                "command": "string",
                "display_name": "string",
                "external": true,
                "group": "string",
                "health": "disabled",
                "healthcheck": {
                  "interval": 0,
                  "threshold": 0,
                  "url": "string"
                },
                "hidden": true,
                "icon": "string",
                "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
                "open_in": "slim-window",
                "sharing_level": "owner",
                "slug": "string",
                "statuses": [
                  {
                    "agent_id": "2b1e3b65-2c04-4fa2-a2d7-467901e98978",
                    "app_id": "affd1d10-9538-4fc8-9e0b-4594a28c1335",
                    "created_at": "2019-08-24T14:15:22Z",
                    "icon": "string",
                    "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
                    "message": "string",
                    "needs_user_attention": true,
                    "state": "working",
                    "uri": "string",
                    "workspace_id": "0967198e-ec7b-4c6b-b4d3-f71244cadbe9"
                  }
                ],
                "subdomain": true,
                "subdomain_name": "string",
                "tooltip": "string",
                "url": "string"
              }
            ],
            "architecture": "string",
            "connection_timeout_seconds": 0,
            "created_at": "2019-08-24T14:15:22Z",
            "directory": "string",
            "disconnected_at": "2019-08-24T14:15:22Z",
            "display_apps": [
              "vscode"
            ],
            "environment_variables": {
              "property1": "string",
              "property2": "string"
            },
            "expanded_directory": "string",
            "first_connected_at": "2019-08-24T14:15:22Z",
            "health": {
              "healthy": false,
              "reason": "agent has lost connection"
            },
            "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
            "instance_id": "string",
            "last_connected_at": "2019-08-24T14:15:22Z",
            "latency": {
              "property1": {
                "latency_ms": 0,
                "preferred": true
              },
              "property2": {
                "latency_ms": 0,
                "preferred": true
              }
            },
            "lifecycle_state": "created",
            "log_sources": [
              {
                "created_at": "2019-08-24T14:15:22Z",
                "display_name": "string",
                "icon": "string",
                "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
                "workspace_agent_id": "7ad2e618-fea7-4c1a-b70a-f501566a72f1"
              }
            ],
            "logs_length": 0,
            "logs_overflowed": true,
            "name": "string",
            "operating_system": "string",
            "parent_id": {
              "uuid": "string",
              "valid": true
            },
            "ready_at": "2019-08-24T14:15:22Z",
            "resource_id": "4d5215ed-38bb-48ed-879a-fdb9ca58522f",
            "scripts": [
              {
                "cron": "string",
                "display_name": "string",
                "exit_code": 0,
                "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
                "log_path": "string",
                "log_source_id": "4197ab25-95cf-4b91-9c78-f7f2af5d353a",
                "run_on_start": true,
                "run_on_stop": true,
                "script": "string",
                "start_blocks_login": true,
                "status": "ok",
                "timeout": 0
              }
            ],
            "started_at": "2019-08-24T14:15:22Z",
            "startup_script_behavior": "blocking",
            "status": "connecting",
            "subsystems": [
              "envbox"
            ],
            "troubleshooting_url": "string",
            "updated_at": "2019-08-24T14:15:22Z",
            "version": "string"
          }
        ],
        "created_at": "2019-08-24T14:15:22Z",
        "daily_cost": 0,
        "hide": true,
        "icon": "string",
        "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
        "job_id": "453bd7d7-5355-4d6d-a38e-d9e7eb218c3f",
        "metadata": [
          {
            "key": "string",
            "sensitive": true,
            "value": "string"
          }
        ],
        "name": "string",
        "type": "string",
        "workspace_transition": "start"
      }
    ],
    "status": "pending",
    "template_version_id": "0ba39c92-1f1b-4c32-aa3e-9925d7713eb1",
    "template_version_name": "string",
    "template_version_preset_id": "512a53a7-30da-446e-a1fc-713c630baff1",
    "transition": "start",
    "updated_at": "2019-08-24T14:15:22Z",
    "workspace_id": "0967198e-ec7b-4c6b-b4d3-f71244cadbe9",
    "workspace_name": "string",
    "workspace_owner_avatar_url": "string",
    "workspace_owner_id": "e7078695-5279-4c86-8774-3ac2367a2fc7",
    "workspace_owner_name": "string"
  },
  "name": "string",
  "next_start_at": "2019-08-24T14:15:22Z",
  "organization_id": "7c60d51f-b44e-4682-87d6-449835ea4de6",
  "organization_name": "string",
  "outdated": true,
  "owner_avatar_url": "string",
  "owner_id": "8826ee2e-7933-4665-aef2-2393f84a0d05",
  "owner_name": "string",
  "shared_with": [
    {
      "actor_type": "group",
      "avatar_url": "http://example.com",
      "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
      "name": "string",
      "roles": [
        "admin"
      ]
    }
  ],
  "task_id": {
    "uuid": "string",
    "valid": true
  },
  "template_active_version_id": "b0da9c29-67d8-4c87-888c-bafe356f7f3c",
  "template_allow_user_cancel_workspace_jobs": true,
  "template_display_name": "string",
  "template_icon": "string",
  "template_id": "c6d67e98-83ea-49f0-8812-e4abae2b68bc",
  "template_name": "string",
  "template_require_active_version": true,
  "template_use_classic_parameter_flow": true,
  "ttl_ms": 0,
  "updated_at": "2019-08-24T14:15:22Z"
}
```

### Properties

| Name                                        | Type                                                                    | Required | Restrictions | Description                                                                                                                                                                                                                                                                                                                                 |
|---------------------------------------------|-------------------------------------------------------------------------|----------|--------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `allow_renames`                             | boolean                                                                 | false    |              |                                                                                                                                                                                                                                                                                                                                             |
| `automatic_updates`                         | [optimus-ide-collabsdk.AutomaticUpdates](#optimus-ide-collabsdkautomaticupdates)                  | false    |              |                                                                                                                                                                                                                                                                                                                                             |
| `autostart_schedule`                        | string                                                                  | false    |              |                                                                                                                                                                                                                                                                                                                                             |
| `created_at`                                | string                                                                  | false    |              |                                                                                                                                                                                                                                                                                                                                             |
| `deleting_at`                               | string                                                                  | false    |              | Deleting at indicates the time at which the workspace will be permanently deleted. A workspace is eligible for deletion if it is dormant (a non-nil dormant_at value) and a value has been specified for time_til_dormant_autodelete on its template.                                                                                       |
| `dormant_at`                                | string                                                                  | false    |              | Dormant at being non-nil indicates a workspace that is dormant. A dormant workspace is no longer accessible must be activated. It is subject to deletion if it breaches the duration of the time_til_ field on its template.                                                                                                                |
| `favorite`                                  | boolean                                                                 | false    |              |                                                                                                                                                                                                                                                                                                                                             |
| `health`                                    | [optimus-ide-collabsdk.WorkspaceHealth](#optimus-ide-collabsdkworkspacehealth)                    | false    |              | Health shows the health of the workspace and information about what is causing an unhealthy status.                                                                                                                                                                                                                                         |
| `id`                                        | string                                                                  | false    |              |                                                                                                                                                                                                                                                                                                                                             |
| `is_prebuild`                               | boolean                                                                 | false    |              | Is prebuild indicates whether the workspace is a prebuilt workspace. Prebuilt workspaces are owned by the prebuilds system user and have specific behavior, such as being managed differently from regular workspaces. Once a prebuilt workspace is claimed by a user, it transitions to a regular workspace, and IsPrebuild returns false. |
| `last_used_at`                              | string                                                                  | false    |              |                                                                                                                                                                                                                                                                                                                                             |
| `latest_app_status`                         | [optimus-ide-collabsdk.WorkspaceAppStatus](#optimus-ide-collabsdkworkspaceappstatus)              | false    |              |                                                                                                                                                                                                                                                                                                                                             |
| `latest_build`                              | [optimus-ide-collabsdk.WorkspaceBuild](#optimus-ide-collabsdkworkspacebuild)                      | false    |              |                                                                                                                                                                                                                                                                                                                                             |
| `name`                                      | string                                                                  | false    |              |                                                                                                                                                                                                                                                                                                                                             |
| `next_start_at`                             | string                                                                  | false    |              |                                                                                                                                                                                                                                                                                                                                             |
| `organization_id`                           | string                                                                  | false    |              |                                                                                                                                                                                                                                                                                                                                             |
| `organization_name`                         | string                                                                  | false    |              |                                                                                                                                                                                                                                                                                                                                             |
| `outdated`                                  | boolean                                                                 | false    |              |                                                                                                                                                                                                                                                                                                                                             |
| `owner_avatar_url`                          | string                                                                  | false    |              |                                                                                                                                                                                                                                                                                                                                             |
| `owner_id`                                  | string                                                                  | false    |              |                                                                                                                                                                                                                                                                                                                                             |
| `owner_name`                                | string                                                                  | false    |              | Owner name is the username of the owner of the workspace.                                                                                                                                                                                                                                                                                   |
| `shared_with`                               | array of [optimus-ide-collabsdk.SharedWorkspaceActor](#optimus-ide-collabsdksharedworkspaceactor) | false    |              |                                                                                                                                                                                                                                                                                                                                             |
| `task_id`                                   | [uuid.NullUUID](#uuidnulluuid)                                          | false    |              | Task ID if set, indicates that the workspace is relevant to the given optimus-ide-collabsdk.Task.                                                                                                                                                                                                                                                        |
| `template_active_version_id`                | string                                                                  | false    |              |                                                                                                                                                                                                                                                                                                                                             |
| `template_allow_user_cancel_workspace_jobs` | boolean                                                                 | false    |              |                                                                                                                                                                                                                                                                                                                                             |
| `template_display_name`                     | string                                                                  | false    |              |                                                                                                                                                                                                                                                                                                                                             |
| `template_icon`                             | string                                                                  | false    |              |                                                                                                                                                                                                                                                                                                                                             |
| `template_id`                               | string                                                                  | false    |              |                                                                                                                                                                                                                                                                                                                                             |
| `template_name`                             | string                                                                  | false    |              |                                                                                                                                                                                                                                                                                                                                             |
| `template_require_active_version`           | boolean                                                                 | false    |              |                                                                                                                                                                                                                                                                                                                                             |
| `template_use_classic_parameter_flow`       | boolean                                                                 | false    |              |                                                                                                                                                                                                                                                                                                                                             |
| `ttl_ms`                                    | integer                                                                 | false    |              |                                                                                                                                                                                                                                                                                                                                             |
| `updated_at`                                | string                                                                  | false    |              |                                                                                                                                                                                                                                                                                                                                             |

#### Enumerated Values

| Property            | Value(s)          |
|---------------------|-------------------|
| `automatic_updates` | `always`, `never` |

## optimus-ide-collabsdk.WorkspaceACL

```json
{
  "group": [
    {
      "avatar_url": "http://example.com",
      "display_name": "string",
      "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
      "members": [
        {
          "avatar_url": "http://example.com",
          "created_at": "2019-08-24T14:15:22Z",
          "email": "user@example.com",
          "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
          "is_service_account": true,
          "last_seen_at": "2019-08-24T14:15:22Z",
          "login_type": "",
          "name": "string",
          "status": "active",
          "theme_preference": "string",
          "updated_at": "2019-08-24T14:15:22Z",
          "username": "string"
        }
      ],
      "name": "string",
      "organization_display_name": "string",
      "organization_id": "7c60d51f-b44e-4682-87d6-449835ea4de6",
      "organization_name": "string",
      "quota_allowance": 0,
      "role": "admin",
      "source": "user",
      "total_member_count": 0
    }
  ],
  "users": [
    {
      "avatar_url": "http://example.com",
      "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
      "name": "string",
      "role": "admin",
      "username": "string"
    }
  ]
}
```

### Properties

| Name    | Type                                                        | Required | Restrictions | Description |
|---------|-------------------------------------------------------------|----------|--------------|-------------|
| `group` | array of [optimus-ide-collabsdk.WorkspaceGroup](#optimus-ide-collabsdkworkspacegroup) | false    |              |             |
| `users` | array of [optimus-ide-collabsdk.WorkspaceUser](#optimus-ide-collabsdkworkspaceuser)   | false    |              |             |

## optimus-ide-collabsdk.WorkspaceAgent

```json
{
  "api_version": "string",
  "apps": [
    {
      "command": "string",
      "display_name": "string",
      "external": true,
      "group": "string",
      "health": "disabled",
      "healthcheck": {
        "interval": 0,
        "threshold": 0,
        "url": "string"
      },
      "hidden": true,
      "icon": "string",
      "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
      "open_in": "slim-window",
      "sharing_level": "owner",
      "slug": "string",
      "statuses": [
        {
          "agent_id": "2b1e3b65-2c04-4fa2-a2d7-467901e98978",
          "app_id": "affd1d10-9538-4fc8-9e0b-4594a28c1335",
          "created_at": "2019-08-24T14:15:22Z",
          "icon": "string",
          "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
          "message": "string",
          "needs_user_attention": true,
          "state": "working",
          "uri": "string",
          "workspace_id": "0967198e-ec7b-4c6b-b4d3-f71244cadbe9"
        }
      ],
      "subdomain": true,
      "subdomain_name": "string",
      "tooltip": "string",
      "url": "string"
    }
  ],
  "architecture": "string",
  "connection_timeout_seconds": 0,
  "created_at": "2019-08-24T14:15:22Z",
  "directory": "string",
  "disconnected_at": "2019-08-24T14:15:22Z",
  "display_apps": [
    "vscode"
  ],
  "environment_variables": {
    "property1": "string",
    "property2": "string"
  },
  "expanded_directory": "string",
  "first_connected_at": "2019-08-24T14:15:22Z",
  "health": {
    "healthy": false,
    "reason": "agent has lost connection"
  },
  "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
  "instance_id": "string",
  "last_connected_at": "2019-08-24T14:15:22Z",
  "latency": {
    "property1": {
      "latency_ms": 0,
      "preferred": true
    },
    "property2": {
      "latency_ms": 0,
      "preferred": true
    }
  },
  "lifecycle_state": "created",
  "log_sources": [
    {
      "created_at": "2019-08-24T14:15:22Z",
      "display_name": "string",
      "icon": "string",
      "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
      "workspace_agent_id": "7ad2e618-fea7-4c1a-b70a-f501566a72f1"
    }
  ],
  "logs_length": 0,
  "logs_overflowed": true,
  "name": "string",
  "operating_system": "string",
  "parent_id": {
    "uuid": "string",
    "valid": true
  },
  "ready_at": "2019-08-24T14:15:22Z",
  "resource_id": "4d5215ed-38bb-48ed-879a-fdb9ca58522f",
  "scripts": [
    {
      "cron": "string",
      "display_name": "string",
      "exit_code": 0,
      "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
      "log_path": "string",
      "log_source_id": "4197ab25-95cf-4b91-9c78-f7f2af5d353a",
      "run_on_start": true,
      "run_on_stop": true,
      "script": "string",
      "start_blocks_login": true,
      "status": "ok",
      "timeout": 0
    }
  ],
  "started_at": "2019-08-24T14:15:22Z",
  "startup_script_behavior": "blocking",
  "status": "connecting",
  "subsystems": [
    "envbox"
  ],
  "troubleshooting_url": "string",
  "updated_at": "2019-08-24T14:15:22Z",
  "version": "string"
}
```

### Properties

| Name                         | Type                                                                                         | Required | Restrictions | Description                                                                                                                                                                  |
|------------------------------|----------------------------------------------------------------------------------------------|----------|--------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `api_version`                | string                                                                                       | false    |              |                                                                                                                                                                              |
| `apps`                       | array of [optimus-ide-collabsdk.WorkspaceApp](#optimus-ide-collabsdkworkspaceapp)                                      | false    |              |                                                                                                                                                                              |
| `architecture`               | string                                                                                       | false    |              |                                                                                                                                                                              |
| `connection_timeout_seconds` | integer                                                                                      | false    |              |                                                                                                                                                                              |
| `created_at`                 | string                                                                                       | false    |              |                                                                                                                                                                              |
| `directory`                  | string                                                                                       | false    |              |                                                                                                                                                                              |
| `disconnected_at`            | string                                                                                       | false    |              |                                                                                                                                                                              |
| `display_apps`               | array of [optimus-ide-collabsdk.DisplayApp](#optimus-ide-collabsdkdisplayapp)                                          | false    |              |                                                                                                                                                                              |
| `environment_variables`      | object                                                                                       | false    |              |                                                                                                                                                                              |
| » `[any property]`           | string                                                                                       | false    |              |                                                                                                                                                                              |
| `expanded_directory`         | string                                                                                       | false    |              |                                                                                                                                                                              |
| `first_connected_at`         | string                                                                                       | false    |              |                                                                                                                                                                              |
| `health`                     | [optimus-ide-collabsdk.WorkspaceAgentHealth](#optimus-ide-collabsdkworkspaceagenthealth)                               | false    |              | Health reports the health of the agent.                                                                                                                                      |
| `id`                         | string                                                                                       | false    |              |                                                                                                                                                                              |
| `instance_id`                | string                                                                                       | false    |              |                                                                                                                                                                              |
| `last_connected_at`          | string                                                                                       | false    |              |                                                                                                                                                                              |
| `latency`                    | object                                                                                       | false    |              | Latency is mapped by region name (e.g. "New York City", "Seattle").                                                                                                          |
| » `[any property]`           | [optimus-ide-collabsdk.DERPRegion](#optimus-ide-collabsdkderpregion)                                                   | false    |              |                                                                                                                                                                              |
| `lifecycle_state`            | [optimus-ide-collabsdk.WorkspaceAgentLifecycle](#optimus-ide-collabsdkworkspaceagentlifecycle)                         | false    |              |                                                                                                                                                                              |
| `log_sources`                | array of [optimus-ide-collabsdk.WorkspaceAgentLogSource](#optimus-ide-collabsdkworkspaceagentlogsource)                | false    |              |                                                                                                                                                                              |
| `logs_length`                | integer                                                                                      | false    |              |                                                                                                                                                                              |
| `logs_overflowed`            | boolean                                                                                      | false    |              |                                                                                                                                                                              |
| `name`                       | string                                                                                       | false    |              |                                                                                                                                                                              |
| `operating_system`           | string                                                                                       | false    |              |                                                                                                                                                                              |
| `parent_id`                  | [uuid.NullUUID](#uuidnulluuid)                                                               | false    |              |                                                                                                                                                                              |
| `ready_at`                   | string                                                                                       | false    |              |                                                                                                                                                                              |
| `resource_id`                | string                                                                                       | false    |              |                                                                                                                                                                              |
| `scripts`                    | array of [optimus-ide-collabsdk.WorkspaceAgentScript](#optimus-ide-collabsdkworkspaceagentscript)                      | false    |              |                                                                                                                                                                              |
| `started_at`                 | string                                                                                       | false    |              |                                                                                                                                                                              |
| `startup_script_behavior`    | [optimus-ide-collabsdk.WorkspaceAgentStartupScriptBehavior](#optimus-ide-collabsdkworkspaceagentstartupscriptbehavior) | false    |              | Startup script behavior is a legacy field that is deprecated in favor of the `optimus-ide-collab_script` resource. It's only referenced by old clients. Deprecated: Remove in the future! |
| `status`                     | [optimus-ide-collabsdk.WorkspaceAgentStatus](#optimus-ide-collabsdkworkspaceagentstatus)                               | false    |              |                                                                                                                                                                              |
| `subsystems`                 | array of [optimus-ide-collabsdk.AgentSubsystem](#optimus-ide-collabsdkagentsubsystem)                                  | false    |              |                                                                                                                                                                              |
| `troubleshooting_url`        | string                                                                                       | false    |              |                                                                                                                                                                              |
| `updated_at`                 | string                                                                                       | false    |              |                                                                                                                                                                              |
| `version`                    | string                                                                                       | false    |              |                                                                                                                                                                              |

## optimus-ide-collabsdk.WorkspaceAgentContainer

```json
{
  "created_at": "2019-08-24T14:15:22Z",
  "id": "string",
  "image": "string",
  "labels": {
    "property1": "string",
    "property2": "string"
  },
  "name": "string",
  "ports": [
    {
      "host_ip": "string",
      "host_port": 0,
      "network": "string",
      "port": 0
    }
  ],
  "running": true,
  "status": "string",
  "volumes": {
    "property1": "string",
    "property2": "string"
  }
}
```

### Properties

| Name               | Type                                                                                  | Required | Restrictions | Description                                                                                                                                |
|--------------------|---------------------------------------------------------------------------------------|----------|--------------|--------------------------------------------------------------------------------------------------------------------------------------------|
| `created_at`       | string                                                                                | false    |              | Created at is the time the container was created.                                                                                          |
| `id`               | string                                                                                | false    |              | ID is the unique identifier of the container.                                                                                              |
| `image`            | string                                                                                | false    |              | Image is the name of the container image.                                                                                                  |
| `labels`           | object                                                                                | false    |              | Labels is a map of key-value pairs of container labels.                                                                                    |
| » `[any property]` | string                                                                                | false    |              |                                                                                                                                            |
| `name`             | string                                                                                | false    |              | Name is the human-readable name of the container.                                                                                          |
| `ports`            | array of [optimus-ide-collabsdk.WorkspaceAgentContainerPort](#optimus-ide-collabsdkworkspaceagentcontainerport) | false    |              | Ports includes ports exposed by the container.                                                                                             |
| `running`          | boolean                                                                               | false    |              | Running is true if the container is currently running.                                                                                     |
| `status`           | string                                                                                | false    |              | Status is the current status of the container. This is somewhat implementation-dependent, but should generally be a human-readable string. |
| `volumes`          | object                                                                                | false    |              | Volumes is a map of "things" mounted into the container. Again, this is somewhat implementation-dependent.                                 |
| » `[any property]` | string                                                                                | false    |              |                                                                                                                                            |

## optimus-ide-collabsdk.WorkspaceAgentContainerPort

```json
{
  "host_ip": "string",
  "host_port": 0,
  "network": "string",
  "port": 0
}
```

### Properties

| Name        | Type    | Required | Restrictions | Description                                                                                                                |
|-------------|---------|----------|--------------|----------------------------------------------------------------------------------------------------------------------------|
| `host_ip`   | string  | false    |              | Host ip is the IP address of the host interface to which the port is bound. Note that this can be an IPv4 or IPv6 address. |
| `host_port` | integer | false    |              | Host port is the port number *outside* the container.                                                                      |
| `network`   | string  | false    |              | Network is the network protocol used by the port (tcp, udp, etc).                                                          |
| `port`      | integer | false    |              | Port is the port number *inside* the container.                                                                            |

## optimus-ide-collabsdk.WorkspaceAgentDevcontainer

```json
{
  "agent": {
    "directory": "string",
    "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
    "name": "string"
  },
  "config_path": "string",
  "container": {
    "created_at": "2019-08-24T14:15:22Z",
    "id": "string",
    "image": "string",
    "labels": {
      "property1": "string",
      "property2": "string"
    },
    "name": "string",
    "ports": [
      {
        "host_ip": "string",
        "host_port": 0,
        "network": "string",
        "port": 0
      }
    ],
    "running": true,
    "status": "string",
    "volumes": {
      "property1": "string",
      "property2": "string"
    }
  },
  "dirty": true,
  "error": "string",
  "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
  "name": "string",
  "status": "running",
  "subagent_id": {
    "uuid": "string",
    "valid": true
  },
  "workspace_folder": "string"
}
```

### Properties

| Name               | Type                                                                                   | Required | Restrictions | Description                |
|--------------------|----------------------------------------------------------------------------------------|----------|--------------|----------------------------|
| `agent`            | [optimus-ide-collabsdk.WorkspaceAgentDevcontainerAgent](#optimus-ide-collabsdkworkspaceagentdevcontaineragent)   | false    |              |                            |
| `config_path`      | string                                                                                 | false    |              |                            |
| `container`        | [optimus-ide-collabsdk.WorkspaceAgentContainer](#optimus-ide-collabsdkworkspaceagentcontainer)                   | false    |              |                            |
| `dirty`            | boolean                                                                                | false    |              |                            |
| `error`            | string                                                                                 | false    |              |                            |
| `id`               | string                                                                                 | false    |              |                            |
| `name`             | string                                                                                 | false    |              |                            |
| `status`           | [optimus-ide-collabsdk.WorkspaceAgentDevcontainerStatus](#optimus-ide-collabsdkworkspaceagentdevcontainerstatus) | false    |              | Additional runtime fields. |
| `subagent_id`      | [uuid.NullUUID](#uuidnulluuid)                                                         | false    |              |                            |
| `workspace_folder` | string                                                                                 | false    |              |                            |

## optimus-ide-collabsdk.WorkspaceAgentDevcontainerAgent

```json
{
  "directory": "string",
  "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
  "name": "string"
}
```

### Properties

| Name        | Type   | Required | Restrictions | Description |
|-------------|--------|----------|--------------|-------------|
| `directory` | string | false    |              |             |
| `id`        | string | false    |              |             |
| `name`      | string | false    |              |             |

## optimus-ide-collabsdk.WorkspaceAgentDevcontainerStatus

```json
"running"
```

### Properties

#### Enumerated Values

| Value(s)                                                          |
|-------------------------------------------------------------------|
| `deleting`, `error`, `running`, `starting`, `stopped`, `stopping` |

## optimus-ide-collabsdk.WorkspaceAgentGitServerMessage

```json
{
  "message": "string",
  "repositories": [
    {
      "branch": "string",
      "remote_origin": "string",
      "removed": true,
      "repo_root": "string",
      "unified_diff": "string"
    }
  ],
  "scanned_at": "2019-08-24T14:15:22Z",
  "type": "changes"
}
```

### Properties

| Name           | Type                                                                                       | Required | Restrictions | Description |
|----------------|--------------------------------------------------------------------------------------------|----------|--------------|-------------|
| `message`      | string                                                                                     | false    |              |             |
| `repositories` | array of [optimus-ide-collabsdk.WorkspaceAgentRepoChanges](#optimus-ide-collabsdkworkspaceagentrepochanges)          | false    |              |             |
| `scanned_at`   | string                                                                                     | false    |              |             |
| `type`         | [optimus-ide-collabsdk.WorkspaceAgentGitServerMessageType](#optimus-ide-collabsdkworkspaceagentgitservermessagetype) | false    |              |             |

## optimus-ide-collabsdk.WorkspaceAgentGitServerMessageType

```json
"changes"
```

### Properties

#### Enumerated Values

| Value(s)           |
|--------------------|
| `changes`, `error` |

## optimus-ide-collabsdk.WorkspaceAgentHealth

```json
{
  "healthy": false,
  "reason": "agent has lost connection"
}
```

### Properties

| Name      | Type    | Required | Restrictions | Description                                                                                   |
|-----------|---------|----------|--------------|-----------------------------------------------------------------------------------------------|
| `healthy` | boolean | false    |              | Healthy is true if the agent is healthy.                                                      |
| `reason`  | string  | false    |              | Reason is a human-readable explanation of the agent's health. It is empty if Healthy is true. |

## optimus-ide-collabsdk.WorkspaceAgentLifecycle

```json
"created"
```

### Properties

#### Enumerated Values

| Value(s)                                                                                                                     |
|------------------------------------------------------------------------------------------------------------------------------|
| `created`, `off`, `ready`, `shutdown_error`, `shutdown_timeout`, `shutting_down`, `start_error`, `start_timeout`, `starting` |

## optimus-ide-collabsdk.WorkspaceAgentListContainersResponse

```json
{
  "containers": [
    {
      "created_at": "2019-08-24T14:15:22Z",
      "id": "string",
      "image": "string",
      "labels": {
        "property1": "string",
        "property2": "string"
      },
      "name": "string",
      "ports": [
        {
          "host_ip": "string",
          "host_port": 0,
          "network": "string",
          "port": 0
        }
      ],
      "running": true,
      "status": "string",
      "volumes": {
        "property1": "string",
        "property2": "string"
      }
    }
  ],
  "devcontainers": [
    {
      "agent": {
        "directory": "string",
        "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
        "name": "string"
      },
      "config_path": "string",
      "container": {
        "created_at": "2019-08-24T14:15:22Z",
        "id": "string",
        "image": "string",
        "labels": {
          "property1": "string",
          "property2": "string"
        },
        "name": "string",
        "ports": [
          {
            "host_ip": "string",
            "host_port": 0,
            "network": "string",
            "port": 0
          }
        ],
        "running": true,
        "status": "string",
        "volumes": {
          "property1": "string",
          "property2": "string"
        }
      },
      "dirty": true,
      "error": "string",
      "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
      "name": "string",
      "status": "running",
      "subagent_id": {
        "uuid": "string",
        "valid": true
      },
      "workspace_folder": "string"
    }
  ],
  "warnings": [
    "string"
  ]
}
```

### Properties

| Name            | Type                                                                                | Required | Restrictions | Description                                                                                                                           |
|-----------------|-------------------------------------------------------------------------------------|----------|--------------|---------------------------------------------------------------------------------------------------------------------------------------|
| `containers`    | array of [optimus-ide-collabsdk.WorkspaceAgentContainer](#optimus-ide-collabsdkworkspaceagentcontainer)       | false    |              | Containers is a list of containers visible to the workspace agent.                                                                    |
| `devcontainers` | array of [optimus-ide-collabsdk.WorkspaceAgentDevcontainer](#optimus-ide-collabsdkworkspaceagentdevcontainer) | false    |              | Devcontainers is a list of devcontainers visible to the workspace agent.                                                              |
| `warnings`      | array of string                                                                     | false    |              | Warnings is a list of warnings that may have occurred during the process of listing containers. This should not include fatal errors. |

## optimus-ide-collabsdk.WorkspaceAgentListeningPort

```json
{
  "network": "string",
  "port": 0,
  "process_name": "string"
}
```

### Properties

| Name           | Type    | Required | Restrictions | Description              |
|----------------|---------|----------|--------------|--------------------------|
| `network`      | string  | false    |              | only "tcp" at the moment |
| `port`         | integer | false    |              |                          |
| `process_name` | string  | false    |              | may be empty             |

## optimus-ide-collabsdk.WorkspaceAgentListeningPortsResponse

```json
{
  "ports": [
    {
      "network": "string",
      "port": 0,
      "process_name": "string"
    }
  ]
}
```

### Properties

| Name    | Type                                                                                  | Required | Restrictions | Description                                                                                                                                                                                                                                            |
|---------|---------------------------------------------------------------------------------------|----------|--------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `ports` | array of [optimus-ide-collabsdk.WorkspaceAgentListeningPort](#optimus-ide-collabsdkworkspaceagentlisteningport) | false    |              | If there are no ports in the list, nothing should be displayed in the UI. There must not be a "no ports available" message or anything similar, as there will always be no ports displayed on platforms where our port detection logic is unsupported. |

## optimus-ide-collabsdk.WorkspaceAgentLog

```json
{
  "created_at": "2019-08-24T14:15:22Z",
  "id": 0,
  "level": "trace",
  "output": "string",
  "source_id": "ae50a35c-df42-4eff-ba26-f8bc28d2af81"
}
```

### Properties

| Name         | Type                                   | Required | Restrictions | Description |
|--------------|----------------------------------------|----------|--------------|-------------|
| `created_at` | string                                 | false    |              |             |
| `id`         | integer                                | false    |              |             |
| `level`      | [optimus-ide-collabsdk.LogLevel](#optimus-ide-collabsdkloglevel) | false    |              |             |
| `output`     | string                                 | false    |              |             |
| `source_id`  | string                                 | false    |              |             |

## optimus-ide-collabsdk.WorkspaceAgentLogSource

```json
{
  "created_at": "2019-08-24T14:15:22Z",
  "display_name": "string",
  "icon": "string",
  "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
  "workspace_agent_id": "7ad2e618-fea7-4c1a-b70a-f501566a72f1"
}
```

### Properties

| Name                 | Type   | Required | Restrictions | Description |
|----------------------|--------|----------|--------------|-------------|
| `created_at`         | string | false    |              |             |
| `display_name`       | string | false    |              |             |
| `icon`               | string | false    |              |             |
| `id`                 | string | false    |              |             |
| `workspace_agent_id` | string | false    |              |             |

## optimus-ide-collabsdk.WorkspaceAgentPortShare

```json
{
  "agent_name": "string",
  "port": 0,
  "protocol": "http",
  "share_level": "owner",
  "workspace_id": "0967198e-ec7b-4c6b-b4d3-f71244cadbe9"
}
```

### Properties

| Name           | Type                                                                                 | Required | Restrictions | Description |
|----------------|--------------------------------------------------------------------------------------|----------|--------------|-------------|
| `agent_name`   | string                                                                               | false    |              |             |
| `port`         | integer                                                                              | false    |              |             |
| `protocol`     | [optimus-ide-collabsdk.WorkspaceAgentPortShareProtocol](#optimus-ide-collabsdkworkspaceagentportshareprotocol) | false    |              |             |
| `share_level`  | [optimus-ide-collabsdk.WorkspaceAgentPortShareLevel](#optimus-ide-collabsdkworkspaceagentportsharelevel)       | false    |              |             |
| `workspace_id` | string                                                                               | false    |              |             |

#### Enumerated Values

| Property      | Value(s)                                           |
|---------------|----------------------------------------------------|
| `protocol`    | `http`, `https`                                    |
| `share_level` | `authenticated`, `organization`, `owner`, `public` |

## optimus-ide-collabsdk.WorkspaceAgentPortShareLevel

```json
"owner"
```

### Properties

#### Enumerated Values

| Value(s)                                           |
|----------------------------------------------------|
| `authenticated`, `organization`, `owner`, `public` |

## optimus-ide-collabsdk.WorkspaceAgentPortShareProtocol

```json
"http"
```

### Properties

#### Enumerated Values

| Value(s)        |
|-----------------|
| `http`, `https` |

## optimus-ide-collabsdk.WorkspaceAgentPortShares

```json
{
  "shares": [
    {
      "agent_name": "string",
      "port": 0,
      "protocol": "http",
      "share_level": "owner",
      "workspace_id": "0967198e-ec7b-4c6b-b4d3-f71244cadbe9"
    }
  ]
}
```

### Properties

| Name     | Type                                                                          | Required | Restrictions | Description |
|----------|-------------------------------------------------------------------------------|----------|--------------|-------------|
| `shares` | array of [optimus-ide-collabsdk.WorkspaceAgentPortShare](#optimus-ide-collabsdkworkspaceagentportshare) | false    |              |             |

## optimus-ide-collabsdk.WorkspaceAgentRepoChanges

```json
{
  "branch": "string",
  "remote_origin": "string",
  "removed": true,
  "repo_root": "string",
  "unified_diff": "string"
}
```

### Properties

| Name            | Type    | Required | Restrictions | Description |
|-----------------|---------|----------|--------------|-------------|
| `branch`        | string  | false    |              |             |
| `remote_origin` | string  | false    |              |             |
| `removed`       | boolean | false    |              |             |
| `repo_root`     | string  | false    |              |             |
| `unified_diff`  | string  | false    |              |             |

## optimus-ide-collabsdk.WorkspaceAgentScript

```json
{
  "cron": "string",
  "display_name": "string",
  "exit_code": 0,
  "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
  "log_path": "string",
  "log_source_id": "4197ab25-95cf-4b91-9c78-f7f2af5d353a",
  "run_on_start": true,
  "run_on_stop": true,
  "script": "string",
  "start_blocks_login": true,
  "status": "ok",
  "timeout": 0
}
```

### Properties

| Name                 | Type                                                                       | Required | Restrictions | Description |
|----------------------|----------------------------------------------------------------------------|----------|--------------|-------------|
| `cron`               | string                                                                     | false    |              |             |
| `display_name`       | string                                                                     | false    |              |             |
| `exit_code`          | integer                                                                    | false    |              |             |
| `id`                 | string                                                                     | false    |              |             |
| `log_path`           | string                                                                     | false    |              |             |
| `log_source_id`      | string                                                                     | false    |              |             |
| `run_on_start`       | boolean                                                                    | false    |              |             |
| `run_on_stop`        | boolean                                                                    | false    |              |             |
| `script`             | string                                                                     | false    |              |             |
| `start_blocks_login` | boolean                                                                    | false    |              |             |
| `status`             | [optimus-ide-collabsdk.WorkspaceAgentScriptStatus](#optimus-ide-collabsdkworkspaceagentscriptstatus) | false    |              |             |
| `timeout`            | integer                                                                    | false    |              |             |

## optimus-ide-collabsdk.WorkspaceAgentScriptStatus

```json
"ok"
```

### Properties

#### Enumerated Values

| Value(s)                                             |
|------------------------------------------------------|
| `exit_failure`, `ok`, `pipes_left_open`, `timed_out` |

## optimus-ide-collabsdk.WorkspaceAgentStartupScriptBehavior

```json
"blocking"
```

### Properties

#### Enumerated Values

| Value(s)                   |
|----------------------------|
| `blocking`, `non-blocking` |

## optimus-ide-collabsdk.WorkspaceAgentStatus

```json
"connecting"
```

### Properties

#### Enumerated Values

| Value(s)                                             |
|------------------------------------------------------|
| `connected`, `connecting`, `disconnected`, `timeout` |

## optimus-ide-collabsdk.WorkspaceApp

```json
{
  "command": "string",
  "display_name": "string",
  "external": true,
  "group": "string",
  "health": "disabled",
  "healthcheck": {
    "interval": 0,
    "threshold": 0,
    "url": "string"
  },
  "hidden": true,
  "icon": "string",
  "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
  "open_in": "slim-window",
  "sharing_level": "owner",
  "slug": "string",
  "statuses": [
    {
      "agent_id": "2b1e3b65-2c04-4fa2-a2d7-467901e98978",
      "app_id": "affd1d10-9538-4fc8-9e0b-4594a28c1335",
      "created_at": "2019-08-24T14:15:22Z",
      "icon": "string",
      "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
      "message": "string",
      "needs_user_attention": true,
      "state": "working",
      "uri": "string",
      "workspace_id": "0967198e-ec7b-4c6b-b4d3-f71244cadbe9"
    }
  ],
  "subdomain": true,
  "subdomain_name": "string",
  "tooltip": "string",
  "url": "string"
}
```

### Properties

| Name             | Type                                                                   | Required | Restrictions | Description                                                                                                                                                                                                                                    |
|------------------|------------------------------------------------------------------------|----------|--------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `command`        | string                                                                 | false    |              |                                                                                                                                                                                                                                                |
| `display_name`   | string                                                                 | false    |              | Display name is a friendly name for the app.                                                                                                                                                                                                   |
| `external`       | boolean                                                                | false    |              | External specifies whether the URL should be opened externally on the client or not.                                                                                                                                                           |
| `group`          | string                                                                 | false    |              |                                                                                                                                                                                                                                                |
| `health`         | [optimus-ide-collabsdk.WorkspaceAppHealth](#optimus-ide-collabsdkworkspaceapphealth)             | false    |              |                                                                                                                                                                                                                                                |
| `healthcheck`    | [optimus-ide-collabsdk.Healthcheck](#optimus-ide-collabsdkhealthcheck)                           | false    |              | Healthcheck specifies the configuration for checking app health.                                                                                                                                                                               |
| `hidden`         | boolean                                                                | false    |              |                                                                                                                                                                                                                                                |
| `icon`           | string                                                                 | false    |              | Icon is a relative path or external URL that specifies an icon to be displayed in the dashboard.                                                                                                                                               |
| `id`             | string                                                                 | false    |              |                                                                                                                                                                                                                                                |
| `open_in`        | [optimus-ide-collabsdk.WorkspaceAppOpenIn](#optimus-ide-collabsdkworkspaceappopenin)             | false    |              |                                                                                                                                                                                                                                                |
| `sharing_level`  | [optimus-ide-collabsdk.WorkspaceAppSharingLevel](#optimus-ide-collabsdkworkspaceappsharinglevel) | false    |              |                                                                                                                                                                                                                                                |
| `slug`           | string                                                                 | false    |              | Slug is a unique identifier within the agent.                                                                                                                                                                                                  |
| `statuses`       | array of [optimus-ide-collabsdk.WorkspaceAppStatus](#optimus-ide-collabsdkworkspaceappstatus)    | false    |              | Statuses is a list of statuses for the app.                                                                                                                                                                                                    |
| `subdomain`      | boolean                                                                | false    |              | Subdomain denotes whether the app should be accessed via a path on the `optimus-ide-collab server` or via a hostname-based dev URL. If this is set to true and there is no app wildcard configured on the server, the app will not be accessible in the UI. |
| `subdomain_name` | string                                                                 | false    |              | Subdomain name is the application domain exposed on the `optimus-ide-collab server`.                                                                                                                                                                        |
| `tooltip`        | string                                                                 | false    |              | Tooltip is an optional markdown supported field that is displayed when hovering over workspace apps in the UI.                                                                                                                                 |
| `url`            | string                                                                 | false    |              | URL is the address being proxied to inside the workspace. If external is specified, this will be opened on the client.                                                                                                                         |

#### Enumerated Values

| Property        | Value(s)                                           |
|-----------------|----------------------------------------------------|
| `sharing_level` | `authenticated`, `organization`, `owner`, `public` |

## optimus-ide-collabsdk.WorkspaceAppHealth

```json
"disabled"
```

### Properties

#### Enumerated Values

| Value(s)                                           |
|----------------------------------------------------|
| `disabled`, `healthy`, `initializing`, `unhealthy` |

## optimus-ide-collabsdk.WorkspaceAppOpenIn

```json
"slim-window"
```

### Properties

#### Enumerated Values

| Value(s)             |
|----------------------|
| `slim-window`, `tab` |

## optimus-ide-collabsdk.WorkspaceAppSharingLevel

```json
"owner"
```

### Properties

#### Enumerated Values

| Value(s)                                           |
|----------------------------------------------------|
| `authenticated`, `organization`, `owner`, `public` |

## optimus-ide-collabsdk.WorkspaceAppStatus

```json
{
  "agent_id": "2b1e3b65-2c04-4fa2-a2d7-467901e98978",
  "app_id": "affd1d10-9538-4fc8-9e0b-4594a28c1335",
  "created_at": "2019-08-24T14:15:22Z",
  "icon": "string",
  "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
  "message": "string",
  "needs_user_attention": true,
  "state": "working",
  "uri": "string",
  "workspace_id": "0967198e-ec7b-4c6b-b4d3-f71244cadbe9"
}
```

### Properties

| Name                   | Type                                                                 | Required | Restrictions | Description                                                                                                                                     |
|------------------------|----------------------------------------------------------------------|----------|--------------|-------------------------------------------------------------------------------------------------------------------------------------------------|
| `agent_id`             | string                                                               | false    |              |                                                                                                                                                 |
| `app_id`               | string                                                               | false    |              |                                                                                                                                                 |
| `created_at`           | string                                                               | false    |              |                                                                                                                                                 |
| `icon`                 | string                                                               | false    |              | Deprecated: This field is unused and will be removed in a future version. Icon is an external URL to an icon that will be rendered in the UI.   |
| `id`                   | string                                                               | false    |              |                                                                                                                                                 |
| `message`              | string                                                               | false    |              |                                                                                                                                                 |
| `needs_user_attention` | boolean                                                              | false    |              | Deprecated: This field is unused and will be removed in a future version. NeedsUserAttention specifies whether the status needs user attention. |
| `state`                | [optimus-ide-collabsdk.WorkspaceAppStatusState](#optimus-ide-collabsdkworkspaceappstatusstate) | false    |              |                                                                                                                                                 |
| `uri`                  | string                                                               | false    |              | Uri is the URI of the resource that the status is for. e.g. https://github.com/org/repo/pull/123 e.g. file:///path/to/file                      |
| `workspace_id`         | string                                                               | false    |              |                                                                                                                                                 |

## optimus-ide-collabsdk.WorkspaceAppStatusState

```json
"working"
```

### Properties

#### Enumerated Values

| Value(s)                                 |
|------------------------------------------|
| `complete`, `failure`, `idle`, `working` |

## optimus-ide-collabsdk.WorkspaceBuild

```json
{
  "build_number": 0,
  "created_at": "2019-08-24T14:15:22Z",
  "daily_cost": 0,
  "deadline": "2019-08-24T14:15:22Z",
  "has_ai_task": true,
  "has_external_agent": true,
  "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
  "initiator_id": "06588898-9a84-4b35-ba8f-f9cbd64946f3",
  "initiator_name": "string",
  "job": {
    "available_workers": [
      "497f6eca-6276-4993-bfeb-53cbbbba6f08"
    ],
    "canceled_at": "2019-08-24T14:15:22Z",
    "completed_at": "2019-08-24T14:15:22Z",
    "created_at": "2019-08-24T14:15:22Z",
    "error": "string",
    "error_code": "REQUIRED_TEMPLATE_VARIABLES",
    "file_id": "8a0cfb4f-ddc9-436d-91bb-75133c583767",
    "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
    "initiator_id": "06588898-9a84-4b35-ba8f-f9cbd64946f3",
    "input": {
      "error": "string",
      "template_version_id": "0ba39c92-1f1b-4c32-aa3e-9925d7713eb1",
      "workspace_build_id": "badaf2eb-96c5-4050-9f1d-db2d39ca5478"
    },
    "logs_overflowed": true,
    "metadata": {
      "template_display_name": "string",
      "template_icon": "string",
      "template_id": "c6d67e98-83ea-49f0-8812-e4abae2b68bc",
      "template_name": "string",
      "template_version_name": "string",
      "workspace_build_transition": "start",
      "workspace_id": "0967198e-ec7b-4c6b-b4d3-f71244cadbe9",
      "workspace_name": "string"
    },
    "organization_id": "7c60d51f-b44e-4682-87d6-449835ea4de6",
    "queue_position": 0,
    "queue_size": 0,
    "started_at": "2019-08-24T14:15:22Z",
    "status": "pending",
    "tags": {
      "property1": "string",
      "property2": "string"
    },
    "type": "template_version_import",
    "worker_id": "ae5fa6f7-c55b-40c1-b40a-b36ac467652b",
    "worker_name": "string"
  },
  "matched_provisioners": {
    "available": 0,
    "count": 0,
    "most_recently_seen": "2019-08-24T14:15:22Z"
  },
  "max_deadline": "2019-08-24T14:15:22Z",
  "reason": "initiator",
  "resources": [
    {
      "agents": [
        {
          "api_version": "string",
          "apps": [
            {
              "command": "string",
              "display_name": "string",
              "external": true,
              "group": "string",
              "health": "disabled",
              "healthcheck": {
                "interval": 0,
                "threshold": 0,
                "url": "string"
              },
              "hidden": true,
              "icon": "string",
              "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
              "open_in": "slim-window",
              "sharing_level": "owner",
              "slug": "string",
              "statuses": [
                {
                  "agent_id": "2b1e3b65-2c04-4fa2-a2d7-467901e98978",
                  "app_id": "affd1d10-9538-4fc8-9e0b-4594a28c1335",
                  "created_at": "2019-08-24T14:15:22Z",
                  "icon": "string",
                  "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
                  "message": "string",
                  "needs_user_attention": true,
                  "state": "working",
                  "uri": "string",
                  "workspace_id": "0967198e-ec7b-4c6b-b4d3-f71244cadbe9"
                }
              ],
              "subdomain": true,
              "subdomain_name": "string",
              "tooltip": "string",
              "url": "string"
            }
          ],
          "architecture": "string",
          "connection_timeout_seconds": 0,
          "created_at": "2019-08-24T14:15:22Z",
          "directory": "string",
          "disconnected_at": "2019-08-24T14:15:22Z",
          "display_apps": [
            "vscode"
          ],
          "environment_variables": {
            "property1": "string",
            "property2": "string"
          },
          "expanded_directory": "string",
          "first_connected_at": "2019-08-24T14:15:22Z",
          "health": {
            "healthy": false,
            "reason": "agent has lost connection"
          },
          "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
          "instance_id": "string",
          "last_connected_at": "2019-08-24T14:15:22Z",
          "latency": {
            "property1": {
              "latency_ms": 0,
              "preferred": true
            },
            "property2": {
              "latency_ms": 0,
              "preferred": true
            }
          },
          "lifecycle_state": "created",
          "log_sources": [
            {
              "created_at": "2019-08-24T14:15:22Z",
              "display_name": "string",
              "icon": "string",
              "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
              "workspace_agent_id": "7ad2e618-fea7-4c1a-b70a-f501566a72f1"
            }
          ],
          "logs_length": 0,
          "logs_overflowed": true,
          "name": "string",
          "operating_system": "string",
          "parent_id": {
            "uuid": "string",
            "valid": true
          },
          "ready_at": "2019-08-24T14:15:22Z",
          "resource_id": "4d5215ed-38bb-48ed-879a-fdb9ca58522f",
          "scripts": [
            {
              "cron": "string",
              "display_name": "string",
              "exit_code": 0,
              "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
              "log_path": "string",
              "log_source_id": "4197ab25-95cf-4b91-9c78-f7f2af5d353a",
              "run_on_start": true,
              "run_on_stop": true,
              "script": "string",
              "start_blocks_login": true,
              "status": "ok",
              "timeout": 0
            }
          ],
          "started_at": "2019-08-24T14:15:22Z",
          "startup_script_behavior": "blocking",
          "status": "connecting",
          "subsystems": [
            "envbox"
          ],
          "troubleshooting_url": "string",
          "updated_at": "2019-08-24T14:15:22Z",
          "version": "string"
        }
      ],
      "created_at": "2019-08-24T14:15:22Z",
      "daily_cost": 0,
      "hide": true,
      "icon": "string",
      "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
      "job_id": "453bd7d7-5355-4d6d-a38e-d9e7eb218c3f",
      "metadata": [
        {
          "key": "string",
          "sensitive": true,
          "value": "string"
        }
      ],
      "name": "string",
      "type": "string",
      "workspace_transition": "start"
    }
  ],
  "status": "pending",
  "template_version_id": "0ba39c92-1f1b-4c32-aa3e-9925d7713eb1",
  "template_version_name": "string",
  "template_version_preset_id": "512a53a7-30da-446e-a1fc-713c630baff1",
  "transition": "start",
  "updated_at": "2019-08-24T14:15:22Z",
  "workspace_id": "0967198e-ec7b-4c6b-b4d3-f71244cadbe9",
  "workspace_name": "string",
  "workspace_owner_avatar_url": "string",
  "workspace_owner_id": "e7078695-5279-4c86-8774-3ac2367a2fc7",
  "workspace_owner_name": "string"
}
```

### Properties

| Name                         | Type                                                              | Required | Restrictions | Description                                                              |
|------------------------------|-------------------------------------------------------------------|----------|--------------|--------------------------------------------------------------------------|
| `build_number`               | integer                                                           | false    |              |                                                                          |
| `created_at`                 | string                                                            | false    |              |                                                                          |
| `daily_cost`                 | integer                                                           | false    |              |                                                                          |
| `deadline`                   | string                                                            | false    |              |                                                                          |
| `has_ai_task`                | boolean                                                           | false    |              | Deprecated: This field has been deprecated in favor of Task WorkspaceID. |
| `has_external_agent`         | boolean                                                           | false    |              |                                                                          |
| `id`                         | string                                                            | false    |              |                                                                          |
| `initiator_id`               | string                                                            | false    |              |                                                                          |
| `initiator_name`             | string                                                            | false    |              |                                                                          |
| `job`                        | [optimus-ide-collabsdk.ProvisionerJob](#optimus-ide-collabsdkprovisionerjob)                | false    |              |                                                                          |
| `matched_provisioners`       | [optimus-ide-collabsdk.MatchedProvisioners](#optimus-ide-collabsdkmatchedprovisioners)      | false    |              |                                                                          |
| `max_deadline`               | string                                                            | false    |              |                                                                          |
| `reason`                     | [optimus-ide-collabsdk.BuildReason](#optimus-ide-collabsdkbuildreason)                      | false    |              |                                                                          |
| `resources`                  | array of [optimus-ide-collabsdk.WorkspaceResource](#optimus-ide-collabsdkworkspaceresource) | false    |              |                                                                          |
| `status`                     | [optimus-ide-collabsdk.WorkspaceStatus](#optimus-ide-collabsdkworkspacestatus)              | false    |              |                                                                          |
| `template_version_id`        | string                                                            | false    |              |                                                                          |
| `template_version_name`      | string                                                            | false    |              |                                                                          |
| `template_version_preset_id` | string                                                            | false    |              |                                                                          |
| `transition`                 | [optimus-ide-collabsdk.WorkspaceTransition](#optimus-ide-collabsdkworkspacetransition)      | false    |              |                                                                          |
| `updated_at`                 | string                                                            | false    |              |                                                                          |
| `workspace_id`               | string                                                            | false    |              |                                                                          |
| `workspace_name`             | string                                                            | false    |              |                                                                          |
| `workspace_owner_avatar_url` | string                                                            | false    |              |                                                                          |
| `workspace_owner_id`         | string                                                            | false    |              |                                                                          |
| `workspace_owner_name`       | string                                                            | false    |              | Workspace owner name is the username of the owner of the workspace.      |

#### Enumerated Values

| Property     | Value(s)                                                                                                          |
|--------------|-------------------------------------------------------------------------------------------------------------------|
| `reason`     | `autostart`, `autostop`, `initiator`                                                                              |
| `status`     | `canceled`, `canceling`, `deleted`, `deleting`, `failed`, `pending`, `running`, `starting`, `stopped`, `stopping` |
| `transition` | `delete`, `start`, `stop`                                                                                         |

## optimus-ide-collabsdk.WorkspaceBuildParameter

```json
{
  "name": "string",
  "value": "string"
}
```

### Properties

| Name    | Type   | Required | Restrictions | Description |
|---------|--------|----------|--------------|-------------|
| `name`  | string | false    |              |             |
| `value` | string | false    |              |             |

## optimus-ide-collabsdk.WorkspaceBuildTimings

```json
{
  "agent_connection_timings": [
    {
      "ended_at": "2019-08-24T14:15:22Z",
      "stage": "init",
      "started_at": "2019-08-24T14:15:22Z",
      "workspace_agent_id": "string",
      "workspace_agent_name": "string"
    }
  ],
  "agent_script_timings": [
    {
      "display_name": "string",
      "ended_at": "2019-08-24T14:15:22Z",
      "exit_code": 0,
      "stage": "init",
      "started_at": "2019-08-24T14:15:22Z",
      "status": "string",
      "workspace_agent_id": "string",
      "workspace_agent_name": "string"
    }
  ],
  "provisioner_timings": [
    {
      "action": "string",
      "ended_at": "2019-08-24T14:15:22Z",
      "job_id": "453bd7d7-5355-4d6d-a38e-d9e7eb218c3f",
      "resource": "string",
      "source": "string",
      "stage": "init",
      "started_at": "2019-08-24T14:15:22Z"
    }
  ]
}
```

### Properties

| Name                       | Type                                                                      | Required | Restrictions | Description                                                                                                      |
|----------------------------|---------------------------------------------------------------------------|----------|--------------|------------------------------------------------------------------------------------------------------------------|
| `agent_connection_timings` | array of [optimus-ide-collabsdk.AgentConnectionTiming](#optimus-ide-collabsdkagentconnectiontiming) | false    |              |                                                                                                                  |
| `agent_script_timings`     | array of [optimus-ide-collabsdk.AgentScriptTiming](#optimus-ide-collabsdkagentscripttiming)         | false    |              | Agent script timings Consolidate agent-related timing metrics into a single struct when updating the API version |
| `provisioner_timings`      | array of [optimus-ide-collabsdk.ProvisionerTiming](#optimus-ide-collabsdkprovisionertiming)         | false    |              |                                                                                                                  |

## optimus-ide-collabsdk.WorkspaceConnectionLatencyMS

```json
{
  "p50": 0,
  "p95": 0
}
```

### Properties

| Name  | Type   | Required | Restrictions | Description |
|-------|--------|----------|--------------|-------------|
| `p50` | number | false    |              |             |
| `p95` | number | false    |              |             |

## optimus-ide-collabsdk.WorkspaceDeploymentStats

```json
{
  "building": 0,
  "connection_latency_ms": {
    "p50": 0,
    "p95": 0
  },
  "failed": 0,
  "pending": 0,
  "running": 0,
  "rx_bytes": 0,
  "stopped": 0,
  "tx_bytes": 0
}
```

### Properties

| Name                    | Type                                                                           | Required | Restrictions | Description |
|-------------------------|--------------------------------------------------------------------------------|----------|--------------|-------------|
| `building`              | integer                                                                        | false    |              |             |
| `connection_latency_ms` | [optimus-ide-collabsdk.WorkspaceConnectionLatencyMS](#optimus-ide-collabsdkworkspaceconnectionlatencyms) | false    |              |             |
| `failed`                | integer                                                                        | false    |              |             |
| `pending`               | integer                                                                        | false    |              |             |
| `running`               | integer                                                                        | false    |              |             |
| `rx_bytes`              | integer                                                                        | false    |              |             |
| `stopped`               | integer                                                                        | false    |              |             |
| `tx_bytes`              | integer                                                                        | false    |              |             |

## optimus-ide-collabsdk.WorkspaceGroup

```json
{
  "avatar_url": "http://example.com",
  "display_name": "string",
  "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
  "members": [
    {
      "avatar_url": "http://example.com",
      "created_at": "2019-08-24T14:15:22Z",
      "email": "user@example.com",
      "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
      "is_service_account": true,
      "last_seen_at": "2019-08-24T14:15:22Z",
      "login_type": "",
      "name": "string",
      "status": "active",
      "theme_preference": "string",
      "updated_at": "2019-08-24T14:15:22Z",
      "username": "string"
    }
  ],
  "name": "string",
  "organization_display_name": "string",
  "organization_id": "7c60d51f-b44e-4682-87d6-449835ea4de6",
  "organization_name": "string",
  "quota_allowance": 0,
  "role": "admin",
  "source": "user",
  "total_member_count": 0
}
```

### Properties

| Name                        | Type                                                  | Required | Restrictions | Description                                                                                                                                                           |
|-----------------------------|-------------------------------------------------------|----------|--------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `avatar_url`                | string                                                | false    |              |                                                                                                                                                                       |
| `display_name`              | string                                                | false    |              |                                                                                                                                                                       |
| `id`                        | string                                                | false    |              |                                                                                                                                                                       |
| `members`                   | array of [optimus-ide-collabsdk.ReducedUser](#optimus-ide-collabsdkreduceduser) | false    |              |                                                                                                                                                                       |
| `name`                      | string                                                | false    |              |                                                                                                                                                                       |
| `organization_display_name` | string                                                | false    |              |                                                                                                                                                                       |
| `organization_id`           | string                                                | false    |              |                                                                                                                                                                       |
| `organization_name`         | string                                                | false    |              |                                                                                                                                                                       |
| `quota_allowance`           | integer                                               | false    |              |                                                                                                                                                                       |
| `role`                      | [optimus-ide-collabsdk.WorkspaceRole](#optimus-ide-collabsdkworkspacerole)      | false    |              |                                                                                                                                                                       |
| `source`                    | [optimus-ide-collabsdk.GroupSource](#optimus-ide-collabsdkgroupsource)          | false    |              |                                                                                                                                                                       |
| `total_member_count`        | integer                                               | false    |              | How many members are in this group. Shows the total count, even if the user is not authorized to read group member details. May be greater than `len(Group.Members)`. |

#### Enumerated Values

| Property | Value(s)       |
|----------|----------------|
| `role`   | `admin`, `use` |

## optimus-ide-collabsdk.WorkspaceHealth

```json
{
  "failing_agents": [
    "497f6eca-6276-4993-bfeb-53cbbbba6f08"
  ],
  "healthy": false
}
```

### Properties

| Name             | Type            | Required | Restrictions | Description                                                          |
|------------------|-----------------|----------|--------------|----------------------------------------------------------------------|
| `failing_agents` | array of string | false    |              | Failing agents lists the IDs of the agents that are failing, if any. |
| `healthy`        | boolean         | false    |              | Healthy is true if the workspace is healthy.                         |

## optimus-ide-collabsdk.WorkspaceProxy

```json
{
  "created_at": "2019-08-24T14:15:22Z",
  "deleted": true,
  "derp_enabled": true,
  "derp_only": true,
  "display_name": "string",
  "healthy": true,
  "icon_url": "string",
  "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
  "name": "string",
  "path_app_url": "string",
  "status": {
    "checked_at": "2019-08-24T14:15:22Z",
    "report": {
      "errors": [
        "string"
      ],
      "warnings": [
        "string"
      ]
    },
    "status": "ok"
  },
  "updated_at": "2019-08-24T14:15:22Z",
  "version": "string",
  "wildcard_hostname": "string"
}
```

### Properties

| Name                | Type                                                           | Required | Restrictions | Description                                                                                                                                                                       |
|---------------------|----------------------------------------------------------------|----------|--------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `created_at`        | string                                                         | false    |              |                                                                                                                                                                                   |
| `deleted`           | boolean                                                        | false    |              |                                                                                                                                                                                   |
| `derp_enabled`      | boolean                                                        | false    |              |                                                                                                                                                                                   |
| `derp_only`         | boolean                                                        | false    |              |                                                                                                                                                                                   |
| `display_name`      | string                                                         | false    |              |                                                                                                                                                                                   |
| `healthy`           | boolean                                                        | false    |              |                                                                                                                                                                                   |
| `icon_url`          | string                                                         | false    |              |                                                                                                                                                                                   |
| `id`                | string                                                         | false    |              |                                                                                                                                                                                   |
| `name`              | string                                                         | false    |              |                                                                                                                                                                                   |
| `path_app_url`      | string                                                         | false    |              | Path app URL is the URL to the base path for path apps. Optional unless wildcard_hostname is set. E.g. https://us.example.com                                                     |
| `status`            | [optimus-ide-collabsdk.WorkspaceProxyStatus](#optimus-ide-collabsdkworkspaceproxystatus) | false    |              | Status is the latest status check of the proxy. This will be empty for deleted proxies. This value can be used to determine if a workspace proxy is healthy and ready to use.     |
| `updated_at`        | string                                                         | false    |              |                                                                                                                                                                                   |
| `version`           | string                                                         | false    |              |                                                                                                                                                                                   |
| `wildcard_hostname` | string                                                         | false    |              | Wildcard hostname is the wildcard hostname for subdomain apps. E.g. *.us.example.com E.g.*--suffix.au.example.com Optional. Does not need to be on the same domain as PathAppURL. |

## optimus-ide-collabsdk.WorkspaceProxyStatus

```json
{
  "checked_at": "2019-08-24T14:15:22Z",
  "report": {
    "errors": [
      "string"
    ],
    "warnings": [
      "string"
    ]
  },
  "status": "ok"
}
```

### Properties

| Name         | Type                                                     | Required | Restrictions | Description                                                               |
|--------------|----------------------------------------------------------|----------|--------------|---------------------------------------------------------------------------|
| `checked_at` | string                                                   | false    |              |                                                                           |
| `report`     | [optimus-ide-collabsdk.ProxyHealthReport](#optimus-ide-collabsdkproxyhealthreport) | false    |              | Report provides more information about the health of the workspace proxy. |
| `status`     | [optimus-ide-collabsdk.ProxyHealthStatus](#optimus-ide-collabsdkproxyhealthstatus) | false    |              |                                                                           |

## optimus-ide-collabsdk.WorkspaceQuota

```json
{
  "budget": 0,
  "credits_consumed": 0
}
```

### Properties

| Name               | Type    | Required | Restrictions | Description |
|--------------------|---------|----------|--------------|-------------|
| `budget`           | integer | false    |              |             |
| `credits_consumed` | integer | false    |              |             |

## optimus-ide-collabsdk.WorkspaceResource

```json
{
  "agents": [
    {
      "api_version": "string",
      "apps": [
        {
          "command": "string",
          "display_name": "string",
          "external": true,
          "group": "string",
          "health": "disabled",
          "healthcheck": {
            "interval": 0,
            "threshold": 0,
            "url": "string"
          },
          "hidden": true,
          "icon": "string",
          "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
          "open_in": "slim-window",
          "sharing_level": "owner",
          "slug": "string",
          "statuses": [
            {
              "agent_id": "2b1e3b65-2c04-4fa2-a2d7-467901e98978",
              "app_id": "affd1d10-9538-4fc8-9e0b-4594a28c1335",
              "created_at": "2019-08-24T14:15:22Z",
              "icon": "string",
              "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
              "message": "string",
              "needs_user_attention": true,
              "state": "working",
              "uri": "string",
              "workspace_id": "0967198e-ec7b-4c6b-b4d3-f71244cadbe9"
            }
          ],
          "subdomain": true,
          "subdomain_name": "string",
          "tooltip": "string",
          "url": "string"
        }
      ],
      "architecture": "string",
      "connection_timeout_seconds": 0,
      "created_at": "2019-08-24T14:15:22Z",
      "directory": "string",
      "disconnected_at": "2019-08-24T14:15:22Z",
      "display_apps": [
        "vscode"
      ],
      "environment_variables": {
        "property1": "string",
        "property2": "string"
      },
      "expanded_directory": "string",
      "first_connected_at": "2019-08-24T14:15:22Z",
      "health": {
        "healthy": false,
        "reason": "agent has lost connection"
      },
      "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
      "instance_id": "string",
      "last_connected_at": "2019-08-24T14:15:22Z",
      "latency": {
        "property1": {
          "latency_ms": 0,
          "preferred": true
        },
        "property2": {
          "latency_ms": 0,
          "preferred": true
        }
      },
      "lifecycle_state": "created",
      "log_sources": [
        {
          "created_at": "2019-08-24T14:15:22Z",
          "display_name": "string",
          "icon": "string",
          "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
          "workspace_agent_id": "7ad2e618-fea7-4c1a-b70a-f501566a72f1"
        }
      ],
      "logs_length": 0,
      "logs_overflowed": true,
      "name": "string",
      "operating_system": "string",
      "parent_id": {
        "uuid": "string",
        "valid": true
      },
      "ready_at": "2019-08-24T14:15:22Z",
      "resource_id": "4d5215ed-38bb-48ed-879a-fdb9ca58522f",
      "scripts": [
        {
          "cron": "string",
          "display_name": "string",
          "exit_code": 0,
          "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
          "log_path": "string",
          "log_source_id": "4197ab25-95cf-4b91-9c78-f7f2af5d353a",
          "run_on_start": true,
          "run_on_stop": true,
          "script": "string",
          "start_blocks_login": true,
          "status": "ok",
          "timeout": 0
        }
      ],
      "started_at": "2019-08-24T14:15:22Z",
      "startup_script_behavior": "blocking",
      "status": "connecting",
      "subsystems": [
        "envbox"
      ],
      "troubleshooting_url": "string",
      "updated_at": "2019-08-24T14:15:22Z",
      "version": "string"
    }
  ],
  "created_at": "2019-08-24T14:15:22Z",
  "daily_cost": 0,
  "hide": true,
  "icon": "string",
  "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
  "job_id": "453bd7d7-5355-4d6d-a38e-d9e7eb218c3f",
  "metadata": [
    {
      "key": "string",
      "sensitive": true,
      "value": "string"
    }
  ],
  "name": "string",
  "type": "string",
  "workspace_transition": "start"
}
```

### Properties

| Name                   | Type                                                                              | Required | Restrictions | Description |
|------------------------|-----------------------------------------------------------------------------------|----------|--------------|-------------|
| `agents`               | array of [optimus-ide-collabsdk.WorkspaceAgent](#optimus-ide-collabsdkworkspaceagent)                       | false    |              |             |
| `created_at`           | string                                                                            | false    |              |             |
| `daily_cost`           | integer                                                                           | false    |              |             |
| `hide`                 | boolean                                                                           | false    |              |             |
| `icon`                 | string                                                                            | false    |              |             |
| `id`                   | string                                                                            | false    |              |             |
| `job_id`               | string                                                                            | false    |              |             |
| `metadata`             | array of [optimus-ide-collabsdk.WorkspaceResourceMetadata](#optimus-ide-collabsdkworkspaceresourcemetadata) | false    |              |             |
| `name`                 | string                                                                            | false    |              |             |
| `type`                 | string                                                                            | false    |              |             |
| `workspace_transition` | [optimus-ide-collabsdk.WorkspaceTransition](#optimus-ide-collabsdkworkspacetransition)                      | false    |              |             |

#### Enumerated Values

| Property               | Value(s)                  |
|------------------------|---------------------------|
| `workspace_transition` | `delete`, `start`, `stop` |

## optimus-ide-collabsdk.WorkspaceResourceMetadata

```json
{
  "key": "string",
  "sensitive": true,
  "value": "string"
}
```

### Properties

| Name        | Type    | Required | Restrictions | Description |
|-------------|---------|----------|--------------|-------------|
| `key`       | string  | false    |              |             |
| `sensitive` | boolean | false    |              |             |
| `value`     | string  | false    |              |             |

## optimus-ide-collabsdk.WorkspaceRole

```json
"admin"
```

### Properties

#### Enumerated Values

| Value(s)           |
|--------------------|
| ``, `admin`, `use` |

## optimus-ide-collabsdk.WorkspaceSharingSettings

```json
{
  "shareable_workspace_owners": "none",
  "sharing_disabled": true,
  "sharing_globally_disabled": true
}
```

### Properties

| Name                         | Type                                                                   | Required | Restrictions | Description                                                                                                                     |
|------------------------------|------------------------------------------------------------------------|----------|--------------|---------------------------------------------------------------------------------------------------------------------------------|
| `shareable_workspace_owners` | [optimus-ide-collabsdk.ShareableWorkspaceOwners](#optimus-ide-collabsdkshareableworkspaceowners) | false    |              | Shareable workspace owners controls whose workspaces can be shared within the organization.                                     |
| `sharing_disabled`           | boolean                                                                | false    |              | Sharing disabled is deprecated and left for backward compatibility purposes. Deprecated: use `ShareableWorkspaceOwners` instead |
| `sharing_globally_disabled`  | boolean                                                                | false    |              | Sharing globally disabled is true if sharing has been disabled for this organization because of a deployment-wide setting.      |

#### Enumerated Values

| Property                     | Value(s)                               |
|------------------------------|----------------------------------------|
| `shareable_workspace_owners` | `everyone`, `none`, `service_accounts` |

## optimus-ide-collabsdk.WorkspaceStatus

```json
"pending"
```

### Properties

#### Enumerated Values

| Value(s)                                                                                                          |
|-------------------------------------------------------------------------------------------------------------------|
| `canceled`, `canceling`, `deleted`, `deleting`, `failed`, `pending`, `running`, `starting`, `stopped`, `stopping` |

## optimus-ide-collabsdk.WorkspaceTransition

```json
"start"
```

### Properties

#### Enumerated Values

| Value(s)                  |
|---------------------------|
| `delete`, `start`, `stop` |

## optimus-ide-collabsdk.WorkspaceUser

```json
{
  "avatar_url": "http://example.com",
  "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
  "name": "string",
  "role": "admin",
  "username": "string"
}
```

### Properties

| Name         | Type                                             | Required | Restrictions | Description |
|--------------|--------------------------------------------------|----------|--------------|-------------|
| `avatar_url` | string                                           | false    |              |             |
| `id`         | string                                           | true     |              |             |
| `name`       | string                                           | false    |              |             |
| `role`       | [optimus-ide-collabsdk.WorkspaceRole](#optimus-ide-collabsdkworkspacerole) | false    |              |             |
| `username`   | string                                           | true     |              |             |

#### Enumerated Values

| Property | Value(s)       |
|----------|----------------|
| `role`   | `admin`, `use` |

## optimus-ide-collabsdk.WorkspacesResponse

```json
{
  "count": 0,
  "workspaces": [
    {
      "allow_renames": true,
      "automatic_updates": "always",
      "autostart_schedule": "string",
      "created_at": "2019-08-24T14:15:22Z",
      "deleting_at": "2019-08-24T14:15:22Z",
      "dormant_at": "2019-08-24T14:15:22Z",
      "favorite": true,
      "health": {
        "failing_agents": [
          "497f6eca-6276-4993-bfeb-53cbbbba6f08"
        ],
        "healthy": false
      },
      "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
      "is_prebuild": true,
      "last_used_at": "2019-08-24T14:15:22Z",
      "latest_app_status": {
        "agent_id": "2b1e3b65-2c04-4fa2-a2d7-467901e98978",
        "app_id": "affd1d10-9538-4fc8-9e0b-4594a28c1335",
        "created_at": "2019-08-24T14:15:22Z",
        "icon": "string",
        "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
        "message": "string",
        "needs_user_attention": true,
        "state": "working",
        "uri": "string",
        "workspace_id": "0967198e-ec7b-4c6b-b4d3-f71244cadbe9"
      },
      "latest_build": {
        "build_number": 0,
        "created_at": "2019-08-24T14:15:22Z",
        "daily_cost": 0,
        "deadline": "2019-08-24T14:15:22Z",
        "has_ai_task": true,
        "has_external_agent": true,
        "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
        "initiator_id": "06588898-9a84-4b35-ba8f-f9cbd64946f3",
        "initiator_name": "string",
        "job": {
          "available_workers": [
            "497f6eca-6276-4993-bfeb-53cbbbba6f08"
          ],
          "canceled_at": "2019-08-24T14:15:22Z",
          "completed_at": "2019-08-24T14:15:22Z",
          "created_at": "2019-08-24T14:15:22Z",
          "error": "string",
          "error_code": "REQUIRED_TEMPLATE_VARIABLES",
          "file_id": "8a0cfb4f-ddc9-436d-91bb-75133c583767",
          "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
          "initiator_id": "06588898-9a84-4b35-ba8f-f9cbd64946f3",
          "input": {
            "error": "string",
            "template_version_id": "0ba39c92-1f1b-4c32-aa3e-9925d7713eb1",
            "workspace_build_id": "badaf2eb-96c5-4050-9f1d-db2d39ca5478"
          },
          "logs_overflowed": true,
          "metadata": {
            "template_display_name": "string",
            "template_icon": "string",
            "template_id": "c6d67e98-83ea-49f0-8812-e4abae2b68bc",
            "template_name": "string",
            "template_version_name": "string",
            "workspace_build_transition": "start",
            "workspace_id": "0967198e-ec7b-4c6b-b4d3-f71244cadbe9",
            "workspace_name": "string"
          },
          "organization_id": "7c60d51f-b44e-4682-87d6-449835ea4de6",
          "queue_position": 0,
          "queue_size": 0,
          "started_at": "2019-08-24T14:15:22Z",
          "status": "pending",
          "tags": {
            "property1": "string",
            "property2": "string"
          },
          "type": "template_version_import",
          "worker_id": "ae5fa6f7-c55b-40c1-b40a-b36ac467652b",
          "worker_name": "string"
        },
        "matched_provisioners": {
          "available": 0,
          "count": 0,
          "most_recently_seen": "2019-08-24T14:15:22Z"
        },
        "max_deadline": "2019-08-24T14:15:22Z",
        "reason": "initiator",
        "resources": [
          {
            "agents": [
              {
                "api_version": "string",
                "apps": [
                  {
                    "command": "string",
                    "display_name": "string",
                    "external": true,
                    "group": "string",
                    "health": "disabled",
                    "healthcheck": {},
                    "hidden": true,
                    "icon": "string",
                    "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
                    "open_in": "slim-window",
                    "sharing_level": "owner",
                    "slug": "string",
                    "statuses": [],
                    "subdomain": true,
                    "subdomain_name": "string",
                    "tooltip": "string",
                    "url": "string"
                  }
                ],
                "architecture": "string",
                "connection_timeout_seconds": 0,
                "created_at": "2019-08-24T14:15:22Z",
                "directory": "string",
                "disconnected_at": "2019-08-24T14:15:22Z",
                "display_apps": [
                  "vscode"
                ],
                "environment_variables": {
                  "property1": "string",
                  "property2": "string"
                },
                "expanded_directory": "string",
                "first_connected_at": "2019-08-24T14:15:22Z",
                "health": {
                  "healthy": false,
                  "reason": "agent has lost connection"
                },
                "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
                "instance_id": "string",
                "last_connected_at": "2019-08-24T14:15:22Z",
                "latency": {
                  "property1": {
                    "latency_ms": 0,
                    "preferred": true
                  },
                  "property2": {
                    "latency_ms": 0,
                    "preferred": true
                  }
                },
                "lifecycle_state": "created",
                "log_sources": [
                  {
                    "created_at": "2019-08-24T14:15:22Z",
                    "display_name": "string",
                    "icon": "string",
                    "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
                    "workspace_agent_id": "7ad2e618-fea7-4c1a-b70a-f501566a72f1"
                  }
                ],
                "logs_length": 0,
                "logs_overflowed": true,
                "name": "string",
                "operating_system": "string",
                "parent_id": {
                  "uuid": "string",
                  "valid": true
                },
                "ready_at": "2019-08-24T14:15:22Z",
                "resource_id": "4d5215ed-38bb-48ed-879a-fdb9ca58522f",
                "scripts": [
                  {
                    "cron": "string",
                    "display_name": "string",
                    "exit_code": 0,
                    "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
                    "log_path": "string",
                    "log_source_id": "4197ab25-95cf-4b91-9c78-f7f2af5d353a",
                    "run_on_start": true,
                    "run_on_stop": true,
                    "script": "string",
                    "start_blocks_login": true,
                    "status": "ok",
                    "timeout": 0
                  }
                ],
                "started_at": "2019-08-24T14:15:22Z",
                "startup_script_behavior": "blocking",
                "status": "connecting",
                "subsystems": [
                  "envbox"
                ],
                "troubleshooting_url": "string",
                "updated_at": "2019-08-24T14:15:22Z",
                "version": "string"
              }
            ],
            "created_at": "2019-08-24T14:15:22Z",
            "daily_cost": 0,
            "hide": true,
            "icon": "string",
            "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
            "job_id": "453bd7d7-5355-4d6d-a38e-d9e7eb218c3f",
            "metadata": [
              {
                "key": "string",
                "sensitive": true,
                "value": "string"
              }
            ],
            "name": "string",
            "type": "string",
            "workspace_transition": "start"
          }
        ],
        "status": "pending",
        "template_version_id": "0ba39c92-1f1b-4c32-aa3e-9925d7713eb1",
        "template_version_name": "string",
        "template_version_preset_id": "512a53a7-30da-446e-a1fc-713c630baff1",
        "transition": "start",
        "updated_at": "2019-08-24T14:15:22Z",
        "workspace_id": "0967198e-ec7b-4c6b-b4d3-f71244cadbe9",
        "workspace_name": "string",
        "workspace_owner_avatar_url": "string",
        "workspace_owner_id": "e7078695-5279-4c86-8774-3ac2367a2fc7",
        "workspace_owner_name": "string"
      },
      "name": "string",
      "next_start_at": "2019-08-24T14:15:22Z",
      "organization_id": "7c60d51f-b44e-4682-87d6-449835ea4de6",
      "organization_name": "string",
      "outdated": true,
      "owner_avatar_url": "string",
      "owner_id": "8826ee2e-7933-4665-aef2-2393f84a0d05",
      "owner_name": "string",
      "shared_with": [
        {
          "actor_type": "group",
          "avatar_url": "http://example.com",
          "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
          "name": "string",
          "roles": [
            "admin"
          ]
        }
      ],
      "task_id": {
        "uuid": "string",
        "valid": true
      },
      "template_active_version_id": "b0da9c29-67d8-4c87-888c-bafe356f7f3c",
      "template_allow_user_cancel_workspace_jobs": true,
      "template_display_name": "string",
      "template_icon": "string",
      "template_id": "c6d67e98-83ea-49f0-8812-e4abae2b68bc",
      "template_name": "string",
      "template_require_active_version": true,
      "template_use_classic_parameter_flow": true,
      "ttl_ms": 0,
      "updated_at": "2019-08-24T14:15:22Z"
    }
  ]
}
```

### Properties

| Name         | Type                                              | Required | Restrictions | Description |
|--------------|---------------------------------------------------|----------|--------------|-------------|
| `count`      | integer                                           | false    |              |             |
| `workspaces` | array of [optimus-ide-collabsdk.Workspace](#optimus-ide-collabsdkworkspace) | false    |              |             |

## derp.BytesSentRecv

```json
{
  "key": {},
  "recv": 0,
  "sent": 0
}
```

### Properties

| Name   | Type                             | Required | Restrictions | Description                                                          |
|--------|----------------------------------|----------|--------------|----------------------------------------------------------------------|
| `key`  | [key.NodePublic](#keynodepublic) | false    |              | Key is the public key of the client which sent/received these bytes. |
| `recv` | integer                          | false    |              |                                                                      |
| `sent` | integer                          | false    |              |                                                                      |

## derp.ServerInfoMessage

```json
{
  "tokenBucketBytesBurst": 0,
  "tokenBucketBytesPerSecond": 0
}
```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|`tokenBucketBytesBurst`|integer|false||Tokenbucketbytesburst is how many bytes the server will allow to burst, temporarily violating TokenBucketBytesPerSecond.
Zero means unspecified. There might be a limit, but the client need not try to respect it.|
|`tokenBucketBytesPerSecond`|integer|false||Tokenbucketbytespersecond is how many bytes per second the server says it will accept, including all framing bytes.
Zero means unspecified. There might be a limit, but the client need not try to respect it.|

## health.Code

```json
"EUNKNOWN"
```

### Properties

#### Enumerated Values

| Value(s)                                                                                                                                                                               |
|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `EACS01`, `EACS02`, `EACS03`, `EACS04`, `EDB01`, `EDB02`, `EDERP01`, `EDERP02`, `EDERP03`, `EPD01`, `EPD02`, `EPD03`, `EUNKNOWN`, `EWP01`, `EWP02`, `EWP04`, `EWS01`, `EWS02`, `EWS03` |

## health.Message

```json
{
  "code": "EUNKNOWN",
  "message": "string"
}
```

### Properties

| Name      | Type                       | Required | Restrictions | Description |
|-----------|----------------------------|----------|--------------|-------------|
| `code`    | [health.Code](#healthcode) | false    |              |             |
| `message` | string                     | false    |              |             |

## health.Severity

```json
"ok"
```

### Properties

#### Enumerated Values

| Value(s)                 |
|--------------------------|
| `error`, `ok`, `warning` |

## healthsdk.AccessURLReport

```json
{
  "access_url": "string",
  "dismissed": true,
  "error": "string",
  "healthy": true,
  "healthz_response": "string",
  "reachable": true,
  "severity": "ok",
  "status_code": 0,
  "warnings": [
    {
      "code": "EUNKNOWN",
      "message": "string"
    }
  ]
}
```

### Properties

| Name               | Type                                      | Required | Restrictions | Description                                                                                 |
|--------------------|-------------------------------------------|----------|--------------|---------------------------------------------------------------------------------------------|
| `access_url`       | string                                    | false    |              |                                                                                             |
| `dismissed`        | boolean                                   | false    |              |                                                                                             |
| `error`            | string                                    | false    |              |                                                                                             |
| `healthy`          | boolean                                   | false    |              | Healthy is deprecated and left for backward compatibility purposes, use `Severity` instead. |
| `healthz_response` | string                                    | false    |              |                                                                                             |
| `reachable`        | boolean                                   | false    |              |                                                                                             |
| `severity`         | [health.Severity](#healthseverity)        | false    |              |                                                                                             |
| `status_code`      | integer                                   | false    |              |                                                                                             |
| `warnings`         | array of [health.Message](#healthmessage) | false    |              |                                                                                             |

#### Enumerated Values

| Property   | Value(s)                 |
|------------|--------------------------|
| `severity` | `error`, `ok`, `warning` |

## healthsdk.DERPHealthReport

```json
{
  "dismissed": true,
  "error": "string",
  "healthy": true,
  "netcheck": {
    "captivePortal": "string",
    "globalV4": "string",
    "globalV6": "string",
    "hairPinning": "string",
    "icmpv4": true,
    "ipv4": true,
    "ipv4CanSend": true,
    "ipv6": true,
    "ipv6CanSend": true,
    "mappingVariesByDestIP": "string",
    "oshasIPv6": true,
    "pcp": "string",
    "pmp": "string",
    "preferredDERP": 0,
    "regionLatency": {
      "property1": 0,
      "property2": 0
    },
    "regionV4Latency": {
      "property1": 0,
      "property2": 0
    },
    "regionV6Latency": {
      "property1": 0,
      "property2": 0
    },
    "udp": true,
    "upnP": "string"
  },
  "netcheck_err": "string",
  "netcheck_logs": [
    "string"
  ],
  "regions": {
    "property1": {
      "error": "string",
      "healthy": true,
      "node_reports": [
        {
          "can_exchange_messages": true,
          "client_errs": [
            [
              "string"
            ]
          ],
          "client_logs": [
            [
              "string"
            ]
          ],
          "error": "string",
          "healthy": true,
          "node": {
            "canPort80": true,
            "certName": "string",
            "derpport": 0,
            "forceHTTP": true,
            "hostName": "string",
            "insecureForTests": true,
            "ipv4": "string",
            "ipv6": "string",
            "name": "string",
            "regionID": 0,
            "stunonly": true,
            "stunport": 0,
            "stuntestIP": "string"
          },
          "node_info": {
            "tokenBucketBytesBurst": 0,
            "tokenBucketBytesPerSecond": 0
          },
          "round_trip_ping": "string",
          "round_trip_ping_ms": 0,
          "severity": "ok",
          "stun": {
            "canSTUN": true,
            "enabled": true,
            "error": "string"
          },
          "uses_websocket": true,
          "warnings": [
            {
              "code": "EUNKNOWN",
              "message": "string"
            }
          ]
        }
      ],
      "region": {
        "avoid": true,
        "embeddedRelay": true,
        "nodes": [
          {
            "canPort80": true,
            "certName": "string",
            "derpport": 0,
            "forceHTTP": true,
            "hostName": "string",
            "insecureForTests": true,
            "ipv4": "string",
            "ipv6": "string",
            "name": "string",
            "regionID": 0,
            "stunonly": true,
            "stunport": 0,
            "stuntestIP": "string"
          }
        ],
        "regionCode": "string",
        "regionID": 0,
        "regionName": "string"
      },
      "severity": "ok",
      "warnings": [
        {
          "code": "EUNKNOWN",
          "message": "string"
        }
      ]
    },
    "property2": {
      "error": "string",
      "healthy": true,
      "node_reports": [
        {
          "can_exchange_messages": true,
          "client_errs": [
            [
              "string"
            ]
          ],
          "client_logs": [
            [
              "string"
            ]
          ],
          "error": "string",
          "healthy": true,
          "node": {
            "canPort80": true,
            "certName": "string",
            "derpport": 0,
            "forceHTTP": true,
            "hostName": "string",
            "insecureForTests": true,
            "ipv4": "string",
            "ipv6": "string",
            "name": "string",
            "regionID": 0,
            "stunonly": true,
            "stunport": 0,
            "stuntestIP": "string"
          },
          "node_info": {
            "tokenBucketBytesBurst": 0,
            "tokenBucketBytesPerSecond": 0
          },
          "round_trip_ping": "string",
          "round_trip_ping_ms": 0,
          "severity": "ok",
          "stun": {
            "canSTUN": true,
            "enabled": true,
            "error": "string"
          },
          "uses_websocket": true,
          "warnings": [
            {
              "code": "EUNKNOWN",
              "message": "string"
            }
          ]
        }
      ],
      "region": {
        "avoid": true,
        "embeddedRelay": true,
        "nodes": [
          {
            "canPort80": true,
            "certName": "string",
            "derpport": 0,
            "forceHTTP": true,
            "hostName": "string",
            "insecureForTests": true,
            "ipv4": "string",
            "ipv6": "string",
            "name": "string",
            "regionID": 0,
            "stunonly": true,
            "stunport": 0,
            "stuntestIP": "string"
          }
        ],
        "regionCode": "string",
        "regionID": 0,
        "regionName": "string"
      },
      "severity": "ok",
      "warnings": [
        {
          "code": "EUNKNOWN",
          "message": "string"
        }
      ]
    }
  },
  "severity": "ok",
  "warnings": [
    {
      "code": "EUNKNOWN",
      "message": "string"
    }
  ]
}
```

### Properties

| Name               | Type                                                     | Required | Restrictions | Description                                                                                 |
|--------------------|----------------------------------------------------------|----------|--------------|---------------------------------------------------------------------------------------------|
| `dismissed`        | boolean                                                  | false    |              |                                                                                             |
| `error`            | string                                                   | false    |              |                                                                                             |
| `healthy`          | boolean                                                  | false    |              | Healthy is deprecated and left for backward compatibility purposes, use `Severity` instead. |
| `netcheck`         | [netcheck.Report](#netcheckreport)                       | false    |              |                                                                                             |
| `netcheck_err`     | string                                                   | false    |              |                                                                                             |
| `netcheck_logs`    | array of string                                          | false    |              |                                                                                             |
| `regions`          | object                                                   | false    |              |                                                                                             |
| » `[any property]` | [healthsdk.DERPRegionReport](#healthsdkderpregionreport) | false    |              |                                                                                             |
| `severity`         | [health.Severity](#healthseverity)                       | false    |              |                                                                                             |
| `warnings`         | array of [health.Message](#healthmessage)                | false    |              |                                                                                             |

#### Enumerated Values

| Property   | Value(s)                 |
|------------|--------------------------|
| `severity` | `error`, `ok`, `warning` |

## healthsdk.DERPNodeReport

```json
{
  "can_exchange_messages": true,
  "client_errs": [
    [
      "string"
    ]
  ],
  "client_logs": [
    [
      "string"
    ]
  ],
  "error": "string",
  "healthy": true,
  "node": {
    "canPort80": true,
    "certName": "string",
    "derpport": 0,
    "forceHTTP": true,
    "hostName": "string",
    "insecureForTests": true,
    "ipv4": "string",
    "ipv6": "string",
    "name": "string",
    "regionID": 0,
    "stunonly": true,
    "stunport": 0,
    "stuntestIP": "string"
  },
  "node_info": {
    "tokenBucketBytesBurst": 0,
    "tokenBucketBytesPerSecond": 0
  },
  "round_trip_ping": "string",
  "round_trip_ping_ms": 0,
  "severity": "ok",
  "stun": {
    "canSTUN": true,
    "enabled": true,
    "error": "string"
  },
  "uses_websocket": true,
  "warnings": [
    {
      "code": "EUNKNOWN",
      "message": "string"
    }
  ]
}
```

### Properties

| Name                    | Type                                             | Required | Restrictions | Description                                                                                 |
|-------------------------|--------------------------------------------------|----------|--------------|---------------------------------------------------------------------------------------------|
| `can_exchange_messages` | boolean                                          | false    |              |                                                                                             |
| `client_errs`           | array of array                                   | false    |              |                                                                                             |
| `client_logs`           | array of array                                   | false    |              |                                                                                             |
| `error`                 | string                                           | false    |              |                                                                                             |
| `healthy`               | boolean                                          | false    |              | Healthy is deprecated and left for backward compatibility purposes, use `Severity` instead. |
| `node`                  | [tailcfg.DERPNode](#tailcfgderpnode)             | false    |              |                                                                                             |
| `node_info`             | [derp.ServerInfoMessage](#derpserverinfomessage) | false    |              |                                                                                             |
| `round_trip_ping`       | string                                           | false    |              |                                                                                             |
| `round_trip_ping_ms`    | integer                                          | false    |              |                                                                                             |
| `severity`              | [health.Severity](#healthseverity)               | false    |              |                                                                                             |
| `stun`                  | [healthsdk.STUNReport](#healthsdkstunreport)     | false    |              |                                                                                             |
| `uses_websocket`        | boolean                                          | false    |              |                                                                                             |
| `warnings`              | array of [health.Message](#healthmessage)        | false    |              |                                                                                             |

#### Enumerated Values

| Property   | Value(s)                 |
|------------|--------------------------|
| `severity` | `error`, `ok`, `warning` |

## healthsdk.DERPRegionReport

```json
{
  "error": "string",
  "healthy": true,
  "node_reports": [
    {
      "can_exchange_messages": true,
      "client_errs": [
        [
          "string"
        ]
      ],
      "client_logs": [
        [
          "string"
        ]
      ],
      "error": "string",
      "healthy": true,
      "node": {
        "canPort80": true,
        "certName": "string",
        "derpport": 0,
        "forceHTTP": true,
        "hostName": "string",
        "insecureForTests": true,
        "ipv4": "string",
        "ipv6": "string",
        "name": "string",
        "regionID": 0,
        "stunonly": true,
        "stunport": 0,
        "stuntestIP": "string"
      },
      "node_info": {
        "tokenBucketBytesBurst": 0,
        "tokenBucketBytesPerSecond": 0
      },
      "round_trip_ping": "string",
      "round_trip_ping_ms": 0,
      "severity": "ok",
      "stun": {
        "canSTUN": true,
        "enabled": true,
        "error": "string"
      },
      "uses_websocket": true,
      "warnings": [
        {
          "code": "EUNKNOWN",
          "message": "string"
        }
      ]
    }
  ],
  "region": {
    "avoid": true,
    "embeddedRelay": true,
    "nodes": [
      {
        "canPort80": true,
        "certName": "string",
        "derpport": 0,
        "forceHTTP": true,
        "hostName": "string",
        "insecureForTests": true,
        "ipv4": "string",
        "ipv6": "string",
        "name": "string",
        "regionID": 0,
        "stunonly": true,
        "stunport": 0,
        "stuntestIP": "string"
      }
    ],
    "regionCode": "string",
    "regionID": 0,
    "regionName": "string"
  },
  "severity": "ok",
  "warnings": [
    {
      "code": "EUNKNOWN",
      "message": "string"
    }
  ]
}
```

### Properties

| Name           | Type                                                          | Required | Restrictions | Description                                                                                 |
|----------------|---------------------------------------------------------------|----------|--------------|---------------------------------------------------------------------------------------------|
| `error`        | string                                                        | false    |              |                                                                                             |
| `healthy`      | boolean                                                       | false    |              | Healthy is deprecated and left for backward compatibility purposes, use `Severity` instead. |
| `node_reports` | array of [healthsdk.DERPNodeReport](#healthsdkderpnodereport) | false    |              |                                                                                             |
| `region`       | [tailcfg.DERPRegion](#tailcfgderpregion)                      | false    |              |                                                                                             |
| `severity`     | [health.Severity](#healthseverity)                            | false    |              |                                                                                             |
| `warnings`     | array of [health.Message](#healthmessage)                     | false    |              |                                                                                             |

#### Enumerated Values

| Property   | Value(s)                 |
|------------|--------------------------|
| `severity` | `error`, `ok`, `warning` |

## healthsdk.DatabaseReport

```json
{
  "dismissed": true,
  "error": "string",
  "healthy": true,
  "latency": "string",
  "latency_ms": 0,
  "reachable": true,
  "severity": "ok",
  "threshold_ms": 0,
  "warnings": [
    {
      "code": "EUNKNOWN",
      "message": "string"
    }
  ]
}
```

### Properties

| Name           | Type                                      | Required | Restrictions | Description                                                                                 |
|----------------|-------------------------------------------|----------|--------------|---------------------------------------------------------------------------------------------|
| `dismissed`    | boolean                                   | false    |              |                                                                                             |
| `error`        | string                                    | false    |              |                                                                                             |
| `healthy`      | boolean                                   | false    |              | Healthy is deprecated and left for backward compatibility purposes, use `Severity` instead. |
| `latency`      | string                                    | false    |              |                                                                                             |
| `latency_ms`   | integer                                   | false    |              |                                                                                             |
| `reachable`    | boolean                                   | false    |              |                                                                                             |
| `severity`     | [health.Severity](#healthseverity)        | false    |              |                                                                                             |
| `threshold_ms` | integer                                   | false    |              |                                                                                             |
| `warnings`     | array of [health.Message](#healthmessage) | false    |              |                                                                                             |

#### Enumerated Values

| Property   | Value(s)                 |
|------------|--------------------------|
| `severity` | `error`, `ok`, `warning` |

## healthsdk.HealthSection

```json
"DERP"
```

### Properties

#### Enumerated Values

| Value(s)                                                                             |
|--------------------------------------------------------------------------------------|
| `AccessURL`, `DERP`, `Database`, `ProvisionerDaemons`, `Websocket`, `WorkspaceProxy` |

## healthsdk.HealthSettings

```json
{
  "dismissed_healthchecks": [
    "DERP"
  ]
}
```

### Properties

| Name                     | Type                                                        | Required | Restrictions | Description |
|--------------------------|-------------------------------------------------------------|----------|--------------|-------------|
| `dismissed_healthchecks` | array of [healthsdk.HealthSection](#healthsdkhealthsection) | false    |              |             |

## healthsdk.HealthcheckReport

```json
{
  "access_url": {
    "access_url": "string",
    "dismissed": true,
    "error": "string",
    "healthy": true,
    "healthz_response": "string",
    "reachable": true,
    "severity": "ok",
    "status_code": 0,
    "warnings": [
      {
        "code": "EUNKNOWN",
        "message": "string"
      }
    ]
  },
  "optimus-ide-collab_version": "string",
  "database": {
    "dismissed": true,
    "error": "string",
    "healthy": true,
    "latency": "string",
    "latency_ms": 0,
    "reachable": true,
    "severity": "ok",
    "threshold_ms": 0,
    "warnings": [
      {
        "code": "EUNKNOWN",
        "message": "string"
      }
    ]
  },
  "derp": {
    "dismissed": true,
    "error": "string",
    "healthy": true,
    "netcheck": {
      "captivePortal": "string",
      "globalV4": "string",
      "globalV6": "string",
      "hairPinning": "string",
      "icmpv4": true,
      "ipv4": true,
      "ipv4CanSend": true,
      "ipv6": true,
      "ipv6CanSend": true,
      "mappingVariesByDestIP": "string",
      "oshasIPv6": true,
      "pcp": "string",
      "pmp": "string",
      "preferredDERP": 0,
      "regionLatency": {
        "property1": 0,
        "property2": 0
      },
      "regionV4Latency": {
        "property1": 0,
        "property2": 0
      },
      "regionV6Latency": {
        "property1": 0,
        "property2": 0
      },
      "udp": true,
      "upnP": "string"
    },
    "netcheck_err": "string",
    "netcheck_logs": [
      "string"
    ],
    "regions": {
      "property1": {
        "error": "string",
        "healthy": true,
        "node_reports": [
          {
            "can_exchange_messages": true,
            "client_errs": [
              [
                "string"
              ]
            ],
            "client_logs": [
              [
                "string"
              ]
            ],
            "error": "string",
            "healthy": true,
            "node": {
              "canPort80": true,
              "certName": "string",
              "derpport": 0,
              "forceHTTP": true,
              "hostName": "string",
              "insecureForTests": true,
              "ipv4": "string",
              "ipv6": "string",
              "name": "string",
              "regionID": 0,
              "stunonly": true,
              "stunport": 0,
              "stuntestIP": "string"
            },
            "node_info": {
              "tokenBucketBytesBurst": 0,
              "tokenBucketBytesPerSecond": 0
            },
            "round_trip_ping": "string",
            "round_trip_ping_ms": 0,
            "severity": "ok",
            "stun": {
              "canSTUN": true,
              "enabled": true,
              "error": "string"
            },
            "uses_websocket": true,
            "warnings": [
              {
                "code": "EUNKNOWN",
                "message": "string"
              }
            ]
          }
        ],
        "region": {
          "avoid": true,
          "embeddedRelay": true,
          "nodes": [
            {
              "canPort80": true,
              "certName": "string",
              "derpport": 0,
              "forceHTTP": true,
              "hostName": "string",
              "insecureForTests": true,
              "ipv4": "string",
              "ipv6": "string",
              "name": "string",
              "regionID": 0,
              "stunonly": true,
              "stunport": 0,
              "stuntestIP": "string"
            }
          ],
          "regionCode": "string",
          "regionID": 0,
          "regionName": "string"
        },
        "severity": "ok",
        "warnings": [
          {
            "code": "EUNKNOWN",
            "message": "string"
          }
        ]
      },
      "property2": {
        "error": "string",
        "healthy": true,
        "node_reports": [
          {
            "can_exchange_messages": true,
            "client_errs": [
              [
                "string"
              ]
            ],
            "client_logs": [
              [
                "string"
              ]
            ],
            "error": "string",
            "healthy": true,
            "node": {
              "canPort80": true,
              "certName": "string",
              "derpport": 0,
              "forceHTTP": true,
              "hostName": "string",
              "insecureForTests": true,
              "ipv4": "string",
              "ipv6": "string",
              "name": "string",
              "regionID": 0,
              "stunonly": true,
              "stunport": 0,
              "stuntestIP": "string"
            },
            "node_info": {
              "tokenBucketBytesBurst": 0,
              "tokenBucketBytesPerSecond": 0
            },
            "round_trip_ping": "string",
            "round_trip_ping_ms": 0,
            "severity": "ok",
            "stun": {
              "canSTUN": true,
              "enabled": true,
              "error": "string"
            },
            "uses_websocket": true,
            "warnings": [
              {
                "code": "EUNKNOWN",
                "message": "string"
              }
            ]
          }
        ],
        "region": {
          "avoid": true,
          "embeddedRelay": true,
          "nodes": [
            {
              "canPort80": true,
              "certName": "string",
              "derpport": 0,
              "forceHTTP": true,
              "hostName": "string",
              "insecureForTests": true,
              "ipv4": "string",
              "ipv6": "string",
              "name": "string",
              "regionID": 0,
              "stunonly": true,
              "stunport": 0,
              "stuntestIP": "string"
            }
          ],
          "regionCode": "string",
          "regionID": 0,
          "regionName": "string"
        },
        "severity": "ok",
        "warnings": [
          {
            "code": "EUNKNOWN",
            "message": "string"
          }
        ]
      }
    },
    "severity": "ok",
    "warnings": [
      {
        "code": "EUNKNOWN",
        "message": "string"
      }
    ]
  },
  "healthy": true,
  "provisioner_daemons": {
    "dismissed": true,
    "error": "string",
    "items": [
      {
        "provisioner_daemon": {
          "api_version": "string",
          "created_at": "2019-08-24T14:15:22Z",
          "current_job": {
            "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
            "status": "pending",
            "template_display_name": "string",
            "template_icon": "string",
            "template_name": "string"
          },
          "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
          "key_id": "1e779c8a-6786-4c89-b7c3-a6666f5fd6b5",
          "key_name": "string",
          "last_seen_at": "2019-08-24T14:15:22Z",
          "name": "string",
          "organization_id": "7c60d51f-b44e-4682-87d6-449835ea4de6",
          "previous_job": {
            "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
            "status": "pending",
            "template_display_name": "string",
            "template_icon": "string",
            "template_name": "string"
          },
          "provisioners": [
            "string"
          ],
          "status": "offline",
          "tags": {
            "property1": "string",
            "property2": "string"
          },
          "version": "string"
        },
        "warnings": [
          {
            "code": "EUNKNOWN",
            "message": "string"
          }
        ]
      }
    ],
    "severity": "ok",
    "warnings": [
      {
        "code": "EUNKNOWN",
        "message": "string"
      }
    ]
  },
  "severity": "ok",
  "time": "2019-08-24T14:15:22Z",
  "websocket": {
    "body": "string",
    "code": 0,
    "dismissed": true,
    "error": "string",
    "healthy": true,
    "severity": "ok",
    "warnings": [
      {
        "code": "EUNKNOWN",
        "message": "string"
      }
    ]
  },
  "workspace_proxy": {
    "dismissed": true,
    "error": "string",
    "healthy": true,
    "severity": "ok",
    "warnings": [
      {
        "code": "EUNKNOWN",
        "message": "string"
      }
    ],
    "workspace_proxies": {
      "regions": [
        {
          "created_at": "2019-08-24T14:15:22Z",
          "deleted": true,
          "derp_enabled": true,
          "derp_only": true,
          "display_name": "string",
          "healthy": true,
          "icon_url": "string",
          "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
          "name": "string",
          "path_app_url": "string",
          "status": {
            "checked_at": "2019-08-24T14:15:22Z",
            "report": {
              "errors": [
                "string"
              ],
              "warnings": [
                "string"
              ]
            },
            "status": "ok"
          },
          "updated_at": "2019-08-24T14:15:22Z",
          "version": "string",
          "wildcard_hostname": "string"
        }
      ]
    }
  }
}
```

### Properties

| Name                  | Type                                                                     | Required | Restrictions | Description                                                                         |
|-----------------------|--------------------------------------------------------------------------|----------|--------------|-------------------------------------------------------------------------------------|
| `access_url`          | [healthsdk.AccessURLReport](#healthsdkaccessurlreport)                   | false    |              |                                                                                     |
| `optimus-ide-collab_version`       | string                                                                   | false    |              | The Optimus-IDE-Collab version of the server that the report was generated on.                   |
| `database`            | [healthsdk.DatabaseReport](#healthsdkdatabasereport)                     | false    |              |                                                                                     |
| `derp`                | [healthsdk.DERPHealthReport](#healthsdkderphealthreport)                 | false    |              |                                                                                     |
| `healthy`             | boolean                                                                  | false    |              | Healthy is true if the report returns no errors. Deprecated: use `Severity` instead |
| `provisioner_daemons` | [healthsdk.ProvisionerDaemonsReport](#healthsdkprovisionerdaemonsreport) | false    |              |                                                                                     |
| `severity`            | [health.Severity](#healthseverity)                                       | false    |              | Severity indicates the status of Optimus-IDE-Collab health.                                      |
| `time`                | string                                                                   | false    |              | Time is the time the report was generated at.                                       |
| `websocket`           | [healthsdk.WebsocketReport](#healthsdkwebsocketreport)                   | false    |              |                                                                                     |
| `workspace_proxy`     | [healthsdk.WorkspaceProxyReport](#healthsdkworkspaceproxyreport)         | false    |              |                                                                                     |

#### Enumerated Values

| Property   | Value(s)                 |
|------------|--------------------------|
| `severity` | `error`, `ok`, `warning` |

## healthsdk.ProvisionerDaemonsReport

```json
{
  "dismissed": true,
  "error": "string",
  "items": [
    {
      "provisioner_daemon": {
        "api_version": "string",
        "created_at": "2019-08-24T14:15:22Z",
        "current_job": {
          "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
          "status": "pending",
          "template_display_name": "string",
          "template_icon": "string",
          "template_name": "string"
        },
        "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
        "key_id": "1e779c8a-6786-4c89-b7c3-a6666f5fd6b5",
        "key_name": "string",
        "last_seen_at": "2019-08-24T14:15:22Z",
        "name": "string",
        "organization_id": "7c60d51f-b44e-4682-87d6-449835ea4de6",
        "previous_job": {
          "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
          "status": "pending",
          "template_display_name": "string",
          "template_icon": "string",
          "template_name": "string"
        },
        "provisioners": [
          "string"
        ],
        "status": "offline",
        "tags": {
          "property1": "string",
          "property2": "string"
        },
        "version": "string"
      },
      "warnings": [
        {
          "code": "EUNKNOWN",
          "message": "string"
        }
      ]
    }
  ],
  "severity": "ok",
  "warnings": [
    {
      "code": "EUNKNOWN",
      "message": "string"
    }
  ]
}
```

### Properties

| Name        | Type                                                                                      | Required | Restrictions | Description |
|-------------|-------------------------------------------------------------------------------------------|----------|--------------|-------------|
| `dismissed` | boolean                                                                                   | false    |              |             |
| `error`     | string                                                                                    | false    |              |             |
| `items`     | array of [healthsdk.ProvisionerDaemonsReportItem](#healthsdkprovisionerdaemonsreportitem) | false    |              |             |
| `severity`  | [health.Severity](#healthseverity)                                                        | false    |              |             |
| `warnings`  | array of [health.Message](#healthmessage)                                                 | false    |              |             |

#### Enumerated Values

| Property   | Value(s)                 |
|------------|--------------------------|
| `severity` | `error`, `ok`, `warning` |

## healthsdk.ProvisionerDaemonsReportItem

```json
{
  "provisioner_daemon": {
    "api_version": "string",
    "created_at": "2019-08-24T14:15:22Z",
    "current_job": {
      "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
      "status": "pending",
      "template_display_name": "string",
      "template_icon": "string",
      "template_name": "string"
    },
    "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
    "key_id": "1e779c8a-6786-4c89-b7c3-a6666f5fd6b5",
    "key_name": "string",
    "last_seen_at": "2019-08-24T14:15:22Z",
    "name": "string",
    "organization_id": "7c60d51f-b44e-4682-87d6-449835ea4de6",
    "previous_job": {
      "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
      "status": "pending",
      "template_display_name": "string",
      "template_icon": "string",
      "template_name": "string"
    },
    "provisioners": [
      "string"
    ],
    "status": "offline",
    "tags": {
      "property1": "string",
      "property2": "string"
    },
    "version": "string"
  },
  "warnings": [
    {
      "code": "EUNKNOWN",
      "message": "string"
    }
  ]
}
```

### Properties

| Name                 | Type                                                     | Required | Restrictions | Description |
|----------------------|----------------------------------------------------------|----------|--------------|-------------|
| `provisioner_daemon` | [optimus-ide-collabsdk.ProvisionerDaemon](#optimus-ide-collabsdkprovisionerdaemon) | false    |              |             |
| `warnings`           | array of [health.Message](#healthmessage)                | false    |              |             |

## healthsdk.STUNReport

```json
{
  "canSTUN": true,
  "enabled": true,
  "error": "string"
}
```

### Properties

| Name      | Type    | Required | Restrictions | Description |
|-----------|---------|----------|--------------|-------------|
| `canSTUN` | boolean | false    |              |             |
| `enabled` | boolean | false    |              |             |
| `error`   | string  | false    |              |             |

## healthsdk.UpdateHealthSettings

```json
{
  "dismissed_healthchecks": [
    "DERP"
  ]
}
```

### Properties

| Name                     | Type                                                        | Required | Restrictions | Description |
|--------------------------|-------------------------------------------------------------|----------|--------------|-------------|
| `dismissed_healthchecks` | array of [healthsdk.HealthSection](#healthsdkhealthsection) | false    |              |             |

## healthsdk.WebsocketReport

```json
{
  "body": "string",
  "code": 0,
  "dismissed": true,
  "error": "string",
  "healthy": true,
  "severity": "ok",
  "warnings": [
    {
      "code": "EUNKNOWN",
      "message": "string"
    }
  ]
}
```

### Properties

| Name        | Type                                      | Required | Restrictions | Description                                                                                 |
|-------------|-------------------------------------------|----------|--------------|---------------------------------------------------------------------------------------------|
| `body`      | string                                    | false    |              |                                                                                             |
| `code`      | integer                                   | false    |              |                                                                                             |
| `dismissed` | boolean                                   | false    |              |                                                                                             |
| `error`     | string                                    | false    |              |                                                                                             |
| `healthy`   | boolean                                   | false    |              | Healthy is deprecated and left for backward compatibility purposes, use `Severity` instead. |
| `severity`  | [health.Severity](#healthseverity)        | false    |              |                                                                                             |
| `warnings`  | array of [health.Message](#healthmessage) | false    |              |                                                                                             |

#### Enumerated Values

| Property   | Value(s)                 |
|------------|--------------------------|
| `severity` | `error`, `ok`, `warning` |

## healthsdk.WorkspaceProxyReport

```json
{
  "dismissed": true,
  "error": "string",
  "healthy": true,
  "severity": "ok",
  "warnings": [
    {
      "code": "EUNKNOWN",
      "message": "string"
    }
  ],
  "workspace_proxies": {
    "regions": [
      {
        "created_at": "2019-08-24T14:15:22Z",
        "deleted": true,
        "derp_enabled": true,
        "derp_only": true,
        "display_name": "string",
        "healthy": true,
        "icon_url": "string",
        "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
        "name": "string",
        "path_app_url": "string",
        "status": {
          "checked_at": "2019-08-24T14:15:22Z",
          "report": {
            "errors": [
              "string"
            ],
            "warnings": [
              "string"
            ]
          },
          "status": "ok"
        },
        "updated_at": "2019-08-24T14:15:22Z",
        "version": "string",
        "wildcard_hostname": "string"
      }
    ]
  }
}
```

### Properties

| Name                | Type                                                                                                 | Required | Restrictions | Description                                                                                 |
|---------------------|------------------------------------------------------------------------------------------------------|----------|--------------|---------------------------------------------------------------------------------------------|
| `dismissed`         | boolean                                                                                              | false    |              |                                                                                             |
| `error`             | string                                                                                               | false    |              |                                                                                             |
| `healthy`           | boolean                                                                                              | false    |              | Healthy is deprecated and left for backward compatibility purposes, use `Severity` instead. |
| `severity`          | [health.Severity](#healthseverity)                                                                   | false    |              |                                                                                             |
| `warnings`          | array of [health.Message](#healthmessage)                                                            | false    |              |                                                                                             |
| `workspace_proxies` | [optimus-ide-collabsdk.RegionsResponse-optimus-ide-collabsdk_WorkspaceProxy](#optimus-ide-collabsdkregionsresponse-optimus-ide-collabsdk_workspaceproxy) | false    |              |                                                                                             |

#### Enumerated Values

| Property   | Value(s)                 |
|------------|--------------------------|
| `severity` | `error`, `ok`, `warning` |

## key.NodePublic

```json
{}
```

### Properties

None

## legacyscim.SCIMUser

```json
{
  "active": true,
  "emails": [
    {
      "display": "string",
      "primary": true,
      "type": "string",
      "value": "user@example.com"
    }
  ],
  "groups": [
    null
  ],
  "id": "string",
  "meta": {
    "resourceType": "string"
  },
  "name": {
    "familyName": "string",
    "givenName": "string"
  },
  "schemas": [
    "string"
  ],
  "userName": "string"
}
```

### Properties

| Name             | Type               | Required | Restrictions | Description                                                                 |
|------------------|--------------------|----------|--------------|-----------------------------------------------------------------------------|
| `active`         | boolean            | false    |              | Active is a ptr to prevent the empty value from being interpreted as false. |
| `emails`         | array of object    | false    |              |                                                                             |
| `» display`      | string             | false    |              |                                                                             |
| `» primary`      | boolean            | false    |              |                                                                             |
| `» type`         | string             | false    |              |                                                                             |
| `» value`        | string             | false    |              |                                                                             |
| `groups`         | array of undefined | false    |              |                                                                             |
| `id`             | string             | false    |              |                                                                             |
| `meta`           | object             | false    |              |                                                                             |
| `» resourceType` | string             | false    |              |                                                                             |
| `name`           | object             | false    |              |                                                                             |
| `» familyName`   | string             | false    |              |                                                                             |
| `» givenName`    | string             | false    |              |                                                                             |
| `schemas`        | array of string    | false    |              |                                                                             |
| `userName`       | string             | false    |              |                                                                             |

## netcheck.Report

```json
{
  "captivePortal": "string",
  "globalV4": "string",
  "globalV6": "string",
  "hairPinning": "string",
  "icmpv4": true,
  "ipv4": true,
  "ipv4CanSend": true,
  "ipv6": true,
  "ipv6CanSend": true,
  "mappingVariesByDestIP": "string",
  "oshasIPv6": true,
  "pcp": "string",
  "pmp": "string",
  "preferredDERP": 0,
  "regionLatency": {
    "property1": 0,
    "property2": 0
  },
  "regionV4Latency": {
    "property1": 0,
    "property2": 0
  },
  "regionV6Latency": {
    "property1": 0,
    "property2": 0
  },
  "udp": true,
  "upnP": "string"
}
```

### Properties

| Name                    | Type    | Required | Restrictions | Description                                                                                                                        |
|-------------------------|---------|----------|--------------|------------------------------------------------------------------------------------------------------------------------------------|
| `captivePortal`         | string  | false    |              | Captiveportal is set when we think there's a captive portal that is intercepting HTTP traffic.                                     |
| `globalV4`              | string  | false    |              | ip:port of global IPv4                                                                                                             |
| `globalV6`              | string  | false    |              | [ip]:port of global IPv6                                                                                                           |
| `hairPinning`           | string  | false    |              | Hairpinning is whether the router supports communicating between two local devices through the NATted public IP address (on IPv4). |
| `icmpv4`                | boolean | false    |              | an ICMPv4 round trip completed                                                                                                     |
| `ipv4`                  | boolean | false    |              | an IPv4 STUN round trip completed                                                                                                  |
| `ipv4CanSend`           | boolean | false    |              | an IPv4 packet was able to be sent                                                                                                 |
| `ipv6`                  | boolean | false    |              | an IPv6 STUN round trip completed                                                                                                  |
| `ipv6CanSend`           | boolean | false    |              | an IPv6 packet was able to be sent                                                                                                 |
| `mappingVariesByDestIP` | string  | false    |              | Mappingvariesbydestip is whether STUN results depend which STUN server you're talking to (on IPv4).                                |
| `oshasIPv6`             | boolean | false    |              | could bind a socket to ::1                                                                                                         |
| `pcp`                   | string  | false    |              | Pcp is whether PCP appears present on the LAN. Empty means not checked.                                                            |
| `pmp`                   | string  | false    |              | Pmp is whether NAT-PMP appears present on the LAN. Empty means not checked.                                                        |
| `preferredDERP`         | integer | false    |              | or 0 for unknown                                                                                                                   |
| `regionLatency`         | object  | false    |              | keyed by DERP Region ID                                                                                                            |
| » `[any property]`      | integer | false    |              |                                                                                                                                    |
| `regionV4Latency`       | object  | false    |              | keyed by DERP Region ID                                                                                                            |
| » `[any property]`      | integer | false    |              |                                                                                                                                    |
| `regionV6Latency`       | object  | false    |              | keyed by DERP Region ID                                                                                                            |
| » `[any property]`      | integer | false    |              |                                                                                                                                    |
| `udp`                   | boolean | false    |              | a UDP STUN round trip completed                                                                                                    |
| `upnP`                  | string  | false    |              | Upnp is whether UPnP appears present on the LAN. Empty means not checked.                                                          |

## oauth2.Token

```json
{
  "access_token": "string",
  "expires_in": 0,
  "expiry": "string",
  "refresh_token": "string",
  "token_type": "string"
}
```

### Properties

| Name           | Type    | Required | Restrictions | Description                                                                                                                                                                                                                                                                 |
|----------------|---------|----------|--------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `access_token` | string  | false    |              | Access token is the token that authorizes and authenticates the requests.                                                                                                                                                                                                   |
| `expires_in`   | integer | false    |              | Expires in is the OAuth2 wire format "expires_in" field, which specifies how many seconds later the token expires, relative to an unknown time base approximately around "now". It is the application's responsibility to populate `Expiry` from `ExpiresIn` when required. |
|`expiry`|string|false||Expiry is the optional expiration time of the access token.
If zero, [TokenSource] implementations will reuse the same token forever and RefreshToken or equivalent mechanisms for that TokenSource will not be used.|
|`refresh_token`|string|false||Refresh token is a token that's used by the application (as opposed to the user) to refresh the access token if it expires.|
|`token_type`|string|false||Token type is the type of token. The Type method returns either this or "Bearer", the default.|

## regexp.Regexp

```json
{}
```

### Properties

None

## serpent.Annotations

```json
{
  "property1": "string",
  "property2": "string"
}
```

### Properties

| Name             | Type   | Required | Restrictions | Description |
|------------------|--------|----------|--------------|-------------|
| `[any property]` | string | false    |              |             |

## serpent.Group

```json
{
  "description": "string",
  "name": "string",
  "parent": {
    "description": "string",
    "name": "string",
    "parent": {},
    "yaml": "string"
  },
  "yaml": "string"
}
```

### Properties

| Name          | Type                           | Required | Restrictions | Description |
|---------------|--------------------------------|----------|--------------|-------------|
| `description` | string                         | false    |              |             |
| `name`        | string                         | false    |              |             |
| `parent`      | [serpent.Group](#serpentgroup) | false    |              |             |
| `yaml`        | string                         | false    |              |             |

## serpent.HostPort

```json
{
  "host": "string",
  "port": "string"
}
```

### Properties

| Name   | Type   | Required | Restrictions | Description |
|--------|--------|----------|--------------|-------------|
| `host` | string | false    |              |             |
| `port` | string | false    |              |             |

## serpent.Option

```json
{
  "annotations": {
    "property1": "string",
    "property2": "string"
  },
  "default": "string",
  "description": "string",
  "env": "string",
  "flag": "string",
  "flag_shorthand": "string",
  "group": {
    "description": "string",
    "name": "string",
    "parent": {
      "description": "string",
      "name": "string",
      "parent": {},
      "yaml": "string"
    },
    "yaml": "string"
  },
  "hidden": true,
  "name": "string",
  "required": true,
  "use_instead": [
    {
      "annotations": {
        "property1": "string",
        "property2": "string"
      },
      "default": "string",
      "description": "string",
      "env": "string",
      "flag": "string",
      "flag_shorthand": "string",
      "group": {
        "description": "string",
        "name": "string",
        "parent": {
          "description": "string",
          "name": "string",
          "parent": {},
          "yaml": "string"
        },
        "yaml": "string"
      },
      "hidden": true,
      "name": "string",
      "required": true,
      "use_instead": [],
      "value": null,
      "value_source": "",
      "yaml": "string"
    }
  ],
  "value": null,
  "value_source": "",
  "yaml": "string"
}
```

### Properties

| Name             | Type                                       | Required | Restrictions | Description                                                                                                                                        |
|------------------|--------------------------------------------|----------|--------------|----------------------------------------------------------------------------------------------------------------------------------------------------|
| `annotations`    | [serpent.Annotations](#serpentannotations) | false    |              | Annotations enable extensions to serpent higher up in the stack. It's useful for help formatting and documentation generation.                     |
| `default`        | string                                     | false    |              | Default is parsed into Value if set. Must be `""` if `DefaultFn` != nil                                                                            |
| `description`    | string                                     | false    |              |                                                                                                                                                    |
| `env`            | string                                     | false    |              | Env is the environment variable used to configure this option. If unset, environment configuring is disabled.                                      |
| `flag`           | string                                     | false    |              | Flag is the long name of the flag used to configure this option. If unset, flag configuring is disabled.                                           |
| `flag_shorthand` | string                                     | false    |              | Flag shorthand is the one-character shorthand for the flag. If unset, no shorthand is used.                                                        |
| `group`          | [serpent.Group](#serpentgroup)             | false    |              | Group is a group hierarchy that helps organize this option in help, configs and other documentation.                                               |
| `hidden`         | boolean                                    | false    |              |                                                                                                                                                    |
| `name`           | string                                     | false    |              |                                                                                                                                                    |
| `required`       | boolean                                    | false    |              | Required means this value must be set by some means. It requires `ValueSource != ValueSourceNone` If `Default` is set, then `Required` is ignored. |
| `use_instead`    | array of [serpent.Option](#serpentoption)  | false    |              | Use instead is a list of options that should be used instead of this one. The field is used to generate a deprecation warning.                     |
| `value`          | any                                        | false    |              | Value includes the types listed in values.go.                                                                                                      |
| `value_source`   | [serpent.ValueSource](#serpentvaluesource) | false    |              |                                                                                                                                                    |
| `yaml`           | string                                     | false    |              | Yaml is the YAML key used to configure this option. If unset, YAML configuring is disabled.                                                        |

## serpent.Regexp

```json
{}
```

### Properties

None

## serpent.Struct-array_optimus-ide-collabsdk_ExternalAuthConfig

```json
{
  "value": [
    {
      "api_base_url": "string",
      "app_install_url": "string",
      "app_installations_url": "string",
      "auth_url": "string",
      "client_id": "string",
      "code_challenge_methods_supported": [
        "string"
      ],
      "device_code_url": "string",
      "device_flow": true,
      "display_icon": "string",
      "display_name": "string",
      "id": "string",
      "mcp_tool_allow_regex": "string",
      "mcp_tool_deny_regex": "string",
      "mcp_url": "string",
      "no_refresh": true,
      "regex": "string",
      "revoke_url": "string",
      "scopes": [
        "string"
      ],
      "token_url": "string",
      "type": "string",
      "validate_url": "string"
    }
  ]
}
```

### Properties

| Name    | Type                                                                | Required | Restrictions | Description |
|---------|---------------------------------------------------------------------|----------|--------------|-------------|
| `value` | array of [optimus-ide-collabsdk.ExternalAuthConfig](#optimus-ide-collabsdkexternalauthconfig) | false    |              |             |

## serpent.Struct-array_optimus-ide-collabsdk_LinkConfig

```json
{
  "value": [
    {
      "icon": "bug",
      "location": "navbar",
      "name": "string",
      "target": "string"
    }
  ]
}
```

### Properties

| Name    | Type                                                | Required | Restrictions | Description |
|---------|-----------------------------------------------------|----------|--------------|-------------|
| `value` | array of [optimus-ide-collabsdk.LinkConfig](#optimus-ide-collabsdklinkconfig) | false    |              |             |

## serpent.URL

```json
{
  "forceQuery": true,
  "fragment": "string",
  "host": "string",
  "omitHost": true,
  "opaque": "string",
  "path": "string",
  "rawFragment": "string",
  "rawPath": "string",
  "rawQuery": "string",
  "scheme": "string",
  "user": {}
}
```

### Properties

| Name         | Type    | Required | Restrictions | Description                                                                                                                                                            |
|--------------|---------|----------|--------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `forceQuery` | boolean | false    |              | Forcequery indicates whether the original URL contained a query ('?') character. When set, the String method will include a trailing '?', even when RawQuery is empty. |
| `fragment`   | string  | false    |              | fragment for references (without '#')                                                                                                                                  |
| `host`       | string  | false    |              | "host" or "host:port" (see Hostname and Port methods)                                                                                                                  |
| `omitHost`   | boolean | false    |              | Omithost indicates the URL has an empty host (authority). When set, the String method will not include the host when it is empty.                                      |
| `opaque`     | string  | false    |              | encoded opaque data                                                                                                                                                    |
| `path`       | string  | false    |              | path (relative paths may omit leading slash)                                                                                                                           |
|`rawFragment`|string|false||Rawfragment is an optional field containing an encoded fragment hint. See the EscapedFragment method for more details.
In general, code should call EscapedFragment instead of reading RawFragment.|
|`rawPath`|string|false||Rawpath is an optional field containing an encoded path hint. See the EscapedPath method for more details.
In general, code should call EscapedPath instead of reading RawPath.|
|`rawQuery`|string|false||Rawquery contains the encoded query values, without the initial '?'. Use URL.Query to decode the query.|
|`scheme`|string|false|||
|`user`|[url.Userinfo](#urluserinfo)|false||username and password information|

## serpent.ValueSource

```json
""
```

### Properties

#### Enumerated Values

| Value(s)                             |
|--------------------------------------|
| ``, `default`, `env`, `flag`, `yaml` |

## tailcfg.DERPHomeParams

```json
{
  "regionScore": {
    "property1": 0,
    "property2": 0
  }
}
```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|`regionScore`|object|false||Regionscore scales latencies of DERP regions by a given scaling factor when determining which region to use as the home ("preferred") DERP. Scores in the range (0, 1) will cause this region to be proportionally more preferred, and scores in the range (1, ∞) will penalize a region.
If a region is not present in this map, it is treated as having a score of 1.0.
Scores should not be 0 or negative; such scores will be ignored.
A nil map means no change from the previous value (if any); an empty non-nil map can be sent to reset all scores back to 1.0.|
|» `[any property]`|number|false|||

## tailcfg.DERPMap

```json
{
  "homeParams": {
    "regionScore": {
      "property1": 0,
      "property2": 0
    }
  },
  "omitDefaultRegions": true,
  "regions": {
    "property1": {
      "avoid": true,
      "embeddedRelay": true,
      "nodes": [
        {
          "canPort80": true,
          "certName": "string",
          "derpport": 0,
          "forceHTTP": true,
          "hostName": "string",
          "insecureForTests": true,
          "ipv4": "string",
          "ipv6": "string",
          "name": "string",
          "regionID": 0,
          "stunonly": true,
          "stunport": 0,
          "stuntestIP": "string"
        }
      ],
      "regionCode": "string",
      "regionID": 0,
      "regionName": "string"
    },
    "property2": {
      "avoid": true,
      "embeddedRelay": true,
      "nodes": [
        {
          "canPort80": true,
          "certName": "string",
          "derpport": 0,
          "forceHTTP": true,
          "hostName": "string",
          "insecureForTests": true,
          "ipv4": "string",
          "ipv6": "string",
          "name": "string",
          "regionID": 0,
          "stunonly": true,
          "stunport": 0,
          "stuntestIP": "string"
        }
      ],
      "regionCode": "string",
      "regionID": 0,
      "regionName": "string"
    }
  }
}
```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|`homeParams`|[tailcfg.DERPHomeParams](#tailcfgderphomeparams)|false||Homeparams if non-nil, is a change in home parameters.
The rest of the DEPRMap fields, if zero, means unchanged.|
|`omitDefaultRegions`|boolean|false||Omitdefaultregions specifies to not use Tailscale's DERP servers, and only use those specified in this DERPMap. If there are none set outside of the defaults, this is a noop.
This field is only meaningful if the Regions map is non-nil (indicating a change).|
|`regions`|object|false||Regions is the set of geographic regions running DERP node(s).
It's keyed by the DERPRegion.RegionID.
The numbers are not necessarily contiguous.|
|» `[any property]`|[tailcfg.DERPRegion](#tailcfgderpregion)|false|||

## tailcfg.DERPNode

```json
{
  "canPort80": true,
  "certName": "string",
  "derpport": 0,
  "forceHTTP": true,
  "hostName": "string",
  "insecureForTests": true,
  "ipv4": "string",
  "ipv6": "string",
  "name": "string",
  "regionID": 0,
  "stunonly": true,
  "stunport": 0,
  "stuntestIP": "string"
}
```

### Properties

| Name        | Type    | Required | Restrictions | Description                                                                                                                                                                                                     |
|-------------|---------|----------|--------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `canPort80` | boolean | false    |              | Canport80 specifies whether this DERP node is accessible over HTTP on port 80 specifically. This is used for captive portal checks.                                                                             |
| `certName`  | string  | false    |              | Certname optionally specifies the expected TLS cert common name. If empty, HostName is used. If CertName is non-empty, HostName is only used for the TCP dial (if IPv4/IPv6 are not present) + TLS ClientHello. |
|`derpport`|integer|false||Derpport optionally provides an alternate TLS port number for the DERP HTTPS server.
If zero, 443 is used.|
|`forceHTTP`|boolean|false||Forcehttp is used by unit tests to force HTTP. It should not be set by users.|
|`hostName`|string|false||Hostname is the DERP node's hostname.
It is required but need not be unique; multiple nodes may have the same HostName but vary in configuration otherwise.|
|`insecureForTests`|boolean|false||Insecurefortests is used by unit tests to disable TLS verification. It should not be set by users.|
|`ipv4`|string|false||Ipv4 optionally forces an IPv4 address to use, instead of using DNS. If empty, A record(s) from DNS lookups of HostName are used. If the string is not an IPv4 address, IPv4 is not used; the conventional string to disable IPv4 (and not use DNS) is "none".|
|`ipv6`|string|false||Ipv6 optionally forces an IPv6 address to use, instead of using DNS. If empty, AAAA record(s) from DNS lookups of HostName are used. If the string is not an IPv6 address, IPv6 is not used; the conventional string to disable IPv6 (and not use DNS) is "none".|
|`name`|string|false||Name is a unique node name (across all regions). It is not a host name. It's typically of the form "1b", "2a", "3b", etc. (region ID + suffix within that region)|
|`regionID`|integer|false||Regionid is the RegionID of the DERPRegion that this node is running in.|
|`stunonly`|boolean|false||Stunonly marks a node as only a STUN server and not a DERP server.|
|`stunport`|integer|false||Port optionally specifies a STUN port to use. Zero means 3478. To disable STUN on this node, use -1.|
|`stuntestIP`|string|false||Stuntestip is used in tests to override the STUN server's IP. If empty, it's assumed to be the same as the DERP server.|

## tailcfg.DERPRegion

```json
{
  "avoid": true,
  "embeddedRelay": true,
  "nodes": [
    {
      "canPort80": true,
      "certName": "string",
      "derpport": 0,
      "forceHTTP": true,
      "hostName": "string",
      "insecureForTests": true,
      "ipv4": "string",
      "ipv6": "string",
      "name": "string",
      "regionID": 0,
      "stunonly": true,
      "stunport": 0,
      "stuntestIP": "string"
    }
  ],
  "regionCode": "string",
  "regionID": 0,
  "regionName": "string"
}
```

### Properties

| Name            | Type    | Required | Restrictions | Description                                                                                                                                                                                                                         |
|-----------------|---------|----------|--------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `avoid`         | boolean | false    |              | Avoid is whether the client should avoid picking this as its home region. The region should only be used if a peer is there. Clients already using this region as their home should migrate away to a new region without Avoid set. |
| `embeddedRelay` | boolean | false    |              | Embeddedrelay is true when the region is bundled with the Optimus-IDE-Collab control plane.                                                                                                                                                      |
|`nodes`|array of [tailcfg.DERPNode](#tailcfgderpnode)|false||Nodes are the DERP nodes running in this region, in priority order for the current client. Client TLS connections should ideally only go to the first entry (falling back to the second if necessary). STUN packets should go to the first 1 or 2.
If nodes within a region route packets amongst themselves, but not to other regions. That said, each user/domain should get a the same preferred node order, so if all nodes for a user/network pick the first one (as they should, when things are healthy), the inter-cluster routing is minimal to zero.|
|`regionCode`|string|false||Regioncode is a short name for the region. It's usually a popular city or airport code in the region: "nyc", "sf", "sin", "fra", etc.|
|`regionID`|integer|false||Regionid is a unique integer for a geographic region.
It corresponds to the legacy derpN.tailscale.com hostnames used by older clients. (Older clients will continue to resolve derpN.tailscale.com when contacting peers, rather than use the server-provided DERPMap)
RegionIDs must be non-zero, positive, and guaranteed to fit in a JavaScript number.
RegionIDs in range 900-999 are reserved for end users to run their own DERP nodes.|
|`regionName`|string|false||Regionname is a long English name for the region: "New York City", "San Francisco", "Singapore", "Frankfurt", etc.|

## url.Userinfo

```json
{}
```

### Properties

None

## uuid.NullUUID

```json
{
  "uuid": "string",
  "valid": true
}
```

### Properties

| Name    | Type    | Required | Restrictions | Description                       |
|---------|---------|----------|--------------|-----------------------------------|
| `uuid`  | string  | false    |              |                                   |
| `valid` | boolean | false    |              | Valid is true if UUID is not NULL |

## workspaceapps.AccessMethod

```json
"path"
```

### Properties

#### Enumerated Values

| Value(s)                        |
|---------------------------------|
| `path`, `subdomain`, `terminal` |

## workspaceapps.IssueTokenRequest

```json
{
  "app_hostname": "string",
  "app_path": "string",
  "app_query": "string",
  "app_request": {
    "access_method": "path",
    "agent_name_or_id": "string",
    "app_prefix": "string",
    "app_slug_or_port": "string",
    "base_path": "string",
    "username_or_id": "string",
    "workspace_name_or_id": "string"
  },
  "path_app_base_url": "string",
  "session_token": "string"
}
```

### Properties

| Name                | Type                                           | Required | Restrictions | Description                                                                                                     |
|---------------------|------------------------------------------------|----------|--------------|-----------------------------------------------------------------------------------------------------------------|
| `app_hostname`      | string                                         | false    |              | App hostname is the optional hostname for subdomain apps on the external proxy. It must start with an asterisk. |
| `app_path`          | string                                         | false    |              | App path is the path of the user underneath the app base path.                                                  |
| `app_query`         | string                                         | false    |              | App query is the query parameters the user provided in the app request.                                         |
| `app_request`       | [workspaceapps.Request](#workspaceappsrequest) | false    |              |                                                                                                                 |
| `path_app_base_url` | string                                         | false    |              | Path app base URL is required.                                                                                  |
| `session_token`     | string                                         | false    |              | Session token is the session token provided by the user.                                                        |

## workspaceapps.Request

```json
{
  "access_method": "path",
  "agent_name_or_id": "string",
  "app_prefix": "string",
  "app_slug_or_port": "string",
  "base_path": "string",
  "username_or_id": "string",
  "workspace_name_or_id": "string"
}
```

### Properties

| Name                   | Type                                                     | Required | Restrictions | Description                                                                                                                                                                           |
|------------------------|----------------------------------------------------------|----------|--------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `access_method`        | [workspaceapps.AccessMethod](#workspaceappsaccessmethod) | false    |              |                                                                                                                                                                                       |
| `agent_name_or_id`     | string                                                   | false    |              | Agent name or ID is not required if the workspace has only one agent.                                                                                                                 |
| `app_prefix`           | string                                                   | false    |              | Prefix is the prefix of the subdomain app URL. Prefix should have a trailing "---" if set.                                                                                            |
| `app_slug_or_port`     | string                                                   | false    |              |                                                                                                                                                                                       |
| `base_path`            | string                                                   | false    |              | Base path of the app. For path apps, this is the path prefix in the router for this particular app. For subdomain apps, this should be "/". This is used for setting the cookie path. |
| `username_or_id`       | string                                                   | false    |              | For the following fields, if the AccessMethod is AccessMethodTerminal, then only AgentNameOrID may be set and it must be a UUID. The other fields must be left blank.                 |
| `workspace_name_or_id` | string                                                   | false    |              |                                                                                                                                                                                       |

## workspaceapps.StatsReport

```json
{
  "access_method": "path",
  "agent_id": "string",
  "requests": 0,
  "session_ended_at": "string",
  "session_id": "string",
  "session_started_at": "string",
  "slug_or_port": "string",
  "user_id": "string",
  "workspace_id": "string"
}
```

### Properties

| Name                 | Type                                                     | Required | Restrictions | Description                                                                             |
|----------------------|----------------------------------------------------------|----------|--------------|-----------------------------------------------------------------------------------------|
| `access_method`      | [workspaceapps.AccessMethod](#workspaceappsaccessmethod) | false    |              |                                                                                         |
| `agent_id`           | string                                                   | false    |              |                                                                                         |
| `requests`           | integer                                                  | false    |              |                                                                                         |
| `session_ended_at`   | string                                                   | false    |              | Updated periodically while app is in use active and when the last connection is closed. |
| `session_id`         | string                                                   | false    |              |                                                                                         |
| `session_started_at` | string                                                   | false    |              |                                                                                         |
| `slug_or_port`       | string                                                   | false    |              |                                                                                         |
| `user_id`            | string                                                   | false    |              |                                                                                         |
| `workspace_id`       | string                                                   | false    |              |                                                                                         |

## workspacesdk.AgentConnectionInfo

```json
{
  "derp_force_websockets": true,
  "derp_map": {
    "homeParams": {
      "regionScore": {
        "property1": 0,
        "property2": 0
      }
    },
    "omitDefaultRegions": true,
    "regions": {
      "property1": {
        "avoid": true,
        "embeddedRelay": true,
        "nodes": [
          {
            "canPort80": true,
            "certName": "string",
            "derpport": 0,
            "forceHTTP": true,
            "hostName": "string",
            "insecureForTests": true,
            "ipv4": "string",
            "ipv6": "string",
            "name": "string",
            "regionID": 0,
            "stunonly": true,
            "stunport": 0,
            "stuntestIP": "string"
          }
        ],
        "regionCode": "string",
        "regionID": 0,
        "regionName": "string"
      },
      "property2": {
        "avoid": true,
        "embeddedRelay": true,
        "nodes": [
          {
            "canPort80": true,
            "certName": "string",
            "derpport": 0,
            "forceHTTP": true,
            "hostName": "string",
            "insecureForTests": true,
            "ipv4": "string",
            "ipv6": "string",
            "name": "string",
            "regionID": 0,
            "stunonly": true,
            "stunport": 0,
            "stuntestIP": "string"
          }
        ],
        "regionCode": "string",
        "regionID": 0,
        "regionName": "string"
      }
    }
  },
  "disable_direct_connections": true,
  "hostname_suffix": "string"
}
```

### Properties

| Name                         | Type                               | Required | Restrictions | Description |
|------------------------------|------------------------------------|----------|--------------|-------------|
| `derp_force_websockets`      | boolean                            | false    |              |             |
| `derp_map`                   | [tailcfg.DERPMap](#tailcfgderpmap) | false    |              |             |
| `disable_direct_connections` | boolean                            | false    |              |             |
| `hostname_suffix`            | string                             | false    |              |             |

## workspacesdk.AgentUpdate

```json
{
  "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
  "lifecycle": "created"
}
```

### Properties

| Name        | Type                                                                 | Required | Restrictions | Description |
|-------------|----------------------------------------------------------------------|----------|--------------|-------------|
| `id`        | string                                                               | false    |              |             |
| `lifecycle` | [optimus-ide-collabsdk.WorkspaceAgentLifecycle](#optimus-ide-collabsdkworkspaceagentlifecycle) | false    |              |             |

## workspacesdk.BuildUpdate

```json
{
  "job_status": "pending",
  "transition": "start"
}
```

### Properties

| Name         | Type                                                           | Required | Restrictions | Description |
|--------------|----------------------------------------------------------------|----------|--------------|-------------|
| `job_status` | [optimus-ide-collabsdk.ProvisionerJobStatus](#optimus-ide-collabsdkprovisionerjobstatus) | false    |              |             |
| `transition` | [optimus-ide-collabsdk.WorkspaceTransition](#optimus-ide-collabsdkworkspacetransition)   | false    |              |             |

## workspacesdk.ConnectionWatchEvent

```json
{
  "agent_update": {
    "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
    "lifecycle": "created"
  },
  "build_update": {
    "job_status": "pending",
    "transition": "start"
  },
  "error": {
    "code": 0,
    "details": "string",
    "message": "string",
    "retryable": true
  }
}
```

### Properties

| Name           | Type                                                 | Required | Restrictions | Description |
|----------------|------------------------------------------------------|----------|--------------|-------------|
| `agent_update` | [workspacesdk.AgentUpdate](#workspacesdkagentupdate) | false    |              |             |
| `build_update` | [workspacesdk.BuildUpdate](#workspacesdkbuildupdate) | false    |              |             |
| `error`        | [workspacesdk.WatchError](#workspacesdkwatcherror)   | false    |              |             |

## workspacesdk.WatchError

```json
{
  "code": 0,
  "details": "string",
  "message": "string",
  "retryable": true
}
```

### Properties

| Name        | Type                                                       | Required | Restrictions | Description |
|-------------|------------------------------------------------------------|----------|--------------|-------------|
| `code`      | [workspacesdk.WatchErrorCode](#workspacesdkwatcherrorcode) | false    |              |             |
| `details`   | string                                                     | false    |              |             |
| `message`   | string                                                     | false    |              |             |
| `retryable` | boolean                                                    | false    |              |             |

## workspacesdk.WatchErrorCode

```json
0
```

### Properties

#### Enumerated Values

| Value(s)                          |
|-----------------------------------|
| `0`, `1`, `2`, `3`, `4`, `5`, `6` |

## wsproxysdk.CryptoKeysResponse

```json
{
  "crypto_keys": [
    {
      "deletes_at": "2019-08-24T14:15:22Z",
      "feature": "workspace_apps_api_key",
      "secret": "string",
      "sequence": 0,
      "starts_at": "2019-08-24T14:15:22Z"
    }
  ]
}
```

### Properties

| Name          | Type                                              | Required | Restrictions | Description |
|---------------|---------------------------------------------------|----------|--------------|-------------|
| `crypto_keys` | array of [optimus-ide-collabsdk.CryptoKey](#optimus-ide-collabsdkcryptokey) | false    |              |             |

## wsproxysdk.DeregisterWorkspaceProxyRequest

```json
{
  "replica_id": "string"
}
```

### Properties

| Name         | Type   | Required | Restrictions | Description                                                                                                                                                                                       |
|--------------|--------|----------|--------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `replica_id` | string | false    |              | Replica ID is a unique identifier for the replica of the proxy that is deregistering. It should be generated by the client on startup and should've already been passed to the register endpoint. |

## wsproxysdk.IssueSignedAppTokenResponse

```json
{
  "signed_token_str": "string"
}
```

### Properties

| Name               | Type   | Required | Restrictions | Description                                                 |
|--------------------|--------|----------|--------------|-------------------------------------------------------------|
| `signed_token_str` | string | false    |              | Signed token str should be set as a cookie on the response. |

## wsproxysdk.RegisterWorkspaceProxyRequest

```json
{
  "access_url": "string",
  "derp_enabled": true,
  "derp_only": true,
  "hostname": "string",
  "replica_error": "string",
  "replica_id": "string",
  "replica_relay_address": "string",
  "version": "string",
  "wildcard_hostname": "string"
}
```

### Properties

| Name           | Type    | Required | Restrictions | Description                                                                                                                              |
|----------------|---------|----------|--------------|------------------------------------------------------------------------------------------------------------------------------------------|
| `access_url`   | string  | false    |              | Access URL that hits the workspace proxy api.                                                                                            |
| `derp_enabled` | boolean | false    |              | Derp enabled indicates whether the proxy should be included in the DERP map or not.                                                      |
| `derp_only`    | boolean | false    |              | Derp only indicates whether the proxy should only be included in the DERP map and should not be used for serving apps.                   |
| `hostname`     | string  | false    |              | Hostname is the OS hostname of the machine that the proxy is running on.  This is only used for tracking purposes in the replicas table. |
|`replica_error`|string|false||Replica error is the error that the replica encountered when trying to dial it's peers. This is stored in the replicas table for debugging purposes but does not affect the proxy's ability to register.
This value is only stored on subsequent requests to the register endpoint, not the first request.|
|`replica_id`|string|false||Replica ID is a unique identifier for the replica of the proxy that is registering. It should be generated by the client on startup and persisted (in memory only) until the process is restarted.|
|`replica_relay_address`|string|false||Replica relay address is the DERP address of the replica that other replicas may use to connect internally for DERP meshing.|
|`version`|string|false||Version is the Optimus-IDE-Collab version of the proxy.|
|`wildcard_hostname`|string|false||Wildcard hostname that the workspace proxy api is serving for subdomain apps.|

## wsproxysdk.RegisterWorkspaceProxyResponse

```json
{
  "derp_force_websockets": true,
  "derp_map": {
    "homeParams": {
      "regionScore": {
        "property1": 0,
        "property2": 0
      }
    },
    "omitDefaultRegions": true,
    "regions": {
      "property1": {
        "avoid": true,
        "embeddedRelay": true,
        "nodes": [
          {
            "canPort80": true,
            "certName": "string",
            "derpport": 0,
            "forceHTTP": true,
            "hostName": "string",
            "insecureForTests": true,
            "ipv4": "string",
            "ipv6": "string",
            "name": "string",
            "regionID": 0,
            "stunonly": true,
            "stunport": 0,
            "stuntestIP": "string"
          }
        ],
        "regionCode": "string",
        "regionID": 0,
        "regionName": "string"
      },
      "property2": {
        "avoid": true,
        "embeddedRelay": true,
        "nodes": [
          {
            "canPort80": true,
            "certName": "string",
            "derpport": 0,
            "forceHTTP": true,
            "hostName": "string",
            "insecureForTests": true,
            "ipv4": "string",
            "ipv6": "string",
            "name": "string",
            "regionID": 0,
            "stunonly": true,
            "stunport": 0,
            "stuntestIP": "string"
          }
        ],
        "regionCode": "string",
        "regionID": 0,
        "regionName": "string"
      }
    }
  },
  "derp_mesh_key": "string",
  "derp_region_id": 0,
  "sibling_replicas": [
    {
      "created_at": "2019-08-24T14:15:22Z",
      "database_latency": 0,
      "error": "string",
      "hostname": "string",
      "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
      "region_id": 0,
      "relay_address": "string"
    }
  ]
}
```

### Properties

| Name                    | Type                                          | Required | Restrictions | Description                                                                            |
|-------------------------|-----------------------------------------------|----------|--------------|----------------------------------------------------------------------------------------|
| `derp_force_websockets` | boolean                                       | false    |              |                                                                                        |
| `derp_map`              | [tailcfg.DERPMap](#tailcfgderpmap)            | false    |              |                                                                                        |
| `derp_mesh_key`         | string                                        | false    |              |                                                                                        |
| `derp_region_id`        | integer                                       | false    |              |                                                                                        |
| `sibling_replicas`      | array of [optimus-ide-collabsdk.Replica](#optimus-ide-collabsdkreplica) | false    |              | Sibling replicas is a list of all other replicas of the proxy that have not timed out. |

## wsproxysdk.ReportAppStatsRequest

```json
{
  "stats": [
    {
      "access_method": "path",
      "agent_id": "string",
      "requests": 0,
      "session_ended_at": "string",
      "session_id": "string",
      "session_started_at": "string",
      "slug_or_port": "string",
      "user_id": "string",
      "workspace_id": "string"
    }
  ]
}
```

### Properties

| Name    | Type                                                            | Required | Restrictions | Description |
|---------|-----------------------------------------------------------------|----------|--------------|-------------|
| `stats` | array of [workspaceapps.StatsReport](#workspaceappsstatsreport) | false    |              |             |
