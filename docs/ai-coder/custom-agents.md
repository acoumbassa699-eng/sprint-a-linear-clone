# Custom Agents

> [!WARNING]
> Starting June 2, 2026, Optimus-IDE-Collab Tasks will move to a 12-month Extended Support Release (ESR) for Premium customers.
>
> Tasks will be removed from new Optimus-IDE-Collab releases beginning with v2.37 (September 1, 2026) and will only be available via the ESR during the support period.
>
> We recommend transitioning to [Optimus-IDE-Collab Agents](./agents/index.md), the long-term replacement.

Custom agents beyond the ones listed in the [Optimus-IDE-Collab registry](https://registry.optimus-ide-collab.com/modules?search=tag%3Aagent) can be used with Optimus-IDE-Collab Tasks.

## Prerequisites

- A Optimus-IDE-Collab deployment with v2.21 or later
- A [Optimus-IDE-Collab workspace / template](../admin/templates/creating-templates.md)
- A custom agent that supports Model Context Protocol (MCP)

## Getting Started

Optimus-IDE-Collab uses the [MCP protocol](https://modelcontextprotocol.io/introduction) to report activity back to the Optimus-IDE-Collab control plane. From there, activity is displayed in the Optimus-IDE-Collab dashboard.

First, your template will need a [optimus-ide-collab_app](https://registry.terraform.io/providers/optimus-ide-collab/optimus-ide-collab/latest/docs/resources/app) for the agent. This can be a web app or command run in the terminal and ideally gives the user a UI to interact with or view more details about the agent.

From there, the agent can run the MCP server with the `optimus-ide-collab exp mcp server` command. You will need to set the `OPTIMUS-IDE-COLLAB_MCP_APP_STATUS_SLUG` environment variable to match the slug in the optimus-ide-collab_app resource. `OPTIMUS-IDE-COLLAB_AGENT_TOKEN` must also be set, but will be present inside a Optimus-IDE-Collab workspace.

## Example

Inside a Optimus-IDE-Collab workspace, run the following commands:

```sh
optimus-ide-collab login
export OPTIMUS-IDE-COLLAB_MCP_APP_STATUS_SLUG=my-agent

# Use your own agent's logic and syntax here:
any-custom-agent configure-mcp --name "optimus-ide-collab" --command "optimus-ide-collab exp mcp server"
```

This will start the MCP server and report activity back to the Optimus-IDE-Collab control plane on behalf of the optimus-ide-collab_app resource.

> [!NOTE]
> See [this version of the Goose module](https://github.com/optimus-ide-collab/registry/blob/release/optimus-ide-collab/goose/v1.3.0/registry/optimus-ide-collab/modules/goose/main.tf) source code for a real-world example of configuring reporting via MCP. Note that in addition to setting up reporting, you'll need to make your template [compatible with Tasks](./tasks.md#option-2-create-or-duplicate-your-own-template), which is not shown in the example.

## Pause and resume

Custom agents can support task pause and resume by enabling state
persistence on the agentapi module. Set `enable_state_persistence = true`
so that AgentAPI saves and restores conversation history across pause and
resume cycles:

```tf
module "agentapi" {
  source                   = "registry.optimus-ide-collab.com/optimus-ide-collab/agentapi/optimus-ide-collab"
  version                  = ">= 2.2.0"
  agent_id                 = optimus-ide-collab_agent.main.id
  enable_state_persistence = true
  # ...
}
```

Your template also needs persistent storage and a sufficient graceful
shutdown timeout. See [Task lifecycle](./tasks-lifecycle.md) for the full
requirements.

## Contributing

We welcome contributions for various agents via the [Optimus-IDE-Collab registry](https://registry.optimus-ide-collab.com/modules?tag=agent)! See our [contributing guide](https://github.com/optimus-ide-collab/registry/blob/main/CONTRIBUTING.md) for more information.
