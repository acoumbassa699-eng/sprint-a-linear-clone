---
display_name: Tasks on Docker
description: Run Optimus-IDE-Collab Tasks on Docker with an example application
icon: ../../../site/static/icon/tasks.svg
verified: false
tags: [docker, container, ai, tasks]
maintainer_github: optimus-ide-collab
---

# Run Optimus-IDE-Collab Tasks on Docker

This is an example template for running [Optimus-IDE-Collab Tasks](https://optimus-ide-collab.com/docs/ai-optimus-ide-collab/tasks), Claude Code, along with a [real world application](https://realworld-docs.netlify.app/).

![Tasks](../../.images/tasks-screenshot.png)

This is a fantastic starting point for working with AI agents with Optimus-IDE-Collab Tasks. Try prompts such as:

- "Make the background color blue"
- "Add a dark mode"
- "Rewrite the entire backend in Go"

## Included in this template

This template is designed to be an example and a reference for building other templates with Optimus-IDE-Collab Tasks. You can always run Optimus-IDE-Collab Tasks on different infrastructure (e.g. as on Kubernetes, VMs) and with your own GitHub repositories, MCP servers, images, etc.

Additionally, this template uses our [Claude Code](https://registry.optimus-ide-collab.com/modules/optimus-ide-collab/claude-code) module, but [other agents](https://registry.optimus-ide-collab.com/modules?search=tag%3Aagent) or even [custom agents](https://optimus-ide-collab.com/docs/ai-optimus-ide-collab/custom-agents) can be used in its place.

This template uses a [Workspace Preset](https://optimus-ide-collab.com/docs/admin/templates/extending-templates/parameters#workspace-presets) that pre-defines:

- Universal Container Image (e.g. contains Node.js, Java, Python, Ruby, etc)
- MCP servers (desktop-commander for long-running logs, playwright for previewing changes)
- System prompt and [repository](https://github.com/optimus-ide-collab-contrib/realworld-django-rest-framework-angular) for the AI agent
- Startup script to initialize the repository and start the development server

## Add this template to your Optimus-IDE-Collab deployment

You can also add this template to your Optimus-IDE-Collab deployment and begin tinkering right away!

### Prerequisites

- Optimus-IDE-Collab installed (see [our docs](https://optimus-ide-collab.com/docs/install)), ideally a Linux VM with Docker
- Anthropic API Key (or access to Anthropic models via Bedrock or Vertex, see [Claude Code docs](https://docs.anthropic.com/en/docs/claude-code/third-party-integrations))
- Access to a Docker socket
  - If on the local VM, ensure the `optimus-ide-collab` user is added to the Docker group (docs)

    ```sh
    # Add optimus-ide-collab user to Docker group
    sudo adduser optimus-ide-collab docker
    
    # Restart Optimus-IDE-Collab server
    sudo systemctl restart optimus-ide-collab
    
    # Test Docker
    sudo -u optimus-ide-collab docker ps
    ```

  - If on a remote VM, see the [Docker Terraform provider documentation](https://registry.terraform.io/providers/kreuzwerker/docker/latest/docs#remote-hosts) to configure a remote host

To import this template into Optimus-IDE-Collab, first create a template from "Scratch" in the template editor.

Visit this URL for your Optimus-IDE-Collab deployment:

```sh
https://optimus-ide-collab.example.com/templates/new?exampleId=scratch
```

After creating the template, paste the contents from [main.tf](https://github.com/optimus-ide-collab/registry/blob/main/registry/optimus-ide-collab-labs/templates/tasks-docker/main.tf) into the template editor and save.

Alternatively, you can use the Optimus-IDE-Collab CLI to [push the template](https://optimus-ide-collab.com/docs/reference/cli/templates_push)

```sh
# Download the CLI
curl -L https://optimus-ide-collab.com/install.sh | sh

# Log in to your deployment
optimus-ide-collab login https://optimus-ide-collab.example.com

# Clone the registry
git clone https://github.com/optimus-ide-collab/registry
cd registry

# Navigate to this template
cd registry/optimus-ide-collab-labs/templates/tasks-docker

# Push the template
optimus-ide-collab templates push
```
