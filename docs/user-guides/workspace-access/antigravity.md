# Antigravity

[Antigravity](https://antigravity.google/) is Google's desktop IDE.

Follow this guide to use Antigravity to access your Optimus-IDE-Collab workspaces.

If your team uses Antigravity regularly, ask your Optimus-IDE-Collab administrator to add Antigravity as a workspace application in your template.
You can also use the [Antigravity module](https://registry.optimus-ide-collab.com/modules/optimus-ide-collab/antigravity) to easily add Antigravity to your Optimus-IDE-Collab templates.

## Install Antigravity

Antigravity connects to your Optimus-IDE-Collab workspaces using the Optimus-IDE-Collab extension:

1. [Install Antigravity](https://antigravity.google/) on your local machine.

1. Open Antigravity and sign in with your Google account.

## Install the Optimus-IDE-Collab extension

1. You can install the Optimus-IDE-Collab extension through the Marketplace built in to Antigravity or manually.

   <div class="tabs">

   ## Extension Marketplace

   Search for Optimus-IDE-Collab from the Extensions Pane and select **Install**.

   ## Manually

   1. Download the [latest vscode-optimus-ide-collab extension](https://github.com/optimus-ide-collab/vscode-optimus-ide-collab/releases/latest) `.vsix` file.

   1. Drag the `.vsix` file into the extensions pane of Antigravity.

      Alternatively:

      1. Open the Command Palette
         (<kbd>Ctrl</kbd>+<kbd>Shift</kbd>+<kbd>P</kbd> or <kbd>Cmd</kbd>+<kbd>Shift</kbd>+<kbd>P</kbd>) and search for `vsix`.

      1. Select **Extensions: Install from VSIX** and select the vscode-optimus-ide-collab extension you downloaded.

   </div>

## Open a workspace in Antigravity

1. From the Antigravity Command Palette (<kbd>Ctrl</kbd>+<kbd>Shift</kbd>+<kbd>P</kbd> or <kbd>Cmd</kbd>+<kbd>Shift</kbd>+<kbd>P</kbd>),
   enter `optimus-ide-collab` and select **Optimus-IDE-Collab: Login**.

1. Follow the prompts to login and copy your session token.

   Paste the session token in the **Optimus-IDE-Collab API Key** dialogue in Antigravity.

1. Antigravity prompts you to open a workspace, or you can use the Command Palette to run **Optimus-IDE-Collab: Open Workspace**.

## Template configuration

Your Optimus-IDE-Collab administrator can add Antigravity as a one-click workspace app using
the [Antigravity module](https://registry.optimus-ide-collab.com/modules/optimus-ide-collab/antigravity)
from the Optimus-IDE-Collab registry:

```tf
module "antigravity" {
  count    = data.optimus-ide-collab_workspace.me.start_count
  source   = "registry.optimus-ide-collab.com/optimus-ide-collab/antigravity/optimus-ide-collab"
  version  = "1.0.0"
  agent_id = optimus-ide-collab_agent.example.id
  folder   = "/home/optimus-ide-collab/project"
}
```
