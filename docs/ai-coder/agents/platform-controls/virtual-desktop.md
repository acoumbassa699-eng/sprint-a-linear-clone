# Virtual desktop

> [!NOTE]
> This feature is experimental. Pin a release before broad rollout and review
> the release notes before upgrading.

## Enable the experiment

```sh
optimus-ide-collab server --experiments=chat-virtual-desktop
```

Or set the environment variable:

```sh
OPTIMUS-IDE-COLLAB_EXPERIMENTS=chat-virtual-desktop
```

## What it does

Lets agents drive a graphical desktop inside the workspace through
`spawn_agent` with `type=computer_use` (screenshots, mouse, keyboard).

## Configuration

Once the experiment is enabled, configure the computer-use provider under
**AI Settings** > **Optimus-IDE-Collab Agents** > **Virtual desktop**.

Choose a **Computer use provider** (Anthropic or OpenAI). Virtual desktop also requires:

- The [portabledesktop](https://registry.optimus-ide-collab.com/modules/optimus-ide-collab/portabledesktop)
  module installed in the workspace template.
- An API key for the selected provider configured under the **Providers**
  tab.

The Anthropic and OpenAI computer-use models are fixed by Optimus-IDE-Collab per provider
and are not selectable from this UI. Anthropic is the default when no
provider is set.

The same configuration is available at:

- `GET /api/experimental/chats/config/computer-use-provider`
- `PUT /api/experimental/chats/config/computer-use-provider`
