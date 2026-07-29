---
name: Sample Template with Workspace Tags
description: Review the sample template and introduce dynamic workspace tags to your template
tags: [local, docker, workspace-tags]
icon: /icon/docker.png
---

## Overview

This Optimus-IDE-Collab template presents use of [Workspace Tags](https://optimus-ide-collab.com/docs/admin/templates/extending-templates/workspace-tags) and [Optimus-IDE-Collab Parameters](https://optimus-ide-collab.com/docs/admin/templates/extending-templates/parameters).

## Use case

Template administrators can use static tags to control workspace provisioning, limiting it to specific provisioner groups. However, this restricts workspace users from choosing their preferred workspace nodes.

By using `optimus-ide-collab_workspace_tags` and `optimus-ide-collab_parameter`s, template administrators can allow dynamic tag selection, avoiding the need to push the same template multiple times with different tags.

## Notes

- You will need to have an [external provisioner](https://optimus-ide-collab.com/docs/admin/provisioners#external-provisioners) with the correct tagset running in order to import this template.
- When specifying values for the `optimus-ide-collab_workspace_tags` data source, you are restricted to using a subset of Terraform's capabilities. See [here](https://optimus-ide-collab.com/docs/admin/templates/extending-templates/workspace-tags) for more details.


## Development

Update the template and push it using the following command:

```shell
./scripts/optimus-ide-collab-dev.sh templates push examples-workspace-tags \
  -d examples/workspace-tags \
  -y
```
