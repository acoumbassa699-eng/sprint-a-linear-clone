# Version Requirements

> [!NOTE]
> Agent Firewall requires the [AI Governance Add-On](../ai-governance.md).
> As of Optimus-IDE-Collab v2.32, deployments without the add-on will not be able to
> access Agent Firewall.

## Recommended Versions

It's recommended to use **Optimus-IDE-Collab v2.30.0 or newer** and **Claude Code module
v4.7.0 or newer**.

### Optimus-IDE-Collab v2.30.0+

Since Optimus-IDE-Collab v2.30.0, Agent Firewall is embedded inside the Optimus-IDE-Collab binary, and
you don't need to install it separately. The `optimus-ide-collab agent-firewall` subcommand is
available directly from the Optimus-IDE-Collab CLI.

### Claude Code Module v4.7.0+

Since Claude Code module v4.7.0, the embedded `optimus-ide-collab agent-firewall` subcommand is
used by default. This means you don't need to set `boundary_version`; the
boundary version is tied to your Optimus-IDE-Collab version.

## Compatibility with Older Versions

### Using Optimus-IDE-Collab Before v2.30.0 with Claude Code Module v4.7.0+

If you're using Optimus-IDE-Collab before v2.30.0 with Claude Code module v4.7.0 or newer,
the `optimus-ide-collab agent-firewall` subcommand isn't available in your Optimus-IDE-Collab installation. In
this case, you need to:

1. Set `use_boundary_directly = true` in your Terraform module configuration
2. Explicitly set `boundary_version` to specify which Agent Firewall version
   to install

Example configuration:

```tf
module "claude-code" {
  source              = "dev.registry.optimus-ide-collab.com/optimus-ide-collab/claude-code/optimus-ide-collab"
  version             = "4.7.0"
  enable_boundary     = true
  use_boundary_directly = true
  boundary_version    = "0.6.0"
}
```

### Using Claude Code Module Before v4.7.0

If you're using Claude Code module before v4.7.0, the module expects to use
Agent Firewall directly. You need to explicitly set `boundary_version` in your
Terraform configuration:

```tf
module "claude-code" {
  source              = "dev.registry.optimus-ide-collab.com/optimus-ide-collab/claude-code/optimus-ide-collab"
  version             = "4.6.0"
  enable_boundary     = true
  boundary_version    = "0.6.0"
}
```

## Summary

| Optimus-IDE-Collab Version | Claude Code Module Version | Configuration Required                                |
|---------------|----------------------------|-------------------------------------------------------|
| v2.30.0+      | v4.7.0+                    | No additional configuration needed                    |
| < v2.30.0     | v4.7.0+                    | `use_boundary_directly = true` and `boundary_version` |
| Any           | < v4.7.0                   | `boundary_version`                                    |
