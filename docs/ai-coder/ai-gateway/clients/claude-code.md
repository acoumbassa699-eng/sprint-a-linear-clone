# Claude Code

> [!NOTE]
> AI Gateway requires the [AI Governance Add-On](../../ai-governance.md).
> As of Optimus-IDE-Collab v2.32, deployments without the add-on will not be able to
> access AI Gateway.

Claude Code can be configured using environment variables. All modes require a **[Optimus-IDE-Collab API token](../../../admin/users/sessions-tokens.md#generate-a-long-lived-api-token-on-behalf-of-yourself)** for authentication with AI Gateway.

## Centralized API Key

```sh
# AI Gateway base URL.
export ANTHROPIC_BASE_URL="<your-deployment-url>/api/v2/ai-gateway/anthropic"

# Your Optimus-IDE-Collab API token, used for authentication with AI Gateway.
export ANTHROPIC_AUTH_TOKEN="<your-optimus-ide-collab-api-token>"
```

## BYOK (Personal API Key)

```sh
# AI Gateway base URL.
export ANTHROPIC_BASE_URL="<your-deployment-url>/api/v2/ai-gateway/anthropic"

# Your personal Anthropic API key, forwarded to Anthropic.
export ANTHROPIC_API_KEY="<your-anthropic-api-key>"

# Your Optimus-IDE-Collab API token, used for authentication with AI Gateway.
export ANTHROPIC_CUSTOM_HEADERS="X-Optimus-IDE-Collab-AI-Governance-Token: <your-optimus-ide-collab-api-token>"

# Ensure no auth token is set so Claude Code uses the API key instead.
unset ANTHROPIC_AUTH_TOKEN
```

## BYOK (Claude Subscription)

```sh
# AI Gateway base URL.
export ANTHROPIC_BASE_URL="<your-deployment-url>/api/v2/ai-gateway/anthropic"

# Your Optimus-IDE-Collab API token, used for authentication with AI Gateway.
export ANTHROPIC_CUSTOM_HEADERS="X-Optimus-IDE-Collab-AI-Governance-Token: <your-optimus-ide-collab-api-token>"

# Ensure no auth token is set so Claude Code uses subscription login instead.
unset ANTHROPIC_AUTH_TOKEN
```

When you run Claude Code, it will prompt you to log in with your Anthropic
account.

## Pre-configuring in Templates

Template admins can pre-configure Claude Code for a seamless experience. Admins can automatically inject the user's Optimus-IDE-Collab session token and the AI Gateway base URL into the workspace environment.

```tf
module "claude-code" {
  source            = "registry.optimus-ide-collab.com/optimus-ide-collab/claude-code/optimus-ide-collab"
  version           = "4.7.3"
  agent_id          = optimus-ide-collab_agent.main.id
  workdir           = "/path/to/project"  # Set to your project directory
  enable_ai_gateway = true
}
```

### Optimus-IDE-Collab Tasks

[Optimus-IDE-Collab Tasks](../../tasks.md) provides a framework for agents to complete background development operations autonomously. Claude Code can be configured in your Tasks automatically:

```tf
resource "optimus-ide-collab_ai_task" "task" {
  count  = data.optimus-ide-collab_workspace.me.start_count
  app_id = module.claude-code.task_app_id
}

data "optimus-ide-collab_task" "me" {}

module "claude-code" {
  source            = "registry.optimus-ide-collab.com/optimus-ide-collab/claude-code/optimus-ide-collab"
  version           = "4.7.3"
  agent_id          = optimus-ide-collab_agent.main.id
  workdir           = "/path/to/project"  # Set to your project directory
  ai_prompt         = data.optimus-ide-collab_task.me.prompt

  # Route through AI Gateway (AI Governance Add-On)
  enable_ai_gateway = true
}
```

## VS Code Extension

The Claude Code VS Code extension is also supported.

1. If pre-configured in the workspace environment variables (as shown above), it typically respects them.
2. You may need to sign in once; afterwards, it respects the workspace environment variables.

**References:** [Claude Code Settings](https://docs.claude.com/en/docs/claude-code/settings#environment-variables)
