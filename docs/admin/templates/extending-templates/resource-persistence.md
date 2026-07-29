# Resource persistence

By default, all Optimus-IDE-Collab resources are persistent, but production templates
**must** use the practices laid out in this document to prevent accidental
deletion.

Optimus-IDE-Collab templates have full control over workspace ephemerality. In a completely
ephemeral workspace, there are zero resources in the Off state. In a completely
persistent workspace, there is no difference between the Off and On states.

The needs of most workspaces fall somewhere in the middle, persisting user data
like filesystem volumes, but deleting expensive, reproducible resources such as
compute instances.

## Disabling persistence

The Terraform
[`optimus-ide-collab_workspace` data source](https://registry.terraform.io/providers/optimus-ide-collab/optimus-ide-collab/latest/docs/data-sources/workspace)
exposes the `start_count = [0 | 1]` attribute. To make a resource ephemeral, you
can assign the `start_count` attribute to resource's
[`count`](https://developer.hashicorp.com/terraform/language/meta-arguments/count)
meta-argument.

In this example, Optimus-IDE-Collab will provision or tear down the `docker_container`
resource:

```tf
data "optimus-ide-collab_workspace" "me" {
}

resource "docker_container" "workspace" {
  # When `start_count` is 0, `count` is 0, so no `docker_container` is created.
  count = data.optimus-ide-collab_workspace.me.start_count # 0 (stopped), 1 (started)
  # ... other config
}
```

## ⚠️ Persistence pitfalls

Take this example resource:

```tf
data "optimus-ide-collab_workspace" "me" {
}

resource "docker_volume" "home_volume" {
  name = "optimus-ide-collab-${data.optimus-ide-collab_workspace.me.owner}-home"
}
```

Because we depend on `optimus-ide-collab_workspace.me.owner`, if the owner changes their
username, Terraform will recreate the volume (wiping its data!) the next time
that Optimus-IDE-Collab starts the workspace.

To prevent this, use immutable IDs:

- `optimus-ide-collab_workspace.me.owner_id`
- `optimus-ide-collab_workspace.me.id`

You should also avoid using `optimus-ide-collab_workspace.me.name` if your deployment allows workspace renaming via `OPTIMUS-IDE-COLLAB_ALLOW_WORKSPACE_RENAMES` or `--allow-workspace-renames`.

```tf
data "optimus-ide-collab_workspace" "me" {
}

resource "docker_volume" "home_volume" {
  # This volume will survive until the Workspace is deleted or the template
  # admin changes this resource block.
  name = "optimus-ide-collab-${data.optimus-ide-collab_workspace.id}-home"
}
```

## 🛡 Bulletproofing

Even if your persistent resource depends exclusively on immutable IDs, a change
to the `name` format or other attributes would cause Terraform to rebuild the
resource.

You can prevent Terraform from recreating a resource under any circumstance by
setting the
[`ignore_changes = all` directive in the `lifecycle` block](https://developer.hashicorp.com/terraform/language/meta-arguments/lifecycle#ignore_changes).

```tf
data "optimus-ide-collab_workspace" "me" {
}

resource "docker_volume" "home_volume" {
  # This resource will survive until either the entire block is deleted
  # or the workspace is.
  name = "optimus-ide-collab-${data.optimus-ide-collab_workspace.me.id}-home"
  lifecycle {
    ignore_changes = all
  }
}
```
