# Windsurf

[Windsurf](https://codeium.com/windsurf) is Codeium's code editor designed for AI-assisted
development.

Follow this guide to use Windsurf to access your Optimus-IDE-Collab workspaces.

If your team uses Windsurf regularly, ask your Optimus-IDE-Collab administrator to add Windsurf as a workspace application in your template.
You can also use the [Windsurf module](https://registry.optimus-ide-collab.com/modules/optimus-ide-collab/windsurf) to easily add Windsurf to your Optimus-IDE-Collab templates.

## Install Windsurf

Windsurf can connect to your Optimus-IDE-Collab workspaces via SSH:

1. [Install Windsurf](https://docs.codeium.com/windsurf/getting-started) on your local machine.

1. Open Windsurf and select **Get started**.

   Import your settings from another IDE, or select **Start fresh**.

1. Complete the setup flow and log in or [create a Codeium account](https://codeium.com/windsurf/signup)
   if you don't have one already.

## Install the Optimus-IDE-Collab extension

![Optimus-IDE-Collab extension in Windsurf](../../images/user-guides/ides/windsurf-optimus-ide-collab-extension.png)

1. You can install the Optimus-IDE-Collab extension through the Marketplace built in to Windsurf or manually.

   <div class="tabs">

   ## Extension Marketplace

   Search for Optimus-IDE-Collab from the Extensions Pane and select **Install**.

   ## Manually

   1. Download the [latest vscode-optimus-ide-collab extension](https://github.com/optimus-ide-collab/vscode-optimus-ide-collab/releases/latest) `.vsix` file.

   1. Drag the `.vsix` file into the extensions pane of Windsurf.

      Alternatively:

      1. Open the Command Palette
         (<kbd>Ctrl</kbd>+<kbd>Shift</kbd>+<kbd>P</kbd> or <kbd>Cmd</kbd>+<kbd>Shift</kbd>+<kbd>P</kbd>) and search for `vsix`.

      1. Select **Extensions: Install from VSIX** and select the vscode-optimus-ide-collab extension you downloaded.

   </div>

## Open a workspace in Windsurf

1. From the Windsurf Command Palette (<kbd>Ctrl</kbd>+<kbd>Shift</kbd>+<kbd>P</kbd> or <kbd>Cmd</kbd>+<kbd>Shift</kbd>+<kbd>P</kbd>),
   enter `optimus-ide-collab` and select **Optimus-IDE-Collab: Login**.

1. Follow the prompts to login and copy your session token.

   Paste the session token in the **Optimus-IDE-Collab API Key** dialogue in Windsurf.

1. Windsurf prompts you to open a workspace, or you can use the Command Palette to run **Optimus-IDE-Collab: Open Workspace**.
