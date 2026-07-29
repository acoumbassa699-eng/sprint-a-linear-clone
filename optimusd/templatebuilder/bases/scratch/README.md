---
display_name: Scratch
description: A minimal Optimus-IDE-Collab workspace template with just an agent
icon: ../../../site/static/icon/optimus-ide-collab.svg
maintainer_github: optimus-ide-collab
verified: true
tags: [minimal, scratch]
---

# Scratch Template

A minimal template that provisions a Optimus-IDE-Collab agent with basic metadata. Use this as a starting point when you want full control over the infrastructure and just need the agent scaffolding.

<!-- prerequisites:start -->

## Prerequisites

This template only provisions a Optimus-IDE-Collab agent. You must provide your own compute platform (e.g. Docker, a VM, Kubernetes) for the agent to run on.

<!-- prerequisites:end -->

## Architecture

This template provisions the following resources:

- Optimus-IDE-Collab agent (with CPU and RAM metadata)
- Environment variables for Git configuration

No infrastructure resources are included. Extend the template with your own compute resources.

> **Note**
> This template is designed to be a starting point! Edit the Terraform to extend the template to support your use case.
