---
display_name: Incus System Container with Docker
description: Develop in an Incus System Container with Docker using Incus
icon: ../../../site/static/icon/lxc.svg
maintainer_github: optimus-ide-collab
verified: true
tags: [incus, lxc, lxd]
---

# Incus System Container with Docker

Develop in an Incus System Container and run nested Docker containers using Incus.

## Architecture

This template uses the [Incus guest API](https://linuxcontainers.org/incus/docs/main/dev-incus/) (`/dev/incus/sock`) to deliver the Optimus-IDE-Collab agent token and URL into the container without any host filesystem coupling. This means:

- **The provisioner does not need to run on the Incus host.** There are no bind mounts or local file writes. All configuration is passed via Incus `user.*` config keys and read from inside the container at runtime.
- **The agent binary is downloaded automatically.** The standard Optimus-IDE-Collab init script fetches the correct binary from the Optimus-IDE-Collab server on every boot, keeping it in sync with the server version.
- **The agent token is refreshed on every start.** Terraform updates the `user.optimus-ide-collab_agent_token` config key each workspace start. A watcher service inside the container listens for config changes via the guest API events endpoint and restarts the agent when a new token arrives.

### Boot sequence

1. **First boot (cloud-init):** Creates the workspace user, writes the bootstrap scripts and systemd units, installs `curl` and `git`, and enables the services. Cloud-init only runs once.
2. **Every boot (systemd):**
   - `optimus-ide-collab-agent-config.service` (oneshot) reads `OPTIMUS-IDE-COLLAB_AGENT_TOKEN` and `OPTIMUS-IDE-COLLAB_AGENT_URL` from the Incus guest API and writes them to `/opt/optimus-ide-collab/init.env`.
   - `optimus-ide-collab-agent.service` loads the env file and runs the Optimus-IDE-Collab init script, which downloads the agent binary and starts it.
   - `optimus-ide-collab-agent-watcher.service` streams config change events from the guest API. If the Incus provider updates the token *after* the container has already booted (a known provider ordering issue), the watcher detects the change, re-fetches the config, and restarts the agent.

### Packages

Essential packages (`curl`, `git`) are installed via cloud-init on first boot, before the agent starts. Additional packages (e.g. `docker.io`) are installed via a non-blocking [`optimus-ide-collab_script`](https://registry.terraform.io/providers/optimus-ide-collab/optimus-ide-collab/latest/docs/resources/script) that runs on each workspace start. It does not block login; users can connect to the workspace immediately while packages install in the background. On subsequent starts, it detects packages are already installed and skips the installation.

## Prerequisites

1. Install [Incus](https://linuxcontainers.org/incus/) on a machine reachable by the Optimus-IDE-Collab provisioner.
2. Allow Optimus-IDE-Collab to access the Incus socket.

   - If you're running Optimus-IDE-Collab as a system service, run `sudo usermod -aG incus-admin optimus-ide-collab` and restart the Optimus-IDE-Collab service.
   - If you're running Optimus-IDE-Collab as a Docker Compose service, get the group ID of the `incus-admin` group by running `getent group incus-admin` and add the following to your `compose.yaml` file:

     ```yaml
     services:
       optimus-ide-collab:
         volumes:
           - /var/lib/incus/unix.socket:/var/lib/incus/unix.socket
         group_add:
           - 996 # Replace with the group ID of the `incus-admin` group
     ```

3. Create a storage pool named `optimus-ide-collab` by running `incus storage create optimus-ide-collab btrfs` (or use another [supported driver](https://linuxcontainers.org/incus/docs/main/reference/storage_drivers/)).

## Usage

> **Note:** This template requires a container image with cloud-init installed, such as `images:debian/13/cloud` or `images:ubuntu/24.04/cloud`. Images are pulled automatically from the [Linux Containers image server](https://images.linuxcontainers.org/).

1. Run `optimus-ide-collab templates push --directory .` from this directory.
2. Create a workspace from the template in the Optimus-IDE-Collab UI.

## Parameters

| Parameter          | Description                                                                                | Default                  |
|--------------------|--------------------------------------------------------------------------------------------|--------------------------|
| **Image**          | Container image with cloud-init. Options: Debian 13, Debian 12, Ubuntu 24.04, Ubuntu 22.04 | `images:debian/13/cloud` |
| **CPU**            | Number of CPUs (1-8)                                                                       | `1`                      |
| **Memory**         | Memory in GB (1-16)                                                                        | `2`                      |
| **Storage pool**   | Incus storage pool name                                                                    | `optimus-ide-collab`                  |
| **Git repository** | Clone a git repo inside the workspace                                                      | *(empty)*                |

## Extending this template

See the [lxc/incus](https://registry.terraform.io/providers/lxc/incus/latest/docs) Terraform provider documentation to add the following features to your Optimus-IDE-Collab template:

- Remote Incus hosts (HTTPS)
- Additional volume mounts
- Custom networks
- GPU passthrough
- More

We also welcome contributions!
