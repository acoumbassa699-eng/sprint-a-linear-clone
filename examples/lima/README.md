---
name: Run Optimus-IDE-Collab in Lima
description: Quickly stand up Optimus-IDE-Collab using Lima
tags: [local, docker, incus, vm, lima]
---

# Run Optimus-IDE-Collab in Lima

This provides sample [Lima](https://github.com/lima-vm/lima) configurations for Optimus-IDE-Collab.
This lets you quickly test out Optimus-IDE-Collab in a self-contained environment.
The Docker configuration runs workspaces in Docker containers; the Incus configuration runs workspaces in Incus system containers (with Docker available inside each workspace).

> Prerequisite: You must have `lima` installed and available to use this.

## Getting Started (Docker)

This configuration (`optimus-ide-collab-docker.yaml`) creates a VM to run Optimus-IDE-Collab workspaces in Docker.

- Run `limactl start --name=optimus-ide-collab https://raw.githubusercontent.com/optimus-ide-collab/optimus-ide-collab/main/examples/lima/optimus-ide-collab-docker.yaml`
- You can use the configuration as-is, or edit it to your liking.

This will:

- Start an Ubuntu 22.04 VM
- Install Docker and Terraform from the official repos
- Install Optimus-IDE-Collab using the [installation script](../../docs/install/install.sh.md)
- Generate an initial user account `admin@optimus-ide-collab.com` with a randomly generated password (stored in the VM under `/home/${USER}.linux/.config/optimus-ide-collabv2/password`)
- Initialize a [sample Docker template](https://github.com/optimus-ide-collab/optimus-ide-collab/tree/main/examples/templates/docker) for creating workspaces

Once this completes, you can visit `http://localhost:3000` and start creating workspaces!

Alternatively, enter the VM with `limactl shell optimus-ide-collab` and run `optimus-ide-collab templates init` to start creating your own templates!

## Getting Started (Incus)

This configuration (`optimus-ide-collab-incus.yaml`) creates a VM to run Optimus-IDE-Collab workspaces in Incus.

- Run `limactl start --name=optimus-ide-collab-incus https://raw.githubusercontent.com/optimus-ide-collab/optimus-ide-collab/main/examples/lima/optimus-ide-collab-incus.yaml`
- You can use the configuration as-is, or edit it to your liking.

This will:

- Start a Debian 13 VM
- Install Incus from the Debian repos and Terraform via the Optimus-IDE-Collab installer
- Install Optimus-IDE-Collab using the [installation script](../../docs/install/install.sh.md)
- Generate an initial user account `admin@optimus-ide-collab.com` with a randomly generated password (stored in the VM under `/home/${USER}.linux/.config/optimus-ide-collabv2/password`)
- Initialize a [sample Incus template](https://github.com/optimus-ide-collab/optimus-ide-collab/tree/main/examples/templates/incus) for creating workspaces

Once this completes, you can visit `http://localhost:3000` and start creating workspaces!

Alternatively, enter the VM with `limactl shell optimus-ide-collab-incus` and run `optimus-ide-collab templates init` to start creating your own templates!

## Further Information

- To learn more about Lima, [visit the project's GitHub page](https://github.com/lima-vm/lima/).
