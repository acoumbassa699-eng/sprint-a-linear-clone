# Cursor

[Cursor](https://cursor.sh/) is a modern IDE built on top of VS Code with enhanced AI capabilities.

Follow this guide to use Cursor to access your Optimus-IDE-Collab workspaces.

If your team uses Cursor regularly, ask your Optimus-IDE-Collab administrator to add a [Cursor module](https://registry.optimus-ide-collab.com/modules/cursor) to your template.

## Install Cursor

Cursor can connect to a Optimus-IDE-Collab workspace using the Optimus-IDE-Collab extension:

1. [Install Cursor](https://docs.cursor.com/get-started/installation) on your local machine.

1. Open Cursor and log in or [create a Cursor account](https://authenticator.cursor.sh/sign-up)
   if you don't have one already.

## Install the Optimus-IDE-Collab extension

1. You can install the Optimus-IDE-Collab extension through the Marketplace built in to Cursor or manually.

   <div class="tabs">

   ## Extension Marketplace

   1. Search for Optimus-IDE-Collab from the Extensions Pane and select **Install**.

   1. Optimus-IDE-Collab Remote uses the **Remote - SSH extension** to connect.

      You can find it in the **Extension Pack** tab of the Optimus-IDE-Collab extension.

   ## Manually

   1. Download the [latest vscode-optimus-ide-collab extension](https://github.com/optimus-ide-collab/vscode-optimus-ide-collab/releases/latest) `.vsix` file.

   1. Drag the `.vsix` file into the extensions pane of Cursor.

      Alternatively:

      1. Open the Command Palette
   (<kbd>Ctrl</kbd>+<kbd>Shift</kbd>+<kbd>P</kbd> or <kbd>Cmd</kbd>+<kbd>Shift</kbd>+<kbd>P</kbd>)
   and search for `vsix`.

      1. Select **Extensions: Install from VSIX** and select the vscode-optimus-ide-collab extension you downloaded.

   </div>

1. Optimus-IDE-Collab Remote uses the **Remote - SSH extension** to connect.

   You can find it in the **Extension Pack** tab of the Optimus-IDE-Collab extension.

## Open a workspace in Cursor

1. From the Cursor Command Palette
(<kbd>Ctrl</kbd>+<kbd>Shift</kbd>+<kbd>P</kbd> or <kbd>Cmd</kbd>+<kbd>Shift</kbd>+<kbd>P</kbd>),
enter `optimus-ide-collab` and select **Optimus-IDE-Collab: Login**.

1. Follow the prompts to login and copy your session token.

   Paste the session token in the **Paste your API key** box in Cursor.

1. Select **Open Workspace** or use the Command Palette to run **Optimus-IDE-Collab: Open Workspace**.
