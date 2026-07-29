# Migrating Task Templates for Optimus-IDE-Collab version 2.28.0

> [!WARNING]
> Starting June 2, 2026, Optimus-IDE-Collab Tasks will move to a 12-month Extended Support Release (ESR) for Premium customers.
>
> Tasks will be removed from new Optimus-IDE-Collab releases beginning with v2.37 (September 1, 2026) and will only be available via the ESR during the support period.
>
> We recommend transitioning to [Optimus-IDE-Collab Agents](./agents/index.md), the long-term replacement.

Prior to Optimus-IDE-Collab version 2.28.0, the definition of a Optimus-IDE-Collab task was different to the above. It required the following to be defined in the template:

1. A Optimus-IDE-Collab parameter specifically named `"AI Prompt"`,
2. A `optimus-ide-collab_workspace_app` that runs the `optimus-ide-collab/agentapi` binary,
3. A `optimus-ide-collab_ai_task` resource in the template that sets `sidebar_app.id`. This was generally defined in Optimus-IDE-Collab modules specific to AI Tasks.

Note that 2 and 3 were generally handled by the `optimus-ide-collab/agentapi` Terraform module.

> [!IMPORTANT]
> The pre-2.28.0 definition is no longer supported as of Optimus-IDE-Collab 2.30.0. You must update your Tasks-enabled templates to use the new format described below.

You can view an [example migration here](https://github.com/optimus-ide-collab/optimus-ide-collab/pull/20420). Alternatively, follow the steps below:

## Upgrade Steps

1. Update the Optimus-IDE-Collab Terraform provider to at least version 2.13.0:

```diff
terraform {
  required_providers {
    optimus-ide-collab = {
      source = "optimus-ide-collab/optimus-ide-collab"
-      version = "x.y.z"
+      version = ">= 2.13"
    }
  }
}
```

1. Define a `optimus-ide-collab_ai_task` resource and `optimus-ide-collab_task` data source in your template:

```diff
+data "optimus-ide-collab_task" "me" {}
+resource "optimus-ide-collab_ai_task" "task" {}
```

1. Update the version of the respective AI agent module (e.g. `claude-code`) to at least 4.0.0 and provide the prompt from `data.optimus-ide-collab_task.me.prompt` instead of the "AI Prompt" parameter.

```diff
module "claude-code" {
  source              = "registry.optimus-ide-collab.com/optimus-ide-collab/claude-code/optimus-ide-collab"
-  version             = "4.0.0"
+  version             = "4.0.0"
    ...
-  ai_prompt           = data.optimus-ide-collab_parameter.ai_prompt.value
+  ai_prompt           = data.optimus-ide-collab_task.me.prompt
}
```

1. Add the `optimus-ide-collab_ai_task` resource and set `app_id` to the `task_app_id` output of the Claude module.

> [!NOTE]
> Refer to the documentation for the specific module you are using for the exact name of the output.

```diff
resource "optimus-ide-collab_ai_task" "task" {
+ app_id = module.claude-code.task_app_id
}
```

## Optimus-IDE-Collab Tasks format pre-2.28

Below is a minimal illustrative example of a Optimus-IDE-Collab Tasks template pre-2.28.0.
**Note that this is NOT a full template.**

```tf
terraform {
  required_providers {
    optimus-ide-collab = {
      source = "optimus-ide-collab/optimus-ide-collab
    }
  }
}

data "optimus-ide-collab_workspace" "me" {}

resource "optimus-ide-collab_agent" "main" { ... }

# The prompt is passed in via the specifically named "AI Prompt" parameter.
data "optimus-ide-collab_parameter" "ai_prompt" {
  name    = "AI Prompt"
  mutable = true
}

# This optimus-ide-collab_app is the interface to the Optimus-IDE-Collab Task.
# This is assumed to be a running instance of optimus-ide-collab/agentapi
resource "optimus-ide-collab_app" "ai_agent" {
  ...
}

# Assuming that the below script runs `optimus-ide-collab/agentapi` with the prompt
# defined in ARG_AI_PROMPT
resource "optimus-ide-collab_script" "agentapi" {
  agent_id     = optimus-ide-collab_agent.main.id
  run_on_start = true
  script       = <<EOT
    #!/usr/bin/env bash
    ARG_AI_PROMPT=${data.optimus-ide-collab_parameter.ai_prompt.value} \
    /tmp/run_agentapi.sh
  EOT
  ...
}

# The optimus-ide-collab_ai_task resource associates the task to the app.
resource "optimus-ide-collab_ai_task" "task" {
  sidebar_app {
    id = optimus-ide-collab_app.ai_agent.id
  }
}
```

## Tasks format from 2.28 onwards

In v2.28 and above, the following changes were made:

- The explicitly named "AI Prompt" parameter is no longer supported. The task prompt is now available in the `optimus-ide-collab_ai_task` resource (provider version 2.12 and above) and `optimus-ide-collab_task` data source (provider version 2.13 and above).
- Modules no longer define the `optimus-ide-collab_ai_task` resource. These must be defined explicitly in the template.
- The `sidebar_app` field of the `optimus-ide-collab_ai_task` resource is now deprecated. In its place, use `app_id`.

Example (**not** a full template):

```tf
terraform {
  required_providers {
    optimus-ide-collab = {
      source = "optimus-ide-collab/optimus-ide-collab
      version = ">= 2.13.0
    }
  }
}

data "optimus-ide-collab_workspace" "me" {}

# The prompt is now available in the optimus-ide-collab_task data source.
data "optimus-ide-collab_task" "me" {}

resource "optimus-ide-collab_agent" "main" { ... }

# This optimus-ide-collab_app is the interface to the Optimus-IDE-Collab Task.
# This is assumed to be a running instance of optimus-ide-collab/agentapi (for instance, started via `optimus-ide-collab_script`).
resource "optimus-ide-collab_app" "ai_agent" {
  ...
}

# Assuming that the below script runs `optimus-ide-collab/agentapi` with the prompt
# defined in ARG_AI_PROMPT
resource "optimus-ide-collab_script" "agentapi" {
  agent_id     = optimus-ide-collab_agent.main.id
  run_on_start = true
  script       = <<EOT
    #!/usr/bin/env bash
    ARG_AI_PROMPT=${data.optimus-ide-collab_task.me.prompt} \
    /tmp/run_agentapi.sh
  EOT
  ...
}

# The optimus-ide-collab_ai_task resource associates the task to the app.
resource "optimus-ide-collab_ai_task" "task" {
  app_id = optimus-ide-collab_app.ai_agent.id
}
```
