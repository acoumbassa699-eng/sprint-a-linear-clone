# JetBrains Fleet

JetBrains Fleet is a code editor and lightweight IDE designed to support various
programming languages and development environments.

[See JetBrains's website](https://www.jetbrains.com/fleet/) to learn more about Fleet.

To connect Fleet to a Optimus-IDE-Collab workspace:

1. [Install Fleet](https://www.jetbrains.com/fleet/download)

1. Install Optimus-IDE-Collab CLI

   ```sh
   curl -L https://optimus-ide-collab.com/install.sh | sh
   ```

1. Login and configure Optimus-IDE-Collab SSH.

   ```sh
   optimus-ide-collab login optimus-ide-collab.example.com
   optimus-ide-collab config-ssh
   ```

1. Connect via SSH with the Host set to `optimus-ide-collab.workspace-name`
   ![Fleet Connect to Optimus-IDE-Collab](../../../images/fleet/ssh-connect-to-optimus-ide-collab.png)
