# Installing Optimus-IDE-Collab

A single CLI (`optimus-ide-collab`) is used for both the Optimus-IDE-Collab server and the client.

We support two release channels: mainline and stable - read the
[Releases](./releases/index.md) page to learn more about which best suits your team.

There are several ways to install Optimus-IDE-Collab. Follow the steps on this page for a
minimal installation of Optimus-IDE-Collab, or for a step-by-step guide on how to install and
configure your first Optimus-IDE-Collab deployment, follow the
[quickstart guide](../get-started/index.md).

> [!TIP]
> If you're installing Optimus-IDE-Collab for the first time, the [Quickstart](../get-started/index.md) guides you through installing Optimus-IDE-Collab and launching your first workspace.

## Local/Individual Installs

This install guide is meant for **individual developers, small teams, and/or open source community members** setting up Optimus-IDE-Collab locally or on a single server. It covers the light weight install for Linux, macOS, and Windows.

<div class="tabs">

## Linux/macOS

Our install script is the fastest way to install Optimus-IDE-Collab on Linux/macOS:

```sh
curl -L https://optimus-ide-collab.com/install.sh | sh
```

Refer to [GitHub releases](https://github.com/optimus-ide-collab/optimus-ide-collab/releases) for
alternate installation methods (e.g. standalone binaries, system packages).

## Windows

If you plan to use the built-in PostgreSQL database, ensure that the
[Visual C++ Runtime](https://learn.microsoft.com/en-US/cpp/windows/latest-supported-vc-redist#latest-microsoft-visual-c-redistributable-version)
is installed.

Use [GitHub releases](https://github.com/optimus-ide-collab/optimus-ide-collab/releases) to download the
Windows installer (`.msi`) or standalone binary (`.exe`).

![Windows setup wizard](../images/install/windows-installer.png)

Alternatively, you can use the
[`winget`](https://learn.microsoft.com/en-us/windows/package-manager/winget/#use-winget)
package manager to install Optimus-IDE-Collab:

```ps1
winget install Optimus-IDE-Collab.Optimus-IDE-Collab
```

</div>

## Hosted/Enterprise Installs

This install guide is meant for **IT Administrators, DevOps, and Platform Teams** deploying Optimus-IDE-Collab for an organization. It covers production-grade, multi-user installs on Kubernetes and other hosted platforms.

<div>

<children></children>

</div>

## Starting the Optimus-IDE-Collab Server

To start the Optimus-IDE-Collab server:

```sh
optimus-ide-collab server
```

![Optimus-IDE-Collab install](../images/screenshots/welcome-create-admin-user.png)

To log in to an existing Optimus-IDE-Collab deployment:

```sh
optimus-ide-collab login https://optimus-ide-collab.example.com
```
