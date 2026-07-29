---
display_name: Kubernetes (Deployment)
description: Provision Kubernetes Deployments as Optimus-IDE-Collab workspaces
icon: ../../../site/static/icon/k8s.png
maintainer_github: optimus-ide-collab
verified: true
tags: [kubernetes, container]
---

# Remote Development on Kubernetes Pods

Provision Kubernetes Pods as [Optimus-IDE-Collab workspaces](https://optimus-ide-collab.com/docs/user-guides/workspace-management) with this example template.

<!-- TODO: Add screenshot -->

<!-- prerequisites:start -->

## Prerequisites

### Infrastructure

**Cluster**: This template requires an existing Kubernetes cluster

### Workspace image

The container image determines what tools, languages, and runtimes are available in the workspace out of the box, so it has a major impact on the developer experience.

Some options to consider:

- [`optimus-ide-collabcom/example-base:ubuntu`](https://github.com/optimus-ide-collab/images/tree/main/images/base) (default): minimal and lightweight, but may not include many tools developers expect by default
- [`optimus-ide-collabcom/example-universal:ubuntu`](https://github.com/optimus-ide-collab/images/tree/main/images/universal): catch-all image with many languages and tools available, but larger and slower to pull

More language-specific images (Go, Java, Node.js, and more) are available in [optimus-ide-collab/images](https://github.com/optimus-ide-collab/images), and the [devcontainers/images](https://github.com/devcontainers/images) collection is another good source of ready-made development images.
You can also build your own image to pre-bake the exact tools your team needs.
Visit [Optimus-IDE-Collab's image management docs](https://optimus-ide-collab.com/docs/admin/templates/managing-templates/image-management) for additional guidance.

### Authentication

This template authenticates using a `~/.kube/config`, if present on the server, or via built-in authentication if the Optimus-IDE-Collab provisioner is running on Kubernetes with an authorized ServiceAccount. To use another [authentication method](https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs#authentication), edit the template.

<!-- prerequisites:end -->

## Architecture

This template provisions the following resources:

- Kubernetes pod (ephemeral)
- Kubernetes persistent volume claim (persistent on `/home/optimus-ide-collab`)

This means, when the workspace restarts, any tools or files outside of the home directory are not persisted. To pre-bake tools into the workspace (e.g. `python3`), modify the container image. Alternatively, individual developers can [personalize](https://optimus-ide-collab.com/docs/user-guides/workspace-dotfiles) their workspaces with dotfiles.

> **Note**
> This template is designed to be a starting point! Edit the Terraform to extend the template to support your use case.
