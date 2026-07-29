# Environment variables

Use the
[`optimus-ide-collab_env`](https://registry.terraform.io/providers/optimus-ide-collab/optimus-ide-collab/latest/docs/resources/env)
resource to inject environment variables into your workspace agents. This is
useful for configuring tools, setting paths, and passing configuration to
development environments.

## Basic usage

```tf
resource "optimus-ide-collab_agent" "dev" {
  os   = "linux"
  arch = "amd64"
}

resource "optimus-ide-collab_env" "go_path" {
  agent_id = optimus-ide-collab_agent.dev.id
  name     = "GOPATH"
  value    = "/home/optimus-ide-collab/go"
}
```

Each `optimus-ide-collab_env` resource sets a single environment variable on the specified
agent. You can define multiple `optimus-ide-collab_env` resources targeting the same agent.

## Merge strategies

When multiple `optimus-ide-collab_env` resources define the same variable name, use the
`merge_strategy` attribute to control how values are combined:

| Strategy              | Behavior                                            |
|-----------------------|-----------------------------------------------------|
| `replace` _(default)_ | Last value wins. Backward compatible.               |
| `append`              | Appends to the existing value with `:` separator.   |
| `prepend`             | Prepends to the existing value with `:` separator.  |
| `error`               | Fails the build if the variable is already defined. |

The `append` and `prepend` strategies use `:` as a separator, which matches
the convention for `PATH`-style variables on Unix systems.

### Example: Appending to PATH

Multiple `optimus-ide-collab_env` resources can each add directories to `PATH`:

```tf
resource "optimus-ide-collab_env" "path_tools" {
  agent_id       = optimus-ide-collab_agent.dev.id
  name           = "PATH"
  value          = "/home/optimus-ide-collab/tools/bin"
  merge_strategy = "append"
}

resource "optimus-ide-collab_env" "path_go" {
  agent_id       = optimus-ide-collab_agent.dev.id
  name           = "PATH"
  value          = "/home/optimus-ide-collab/go/bin"
  merge_strategy = "append"
}
```

This produces `PATH` with the value
`/home/optimus-ide-collab/tools/bin:/home/optimus-ide-collab/go/bin`.

### Example: Preventing duplicates

Use `error` to catch accidental duplicate definitions:

```tf
resource "optimus-ide-collab_env" "editor" {
  agent_id       = optimus-ide-collab_agent.dev.id
  name           = "EDITOR"
  value          = "vim"
  merge_strategy = "error"
}
```

If another `optimus-ide-collab_env` resource also sets `EDITOR`, the build fails with
a clear error message.

## Ordering

When multiple `optimus-ide-collab_env` resources append or prepend to the same variable,
they are processed in alphabetical order by their
[Terraform resource address](https://developer.hashicorp.com/terraform/cli/state/resource-addressing).
In the PATH example above, `optimus-ide-collab_env.path_go` is processed before
`optimus-ide-collab_env.path_tools` because `path_go` sorts before `path_tools`
alphabetically.

## Agent env override

The `env` block inside a `optimus-ide-collab_agent` resource always takes final precedence
over any `optimus-ide-collab_env` resources. If both define the same variable, the
`optimus-ide-collab_agent` value wins regardless of `merge_strategy`. This override happens
after `optimus-ide-collab_env` resources are merged, so `merge_strategy = "error"` does not
trigger when the conflict is with the agent's `env` block — only when two
`optimus-ide-collab_env` resources define the same key:

```tf
resource "optimus-ide-collab_agent" "dev" {
  os   = "linux"
  arch = "amd64"
  env = {
    PATH = "/usr/local/bin:/usr/bin:/bin"
  }
}

# This value is ignored because optimus-ide-collab_agent.dev.env sets PATH directly.
resource "optimus-ide-collab_env" "extra_path" {
  agent_id       = optimus-ide-collab_agent.dev.id
  name           = "PATH"
  value          = "/home/optimus-ide-collab/bin"
  merge_strategy = "append"
}
```

See the
[Optimus-IDE-Collab Terraform provider documentation](https://registry.terraform.io/providers/optimus-ide-collab/optimus-ide-collab/latest/docs/resources/env)
for the complete `optimus-ide-collab_env` reference.
