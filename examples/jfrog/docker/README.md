---
name: JFrog and Docker
description: Develop inside Docker containers using your local daemon
tags: [local, docker, jfrog]
icon: /icon/docker.png
---

# Docker

To get started, run `optimus-ide-collab templates init`. When prompted, select this template.
Follow the on-screen instructions to proceed.

## Editing the image

Edit the `Dockerfile` and run `optimus-ide-collab templates push` to update workspaces.

## code-server

`code-server` is installed via the `startup_script` argument in the `optimus-ide-collab_agent`
resource block. The `optimus-ide-collab_app` resource is defined to access `code-server` through
the dashboard UI over `localhost:13337`.

## Next steps

Check out our [Docker](../../templates/docker/) template for a more fully featured Docker
example.
