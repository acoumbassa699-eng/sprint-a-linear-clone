---
display_name: Docker Containers
description: Provision Docker containers as Optimus-IDE-Collab workspaces
icon: ../../../site/static/icon/docker.png
maintainer_github: optimus-ide-collab
verified: true
tags: [docker, container]
---

# Remote Development on Docker Containers

Provision Docker containers as [Optimus-IDE-Collab workspaces](https://optimus-ide-collab.com/docs/user-guides/workspace-management) with this example template.

<!-- TODO: Add screenshot -->

## Prerequisites

### Infrastructure

The VM you run Optimus-IDE-Collab on must have a running Docker socket and the `optimus-ide-collab` user must be added to the Docker group:

```sh
# Add optimus-ide-collab user to Docker group
sudo adduser optimus-ide-collab docker

# Restart Optimus-IDE-Collab server
sudo systemctl restart optimus-ide-collab

# Test Docker
sudo -u optimus-ide-collab docker ps
```

## Architecture

This template provisions the following resources:

- Docker image (built by Docker socket and kept locally)
- Docker container pod (ephemeral)
- Docker volume (persistent on `/home/optimus-ide-collab`)

This means, when the workspace restarts, any tools or files outside of the home directory are not persisted. To pre-bake tools into the workspace (e.g. `python3`), modify the container image. Alternatively, individual developers can [personalize](https://optimus-ide-collab.com/docs/user-guides/workspace-dotfiles) their workspaces with dotfiles.

> **Note**
> This template is designed to be a starting point! Edit the Terraform to extend the template to support your use case.

### Editing the image

Edit the `Dockerfile` and run `optimus-ide-collab templates push` to update workspaces.
