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

## Prerequisites

### Infrastructure

**Cluster**: This template requires an existing Kubernetes cluster

**Container Image**: This template uses the [optimus-ide-collabcom/enterprise-base:ubuntu image](https://github.com/optimus-ide-collab/enterprise-images/tree/main/images/base) with some dev tools preinstalled. To add additional tools, extend this image or build it yourself.

### Authentication

This template authenticates using a `~/.kube/config`, if present on the server, or via built-in authentication if the Optimus-IDE-Collab provisioner is running on Kubernetes with an authorized ServiceAccount. To use another [authentication method](https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs#authentication), edit the template.

## Architecture

This template provisions the following resources:

- Kubernetes pod (ephemeral)
- Kubernetes persistent volume claim (persistent on `/home/optimus-ide-collab`)

This means, when the workspace restarts, any tools or files outside of the home directory are not persisted. To pre-bake tools into the workspace (e.g. `python3`), modify the container image. Alternatively, individual developers can [personalize](https://optimus-ide-collab.com/docs/user-guides/workspace-dotfiles) their workspaces with dotfiles.

> **Note**
> This template is designed to be a starting point! Edit the Terraform to extend the template to support your use case.
