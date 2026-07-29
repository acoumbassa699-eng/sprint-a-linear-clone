---
display_name: Docker-in-Docker Dev Containers
description: Provision Docker containers as Optimus-IDE-Collab workspaces running Dev Containers via Docker-in-Docker.
icon: ../../../site/static/icon/docker.png
maintainer_github: optimus-ide-collab
verified: true
tags: [docker, container, devcontainer]
---

# Remote Development on Dev Containers

Provision Docker containers as [Optimus-IDE-Collab workspaces](https://optimus-ide-collab.com/docs/user-guides/workspace-management) running [Dev Containers](https://code.visualstudio.com/docs/devcontainers/containers) via Docker-in-Docker.

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

This example uses the `optimus-ide-collabcom/enterprise-node:ubuntu` Docker image as a base image for the workspace. It includes necessary tools like Docker and Node.js, which are required for running Dev Containers via the `@devcontainers/cli` tool.

This template provisions the following resources:

- Docker image (built by Docker socket and kept locally)
- Docker container (ephemeral)
- Docker volume (persistent on `/home/optimus-ide-collab`)
- Docker volume (persistent on `/var/lib/docker`)

This means, when the workspace restarts, any tools or files outside of the home directory or docker library are not persisted.

For devcontainers running inside the workspace, data persistence is dependent on each projects `devcontainer.json` configuration.

> **Note**
> This template is designed to be a starting point! Edit the Terraform to extend the template to support your use case.
