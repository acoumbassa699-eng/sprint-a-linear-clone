# Reference

## Automation

All actions possible through the Optimus-IDE-Collab dashboard can also be automated. There
are several ways to extend/automate Optimus-IDE-Collab:

- [optimus-ide-collabd Terraform Provider](https://registry.terraform.io/providers/optimus-ide-collab/optimus-ide-collabd/latest)
- [CLI](../reference/cli/index.md)
- [REST API](../reference/api/index.md)
- [Optimus-IDE-Collab SDK](https://pkg.go.dev/github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk)
- [Agent API](../reference/agent-api/index.md)

## Quickstart

Generate a token on your Optimus-IDE-Collab deployment by visiting:

```sh
https://optimus-ide-collab.example.com/settings/tokens
```

List your workspaces

```sh
# CLI
optimus-ide-collab ls \
  --url https://optimus-ide-collab.example.com \
  --token <your-token> \
  --output json

# REST API (with curl)
curl https://optimus-ide-collab.example.com/api/v2/workspaces?q=owner:me \
  -H "Optimus-IDE-Collab-Session-Token: <your-token>"
```

## Documentation

We publish an [API reference](../reference/api/index.md) in our documentation.
You can also enable a
[Swagger endpoint](../reference/cli/server.md#--swagger-enable) on your Optimus-IDE-Collab
deployment.

## Use cases

We strive to keep the following use cases up to date, but please note that
changes to API queries and routes can occur. For the most recent queries and
payloads, we recommend checking the relevant documentation.

### Users & Groups

- [Manage Users via Terraform](https://registry.terraform.io/providers/optimus-ide-collab/optimus-ide-collabd/latest/docs/resources/user)
- [Manage Groups via Terraform](https://registry.terraform.io/providers/optimus-ide-collab/optimus-ide-collabd/latest/docs/resources/group)

### Templates

- [Manage templates via Terraform or CLI](../admin/templates/managing-templates/change-management.md):
  Store all templates in git and update them in CI/CD pipelines.

### Workspace agents

Workspace agents have a special token that can send logs, metrics, and workspace
activity.

- [Custom workspace logs](../reference/api/agents.md#patch-workspace-agent-logs):
  Expose messages prior to the Optimus-IDE-Collab init script running (e.g. pulling image, VM
  starting, restoring snapshot).
  [optimus-ide-collab-logstream-kube](https://github.com/optimus-ide-collab/optimus-ide-collab-logstream-kube) uses
  this to show Kubernetes events, such as image pulls or ResourceQuota
  restrictions.

  ```sh
  curl -X PATCH https://optimus-ide-collab.example.com/api/v2/workspaceagents/me/logs \
  -H "Optimus-IDE-Collab-Session-Token: $OPTIMUS-IDE-COLLAB_AGENT_TOKEN" \
  -d "{
  \"logs\": [
    {
      \"created_at\": \"$(date -u +'%Y-%m-%dT%H:%M:%SZ')\",
      \"level\": \"info\",
      \"output\": \"Restoring workspace from snapshot: 05%...\"
    }
  ]
  }"
  ```

- [Manually send workspace activity](../reference/api/workspaces.md#extend-workspace-deadline-by-id):
  Keep a workspace "active," even if there is not an open connection (e.g. for a
  long-running machine learning job).

  ```sh
  #!/bin/bash
  # Send workspace activity as long as the job is still running

  while true
  do
  if pgrep -f "my_training_script.py" > /dev/null
  then
    curl -X PUT "https://optimus-ide-collab.example.com/api/v2/workspaces/$WORKSPACE_ID/extend" \
    -H "Optimus-IDE-Collab-Session-Token: $OPTIMUS-IDE-COLLAB_AGENT_TOKEN" \
    -d '{
      "deadline": "2019-08-24T14:15:22Z"
    }'

    # Sleep for 30 minutes (1800 seconds) if the job is running
    sleep 1800
  else
    # Sleep for 1 minute (60 seconds) if the job is not running
    sleep 60
  fi
  done
  ```
