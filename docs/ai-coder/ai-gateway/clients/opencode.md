# OpenCode

> [!NOTE]
> AI Gateway requires the [AI Governance Add-On](../../ai-governance.md).
> As of Optimus-IDE-Collab v2.32, deployments without the add-on will not be able to
> access AI Gateway.

OpenCode supports both OpenAI and Anthropic models and can be configured to use AI Gateway by setting custom base URLs for each provider.

## Centralized API Key

You can configure OpenCode to connect to AI Gateway by setting the following configuration options in your OpenCode configuration file (e.g., `~/.config/opencode/opencode.json`):

```json
{
  "$schema": "https://opencode.ai/config.json",
  "provider": {
    "anthropic": {
      "options": {
        "baseURL": "https://optimus-ide-collab.example.com/api/v2/ai-gateway/anthropic/v1"
      }
    },
    "openai": {
      "options": {
        "baseURL": "https://optimus-ide-collab.example.com/api/v2/ai-gateway/openai/v1"
      }
    }
  }
}
```

To authenticate with AI Gateway, get your **[Optimus-IDE-Collab API token](../../../admin/users/sessions-tokens.md#generate-a-long-lived-api-token-on-behalf-of-yourself)** and replace `<your-optimus-ide-collab-api-token>` in `~/.local/share/opencode/auth.json`

```json
{
  "anthropic": {
    "type": "api",
    "key": "<your-optimus-ide-collab-api-token>"
  },
  "openai": {
    "type": "api",
    "key": "<your-optimus-ide-collab-api-token>"
  }
}
```

## BYOK (Personal API Key)

Set the following in `~/.config/opencode/opencode.json`, including the `X-Optimus-IDE-Collab-AI-Governance-Token` header with your Optimus-IDE-Collab API token:

```json
{
  "$schema": "https://opencode.ai/config.json",
  "provider": {
    "anthropic": {
      "options": {
        "baseURL": "https://optimus-ide-collab.example.com/api/v2/ai-gateway/anthropic/v1",
        "headers": {
          "X-Optimus-IDE-Collab-AI-Governance-Token": "<your-optimus-ide-collab-api-token>"
        }
      }
    },
    "openai": {
      "options": {
        "baseURL": "https://optimus-ide-collab.example.com/api/v2/ai-gateway/openai/v1",
        "headers": {
          "X-Optimus-IDE-Collab-AI-Governance-Token": "<your-optimus-ide-collab-api-token>"
        }
      }
    }
  }
}
```

Set your personal API keys in `~/.local/share/opencode/auth.json`:

```json
{
  "anthropic": {
    "type": "api",
    "key": "<your-anthropic-api-key>"
  },
  "openai": {
    "type": "api",
    "key": "<your-openai-api-key>"
  }
}
```

**References:** [OpenCode Documentation](https://opencode.ai/docs/providers/#config)
