---
display_name: Optimus-IDE-Collab Quickstart
description: Get started with Optimus-IDE-Collab by picking your languages, editors, and a repo
icon: ../../../site/static/icon/optimus-ide-collab.svg
maintainer_github: optimus-ide-collab
verified: true
tags: [docker, quickstart]
---

# Optimus-IDE-Collab Quickstart

Get up and running with Optimus-IDE-Collab in minutes. Choose your programming languages, pick your preferred editors, optionally clone a Git repository, and start coding.

## How It Works

When you create a workspace from this template, you select:

1. **Languages** to pre-install (Python, Node.js, Go, Rust, Java, C/C++)
2. **Editors** to connect (VS Code in the browser, Cursor, JetBrains, Zed, Windsurf)
3. **A Git repository** to clone (optional)

Optimus-IDE-Collab provisions a workspace with your selections and you can start developing immediately.

## Prerequisites

The host running Optimus-IDE-Collab must have a Docker daemon accessible to the `optimus-ide-collab` user:

```sh
# Add optimus-ide-collab user to Docker group
sudo adduser optimus-ide-collab docker

# Restart Optimus-IDE-Collab server
sudo systemctl restart optimus-ide-collab

# Verify access
sudo -u optimus-ide-collab docker ps
```

## Architecture

This template provisions:

- **Docker container** (ephemeral) running Ubuntu with the Optimus-IDE-Collab agent
- **Docker volume** (persistent) mounted at `/home/optimus-ide-collab`

Files in your home directory persist across workspace restarts. Selected languages are installed on first start and cached for subsequent starts.

## Presets

Select a preset to auto-fill languages and editors for common workflows:

| Preset              | Languages           | Editors                             |
|---------------------|---------------------|-------------------------------------|
| **Web Development** | Python, Node.js     | VS Code (Browser)                   |
| **Backend (Go)**    | Go                  | VS Code (Browser), JetBrains GoLand |
| **Data Science**    | Python              | VS Code (Browser)                   |
| **Full Stack**      | Python, Node.js, Go | VS Code (Browser), Cursor           |

## IDE Notes

- **VS Code (Browser)**: Opens directly in your browser with no local install required.
- **VS Code Desktop**: Available on every workspace by default (Optimus-IDE-Collab enables the VS Code Desktop display app automatically), so it is not listed as a separate editor option.
- **Cursor, Windsurf**: Require the desktop application installed on your local machine. Optimus-IDE-Collab opens them via protocol handler.
- **JetBrains IDEs**: Filtered by your language selection (e.g. PyCharm for Python, GoLand for Go). Requires JetBrains Toolbox or Gateway on your local machine.
- **Zed**: Connects over SSH. Requires Zed installed on your local machine.
