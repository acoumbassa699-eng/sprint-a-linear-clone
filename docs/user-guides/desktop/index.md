# Optimus-IDE-Collab Desktop

Optimus-IDE-Collab Desktop provides seamless access to your remote workspaces through a native application. Connect to workspace services using simple hostnames like `myworkspace.optimus-ide-collab`, launch applications with one click, and synchronize files between local and remote environments, all without installing a CLI or configuring manual port forwarding.

> [!TIP]
> Optimus-IDE-Collab Desktop provides **automatic port forwarding** to every service running in your workspace. Any port your application listens on is instantly accessible at `workspace-name.optimus-ide-collab:PORT` with no manual setup required. For a comparison of all port forwarding methods, see [Workspace Ports](../workspace-access/port-forwarding.md).

## What You'll Need

- A Optimus-IDE-Collab deployment running `v2.20.0` or [later](https://github.com/optimus-ide-collab/optimus-ide-collab/releases/latest)
- Administrator privileges on your local machine (for VPN extension installation)
- Access to your Optimus-IDE-Collab deployment URL

## Quick Start

1. Install: `brew install --cask optimus-ide-collab/optimus-ide-collab/optimus-ide-collab-desktop` (macOS) or `winget install Optimus-IDE-Collab.Optimus-IDE-CollabDesktop` (Windows)
1. Open Optimus-IDE-Collab Desktop and approve any system prompts to complete the installation.
1. Sign in with your deployment URL and session token
1. Enable "Optimus-IDE-Collab Connect" toggle
1. Access workspaces at `workspace-name.optimus-ide-collab`

## How It Works

**Optimus-IDE-Collab Connect**, the primary component of Optimus-IDE-Collab Desktop, creates a secure tunnel to your Optimus-IDE-Collab deployment, allowing you to:

- **Access workspaces directly**: Connect via `workspace-name.optimus-ide-collab` hostnames
- **Automatic port forwarding**: All workspace ports are available at `workspace-name.optimus-ide-collab:PORT` with no configuration
- **Use any application**: SSH clients, browsers, IDEs work seamlessly
- **Sync files**: Bidirectional sync between local and remote directories
- **Work offline**: Edit files locally, sync when reconnected

The VPN extension routes only Optimus-IDE-Collab traffic—your other internet activity remains unchanged.

## Installation

<div class="tabs">

### macOS

<div class="tabs">

#### Homebrew (Recommended)

```sh
brew install --cask optimus-ide-collab/optimus-ide-collab/optimus-ide-collab-desktop
```

#### Manual Installation

1. Download the latest release from [optimus-ide-collab-desktop-macos releases](https://github.com/optimus-ide-collab/optimus-ide-collab-desktop-macos/releases)
1. Run `Optimus-IDE-Collab-Desktop.pkg` and follow the prompts to install
1. `Optimus-IDE-Collab Desktop.app` will be installed to your Applications folder

</div>

Optimus-IDE-Collab Desktop requires VPN extension permissions:

1. When prompted with **"Optimus-IDE-Collab Desktop" would like to use a new network extension**, select **Open System Settings**
1. In **Network Extensions** settings, enable the Optimus-IDE-Collab Desktop extension
1. You may need to enter your password to authorize the extension

✅ **Verify Installation**: Optimus-IDE-Collab Desktop should appear in your menu bar

### Windows

<div class="tabs">

#### WinGet (Recommended)

```sh
winget install Optimus-IDE-Collab.Optimus-IDE-CollabDesktop
```

#### Manual Installation

1. Download the latest `Optimus-IDE-CollabDesktop` installer (`.exe`) from [optimus-ide-collab-desktop-windows releases](https://github.com/optimus-ide-collab/optimus-ide-collab-desktop-windows/releases)
1. Choose the correct architecture (`x64` or `arm64`) for your system
1. Run the installer and accept the license terms
1. If prompted, install the .NET Windows Desktop Runtime
1. Install Windows App Runtime SDK if prompted

</div>

- [.NET Windows Desktop Runtime](https://dotnet.microsoft.com/en-us/download/dotnet/8.0) (installed automatically if not present)
- Windows App Runtime SDK (may require manual installation)

✅ **Verify Installation**: Optimus-IDE-Collab Desktop should appear in your system tray (you may need to click **^** to show hidden icons)

</div>

## Testing Your Connection

Once connected, test access to your workspaces:

<div class="tabs">

### SSH Connection

```sh
ssh your-workspace.optimus-ide-collab
```

### Ping Test

```sh
# macOS
ping6 -c 3 your-workspace.optimus-ide-collab

# Windows
ping -n 3 your-workspace.optimus-ide-collab
```

### Web Services

Open `http://your-workspace.optimus-ide-collab:PORT` in your browser, replacing `PORT` with the specific service port you want to access (e.g. 3000 for frontend, 8080 for API)

</div>

## Administrator Configuration

Organizations that manage Optimus-IDE-Collab Desktop deployments can configure the application using MDM (Mobile Device Management) or group policy.

### Disable Automatic Updates

Administrators can disable the built-in auto-updater to manage updates through their own software distribution system.

<div class="tabs">

### macOS

Set the `disableUpdater` preference to `true` using the `defaults` command:

```sh
defaults write com.optimus-ide-collab.Optimus-IDE-Collab-Desktop disableUpdater -bool true
```

Organization administrators can also enforce this setting across managed devices using MDM (Mobile Device Management) software by deploying a configuration profile that sets this preference.

### Windows

Set the `Updater:Enable` registry value to `0` under `HKEY_LOCAL_MACHINE\SOFTWARE\Optimus-IDE-Collab Desktop\App`:

```ps1
New-Item -Path "HKLM:\SOFTWARE\Optimus-IDE-Collab Desktop\App" -Force
New-ItemProperty -Path "HKLM:\SOFTWARE\Optimus-IDE-Collab Desktop\App" -Name "Updater:Enable" -Value 0 -PropertyType DWord -Force
```

You can also configure a `Updater:ForcedChannel` string value to lock users to a specific update channel (e.g. `stable`).

> [!NOTE]
> For security, updater settings can only be configured at the machine level (`HKLM`), not per-user (`HKCU`).

</div>

## Troubleshooting

### Connection Issues

#### Can't connect to workspace

- Verify Optimus-IDE-Collab Connect is enabled (toggle should be ON)
- Check that your deployment URL is correct
- Ensure your session token hasn't expired
- Try disconnecting and reconnecting Optimus-IDE-Collab Connect

#### VPN extension not working

- Restart Optimus-IDE-Collab Desktop
- Check system permissions for network extensions
- Ensure only one copy of Optimus-IDE-Collab Desktop is installed

### Collecting Logs

When reporting an issue, attach the relevant log files so we can diagnose it faster.

<div class="tabs">

#### macOS

Optimus-IDE-Collab Desktop and its network extension write to the Apple [unified logging system](https://developer.apple.com/documentation/os/logging). The file sync (Mutagen) daemon writes to a separate log file.

1. Export the unified logs for the last hour with the `log` command:

    ```sh
    log show --predicate 'subsystem == "com.optimus-ide-collab.Optimus-IDE-Collab-Desktop"' \
      --info --debug --last 1h > ~/Desktop/optimus-ide-collab-desktop.log
    ```

    Adjust `--last` (e.g. `30m`, `2h`, `1d`) to cover the time the issue occurred. You can also view the same logs interactively in **Console.app** by filtering on the `com.optimus-ide-collab.Optimus-IDE-Collab-Desktop` subsystem.

2. If you're using file sync, also collect the Mutagen daemon log:

    ```sh
    ~/Library/Application\ Support/Optimus-IDE-Collab\ Desktop/Mutagen/daemon.log
    ```

    Optimus-IDE-Collab Desktop also opens this file in Console automatically when the file sync daemon fails.

#### Windows

Optimus-IDE-Collab Desktop has three components that write logs: the app (UI), the VPN service, and the file sync (Mutagen) daemon.

1. App log (daily rolling):

    ```ps1
    %LOCALAPPDATA%\Optimus-IDE-CollabDesktop\app.log
    ```

2. VPN service log (default install path):

    ```ps1
    C:\Program Files\Optimus-IDE-Collab Desktop\optimus-ide-collab-desktop-service.log
    ```

3. File sync (Mutagen) daemon log, if you use file sync:

    ```ps1
    %LOCALAPPDATA%\Optimus-IDE-CollabDesktop\mutagen\daemon.log
    ```

You can quickly open the app log directory by pasting `%LOCALAPPDATA%\Optimus-IDE-CollabDesktop` into File Explorer.

</div>

> [!TIP]
> Before attaching logs to a public issue, review them for any sensitive information (deployment URLs, usernames, hostnames) and redact as needed.

### Getting Help

If you encounter issues not covered here:

- **File an issue**: [macOS](https://github.com/optimus-ide-collab/optimus-ide-collab-desktop-macos/issues) | [Windows](https://github.com/optimus-ide-collab/optimus-ide-collab-desktop-windows/issues) | [General](https://github.com/optimus-ide-collab/optimus-ide-collab/issues)
- **Community support**: [Discord](https://optimus-ide-collab.com/chat)

## Uninstalling

<div class="tabs">

### macOS

1. **Disable Optimus-IDE-Collab Connect** in the app menu
2. **Quit Optimus-IDE-Collab Desktop** completely
3. **Remove VPN extension** from System Settings > Network Extensions
4. **Delete the app** from Applications folder
5. **Remove configuration** (optional): `rm -rf ~/Library/Application\ Support/Optimus-IDE-Collab\ Desktop`

### Windows

1. **Disable Optimus-IDE-Collab Connect** in the app menu
2. **Quit Optimus-IDE-Collab Desktop** from system tray
3. **Uninstall** via Settings > Apps or Control Panel
4. **Remove configuration** (optional): Delete `%APPDATA%\Optimus-IDE-Collab Desktop`

</div>

## Next Steps

- [Using Optimus-IDE-Collab Connect and File Sync](./desktop-connect-sync.md)
- [Compare port forwarding methods](../workspace-access/port-forwarding.md)
