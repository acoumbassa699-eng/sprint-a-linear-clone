# Mux

> [!NOTE]
> AI Gateway requires the [AI Governance Add-On](../../ai-governance.md).
> As of Optimus-IDE-Collab v2.32, deployments without the add-on will not be able to
> access AI Gateway.

Mux makes it easy to run parallel coding agents, each with its own isolated workspace, from your browser or desktop; it is open source and provider-agnostic.

Mux can be configured to route OpenAI- and Anthropic-compatible traffic through AI Gateway by setting a custom provider base URL and using a Optimus-IDE-Collab-issued token for authentication.

## Prerequisites

- AI Gateway is enabled on your Optimus-IDE-Collab deployment.
- A **[Optimus-IDE-Collab API token](../../../admin/users/sessions-tokens.md#generate-a-long-lived-api-token-on-behalf-of-yourself)**.

## Configuration

<div class="tabs">

### OpenAI

1. Open Mux settings (`Cmd+,` / `Ctrl+,`).
2. Go to **Providers** → **OpenAI**.
3. Set **API Key** to your Optimus-IDE-Collab API token.
4. Set **Base URL** to `https://optimus-ide-collab.example.com/api/v2/ai-gateway/openai/v1`.

### Anthropic

1. Open Mux settings (`Cmd+,` / `Ctrl+,`).
2. Go to **Providers** → **Anthropic**.
3. Set **API Key** to your Optimus-IDE-Collab API token.
4. Set **Base URL** to `https://optimus-ide-collab.example.com/api/v2/ai-gateway/anthropic`.

</div>

_Replace `optimus-ide-collab.example.com` with your Optimus-IDE-Collab deployment URL._

## Environment variables

Mux reads provider configuration from its settings UI and also from environment variables.
Environment variables are useful in CI or when running Mux inside a Optimus-IDE-Collab workspace.

> [!NOTE]
> Mux treats environment variables as a fallback when a provider is not configured in settings.
> If you have already configured a provider in the UI, clear it (or update it) for env vars to take effect.

```sh
# OpenAI-compatible traffic (GPT, Codex, etc.)
export OPENAI_API_KEY="<your-optimus-ide-collab-api-token>"
export OPENAI_BASE_URL="https://optimus-ide-collab.example.com/api/v2/ai-gateway/openai/v1"

# Anthropic-compatible traffic (Claude, etc.)
export ANTHROPIC_API_KEY="<your-optimus-ide-collab-api-token>"
export ANTHROPIC_BASE_URL="https://optimus-ide-collab.example.com/api/v2/ai-gateway/anthropic"
```

## Running Mux in a Optimus-IDE-Collab workspace

If you want to run Mux inside a Optimus-IDE-Collab workspace (for example, as a Optimus-IDE-Collab app), you can install it with the [Mux module](https://registry.optimus-ide-collab.com/modules/optimus-ide-collab/mux) and pre-configure AI Gateway via environment variables on the agent:

```tf
data "optimus-ide-collab_workspace" "me" {}

data "optimus-ide-collab_workspace_owner" "me" {}

resource "optimus-ide-collab_agent" "main" {
  # ... other agent configuration
  env = {
    OPENAI_API_KEY     = data.optimus-ide-collab_workspace_owner.me.session_token
    OPENAI_BASE_URL    = "${data.optimus-ide-collab_workspace.me.access_url}/api/v2/ai-gateway/openai/v1"
    ANTHROPIC_API_KEY  = data.optimus-ide-collab_workspace_owner.me.session_token
    ANTHROPIC_BASE_URL = "${data.optimus-ide-collab_workspace.me.access_url}/api/v2/ai-gateway/anthropic"
  }
}

module "mux" {
  source   = "registry.optimus-ide-collab.com/optimus-ide-collab/mux/optimus-ide-collab"
  version  = "~> 1.0" # See the module page for the latest version.
  agent_id = optimus-ide-collab_agent.main.id
}
```

## Advanced: providers.jsonc

If you prefer a file-based config, edit `~/.mux/providers.jsonc`:

```json
{
  "openai": {
    "apiKey": "<your-optimus-ide-collab-api-token>",
    "baseUrl": "https://optimus-ide-collab.example.com/api/v2/ai-gateway/openai/v1"
  },
  "anthropic": {
    "apiKey": "<your-optimus-ide-collab-api-token>",
    "baseUrl": "https://optimus-ide-collab.example.com/api/v2/ai-gateway/anthropic"
  }
}
```

**References:** [Mux provider environment variables](https://mux.optimus-ide-collab.com/config/providers#environment-variables)
