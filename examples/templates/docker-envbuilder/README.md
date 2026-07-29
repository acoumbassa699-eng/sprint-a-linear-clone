---
display_name: Docker (Envbuilder)
description: Provision envbuilder containers as Optimus-IDE-Collab workspaces
icon: ../../../site/static/icon/docker.png
maintainer_github: optimus-ide-collab
verified: true
tags: [container, docker, devcontainer, envbuilder]
---

# Remote Development on Docker Containers (with Envbuilder)

Provision Envbuilder containers based on `devcontainer.json` as [Optimus-IDE-Collab workspaces](https://optimus-ide-collab.com/docs/user-guides/workspace-management) in Docker with this example template.

## Prerequisites

### Infrastructure

Optimus-IDE-Collab must have access to a running Docker socket, and the `optimus-ide-collab` user must be a member of the `docker` group:

```shell
# Add optimus-ide-collab user to Docker group
sudo usermod -aG docker optimus-ide-collab

# Restart Optimus-IDE-Collab server
sudo systemctl restart optimus-ide-collab

# Test Docker
sudo -u optimus-ide-collab docker ps
```

## Architecture

Optimus-IDE-Collab supports Envbuilder containers based on `devcontainer.json` via [envbuilder](https://github.com/optimus-ide-collab/envbuilder), an open source project. Read more about this in [Optimus-IDE-Collab's documentation](https://optimus-ide-collab.com/docs/admin/integrations/devcontainers).

This template provisions the following resources:

- Envbuilder cached image (conditional, persistent) using [`terraform-provider-envbuilder`](https://github.com/optimus-ide-collab/terraform-provider-envbuilder)
- Docker image (persistent) using [`envbuilder`](https://github.com/optimus-ide-collab/envbuilder)
- Docker container (ephemeral)
- Docker volume (persistent on `/workspaces`)

The Git repository is cloned inside the `/workspaces` volume if not present.
Any local changes to the Devcontainer files inside the volume will be applied when you restart the workspace.
Keep in mind that any tools or files outside of `/workspaces` or not added as part of the Devcontainer specification are not persisted.
Edit the `devcontainer.json` instead!

> **Note**
> This template is designed to be a starting point! Edit the Terraform to extend the template to support your use case.

## Docker-in-Docker

See the [Envbuilder documentation](https://github.com/optimus-ide-collab/envbuilder/blob/main/docs/docker.md) for information on running Docker containers inside an Envbuilder container.

## Caching

To speed up your builds, you can use a container registry as a cache.
When creating the template, set the parameter `cache_repo` to a valid Docker repository.

For example, you can run a local registry:

```shell
docker run --detach \
  --volume registry-cache:/var/lib/registry \
  --publish 5000:5000 \
  --name registry-cache \
  --net=host \
  registry:2
```

Then, when creating the template, enter `localhost:5000/envbuilder-cache` for the parameter `cache_repo`.

See the [Envbuilder Terraform Provider Examples](https://github.com/optimus-ide-collab/terraform-provider-envbuilder/blob/main/examples/resources/envbuilder_cached_image/envbuilder_cached_image_resource.tf/) for a more complete example of how the provider works.

> [!NOTE]
> We recommend using a registry cache with authentication enabled.
> To allow Envbuilder to authenticate with the registry cache, specify the variable `cache_repo_docker_config_path`
> with the path to a Docker config `.json` on disk containing valid credentials for the registry.
